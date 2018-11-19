package capture

import (
	"fmt"
	"io"
	"strings"

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
	useSudo    bool
}

// SudoConfig contains config params regarding sudo usage.
// Use is a bool which indicates if tcpdump should be started with sudo
// Username is pointer to a string with the username, which tcpdump will use
// to drop privilege to (-Z option ). If UseSudo is false, Username is nil
type SudoConfig struct {
	Use      bool
	Username *string
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(name string, outer *output.MultiOutput, subsc CapturerEventChan, trans captureTransport, sudo SudoConfig) Capturer {
	const sudoCmd = "sudo -n "
	const captureCmd = "tcpdump -U -s0 -w - 'ip and not port 22'"
	const dropPriviledges = " -Z "
	const runInBackground = " & "

	var cmd strings.Builder
	if sudo.Use == true {
		cmd.WriteString(sudoCmd)
	}
	cmd.WriteString(captureCmd)
	if sudo.Use == true {
		cmd.WriteString(dropPriviledges)
		cmd.WriteString(*sudo.Username)
	}
	cmd.WriteString(runInBackground)
	cmd.WriteString(cmdGetPid())

	return &Tcpdump{
		fmt.Sprintf("<%s>", name),
		cmd.String(),
		newStdErrHandler(),
		outer,
		subsc,
		trans,
		sudo.Use,
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
func (capt *Tcpdump) Stop() error {
	pid := capt.pid.GetPid()
	// Clear PID to indicate an expected kill
	capt.pid.ClearPid()

	var cmd string
	if capt.useSudo == true {
		// the sudo process runs as root and it can't be killed with a regular user
		// that's why the child process is killed
		cmd = fmt.Sprintf("kill `ps --ppid %d -o pid=`", pid)

	} else {
		cmd = fmt.Sprintf("kill %d", pid)
	}

	err := capt.trans.Run(cmd, nil, nil)
	if err != nil {
		return err
	}

	capt.out.Close()

	rplog.Info("capture.Tcpdump: '%s' executed successfully for capturer %s", cmd, capt.Name())
	return nil
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
