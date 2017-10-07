package file

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
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/garfieldius/t3ll/log"
	"github.com/kr/pretty"
)

var (
	xmlStart = []byte(`<?xml version="1.0" encoding="utf-8" standalone="yes" ?>` + "\n")
)

func Save(l *Labels) error {
	var lastErr error
	var jobs []converter

	switch l.Type {
	case XmlLegacy:
		log.Msg("Create legacy converter for %s", l.FromFile)
		jobs = []converter{&LocallangConverter{src: l}}
		break
	case XmlXliff:
		jobs = make([]converter, len(l.Langs))

		for i, lang := range l.Langs {
			jobs[i] = &XliffConverter{
				src:  l,
				lang: lang,
			}
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
		} else {
			log.Msg("No error during save")
		}
	}

	return lastErr
}

type LocallangConverter struct {
	src *Labels
}

func (c *LocallangConverter) Xml() LangFile {
	data := &T3Data{
		Langs: make([]*T3Lang, 0, len(c.src.Langs)),
	}

	for _, l := range c.src.Langs {
		lang := &T3Lang{
			Typ:    "index",
			Lang:   l,
			Labels: make([]*T3Label, 0, len(c.src.Data)),
		}

		for _, label := range c.src.Data {
			for _, trans := range label.Trans {
				if trans.Lng == l {
					lang.Labels = append(lang.Labels, &T3Label{
						Key: label.Id,
						Cnt: trans.Content,
					})
					break
				}
			}
		}

		if len(lang.Labels) > 0 {
			data.Langs = append(data.Langs, lang)
		}
	}

	res := &T3Root{Data: data}

	if s, ok := from.(*T3Root); ok {
		res.Meta = s.Meta
	}

	log.Msg("Converted data to %# v", pretty.Formatter(res))

	return res
}

func (c *LocallangConverter) File() string {
	return c.src.FromFile
}

type XliffConverter struct {
	src  *Labels
	lang string
}

func (c *XliffConverter) Xml() LangFile {
	b := &XliffBody{
		Units: make([]*XliffUnit, 0, len(c.src.Data)),
	}
	f := &XliffFile{
		SrcLang: "en",
		Date:    time.Now().Format(time.RFC3339),
		Body:    b,
	}
	l := "en"

	if c.lang != "en" {
		f.ToLang = c.lang
		l = c.lang
	}

	for _, label := range c.src.Data {
		for _, trans := range label.Trans {
			if trans.Lng == l {
				u := &XliffUnit{
					Id: label.Id,
				}

				if c.lang == "en" {
					u.Src = trans.Content
				} else {
					if trans.Content == "" {
						continue
					}
					u.To = trans.Content
					for _, orig := range label.Trans {
						if orig.Lng == "en" {
							u.Src = orig.Content
						}
					}
				}

				b.Units = append(b.Units, u)
			}
		}
	}

	if len(b.Units) < 1 {
		return nil
	}

	x := &XliffRoot{
		Src:     c.src.FromFile,
		Lang:    c.lang,
		Version: "1.0",
		File:    f,
	}

	if orig, ok := from.(*Xliff); ok {
		if orig.Files[0].File.Header != nil {
			f.Header = orig.Files[0].File.Header
		}

		if orig.Files[0].File.Name != "" {
			f.Name = orig.Files[0].File.Name
		}

		if orig.Files[0].File.Orig != "" {
			f.Orig = orig.Files[0].File.Orig
		}
	}

	return x
}

func (c *XliffConverter) File() string {
	d := filepath.Dir(c.src.FromFile)
	f := xliffLangPrefix.ReplaceAllString(filepath.Base(c.src.FromFile), "")
	p := ""

	if c.lang != "en" {
		p = c.lang + "."
	}

	full := d + string(filepath.Separator) + p + f
	log.Msg("Nomalized from %s to %s", c.src.FromFile, full)

	return full
}

func doSave(d converter, done chan error) {
	if langData := d.Xml(); langData != nil {
		buf, err := xml.MarshalIndent(langData, "", "    ")
		if err != nil {
			done <- err
			return
		}

		filename := d.File()
		if strings.HasSuffix(filename, ".xlf") {
			buf = append(xmlStart, buf...)
		}

		if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
			done <- err
			return
		}
	}

	done <- nil
}

type converter interface {
	Xml() LangFile
	File() string
}
