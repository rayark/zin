package zin

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MuxGroup struct {
	basePath string
}

type StdFuncMiddlware func(http.HandlerFunc) http.HandlerFunc
type Middleware func(httprouter.Handle) httprouter.Handle

func WrapM(sm StdFuncMiddlware) Middleware {

	return func(h httprouter.Handle) httprouter.Handle {
		var params httprouter.Params

		stdh := func(w http.ResponseWriter, r *http.Request) {
			h(w, r, params)
		}

		stdh2 := sm(stdh)

		h2 := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			params = p
			stdh2(w, r)
		}

		return h2
	}
}
