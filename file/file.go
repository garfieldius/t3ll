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
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/garfieldius/t3ll/log"
	"github.com/kr/pretty"
)

type XmlLayout string

const (
	XmlXliff  XmlLayout = "xlf"
	XmlLegacy XmlLayout = "xml"
)

var (
	xliffLangPrefix = regexp.MustCompile(`^[a-z]{2,3}\.`)
	from            LangFile
)

func New(name string) (*Labels, error) {
	l := Labels{
		FromFile: name,
		Langs:    []string{"en"},
		Data: []*Label{{
			Id: "new.1",
			Trans: []*Translation{{
				Lng:     "en",
				Content: "",
			}},
		}},
	}

	switch {
	case strings.HasSuffix(name, ".xml"):
		l.Type = XmlLegacy
		log.Msg("Using legacy XML for %s", name)
		break
	case strings.HasSuffix(name, ".xlf") || strings.HasSuffix(name, ".xlif") || strings.HasSuffix(name, ".xliff"):
		l.Type = XmlXliff
		log.Msg("Using XLIF for %s", name)

		base := path.Base(name)
		if xliffLangPrefix.MatchString(base) {
			lang := base[0:strings.Index(base, ".")]
			l.Langs = append(l.Langs, lang)
			l.Data[0].Trans = append(l.Data[0].Trans, &Translation{Lng: lang, Content: ""})
		}
		break

	default:
		return nil, errors.New("Invalid file suffix")
	}

	return &l, nil
}

func Open(src string) (*Labels, error) {
	if len(src) < 4 {
		return nil, errors.New("Filename cannot have less than 4 chars")
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
		tree.Src = abs
		from = tree
		log.Msg("Unmarshaled %s into %# v", abs, pretty.Formatter(tree))
		return tree.Labels(), nil

	case strings.HasSuffix(abs, ".xlf") || strings.HasSuffix(abs, ".xlif") || strings.HasSuffix(abs, ".xliff"):
		xlif := new(XliffRoot)
		err = xml.Unmarshal(data, xlif)

		if err != nil {
			return nil, err
		}

		name := filepath.Base(abs)
		xlif.Src = abs
		all := &Xliff{
			StartSrc: abs,
			Files:    []*XliffRoot{xlif},
		}

		if xliffLangPrefix.MatchString(name) {
			xlif.Lang = name[0:strings.Index(name, ".")]
			name = name[strings.Index(name, ".")+1:]
		} else {
			xlif.Lang = "en"
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
				xlif.Lang = n[0:strings.Index(n, ".")]
			} else {
				xlif.Lang = "en"
			}

			xlif.Src = targetPath

			all.Files = append(all.Files, xlif)
		}

		log.Msg("Marshalled Xlif into %# v", pretty.Formatter(all))
		from = all
		return all.Labels(), nil
	}

	return nil, errors.New("Cannot read file " + src)
}

type LangFile interface {
	Labels() *Labels
}

type Labels struct {
	Type     XmlLayout `json:"format"`
	FromFile string    `json:"-"`
	Langs    []string  `json:"languages"`
	Data     []*Label  `json:"labels"`
}

type Label struct {
	Id    string         `json:"id"`
	Trans []*Translation `json:"trans"`
}

type Translation struct {
	Content string `json:"content"`
	Lng     string `json:"lng"`
}
