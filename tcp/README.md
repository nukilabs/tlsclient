# tcp — TC eBPF SYN fingerprint spoofing

Linux-only. Rewrites outbound TCP SYN packets after the kernel emits them but
before they leave the NIC, so the on-wire fingerprint (TTL, window, options)
matches a configured target OS profile while the kernel's TCP state machine
keeps doing all the real work.

## Usage

```go
import (
    "github.com/nukilabs/tlsclient"
    "github.com/nukilabs/tlsclient/profiles"
    "github.com/nukilabs/tlsclient/tcp"
)

client := tlsclient.New(profiles.Chrome(120),
    tlsclient.WithTCP(tcp.Windows),
)
res, err := client.Get("https://example.com/")
```

The first call to a `WithTCP`-enabled Client lazily attaches the BPF program to
the default-route interface (read from `/proc/net/route`). Override before
the first dial:

```go
tcp.SetInterface("ens3")          // programmatic
// or
TCP_IFACE=ens3 ./yourbinary    // env var
```

Built-in profiles: `tcp.Windows`, `tcp.MacOS`, `tcp.IOS`. They're educated
defaults — replace with pcap-derived values for your reference target.

For full control, construct your own `Profile`:

```go
custom := tcp.Profile{
    WindowSize: 65535,
    TTL:        128,
    Options: []tcp.Option{
        {Kind: tcp.OptKindMSS, Data: []byte{0x05, 0xb4}}, // 1460
        {Kind: tcp.OptKindNOP},
        {Kind: tcp.OptKindWScale, Data: []byte{0x08}},
        {Kind: tcp.OptKindNOP},
        {Kind: tcp.OptKindNOP},
        {Kind: tcp.OptKindSACKPerm},
    },
}
client := tlsclient.New(httpprofiles.Chrome(120), tlsclient.WithTCP(custom))
```

## What gets rewritten

For SYN packets (SYN=1, ACK=0, no payload) on the egress hook of the attached
interface, when `skb->mark` matches a registered profile:

- IP TTL
- IP total length (only when the new options block is a different length)
- TCP window
- TCP data offset (only on length change)
- TCP options block (replaced wholesale with `profile.Options`)
- IP and TCP checksums (incrementally, via `bpf_l3_csum_replace` /
  `bpf_l4_csum_replace`; options use `bpf_csum_diff` over zero-padded buffers)

A zero or unknown mark passes through untouched. SYN-with-payload (TCP Fast
Open data) is detected and skipped — the resize logic assumes options are at
the tail of the skb.

## Privileges

Loading the BPF program and attaching to a TC egress hook requires
**CAP_NET_ADMIN + CAP_BPF** (or just running as root). On modern kernels
RLIMIT_MEMLOCK accounting is replaced by per-cgroup memcg, but
`rlimit.RemoveMemlock` is still called for older kernel compatibility.

`SO_MARK` (set per-connection by the dialer) requires `CAP_NET_ADMIN`.

## Kernel requirements

Uses **TCX** (`link.AttachTCX`), so **Linux ≥ 6.6**. Older kernels return
`ENOSYS`/`EOPNOTSUPP` at attach time. A clsact-qdisc-based fallback via raw
netlink is doable in ~150 lines if needed — open an issue.

## Docker

```
docker run --cap-add=NET_ADMIN --cap-add=BPF ...
```

(Or `--privileged`.) The host kernel must be ≥ 6.6.

## Building the BPF object

Generated outputs (`synrewrite_*_bpfel.go` and `.o`) are committed so users
without clang installed can `go build`. Regenerate after editing
`bpf/syn_rewrite.c`:

```bash
sudo apt install clang llvm
go generate ./tcp/...
```

## Verifying the rewrite

### tcpdump

```bash
sudo tcpdump -nn -vvv -i eth0 \
  'tcp[tcpflags] & tcp-syn != 0 and not tcp[tcpflags] & tcp-ack != 0'
```

The SYN should show the configured TTL, window, and options:

```
... ttl 128, ... Flags [S], ..., win 65535, options [mss 1460,nop,wscale 8,nop,nop,sackOK]
```

### Live fingerprint readout

`https://robinsamuel.dev/` returns a JSON describing what TCP fingerprint the
server saw (`device.os`, `tcp.window`, `tcp.window_scale`, `tcp.options`,
`tcp.maximum_segment_size`). Useful as a black-box check.

### JA4T

Run a JA4T-aware capture (Zeek + JA4 plugin, `ja4t.py`) on the SYN and compare
against the reference for the target OS.

## The kernel-RST gotcha

Some SYN-rewriting implementations use a TUN device or raw sockets to send
crafted SYNs. In those setups, when the SYN-ACK comes back, the kernel has no
matching socket — it sends a RST. **This implementation doesn't have that
problem**: the kernel emits the SYN through its own socket (so the connection
is in `SYN_SENT` state); we only mutate egress bytes. The kernel still owns
the socket and its 4-tuple, so SYN-ACKs match.

If you do see RSTs, look elsewhere: bad checksum (rewrite bug), upstream
firewall, or the target server actively rejecting the SYN fingerprint.

## Why Linux only

This is a TC (Traffic Control) eBPF program attached to an egress hook on a
Linux netdev. macOS uses Network Extensions (PF/divert sockets); Windows
uses WFP. Both would require entirely separate implementations and aren't
planned. On non-Linux the public API compiles but `WithTCP` is a silent
no-op: dials proceed with the host kernel's normal TCP fingerprint.

## Tests

```bash
sudo -E go test -v ./tcp
```

`TestSynRewriteOnLoopback` attaches the program to `lo`, captures SYNs via
AF_PACKET, and asserts the rewritten bytes match the profile. Skipped when
not root.

## Not in scope

- Userspace TCP stack — we just rewrite bytes on egress.
- Receive-path rewrites — SYN-only.
- macOS / Windows — stubbed to return an error.
- Spoofing the proxy-to-target hop when an outbound HTTP/SOCKS proxy is
  configured. The local-to-proxy SYN does get marked and rewritten;
  beyond the proxy, no.
