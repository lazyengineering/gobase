// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
)

// Important metadata
var (
	ServerAddr  = flag.String("addr", ":5050", "Server Address to listen on")
	StaticDir   = flag.String("static", "static", "Static Assets folder")
	NoTimestamp = flag.Bool("noTimestamp", false, "When set to true, removes timestamp from log statements")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

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

	// Actual Web Application Handlers
	HandleFunc("/", hello)
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

func NoIndex(h http.Handler) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, q *http.Request) {
		if strings.HasSuffix(q.URL.Path, "/") {
			Error403(r, q)
			return
		}
		h.ServeHTTP(r, q)
	})
}

func main() {
	log.Println("\x1b[32mlistening at \x1b[1;32m" + *ServerAddr + "\x1b[32m...\x1b[0m")
	log.Fatalln("Fatal Error:", http.ListenAndServe(*ServerAddr, nil))
}

func Error500(res http.ResponseWriter, req *http.Request, err error) {
	log.Println("\x1b[1;31mError:\x1b[0m", req.URL.String(), err)
	http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
}

func Error403(res http.ResponseWriter, req *http.Request) {
	log.Println("\x1b[1;31mNot Allowed:\x1b[0m", req.URL.String())
	http.Error(res, "Nothing to see here.", http.StatusForbidden)
}

func Error404(res http.ResponseWriter, req *http.Request) {
	log.Println("\x1b[1;31mNot Found:\x1b[0m", req.URL.String())
	http.Error(res, "We can't seem to find that.", http.StatusNotFound)
}

type Nav struct {
	*http.Request
}

func (n Nav) IsCurrent(p string) bool {
	return p == n.Request.URL.Path
}

func hello(res http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.URL.Path, "/") { // no sub-paths
		Error404(res, req)
		return
	}
	t, err := LoadTemplates("templates/hello/*.html")
	if err != nil {
		Error500(res, req, err)
		return
	}
	err = t.ExecuteTemplate(res, "bootstrap.html", map[string]interface{}{
		"Title":     "Hello World",
		"BodyClass": "hello",
		"Nav":       Nav{req},
	})
	if err != nil {
		Error500(res, req, err)
		return
	}
}
