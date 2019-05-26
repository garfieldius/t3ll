// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"net/http"
	"strconv"
)

type response struct {
	status int
	ctype  string
	body   []byte
}

func (r response) write(w http.ResponseWriter) {
	if r.status > 199 && r.status < 504 {
		w.WriteHeader(r.status)
	} else {
		w.WriteHeader(200)
	}

	if r.ctype == "" {
		r.ctype = "application/json"
	}

	w.Header().Set("Content-Type", r.ctype+";charset=UTF-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(r.body)))
	w.Write(r.body)
}