package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin"
)

// NewHMACAuthenticator returns a wrapper middleware to add with
func NewHMACAuthenticator(headerKey, secret string) zin.Middleware {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil
	}
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			deferWriter := NewDeferWriter(w)
			defer deferWriter.WriteAll()

			h(deferWriter, r, p)
			hmac := computeHMAC256(deferWriter.Bytes(), key)
			deferWriter.Header().Set(headerKey, hmac)
		}
	}
}

func computeHMAC256(msg, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
