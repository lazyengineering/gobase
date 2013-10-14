// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	serverAddr   string
	staticDir    string
	rawTemplates map[string]string
	noTimestamp  bool
)

func init() {
	var err error

	// grab important metadata
	flag.StringVar(&serverAddr, "addr", ":5050", "Server Address to listen on")
	flag.StringVar(&staticDir, "static", "static", "Static Assets folder")
	flag.BoolVar(&noTimestamp, "noTimestamp", false, "When set to true, removes timestamp from log statements")
	flag.Parse()

	if noTimestamp {
		log.SetFlags(0)
	}

	// Static Asset Serving
	staticServer := http.FileServer(http.Dir(staticDir))
	logHandle("/js/", staticServer)
	logHandle("/css/", staticServer)
	logHandle("/fonts/", staticServer)
	logHandle("/img/", staticServer)
	logHandle("/favicon.ico", staticServer)

	// Raw Templates Loaded into memory
	var raw []byte
	rawTemplates = make(map[string]string)
	for _, filename := range []string{"layout.html", "hello.html"} {
		raw, err = ioutil.ReadFile("templates/" + filename)
		if err != nil {
			log.Fatalln("Fatal Error:", err)
			return
		}
		rawTemplates[filename] = string(raw)
	}

	// Baseline Template
	// Any use of this template needs to have a "body" template defined
	baseTemplate, err := template.New("layout").Parse(string(rawTemplates["layout.html"]))
	if err != nil {
		log.Fatalln("Fatal Error:", err)
		return
	}

	// Actual Web Application Handlers
	logHandle("/", handleFuncTemplate(baseTemplate, hello))
}

func logHandle(path string, h http.Handler) {
	http.HandleFunc(path, func(r http.ResponseWriter, q *http.Request) {
		log.Println("Serve:", q.URL.String())
		h.ServeHTTP(r, q)
	})
}

func main() {
	log.Println("listening at " + serverAddr + "...")
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Fatal Error:", err)
	}
}

type templateHandler func(http.ResponseWriter, *http.Request, *template.Template)

func handleFuncTemplate(t *template.Template, handler templateHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		plate, err := t.Clone()
		if err != nil {
			log.Println("Error:", err)
			http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
			return
		}
		handler(res, req, plate)
	}
}

func hello(res http.ResponseWriter, req *http.Request, t *template.Template) {
	_, err := t.New("body").Parse(rawTemplates["hello.html"])
	if err != nil {
		log.Println("Error: ", err)
		http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
		return
	}
  err = t.Execute(res, map[string]interface{}{"Title":"Hello World"})
	if err != nil {
		log.Println("Error: ", err)
		http.Error(res, "We seem to have an error on our end.", http.StatusInternalServerError)
		return
	}
}
