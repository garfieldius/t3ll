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
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/garfieldius/t3ll/file"
	"github.com/garfieldius/t3ll/log"
	"github.com/hydrogen18/stoppableListener"
)

var (
	wg       sync.WaitGroup
	listener *stoppableListener.StoppableListener
	data     *file.Labels
	stop     chan bool
)

func Start(start *file.Labels, onstop chan bool) (string, error) {
	data = start
	stop = onstop

	var originalListener net.Listener
	var err error
	var port string

	for i := 2000; i < 8000; i++ {
		port = ":" + strconv.Itoa(i)
		originalListener, err = net.Listen("tcp", port)

		if err != nil {
			log.Msg("Cannot start server on port %s: %s", port, err)
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}

	if originalListener == nil {
		if err == nil {
			err = errors.New("Cannot start listener")
		}
		return "", err
	}

	listener, err = stoppableListener.New(originalListener)
	if err != nil {
		return "", err
	}

	server := http.Server{}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/quit", quitHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(listener)
	}()

	return "http://localhost" + port + "/", nil
}

func Stop() {
	if listener == nil {
		return
	}

	log.Msg("Stopping HTTP server")

	listener.Stop()
	listener = nil

	wg.Wait()

	time.AfterFunc(10*time.Millisecond, func() {
		stop <- true
	})
}
