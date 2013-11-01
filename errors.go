// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"log"
	"net/http"
)

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
