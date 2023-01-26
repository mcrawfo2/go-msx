// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idempotency

import "net/http"

type CachedWebData struct {
	Req  CachedRequest  `json:"req"`
	Resp CachedResponse `json:"resp"`
}

type CachedRequest struct {
	Method     string `json:"method"`
	RequestURI string `json:"requestURI"`
}

type CachedResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       []byte      `json:"data"`
	Header     http.Header `json:"header"`
}
