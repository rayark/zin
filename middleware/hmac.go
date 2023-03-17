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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HMACSHA1Signer returns a middleware wrapper to add hmac signing string in
// response header
func HMACSHA1Signer(hmacHeaderKey, nounceHeaderKey string, secret []byte) middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			deferWriter := NewDeferWriter(w)
			defer deferWriter.WriteAll()

			key := secret
			if nounceHeaderKey != "" {
				nounceInHex := r.Header.Get(nounceHeaderKey)
				nounce, err := hex.DecodeString(nounceInHex)
				if err == nil {
					key = append(key, nounce...)
				}
			}

			h(deferWriter, r, p)
			hmacSignature := generateSignature(deferWriter.Bytes(), key)
			deferWriter.Header().Set(hmacHeaderKey, hmacSignature)
		}
	}
}

func generateSignature(msg, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(msg)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}
