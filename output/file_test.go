/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package output

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func getTmpDir() string {
	baseDir := "/tmp/tranqap-tests"

	for i := 0; i < 10; i++ {
		id := rand.Uint32()
		dirname := fmt.Sprintf("%s-%d", baseDir, id)
		_, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			if err := os.Mkdir(dirname, 0755); err != nil {
				panic(err)
			}
			return dirname
		}
	}
	panic("Can't generate non existent directory name.")
}

func cleanup(dirname string) {
	if err := os.RemoveAll(dirname); err != nil {
		panic(err)
	}
}

func prepDestDirIsFile(t *testing.T) {
	dir := getTmpDir()
	defer cleanup(dir)

	target := "target_dir"
	path := fmt.Sprintf("%s/%s", dir, target)

	// Create a file in the target dir
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	// And pass it to prepareDestDir
	err = prepareDestDir(path)

	if err == nil {
		t.Errorf("Destination exists. Expected error, got nil")
	}
}

func prepDestDirIsDir(t *testing.T) {
	dir := getTmpDir()
	defer cleanup(dir)

	// And pass it to prepareDestDir
	err := prepareDestDir(dir)

	if err != nil {
		t.Errorf("Destination is DIR. Expected nil, got error: %s", err)
	}
}

func prepDestDirDoesntExist(t *testing.T) {
	dir := getTmpDir()
	defer cleanup(dir)

	target := "target_dir"

	// And pass it to prepareDestDir
	dest := dir + "/" + target
	err := prepareDestDir(dest)

	if err != nil {
		t.Errorf("Destination doesn't exist. Expected nil, got error: %s", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("Destination dir is not created.")
	}
}

func TestOutputPrepareDestDir(t *testing.T) {
	t.Run("prepareDestDir(): When destDir is file", prepDestDirIsFile)
	t.Run("prepareDestDir(): When destDir is dir", prepDestDirIsDir)
	t.Run("prepareDestDir(): When destDir doesn't exist", prepDestDirDoesntExist)
}

func TestOutputOpenFile(t *testing.T) {
	dir := getTmpDir()
	defer cleanup(dir)

	filepattern := "testf"
	rotationCount := 3
	loopRange := rotationCount + 2

	// On iteration 0, the destination DIR is empty, new file is created
	for i := 0; i < loopRange; i++ {
		f, err := openFile(dir, filepattern, rotationCount)
		if err != nil {
			t.Errorf("Error opening file on iteration %d: %s", i, err)
		}
		defer f.Close()

		expectedFName := fmt.Sprintf("%s/%s.pcap", dir, filepattern)
		if _, err := os.Stat(expectedFName); os.IsNotExist(err) {
			t.Errorf("Error on iteration %d: %s doesn't exist", i, expectedFName)
		}

		for j := 1; j < loopRange; j++ {
			expectedFName := fmt.Sprintf("%s/%s.%d.pcap", dir, filepattern, j)
			_, err := os.Stat(expectedFName)
			exists := !os.IsNotExist(err)

			// on i-th iteration files from 1 to i should be created
			if j <= i && j <= rotationCount {
				if exists == false {
					t.Errorf("On %d-th iteration there is no file %s", i, expectedFName)
				}
			} else if j > i || j > rotationCount {
				if exists == true {
					t.Errorf("On %d-th iteration, there is a file %s, which is wrong!", i, expectedFName)
				}
			} else {
				panic(fmt.Sprintf("Should not be here: i=%d, j=%d", i, j))
			}
		}
	}

}
