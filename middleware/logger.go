package middleware

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func Logger(h httprouter.Handle) httprouter.Handle {
	var buf bytes.Buffer
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		buf.Reset()
		proxyWriter := NewProxyWriter(w)
		t1 := time.Now()
		h(proxyWriter, r, p)
		t2 := time.Now()

		method := r.Method
		url := r.URL.String()
		sourceAddr := findRemoteAddr(r)
		elapsed := t2.Sub(t1).Seconds()
		status := proxyWriter.Status()

		entry := log.WithFields(log.Fields{
			"method":      method,
			"url":         url,
			"source_addr": sourceAddr,
			"elapsed":     elapsed,
			"status":      status,
		})

		summary := fmt.Sprintf("%d %s %s from %s", status, method, url, sourceAddr)

		if elapsed > 0.5 {
			summary = summary + fmt.Sprintf(" (%f sec)", elapsed)
		}

		if status <= 399 && elapsed <= 0.5 {
			entry.Debug(summary)
		} else if status <= 499 && elapsed < 5 {
			entry.Warn(summary)
		} else if status <= 599 {
			entry.Error(summary)
		}
	}
}

func findRemoteAddr(r *http.Request) string {
	addr := r.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = r.RemoteAddr
	}

	return addr
}
