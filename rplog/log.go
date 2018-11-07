package rplog

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

var rpcapLog *LogFile

// Init bootstraps the logger
func Init(fname string) error {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file %s: %s", fname, err)
	}

	rpcapLog = &LogFile{f, log.New(f, "", log.LstdFlags)}

	return nil
}

func (l *LogFile) logError(format string, a ...interface{}) {
	var msgFormat strings.Builder

	if string(format[len(format)-1]) != "\n" {
		fmt.Fprintf(&msgFormat, "ERROR: %s\n", format)
	} else {
		fmt.Fprintf(&msgFormat, "ERROR: %s", format)
	}

	fmt.Fprintf(os.Stdout, msgFormat.String(), a...)
	l.logger.Printf(msgFormat.String(), a...)
}

func (l *LogFile) logInfo(format string, a ...interface{}) {
	var msgFormat strings.Builder
	fmt.Fprintf(&msgFormat, "INFO: %s", format)
	l.logger.Printf(msgFormat.String(), a...)
}

//
// Exported wrappers
//

// Error logs with prefix ERROR in file and stdout
func Error(format string, a ...interface{}) {
	if rpcapLog == nil {
		return
	}

	rpcapLog.logError(format, a...)
}

// Info logs only in file
func Info(format string, a ...interface{}) {
	if rpcapLog == nil {
		return
	}

	rpcapLog.logInfo(format, a...)
}
