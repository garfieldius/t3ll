package server

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

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/garfieldius/t3ll/file"
	"github.com/garfieldius/t3ll/log"
)

var (
	srv  *http.Server
	wg   sync.WaitGroup
	data *file.Labels
)

func Start(start *file.Labels) (string, error) {
	data = start

	var l net.Listener
	var err error
	var port string

	for i := 2000; i < 8000; i++ {
		port = ":" + strconv.Itoa(i)
		l, err = net.Listen("tcp", port)

		if err != nil {
			log.Msg("Cannot start server on port %s: %s", port, err)
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}

	if l == nil {
		return "", errors.New("Cannot start server")
	}

	srv = &http.Server{}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/quit", quitHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Serve(l)
	}()

	return "http://localhost" + port + "/", nil
}

func Stop() {
	if srv == nil {
		return
	}

	log.Msg("Stopping HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
