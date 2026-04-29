/* SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause */
/*
 * Minimal subset of libbpf's bpf_helpers.h.
 */

#ifndef __BPF_HELPERS_H__
#define __BPF_HELPERS_H__

#include "vmlinux.h"
#include "bpf_helper_defs.h"

#define SEC(name) __attribute__((section(name), used))

#define __uint(name, val)  int   (*name)[val]
#define __type(name, val)  typeof(val) *name
#define __array(name, val) typeof(val) *name[]

#define __always_inline inline __attribute__((always_inline))

#ifndef offsetof
#define offsetof(t, m) __builtin_offsetof(t, m)
#endif

/* TC actions (from <linux/pkt_cls.h>) */
#define TC_ACT_OK       0
#define TC_ACT_SHOT     2

/* IP protocols */
#define ETH_P_IP        0x0800
#define ETH_P_IPV6      0x86DD
#define IPPROTO_TCP     6

/* bpf_l4_csum_replace flags (from <linux/bpf.h>) */
#define BPF_F_RECOMPUTE_CSUM    (1ULL << 0)
#define BPF_F_PSEUDO_HDR        (1ULL << 4)
#define BPF_F_HDR_FIELD_MASK    0xfULL

#endif /* __BPF_HELPERS_H__ */
