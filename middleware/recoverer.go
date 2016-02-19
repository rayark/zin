package middleware

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"runtime/debug"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Recoverer prints a request ID if one is provided.
func Recoverer(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				printPanic(err)
				debug.PrintStack()
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		h(w, r, p)
	}
}

func printPanic(err interface{}) {
	var buf bytes.Buffer
	cW(&buf, bRed, "panic: %+v", err)
	log.Print(buf.String())
}
