package capture

import (
	"fmt"
	"io"
	"testing"

	"github.com/tdimitrov/rpcap/output"
)

type outputMock struct {
	isClosed bool
}

func (out *outputMock) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (out *outputMock) Close() {
	out.isClosed = true
}

type transportMock struct {
	active        bool
	failOnRun     bool
	failOnConnect bool
	finish        chan struct{}
}

func (trans *transportMock) IsActive() bool {
	return trans.active
}

func (trans *transportMock) Connect() error {
	if trans.failOnConnect == true {
		return fmt.Errorf("Something went wrong")
	}

	trans.active = true
	return nil
}

func (trans *transportMock) GetRemoteIP() *string {
	ret := "127.0.0.1"
	return &ret
}

func (trans *transportMock) GetRemotePort() *int {
	ret := 22
	return &ret
}

func (trans *transportMock) Run(cmd string, stdout io.Writer, stderr io.Writer) error {
	if trans.failOnRun == true {
		return fmt.Errorf("Something went wrong")
	}

	<-trans.finish
	return nil
}

func createTestInstances() (CapturerEventChan, *transportMock, Capturer, *outputMock) {
	events := make(CapturerEventChan)
	trans := transportMock{false, false, false, make(chan struct{}, 1)}
	out := &outputMock{false}

	inst := NewTcpdump("Test Instance", output.NewMultiOutput(out), events, &trans, SudoConfig{false, nil})

	return events, &trans, inst, out
}

func TestTcpdumpStop(t *testing.T) {
	events, trans, inst, out := createTestInstances()

	if inst.Start() != nil {
		t.Errorf("Unexpected Start() failure\n")
	}

	// Set PID to invalid value, which means that Stop() command was issued
	*inst.(*Tcpdump).pid.pid = -1

	trans.finish <- struct{}{}
	ev := <-events

	if ev.from != inst.Name() {
		t.Errorf("Got event with wrong from")
	}

	if ev.event != CapturerStopped {
		t.Errorf("Got wrong event type")
	}

	if out.isClosed == false {
		t.Errorf("Outputer is not closed")
	}
}

func TestTcpdumpDie(t *testing.T) {
	events, trans, inst, out := createTestInstances()

	if inst.Start() != nil {
		t.Errorf("Unexpected Start() failure\n")
	}

	// Set PID to valid value, which means that the PID was not cleared by Stop() command
	*inst.(*Tcpdump).pid.pid = 10

	trans.finish <- struct{}{}
	ev := <-events

	if ev.from != inst.Name() {
		t.Errorf("Got event with wrong from")
	}

	if ev.event != CapturerDead {
		t.Errorf("Got wrong event type")
	}

	if out.isClosed == false {
		t.Errorf("Outputer is not closed")
	}
}

func TestTcpdumpFailOnRun(t *testing.T) {
	events, trans, inst, out := createTestInstances()

	trans.failOnRun = true

	if inst.Start() != nil {
		t.Errorf("Unexpected Start() failure\n")
	}

	ev := <-events

	if ev.from != inst.Name() {
		t.Errorf("Got event with wrong from")
	}

	if ev.event != CapturerDead {
		t.Errorf("Got wrong event type")
	}

	if out.isClosed == false {
		t.Errorf("Outputer is not closed")
	}
}

func TestTcpdumpFailOnConnect(t *testing.T) {
	_, trans, inst, out := createTestInstances()

	trans.failOnConnect = true

	if inst.Start() == nil {
		t.Errorf("Start() should return error\n")
	}

	if out.isClosed == false {
		t.Errorf("Outputer is not closed")
	}
}
