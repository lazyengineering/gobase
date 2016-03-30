// Copyright 2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// Cache sets the correct CacheControl headers according to ttl.
// A ttl <= 0 will result in a no-cache value.
// As with any middleware here, Header values can be overwritten by the next handler.
func Cache(ttl time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
		if ttl <= 0 {
			r.Header().Set("Cache-Control", "no-cache")
			r.Header().Set("Expires", time.Now().Add(-2*time.Minute).Format(time.RFC1123))
		} else {
			r.Header().Set("Cache-Control", "public, max-age="+strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64))
			r.Header().Set("Expires", time.Now().Add(ttl).Format(time.RFC1123))
		}
		h.ServeHTTP(r, q)
	})
}
