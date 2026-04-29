/* SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause */

#ifndef __BPF_ENDIAN_H__
#define __BPF_ENDIAN_H__

#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
# define bpf_htons(x) ((__be16)__builtin_bswap16((__u16)(x)))
# define bpf_ntohs(x) ((__u16)__builtin_bswap16((__be16)(x)))
# define bpf_htonl(x) ((__be32)__builtin_bswap32((__u32)(x)))
# define bpf_ntohl(x) ((__u32)__builtin_bswap32((__be32)(x)))
#elif __BYTE_ORDER__ == __ORDER_BIG_ENDIAN__
# define bpf_htons(x) ((__be16)(x))
# define bpf_ntohs(x) ((__u16)(x))
# define bpf_htonl(x) ((__be32)(x))
# define bpf_ntohl(x) ((__u32)(x))
#else
# error "Unknown __BYTE_ORDER__"
#endif

#endif /* __BPF_ENDIAN_H__ */
