package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CacheControl(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Add("Cache-Control", "public, s-maxage=86400")
		h(w, r, p)
	}
}
