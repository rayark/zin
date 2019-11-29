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
