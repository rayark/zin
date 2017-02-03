package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"compress/gzip"

	"io/ioutil"

	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin"
)

func TestWithGzip(t *testing.T) {
	origText := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	req, err := http.NewRequest("GET", "/gzip", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	router := httprouter.New()
	grp := zin.NewGroup("/", Compressor)
	grp.R(router.GET, "/gzip", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte(origText))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got '%d'", rec.Code)
	}

	if rec.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding: gzip got '%s'", rec.HeaderMap.Get("Content-Encoding"))
	}

	origSize := len(rec.Body.Bytes())

	gzr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(gzr)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != origText {
		t.Fatalf(`expected "%s" got "%s"`, origText, string(body))
	}

	log.Printf("Code: %d with Content-size: (%d/%d)", rec.Code, origSize, len(body))
}

func TestWithoutGzip(t *testing.T) {

}

func TestNoBody(t *testing.T) {

}
