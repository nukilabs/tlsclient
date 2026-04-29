// SPDX-License-Identifier: Apache-2.0
//
// syn_rewrite.c -- TC egress eBPF program that rewrites outbound IPv4 TCP
// SYN packets to match a configured TCP fingerprint profile.
//
// Profile selection is by skb->mark. The Go side calls SO_MARK on the
// socket; the kernel propagates that to skb->mark on egress; this program
// looks up the matching tcp_profile in `profiles` (a hash map keyed by
// mark). A zero or unknown mark passes through untouched.
//
// SYN packets carry no payload (TFO data carrying is detected and skipped),
// so options sit at the tail of the skb and bpf_skb_change_tail is the
// right helper for resize when the new options block is a different length
// than the kernel-emitted one.
//
// Rewritten fields, in order:
//   1. IP TTL                  (l3 csum incremental)
//   2. IP tot_len              (only on resize; l3 csum incremental,
//                               TCP pseudo-header length l4 csum incremental)
//   3. TCP window               (l4 csum incremental)
//   4. TCP doff                (only on resize; l4 csum incremental)
//   5. TCP options block       (l4 csum via bpf_csum_diff over old vs. new
//                               40-byte zero-padded buffers)

#include "vmlinux.h"
#include "bpf_helpers.h"
#include "bpf_endian.h"

#define ETH_HLEN         14
#define TCP_HDR_LEN      20
#define TCP_OPTIONS_MAX  40

// Mirrors the Go-side Profile. MSS and window scale aren't fields because
// they're already encoded in the options bytes (kind 2 and kind 3).
struct tcp_profile {
    __u16 window_size;  // host order; converted to network when written
    __u8  ttl;
    __u8  options_len;
    __u8  options[TCP_OPTIONS_MAX];
};

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, __u32);
    __type(value, struct tcp_profile);
    __uint(max_entries, 256);
} profiles SEC(".maps");

char _license[] SEC("license") = "Apache-2.0";

SEC("tc")
int syn_rewrite(struct __sk_buff *skb) {
    __u32 mark = skb->mark;
    if (mark == 0)
        return TC_ACT_OK;

    struct tcp_profile *p = bpf_map_lookup_elem(&profiles, &mark);
    if (!p)
        return TC_ACT_OK;
    if (p->options_len > TCP_OPTIONS_MAX || (p->options_len & 0x3))
        return TC_ACT_OK;

    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;

    if (data + ETH_HLEN > data_end)
        return TC_ACT_OK;
    struct ethhdr *eth = data;
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return TC_ACT_OK;

    if (data + ETH_HLEN + sizeof(struct iphdr) > data_end)
        return TC_ACT_OK;
    struct iphdr *ip = data + ETH_HLEN;
    if (ip->protocol != IPPROTO_TCP)
        return TC_ACT_OK;

    // Only handle IP headers without options (ihl == 5, 20 bytes). The
    // unprivileged BPF verifier forbids pointer arithmetic with runtime
    // values, so we keep all packet-pointer offsets compile-time constant.
    // Outbound SYNs from the kernel almost never carry IP options.
    if (ip->ihl != 5)
        return TC_ACT_OK;

    if (data + ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN > data_end)
        return TC_ACT_OK;

    struct tcphdr *tcp = data + ETH_HLEN + sizeof(struct iphdr);
    if (!tcp->syn || tcp->ack)
        return TC_ACT_OK;
    if (tcp->doff < 5)
        return TC_ACT_OK;

    __u32 cur_opts_len = (__u32)tcp->doff * 4 - TCP_HDR_LEN;
    if (cur_opts_len > TCP_OPTIONS_MAX)
        return TC_ACT_OK;

    // Reject SYN-with-payload (TCP Fast Open data, etc.). Our resize logic
    // assumes options sit at the very tail of the skb.
    __u32 expected_len = ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN + cur_opts_len;
    if (skb->len != expected_len)
        return TC_ACT_OK;

    __u32 new_opts_len = p->options_len;
    int   delta = (int)new_opts_len - (int)cur_opts_len;

    // Capture pre-rewrite values; pointers into skb data become invalid
    // after change_tail, so anything we still need we read now.
    __u8  old_ttl       = ip->ttl;
    __u8  ip_proto      = ip->protocol;
    __u16 old_tot_len_n = ip->tot_len;
    __u16 old_window_n  = tcp->window;

    __u16 old_tot_len_h = bpf_ntohs(old_tot_len_n);
    __u16 new_tot_len_h = old_tot_len_h + delta;
    __u16 new_tot_len_n = bpf_htons(new_tot_len_h);
    __u16 old_tcp_len_h = old_tot_len_h - sizeof(struct iphdr);
    __u16 new_tcp_len_h = old_tcp_len_h + delta;
    __u8  new_doff      = (TCP_HDR_LEN + new_opts_len) / 4;

    // Snapshot old options bytes (zero-padded to TCP_OPTIONS_MAX) for the
    // checksum diff. The verifier rejects bpf_skb_load_bytes with a size
    // argument whose range includes zero, so we read in fixed 4-byte
    // chunks under a bounded unrolled loop. SYN packets always have
    // 4-byte-aligned options (kernel pads with NOPs).
    __u8 old_opts[TCP_OPTIONS_MAX] = {};
    {
        __u32 opts_off = ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN;
        #pragma unroll
        for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
            __u32 byte_off = i * 4;
            if (byte_off >= cur_opts_len)
                break;
            __u32 chunk;
            if (bpf_skb_load_bytes(skb, opts_off + byte_off, &chunk, 4) < 0)
                return TC_ACT_OK;
            __builtin_memcpy(&old_opts[byte_off], &chunk, 4);
        }
    }

    if (delta != 0) {
        if (bpf_skb_change_tail(skb, skb->len + delta, 0) < 0)
            return TC_ACT_OK;
    }

    __u32 ip_off       = ETH_HLEN;
    __u32 tcp_off      = ETH_HLEN + sizeof(struct iphdr);
    __u32 ip_csum_off  = ip_off  + offsetof(struct iphdr,  check);
    __u32 tcp_csum_off = tcp_off + offsetof(struct tcphdr, check);

    // ---- IP TTL (TTL byte sits with proto byte in the same 16-bit word) ----
    {
        __u16 old_word = bpf_htons(((__u16)old_ttl << 8) | ip_proto);
        __u16 new_word = bpf_htons(((__u16)p->ttl  << 8) | ip_proto);
        bpf_l3_csum_replace(skb, ip_csum_off, old_word, new_word, 2);
        bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, ttl),
                            &p->ttl, 1, 0);
    }

    // ---- IP tot_len + TCP pseudo-header length (only on resize) ----
    if (delta != 0) {
        bpf_l3_csum_replace(skb, ip_csum_off,
                            old_tot_len_n, new_tot_len_n, 2);
        bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, tot_len),
                            &new_tot_len_n, 2, 0);
        // TCP pseudo-header carries the segment length (IP tot_len - 20).
        bpf_l4_csum_replace(skb, tcp_csum_off,
                            bpf_htons(old_tcp_len_h),
                            bpf_htons(new_tcp_len_h),
                            2 | BPF_F_PSEUDO_HDR);
    }

    // ---- TCP window ----
    {
        __u16 new_window_n = bpf_htons(p->window_size);
        bpf_l4_csum_replace(skb, tcp_csum_off,
                            old_window_n, new_window_n, 2);
        bpf_skb_store_bytes(skb, tcp_off + offsetof(struct tcphdr, window),
                            &new_window_n, 2, 0);
    }

    // ---- TCP doff (only on resize) ----
    // doff is the upper nibble of byte 12 (network order). Read bytes 12-13
    // into a byte array so the indexing is byte-order-independent (a __u16
    // load would put packet[12] in the low byte on LE and the high byte on
    // BE), then patch byte 12 only and feed bpf_htons-formed network-order
    // words to the checksum helper.
    if (delta != 0) {
        __u8 buf[2] = {0};
        if (bpf_skb_load_bytes(skb, tcp_off + 12, buf, 2) < 0)
            return TC_ACT_OK;
        __u8 doffres_old = buf[0];
        __u8 flags_byte  = buf[1];
        __u8 doffres_new = (new_doff << 4) | (doffres_old & 0x0f);
        __be16 from_be = bpf_htons(((__u16)doffres_old << 8) | flags_byte);
        __be16 to_be   = bpf_htons(((__u16)doffres_new << 8) | flags_byte);
        bpf_l4_csum_replace(skb, tcp_csum_off, from_be, to_be, 2);
        bpf_skb_store_bytes(skb, tcp_off + 12, &doffres_new, 1, 0);
    }

    // ---- TCP options ----
    if (cur_opts_len > 0 || new_opts_len > 0) {
        // Pad both buffers to TCP_OPTIONS_MAX with zeros. Zero bytes don't
        // contribute to the checksum, so the resulting diff is correct
        // even though we're checksum-diffing 40 bytes when the live region
        // may be smaller.
        __u8 new_opts[TCP_OPTIONS_MAX] = {};
        #pragma unroll
        for (int i = 0; i < TCP_OPTIONS_MAX; i++) {
            new_opts[i] = (i < (int)new_opts_len) ? p->options[i] : 0;
        }

        __s64 csum_delta = bpf_csum_diff((__be32 *)old_opts, TCP_OPTIONS_MAX,
                                         (__be32 *)new_opts, TCP_OPTIONS_MAX,
                                         0);
        if (csum_delta < 0)
            return TC_ACT_OK;
        bpf_l4_csum_replace(skb, tcp_csum_off, 0, (__u32)csum_delta, 0);

        // Same bounded-chunk pattern for the write side: variable-size
        // bpf_skb_store_bytes also runs into "could be zero" verifier
        // pushback on recent kernels.
        {
            __u32 opts_off = tcp_off + TCP_HDR_LEN;
            #pragma unroll
            for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
                __u32 byte_off = i * 4;
                if (byte_off >= new_opts_len)
                    break;
                __u32 chunk;
                __builtin_memcpy(&chunk, &p->options[byte_off], 4);
                if (bpf_skb_store_bytes(skb, opts_off + byte_off,
                                        &chunk, 4, 0) < 0)
                    return TC_ACT_OK;
            }
        }
    }

    return TC_ACT_OK;
}
