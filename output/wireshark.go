package output

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type wsharkOutput struct {
	pid   int
	read  *io.PipeReader
	write *io.PipeWriter
}

// NewWsharkOutput constructs wsharkOutput object
func NewWsharkOutput() (Outputer, error) {
	// Create pipe
	r, w := io.Pipe()

	// Fork wireshark
	cmd := exec.Command("wireshark", "-k", "-i", "-")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = r

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	pid := cmd.Process.Pid

	return &wsharkOutput{pid, r, w}, nil
}

func (pw wsharkOutput) Write(p []byte) (n int, err error) {
	n, err = pw.write.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing: %v", err)
		fmt.Println(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *wsharkOutput) Close() {
	pw.write.Close()
	pw.read.Close()
}
