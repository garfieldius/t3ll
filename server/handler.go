// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/garfieldius/t3ll/labels"
	"github.com/garfieldius/t3ll/log"
)

var (
	notFound    = []byte(`{"success":false,"message":"Resource not found"}`)
	saveSuccess = []byte(`{"success":true,"message":"File saved successfully"}`)
	saveError   = []byte(`{"success":false,"message":"Error during save"}`)
	invalidCSV  = []byte(`{"success":false,"message":"Invalid CSV data"}`)
	heartbeat   = []byte(`{"success":true,"message":"OK"}`)
)

type handler struct {
	state *labels.Labels
	mu    *sync.Mutex
	quitSig chan struct{}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	d := response{
		status: 404,
		body:   notFound,
	}

	switch {
	case r.URL.Path == "/hb":
		d.body = heartbeat
		d.status = 200
		break
	case r.Method == "GET" && r.URL.Path == "/":
		d.body = html
		d.ctype = "text/html"
		break
	case r.Method == "GET" && r.URL.Path == "/quit":
		log.Msg("Received quit signal")
		h.quitSig <- struct{}{}
		break
	case r.Method == "GET" && r.URL.Path == "/csv":
		buf := new(bytes.Buffer)
		err := writeCsv(h.state, buf)
		if err != nil {
			d.status = 500
			d.body = invalidCSV
		} else {
			d.ctype = "text/csv"
			d.body = buf.Bytes()
		}
		break
	case r.Method == "POST" && r.URL.Path == "/csv":
		newState, err := readCsv(r.Body, h.state, r.URL.Query().Get("mode"))
		if err != nil {
			d.status = 400
			d.body = invalidCSV
		} else {
			h.state = newState
			d.body = saveSuccess
		}
		break
	case r.Method == "GET" && r.URL.Path == "/data":
		d.body, _ = json.Marshal(h.state)
		break
	case r.Method == "POST" && r.URL.Path == "/save":
		src := []byte(r.FormValue("data"))
		newState, err := saveHandler(src, r.FormValue("format"), h.state.FromFile, h.state.Type)
		if err != nil {
			log.Err("Cannot save data: %s", err)
			d.body = saveError
			d.status = 400
		} else {
			h.state = newState
			d.body = saveSuccess
		}
		break
	}

	d.write(w)
}
