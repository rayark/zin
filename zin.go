package zin

import (
	"context"
	"net/http"
	"path"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type StdMiddlware func(http.Handler) http.Handler
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

func WrapS(sm StdMiddlware) Middleware {

	return func(h httprouter.Handle) httprouter.Handle {
		var params httprouter.Params

		stdh := func(w http.ResponseWriter, r *http.Request) {
			h(w, r, params)
		}
		stdh2 := sm(http.HandlerFunc(stdh))

		h2 := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			params = p
			stdh2.ServeHTTP(w, r)
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

func (g *MuxGroup) R(r RegisterFunc, p string, handle httprouter.Handle) {
	route := g.Path(p)
	m := append(g.middlewares, addRouteToCtxMiddleware(route))
	r(route, makePooledHandle(m, handle))
}

func (g *MuxGroup) Path(p string) string {
	return pathJoin(g.basePath, p)
}

func (g *MuxGroup) NotFound(h http.Handler) http.Handler {
	handle := makePooledHandle(g.middlewares, WrapH(h))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, nil)
	})
}

type zinContextKey string

const MatchedRoutePathKey = zinContextKey("MatchRoutePath")

func addRouteToCtxMiddleware(route string) Middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, MatchedRoutePathKey, route)

			h(w, r.WithContext(ctx), p)
		}
	}
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
	path := path.Join(base, r)
	if len(r) > 0 && r[len(r)-1] == '/' && path[len(r)-1] != '/' {
		path = path + "/"
	}
	return path
}

func (g *MuxGroup) Group(path string, middlewares ...Middleware) *MuxGroup {
	return NewGroup(pathJoin(g.basePath, path), append(g.middlewares, middlewares...)...)
}
