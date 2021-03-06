/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package capture

import (
	"testing"

	"github.com/tdimitrov/tranqap/internal/output"
)

type capturerMock struct {
	isStarted bool
	name      string
}

func newCapturerMock() *capturerMock {
	return &capturerMock{true, "Fake Capturer"}
}

func (capt *capturerMock) Start() error {
	capt.isStarted = true
	return nil
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
	if storage.capturers[capt.Name()] != capt {
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

	ch <- CapturerEvent{capt.Name(), CapturerStopped} // simulate a capturer stop
	storage.StopAll()                                 // because StopAll() will wait for event

	if capt.isStarted == true { // check if capturer is actually stopped
		t.Errorf("Saved capturer is not stopped\n")
	}

	storage.Close()

	if cnt := len(storage.capturers); cnt != 0 {
		t.Errorf("Error occurred during StopAll(). There are still %d capturers in the storage\n", cnt)
	}
}
