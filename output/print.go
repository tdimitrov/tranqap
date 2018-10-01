package output

import (
	"fmt"
)

type printOutput struct {
}

// NewPrintOutput creates new printOutput instance. It logs all data to stdout
// Used only for debug purposes
func NewPrintOutput() (Outputer, error) {
	return printOutput{}, nil
}

func (pw printOutput) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func (pw printOutput) Close() {
}
