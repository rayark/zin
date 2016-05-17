package middleware

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type LoggerHandler struct {
	handler http.Handler
}

func LoggerH(h http.Handler) LoggerHandler {
	return LoggerHandler{handler: h}
}

func (lh LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyWriter := NewProxyWriter(w)
	t1 := time.Now()
	lh.handler.ServeHTTP(proxyWriter, r)
	t2 := time.Now()
	logResult(proxyWriter, r, t2.Sub(t1))
}

func Logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		proxyWriter := NewProxyWriter(w)
		t1 := time.Now()
		h(proxyWriter, r, p)
		t2 := time.Now()
		logResult(proxyWriter, r, t2.Sub(t1))
	}
}

func findRemoteAddr(r *http.Request) string {
	addr := r.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = r.RemoteAddr
	}

	return addr
}

func logResult(proxyWriter *ProxyWriter, r *http.Request, t time.Duration) {
	method := r.Method
	url := r.URL.String()
	sourceAddr := findRemoteAddr(r)
	elapsed := t.Seconds()
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
		entry.Info(summary)
	} else if status <= 499 && elapsed < 5 {
		entry.Warn(summary)
	} else if status <= 599 {
		entry.Error(summary)
	}
}
