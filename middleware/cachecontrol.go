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
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CacheControl(age int) func(h httprouter.Handle) httprouter.Handle {
	value := fmt.Sprintf("public, s-maxage=%d", age)
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Header().Add("Cache-Control", value)
			h(w, r, p)
		}
	}
}
