package file

/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
		Data:     make(map[string]map[string]string),
	}

	for _, file := range x.Files {
		langCode := file.File.ToLang

		if langCode == "" {
			langCode = file.File.SrcLang
		}

		codes[langCode] = true

		for _, unit := range file.File.Body.Units {
			l := unit.To

			if l == "" {
				l = unit.Src
			}

			if _, ok := labels.Data[unit.Id]; !ok {
				labels.Data[unit.Id] = make(map[string]string)
			}

			labels.Data[unit.Id][langCode] = l
		}
	}

	for langCode := range codes {
		labels.Langs = append(labels.Langs, langCode)
	}

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
	Type   string `xml:"type,omitempty"`
	Header string `xml:"header,omitempty"`
	Desc   string `xml:"description,omitempty"`
	Gen    string `xml:"generator,omitempty"`
	Mod    string `xml:"module,omitempty"`
}

type XliffBody struct {
	Units []*XliffUnit `xml:"trans-unit,omitempty"`
}

type XliffUnit struct {
	Id  string `xml:"id,attr,omitempty"`
	Src string `xml:"source,omitempty"`
	To  string `xml:"target,omitempty"`
}
