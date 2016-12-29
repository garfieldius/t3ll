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
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		notFoundHandler(res, req)
		return
	}

	initState, err := json.Marshal(data)

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

	if setType := req.FormValue("format"); setType == "xlif" {
		newData.Type = file.XmlXliff
	} else {
		newData.Type = data.Type
	}

	newData.FromFile = data.FromFile
	err = file.Save(newData)

	if err != nil {
		log.Err("Cannot save data to file: %s", err)
		quitWithError(res)
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusCreated)
	res.Write(saveSuccess)
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
