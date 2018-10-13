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
		command, _ := reader.ReadString('\n')

		if processCmd(strings.TrimSuffix(command, "\n")) == cmdExit {
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
