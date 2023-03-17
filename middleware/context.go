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
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type zinContextKey int

const MatchedRoutePathKey zinContextKey = iota

func AddRouteToContext(route string) middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, MatchedRoutePathKey, route)

			h(w, r.WithContext(ctx), p)
		}
	}
}

func GetRouteFromContext(ctx context.Context) (string, bool) {
	route, ok := ctx.Value(MatchedRoutePathKey).(string)
	return route, ok
}
