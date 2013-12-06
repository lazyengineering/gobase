// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lazyengineering/gobase/envflag"
	"github.com/lazyengineering/gobase/layouts"
)

// Important metadata
var (
	ServerAddr   = flag.String("server-addr", ":5050", "Server Address to listen on")
	GATrackingID = flag.String("ga-tracking-id", "", "Google Analytics Tracking ID")
)

var Layout *layouts.Layout

func init() {
	t := time.Now() // measure bootstrap time
	defer func() {
		log.Printf("\x1b[1;32mBootstrapped:\x1b[0m \x1b[34m%8d\x1b[0mµs", time.Since(t).Nanoseconds()/1000)
	}()

	// Flags only useful in init
	var (
		NoTimestamp        = flag.Bool("no-timestamp", false, "When set to true, removes timestamp from log statements")
		StaticDir          = flag.String("static-dir", "static", "Static Assets folder")
		LayoutTemplateGlob = flag.String("layouts", "static/templates/layouts/*.html", "Pattern for layout templates")
		HelperTemplateGlob = flag.String("helpers", "static/templates/helpers/*.html", "Pattern for helper templates")
	)

	// To Parse flags, looking for command-line, then ENV, then defaults
	envflag.Parse(envflag.FlagMap{
		"server-addr": envflag.Flag{
			Name:   "PORT",
			Filter: func(s string) string { return ":" + s },
		},
	})

	if *NoTimestamp {
		log.SetFlags(0)
	}

	// Static Asset Serving
	staticServer := NoIndex(func(h http.Handler) http.Handler {
		// add 1 day caching headers to static assets
		return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
			r.Header().Set("Cache-Control", "public, max-age=86400")
			r.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(time.RFC1123))
			h.ServeHTTP(r, q)
		})
	}(http.FileServer(http.Dir(*StaticDir))))
	Handle("/js/", staticServer)
	Handle("/css/", staticServer)
	Handle("/fonts/", staticServer)
	Handle("/img/", staticServer)
	Handle("/favicon.ico", staticServer)

	// Layouts
	var err error
	Layout, err = layouts.New(layouts.BasicFunctionMap(), "bootstrap.html", *LayoutTemplateGlob, *HelperTemplateGlob)
	if err != nil {
		// this is a fatal condition
		panic(err)
	}

	// Actual Web Application Handlers
	HandleNoSubPaths("/", Layout.Act(hello, Error500, layouts.NoVolatility, "static/templates/hello/*.html"))
}

// Log and Handle http requests
func Handle(path string, h http.Handler) {
	if strings.HasSuffix(path, "/") { // redirect for directories
		indexRedirect := http.RedirectHandler(path, http.StatusMovedPermanently)
		Handle(path+"index.html", indexRedirect)
		Handle(path+"index.htm", indexRedirect)
		Handle(path+"index.php", indexRedirect) // not that anybody would think...
	}
	http.HandleFunc(path, func(r http.ResponseWriter, q *http.Request) {
		t := time.Now()
		h.ServeHTTP(r, q)
		s := time.Since(t).Nanoseconds() / 1000 // time in µs
		message := "\x1b[1;36mServed: \x1b[0m"
		if s > 10000 { // > 10ms is a "long" request, mark in red with a *
			message = "\x1b[1;31mServed*:\x1b[0m"
		}
		log.Printf("%s \x1b[34m%8d\x1b[0mµs \x1b[33m%s\x1b[0m", message, s, q.URL.String())
	})
}

func HandleFunc(path string, h http.HandlerFunc) {
	Handle(path, http.HandlerFunc(h))
}

func HandleNoSubPaths(path string, h http.Handler) {
	Handle(path, NoSubPaths(path, h))
}

func NoIndex(h http.Handler) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
		if strings.HasSuffix(q.URL.Path, "/") {
			Error403(r, q)
			return
		}
		h.ServeHTTP(r, q)
	})
}

func NoSubPaths(path string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
		if q.URL.Path != path {
			Error404(r, q)
			return
		}
		h.ServeHTTP(r, q)
	})
}

func main() {
	log.Println("\x1b[32mlistening at \x1b[1;32m" + *ServerAddr + "\x1b[32m...\x1b[0m")
	log.Fatalln("Fatal Error:", http.ListenAndServe(*ServerAddr, nil))
}

type Nav struct {
	*http.Request
}

func (n Nav) IsCurrent(p string) bool {
	return p == n.Request.URL.Path
}

func hello(req *http.Request) (map[string]interface{}, error) {
	return map[string]interface{}{
		"Title":        "Hello World",
		"BodyClass":    "hello",
		"Nav":          Nav{req},
		"GATrackingID": *GATrackingID,
	}, nil
}
