package tcp

import (
	"fmt"

	"github.com/cilium/ebpf"
)

// Mark 0 is reserved by syn_rewrite.c to mean "no profile, pass through".
// Allocation starts at 1.
const reservedZeroMark uint32 = 0

// registerProfileLocked builds the BPF map entry for p, dedupes against
// previously registered entries by value, writes (or refreshes) the entry,
// and returns the fwmark. Caller must hold m.mu.
func (m *manager) registerProfileLocked(p Profile) (uint32, error) {
	entry, err := buildProfileEntry(p)
	if err != nil {
		return 0, err
	}

	mark, ok := m.byEntry[entry]
	if !ok {
		m.nextMark++
		if m.nextMark == reservedZeroMark {
			m.nextMark = 1
		}
		mark = m.nextMark
		if m.byEntry == nil {
			m.byEntry = make(map[synrewriteTcpProfile]uint32)
		}
		m.byEntry[entry] = mark
	}

	if err := m.objs.Profiles.Update(&mark, &entry, ebpf.UpdateAny); err != nil {
		return 0, fmt.Errorf("tcp: write profile to BPF map (mark %d): %w", mark, err)
	}
	return mark, nil
}

// buildProfileEntry converts the public Profile struct into the C-layout
// hash-map value expected by syn_rewrite.c. The structured option list is
// serialized (NOP-padded to a 4-byte boundary) here.
func buildProfileEntry(p Profile) (synrewriteTcpProfile, error) {
	var entry synrewriteTcpProfile
	options := buildOptions(p.Options)
	if len(options) > len(entry.Options) {
		return entry, fmt.Errorf("tcp: serialized options length %d exceeds max %d",
			len(options), len(entry.Options))
	}
	entry.WindowSize = p.WindowSize
	entry.Ttl = p.TTL
	entry.OptionsLen = uint8(len(options))
	copy(entry.Options[:], options)
	return entry, nil
}

// buildOptions serializes a list of TCP options.
//
// Each entry is encoded as:
//
//	NOP / End:     [Kind]                 (Data ignored, length 1)
//	others:        [Kind, 2+len(Data), Data...]
//
// The result is right-padded with NOPs to a 4-byte boundary so it satisfies
// the BPF program's alignment check on the options block.
func buildOptions(opts []Option) []byte {
	var out []byte
	for _, o := range opts {
		switch o.Kind {
		case OptKindNOP, OptKindEnd:
			out = append(out, o.Kind)
		default:
			out = append(out, o.Kind, byte(2+len(o.Data)))
			out = append(out, o.Data...)
		}
	}
	for len(out)%4 != 0 {
		out = append(out, OptKindNOP)
	}
	return out
}
