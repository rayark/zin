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

package zin

import (
	"net/http"
	"path"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin/v2/middleware"
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

// Deprecated. Use Group or Pack.
// https://stackoverflow.com/questions/53572736/append-to-a-new-slice-affect-original-slice
func (g *MuxGroup) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

type RegisterFunc func(path string, handle httprouter.Handle)

func (g *MuxGroup) R(r RegisterFunc, p string, handle httprouter.Handle) {
	route := g.Path(p)
	m := safeAppend(g.middlewares, middleware.AddRouteToContext(route))
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

// Group returns new MuxGroup with appending inputs middlewares to the end of current muxgroup's middleware
func (g *MuxGroup) Group(path string, middlewares ...Middleware) *MuxGroup {
	return NewGroup(pathJoin(g.basePath, path), safeAppend(g.middlewares, middlewares...)...)
}

// Pack returns new MuxGroup with appending current muxgroup's middlewares to the end of input middlewarse
func (g *MuxGroup) Pack(path string, middlewares ...Middleware) *MuxGroup {
	return NewGroup(pathJoin(g.basePath, path), safeAppend(middlewares, g.middlewares...)...)
}

func safeAppend(middlewaresA []Middleware, middlewaresB ...Middleware) []Middleware {
	mws := make([]Middleware, 0, len(middlewaresA)+len(middlewaresB))
	mws = append(mws, middlewaresA...)
	mws = append(mws, middlewaresB...)
	return mws
}
