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

	// Create file output
	f := output.NewFileOutput("test.pcap")
	if f == nil {
		fmt.Println("Can't create File output.")
		return cmdErr
	}

	// Create multioutput and attach the file output to it
	m := output.NewMultiOutput(f)
	if m == nil {
		fmt.Println("Can't create Multi output.")
		return cmdErr
	}

	// Create capturer
	capt := capture.NewTcpdump(*d, c, m)
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

func cmdWireshark() int {
	// Prepare a factory function, which creates Wireshark Outputer
	factFn := func(p output.MOEventChan) output.Outputer {
		return output.NewWsharkOutput(p)
	}

	capturers.AddNewOutput(factFn)

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

	case "wireshark":
		return cmdWireshark()

	default:
		fmt.Println("No such command", cmd)
		return cmdErr
	}
}
