// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package labels

import (
	"encoding/xml"

	"github.com/kr/pretty"

	"github.com/garfieldius/t3ll/log"
)

// T3Root is the root element of a legacy XML
type T3Root struct {
	XMLName    string  `xml:"T3locallang"`
	Meta       *T3Meta `xml:"meta,omitempty"`
	Data       *T3Data `xml:"data"`
	SourceFile string  `xml:"-"`
}

// T3Meta is the meta block of a legacy XML
type T3Meta struct {
	Type    string `xml:"type,attr,omitempty"`
	For     string `xml:"type,omitempty"`
	Content string `xml:"description,omitempty"`
}

// T3Data is the blocks of translations
type T3Data struct {
	Type      string    `xml:"type,attr,omitempty"`
	Languages []*T3Lang `xml:"languageKey,omitempty"`
}

// T3Lang is a list of labels for a single language
type T3Lang struct {
	Type     string     `xml:"type,attr,omitempty"`
	Language string     `xml:"index,attr,omitempty"`
	Labels   []*T3Label `xml:"label,omitempty"`
}

// T3Label is a singe label and key pair
type T3Label struct {
	Key     string `xml:"index,attr"`
	Content string `xml:",innerxml"`
}

// MarshalXML will marshal a labels content to a XML element
func (l *T3Label) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if l.Content != "" {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "index"}, Value: l.Key}}
		e.EncodeToken(start)
		e.EncodeToken(xml.CharData(l.Content))
		e.EncodeToken(xml.EndElement{Name: start.Name})
	}
	return nil
}

// Labels returns all Labels of a XML
func (t T3Root) Labels() *Labels {
	data := &Labels{
		Type:      XMLLegacy,
		FromFile:  t.SourceFile,
		Languages: make([]string, 0),
		Data:      make([]*Label, 0),
	}

	for _, lang := range t.Data.Languages {
		langKey := lang.Language
		if langKey == "default" {
			langKey = "en"
		}
		data.Languages = append(data.Languages, langKey)

		for _, label := range lang.Labels {
			found := false
			trans := &Translation{Language: langKey, Content: label.Content}

			for _, l := range data.Data {
				if l.ID == label.Key {
					l.Translations = append(l.Translations, trans)
					found = true
					break
				}
			}

			if !found {
				newTrans := &Label{ID: label.Key, Translations: []*Translation{trans}}
				data.Data = append(data.Data, newTrans)
			}
		}
	}

	log.Msg("Converted ll tree into %# v", pretty.Formatter(data))

	return data
}

// LocallangConverter translates a Labels structure to a T3Root structure
// that can be marshalled to XML
type LocallangConverter struct {
	src *Labels
}

// XML will convert the Labels in src into a LangFile
func (c *LocallangConverter) XML() LangFile {
	data := &T3Data{
		Languages: make([]*T3Lang, 0, len(c.src.Languages)),
	}

	for _, l := range c.src.Languages {
		lang := &T3Lang{
			Type:     "index",
			Language: l,
			Labels:   make([]*T3Label, 0, len(c.src.Data)),
		}

		for _, label := range c.src.Data {
			for _, trans := range label.Translations {
				if trans.Language == l {
					lang.Labels = append(lang.Labels, &T3Label{
						Key:     label.ID,
						Content: trans.Content,
					})
					break
				}
			}
		}

		if len(lang.Labels) > 0 {
			data.Languages = append(data.Languages, lang)
		}
	}

	res := &T3Root{Data: data}

	log.Msg("Converted data to %# v", pretty.Formatter(res))

	return res
}

// File returns the file name and path of the source
func (c *LocallangConverter) File() string {
	return c.src.FromFile
}
