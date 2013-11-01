// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
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
	staticServer := http.FileServer(http.Dir(*StaticDir))
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
	http.HandleFunc(path, func(r http.ResponseWriter, q *http.Request) {
		t := time.Now()
		h.ServeHTTP(r, q)
		log.Printf("\x1b[1;36mServed:\x1b[0m \x1b[34m%6d\x1b[0mµs \x1b[33m%s\x1b[0m", time.Since(t).Nanoseconds()/1000, q.URL.String())
	})
}

func HandleFunc(path string, h http.HandlerFunc) {
	Handle(path, http.HandlerFunc(h))
}

func main() {
	log.Println("\x1b[32mlistening at \x1b[1;32m" + *ServerAddr + "\x1b[32m...\x1b[0m")
	log.Fatalln("Fatal Error:", http.ListenAndServe(*ServerAddr, nil))
}

func Error500(res http.ResponseWriter, req *http.Request, err error) {
	log.Println("\x1b[1;31mError:\x1b[0m", req.URL.String(), err)
	http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
}

func hello(res http.ResponseWriter, req *http.Request) {
	t, err := LoadTemplates("templates/hello/*.html")
	if err != nil {
		Error500(res, req, err)
		return
	}
	err = t.ExecuteTemplate(res, "bootstrap.html", map[string]interface{}{
		"Title":     "Hello World",
		"BodyClass": "hello",
	})
	if err != nil {
		Error500(res, req, err)
		return
	}
}
