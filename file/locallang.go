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

	"github.com/garfieldius/t3ll/log"
	"github.com/kr/pretty"
)

type T3Root struct {
	XMLName string  `xml:"T3locallang"`
	Meta    *T3Meta `xml:"meta,omitempty"`
	Data    *T3Data `xml:"data"`
	Src     string  `xml:"-"`
}

type T3Meta struct {
	Typ  string `xml:"type,attr,omitempty"`
	For  string `xml:"type,omitempty"`
	Desc string `xml:"description,omitempty"`
}

type T3Data struct {
	Typ   string    `xml:"type,attr,omitempty"`
	Langs []*T3Lang `xml:"languageKey,omitempty"`
}

type T3Lang struct {
	Typ    string     `xml:"type,attr,omitempty"`
	Lang   string     `xml:"index,attr,omitempty"`
	Labels []*T3Label `xml:"label,omitempty"`
}

type T3Label struct {
	Key string `xml:"index,attr"`
	Cnt string `xml:",innerxml"`
}

func (l *T3Label) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if l.Cnt != "" {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "index"}, Value: l.Key}}
		e.EncodeToken(start)
		e.EncodeToken(xml.CharData(l.Cnt))
		e.EncodeToken(xml.EndElement{Name: start.Name})
	}
	return nil
}

func (t T3Root) Labels() *Labels {
	data := &Labels{
		Type:     XmlLegacy,
		FromFile: t.Src,
		Langs:    make([]string, 0),
		Data:     make([]*Label, 0),
	}

	for _, lang := range t.Data.Langs {
		langKey := lang.Lang

		if langKey == "default" {
			langKey = "en"
		}

		data.Langs = append(data.Langs, langKey)

		for _, label := range lang.Labels {
			found := false
			trans := &Translation{Lng: langKey, Content: label.Cnt}

			for _, l := range data.Data {
				if l.Id == label.Key {
					l.Trans = append(l.Trans, trans)
					found = true
					break
				}
			}

			if !found {
				newTrans := &Label{Id: label.Key, Trans: []*Translation{trans}}
				data.Data = append(data.Data, newTrans)
			}
		}
	}

	log.Msg("Converted ll tree into %# v", pretty.Formatter(data))

	return data
}
