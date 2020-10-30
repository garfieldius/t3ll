// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package labels

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/garfieldius/t3ll/log"
)

var (
	xmlStart = []byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>` + "\n")
)

// Save will marshal the labels to the given XML type and store its result to disk
func (l *Labels) Save() error {
	var lastErr error
	var jobs []converter

	switch l.Type {
	case XMLLegacy:
		log.Msg("Create legacy converter for %s", l.FromFile)
		jobs = []converter{&LocallangConverter{src: l}}
		break
	case XMLXliffv1:
		jobs = make([]converter, len(l.Languages))

		for i, lang := range l.Languages {
			c := &XliffConverter{
				src:  l,
				lang: lang,
			}

			if l.SrcXlif != nil && len(l.SrcXlif.Files) >= i+1 {
				c.file = l.SrcXlif.Files[i].File
			}

			jobs[i] = c
		}
		break
	}

	msgs := make(chan error, len(jobs))

	for _, j := range jobs {
		go doSave(j, msgs)
	}

	for i := 0; i < len(jobs); i++ {
		err := <-msgs
		if err != nil {
			lastErr = err
			log.Err("Error during save op: %s", err)
		}
	}

	return lastErr
}

func doSave(d converter, done chan error) {
	if langData := d.XML(); langData != nil {
		buf, err := xml.MarshalIndent(langData, "", langData.IndentChar())
		if err != nil {
			done <- err
			return
		}

		filename := d.File()
		if strings.HasSuffix(filename, ".xlf") {
			buf = append(xmlStart, buf...)
		}

		buf = append(buf, byte("\n"[0]))

		if werr := ioutil.WriteFile(filename, buf, 0644); werr != nil {
			done <- werr
			return
		}
	}

	done <- nil
}

type converter interface {
	XML() LangFile
	File() string
}
