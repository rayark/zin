package middleware

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CacheControl(age int) func(h httprouter.Handle) httprouter.Handle {
	value := fmt.Sprintf("public, s-maxage=%s", age)
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Header().Add("Cache-Control", value)
			h(w, r, p)
		}
	}
}
