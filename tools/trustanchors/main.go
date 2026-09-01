// Command trustanchors regenerates profiles/trustanchors.go from Chromium's root store.
//
// The trust anchor identifiers Chrome advertises in the trust_anchors extension come from the
// Chrome Root Store, which ships through the component updater and therefore moves independently
// of browser releases: two installs on the same Chrome version can advertise different sets, and
// the same install changes set without a browser update. So the table has to be refreshed on the
// root store's schedule rather than on Chrome's.
//
// Refresh the table:
//
//	go run ./tools/trustanchors
//
// Report drift without writing, which is what CI does:
//
//	go run ./tools/trustanchors -check
//
// Note that root_store.textproto on Chromium main is the source that later becomes a component,
// so it leads what real installs advertise by days to weeks. Prefer landing a change here once
// the component has actually rolled out.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const sourcepath = "net/data/ssl/chrome_root_store/root_store.textproto"

var (
	ref   = flag.String("ref", "main", "Chromium ref to read the root store from")
	out   = flag.String("out", "profiles/trustanchors.go", "file to write")
	check = flag.Bool("check", false, "report drift and exit non-zero instead of writing")
)

var (
	anchorre  = regexp.MustCompile(`trust_anchor_id:\s*"([^"]*)"(?:[^\S\n]*#[^\S\n]*([0-9.]+))?`)
	versionre = regexp.MustCompile(`(?m)^version_major:\s*(\d+)`)
)

type anchor struct {
	id  []byte
	oid string
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "trustanchors:", err)
		os.Exit(1)
	}
}

func run() error {
	src, err := fetch(*ref)
	if err != nil {
		return err
	}

	version := "unknown"
	if m := versionre.FindSubmatch(src); m != nil {
		version = string(m[1])
	}

	anchors, err := parse(src)
	if err != nil {
		return err
	}
	// An upstream format change that silently matched nothing would empty the table and take
	// the extension with it, so treat that as a failure rather than a valid result.
	if len(anchors) == 0 {
		return errors.New("no trust anchor identifiers found, the upstream format has probably changed")
	}

	want, err := render(anchors, version)
	if err != nil {
		return err
	}

	got, err := os.ReadFile(*out)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if bytes.Equal(got, want) {
		fmt.Printf("up to date: %d anchors, root store version_major %s\n", len(anchors), version)
		return nil
	}

	if *check {
		return fmt.Errorf("%s is out of date: upstream has %d anchors at version_major %s, run go run ./tools/trustanchors",
			*out, len(anchors), version)
	}
	if err := os.WriteFile(*out, want, 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s: %d anchors, root store version_major %s\n", *out, len(anchors), version)
	return nil
}

func fetch(ref string) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/chromium/chromium/%s/%s", ref, sourcepath)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, res.Status)
	}
	return io.ReadAll(res.Body)
}

func parse(src []byte) ([]anchor, error) {
	var anchors []anchor
	for _, m := range anchorre.FindAllSubmatch(src, -1) {
		id, err := unescape(string(m[1]))
		if err != nil {
			return nil, fmt.Errorf("trust_anchor_id %q: %w", m[1], err)
		}
		if len(id) == 0 || len(id) > 255 {
			return nil, fmt.Errorf("trust_anchor_id %q has an unusable length %d", m[1], len(id))
		}
		anchors = append(anchors, anchor{id: id, oid: string(m[2])})
	}
	// Sorted for stable diffs.
	slices.SortFunc(anchors, func(a, b anchor) int { return bytes.Compare(a.id, b.id) })
	anchors = slices.CompactFunc(anchors, func(a, b anchor) bool { return bytes.Equal(a.id, b.id) })
	return anchors, nil
}

// unescape decodes a protobuf text-format string literal.
func unescape(s string) ([]byte, error) {
	var out []byte
	for i := 0; i < len(s); {
		if s[i] != '\\' {
			out = append(out, s[i])
			i++
			continue
		}
		if i+1 >= len(s) {
			return nil, errors.New("trailing backslash")
		}
		switch c := s[i+1]; c {
		case 'x', 'X':
			if i+4 > len(s) {
				return nil, errors.New("truncated hex escape")
			}
			v, err := strconv.ParseUint(s[i+2:i+4], 16, 8)
			if err != nil {
				return nil, fmt.Errorf("bad hex escape: %w", err)
			}
			out = append(out, byte(v))
			i += 4
		case 'n':
			out = append(out, '\n')
			i += 2
		case 'r':
			out = append(out, '\r')
			i += 2
		case 't':
			out = append(out, '\t')
			i += 2
		case '\\', '"', '\'':
			out = append(out, c)
			i += 2
		default:
			if c < '0' || c > '7' || i+4 > len(s) {
				return nil, fmt.Errorf("unsupported escape \\%c", c)
			}
			v, err := strconv.ParseUint(s[i+1:i+4], 8, 8)
			if err != nil {
				return nil, fmt.Errorf("bad octal escape: %w", err)
			}
			out = append(out, byte(v))
			i += 4
		}
	}
	return out, nil
}

func render(anchors []anchor, version string) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by tools/trustanchors. DO NOT EDIT.\n")
	fmt.Fprintf(&buf, "// Source: chromium %s, %s (version_major %s)\n\n", *ref, sourcepath, version)
	buf.WriteString("package profiles\n\n")
	buf.WriteString(`// chromerootstoreanchors are the Chrome Root Store trust anchors that carry a trust anchor
// identifier, which is the set Chrome advertises in the trust_anchors extension
// (draft-ietf-tls-trust-anchor-ids). Each entry is a relative object identifier in DER base-128
// form. The leading arcs are the CA's IANA private enterprise number: 11129 is Google Trust
// Services and 44947 is ISRG / Let's Encrypt.
//
// The set is shared by every Chrome profile, because it follows the root store rather than the
// browser version.
`)
	buf.WriteString("var chromerootstoreanchors = [][]byte{\n")
	for _, a := range anchors {
		parts := make([]string, len(a.id))
		for i, b := range a.id {
			parts[i] = fmt.Sprintf("0x%02x", b)
		}
		fmt.Fprintf(&buf, "\t{%s},", strings.Join(parts, ", "))
		if a.oid != "" {
			fmt.Fprintf(&buf, " // %s", a.oid)
		}
		buf.WriteString("\n")
	}
	buf.WriteString("}\n")

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("formatting generated source: %w", err)
	}
	return src, nil
}
