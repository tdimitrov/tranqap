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
	cmd :=
		`func rpcap_kill {
			kill %d
			R1=$?
			sleep 0.1
			kill -0 %d
			R2=$?
			echo $R1 $R2
		}
		rpcap_kill`
	return fmt.Sprintf(cmd, pid, pid)
	//return fmt.Sprintf("kill %d ; R1=$?; sleep 0.1; kill -0 %d; R2=$?; echo $R1 $R2", pid, pid)
}

const (
	// EvKillSuccess is returned when the process is killed successfully. R1=0; R2=NonZero
	EvKillSuccess = iota
	// EvKillNotRuning is returned when the process that should be killed is not runnung
	EvKillNotRuning = iota
	// EvKillNotResponding is returned when the process is running, but doesn't die after the kill command
	EvKillNotResponding = iota
	// EvKillError is returned when the process is not running, but the 2nd kill doesn't return error.
	// This could happen in very strange circumstances.
	EvKillError = iota
	// EvKillExecError is returned when result from the bash command is malformed.
	EvKillExecError = iota
)

// KillResToStr returns a string with the description of the error code used by CmdKillPid (described in the consts above)
func KillResToStr(res int) string {
	switch res {
	case EvKillError:
		return "Unknown error"
	case EvKillNotResponding:
		return "PID not responding"
	case EvKillNotRuning:
		return "PID not running"
	case EvKillSuccess:
		return "Success (no error)"
	case EvKillExecError:
		return "Execution error"
	}

	return "Unknown error code"
}

type killPidHandler struct {
	result chan<- int // Sends EvKill*** from the consts above
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
		pw.result <- EvKillExecError
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

	pw.sendResult(r1, r2)

	close(pw.result)

	return len(p), nil
}

func (pw killPidHandler) sendResult(r1, r2 int) {
	if r1 == 0 {
		if r2 == 0 {
			pw.result <- EvKillNotResponding
		} else {
			pw.result <- EvKillSuccess
		}
	} else {
		if r2 == 0 {
			pw.result <- EvKillError
		} else {
			pw.result <- EvKillNotRuning
		}

	}
}

func (pw killPidHandler) Close() {
}
