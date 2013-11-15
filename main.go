// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lazyengineering/gobase/templates"
)

// Important metadata
var (
	ServerAddr   = flag.String("server-addr", ":5050", "Server Address to listen on")
	GATrackingID = flag.String("ga-tracking-id", "", "Google Analytics Tracking ID")
)

var Templates templates.Collection

func init() {
	var (
		NoTimestamp        = flag.Bool("no-timestamp", false, "When set to true, removes timestamp from log statements")
		StaticDir          = flag.String("static-dir", "static", "Static Assets folder")
		LayoutTemplateGlob = flag.String("layouts", "static/templates/layouts/*.html", "Pattern for layout templates")
		HelperTemplateGlob = flag.String("helpers", "static/templates/helpers/*.html", "Pattern for helper templates")
	)

	// set from environment where available before parsing (allows flags to overrule env)
	flag.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "server-addr": // special case because it doesn't map directly
			if port := os.Getenv("PORT"); len(port) > 0 {
				f.Value.Set(":" + port)
			}
		default:
			if v := os.Getenv(strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))); len(v) > 0 {
				f.Value.Set(v)
			}
		}
	})

	flag.Parse()

	if *NoTimestamp {
		log.SetFlags(0)
	}

	// Static Asset Serving
	staticServer := NoIndex(http.FileServer(http.Dir(*StaticDir)))
	Handle("/js/", staticServer)
	Handle("/css/", staticServer)
	Handle("/fonts/", staticServer)
	Handle("/img/", staticServer)
	Handle("/favicon.ico", staticServer)

	// Templates
	Templates = templates.Collection{
		LayoutGlob: *LayoutTemplateGlob,
		HelperGlob: *HelperTemplateGlob,
		Functions:  templates.BasicFunctionMap(),
	}

	// Actual Web Application Handlers
	HandleNoSubPaths("/", http.HandlerFunc(hello))
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
		log.Printf("\x1b[1;36mServed:\x1b[0m \x1b[34m%6d\x1b[0mµs \x1b[33m%s\x1b[0m", time.Since(t).Nanoseconds()/1000, q.URL.String())
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

func hello(res http.ResponseWriter, req *http.Request) {
	t, err := Templates.Load("static/templates/hello/*.html")
	if err != nil {
		Error500(res, req, err)
		return
	}
	// write to a buffer to eliminate any half-executed templates from being written
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "bootstrap.html", map[string]interface{}{
		"Title":        "Hello World",
		"BodyClass":    "hello",
		"Nav":          Nav{req},
		"GATrackingID": *GATrackingID,
	})
	if err != nil {
		Error500(res, req, err)
		return
	}
	_, err = b.WriteTo(res)
	if err != nil {
		Error500(res, req, err)
		return
	}
}
