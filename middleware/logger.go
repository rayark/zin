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
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type LoggerHandler struct {
	handler http.Handler
	entry   LogEntry
}

func LoggerH(h http.Handler, entry LogEntry) LoggerHandler {
	return LoggerHandler{handler: h, entry: entry}
}

func (lh LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyWriter := NewProxyWriter(w)
	t1 := time.Now()
	lh.handler.ServeHTTP(proxyWriter, r)
	t2 := time.Now()
	logResult(proxyWriter, r, t2.Sub(t1), lh.entry)
}

func Logger(entry LogEntry) func(httprouter.Handle) httprouter.Handle {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			proxyWriter := NewProxyWriter(w)
			t1 := time.Now()
			h(proxyWriter, r, p)
			t2 := time.Now()
			logResult(proxyWriter, r, t2.Sub(t1), entry)
		}
	}
}

func findRemoteAddr(r *http.Request) string {
	addr := r.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return addr
}

func logResult(proxyWriter *ProxyWriter, r *http.Request, t time.Duration, log LogEntry) {
	ctx := r.Context()

	method := r.Method
	uri := r.URL.String()
	route, _ := GetRouteFromContext(ctx)
	sourceAddr := findRemoteAddr(r)
	msec := t.Milliseconds()
	status := proxyWriter.Status()
	uagent := r.Header.Get("User-Agent")

	entry := log.
		WithField("method", method).
		WithField("uri", uri).
		WithField("route", route).
		WithField("addr", sourceAddr).
		WithField("msec", msec).
		WithField("status", strconv.Itoa(status)).
		WithField("uagent", uagent)

	summary := fmt.Sprintf("%d %s %s from %s", status, method, uri, sourceAddr)

	if msec > 500 {
		summary = summary + fmt.Sprintf(" (%d msec)", msec)
	}

	if status <= 399 && msec <= 500 {
		entry.Infof(summary)
	} else if status <= 499 {
		entry.Warningf(summary)
	} else if status <= 599 {
		entry.Errorf(summary)
	}
}
