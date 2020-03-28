// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"encoding/csv"
	"io"
	"sort"

	"github.com/garfieldius/t3ll/labels"
)

func writeCsv(src *labels.Labels, to io.Writer) error {
	w := csv.NewWriter(to)
	codes := make([]string, 0)

	for _, lang := range src.Languages {
		if lang != "en" {
			codes = append(codes, lang)
		}
	}

	sort.Strings(codes)
	codes = append([]string{"en"}, codes...)
	w.Write(append([]string{"key"}, codes...))

	for _, label := range src.Data {
		row := []string{label.ID}

		for _, c := range codes {
			for _, t := range label.Translations {
				if t.Language == c {
					row = append(row, t.Content)
				}
			}
		}

		err := w.Write(row)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func readCsv(from io.Reader, data *labels.Labels, mode string) (*labels.Labels, error) {
	r := csv.NewReader(from)
	newData := new(labels.Labels)
	newData.FromFile = data.FromFile
	newData.Type = data.Type
	newData.Data = make([]*labels.Label, 0)
	newData.SrcXlif = data.SrcXlif
	newData.SrcLegacy = data.SrcLegacy

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		if len(newData.Languages) == 0 {
			newData.Languages = row[1:]
		} else {
			l := new(labels.Label)
			l.ID = row[0]
			l.Translations = make([]*labels.Translation, 0, len(newData.Languages))

			for i, c := range row[1:] {
				if c == "" {
					continue
				}
				l.Translations = append(l.Translations, &labels.Translation{
					Language: newData.Languages[i],
					Content:  c,
				})
			}

			newData.Data = append(newData.Data, l)
		}
	}

	if mode != "replace" {
		newData = mergeLabels(data, newData)
	}

	if err := newData.Save(); err != nil {
		return nil, err
	}

	return data, nil
}

func mergeLabels(a, b *labels.Labels) *labels.Labels {
	res := &labels.Labels{
		Languages: []string{"en"},
		Data:      make([]*labels.Label, 0),
		Type:      a.Type,
		FromFile:  a.FromFile,
	}

	if b.Type != "" {
		res.Type = b.Type
	}

	if b.FromFile != "" {
		res.FromFile = b.FromFile
	}

	for _, l := range a.Languages {
		if l != "en" {
			res.Languages = append(res.Languages, l)
		}
	}

	for _, l := range b.Languages {
		found := false

		for _, la := range a.Languages {
			if la == l {
				found = true
				break
			}
		}

		if !found {
			res.Languages = append(res.Languages, l)
		}
	}

	for _, la := range a.Data {
		ln := &labels.Label{
			ID:           la.ID,
			Translations: make([]*labels.Translation, 0),
		}

		for _, ta := range la.Translations {
			tn := &labels.Translation{
				Language: ta.Language,
				Content:  ta.Content,
			}

			ln.Translations = append(ln.Translations, tn)
		}

		res.Data = append(res.Data, ln)
	}

	for i, lb := range b.Data {
		var ln *labels.Label

		for _, la := range res.Data {
			if la.ID == lb.ID {
				ln = lb
				break
			}
		}

		if ln == nil {
			ln = &labels.Label{
				ID:           lb.ID,
				Translations: make([]*labels.Translation, 0),
			}

			old := res.Data
			res.Data = make([]*labels.Label, 0)
			k := 0
			added := false

			for j := 0; j < len(old)+1; j++ {
				if j == i {
					res.Data = append(res.Data, ln)
					added = true
				} else {
					res.Data = append(res.Data, old[k])
					k++
				}
			}

			if !added {
				res.Data = append(res.Data, ln)
			}
		}

		for _, tb := range lb.Translations {
			found := false
			for _, tn := range ln.Translations {
				if tn.Language == tb.Language {
					tn.Content = tb.Content
					found = true
					break
				}
			}

			if !found {
				ln.Translations = append(ln.Translations, &labels.Translation{
					Language: tb.Language,
					Content:  tb.Content,
				})
			}
		}
	}

	return res
}
