package capture

import (
	"fmt"

	"github.com/tdimitrov/rpcap/output"
	"golang.org/x/crypto/ssh"
)

// Tcpdump is Capturer implementation for tcpdump
type Tcpdump struct {
	dest       string
	config     ssh.ClientConfig
	captureCmd string

	client  *ssh.Client
	session *ssh.Session

	output []output.Outputer
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(dest string, config *ssh.ClientConfig, outputs []output.Outputer) Capturer {
	return &Tcpdump{dest, *config, "sudo tcpdump -U -s0 -w - 'ip and not port 22'", nil, nil, outputs}
}

// Start method connects the ssh client to the destination and start capturing
func (capt *Tcpdump) Start() bool {
	if capt.session != nil || capt.client != nil {
		fmt.Println("There is an active session for this capture")
		return false
	}

	var err error
	capt.client, err = connect(capt.dest, &capt.config)
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return false
	}

	go capt.startSession()

	return true
}

// Stop terminates the capture
func (capt *Tcpdump) Stop() bool {
	err := capt.session.Signal(ssh.SIGINT)

	if err != nil {
		fmt.Println(err)
	}
	return false
}

func (capt *Tcpdump) startSession() bool {
	//fmt.Println(client.LocalAddr().(*net.TCPAddr).IP)
	var err error
	capt.session, err = capt.client.NewSession()
	if err != nil {
		fmt.Println("Error creating session!")
		return false
	}

	defer capt.session.Close()

	writer := capt.output[0]
	defer writer.Close()

	capt.session.Stdout = writer

	err = capt.session.Start(capt.captureCmd)
	if err != nil {
		fmt.Println("Error running command!")
		return false
	}

	capt.session.Wait()

	return true
}
