/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package output

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tdimitrov/rpcap/rplog"
)

type fileOutput struct {
	fd *os.File
}

// NewFileOutput constructs fileOutput object
func NewFileOutput(destDir string, filePattern string, rotationCnt int) Outputer {
	fd, err := openFile(destDir, filePattern, rotationCnt)
	if err != nil {
		return nil
	}

	return &fileOutput{fd}
}

func (pw fileOutput) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file: %v", err)
		rplog.Info(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *fileOutput) Close() {
	pw.fd.Close()
}

func openFile(destDir string, filePattern string, rotationCnt int) (*os.File, error) {
	// If destination dir doesn't exist - create it
	err := prepareDestDir(destDir)
	if err != nil {
		return nil, err
	}

	filePath := destDir + "/" + filePattern + ".pcap"

	// If file does not exist - create it and return
	if !fileExists(filePath) {
		fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			rplog.Error("Error opening file: %s", err)
			return nil, err
		}
		return fd, err
	}

	// Split filename and extenson
	ext := path.Ext(filePath) // this includes the dot, e.g. ".pcap"
	basename := strings.Replace(filePath, ext, "", 1)

	// Remove the last file
	lastFile := fmt.Sprintf("%s.%d%s", basename, rotationCnt, ext)
	if fileExists(lastFile) {
		err := os.Remove(lastFile)
		if err != nil {
			rplog.Error("Error removing %v during file rotation: %v\n", lastFile, err)
			return nil, err
		}
	}

	// Shift the rest
	for n := rotationCnt; n > 1; n-- {
		old := basename + "." + strconv.Itoa(n-1) + ext
		if fileExists(old) {
			new := basename + "." + strconv.Itoa(n) + ext
			err := os.Rename(old, new)
			if err != nil {
				rplog.Error("Error rotating %v to %v: %v\n", old, new, err)
				continue
			}
		}
	}

	// Move the last file
	if fileExists(filePath) {
		newName := basename + "." + strconv.Itoa(1) + ext
		err := os.Rename(filePath, newName)
		if err != nil {
			rplog.Error("Error rotating %v to %v: %v\n", filePath, newName, err)
		}
	}

	// And finally create the new file
	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		rplog.Error("Error creating file: %s", err)
		return nil, err
	}

	return fd, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

func prepareDestDir(destDir string) error {
	fi, err := os.Stat(destDir)

	if os.IsNotExist(err) {
		// Destination doesn't exist - create it
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return err
		}

		return nil
	}

	if fi.IsDir() {
		// Destination exists and is a DIR
		return nil
	}

	return fmt.Errorf("Destination dir path (%s) points to a file", destDir)
}
