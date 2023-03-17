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
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin/v2"
)

const (
	hmacHeaderKey   = "HMAC-Signature-Hash"
	nounceHeaderKey = "Nounce-For-HMAC"
	nounceString    = "this is nounce"
	secretString    = "ThisIsSecret"
	bodyContent     = "this is http body content"
)

func middlewareHMACTest(t *testing.T, reqHeaders map[string]string, nouceHeaderKey string) *httptest.ResponseRecorder {
	path := "/hmac"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	msg := []byte(bodyContent)
	key := []byte(secretString)

	hmacWrapper := HMACSHA1Signer(hmacHeaderKey, nouceHeaderKey, key)
	router := httprouter.New()
	grp := zin.NewGroup("/", hmacWrapper)
	grp.R(router.GET, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Language", "klingon")
		w.Write(msg)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestHMACSignature(t *testing.T) {
	reqHeaders := map[string]string{}

	rec := middlewareHMACTest(t, reqHeaders, "")

	valueExpected := generateSignature([]byte(bodyContent), []byte(secretString))

	// check if regular header value consistent or not
	if rec.HeaderMap.Get("Content-Language") != "klingon" {
		t.Fatalf("standard header inconsistent")
	}

	// check if "hmac" header value identical to the expected one
	if rec.HeaderMap.Get(hmacHeaderKey) != valueExpected {
		t.Fatalf("appended header inconsistent")
	}

	// check if response body consistent or not
	if string(rec.Body.Bytes()) != bodyContent {
		t.Fatalf("content body inconsistent")
	}
}

func TestHMACSignatureWithNounce(t *testing.T) {
	nouceInHex := hex.EncodeToString([]byte(nounceString))
	reqHeaders := map[string]string{
		nounceHeaderKey: nouceInHex,
	}

	rec := middlewareHMACTest(t, reqHeaders, nounceHeaderKey)

	key := append([]byte(secretString), []byte(nounceString)...)
	valueExpected := generateSignature([]byte(bodyContent), key)

	// check if regular header value consistent or not
	if rec.HeaderMap.Get("Content-Language") != "klingon" {
		t.Fatalf("standard header inconsistent")
	}

	// check if "hmac" header value identical to the expected one
	if rec.HeaderMap.Get(hmacHeaderKey) != valueExpected {
		t.Fatalf("appended header inconsistent")
	}

	// check if response body consistent or not
	if string(rec.Body.Bytes()) != bodyContent {
		t.Fatalf("content body inconsistent")
	}
}
