package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/tdimitrov/rpcap/rplog"
)

func main() {
	go handleSIGINT()

	if err := rplog.Init("rpcap.log"); err != nil {
		fmt.Printf("Error initialising logger: %s\nLog file won't be generated", err)
	}

	rplog.Info("Program started.")

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

		if processCmd(strings.TrimSpace(strings.TrimSuffix(command, "\n"))) == cmdExit {
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
