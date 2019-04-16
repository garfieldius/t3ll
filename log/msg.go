// +build !debug

// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package log

// Msg does not do anything in default build
func Msg(msg string, a ...interface{}) {
	// Noop
}
