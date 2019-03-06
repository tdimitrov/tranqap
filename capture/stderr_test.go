/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package capture

import (
	"fmt"
	"testing"
)

func TestPidSuccess(t *testing.T) {
	inst := newStdErrHandler()

	expectedPid := 348

	buf := []byte(fmt.Sprintf("%s%d\n", pidPrefix, expectedPid))
	inst.Write(buf)

	pid := inst.GetPid()

	if pid != expectedPid {
		t.Errorf("Expected value %d, but received %d\n", expectedPid, pid)
	}
}

func TestPidMalformedValue(t *testing.T) {
	inst := newStdErrHandler()

	buf := []byte(fmt.Sprintf("%sgibberish\n", pidPrefix))
	inst.Write(buf)

	pid := inst.GetPid()

	if pid != -1 {
		t.Errorf("Expected value -1, but received %d\n", pid)
	}
}

func TestPidSMalformedPrefix(t *testing.T) {
	inst := newStdErrHandler()
	expectedPid := 348

	buf := []byte(fmt.Sprintf("Gibberish:%d\n", expectedPid))
	inst.Write(buf)

	pid := inst.GetPid()

	if pid != -1 {
		t.Errorf("Expected value -1, but received %d\n", pid)
	}
}
