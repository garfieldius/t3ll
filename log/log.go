// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

// Package log is a convenience wrapper around go's log and fmt packages
// It is not, by intention, a sophisticated logging package
package log

import (
	"fmt"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

func now() string {
	return time.Now().Format(timeLayout)
}

func format(severity, msg string, a ...interface{}) string {
	full := fmt.Sprintf("%s [%s] %s", now(), severity, msg)
	return fmt.Sprintf(full, a...)
}
