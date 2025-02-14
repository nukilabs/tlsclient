package tlsclient

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
	"github.com/nukilabs/http"
)

const (
	zlibMethodDeflate = 0x78
	zlibLevelDefault  = 0x9C
	zlibLevelLow      = 0x01
	zlibLevelMedium   = 0x5E
	zlibLevelBest     = 0xDA
)

func DecompressBody(res *http.Response) {
	ce := res.Header.Get("Content-Encoding")
	switch ce {
	case "gzip":
		res.Body = &gzipReader{body: res.Body}
	case "br":
		res.Body = &brReader{body: res.Body}
	case "deflate":
		// read zlib header
		var header [2]byte
		if _, err := io.ReadFull(res.Body, header[:]); err != nil {
			return
		}
		// reset body to include header
		res.Body = io.NopCloser(io.MultiReader(bytes.NewReader(header[:]), res.Body))
		// check for zlib header
		if header[0] == zlibMethodDeflate && (header[1] == zlibLevelDefault || header[1] == zlibLevelLow || header[1] == zlibLevelMedium || header[1] == zlibLevelBest) {
			res.Body = &zlibDeflateReader{body: res.Body}
		} else if header[0] == zlibMethodDeflate {
			res.Body = &deflateReader{body: res.Body}
		}
		return
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

// zlibDeflateReader wraps a response body so it can lazily
// call zlib.NewReader on the first call to Read
type zlibDeflateReader struct {
	body io.ReadCloser
	r    io.ReadCloser
	err  error
}

func (z *zlibDeflateReader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	if z.r == nil {
		z.r, err = zlib.NewReader(z.body)
		if err != nil {
			z.err = err
			return 0, z.err
		}
	}
	return z.r.Read(p)
}

func (z *zlibDeflateReader) Close() error {
	return z.r.Close()
}

// deflateReader wraps a response body so it can lazily
// call flate.NewReader on the first call to Read
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
		dr.r = flate.NewReader(dr.body)
	}
	return dr.r.Read(p)
}

func (dr *deflateReader) Close() error {
	return dr.r.Close()
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
	return z.body.Close()
}
