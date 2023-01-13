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
	"bufio"
	"net"
	"net/http"
)

type ProxyWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewProxyWriter(w http.ResponseWriter) *ProxyWriter {
	return &ProxyWriter{
		ResponseWriter: w,
	}
}
func (w *ProxyWriter) Status() int {
	return w.status
}

func (w *ProxyWriter) Size() int {
	return w.size
}

func (w *ProxyWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}

func (w *ProxyWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *ProxyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
