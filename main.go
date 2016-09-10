package main

/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/garfieldius/t3ll/browser"
	"github.com/garfieldius/t3ll/file"
	"github.com/garfieldius/t3ll/log"
	"github.com/garfieldius/t3ll/server"
	"github.com/kr/pretty"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough input arguments: %# v", pretty.Formatter(os.Args))
		return
	}

	var tree *file.Labels

	for i := 1; i < len(os.Args); i++ {
		filename := os.Args[i]
		log.Msg("Checking argument %s", filename)

		if len(filename) < 4 || filename[0:1] == "-" {
			continue
		}

		f, err := file.Open(filename)
		if err == nil {
			tree = f
			break
		}

		log.Msg("Cannot open %s: %s", filename, err)
	}

	if tree == nil {
		log.Fatal("No valid file argument given")
	}

	stop := make(chan bool)
	url, err := server.Start(tree, stop)
	if err != nil {
		log.Fatal("Cannot start server: %s", err)
	}

	err = browser.Start(url)
	if err != nil {
		server.Stop()
		log.Fatal("Cannot start browser")
	}

	cancel := make(chan os.Signal)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	select {
	case <-cancel:
		browser.Stop()
		server.Stop()
		return
	case <-stop:
		return
	}
}
