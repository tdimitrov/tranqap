package capture

import (
	"fmt"

	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/rplog"
	"github.com/tdimitrov/rpcap/shell"
	"golang.org/x/crypto/ssh"
)

// Tcpdump is Capturer implementation for tcpdump
type Tcpdump struct {
	name       string
	dest       string
	config     ssh.ClientConfig
	captureCmd string
	client     *ssh.Client
	session    *ssh.Session
	pid        *shell.StdErrHandler
	out        *output.MultiOutput
	onDie      CapturerEventChan
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(name string, dest string, config *ssh.ClientConfig, outer *output.MultiOutput, subsc CapturerEventChan) Capturer {
	const captureCmd = "tcpdump -U -s0 -w - 'ip and not port 22'"
	const runInBackground = " & "
	//const stderrToDevNull = " 2> /dev/null "

	return &Tcpdump{
		fmt.Sprintf("<%s>", name),
		dest,
		*config,
		captureCmd + runInBackground + shell.CmdGetPid(),
		nil,
		nil,
		shell.NewStdErrHandler(),
		outer,
		subsc,
	}
}

// Start method connects the ssh client to the destination and start capturing
func (capt *Tcpdump) Start() bool {
	if capt.session != nil || capt.client != nil {
		rplog.Error("There is an active session for capturer %s", capt.Name())
		return false
	}

	var err error
	capt.client, err = connect(capt.dest, &capt.config)
	if err != nil {
		rplog.Error("Error connecting to target %s: %s", capt.Name(), err)
		return false
	}

	go capt.startSession()

	return true
}

// Stop terminates the capture
func (capt *Tcpdump) Stop() bool {
	sess, err := capt.client.NewSession()
	if err != nil {
		rplog.Error("capture.Tcpdump: Error creating session for Stop() for capturer %s", capt.Name())
		return false
	}

	defer sess.Close()

	pid := capt.pid.GetPid()
	// Clear PID to indicate an expected kill
	capt.pid.ClearPid()

	err = sess.Start(fmt.Sprintf("kill %d", pid))
	if err != nil {
		rplog.Error("capture.Tcpdump: Error starting kill command for capturer %s", capt.Name())
		return false
	}

	sess.Wait()

	capt.out.Close()

	rplog.Info("capture.Tcpdump: Kill executed successfully for capturer %s", capt.Name())
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
		rplog.Error("Error creating initial session for capturer %s", capt.Name())
		return false
	}

	defer capt.session.Close()
	defer capt.out.Close()

	capt.session.Stdout = capt.out
	capt.session.Stderr = *capt.pid

	err = capt.session.Start(capt.captureCmd)
	if err != nil {
		rplog.Error("Error running tcpdum command for capturer %s", capt.Name())
		return false
	}

	capt.session.Wait()

	if capt.pid.GetPid() != -1 {
		// PID is not cleared - this is unexpected stop
		capt.onDie <- CapturerEvent{capt, CapturerDead}
		rplog.Error("capture.Tcpdump: Capturer %s died unexpectedly. Dumping stderr:", capt.Name())
		capt.pid.DumpStdErr()
	} else {
		rplog.Info("capture.Tcpdump: Capturer %s killed by command", capt.Name())
		capt.onDie <- CapturerEvent{capt, CapturerStopped}
	}

	return true
}

// Name returns the name of the capturer's target (used only for logging purposes)
func (capt *Tcpdump) Name() string {
	return capt.name
}
