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
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

type RecovererHandler struct {
	handler http.Handler
	entry   LogEntry
}

func RecovererH(h http.Handler, entry LogEntry) RecovererHandler {
	return RecovererHandler{handler: h, entry: entry}
}

func (rh *RecovererHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rh.entry.
				WithField("call_stack", string(debug.Stack())).
				Errorf("panic: %+v", err)
			http.Error(w, http.StatusText(500), 500)
		}
	}()

	rh.handler.ServeHTTP(w, r)
}

func Recoverer(entry LogEntry) func(httprouter.Handle) httprouter.Handle {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			defer func() {
				if err := recover(); err != nil {
					entry.
						WithField("call_stack", string(debug.Stack())).
						Errorf("panic: %+v", err)
					http.Error(w, http.StatusText(500), 500)
				}
			}()

			h(w, r, p)
		}
	}
}
