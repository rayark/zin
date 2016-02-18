package zin

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestMakeHandle(t *testing.T) {
	data := ""
	m1 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m1")
			data = data + "A"
			h(w, r, p)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m2")
			data = data + "B"
			h(w, r, p)
		}
	}

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		t.Log("handle")
		data = data + "C"
	}

	fh := makeHandle([]Middleware{m1, m2}, h)

	if data != "" {
		t.Fail()
	}

	fh(nil, nil, nil)

	if data != "BAC" {
		t.Fail()
	}
}

func TestMakePooledHandle(t *testing.T) {

	m1 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			p[0].Value += "A"
			h(w, r, p)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			p[0].Value += "B"
			h(w, r, p)
		}
	}

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		p[0].Value += "C"
	}

	fh := makePooledHandle([]Middleware{m1, m2}, h)

	var wg sync.WaitGroup
	var data []string
	mutex := &sync.Mutex{}

	r := func() {
		params := []httprouter.Param{{Key: "Key", Value: ""}}
		fh(nil, nil, params)
		mutex.Lock()
		data = append(data, params[0].Value)
		mutex.Unlock()

		params[0].Value = ""
		fh(nil, nil, params)
		mutex.Lock()
		data = append(data, params[0].Value)
		mutex.Unlock()

		wg.Done()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go r()
	}

	wg.Wait()

	if len(data) != 200 {
		t.Fail()
	}

	for _, v := range data {
		if v != "BAC" {
			t.Fail()
		}
	}
}

func TestWrapM(t *testing.T) {
	params := []httprouter.Param{{Key: "Key", Value: ""}}
	data := "X"
	m1 := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			data += params[0].Value + "A"
			h(w, r)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			p[0].Value += "B"
			h(w, r, p)
		}
	}

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		p[0].Value += "C"
	}

	fh := makeHandle([]Middleware{WrapM(m1), m2}, h)

	fh(nil, nil, params)

	if data != "XBA" {
		t.Fail()
	}

	if params[0].Value != "BC" {
		t.Fail()
	}
}

func TestGroup(t *testing.T) {

	data := ""
	m1 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m1")
			data = data + "A"
			h(w, r, p)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m2")
			data = data + "B"
			h(w, r, p)
		}
	}

	router := httprouter.New()

	group := NewGroup("/test", m1, m2)

	group.R(router.GET, "admin/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data = data + "C"
		fmt.Fprint(w, data)
	})

	r, err := http.NewRequest("GET", "http://example.com/test/admin", nil)
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Body.String() != "BAC" {
		t.Fail()
	}

	if data != "BAC" {
		t.Fail()
	}
}

func TestChildGroup(t *testing.T) {

	data := ""
	m1 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m1")
			data = data + "A"
			h(w, r, p)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m2")
			data = data + "B"
			h(w, r, p)
		}
	}

	m3 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			t.Log("m3")
			data = data + "C"
			h(w, r, p)
		}
	}

	router := httprouter.New()

	group := NewGroup("/test", m1, m2)
	child := group.Group("/admin", m3)

	child.R(router.GET, "/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data = data + "D"
		t.Log(data)
		fmt.Fprint(w, data)
	})

	child.R(router.GET, "/ok", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data = data + "E"
		t.Log(data)
		fmt.Fprint(w, data)
	})

	r, err := http.NewRequest("GET", "http://example.com/test/admin", nil)
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Body.String() != "CBAD" {
		t.Fail()
	}

	if data != "CBAD" {
		t.Fail()
	}

	data = ""

	r, err = http.NewRequest("GET", "http://example.com/test/admin/ok", nil)
	if err != nil {
		panic(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Body.String() != "CBAE" {
		t.Fail()
	}

	if data != "CBAE" {
		t.Fail()
	}

}

func BenchmarkMakePooledHandle(t *testing.B) {
	data := 0
	m1 := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			data += 1
			h(w, r)
		}
	}

	m2 := func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			data -= 1
			h(w, r, p)
		}
	}

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data += 1
	}

	fh := makePooledHandle([]Middleware{WrapM(m1), m2}, h)
	params := []httprouter.Param{{Key: "Key", Value: ""}}

	var wg sync.WaitGroup
	r := func() {
		fh(nil, nil, params)
		wg.Done()
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go r()
		}
		wg.Wait()
	}

}
