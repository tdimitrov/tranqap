package shell

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// pidPrefix is the string, which is put in front of the PID, when it is transmitted over stderr
const pidPrefix = "MY_PID_IS:"

// CmdGetPid returns a Bash one-liner, which does the following:
// 1. Saves the PID of the last command executed in a Bash variable. This is supposed to be the capture command
// 2. Echoes the PID to stderr, so that it can be saved by the Capturer. Stderr is used, because PCAP data is
//		transmitted over stdout
// 3. Waits the PID to finish, so that the session remains active until stop command is sent from rpcap shell
func CmdGetPid() string {
	return "RPCAP_MY_PID=$! ; echo " + pidPrefix + " $RPCAP_MY_PID >&2 ; wait $RPCAP_MY_PID"
}

// StdErrHandler parses PID of the Capturer from stderr and saves all stderr messages in a string slice
// If needed these messages are dumped to the user
type StdErrHandler struct {
	pid        *int // Pointer, because Write() has got value receiver. Requirement of ssh lib
	pidLock    *sync.Mutex
	errLog     *[]string
	errLogLock *sync.Mutex
}

// NewStdErrHandler creates new pidOutput instance.
// It reads a PID from the buffer, passed to Write(). pidOutput is used to parse
// the PID of the capturer so that it can be stopped on user request.
// It's input parameter is a channel, used to return the PID as an integer
func NewStdErrHandler() *StdErrHandler {
	pid := -1
	return &StdErrHandler{&pid, &sync.Mutex{}, new([]string), &sync.Mutex{}}
}

func (pw StdErrHandler) Write(p []byte) (n int, err error) {
	data := string(p)

	if strings.HasPrefix(data, pidPrefix) {
		// The PID is sent. Parse it and send it over the channel
		data := strings.Replace(data, pidPrefix, "", 1)
		pid, err := strconv.Atoi(strings.Trim(data, "\n\t "))
		if err != nil {
			fmt.Println("Expected PID, received: ", data)
			pid = -1
		}
		pw.pidLock.Lock()
		*pw.pid = pid
		pw.pidLock.Unlock()
	} else {
		// Prefix not found in response. Save the output.
		pw.errLogLock.Lock()
		*pw.errLog = append(*pw.errLog, data)
		pw.errLogLock.Unlock()
	}

	return len(p), nil
}

func (pw *StdErrHandler) GetPid() int {
	pw.pidLock.Lock()
	pid := *pw.pid
	pw.pidLock.Unlock()
	return pid
}

func (pw *StdErrHandler) ClearPid() {
	pw.pidLock.Lock()
	*pw.pid = -1
	pw.pidLock.Unlock()
}

func (pw *StdErrHandler) DumpStdErr() {
	pw.errLogLock.Lock()
	for _, errmsg := range *pw.errLog {
		fmt.Println(errmsg)
	}
	pw.errLogLock.Unlock()
}
