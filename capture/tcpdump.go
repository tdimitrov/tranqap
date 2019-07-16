/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package capture

import (
	"fmt"
	"io"
	"strings"

	"github.com/tdimitrov/tranqap/output"
	"github.com/tdimitrov/tranqap/rplog"
)

type captureTransport interface {
	IsActive() bool
	Connect() error
	Run(cmd string, stdout io.Writer, stderr io.Writer) error
	GetRemoteIP() *string
	GetRemotePort() *int
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
	filter     FilterConfig
}

// SudoConfig contains config params regarding sudo usage.
// Use is a bool which indicates if tcpdump should be started with sudo
// Username is pointer to a string with the username, which tcpdump will use
// to drop privilege to (-Z option ). If UseSudo is false, Username is nil
type SudoConfig struct {
	Use      bool
	Username *string
}

// FilterConfig contains Port, which is set as tcpdump capture
// filter. This is useful for the cases when the target is behind NAT and is
// accessed via port and/or IP redirection. If so, the port used to connect
// to the target will differ from the actual on which the SSH service is
// bound.
type FilterConfig struct {
	Port *int
}

// NewTcpdump creates Tcpdump Capturer
func NewTcpdump(name string, outer *output.MultiOutput, subsc CapturerEventChan, trans captureTransport, sudo SudoConfig, filter FilterConfig) Capturer {
	const sudoCmd = "sudo -n "
	const captureCmd = "tcpdump -U -s0 -i any -w - 'not port %d'"
	const dropPrivileges = " -Z "
	const runInBackground = " & "

	var cmd strings.Builder
	if sudo.Use == true {
		cmd.WriteString(sudoCmd)
	}
	cmd.WriteString(captureCmd)
	if sudo.Use == true {
		cmd.WriteString(dropPrivileges)
		cmd.WriteString(*sudo.Username)
	}
	cmd.WriteString(runInBackground)
	cmd.WriteString(cmdGetPid())

	return &Tcpdump{
		name,
		cmd.String(),
		newStdErrHandler(),
		outer,
		subsc,
		trans,
		sudo.Use,
		filter,
	}
}

// Start method connects the ssh client to the destination and start capturing
func (capt *Tcpdump) Start() error {
	if capt.trans.IsActive() {
		return fmt.Errorf("There is an active session for capturer %s", capt.Name())
	}

	if err := capt.trans.Connect(); err != nil {
		capt.out.Close()
		return fmt.Errorf("Error connecting to %s: %s", capt.Name(), err)
	}

	go capt.startSession()

	rplog.Info("Connected to %s and started a session.", capt.Name())

	return nil
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
		return fmt.Errorf("Error running kill command: %s", err)
	}

	rplog.Info("Kill executed successfully for %s", capt.Name())
	return nil
}

// AddOutputer calls AddMember of the MultiOutput instance of Tcpdump
func (capt *Tcpdump) AddOutputer(newOutputerFn output.OutputerFactory) error {
	return capt.out.AddExtMember(newOutputerFn)
}

func (capt *Tcpdump) startSession() {
	var err error

	defer capt.out.Close()

	// Prepare tcpdump filter - by default use the values from transport
	port := capt.trans.GetRemotePort()
	if capt.filter.Port != nil {
		// unless there is an explicitly set filter
		port = capt.filter.Port
	}

	if port == nil {
		rplog.Error("Session error for %s. Can't get remote port from transport.", capt.Name())
		capt.onDie <- CapturerEvent{capt.Name(), CapturerDead}
		return
	}
	cmd := fmt.Sprintf(capt.captureCmd, *port)

	// Run capturer
	err = capt.trans.Run(cmd, capt.out, capt.pid)
	if err != nil {
		rplog.Error("Session error for %s. Can't run tcpdump command: %s", capt.Name(), err)
		capt.onDie <- CapturerEvent{capt.Name(), CapturerDead}
		return
	}

	if capt.pid.GetPid() != -1 {
		// PID is not cleared - this is unexpected stop
		capt.onDie <- CapturerEvent{capt.Name(), CapturerDead}
		rplog.Error("Session error for %s. Process died unexpectedly. Dumping stderr:\n%s",
			capt.Name(), capt.pid.DumpStdErr())
		rplog.Feedback("Capturer %s died. stderr:\n%s", capt.Name(), capt.pid.DumpStdErr())
	} else {
		rplog.Info("Session info for %s: process killed by command", capt.Name())
		capt.onDie <- CapturerEvent{capt.Name(), CapturerStopped}
	}

	return
}

// Name returns the name of the capturer's target (used only for logging purposes)
func (capt *Tcpdump) Name() string {
	return capt.name
}
