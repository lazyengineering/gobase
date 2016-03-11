// Copyright 2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package middleware

import (
	"io"
	"testing"
	"time"
)

// Cache should set the Cache-Control and Expires Headers to match ttl before calling h.ServeHTTP
// Test Cases:
//   - ttl = -1h, 0, 1h
//   - h = simple, sets different cache
// TODO(jesse@jessecarl.com): Define behavior for Cache based on ttl values. We need to be able to correctly handle positive, negative, and zero values.
func TestCache(t *testing.T) {

}

func simpleHandler(w http.ResponseWriter, q *http.Request) {
	io.WriteString(w, "Hello, simple Handler")
}

func cachingHandler(w http.ResponseWriter, q *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=120")
	w.Header().Set("Expires", time.Now().Add(120*time.Second).Format(time.RFC1123))
	io.WriteString(w, "Hello, caching Handler")
}
