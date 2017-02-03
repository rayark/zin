package middleware

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}

// Compressor compress the response body if the header of request
// contained `Accept-Encoding`
func Compressor(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		switch r.Header.Get("Accept-Encoding") {
		case "gzip":
			w.Header().Set("Content-Encoding", "gzip")
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			gzRespWriter := gzipResponseWriter{Writer: gzw, ResponseWriter: w}
			h(gzRespWriter, r, p)
		case "deflate":
			// Currently only support "gzip" encoding, "deflate" is not
			// implemented yet.
			h(w, r, p)
		default:
			h(w, r, p)
		}
	}
}
