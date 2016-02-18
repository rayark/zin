package middleware

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func Logger(h httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var buf bytes.Buffer
		proxyWriter := NewProxyWriter(w)
		printStart(&buf, r)
		t1 := time.Now()
		h(proxyWriter, r, p)
		t2 := time.Now()
		printEnd(&buf, proxyWriter, t2.Sub(t1))
		log.Print(buf.String())
	}
}

func printStart(buf *bytes.Buffer, r *http.Request) {
	buf.WriteString("Started ")
	cW(buf, bMagenta, "%s ", r.Method)
	cW(buf, nBlue, "%q ", r.URL.String())
	buf.WriteString("from ")
	buf.WriteString(r.RemoteAddr)
}

func printEnd(buf *bytes.Buffer, w *ProxyWriter, dt time.Duration) {

	buf.WriteString("Returning ")
	status := w.Status()
	if status < 200 {
		cW(buf, bBlue, "%03d", status)
	} else if status < 300 {
		cW(buf, bGreen, "%03d", status)
	} else if status < 400 {
		cW(buf, bCyan, "%03d", status)
	} else if status < 500 {
		cW(buf, bYellow, "%03d", status)
	} else {
		cW(buf, bRed, "%03d", status)
	}

	buf.WriteString(" in ")
	if dt < 500*time.Millisecond {
		cW(buf, nGreen, "%s", dt)
	} else if dt < 5*time.Second {
		cW(buf, nYellow, "%s", dt)
	} else {
		cW(buf, nRed, "%s", dt)
	}
}
