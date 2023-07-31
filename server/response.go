// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"github.com/garfieldius/t3ll/log"
	"github.com/kr/pretty"
	"net/http"
	"strconv"
)

type response struct {
	status int
	ctype  string
	body   []byte
	dlName string
}

func (r response) write(w http.ResponseWriter) {
	if r.ctype == "" {
		r.ctype = "application/json"
	}

	w.Header().Set("Content-Type", r.ctype+"; charset=UTF-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(r.body)))

	if r.dlName != "" {
		log.Msg("Sending content-disposition header")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+r.dlName+"\"")
	}

	log.Msg("Respond to % #v with headers % #v", pretty.Formatter(r), pretty.Formatter(w.Header()))

	if r.status > 199 && r.status < 504 {
		w.WriteHeader(r.status)
	} else {
		w.WriteHeader(200)
	}

	w.Write(r.body)
}
