/* SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause */
/*
 * Minimal vmlinux.h for the syn_rewrite TC eBPF program.
 *
 * A real BTF-derived vmlinux.h is ~5MB and is generated with:
 *   bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
 *
 * This hand-written subset contains only the kernel UAPI types the program
 * touches (packet headers + __sk_buff). All of these are stable kernel UAPI
 * with no field-offset drift, so we get cross-kernel compatibility without
 * relying on BTF/CO-RE relocations. If you need to access kernel-internal
 * structs later, replace this file with the full bpftool-generated one.
 */

#ifndef __VMLINUX_H__
#define __VMLINUX_H__

typedef unsigned char      __u8;
typedef unsigned short     __u16;
typedef unsigned int       __u32;
typedef unsigned long long __u64;
typedef signed char        __s8;
typedef signed short       __s16;
typedef signed int         __s32;
typedef signed long long   __s64;

typedef __u16 __be16;
typedef __u32 __be32;
typedef __u16 __sum16;
typedef __u32 __wsum;

/* Packet headers (Linux UAPI) */

struct ethhdr {
	unsigned char h_dest[6];
	unsigned char h_source[6];
	__be16        h_proto;
} __attribute__((packed));

struct iphdr {
#if defined(__LITTLE_ENDIAN_BITFIELD) || __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
	__u8    ihl: 4,
	        version: 4;
#else
	__u8    version: 4,
	        ihl: 4;
#endif
	__u8    tos;
	__be16  tot_len;
	__be16  id;
	__be16  frag_off;
	__u8    ttl;
	__u8    protocol;
	__sum16 check;
	__be32  saddr;
	__be32  daddr;
};

struct tcphdr {
	__be16 source;
	__be16 dest;
	__be32 seq;
	__be32 ack_seq;
#if defined(__LITTLE_ENDIAN_BITFIELD) || __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
	__u16  res1: 4,
	       doff: 4,
	       fin:  1,
	       syn:  1,
	       rst:  1,
	       psh:  1,
	       ack:  1,
	       urg:  1,
	       ece:  1,
	       cwr:  1;
#else
	__u16  doff: 4,
	       res1: 4,
	       cwr:  1,
	       ece:  1,
	       urg:  1,
	       ack:  1,
	       psh:  1,
	       rst:  1,
	       syn:  1,
	       fin:  1;
#endif
	__be16 window;
	__sum16 check;
	__be16 urg_ptr;
};

/* BPF program context for SCHED_CLS / SCHED_ACT (TC programs) */
struct __sk_buff {
	__u32 len;
	__u32 pkt_type;
	__u32 mark;
	__u32 queue_mapping;
	__u32 protocol;
	__u32 vlan_present;
	__u32 vlan_tci;
	__u32 vlan_proto;
	__u32 priority;
	__u32 ingress_ifindex;
	__u32 ifindex;
	__u32 tc_index;
	__u32 cb[5];
	__u32 hash;
	__u32 tc_classid;
	__u32 data;
	__u32 data_end;
	__u32 napi_id;
	__u32 family;
	__u32 remote_ip4;
	__u32 local_ip4;
	__u32 remote_ip6[4];
	__u32 local_ip6[4];
	__u32 remote_port;
	__u32 local_port;
	__u32 data_meta;
};

enum bpf_map_type {
	BPF_MAP_TYPE_UNSPEC = 0,
	BPF_MAP_TYPE_HASH   = 1,
	BPF_MAP_TYPE_ARRAY  = 2,
};

#endif /* __VMLINUX_H__ */
