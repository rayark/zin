package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin"
	log "github.com/sirupsen/logrus"
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
		addr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return addr
}

func logResult(proxyWriter *ProxyWriter, r *http.Request, t time.Duration) {
	ctx := r.Context()

	method := r.Method
	uri := r.URL.String()
	route, _ := ctx.Value(zin.MatchedRoutePathKey).(string)
	sourceAddr := findRemoteAddr(r)
	msec := t.Nanoseconds() / 1000000
	status := proxyWriter.Status()
	uagent := r.Header.Get("User-Agent")

	entry := log.WithFields(log.Fields{
		"method": method,
		"uri":    uri,
		"route":  route,
		"addr":   sourceAddr,
		"msec":   msec,
		"status": strconv.Itoa(status),
		"uagent": uagent,
	})

	summary := fmt.Sprintf("%d %s %s from %s", status, method, uri, sourceAddr)

	if msec > 500 {
		summary = summary + fmt.Sprintf(" (%d msec)", msec)
	}

	if status <= 399 && msec <= 500 {
		entry.Info(summary)
	} else if status <= 499 {
		entry.Warn(summary)
	} else if status <= 599 {
		entry.Error(summary)
	}
}
