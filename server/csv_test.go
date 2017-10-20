package server

import (
	"bytes"
	"testing"

	"github.com/garfieldius/t3ll/file"
	"github.com/kr/pretty"
)

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

func TestCsvExport(t *testing.T) {
	l := file.Labels{
		Langs: []string{"fr", "en", "it"},
		Data: []*file.Label{
			&file.Label{
				Id: "l1",
				Trans: []*file.Translation{
					&file.Translation{
						Lng:     "en",
						Content: "l1.en",
					},
					&file.Translation{
						Lng:     "fr",
						Content: "l1.fr",
					},
					&file.Translation{
						Lng:     "it",
						Content: "l1.it",
					},
				},
			},
			&file.Label{
				Id: "l2",
				Trans: []*file.Translation{
					&file.Translation{
						Lng:     "it",
						Content: "l2.it",
					},
					&file.Translation{
						Lng:     "fr",
						Content: "l2.fr",
					},
					&file.Translation{
						Lng:     "en",
						Content: "l2.en",
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

	data = &file.Labels{}
	err := readCsv(src, "replace")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if data == nil {
		t.Error("Expected data to be set, but is nil")
	}

	if data.Langs == nil || len(data.Langs) != 3 {
		t.Errorf("Invalid number of languages: %# v", pretty.Formatter(data.Langs))
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

	err := readCsv(existing, "replace")
	if err != nil {
		t.Errorf("Cannot set existing labels: %s", err)
	}

	additional := bytes.NewBufferString(`key,en,fr,it
l1,main.en,main.fr,main.it
l2,l2.en,l2.fr,l2.it
l3,l3.en,l3.fr,l3.it
`)

	err = readCsv(additional, "")
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
