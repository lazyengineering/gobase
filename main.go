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
		url := q.URL.String()
		t := time.Now()
		h.ServeHTTP(r, q)
		log.Println("Served:", url, "in", time.Since(t).Nanoseconds()/1000, "µs")
	})
}

func HandleFunc(path string, h http.HandlerFunc) {
	Handle(path, http.HandlerFunc(h))
}

func main() {
	log.Println("listening at " + *ServerAddr + "...")
	err := http.ListenAndServe(*ServerAddr, nil)
	if err != nil {
		log.Fatal("Fatal Error:", err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	t, err := LoadTemplates("templates/hello/*.html")
	if err != nil {
		log.Println("Error: ", err)
		http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(res, "bootstrap.html", map[string]interface{}{"Title": "Hello World"})
	if err != nil {
		log.Println("Error: ", err)
		http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
		return
	}
}
