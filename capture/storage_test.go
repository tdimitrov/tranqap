package capture

import (
	"testing"

	"github.com/tdimitrov/rpcap/output"
)

type capturerMock struct {
	isStarted bool
	name      string
}

func newCapturerMock() *capturerMock {
	return &capturerMock{true, "Fake Capturer"}
}

func (capt *capturerMock) Start() bool {
	capt.isStarted = true
	return true
}

func (capt *capturerMock) Stop() error {
	capt.isStarted = false
	return nil
}

func (capt *capturerMock) AddOutputer(newOutputer output.OutputerFactory) error {
	return nil
}

func (capt capturerMock) Name() string {
	return capt.name
}

func TestStorageAdd(t *testing.T) {
	storage := NewStorage()

	// Initially storage should be empty
	if len(storage.capturers) != 0 {
		t.Errorf("Invalid size of storage.capturers. Expected 0, got %d\n", len(storage.capturers))
	}

	// Add one
	capt := newCapturerMock()
	storage.Add(capt)

	// Check if it is saved
	if len(storage.capturers) != 1 {
		t.Errorf("Invalid size of storage.capturers. Expected 1, got %d\n", len(storage.capturers))
	}

	// Check if it is the same
	if storage.capturers[0] != capt {
		t.Errorf("Saved capturer doesn't match with the created capturer\n")
	}
}

func TestStorageStopAll(t *testing.T) {
	storage := NewStorage()

	// Initially storage should be empty
	if len(storage.capturers) != 0 {
		t.Errorf("Invalid size of storage.capturers. Expected 0, got %d\n", len(storage.capturers))
	}

	// Add one
	capt := newCapturerMock()
	storage.Add(capt)

	// Stop all
	ch := storage.GetChan() // on this chan events from the Capturers are received

	ch <- CapturerEvent{capt, CapturerStopped} // simulate a capturer stop
	storage.StopAll()                          // because StopAll() will wait for event

	if capt.isStarted == true { // check if capturer is actually stopped
		t.Errorf("Saved capturer is not stopped\n")
	}

	storage.Close()

	if cnt := len(storage.capturers); cnt != 0 {
		t.Errorf("Error occured during StopAll(). There are still %d capturers in the storage\n", cnt)
	}
}
