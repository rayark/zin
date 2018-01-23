package middleware

import (
	"bytes"
	"net/http"
)

type DeferWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func NewDeferWriter(w http.ResponseWriter) *DeferWriter {
	return &DeferWriter{
		ResponseWriter: w,
		buf:            new(bytes.Buffer),
	}
}

func (w *DeferWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *DeferWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *DeferWriter) WriteAll() {
	w.ResponseWriter.Write(w.buf.Bytes())
}
