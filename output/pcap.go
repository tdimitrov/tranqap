package output

import (
	"fmt"
	"os"
)

type pcapOutput struct {
	fd *os.File
}

// NewPcapOutput constructs pcapOutput object
func NewPcapOutput(fname string) Outputer {
	fd, err := os.OpenFile(fname, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil
	}

	return &pcapOutput{fd}
}

func (pw pcapOutput) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
	return n, err
}

func (pw *pcapOutput) Close() {
	pw.fd.Close()
}
