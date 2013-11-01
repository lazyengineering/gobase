// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"log"
	"net/http"
)

const page500 = `<html>
<head><title>Oops! Something went wrong!</title></head>
<body>
<h1>Oops! Something went wrong!</h1>
<p>Don't worry, we'll be looking into the problem, and have everything back to you in tip-top shape soon enough</p>
</body>
</html>
`

func Error500(res http.ResponseWriter, req *http.Request, err error) {
	log.Println("\x1b[1;31mError:\x1b[0m", req.URL.String(), err)
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(page500))
}

const page404 = `<html>
<head><title>Oops! We can't seem to find that!</title></head>
<body>
<h1>We can't seem to locate the resource you are looking for.</h1>
</body>
</html>`

func Error404(res http.ResponseWriter, req *http.Request) {
	log.Println("\x1b[1;31mNot Found:\x1b[0m", req.URL.String())
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(page404))
}

const page403 = `<html>
<head><title>No! You can't be here!</title></head>
<body>
<h1>You have stumbled upon a locked door, turn around and go the other way.</h1>
</body>
</html>`

func Error403(res http.ResponseWriter, req *http.Request) {
	log.Println("\x1b[1;31mNot Allowed:\x1b[0m", req.URL.String())
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusForbidden)
	res.Write([]byte(page403))
}
