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
	log "github.com/sirupsen/logrus"
)

type RecovererHandler struct {
	handler http.Handler
}

func RecovererH(h http.Handler) RecovererHandler {
	return RecovererHandler{handler: h}
}

func (rh *RecovererHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"call_stack": string(debug.Stack()),
			}).Errorf("panic: %+v", err)
			http.Error(w, http.StatusText(500), 500)
		}
	}()

	rh.ServeHTTP(w, r)
}

func Recoverer(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"call_stack": string(debug.Stack()),
				}).Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		h(w, r, p)
	}
}
