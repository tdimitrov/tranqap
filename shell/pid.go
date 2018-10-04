package shell

import (
	"fmt"
	"strconv"
	"strings"
)

// PidPrefix is the string, which is put in front of the PID, when it is transmitted over stderr
const PidPrefix = "MY_PID_IS:"

// CmdHandlePid returns a Bash one-liner, which does the following:
// 1. Saves the PID of the last command executed in a Bash variable. This is supposed to be the capture command
// 2. Echoes the PID to stderr, so that it can be saved by the Capturer. Stderr is used, because PCAP data is
//		transmitted over stdout
// 3. Waits the PID to finish, so that the session remains active until stop command is sent from rpcap shell
func CmdHandlePid() string {
	return "RPCAP_MY_PID=$! ; echo " + PidPrefix + " $RPCAP_MY_PID >&2 ; wait $RPCAP_MY_PID"
}

type pidOutput struct {
	result chan<- int
}

// NewPidOutput creates new pidOutput instance.
// It reads a PID from the buffer, passed to Write(). pidOutput is used to parse
// the PID of the capturer so that it can be stopped on user request.
// It's input parameter is a channel, used to return the PID as an integer
func NewPidOutput(pid chan<- int) (CmdHandler, error) {
	return pidOutput{pid}, nil
}

func (pw pidOutput) Write(p []byte) (n int, err error) {
	data := string(p)

	if strings.HasPrefix(data, PidPrefix) {
		// The PID is sent. Parse it and send it over the channel
		data := strings.Replace(data, PidPrefix, "", 1)
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
