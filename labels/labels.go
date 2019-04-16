// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

// Package labels offes types and functions to load and save XLIFF v1 and
// TYPO3 locallang XML files
package labels

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kr/pretty"

	"github.com/garfieldius/t3ll/log"
)

// XMLType sets the type of XML schema a file has
type XMLType string

const (
	// XMLXliffv1 is the new XLIF schema
	XMLXliffv1 XMLType = "xlf"
	// XMLLegacy is the old TYPO3 schema
	XMLLegacy XMLType = "xml"
)

var xliffLangPrefix = regexp.MustCompile(`^[a-z]{2,3}\.`)

// New create a new Labels object for the given file
// That file should not exist yet, because its content will
// be overwritten upon the first save
func New(name string) (*Labels, error) {
	l := Labels{
		FromFile:  name,
		Languages: []string{"en"},
		Data: []*Label{{
			ID: "new.1",
			Translations: []*Translation{{
				Language: "en",
				Content:  "",
			}},
		}},
	}

	switch {
	case strings.HasSuffix(name, ".xml"):
		l.Type = XMLLegacy
		log.Msg("Using legacy XML for %s", name)
		break
	case strings.HasSuffix(name, ".xlf") || strings.HasSuffix(name, ".xlif") || strings.HasSuffix(name, ".xliff"):
		l.Type = XMLXliffv1
		log.Msg("Using XLIF for %s", name)
		base := path.Base(name)
		if xliffLangPrefix.MatchString(base) {
			lang := base[0:strings.Index(base, ".")]
			l.Languages = append(l.Languages, lang)
			l.Data[0].Translations = append(l.Data[0].Translations, &Translation{Language: lang, Content: ""})
		}
		break
	default:
		return nil, errors.New("invalid file suffix")
	}

	return &l, nil
}

// Open creates a new Labels object with the content of the given file
func Open(src string) (*Labels, error) {
	if len(src) < 4 {
		return nil, errors.New("filename cannot have less than 4 chars")
	}

	abs, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(abs)
	if err != nil {
		log.Msg("Cannot stat %s: %s", abs, err)
		log.Msg("Naively assuming file does not exist and create one")
		return New(abs)
	}

	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasSuffix(abs, ".xml"):
		tree := new(T3Root)
		err = xml.Unmarshal(data, tree)

		if err != nil {
			return nil, err
		}
		tree.SourceFile = abs
		log.Msg("Unmarshaled %s into %# v", abs, pretty.Formatter(tree))
		return tree.Labels(), nil

	case strings.HasSuffix(abs, ".xlf") || strings.HasSuffix(abs, ".xlif") || strings.HasSuffix(abs, ".xliff"):
		xlif := new(XliffRoot)
		err = xml.Unmarshal(data, xlif)

		if err != nil {
			return nil, err
		}

		name := filepath.Base(abs)
		xlif.SourceFile = abs
		all := &Xliff{
			SourceFile: abs,
			Files:      []*XliffRoot{xlif},
		}

		if xliffLangPrefix.MatchString(name) {
			xlif.Language = name[0:strings.Index(name, ".")]
			name = name[strings.Index(name, ".")+1:]
		} else {
			xlif.Language = "en"
		}

		dir := filepath.Dir(abs)
		start := filepath.Base(abs)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, info := range files {
			targetPath := filepath.Join(dir, info.Name())

			if targetPath == abs {
				continue
			}

			if !strings.Contains(info.Name(), start) {
				continue
			}

			if info.IsDir() || !strings.HasSuffix(targetPath, name) {
				log.Msg("Ignoring entry %s", targetPath)
				continue
			}

			data, err := ioutil.ReadFile(targetPath)
			if err != nil {
				log.Msg("Cannot read file %s: %s", err)
				continue
			}

			xlif := new(XliffRoot)
			err = xml.Unmarshal(data, xlif)
			if err != nil {
				log.Err("Cannot unmarshal data of file %s: %s", targetPath, err)
				continue
			}

			n := filepath.Base(targetPath)
			if xliffLangPrefix.MatchString(n) {
				xlif.Language = n[0:strings.Index(n, ".")]
			} else {
				xlif.Language = "en"
			}

			xlif.SourceFile = targetPath

			all.Files = append(all.Files, xlif)
		}

		log.Msg("Marshalled Xlif into %# v", pretty.Formatter(all))
		return all.Labels(), nil
	}

	return nil, errors.New("Cannot read file " + src)
}

// LangFile describes a reader for Labels from a XML source
type LangFile interface {
	Labels() *Labels
}

// Labels is the root object of all translations
type Labels struct {
	Type      XMLType  `json:"format"`
	FromFile  string   `json:"-"`
	Languages []string `json:"languages"`
	Data      []*Label `json:"labels"`
}

// Label is a single label, containing one or more translations
type Label struct {
	ID           string         `json:"id"`
	Translations []*Translation `json:"trans"`
}

// Translation is the text of a label in the given language
type Translation struct {
	Content  string `json:"content"`
	Language string `json:"lng"`
}
