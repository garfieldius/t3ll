// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"encoding/json"

	"github.com/garfieldius/t3ll/labels"
	"github.com/garfieldius/t3ll/log"
)

func saveHandler(data []byte, format string, cur *labels.Labels) (*labels.Labels, error) {
	newData := new(labels.Labels)

	if err := json.Unmarshal(data, newData); err != nil {
		return nil, err
	}

	newData.SrcXlif = cur.SrcXlif
	newData.SrcLegacy = cur.SrcLegacy
	newData.FromFile = cur.FromFile
	newData.Type = cur.Type

	if format == string(labels.XMLXliffv1) {
		log.Msg("Converting to xliff")
		newData.Type = labels.XMLXliffv1
		newData.FromFile = newData.FromFile[:len(newData.FromFile)-3] + "xlf"
	}

	if err := newData.Save(); err != nil {
		return nil, err
	}

	return newData, nil
}
