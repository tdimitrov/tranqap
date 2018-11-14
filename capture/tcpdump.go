package capture

import (
	"fmt"
	"io"

	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/rplog"
)

type captureTransport interface {
	IsActive() bool
	Connect() error
	Run(cmd string, stdout io.Writer, stderr io.Writer) error
}

// Tcpdump is Capturer implementation for tcpdump
type Tcpdump struct {
	name       string
	captureCmd string
	pid        *stdErrHandler
	out        *output.MultiOutput
	onDie      CapturerEventChan
	trans      captureTransport
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(name string, outer *output.MultiOutput, subsc CapturerEventChan, trans captureTransport) Capturer {
	const captureCmd = "tcpdump -U -s0 -w - 'ip and not port 22'"
	const runInBackground = " & "
	//const stderrToDevNull = " 2> /dev/null "

	return &Tcpdump{
		fmt.Sprintf("<%s>", name),
		captureCmd + runInBackground + cmdGetPid(),
		newStdErrHandler(),
		outer,
		subsc,
		trans,
	}
}

// Start method connects the ssh client to the destination and start capturing
func (capt *Tcpdump) Start() bool {
	if capt.trans.IsActive() {
		rplog.Error("There is an active session for capturer %s", capt.Name())
		return false
	}

	if err := capt.trans.Connect(); err != nil {
		capt.out.Close()
		rplog.Error("Error connecting to target %s: %s", capt.Name(), err)
		return false
	}

	go capt.startSession()

	return true
}

// Stop terminates the capture
func (capt *Tcpdump) Stop() bool {
	pid := capt.pid.GetPid()
	// Clear PID to indicate an expected kill
	capt.pid.ClearPid()

	err := capt.trans.Run(fmt.Sprintf("kill %d", pid), nil, nil)
	if err != nil {
		rplog.Error("capture.Tcpdump: Error starting kill command for capturer %s: %s", capt.Name(), err)
		return false
	}

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

	defer capt.out.Close()

	err = capt.trans.Run(capt.captureCmd, capt.out, capt.pid)
	if err != nil {
		rplog.Error("Error running tcpdump command for capturer %s: ", capt.Name(), err)
		capt.onDie <- CapturerEvent{capt, CapturerDead}
		return false
	}

	if capt.pid.GetPid() != -1 {
		// PID is not cleared - this is unexpected stop
		capt.onDie <- CapturerEvent{capt, CapturerDead}
		rplog.Error("capture.Tcpdump: Capturer %s died unexpectedly. Dumping stderr:\n%s",
			capt.Name(), capt.pid.DumpStdErr())
		rplog.Feedback("Capturer %s died. stderr:\n%s", capt.Name(), capt.pid.DumpStdErr())
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
