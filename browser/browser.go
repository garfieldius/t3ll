// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

// Package browser offers a type to open a chromium or google chrome
// browser window in app-mode with the given URL as entrypoint
package browser

import (
	"context"
	"github.com/garfieldius/t3ll/log"
	"github.com/go-rod/rod/lib/launcher"
	"os/exec"
)

// Browser handles a chromium or google chrome process
type Browser struct {
	cancel context.CancelFunc
	Done   chan error
}

// Start will create a new browser window
func (b *Browser) Start(url string) error {
	br := launcher.NewBrowser()
	l := launcher.NewAppMode(url)

	bin, has := launcher.LookPath()
	if !has {
		downloadedBin, err := br.Get()

		if err != nil {
			return err
		}
		log.Msg("Using downloaded chrome located at %s", downloadedBin)
		bin = downloadedBin
	} else {
		log.Msg("Using already installed chrome located at %s", bin)
	}

	l.Bin(bin)

	ctx, cancel := context.WithCancel(context.Background())
	chromeParams := l.FormatArgs()

	log.Msg("Running %v", chromeParams)

	cmd := exec.CommandContext(ctx, bin, chromeParams...)
	err := cmd.Start()
	if err != nil {
		log.Err("Did not start: %s", err)
		cancel()
		return err
	}

	b.cancel = cancel
	b.Done = make(chan error)

	go func() {
		err := cmd.Wait()
		log.Msg("Finished browser process with %s", err)
		if err != nil && !cmd.ProcessState.Success() {
			log.Err("Browser quit unexpectedly")
			b.Done <- err
		} else {
			b.Done <- nil
		}
	}()

	return nil
}

// Stop will stop the current browser process
func (b *Browser) Stop() {
	if b.cancel == nil {
		return
	}

	b.cancel()
	b.cancel = nil
}
