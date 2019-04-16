// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

// Package browser offers a type to open a chromium or google chrome
// browser window in app-mode with the given URL as entrypoint
package browser

import (
	"context"
	"os/exec"

	"github.com/garfieldius/t3ll/log"
)

// Browser handles a chromium or google chrome process
type Browser struct {
	cancel context.CancelFunc
	Done   chan error
}

// Start will create a new browser process
func (b *Browser) Start(url string) error {
	bin, err := lookup()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	log.Msg("Running %s --disable-background-mode --disable-plugins --disable-plugins-discovery --reset-variation-state --single-tab-mode --app=%s", bin, url)
	cmd := exec.CommandContext(ctx, bin, "--disable-background-mode", "--disable-plugins", "--disable-plugins-discovery", "--reset-variation-state", "--single-tab-mode", "--app="+url)
	err = cmd.Start()
	if err != nil {
		cancel()
		return err
	}
	b.cancel = cancel
	b.Done = make(chan error)

	go func() {
		if err := cmd.Wait(); err != nil && !cmd.ProcessState.Success() {
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
