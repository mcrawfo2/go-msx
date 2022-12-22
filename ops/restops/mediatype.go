// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

// Simple media types
const (
	MediaTypeJson           = "application/json"
	MediaTypeXml            = "text/xml"
	MediaTypeBinary         = "application/octet-stream"
	MediaTypeFormUrlencoded = "application/x-www-form-urlencoded"
	MediaTypeMultipartForm  = "multipart/form-data"
	MediaTypeMultipartMixed = "multipart/mixed"
	MediaTypeTextPlain      = "text/plain"
)

// Content types
const (
	ContentTypeJson = MediaTypeJson + "; charset=utf-8"
	ContentTypeXml  = MediaTypeXml + "; charset=utf-8"
)

// Content encodings
const (
	ContentEncodingBase64  = "base64"
	ContentEncodingGzip    = "gzip"
	ContentEncodingDeflate = "deflate"
	ContentEncodingNone    = ""
)
