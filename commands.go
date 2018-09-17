package main

import (
	"fmt"

	"github.com/tdimitrov/rpcap/capture"
	"github.com/tdimitrov/rpcap/output"
)

const (
	cmdErr  = iota
	cmdOk   = iota
	cmdExit = iota
)

func cmdStart() int {
	// Get configuration
	t, err := getTarget("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return cmdErr
	}

	// Create SSH client config from configuration
	c, err := getClientConfig(&t)
	if err != nil {
		fmt.Println("Error parsing client configuration: ", err)
		return cmdErr
	}

	// Get destination from configuration
	d := getDest(&t)
	if err != nil {
		fmt.Println("Error parsing destination: ", err)
		return cmdErr
	}

	// Create capturer
	capt := capture.NewTcpdump(d, c)

	return cmdOk
}

func processCmd(cmd string) int {
	switch cmd {
	case "exit":
		return cmdExit

	case "quit":
		return cmdExit

	case "start":
		return cmdStart()
	}

	return cmdErr
}
