// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package labels

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/copier"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

// Xliff is the root of all languages and labels
type Xliff struct {
	SourceFile string
	Files      []*XliffRoot
}

// Labels converts a Xliff tree into a Labels tree
func (x *Xliff) Labels() *Labels {
	codes := make(map[string]bool)
	labels := &Labels{
		Type:      XMLXliffv1,
		SrcXlif:   x,
		File:      extPathOfFile(x.SourceFile),
		FromFile:  x.SourceFile,
		Languages: make([]string, 0),
		Data:      make([]*Label, 0),
	}

	for _, file := range x.Files {
		langCode := file.File.ToLang

		if langCode == "" {
			langCode = file.File.SrcLang
		}

		codes[langCode] = true

		for _, unit := range file.File.Body.Units {
			t := &Translation{
				Language: langCode,
				Content:  unit.To,
				Approved: unit.Approved,
			}

			if langCode == "en" {
				t.Content = unit.Src
				labels.File = extPathOfFile(file.SourceFile)
			}

			found := false

			for _, label := range labels.Data {
				if label.ID == unit.ID {
					found = true
					label.Translations = append(label.Translations, t)
				}
			}

			if !found {
				labels.Data = append(labels.Data, &Label{
					ID:           unit.ID,
					Translations: []*Translation{t},
				})
			}
		}
	}

	for langCode := range codes {
		labels.Languages = append(labels.Languages, langCode)
	}

	log.Debugf("Converted xlif into %# v", pretty.Formatter(labels))

	return labels
}

// XliffRoot is the virtual root node of a XML document
type XliffRoot struct {
	XMLName    string     `xml:"xliff"`
	Version    string     `xml:"version,attr,omitempty"`
	XMLNS      string     `xml:"xmlns,attr,omitempty"`
	File       *XliffFile `xml:"file,omitempty"`
	Language   string     `xml:"-"`
	SourceFile string     `xml:"-"`
}

// Labels in this implementation has no function, because the
// calling party handles the conversion
func (x *XliffRoot) Labels() *Labels {
	return nil // Noop
}

// IndentChar determines the string for indentation or
// XML tags
func (x *XliffRoot) IndentChar() string {
	return indentOfFile(x.SourceFile)
}

// XliffFile is the actual root node of a Xliff document
type XliffFile struct {
	SrcLang  string       `xml:"source-language,attr"`
	ToLang   string       `xml:"target-language,attr,omitempty"`
	Orig     string       `xml:"original,attr"`
	Datatype string       `xml:"datatype,attr"`
	Name     string       `xml:"product-name,attr,omitempty"`
	Date     string       `xml:"date,attr,omitempty"`
	Header   *XliffHeader `xml:"header,omitempty"`
	Body     *XliffBody   `xml:"body,omitempty"`
}

// XliffHeader is a collection of common metadata
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

// XliffBody contains the units of a Xliff document
type XliffBody struct {
	Units []*XliffUnit `xml:"trans-unit,omitempty"`
}

// XliffUnit is a single ID + Content unit of a Xliff document
type XliffUnit struct {
	ID       string `xml:"id,attr,omitempty"`
	Approved string `xml:"approved,attr,omitempty"`
	Src      string `xml:"source"`
	To       string `xml:"target,omitempty"`
}

// XliffConverter will convert a Xliff structure to struct Labels
type XliffConverter struct {
	src  *Labels
	file *XliffFile
	lang string
}

// XML will return the Xliff structure of a Labels structure
func (c *XliffConverter) XML() LangFile {
	b := &XliffBody{
		Units: make([]*XliffUnit, 0, len(c.src.Data)),
	}
	f := &XliffFile{
		SrcLang:  "en",
		Date:     time.Now().Format(time.RFC3339),
		Body:     b,
		Datatype: "plaintext",
	}

	if c.file != nil {
		if c.file.Datatype != "" {
			f.Datatype = c.file.Datatype
		}
		if c.file.Orig != "" {
			f.Orig = c.file.Orig
		}
		f.Header = c.file.Header
		f.Name = c.file.Name
	}

	l := "en"

	if c.lang != "en" {
		f.ToLang = c.lang
		l = c.lang
	}

	for _, label := range c.src.Data {
		for _, trans := range label.Translations {
			if trans.Language == l {
				a := trans.Approved

				if a != "no" && a != "yes" {
					a = ""
				}

				u := &XliffUnit{
					ID:       label.ID,
					Approved: a,
				}

				if c.lang == "en" {
					u.Src = trans.Content
				} else {
					if trans.Content == "" {
						continue
					}
					u.To = trans.Content
					for _, orig := range label.Translations {
						if orig.Language == "en" {
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

	x := XliffRoot{
		SourceFile: c.src.FromFile,
		Language:   c.lang,
		Version:    "1.2",
		XMLNS:      "urn:oasis:names:tc:xliff:document:1.2",
		File:       f,
	}

	cp := XliffRoot{}

	if err := copier.Copy(&cp, &x); err == nil {
		if oldFileData, err := os.ReadFile(c.File()); err == nil {
			oldXlif := new(XliffRoot)
			if err = xml.Unmarshal(oldFileData, oldXlif); err == nil {
				if oldXlif.File != nil && oldXlif.File.Date != "" {
					cp.File.Date = oldXlif.File.Date
					oldXml, _ := xml.MarshalIndent(&oldXlif, "", "  ")
					newXml, _ := xml.MarshalIndent(&cp, "", "  ")
					if oldXml != nil && newXml != nil && bytes.Equal(oldXml, newXml) {
						return nil
					}
				}
			}
		}
	}

	return &x
}

// File determines the desired target file of this document
func (c *XliffConverter) File() string {
	d := filepath.Dir(c.src.FromFile)
	f := xliffLangPrefix.ReplaceAllString(filepath.Base(c.src.FromFile), "")
	p := ""

	if c.lang != "en" {
		p = c.lang + "."
	}

	full := d + string(filepath.Separator) + p + f
	log.Infof("Normalized target file path from %s to %s", c.src.FromFile, full)

	return full
}
