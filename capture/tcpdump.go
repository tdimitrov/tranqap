package capture

import (
	"fmt"
	"sync"

	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/shell"
	"golang.org/x/crypto/ssh"
)

type atomicPid struct {
	pid int
	mut sync.Mutex
}

func (p *atomicPid) Set(val int) {
	p.mut.Lock()
	p.pid = val
	p.mut.Unlock()
}

func (p *atomicPid) Get() int {
	p.mut.Lock()
	val := p.pid
	p.mut.Unlock()
	return val
}

// Tcpdump is Capturer implementation for tcpdump
type Tcpdump struct {
	dest       string
	config     ssh.ClientConfig
	captureCmd string
	client     *ssh.Client
	session    *ssh.Session
	pid        atomicPid
	out        *output.MultiOutput
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(dest string, config *ssh.ClientConfig, outer *output.MultiOutput) Capturer {
	const captureCmd = "sudo tcpdump -U -s0 -w - 'ip and not port 22'"
	return &Tcpdump{
		dest,
		*config,
		captureCmd + shell.StderrToDevNull + shell.RunInBackground + shell.CmdGetPid(),
		nil,
		nil,
		atomicPid{},
		outer,
	}
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

	results := make(chan int, 2)
	sess.Stdout = shell.NewKillPidHandler(results)

	err = sess.Start(shell.CmdKillPid(capt.pid.Get()))
	if err != nil {
		fmt.Println("Error running stop command!")
		return false
	}

	sess.Wait()

	r1 := <-results
	r2 := <-results

	if r1 != 0 {
		fmt.Println("Process is not running")
	} else if r2 == 0 {
		fmt.Println("Process is not responding")
	}

	capt.out.Close()

	return true
}

// AddOutputer calls AddMember of the MultiOutput instance of Tcpdump
func (capt *Tcpdump) AddOutputer(newOutputerFn output.OutputerFactory) error {
	return capt.out.AddMember(newOutputerFn)
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

	writer := capt.out
	defer writer.Close()

	chanPid := make(chan int)

	capt.session.Stdout = writer
	capt.session.Stderr, _ = shell.NewGetPidHandler(chanPid)

	err = capt.session.Start(capt.captureCmd)
	if err != nil {
		fmt.Println("Error running command!")
		return false
	}

	capt.pid.Set(<-chanPid)

	capt.session.Wait()

	return true
}
