// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package main

import (
	"fmt"
	"os"

	"github.com/garfieldius/t3ll/log"
)

func main() {
	os.Exit(run())
}

func run() int {
	var exitCode = 0
	defer func() {
		if err := recover(); err != nil {
			log.Err("%s", err)
			exitCode = 4
		}
	}()

	app := App{}
	file, state, err := app.Init()

	switch {
	case err != nil:
		log.Err("%s", err)
		return 2
	case file == "version":
		fmt.Printf(versionText, Version)
		return 0
	case file == "help":
		fmt.Printf(helpText, Version)
		return 0
	default:
		if err := app.Run(state); err != nil {
			log.Err("%s", err)
			return 1
		}
	}

	return exitCode
}
