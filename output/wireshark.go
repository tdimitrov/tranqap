package output

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type wsharkOutput struct {
	pipeFd *os.File
	pid    int
}

// NewWsharkOutput constructs wsharkOutput object
func NewWsharkOutput(pipeFname string) (Outputer, error) {
	// First create the pipe
	err := syscall.Mknod(pipeFname, syscall.S_IFIFO|0666, 0)
	if err != nil {
		return nil, err
	}

	// Then open it
	fd, err := os.OpenFile(pipeFname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	// And finally fork wireshark
	// TODO: fork wireshark
	pid, err := forkWireshark()
	if err != nil {
		return nil, err
	}

	return &wsharkOutput{fd, pid}, nil
}

func (pw wsharkOutput) Write(p []byte) (n int, err error) {
	n, err = pw.pipeFd.Write(p)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file: %v", err)
		fmt.Println(msg)
		return n, errors.New(msg)
	}
	return n, nil
}

func (pw *wsharkOutput) Close() {
	pw.pipeFd.Close()
}

func forkWireshark() (int, error) {
	const bin = "wireshark"

	cmd := exec.Command("wireshark", "-k", "-i", "/tmp/test.pipe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return -1, err
	}

	fmt.Println(cmd.Process.Pid)

	return cmd.Process.Pid, nil
}

func forkWireshark1() (int, error) {
	const bin = "wireshark"

	binary, err := exec.LookPath(bin)
	if err != nil {
		return -1, err
	}

	// From execve(2) manpage: By convention, the first of these strings
	// (i.e., argv[0])  should  contain the filename associated with the
	// file being executed.
	args := []string{bin, "-k", "-i"}
	attr := syscall.ProcAttr{"/tmp", os.Environ(), []uintptr{}, nil}

	//err = syscall.Exec(, env)
	pid, err := syscall.ForkExec(binary, args, &attr)
	if err != nil {
		return -1, err
	}

	return pid, nil
}
