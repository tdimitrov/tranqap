package output

import (
	"fmt"
	"strconv"
	"strings"
)

type pidOutput struct {
	result chan<- int
}

// NewPidOutput creates new pidOutput instance.
// It reads a PID from the buffer, passed to Write(). pidOutput is used to parse
// the PID of the capturer so that it can be stopped on user request.
// It's input parameter is a channel, used to return the PID as an integer
func NewPidOutput(pid chan<- int) (Outputer, error) {
	return pidOutput{pid}, nil
}

func (pw pidOutput) Write(p []byte) (n int, err error) {
	data := string(p)

	if strings.HasPrefix(data, "MY_PID_IS:") {
		// The PID is sent. Parse it and send it over the channel
		data := strings.Replace(data, "MY_PID_IS:", "", 1)
		pid, err := strconv.Atoi(strings.Trim(data, "\n\t "))
		if err != nil {
			fmt.Println("Expected PID, received", data)
		}

		pw.result <- pid
		close(pw.result)
	} else {
		// It's something else. Log it on the screen
		fmt.Println(data)
	}

	return len(p), nil
}

func (pw pidOutput) Close() {
}
