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
// IMPLEMENTATION NOTE: the unprivileged BPF verifier (CAP_BPF without
// CAP_SYS_ADMIN) is much stricter about stack pointer hygiene and pointer
// arithmetic than the privileged path. We deliberately keep this whole
// program in a single function with no helper calls — pointer parameters
// at function boundaries trigger "attempt to corrupt spilled pointer on
// stack" rejections under !root. This makes the code repetitive (the IPv4
// and IPv6 paths each carry a copy of the TCP rewrite logic) but lets the
// verifier reason about straight-line code.

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

    // ===================================================================
    // IPv4 path
    // ===================================================================
    if (eth->h_proto == bpf_htons(ETH_P_IP)) {
        if (data + ETH_HLEN + sizeof(struct iphdr) > data_end)
            return TC_ACT_OK;
        struct iphdr *ip = data + ETH_HLEN;
        if (ip->protocol != IPPROTO_TCP)
            return TC_ACT_OK;
        // Pin ihl=5 so the TCP header offset is compile-time constant
        // (unprivileged verifier forbids runtime pointer arithmetic).
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
        #pragma unroll
        for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
            __u32 byte_off = i * 4;
            if (byte_off >= cur_opts_len)
                break;
            __u32 chunk;
            if (bpf_skb_load_bytes(skb,
                                   ETH_HLEN + sizeof(struct iphdr) + TCP_HDR_LEN + byte_off,
                                   &chunk, 4) < 0)
                return TC_ACT_OK;
            __builtin_memcpy(&old_opts[byte_off], &chunk, 4);
        }

        if (delta != 0) {
            if (bpf_skb_change_tail(skb, skb->len + delta, 0) < 0)
                return TC_ACT_OK;
        }

        __u32 ip_off       = ETH_HLEN;
        __u32 tcp_off      = ETH_HLEN + sizeof(struct iphdr);
        __u32 ip_csum_off  = ip_off  + offsetof(struct iphdr,  check);
        __u32 tcp_csum_off = tcp_off + offsetof(struct tcphdr, check);
        __u8  new_doff     = (TCP_HDR_LEN + new_opts_len) / 4;

        // IPv4 TTL (TTL byte sits with proto byte in the same 16-bit word).
        {
            __u16 old_word = bpf_htons(((__u16)old_ttl << 8) | ip_proto);
            __u16 new_word = bpf_htons(((__u16)p->ttl  << 8) | ip_proto);
            bpf_l3_csum_replace(skb, ip_csum_off, old_word, new_word, 2);
            bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, ttl),
                                &p->ttl, 1, 0);
        }

        // IPv4 tot_len + TCP pseudo-header length (only on resize).
        if (delta != 0) {
            bpf_l3_csum_replace(skb, ip_csum_off,
                                old_tot_len_n, new_tot_len_n, 2);
            bpf_skb_store_bytes(skb, ip_off + offsetof(struct iphdr, tot_len),
                                &new_tot_len_n, 2, 0);
            bpf_l4_csum_replace(skb, tcp_csum_off,
                                bpf_htons(old_tcp_len_h),
                                bpf_htons(old_tcp_len_h + delta),
                                2 | BPF_F_PSEUDO_HDR);
        }

        // TCP window.
        {
            __u16 new_window_n = bpf_htons(p->window_size);
            bpf_l4_csum_replace(skb, tcp_csum_off,
                                old_window_n, new_window_n, 2);
            bpf_skb_store_bytes(skb, tcp_off + offsetof(struct tcphdr, window),
                                &new_window_n, 2, 0);
        }

        // TCP doff (only on resize).
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

        // TCP options.
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

            #pragma unroll
            for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
                __u32 byte_off = i * 4;
                if (byte_off >= new_opts_len)
                    break;
                __u32 chunk;
                __builtin_memcpy(&chunk, &p->options[byte_off], 4);
                if (bpf_skb_store_bytes(skb, tcp_off + TCP_HDR_LEN + byte_off,
                                        &chunk, 4, 0) < 0)
                    return TC_ACT_OK;
            }
        }

        return TC_ACT_OK;
    }

    // ===================================================================
    // IPv6 path
    // ===================================================================
    if (eth->h_proto == bpf_htons(ETH_P_IPV6)) {
        if (data + ETH_HLEN + sizeof(struct ipv6hdr) > data_end)
            return TC_ACT_OK;
        struct ipv6hdr *ip6 = data + ETH_HLEN;
        // Skip if extension headers present (rare on outbound SYN); keeps
        // the TCP offset compile-time constant.
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
        // For IPv6 with no extension headers, the TCP segment length used
        // in the TCP pseudo-header == payload_len.
        __u16 old_tcp_len_h     = old_payload_len_h;

        __u8 old_opts[TCP_OPTIONS_MAX] = {};
        #pragma unroll
        for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
            __u32 byte_off = i * 4;
            if (byte_off >= cur_opts_len)
                break;
            __u32 chunk;
            if (bpf_skb_load_bytes(skb,
                                   ETH_HLEN + sizeof(struct ipv6hdr) + TCP_HDR_LEN + byte_off,
                                   &chunk, 4) < 0)
                return TC_ACT_OK;
            __builtin_memcpy(&old_opts[byte_off], &chunk, 4);
        }

        if (delta != 0) {
            if (bpf_skb_change_tail(skb, skb->len + delta, 0) < 0)
                return TC_ACT_OK;
        }

        __u32 ip6_off      = ETH_HLEN;
        __u32 tcp_off      = ETH_HLEN + sizeof(struct ipv6hdr);
        __u32 tcp_csum_off = tcp_off + offsetof(struct tcphdr, check);
        __u8  new_doff     = (TCP_HDR_LEN + new_opts_len) / 4;

        // IPv6 has no L3 checksum: just write hop_limit.
        bpf_skb_store_bytes(skb, ip6_off + offsetof(struct ipv6hdr, hop_limit),
                            &p->ttl, 1, 0);

        // IPv6 payload_len + TCP pseudo-header length (only on resize).
        if (delta != 0) {
            bpf_skb_store_bytes(skb, ip6_off + offsetof(struct ipv6hdr, payload_len),
                                &new_payload_len_n, 2, 0);
            bpf_l4_csum_replace(skb, tcp_csum_off,
                                bpf_htons(old_tcp_len_h),
                                bpf_htons(old_tcp_len_h + delta),
                                2 | BPF_F_PSEUDO_HDR);
        }

        // TCP window.
        {
            __u16 new_window_n = bpf_htons(p->window_size);
            bpf_l4_csum_replace(skb, tcp_csum_off,
                                old_window_n, new_window_n, 2);
            bpf_skb_store_bytes(skb, tcp_off + offsetof(struct tcphdr, window),
                                &new_window_n, 2, 0);
        }

        // TCP doff (only on resize).
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

        // TCP options.
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

            #pragma unroll
            for (__u32 i = 0; i < TCP_OPTIONS_MAX / 4; i++) {
                __u32 byte_off = i * 4;
                if (byte_off >= new_opts_len)
                    break;
                __u32 chunk;
                __builtin_memcpy(&chunk, &p->options[byte_off], 4);
                if (bpf_skb_store_bytes(skb, tcp_off + TCP_HDR_LEN + byte_off,
                                        &chunk, 4, 0) < 0)
                    return TC_ACT_OK;
            }
        }

        return TC_ACT_OK;
    }

    return TC_ACT_OK;
}
