package output

import (
	"strconv"
	"strings"
)

type killOutput struct {
	result chan<- int
}

// NewKillOutput creates new killOutput instance.
// It reads the result codes of both kill commands and reports any errors
func NewKillOutput(res chan<- int) Outputer {
	return killOutput{res}
}

func (pw killOutput) Write(p []byte) (n int, err error) {
	data := string(p)
	data = strings.Trim(data, "\n\t ")

	results := strings.Split(data, " ")
	if len(results) != 2 {
		pw.result <- -1
		pw.result <- -1
		close(pw.result)

		return len(p), nil
	}

	r1, err := strconv.Atoi(results[0])
	if err != nil {
		r1 = -1
	}

	r2, err := strconv.Atoi(results[1])
	if err != nil {
		r2 = -1
	}

	pw.result <- r1
	pw.result <- r2

	close(pw.result)

	return len(p), nil
}

func (pw killOutput) Close() {
}
