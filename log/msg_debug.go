//go:build debug
// +build debug

// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package log

import (
	"fmt"
	"os"
)

// Msg prints a message to stdout in debug build
func Msg(msg string, a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stdout, format("INFO", msg, a...))
}

// Err prints a formatted message to stderr
func Err(msg string, a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, format("ERROR", msg, a...))
}
