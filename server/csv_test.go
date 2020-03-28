// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"bytes"
	"testing"

	"github.com/kr/pretty"

	"github.com/garfieldius/t3ll/labels"
)

func TestCsvExport(t *testing.T) {
	l := labels.Labels{
		Languages: []string{"fr", "en", "it"},
		Data: []*labels.Label{
			&labels.Label{
				ID: "l1",
				Translations: []*labels.Translation{
					&labels.Translation{
						Language: "en",
						Content:  "l1.en",
					},
					&labels.Translation{
						Language: "fr",
						Content:  "l1.fr",
					},
					&labels.Translation{
						Language: "it",
						Content:  "l1.it",
					},
				},
			},
			&labels.Label{
				ID: "l2",
				Translations: []*labels.Translation{
					&labels.Translation{
						Language: "it",
						Content:  "l2.it",
					},
					&labels.Translation{
						Language: "fr",
						Content:  "l2.fr",
					},
					&labels.Translation{
						Language: "en",
						Content:  "l2.en",
					},
				},
			},
		},
	}

	buf := bytes.NewBuffer(nil)

	err := writeCsv(&l, buf)

	if err != nil {
		t.Errorf("Cannot write CSV: %s", err)
	}

	expected := `key,en,fr,it
l1,l1.en,l1.fr,l1.it
l2,l2.en,l2.fr,l2.it
`
	actual := buf.String()

	if actual != expected {
		t.Errorf("Invalid result in CSV Writer, expected\n%s\n\nbut was\n%s", expected, actual)
	}
}

func TestCsvDataReplace(t *testing.T) {
	src := bytes.NewBufferString(`key,en,fr,it
l1,l1.en,l1.fr,l1.it
l2,l2.en,l2.fr,l2.it
`)

	data, err := readCsv(src, &labels.Labels{}, "replace")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if data == nil {
		t.Error("Expected data to be set, but is nil")
	}

	if data.Languages == nil || len(data.Languages) != 3 {
		t.Errorf("Invalid number of languages: %# v", pretty.Formatter(data.Languages))
	}

	if data.Data == nil || len(data.Data) != 2 {
		t.Errorf("Invalid number of labels: %# v", pretty.Formatter(data.Data))
	}
}

func TestCsvDataMerge(t *testing.T) {
	existing := bytes.NewBufferString(`key,en,fr,it
l1,l1.en,l1.fr,l1.it
l2,l2.en,l2.fr,l2.it
l4,l4.en,l4.fr,l4.it
`)

	data, err := readCsv(existing, &labels.Labels{}, "replace")
	if err != nil {
		t.Errorf("Cannot set existing labels: %s", err)
	}

	additional := bytes.NewBufferString(`key,en,fr,it
l1,main.en,main.fr,main.it
l2,l2.en,l2.fr,l2.it
l3,l3.en,l3.fr,l3.it
`)

	data, err = readCsv(additional, data, "")
	if err != nil {
		t.Errorf("Cannot merge new labels labels: %s", err)
	}

	expected := `key,en,fr,it
l1,main.en,main.fr,main.it
l2,l2.en,l2.fr,l2.it
l3,l3.en,l3.fr,l3.it
l4,l4.en,l4.fr,l4.it
`

	buf := bytes.NewBuffer(nil)
	err = writeCsv(data, buf)
	actual := buf.String()

	if actual != expected {
		t.Errorf("Invalid result in CSV Merge, expected\n%s\n\nbut was\n%s", expected, actual)
	}
}
