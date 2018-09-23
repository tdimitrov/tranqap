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
		fmt.Println("Error loading configuration. ", err)
		return cmdErr
	}

	// Create SSH client config and destination from configuration
	c, d, err := getClientConfig(&t)
	if err != nil {
		fmt.Println("Error parsing client configuration. ", err)
		return cmdErr
	}

	// Create outputer
	o, err := output.NewPcapOutput("test.pcap")
	if err != nil {
		fmt.Println("Can't create PCAP output.", err)
		return cmdErr
	}

	// Create capturer
	capt := capture.NewTcpdump(*d, c, []output.Outputer{o})
	if capt == nil {
		fmt.Println("Error creating capturer.")
		return cmdErr
	}

	if capt.Start() == false {
		fmt.Println("Error starting capture")
		return cmdErr
	}

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
