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

package middleware

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin/v2"
)

func middlewareCompressorTest(t *testing.T, reqHeaders map[string]string, respBody string) *httptest.ResponseRecorder {
	path := "/gzip"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	router := httprouter.New()
	grp := zin.NewGroup("/", Compressor)
	grp.R(router.GET, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte(respBody))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestWithGzip(t *testing.T) {
	reqHeaders := map[string]string{
		"Accept-Encoding": "gzip",
	}
	respBody := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	rec := middlewareCompressorTest(t, reqHeaders, respBody)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got '%d'", rec.Code)
	}

	if rec.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding: gzip got '%s'", rec.HeaderMap.Get("Content-Encoding"))
	}

	respBodySize := len(rec.Body.Bytes())

	gzr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(gzr)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != respBody {
		t.Fatalf(`expected "%s" got "%s"`, respBody, string(body))
	}

	log.Printf("Code: %d with Content-size(trans/orig): (%d/%d)", rec.Code, respBodySize, len(body))
}

func TestWithoutGzip(t *testing.T) {
	reqHeaders := map[string]string{}
	respBody := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	rec := middlewareCompressorTest(t, reqHeaders, respBody)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got '%d'", rec.Code)
	}

	if rec.HeaderMap.Get("Content-Encoding") != "" {
		t.Fatalf("expected Content-Encoding: '' got '%s'", rec.HeaderMap.Get("Content-Encoding"))
	}

	respBodySize := len(rec.Body.Bytes())

	if rec.Body.String() != respBody {
		t.Fatalf(`expected "%s" got "%s"`, respBody, rec.Body.String())
	}

	log.Printf("Code: %d with Content-size(trans/orig): (%d/%d)", rec.Code, respBodySize, len(respBody))
}

func TestNoBody(t *testing.T) {
	reqHeaders := map[string]string{
		"Accept-Encoding": "gzip",
	}
	respBody := ""

	rec := middlewareCompressorTest(t, reqHeaders, respBody)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got '%d'", rec.Code)
	}

	if rec.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding: gzip got '%s'", rec.HeaderMap.Get("Content-Encoding"))
	}

	respBodySize := len(rec.Body.Bytes())

	gzr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(gzr)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != respBody {
		t.Fatalf(`expected "%s" got "%s"`, respBody, string(body))
	}

	log.Printf("Code: %d with Content-size(trans/orig): (%d/%d)", rec.Code, respBodySize, len(body))
}
