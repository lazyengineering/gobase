// Copyright 2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Cache should set the Cache-Control and Expires Headers to match ttl before calling h.ServeHTTP
// Test Cases:
//   - ttl = -1h, 0, 1h
//   - h = simple, sets different cache
// TODO(jesse@jessecarl.com): Define behavior for Cache based on ttl values. We need to be able to correctly handle positive, negative, and zero values.
func TestCache(t *testing.T) {
	type expiresTest struct {
		Test   func(time.Time) bool // use to create custom comparison
		String func() string
	}
	type testCase struct {
		// Input
		Ttl time.Duration
		H   http.Handler // simpleHandler or cachingHandler

		// Expectations
		Expires      expiresTest
		CacheControl string // should be a specific string
	}

	// ttl h      | EXPIRES CACHE-CONTROL
	// -1h simple | Past    no-cache
	// 0   simple | Past    no-cache
	// 1h  simple | ~1h     public, max-age=3600
	// -1h cache  | ~2m     public, max-age=120
	// 0   cache  | ~2m     public, max-age=120
	// 1h  cache  | ~2m     public, max-age=120
	testCases := []testCase{
		{ // -1h simple | Past    no-cache
			Ttl: -time.Hour,
			H:   http.HandlerFunc(simpleHandler),
			Expires: expiresTest{
				Test: func(rt time.Time) bool {
					return rt.Before(time.Now())
				},
				String: func() string { return "Past" },
			},
			CacheControl: "no-cache",
		},
		{ // 0   simple | Past    no-cache
			Ttl: 0,
			H:   http.HandlerFunc(simpleHandler),
			Expires: expiresTest{
				Test: func(rt time.Time) bool {
					return rt.Before(time.Now())
				},
				String: func() string { return "Past" },
			},
			CacheControl: "no-cache",
		},
		{ // 1h  simple | ~1h     public, max-age=3600
			Ttl: time.Hour,
			H:   http.HandlerFunc(simpleHandler),
			// let's try for less than 1 second (very doubtful that it would be off by more than that)
			Expires: expiresTest{
				Test:   approxTimeEq(time.Hour, time.Second),
				String: func() string { return "1h" },
			},
			CacheControl: "public, max-age=3600",
		},
		{ // -1h cache  | ~2m     public, max-age=120
			Ttl: -time.Hour,
			H:   http.HandlerFunc(cachingHandler),
			Expires: expiresTest{
				Test:   approxTimeEq(2*time.Minute, time.Second),
				String: func() string { return "2m" },
			},
			CacheControl: "public, max-age=120",
		},
		{ // 0   cache  | ~2m     public, max-age=120
			Ttl: time.Duration(0),
			H:   http.HandlerFunc(cachingHandler),
			Expires: expiresTest{
				Test:   approxTimeEq(2*time.Minute, time.Second),
				String: func() string { return "2m" },
			},
			CacheControl: "public, max-age=120",
		},
		{ // 1h  cache  | ~2m     public, max-age=120
			Ttl: time.Hour,
			H:   http.HandlerFunc(cachingHandler),
			Expires: expiresTest{
				Test:   approxTimeEq(2*time.Minute, time.Second),
				String: func() string { return "2m" },
			},
			CacheControl: "public, max-age=120",
		},
	}

	for idx, tc := range testCases {
		service := httptest.NewServer(Cache(tc.Ttl, tc.H))
		r, err := http.Get(service.URL)
		if err != nil {
			t.Error(err)
		}

		// expires
		if expires, err := time.Parse(time.RFC1123, r.Header.Get("Expires")); err != nil {
			t.Error(err)
		} else if !tc.Expires.Test(expires) {
			t.Error("test\t", idx, "\texpected: Expires:", tc.Expires.String(), "\tactual: Expires:", time.Since(expires).String())
		}

		// cache-control
		if cc := r.Header.Get("Cache-Control"); cc != tc.CacheControl {
			t.Error("test\t", idx, "\texpected: Cache-Control:", tc.CacheControl, "\tactual: Cache-Control:", cc)
		}
	}
}

func simpleHandler(w http.ResponseWriter, q *http.Request) {
	io.WriteString(w, "Hello, simple Handler")
}

func cachingHandler(w http.ResponseWriter, q *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=120")
	w.Header().Set("Expires", time.Now().Add(120*time.Second).Format(time.RFC1123))
	io.WriteString(w, "Hello, caching Handler")
}

func approxTimeEq(d, by time.Duration) func(time.Time) bool {
	return func(t time.Time) bool {
		actual := -time.Since(t)
		return actual > (d-by) && actual < (d+by)
	}
}
