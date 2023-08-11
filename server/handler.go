// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/garfieldius/t3ll/labels"

	log "github.com/sirupsen/logrus"
)

var (
	notFound    = []byte(`{"success":false,"message":"Resource not found"}`)
	saveSuccess = []byte(`{"success":true,"message":"File saved successfully"}`)
	saveError   = []byte(`{"success":false,"message":"Error during save"}`)
	invalidCSV  = []byte(`{"success":false,"message":"Invalid CSV data"}`)
)

//go:embed index.html
var html []byte

type handlerState struct {
	mu      *sync.Mutex
	state   *labels.Labels
	quitSig chan struct{}
}

type handler struct {
	st *handlerState
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.st.mu.Lock()
	defer h.st.mu.Unlock()

	d := response{
		status: 404,
		body:   notFound,
	}

	switch {
	case r.Method == "GET" && r.URL.Path == "/":
		d.body = html
		d.status = 200
		d.ctype = "text/html"
		break
	case r.Method == "GET" && r.URL.Path == "/csv":
		buf := new(bytes.Buffer)
		err := writeCsv(h.st.state, buf)
		if err != nil {
			d.status = 500
			d.body = invalidCSV
		} else {
			name := strings.Replace(h.st.state.File, "\\", "/", -1)

			if strings.Contains(name, "/") {
				name = filepath.Base(name)
			}

			name = strings.TrimSuffix(name, filepath.Ext(name)) + ".csv"

			d.status = 200
			d.ctype = "text/csv"
			d.dlName = name
			d.body = buf.Bytes()
			log.Infof("Sending CSV as %s", d.dlName)
		}
		break
	case r.Method == "POST" && r.URL.Path == "/csv":
		newState, err := readCsv(r.Body, h.st.state, r.URL.Query().Get("mode"))
		if err != nil {
			d.status = 400
			d.body = invalidCSV
		} else {
			d.status = 200
			h.st.state = newState
			d.body = saveSuccess
		}
		break
	case r.Method == "GET" && r.URL.Path == "/data":
		d.body, _ = json.Marshal(h.st.state)
		d.status = 200
		break
	case r.Method == "POST" && r.URL.Path == "/save":
		src := []byte(r.FormValue("data"))
		newState, err := saveHandler(src, r.FormValue("format"), h.st.state)
		if err != nil {
			log.Errorf("Cannot save data: %s", err)
			d.body = saveError
			d.status = 400
		} else {
			d.status = 200
			h.st.state = newState
			d.body = saveSuccess
		}
		break
	}

	d.write(w)
}
