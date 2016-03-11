// Copyright 2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// Adds Cache headers to every response. This is useful when serving static data
// Headers are Set before calling `h.ServeHTTP()`
func Cache(ttl time, h http.Handler) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
		r.Header().Set("Cache-Control", "public, max-age="+strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64))
		r.Header().Set("Expires", time.Now().Add(ttl).Format(time.RFC1123))
		h.ServeHTTP(r, q)
	})
}
