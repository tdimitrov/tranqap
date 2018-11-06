package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogFile is based on log package. Supports log levels and printing messages to stdout
type LogFile struct {
	file   *os.File
	logger *log.Logger
}

// NewLogFile creates instance of LogFile
func NewLogFile(fname string) *LogFile {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("Error opening log file %s: %s\n", fname, err)
		return nil
	}

	logger := log.New(f, "", log.LstdFlags)

	return &LogFile{f, logger}
}

// Error logs with prefix ERROR in file and stdout
func (l *LogFile) Error(format string, a ...interface{}) {
	var msgFormat strings.Builder

	if string(format[len(format)-1]) != "\n" {
		fmt.Fprintf(&msgFormat, "ERROR: %s\n", format)
	} else {
		fmt.Fprintf(&msgFormat, "ERROR: %s", format)
	}

	fmt.Fprintf(os.Stdout, msgFormat.String(), a...)
	l.logger.Printf(msgFormat.String(), a...)
}

// Info logs only in file
func (l *LogFile) Info(format string, a ...interface{}) {
	var msgFormat strings.Builder
	fmt.Fprintf(&msgFormat, "INFO: %s", format)
	l.logger.Printf(msgFormat.String(), a...)
}
