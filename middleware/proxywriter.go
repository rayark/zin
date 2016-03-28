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
