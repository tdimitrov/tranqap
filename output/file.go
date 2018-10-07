package output

import (
	"errors"
	"fmt"
	"os"
)

type fileOutput struct {
	fd *os.File
}

// NewFileOutput constructs fileOutput object
func NewFileOutput(fname string) (Outputer, error) {
	fd, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	return &fileOutput{fd}, nil
}

func (pw fileOutput) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file: %v", err)
		fmt.Println(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *fileOutput) Close() {
	pw.fd.Close()
}
