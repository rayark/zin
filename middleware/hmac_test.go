package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin"
)

func TestHMACAuthenticator(t *testing.T) {
	path := "/hmac"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	hmacHeaderKey := "HMAC-Authenticate-Hash"
	secretInBytes := []byte("ThisIsSecret")
	secretInBase64 := base64.StdEncoding.EncodeToString(secretInBytes)
	message := []byte("this is http body content")
	hmac := computeHMAC256(message, secretInBytes)

	hmacWrapper := NewHMACAuthenticator(hmacHeaderKey, secretInBase64)
	router := httprouter.New()
	grp := zin.NewGroup("/", hmacWrapper)
	grp.R(router.GET, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Language", "klingon")
		w.Write(message)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// check if regular header value consistent or not
	if rec.HeaderMap.Get("Content-Language") != "klingon" {
		t.Fatalf("standard header inconsistent")
	}

	// check if "hmac" header value identical to the expected one
	if rec.HeaderMap.Get(hmacHeaderKey) != hmac {
		t.Fatalf("appended header inconsistent")
	}

	// check if response body consistent or not
	if string(rec.Body.Bytes()) != string(message) {
		t.Fatalf("content body inconsistent")
	}
}
