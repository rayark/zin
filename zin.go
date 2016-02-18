package zin

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
	"sync"
)

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

func WrapF(f http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		f(w, r)
	}
}

func WrapH(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

type MuxGroup struct {
	basePath    string
	middlewares []Middleware
}

func NewGroup(basePath string, middlewares ...Middleware) *MuxGroup {
	return &MuxGroup{
		basePath:    basePath,
		middlewares: middlewares,
	}
}

func (g *MuxGroup) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

type RegisterFunc func(path string, handle httprouter.Handle)

func (g *MuxGroup) R(r RegisterFunc, path string, handle httprouter.Handle) {
	r(pathJoin(g.basePath, path), makePooledHandle(g.middlewares, handle))
}

func makePooledHandle(middlewares []Middleware, handle httprouter.Handle) httprouter.Handle {

	pool := sync.Pool{}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		obj := pool.Get()
		var h httprouter.Handle

		if obj != nil {
			h = obj.(httprouter.Handle)
		} else {
			h = makeHandle(middlewares, handle)
		}

		h(w, r, p)
		pool.Put(h)
	}
}

func makeHandle(middlewares []Middleware, handle httprouter.Handle) httprouter.Handle {
	var h = handle
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func pathJoin(base string, r string) string {
	return path.Join(base, r)
}

func (g *MuxGroup) Group(path string, middlewares ...Middleware) *MuxGroup {
	return NewGroup(pathJoin(g.basePath, path), append(g.middlewares, middlewares...)...)
}
