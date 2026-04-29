// SPDX-License-Identifier: Apache-2.0
//
// syn_rewrite.c -- TC egress eBPF program that rewrites outbound TCP SYN
// packets (IPv4 and IPv6) to match a configured TCP fingerprint profile.
//
// Profile selection is by skb->mark. The Go side sets SO_MARK on the
// socket; the kernel propagates that to skb->mark on egress; this program
// looks up the matching tcp_profile in `profiles` (a hash map keyed by
// mark). A zero or unknown mark passes through untouched.
//
// SYN packets carry no payload (TFO data carrying is detected and skipped),
// so options sit at the tail of the skb and bpf_skb_change_tail is the
// right helper for resize when the new options block is a different length
// than the kernel-emitted one.
//
// Rewritten fields:
//   IPv4: TTL (l3 csum), tot_len (l3 csum + TCP pseudo-hdr length l4 csum),
//         TCP window, TCP doff, TCP options
//   IPv6: hop_limit, payload_len, TCP pseudo-hdr length l4 csum, then
//         TCP window/doff/options identical to IPv4

#include "vmlinux.h"
#include "bpf_helpers.h"
#include "bpf_endian.h"

#define ETH_HLEN          14
#define IPV6_HDR_LEN      40
#define TCP_HDR_LEN       20
#define TCP_OPTIONS_MAX   40

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

// rewrite_tcp performs the TCP-layer rewrite shared by both IPv4 and IPv6:
// pseudo-header length update, window, doff, and options. tcp_off is the
// byte offset of the TCP header from the start of the skb. old_tcp_len_h
// is the host-order TCP segment length (= IPv4 tot_len - ip_hl, or IPv6
// payload_len when no extension headers).
static __always_inline int rewrite_tcp(struct __sk_buff *skb,
                                       struct tcp_profile *p,
                                       __u32 tcp_off,
                                       __u8  *old_opts,
                                       __u32 cur_opts_len,
                                       __u32 new_opts_len,
                                       int   delta,
                                       __u16 old_tcp_len_h,
                                       __u16 old_window_n) {
    __u32 tcp_csum_off = tcp_off + offsetof(struct tcphdr, check);
    __u8  new_doff     = (TCP_HDR_LEN + new_opts_len) / 4;

    // Pseudo-header segment length update (only on resize).
    if (delta != 0) {
        __u16 new_tcp_len_h = old_tcp_len_h + delta;
        bpf_l4_csum_replace(skb, tcp_csum_off,
                            bpf_htons(old_tcp_len_h),
                            bpf_htons(new_tcp_len_h),
                            2 | BPF_F_PSEUDO_HDR);
    }

    // TCP window
    {
        __u16 new_window_n = bpf_htons(p->window_size);
        bpf_l4_csum_replace(skb, tcp_csum_off,
                            old_window_n, new_window_n, 2);
        bpf_skb_store_bytes(skb, tcp_off + offsetof(struct tcphdr, window),
                            &new_window_n, 2, 0);
    }

    // TCP doff (only on resize). doff sits in the upper nibble of byte 12.
    // Read bytes 12-13 into a byte array (byte-order-independent), patch
    // byte 12 only, feed bpf_htons-formed network-order words to the csum
    // helper.
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

    // TCP options. Pad both old and new buffers to TCP_OPTIONS_MAX with
    // zeros for the csum diff (zero bytes contribute nothing). Then write
    // the new bytes in fixed 4-byte chunks.
    if (cur_opts_len > 0 || new_opts_len > 0) {
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

    return TC_ACT_OK;
}

// snapshot_old_options reads cur_opts_len bytes from the TCP options region
// into the (zero-padded) buffer old_opts. The verifier rejects variable-
// size bpf_skb_load_bytes calls, so we do this in fixed 4-byte chunks.
// SYN packets always emit 4-byte-aligned options.
static __always_inline int snapshot_old_options(struct __sk_buff *skb,
                                                __u32 opts_off,
                                                __u32 cur_opts_len,
                                                __u8 old_opts[TCP_OPTIONS_MAX]) {
    #pragma unroll
    for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
        __u32 byte_off = i * 4;
        if (byte_off >= cur_opts_len)
            break;
        __u32 chunk;
        if (bpf_skb_load_bytes(skb, opts_off + byte_off, &chunk, 4) < 0)
            return -1;
        __builtin_memcpy(&old_opts[byte_off], &chunk, 4);
    }
    return 0;
}

static __always_inline int rewrite_ipv4(struct __sk_buff *skb,
                                        struct tcp_profile *p,
                                        void *data, void *data_end) {
    if (data + ETH_HLEN + sizeof(struct iphdr) > data_end)
        return TC_ACT_OK;
    struct iphdr *ip = data + ETH_HLEN;
    if (ip->protocol != IPPROTO_TCP)
        return TC_ACT_OK;

    // Unprivileged BPF verifier forbids pointer arithmetic with runtime
    // values; pin ihl to 5 so the TCP offset is a compile-time constant.
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

    __u32 expected_len = ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN + cur_opts_len;
    if (skb->len != expected_len)
        return TC_ACT_OK;

    __u32 new_opts_len = p->options_len;
    int   delta        = (int)new_opts_len - (int)cur_opts_len;

    __u8  old_ttl       = ip->ttl;
    __u8  ip_proto      = ip->protocol;
    __u16 old_tot_len_n = ip->tot_len;
    __u16 old_window_n  = tcp->window;
    __u16 old_tot_len_h = bpf_ntohs(old_tot_len_n);
    __u16 old_tcp_len_h = old_tot_len_h - sizeof(struct iphdr);
    __u16 new_tot_len_n = bpf_htons(old_tot_len_h + delta);

    __u8 old_opts[TCP_OPTIONS_MAX] = {};
    if (snapshot_old_options(skb,
                             ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN,
                             cur_opts_len, old_opts) < 0)
        return TC_ACT_OK;

    if (delta != 0) {
        if (bpf_skb_change_tail(skb, skb->len + delta, 0) < 0)
            return TC_ACT_OK;
    }

    __u32 ip_off       = ETH_HLEN;
    __u32 tcp_off      = ETH_HLEN + sizeof(struct iphdr);
    __u32 ip_csum_off  = ip_off + offsetof(struct iphdr, check);

    // IPv4 TTL (TTL byte sits with proto byte in the same 16-bit word).
    {
        __u16 old_word = bpf_htons(((__u16)old_ttl << 8) | ip_proto);
        __u16 new_word = bpf_htons(((__u16)p->ttl  << 8) | ip_proto);
        bpf_l3_csum_replace(skb, ip_csum_off, old_word, new_word, 2);
        bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, ttl),
                            &p->ttl, 1, 0);
    }

    // IPv4 tot_len (only on resize).
    if (delta != 0) {
        bpf_l3_csum_replace(skb, ip_csum_off,
                            old_tot_len_n, new_tot_len_n, 2);
        bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, tot_len),
                            &new_tot_len_n, 2, 0);
    }

    return rewrite_tcp(skb, p, tcp_off, old_opts,
                       cur_opts_len, new_opts_len, delta,
                       old_tcp_len_h, old_window_n);
}

static __always_inline int rewrite_ipv6(struct __sk_buff *skb,
                                        struct tcp_profile *p,
                                        void *data, void *data_end) {
    if (data + ETH_HLEN + sizeof(struct ipv6hdr) > data_end)
        return TC_ACT_OK;
    struct ipv6hdr *ip6 = data + ETH_HLEN;
    // Skip if next header isn't TCP — extension headers would shift the
    // TCP offset and we keep all packet-pointer offsets compile-time
    // constant (unprivileged BPF requirement).
    if (ip6->nexthdr != IPPROTO_TCP)
        return TC_ACT_OK;

    if (data + ETH_HLEN + sizeof(struct ipv6hdr) + TCP_HDR_LEN > data_end)
        return TC_ACT_OK;

    struct tcphdr *tcp = data + ETH_HLEN + sizeof(struct ipv6hdr);
    if (!tcp->syn || tcp->ack)
        return TC_ACT_OK;
    if (tcp->doff < 5)
        return TC_ACT_OK;

    __u32 cur_opts_len = (__u32)tcp->doff * 4 - TCP_HDR_LEN;
    if (cur_opts_len > TCP_OPTIONS_MAX)
        return TC_ACT_OK;

    __u32 expected_len = ETH_HLEN + sizeof(struct ipv6hdr) + TCP_HDR_LEN + cur_opts_len;
    if (skb->len != expected_len)
        return TC_ACT_OK;

    __u32 new_opts_len = p->options_len;
    int   delta        = (int)new_opts_len - (int)cur_opts_len;

    __u16 old_payload_len_n = ip6->payload_len;
    __u16 old_window_n      = tcp->window;
    __u16 old_payload_len_h = bpf_ntohs(old_payload_len_n);
    __u16 new_payload_len_n = bpf_htons(old_payload_len_h + delta);
    // For IPv6 with no extension headers, the TCP segment length used in
    // the TCP pseudo-header == payload_len.
    __u16 old_tcp_len_h     = old_payload_len_h;

    __u8 old_opts[TCP_OPTIONS_MAX] = {};
    if (snapshot_old_options(skb,
                             ETH_HLEN + sizeof(struct ipv6hdr) + TCP_HDR_LEN,
                             cur_opts_len, old_opts) < 0)
        return TC_ACT_OK;

    if (delta != 0) {
        if (bpf_skb_change_tail(skb, skb->len + delta, 0) < 0)
            return TC_ACT_OK;
    }

    __u32 ip6_off = ETH_HLEN;
    __u32 tcp_off = ETH_HLEN + sizeof(struct ipv6hdr);

    // IPv6 has no L3 checksum: just write hop_limit.
    bpf_skb_store_bytes(skb, ip6_off + offsetof(struct ipv6hdr, hop_limit),
                        &p->ttl, 1, 0);

    // IPv6 payload_len (only on resize) — also no L3 checksum.
    if (delta != 0) {
        bpf_skb_store_bytes(skb, ip6_off + offsetof(struct ipv6hdr, payload_len),
                            &new_payload_len_n, 2, 0);
    }

    return rewrite_tcp(skb, p, tcp_off, old_opts,
                       cur_opts_len, new_opts_len, delta,
                       old_tcp_len_h, old_window_n);
}

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

    if (eth->h_proto == bpf_htons(ETH_P_IP))
        return rewrite_ipv4(skb, p, data, data_end);
    if (eth->h_proto == bpf_htons(ETH_P_IPV6))
        return rewrite_ipv6(skb, p, data, data_end);

    return TC_ACT_OK;
}
