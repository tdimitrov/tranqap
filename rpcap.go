package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func main() {
	go handleSIGINT()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("rpcap> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			// This should be io.EOF
			fmt.Println() // Print an empty line not to mess the prompt of the shell
			capturers.StopAll()
			break
		}

		if processCmd(strings.TrimSuffix(command, "\n")) == cmdExit {
			capturers.StopAll()
			break
		}
	}
}

func handleSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	_ = <-c
	capturers.StopAll()

	os.Exit(1)
}
