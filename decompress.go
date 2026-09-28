package tlsclient

import (
	"bufio"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
	"github.com/nukilabs/http"
)

func DecompressBody(res *http.Response) {
	ce := strings.ToLower(strings.TrimSpace(res.Header.Get("Content-Encoding")))
	switch ce {
	case "gzip":
		res.Body = &gzipReader{body: res.Body}
	case "br":
		res.Body = &brReader{body: res.Body}
	case "deflate":
		res.Body = &deflateReader{body: res.Body}
	case "zstd":
		res.Body = &zstdReader{body: res.Body}
	default:
		return
	}
	res.Header.Del("Content-Encoding")
	res.Header.Del("Content-Length")
	res.Uncompressed = true
	res.ContentLength = -1
}

// gzipReader wraps a response body so it can lazily
// call gzip.NewReader on the first call to Read
type gzipReader struct {
	body io.ReadCloser
	r    *gzip.Reader
	err  error
}

func (gz *gzipReader) Read(p []byte) (n int, err error) {
	if gz.err != nil {
		return 0, gz.err
	}
	if gz.r == nil {
		gz.r, err = gzip.NewReader(gz.body)
		if err != nil {
			gz.err = err
			return 0, err
		}
	}
	return gz.r.Read(p)
}

func (gz *gzipReader) Close() error {
	return gz.body.Close()
}

// brReader wraps a response body so it can lazily
// call brotli.NewReader on the first call to Read
type brReader struct {
	body io.ReadCloser
	r    *brotli.Reader
	err  error
}

func (br *brReader) Read(p []byte) (n int, err error) {
	if br.err != nil {
		return 0, br.err
	}
	if br.r == nil {
		br.r = brotli.NewReader(br.body)
	}
	return br.r.Read(p)
}

func (br *brReader) Close() error {
	return br.body.Close()
}

// deflateReader wraps a response body so it can lazily decide on the
// first call to Read whether the body is zlib-wrapped (RFC 1950) or raw
// deflate (RFC 1951). Servers send both for "Content-Encoding: deflate".
type deflateReader struct {
	body io.ReadCloser
	r    io.ReadCloser
	err  error
}

func (dr *deflateReader) Read(p []byte) (n int, err error) {
	if dr.err != nil {
		return 0, dr.err
	}
	if dr.r == nil {
		br := bufio.NewReader(dr.body)
		header, err := br.Peek(2)
		if len(header) == 0 {
			// empty body
			dr.err = err
			return 0, err
		}
		if len(header) == 2 && isZlibHeader(header[0], header[1]) {
			dr.r, err = zlib.NewReader(br)
			if err != nil {
				dr.err = err
				return 0, err
			}
		} else {
			dr.r = flate.NewReader(br)
		}
	}
	return dr.r.Read(p)
}

func (dr *deflateReader) Close() error {
	if dr.r != nil {
		dr.r.Close()
	}
	return dr.body.Close()
}

// isZlibHeader reports whether cmf and flg form a valid zlib header:
// compression method 8 (deflate), window size of at most 32K and a
// correct FCHECK.
func isZlibHeader(cmf, flg byte) bool {
	return cmf&0x0f == 8 && cmf>>4 <= 7 && (uint16(cmf)<<8|uint16(flg))%31 == 0
}

// zstdReader wraps a response body so it can lazily
// call zstd.NewReader on the first call to Read
type zstdReader struct {
	body io.ReadCloser
	r    *zstd.Decoder
	err  error
}

func (z *zstdReader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	if z.r == nil {
		z.r, err = zstd.NewReader(z.body)
		if err != nil {
			z.err = err
			return 0, z.err
		}
	}
	return z.r.Read(p)
}

func (z *zstdReader) Close() error {
	if z.r != nil {
		z.r.Close()
	}
	return z.body.Close()
}
