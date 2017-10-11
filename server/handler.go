package server

/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/garfieldius/t3ll/file"
	"github.com/garfieldius/t3ll/log"
)

var (
	markerData  = []byte(`"{DATA}"`)
	markerTitle = []byte("{TITLE}")
	notFound    = []byte(`{"message":"Resource not found"}`)
	saveSuccess = []byte(`{"success":true,"message":"File saved successfully"}`)
	saveError   = []byte(`{"success":false,"message":"Error during save"}`)
	invalidCSV  = []byte(`{"success":false,"message":"Invalid CSV data"}`)
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		notFoundHandler(res, req)
		return
	}

	mu.Lock()
	initState, err := json.Marshal(data)
	mu.Unlock()

	if err != nil {
		log.Fatal("Cannot marshall data to JSON: %s", err)
	}

	filename := []byte(path.Base(data.FromFile))

	response := MustAsset("editor.html")
	response = bytes.Replace(response, markerData, initState, -1)
	response = bytes.Replace(response, markerTitle, filename, -1)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func quitHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/quit" {
		notFoundHandler(res, req)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte{})

	Stop()
	stop <- true
}

func saveHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/save" {
		notFoundHandler(res, req)
		return
	}

	newData := new(file.Labels)
	err := json.Unmarshal([]byte(req.FormValue("data")), newData)

	if err != nil {
		log.Err("Cannot unmarshal data: %s", err)
		quitWithError(res)
		return
	}

	newData.FromFile = data.FromFile
	newData.Type = data.Type

	if setType := req.FormValue("format"); setType == "xlif" {
		log.Msg("Converting to xliff")
		newData.Type = file.XmlXliff
		newData.FromFile = data.FromFile[:len(data.FromFile)-3] + "xlf"
	}

	err = file.Save(newData)

	if err != nil {
		log.Err("Cannot save data to file: %s", err)
		quitWithError(res)
		return
	}

	mu.Lock()
	data = newData
	mu.Unlock()

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusCreated)
	res.Write(saveSuccess)
}

func csvHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/csv" {
		notFoundHandler(res, req)
		return
	}

	switch req.Method {
	case "POST":
		mu.Lock()
		err := readCsv(req.Body, req.URL.Query().Get("mode"))
		mu.Unlock()
		if err != nil {
			quitWithBadRequest(res, invalidCSV)
		} else {
			res.Header().Set("Content-Type", "application/json;charset=UTF-8")
			res.Write([]byte(`{"success":true,"message":"CSV imported"}`))
		}
		break
	case "GET":
		buf := new(bytes.Buffer)
		mu.Lock()
		err := writeCsv(data, buf)
		mu.Unlock()
		if err != nil {
			quitWithBadRequest(res, invalidCSV)
		} else {
			res.Header().Set("Content-Type", "text/csv;charset=UTF-8")
			res.Header().Set("Content-Disposition", "attachment; filename=locallang.csv")
			io.Copy(res, buf)
		}
		break

	default:
		notFoundHandler(res, req)
	}
}

func quitWithBadRequest(res http.ResponseWriter, msg []byte) {
	res.WriteHeader(http.StatusBadRequest)
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.Write(msg)
}

func quitWithError(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusInternalServerError)
	res.Write(saveError)
}

func notFoundHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusNotFound)
	res.Write(notFound)
}
