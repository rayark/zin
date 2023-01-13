/* (C)2023 Rayark Inc. - All Rights Reserved
 * Rayark Confidential
 *
 * NOTICE: The intellectual and technical concepts contained herein are
 * proprietary to or under control of Rayark Inc. and its affiliates.
 * The information herein may be covered by patents, patents in process,
 * and are protected by trade secret or copyright law.
 * You may not disseminate this information or reproduce this material
 * unless otherwise prior agreed by Rayark Inc. in writing.
 */

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
