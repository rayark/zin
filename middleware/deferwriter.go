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
	"bytes"
	"net/http"
)

type DeferWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func NewDeferWriter(w http.ResponseWriter) *DeferWriter {
	return &DeferWriter{
		ResponseWriter: w,
		buf:            new(bytes.Buffer),
	}
}

func (w *DeferWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *DeferWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *DeferWriter) WriteAll() {
	w.ResponseWriter.Write(w.buf.Bytes())
}
