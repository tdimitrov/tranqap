package shell

import (
	"fmt"
	"testing"
)

func getPidInstance() (chan int, getPidHandler) {
	ch := make(chan int, 1)
	inst := getPidHandler{ch}

	return ch, inst
}
func TestPidSuccess(t *testing.T) {
	ch, inst := getPidInstance()

	expectedPid := 348

	buf := []byte(fmt.Sprintf("%s%d\n", pidPrefix, expectedPid))
	inst.Write(buf)

	pid := <-ch

	if pid != expectedPid {
		t.Errorf("Expected value %d, but received %d\n", expectedPid, pid)
	}
}

func TestPidMalformedValue(t *testing.T) {
	ch, inst := getPidInstance()

	buf := []byte(fmt.Sprintf("%sgibberish\n", pidPrefix))
	inst.Write(buf)

	pid := <-ch

	if pid != -1 {
		t.Errorf("Expected value -1, but received %d\n", pid)
	}
}

func TestPidSMalformedPrefix(t *testing.T) {
	ch, inst := getPidInstance()

	expectedPid := 348

	buf := []byte(fmt.Sprintf("Gibberish:%d\n", expectedPid))
	inst.Write(buf)

	pid := <-ch

	if pid != -1 {
		t.Errorf("Expected value -1, but received %d\n", pid)
	}
}
