// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/garfieldius/t3ll/labels"
	"github.com/garfieldius/t3ll/log"
)

// Server is a http server sending and receiving data of the app
type Server struct {
	srv  *http.Server
	Done chan error
	Host string
}

// Start creates a http server listener
func (s *Server) Start(state *labels.Labels, quitSig chan struct{}) error {
	var l net.Listener
	var err error

	for i := 1025; i < 65000; i++ {
		s.Host = "127.0.0.1:" + strconv.Itoa(i)
		l, err = net.Listen("tcp", s.Host)

		if err != nil {
			log.Msg("Cannot start server on %s: %s", s.Host, err)
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}

	if l == nil {
		return errors.New("Cannot start server")
	}

	s.srv = &http.Server{}
	s.srv.Handler = handler{state: state, quitSig: quitSig, mu: new(sync.Mutex)}
	s.Done = make(chan error)

	go func() {
		if err := s.srv.Serve(l); err != nil && err != http.ErrServerClosed {
			s.Done <- err
		} else {
			s.Done <- nil
		}
	}()

	return nil
}

// Stop will stop the running server
func (s *Server) Stop() {
	if s.srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.srv.Close()
	}
	cancel()

	s.srv = nil
}
