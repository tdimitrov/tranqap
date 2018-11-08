package output

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/tdimitrov/rpcap/rplog"
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
		rplog.Info(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *wsharkOutput) Close() {
	pw.stdin.Close()
}
