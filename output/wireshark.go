/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package output

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/tdimitrov/tranqap/tqlog"
)

type wsharkOutput struct {
	pid   int
	stdin io.WriteCloser
	event MOEventChan
}

// NewWsharkOutput constructs wsharkOutput object
func NewWsharkOutput(eventCh MOEventChan) Outputer {
	// Fork wireshark
	cmd := exec.Command("wireshark", "-k", "-i", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	err = cmd.Start()
	if err != nil {
		return nil
	}

	ret := &wsharkOutput{cmd.Process.Pid, stdin, eventCh}

	go func() {
		cmd.Wait()
		eventCh <- MultiOutputEvent{ret, OutputerDead}
	}()

	return ret
}

func (pw wsharkOutput) Write(p []byte) (n int, err error) {
	n, err = pw.stdin.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing: %v", err)
		tqlog.Info(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *wsharkOutput) Close() {
	pw.stdin.Close()
}
