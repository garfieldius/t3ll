package browser

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
	"os/exec"
	"sync"

	"github.com/garfieldius/t3ll/log"
	"github.com/garfieldius/t3ll/server"
)

var (
	wg sync.WaitGroup
	c  *exec.Cmd
)

func Start(url string) error {
	cmd, args := open(url)
	c = exec.Command(cmd, args...)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Run()
	}()

	return nil
}

func Stop() {
	if c == nil {
		return
	}

	log.Msg("Stopping browser process")

	c.Process.Kill()
	c = nil

	wg.Wait()
	server.Stop()
}
