// +build debug

// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"io/ioutil"

	"github.com/garfieldius/t3ll/log"
)

var Dir string
var html []byte

func init() {
	log.Msg("Loading HTML from %s", Dir)
	data, err := ioutil.ReadFile(Dir + "/frontend/build/index.html")
	if err != nil {
		panic(err)
	}
	html = data
}
