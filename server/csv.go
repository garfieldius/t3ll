package server

import (
	"encoding/csv"
	"io"
	"sort"

	"github.com/garfieldius/t3ll/file"
)

func writeCsv(src *file.Labels, to io.Writer) error {
	w := csv.NewWriter(to)
	codes := make([]string, 0)

	for _, lang := range src.Langs {
		if lang != "en" {
			codes = append(codes, lang)
		}
	}

	sort.Strings(codes)
	codes = append([]string{"en"}, codes...)

	w.Write(append([]string{"key"}, codes...))

	for _, label := range src.Data {
		row := []string{label.Id}

		for _, c := range codes {
			for _, t := range label.Trans {
				if t.Lng == c {
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

func readCsv(from io.Reader, mode string) error {
	r := csv.NewReader(from)
	newData := new(file.Labels)
	newData.FromFile = data.FromFile
	newData.Type = data.Type
	newData.Data = make([]*file.Label, 0)

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if len(newData.Langs) == 0 {
			newData.Langs = row[1:]
		} else {
			l := new(file.Label)
			l.Id = row[0]
			l.Trans = make([]*file.Translation, 0, len(newData.Langs))

			for i, c := range row[1:] {
				if c == "" {
					continue
				}
				l.Trans = append(l.Trans, &file.Translation{
					Lng:     newData.Langs[i],
					Content: c,
				})
			}

			newData.Data = append(newData.Data, l)
		}
	}

	if mode == "replace" {
		data = newData
	} else {
		data = mergeLabels(data, newData)
	}

	return file.Save(data)
}

func mergeLabels(a, b *file.Labels) *file.Labels {
	res := &file.Labels{
		Langs:    []string{"en"},
		Data:     make([]*file.Label, 0),
		Type:     a.Type,
		FromFile: a.FromFile,
	}

	if b.Type != "" {
		res.Type = b.Type
	}

	if b.FromFile != "" {
		res.FromFile = b.FromFile
	}

	for _, l := range a.Langs {
		if l != "en" {
			res.Langs = append(res.Langs, l)
		}
	}

	for _, l := range b.Langs {
		found := false

		for _, la := range a.Langs {
			if la == l {
				found = true
				break
			}
		}

		if !found {
			res.Langs = append(res.Langs, l)
		}
	}

	for _, la := range a.Data {
		ln := &file.Label{
			Id:    la.Id,
			Trans: make([]*file.Translation, 0),
		}

		for _, ta := range la.Trans {
			tn := &file.Translation{
				Lng:     ta.Lng,
				Content: ta.Content,
			}

			ln.Trans = append(ln.Trans, tn)
		}

		res.Data = append(res.Data, ln)
	}

	for i, lb := range b.Data {
		var ln *file.Label

		for _, la := range res.Data {
			if la.Id == lb.Id {
				ln = lb
				break
			}
		}

		if ln == nil {
			ln = &file.Label{
				Id:    lb.Id,
				Trans: make([]*file.Translation, 0),
			}

			old := res.Data
			res.Data = make([]*file.Label, len(old)+1)
			k := 0
			added := false

			for j := 0; j < len(old)+1; j++ {
				if j == i {
					res.Data[j] = ln
					added = true
				} else {
					res.Data[j] = old[k]
					k++
				}
			}

			if !added {
				res.Data = append(res.Data, ln)
			}
		}

		for _, tb := range lb.Trans {
			found := false
			for _, tn := range ln.Trans {
				if tn.Lng == tb.Lng {
					tn.Content = tb.Content
					found = true
					break
				}
			}

			if !found {
				ln.Trans = append(ln.Trans, &file.Translation{
					Lng:     tb.Lng,
					Content: tb.Content,
				})
			}
		}
	}

	return res
}
