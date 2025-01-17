// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

// Package labels offes types and functions to load and save XLIFF v1 and
// TYPO3 locallang XML files
package labels

import (
	"encoding/xml"
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kr/pretty"

	log "github.com/sirupsen/logrus"
)

// XMLType sets the type of XML schema a file has
type XMLType string

const (
	// XMLXliffv1 is the new XLIF schema
	XMLXliffv1 XMLType = "xlf"
	// XMLLegacy is the old TYPO3 schema
	XMLLegacy XMLType = "xml"
)

var xliffLangPrefix = regexp.MustCompile(`^(af|ar|bs|bg|ca|ch|cs|cy|da|de|el|eo|es|et|eu|fa|fi|fo|fr|fr_CA|gl|he|hi|hr|hu|is|it|ja|ka|kl|km|ko|lb|lt|lv|mi|mk|ms|nl|no|nb|pl|pt|pt_BR|ro|ru|rw|sk|sl|sn|sq|sr|sv|th|tr|uk|vi|zh|zh_CN|zh_HK|zh_Hans_CN)\.`)

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
				Approved: "",
				Content:  "",
			}},
		}},
	}

	switch {
	case strings.HasSuffix(name, ".xml"):
		l.Type = XMLLegacy
		log.Infof("Using legacy XML for %s", name)
		break
	case strings.HasSuffix(name, ".xlf") || strings.HasSuffix(name, ".xlif") || strings.HasSuffix(name, ".xliff"):
		l.Type = XMLXliffv1
		log.Infof("Using XLIF for %s", name)
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
		log.Warnf("Cannot stat %s: %s", abs, err)
		log.Infof("Naively assuming file does not exist and create one")
		return New(abs)
	}

	data, err := os.ReadFile(abs)
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
		log.Debugf("Unmarshalled %s into %# v", abs, pretty.Formatter(tree))
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
		} else if xlif.File.ToLang != "" {
			xlif.Language = xlif.File.ToLang
		} else {
			xlif.Language = xlif.File.SrcLang
		}

		dir := filepath.Dir(abs)
		start := filepath.Base(abs)
		files, err := os.ReadDir(dir)
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
				log.Infof("Ignoring entry %s", targetPath)
				continue
			}

			data, err := os.ReadFile(targetPath)
			if err != nil {
				log.Warnf("Cannot read file %s: %s", err)
				continue
			}

			xlif := new(XliffRoot)
			err = xml.Unmarshal(data, xlif)
			if err != nil {
				log.Errorf("Cannot unmarshal data of file %s: %s", targetPath, err)
				continue
			}

			n := filepath.Base(targetPath)
			if xliffLangPrefix.MatchString(n) {
				xlif.Language = n[0:strings.Index(n, ".")]
			} else if xlif.File.ToLang != "" {
				xlif.Language = xlif.File.ToLang
			} else {
				xlif.Language = xlif.File.SrcLang
			}

			xlif.SourceFile = targetPath

			all.Files = append(all.Files, xlif)
		}

		log.Debugf("Marshalled Xlif into %# v", pretty.Formatter(all))
		return all.Labels(), nil
	}

	return nil, errors.New("Cannot read file " + src)
}

var indentTest = regexp.MustCompile("\n[ \t]+<")
var indentClean = regexp.MustCompile(`[^\t ]+`)

// indentOfFile checks for the indentation of the first tag
func indentOfFile(filename string) string {
	if data, err := os.ReadFile(filename); err == nil {
		for _, match := range indentTest.FindAll(data, -1) {
			if len(match) > 2 {
				return indentClean.ReplaceAllString(string(match), "")
			}
		}
	}
	return "	"
}

var winRootTest = regexp.MustCompile(`^[a-zA-Z]:(\\+)?$`)

func extPathOfFile(file string) string {
	remainder := file
	previous := ""

	for len(remainder) > 4 && !winRootTest.MatchString(remainder) {
		previous = remainder
		remainder = filepath.Dir(remainder)
		if previous == remainder {
			log.Infof("Reached root, not checking further")
			break
		}

		files := []string{"ext_emconf.php", "composer.json"}

		for _, f := range files {
			log.Infof("Checking for %s in %s", f, remainder)
			if stat, err := os.Stat(remainder + "/" + f); err == nil && stat != nil && !stat.IsDir() {
				parent := filepath.Dir(remainder)
				return "EXT:" + file[len(parent)+1:]
			}
		}
	}

	return filepath.Base(file)
}

// LangFile describes a reader for Labels from a XML source
type LangFile interface {
	Labels() *Labels
	IndentChar() string
}

// Labels is the root object of all translations
type Labels struct {
	File      string   `json:"documentTitle"`
	Type      XMLType  `json:"format"`
	FromFile  string   `json:"-"`
	SrcXlif   *Xliff   `json:"-"`
	SrcLegacy *T3Root  `json:"-"`
	Languages []string `json:"languages"`
	Data      []*Label `json:"labels"`
}

// CopyData creates a copy of all labels and translations
// keeping the (pointer) values in Data unchanged
func (l *Labels) CopyData() []*Label {
	c := make([]*Label, len(l.Data))

	for i, label := range l.Data {
		t := make([]*Translation, len(label.Translations))

		for j, trans := range label.Translations {
			t[j] = &Translation{
				Language: trans.Language,
				Content:  trans.Content,
				Approved: trans.Approved,
			}
		}

		c[i] = &Label{
			ID:           label.ID,
			Translations: t,
		}
	}

	return c
}

// Label is a single label, containing one or more translations
type Label struct {
	ID           string         `json:"id"`
	Translations []*Translation `json:"trans"`
}

// Translation is the text of a label in the given language
type Translation struct {
	Content  string `json:"content"`
	Approved string `json:"approved"`
	Language string `json:"lng"`
}
