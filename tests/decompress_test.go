package tests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"hash/adler32"
	"io"
	"strconv"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/nukilabs/http"
	"github.com/nukilabs/tlsclient"
)

var plain = bytes.Repeat([]byte("hello, deflate! "), 64)

type trackingBody struct {
	io.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

func newResponse(encoding string, body []byte) (*http.Response, *trackingBody) {
	tb := &trackingBody{Reader: bytes.NewReader(body)}
	res := &http.Response{
		Header:        http.Header{},
		Body:          tb,
		ContentLength: int64(len(body)),
	}
	res.Header.Set("Content-Encoding", encoding)
	res.Header.Set("Content-Length", strconv.Itoa(len(body)))
	return res, tb
}

func rawDeflate(t *testing.T, data []byte) []byte {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func zlibLevel(t *testing.T, data []byte, level int) []byte {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, level)
	if err != nil {
		t.Fatal(err)
	}
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

// zlibSmallWindow builds a zlib stream advertising a 4K window (CMF 0x48),
// which compress/zlib cannot produce but decodes fine.
func zlibSmallWindow(t *testing.T, data []byte) []byte {
	out := []byte{0x48, 0x89}
	out = append(out, rawDeflate(t, data)...)
	return binary.BigEndian.AppendUint32(out, adler32.Checksum(data))
}

func TestDecompressDeflate(t *testing.T) {
	cases := map[string][]byte{
		"raw":         rawDeflate(t, plain),
		"zlib":        zlibLevel(t, plain, zlib.DefaultCompression),
		"zlib-fast":   zlibLevel(t, plain, zlib.BestSpeed),
		"zlib-best":   zlibLevel(t, plain, zlib.BestCompression),
		"zlib-none":   zlibLevel(t, plain, zlib.NoCompression),
		"zlib-window": zlibSmallWindow(t, plain),
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			res, tb := newResponse("deflate", body)
			tlsclient.DecompressBody(res)
			got, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, plain) {
				t.Fatalf("body not decoded: got %d bytes, want %d", len(got), len(plain))
			}
			if err := res.Body.Close(); err != nil {
				t.Fatal(err)
			}
			if !tb.closed {
				t.Error("underlying body not closed")
			}
			if res.Header.Get("Content-Encoding") != "" || res.Header.Get("Content-Length") != "" {
				t.Errorf("stale headers: %v", res.Header)
			}
			if !res.Uncompressed || res.ContentLength != -1 {
				t.Errorf("Uncompressed=%v ContentLength=%d", res.Uncompressed, res.ContentLength)
			}
		})
	}
}

func TestDecompressDeflateCorrupt(t *testing.T) {
	body := zlibLevel(t, plain, zlib.DefaultCompression)
	body[len(body)-1] ^= 0xff // break the adler32 checksum
	res, _ := newResponse("deflate", body)
	tlsclient.DecompressBody(res)
	if _, err := io.ReadAll(res.Body); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestDecompressDeflateEmpty(t *testing.T) {
	res, tb := newResponse("deflate", nil)
	tlsclient.DecompressBody(res)
	got, err := io.ReadAll(res.Body)
	if err != nil || len(got) != 0 {
		t.Fatalf("got %q, %v", got, err)
	}
	res.Body.Close()
	if !tb.closed {
		t.Error("underlying body not closed")
	}
}

func TestDecompressCloseBeforeRead(t *testing.T) {
	bodies := map[string][]byte{
		"deflate": rawDeflate(t, plain),
		"gzip":    nil,
		"br":      nil,
		"zstd":    nil,
	}
	for encoding, body := range bodies {
		t.Run(encoding, func(t *testing.T) {
			res, tb := newResponse(encoding, body)
			tlsclient.DecompressBody(res)
			if err := res.Body.Close(); err != nil {
				t.Fatal(err)
			}
			if !tb.closed {
				t.Error("underlying body not closed")
			}
		})
	}
}

func TestDecompressOther(t *testing.T) {
	var gz, br bytes.Buffer
	gw := gzip.NewWriter(&gz)
	gw.Write(plain)
	gw.Close()
	bw := brotli.NewWriter(&br)
	bw.Write(plain)
	bw.Close()
	zw, err := zstd.NewWriter(nil)
	if err != nil {
		t.Fatal(err)
	}
	zs := zw.EncodeAll(plain, nil)

	cases := map[string][]byte{
		"gzip": gz.Bytes(),
		"br":   br.Bytes(),
		"zstd": zs,
		"GZIP": gz.Bytes(),
	}
	for encoding, body := range cases {
		t.Run(encoding, func(t *testing.T) {
			res, tb := newResponse(encoding, body)
			tlsclient.DecompressBody(res)
			got, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, plain) {
				t.Fatal("body not decoded")
			}
			res.Body.Close()
			if !tb.closed {
				t.Error("underlying body not closed")
			}
		})
	}
}
