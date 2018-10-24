package output

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type fileOutput struct {
	fd *os.File
}

// NewFileOutput constructs fileOutput object
func NewFileOutput(fname string) Outputer {
	fd, err := openFile(fname)
	if err != nil {
		return nil
	}

	return &fileOutput{fd}
}

func (pw fileOutput) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file: %v", err)
		fmt.Println(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *fileOutput) Close() {
	pw.fd.Close()
}

func openFile(filePath string) (*os.File, error) {
	const MaxFileCount int = 10

	// If file does not exist - create it and return
	if !fileExists(filePath) {
		fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return fd, nil
		}
	}

	// Split filename and extenson
	ext := path.Ext(filePath) // this includes the dot, e.g. ".pcap"
	basename := strings.Replace(filePath, ext, "", 1)

	// Remove the last file
	lastFile := fmt.Sprintf("%s.%d%s", basename, MaxFileCount, ext)
	if fileExists(lastFile) {
		err := os.Remove(lastFile)
		if err != nil {
			fmt.Printf("Error removing %v during file rotation: %v\n", lastFile, err)
			return nil, err
		}
	}

	// Shift the rest
	for n := MaxFileCount; n > 1; n-- {
		old := basename + "." + strconv.Itoa(n-1) + ext
		if fileExists(old) {
			new := basename + "." + strconv.Itoa(n) + ext
			err := os.Rename(old, new)
			if err != nil {
				fmt.Printf("Error rotating %v to %v: %v\n", old, new, err)
				continue
			}
		}
	}

	// Move the last file
	if fileExists(filePath) {
		newName := basename + "." + strconv.Itoa(1) + ext
		err := os.Rename(filePath, newName)
		if err != nil {
			fmt.Printf("Error rotating %v to %v: %v\n", filePath, newName, err)
		}
	}

	// And finally create the new file
	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("Error creating file!", err)
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
