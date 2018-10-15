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

var capturers capture.Storage

func cmdStart() int {
	// Check if there is a running job
	if capturers.Count() != 0 {
		fmt.Println("There is already a running capture.")
		return cmdErr
	}

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

	// Create outputers
	f, err := output.NewFileOutput("test.pcap")
	if err != nil {
		fmt.Println("Can't create File output.", err)
		return cmdErr
	}

	w, err := output.NewWsharkOutput()
	if err != nil {
		fmt.Println("Can't create Wireshark output.", err)
		return cmdErr
	}

	o, err := newMultiOutput(f, w)
	if err != nil {
		fmt.Println("Can't create multi output.", err)
		return cmdErr
	}

	// Create capturer
	capt := capture.NewTcpdump(*d, c, o)
	if capt == nil {
		fmt.Println("Error creating capturer.")
		return cmdErr
	}

	if capt.Start() == false {
		fmt.Println("Error starting capture")
		return cmdErr
	}

	capturers.Add(capt)

	return cmdOk
}

func cmdStop() int {
	// Check if there is a running job
	if capturers.Count() == 0 {
		fmt.Println("There are no running captures.")
		return cmdErr
	}

	capturers.StopAll()
	capturers.Clear()

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

	case "stop":
		return cmdStop()

	default:
		fmt.Println("No such command", cmd)
		return cmdErr
	}
}
