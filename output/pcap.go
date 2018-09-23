package output

import (
	"errors"
	"fmt"
	"os"
)

type pcapOutput struct {
	fd *os.File
}

// NewPcapOutput constructs pcapOutput object
func NewPcapOutput(fname string) (Outputer, error) {
	fd, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	return &pcapOutput{fd}, nil
}

func (pw pcapOutput) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file: %v", err)
		fmt.Println(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *pcapOutput) Close() {
	pw.fd.Close()
}
