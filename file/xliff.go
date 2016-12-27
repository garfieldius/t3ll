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
	"github.com/garfieldius/t3ll/log"
	"github.com/kr/pretty"
)

type Xliff struct {
	StartSrc string
	Files    []*XliffRoot
}

func (x *Xliff) Labels() *Labels {
	codes := make(map[string]bool)
	labels := &Labels{
		Type:     XmlXliff,
		FromFile: x.StartSrc,
		Langs:    make([]string, 0),
		Data:     make([]*Label, 0),
	}

	for _, file := range x.Files {
		langCode := file.File.ToLang

		if langCode == "" {
			langCode = file.File.SrcLang
		}

		codes[langCode] = true

		for _, unit := range file.File.Body.Units {
			t := &Translation{
				Lng:     langCode,
				Content: unit.To,
			}

			if langCode == "en" {
				t.Content = unit.Src
			}

			found := false

			for _, label := range labels.Data {
				if label.Id == unit.Id {
					found = true
					label.Trans = append(label.Trans, t)
				}
			}

			if !found {
				labels.Data = append(labels.Data, &Label{
					Id:    unit.Id,
					Trans: []*Translation{t},
				})
			}
		}
	}

	for langCode := range codes {
		labels.Langs = append(labels.Langs, langCode)
	}

	log.Msg("Converted xlif into %# v", pretty.Formatter(labels))

	return labels
}

type XliffRoot struct {
	XMLName string     `xml:"xliff"`
	Version string     `xml:"version,attr,omitempty"`
	File    *XliffFile `xml:"file,omitempty"`
	Lang    string     `xml:"-"`
	Src     string     `xml:"-"`
}

func (x *XliffRoot) Labels() *Labels {
	return nil // Noop
}

type XliffFile struct {
	SrcLang string       `xml:"source-language,attr,omitempty"`
	ToLang  string       `xml:"target-language,attr,omitempty"`
	Orig    string       `xml:"original,attr,omitempty"`
	Name    string       `xml:"product-name,attr,omitempty"`
	Date    string       `xml:"date,attr,omitempty"`
	Header  *XliffHeader `xml:"header,omitempty"`
	Body    *XliffBody   `xml:"body,omitempty"`
}

type XliffHeader struct {
	Skl      string `xml:"skl,omitempty"`
	PhaseGrp string `xml:"phase-group,omitempty"`
	Glossary string `xml:"glossary,omitempty"`
	Ref      string `xml:"reference,omitempty"`
	CountGrp string `xml:"count-group,omitempty"`
	Tool     string `xml:"tool,omitempty"`
	PopGrp   string `xml:"pop-group,omitempty"`
	Note     string `xml:"note,omitempty"`

	// Non-standard, but common in TYPO3
	Type        string `xml:"type,omitempty"`
	Description string `xml:"description,omitempty"`
	AuthName    string `xml:"authorName,omitempty"`
	AuthMail    string `xml:"authorEmail,omitempty"`
	AuthCompany string `xml:"authorCompany,omitempty"`
	Generator   string `xml:"generator,omitempty"`
}

type XliffBody struct {
	Units []*XliffUnit `xml:"trans-unit,omitempty"`
}

type XliffUnit struct {
	Id  string `xml:"id,attr,omitempty"`
	Src string `xml:"source,omitempty"`
	To  string `xml:"target,omitempty"`
}
