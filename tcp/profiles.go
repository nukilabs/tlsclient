package tcp

// Profile describes the TCP-layer fingerprint of an outgoing SYN.
//
// Options is the structured option list. It is serialized (NOP-padded to a
// 4-byte boundary) at registration time and the resulting bytes replace
// the kernel's TCP options region in the SYN. The serialized length must
// be at most 40 bytes (max TCP header is 60; 20 fixed + 40 options).
//
// MSS and window-scale aren't separate fields because they're already
// encoded as Option entries (kinds 2 and 3). Only the on-wire SYN window
// and the IP TTL are kept as scalar fields.
//
// The Manager dedupes registrations by the serialized form, so passing
// equal profiles to multiple Clients reuses the same fwmark.
type Profile struct {
	WindowSize uint16 // pre-scaling, what goes in the SYN
	TTL        uint8
	Options    []Option
}

// Option is a single TCP option entry. Serialization at the wire layer
// emits one byte (Kind) for NOP / End-of-list, otherwise [Kind,
// 2+len(Data), Data...].
type Option struct {
	Kind uint8
	Data []byte
}

// TCP option kind numbers from the IANA TCP Parameters registry:
// https://www.iana.org/assignments/tcp-parameters/tcp-parameters.xhtml
//
// Active / standardized:
const (
	OptKindEnd         uint8 = 0  // RFC 9293 - End of Option List
	OptKindNOP         uint8 = 1  // RFC 9293 - No-Operation
	OptKindMSS         uint8 = 2  // RFC 9293 - Maximum Segment Size
	OptKindWScale      uint8 = 3  // RFC 7323 - Window Scale (WSOPT)
	OptKindSACKPerm    uint8 = 4  // RFC 2018 - SACK Permitted
	OptKindSACK        uint8 = 5  // RFC 2018 - SACK
	OptKindTimestamp   uint8 = 8  // RFC 7323 - Timestamps (TSOPT)
	OptKindQuickStart  uint8 = 27 // RFC 4782 - Quick-Start Response
	OptKindUserTimeout uint8 = 28 // RFC 5482 - User Timeout Option
	OptKindTCPAO       uint8 = 29 // RFC 5925 - TCP Authentication Option
	OptKindMPTCP       uint8 = 30 // RFC 8684 - Multipath TCP
	OptKindTFOCookie   uint8 = 34 // RFC 7413 - TCP Fast Open Cookie
	OptKindTCPENO      uint8 = 69 // RFC 8547 - Encryption Negotiation (TCP-ENO)
)

// Historical / obsolete (kept for fingerprint fidelity against legacy stacks
// that may still emit these in captures).
const (
	OptKindEcho          uint8 = 6   // RFC 1072 - obsoleted by Timestamp
	OptKindEchoReply     uint8 = 7   // RFC 1072 - obsoleted by Timestamp
	OptKindPOCPermitted  uint8 = 9   // RFC 1693 - Partial Order Conn. Permitted
	OptKindPOCProfile    uint8 = 10  // RFC 1693 - Partial Order Service Profile
	OptKindCC            uint8 = 11  // RFC 1644 - T/TCP CC
	OptKindCCNew         uint8 = 12  // RFC 1644 - T/TCP CC.NEW
	OptKindCCEcho        uint8 = 13  // RFC 1644 - T/TCP CC.ECHO
	OptKindAltChksumReq  uint8 = 14  // RFC 1146 - Alternate Checksum Request
	OptKindAltChksumData uint8 = 15  // RFC 1146 - Alternate Checksum Data
	OptKindSkeeter       uint8 = 16  // Knowles, Stev
	OptKindBubba         uint8 = 17  // Knowles, Stev
	OptKindTrailerChksum uint8 = 18  // Subbu & Monroe
	OptKindMD5Sig        uint8 = 19  // RFC 2385 - obsoleted by TCP-AO
	OptKindSCPSCap       uint8 = 20  // SCPS Capabilities
	OptKindSNACK         uint8 = 21  // Selective Negative Acknowledgements
	OptKindRecBoundaries uint8 = 22  // Record Boundaries
	OptKindCorruption    uint8 = 23  // Corruption experienced
	OptKindSNAP          uint8 = 24  // SNAP (Sub-Network Access Protocol)
	OptKindCompression   uint8 = 26  // TCP Compression Filter
	OptKindAccECN0       uint8 = 172 // RFC 9690 - Accurate ECN Order 0 (TEMP)
	OptKindAccECN1       uint8 = 173 // RFC 9690 - Accurate ECN Order 1 (TEMP)
)

// Experimental (RFC 4727 / RFC 6994 — use with magic kind sub-numbers).
const (
	OptKindExp1 uint8 = 253 // RFC 4727
	OptKindExp2 uint8 = 254 // RFC 4727
)

// Built-in OS-level fingerprint profiles. The TCP fingerprint is a
// property of the OS network stack, not the browser, so these are named
// by OS.
var (
	// Windows: TTL 128, window 65535, MSS=1460,
	// WS=8 in the options.
	Windows = Profile{
		WindowSize: 65535,
		TTL:        128,
		Options: []Option{
			{Kind: OptKindMSS, Data: []byte{0x05, 0xb4}}, // 1460
			{Kind: OptKindNOP},
			{Kind: OptKindWScale, Data: []byte{0x08}},
			{Kind: OptKindNOP},
			{Kind: OptKindNOP},
			{Kind: OptKindSACKPerm},
		},
	}

	// Apple: TTL 64, window 65535, MSS=1460,
	// WS=6 in the options, plus timestamps.
	Apple = Profile{
		WindowSize: 65535,
		TTL:        64,
		Options: []Option{
			{Kind: OptKindMSS, Data: []byte{0x05, 0xb4}},
			{Kind: OptKindNOP},
			{Kind: OptKindWScale, Data: []byte{0x06}},
			{Kind: OptKindNOP},
			{Kind: OptKindNOP},
			{Kind: OptKindTimestamp, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0}},
			{Kind: OptKindSACKPerm},
			{Kind: OptKindEnd},
		},
	}
)
