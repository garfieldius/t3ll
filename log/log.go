// Package log is a convenience wrapper around go's log and fmt packages

// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package log

import (
	"fmt"
	"os"
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

// Err prints a formatted message to stderr
func Err(msg string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, format("ERROR", msg, a...))
}
