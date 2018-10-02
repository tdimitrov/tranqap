package capture

import (
	"fmt"

	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/shell"
	"golang.org/x/crypto/ssh"
)

// Tcpdump is Capturer implementation for tcpdump
type Tcpdump struct {
	dest       string
	config     ssh.ClientConfig
	captureCmd string

	client  *ssh.Client
	session *ssh.Session
	pid     int
	output  []output.Outputer
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(dest string, config *ssh.ClientConfig, outputs []output.Outputer) Capturer {
	const captureCmd = "sudo tcpdump -U -s0 -w - 'ip and not port 22'"
	return &Tcpdump{dest, *config, captureCmd + shell.StderrToDevNull + shell.RunInBackground + shell.CmdHandlePid(), nil, nil, -1, outputs}
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
	sess, err := capt.client.NewSession()
	if err != nil {
		fmt.Println("Error creating session for stop command!")
		return false
	}

	defer sess.Close()

	err = sess.Start(fmt.Sprintf("sudo kill %d", capt.pid))
	if err != nil {
		fmt.Println("Error running stop command!")
		return false
	}

	sess.Wait()

	return true
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

	chanPid := make(chan int)

	capt.session.Stdout = writer
	capt.session.Stderr, _ = output.NewPidOutput(chanPid)

	err = capt.session.Start(capt.captureCmd)
	if err != nil {
		fmt.Println("Error running command!")
		return false
	}

	capt.pid = <-chanPid

	capt.session.Wait()

	return true
}
