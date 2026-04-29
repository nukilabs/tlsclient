/* SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause */
/*
 * Minimal subset of libbpf's bpf_helper_defs.h. The real header is
 * auto-generated and uses a `(void *)<helper_id>` initializer, which
 * recent clangd configurations reject as a void*-to-function-pointer
 * conversion. We sidestep that by casting the integer helper ID directly
 * to the matching function pointer type via __u64 -- the clang BPF
 * backend still pattern-matches the constant-integer initializer and
 * emits a real BPF helper call instruction.
 */

#ifndef __BPF_HELPER_DEFS_H__
#define __BPF_HELPER_DEFS_H__

#include "vmlinux.h"

#define DEFINE_BPF_HELPER(NAME, ID, RET, ...)              \
	typedef RET NAME##_fn(__VA_ARGS__);                    \
	static NAME##_fn *NAME = (NAME##_fn *)(__u64)(ID)

/* BPF_FUNC_* helper IDs (from <linux/bpf.h>) */

DEFINE_BPF_HELPER(bpf_map_lookup_elem, 1, void *,
                  void *map, const void *key);

DEFINE_BPF_HELPER(bpf_skb_store_bytes, 9, long,
                  struct __sk_buff *skb, __u32 offset,
                  const void *from, __u32 len, __u64 flags);

DEFINE_BPF_HELPER(bpf_l3_csum_replace, 10, long,
                  struct __sk_buff *skb, __u32 offset,
                  __u64 from, __u64 to, __u64 size);

DEFINE_BPF_HELPER(bpf_l4_csum_replace, 11, long,
                  struct __sk_buff *skb, __u32 offset,
                  __u64 from, __u64 to, __u64 flags);

DEFINE_BPF_HELPER(bpf_skb_load_bytes, 26, long,
                  const struct __sk_buff *skb, __u32 offset,
                  void *to, __u32 len);

DEFINE_BPF_HELPER(bpf_csum_diff, 28, __s64,
                  __be32 *from, __u32 from_size,
                  __be32 *to, __u32 to_size, __wsum seed);

DEFINE_BPF_HELPER(bpf_skb_change_tail, 38, long,
                  struct __sk_buff *skb, __u32 len, __u64 flags);

#undef DEFINE_BPF_HELPER

#endif /* __BPF_HELPER_DEFS_H__ */
