package shell

import "testing"

func getKillInstance() (chan int, killPidHandler) {
	ch := make(chan int, 1)
	inst := killPidHandler{ch}

	return ch, inst
}
func TestKillSuccess(t *testing.T) {
	ch, inst := getKillInstance()

	buf := []byte("0 15\n")
	inst.Write(buf)

	r := <-ch
	if EvKillSuccess != r {
		t.Errorf("Expected success, but received %s\n", KillResToStr(r))
	}

}

func TestPidNotResponding(t *testing.T) {
	ch, inst := getKillInstance()

	buf := []byte("0 0\n")
	inst.Write(buf)

	r := <-ch
	if EvKillNotResponding != r {
		t.Errorf("Expected not responding, but received %s\n", KillResToStr(r))
	}

}

func TestPidNotRunning(t *testing.T) {
	ch, inst := getKillInstance()

	buf := []byte("13 11\n")
	inst.Write(buf)

	r := <-ch
	if EvKillNotRuning != r {
		t.Errorf("Expected not running, but received %s\n", KillResToStr(r))
	}

}

func TestKillMalformed(t *testing.T) {
	ch, inst := getKillInstance()

	buf := []byte("1311\n")
	inst.Write(buf)

	r := <-ch
	if EvKillExecError != r {
		t.Errorf("Expected exec error, but received %s\n", KillResToStr(r))
	}

}
