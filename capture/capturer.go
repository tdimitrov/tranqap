package capture

import (
	"golang.org/x/crypto/ssh"
)

// Capturer interface represents a general capturer. There are concrete implementations
// for tcpdump. In the future more can be added, e.g. tshark, dumpcap, etc.
type Capturer interface {
	Start() bool
	Stop() bool
}

func connect(dest string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {

	client, err := ssh.Dial("tcp", dest, clientConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

// bashCmdHandlePid() returns a Bash one-liner, which does the following:
// 1. Saves the PID of the last command executed in a Bash variable. This is supposed to be the capture command
// 2. Echoes the PID to stderr, so that it can be saved by the Capturer. Stderr is used, because PCAP data is
//		transmitted over stdout
// 3. Waits the PID to finish, so that the session remains active until stop command is sent from rpcap shell
func bashCmdHandlePid() string {
	return "RPCAP_MY_PID=$! ; echo $RPCAP_MY_PID >&2 ; wait $RPCAP_MY_PID"
}
