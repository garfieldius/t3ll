// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package main

import (
	"errors"
	"os"
	"os/signal"

	"github.com/garfieldius/t3ll/browser"
	"github.com/garfieldius/t3ll/labels"
	"github.com/garfieldius/t3ll/log"
	"github.com/garfieldius/t3ll/server"
)

// App represents the application consisting of server and browser process
type App struct {
}

// Run starts server and browser with the currently set state
func (a *App) Run(state *labels.Labels) error {
	s := server.Server{}
	b := browser.Browser{}
	quitSig := make(chan struct{})
	if err := s.Start(state, quitSig); err != nil {
		return err
	}
	if err := b.Start("http://" + s.Host + "/"); err != nil {
		s.Stop()
		return err
	}

	cancel := make(chan os.Signal)
	signal.Notify(cancel, os.Interrupt)

	select {
	case <-quitSig:
	case <-cancel:
		log.Msg("Received quit signal")
		s.Stop()
		b.Stop()
		return nil

	case err := <-s.Done:
		log.Msg("Server quit, stop browser and quit")
		if err != nil {
			log.Err("Server stopped with error: %s", err)
		}
		b.Stop()
		return nil

	case err := <-b.Done:
		log.Msg("Browser quit, stop server and quit")
		if err != nil {
			log.Err("Browser stopped with error: %s", err)
		}
		s.Stop()
		return nil
	}

	return nil
}

// Init will try to find a file to edit based on
// process arguments
func (a *App) Init() (string, *labels.Labels, error) {
	if len(os.Args) < 2 {
		return "", nil, errors.New("not enough arguments")
	}

	for _, arg := range os.Args[1:] {
		if arg == "version" || arg == "help" {
			return arg, nil, nil
		}

		log.Msg("Checking argument %s", arg)

		if len(arg) < 4 || arg[0:1] == "-" {
			continue
		}

		s, err := labels.Open(arg)
		if err != nil {
			log.Msg("Cannot open %s: %s", arg, err)
			continue
		}

		return arg, s, nil
	}

	return "", nil, errors.New("no valid arg given")
}
