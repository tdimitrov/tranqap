package shell

import (
	"fmt"
	"strconv"
	"strings"
)

// CmdKillPid returns a Bash one-liner, which does the following:
// 1. Sends SIGINIT to pid, so the process can terminate gracefully
// 2. Send signal 0 to the same pid, to be sure that the process has terminated
// 3. Returns the result codes from both command to the caller via stdout
func CmdKillPid(pid int) string {
	return fmt.Sprintf("sudo kill %d ; R1=$?; kill -0 %d; R2=$?; echo $R1 $R2", pid, pid)
}

type killPidHandler struct {
	result chan<- int
}

// NewKillPidHandler creates new killOutput instance.
// It reads the result codes of both kill commands and reports any errors
func NewKillPidHandler(res chan<- int) CmdHandler {
	return killPidHandler{res}
}

func (pw killPidHandler) Write(p []byte) (n int, err error) {
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

func (pw killPidHandler) Close() {
}
