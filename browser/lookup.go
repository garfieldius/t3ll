// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package browser

import (
	"errors"
	"os"
	"os/exec"
)

func lookup() (string, error) {
	for _, f := range files {
		bin, err := exec.LookPath(f)
		if err == nil && bin != "" {
			return bin, nil
		}
	}

	for _, p := range paths {
		stat, err := os.Stat(p)
		if err == nil && !stat.IsDir() {
			return p, nil
		}
	}

	return "", errors.New("no chromium or google chrome installation found")
}
