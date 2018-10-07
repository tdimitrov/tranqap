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

var capturers []capture.Capturer

func cmdStart() int {
	// Check if there is a running job
	if len(capturers) != 0 {
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

	// Create outputer
	//o, err := output.NewFileOutput("test.pcap")
	o, err := output.NewWsharkOutput("/tmp/test.pipe")
	if err != nil {
		fmt.Println("Can't create output.", err)
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

	capturers = append(capturers, capt)

	return cmdOk
}

func cmdStop() int {
	// Check if there is a running job
	if len(capturers) == 0 {
		fmt.Println("There are no running captures.")
		return cmdErr
	}

	// Stop all capturers
	for _, c := range capturers {
		c.Stop()
	}

	// Clear the slice
	capturers = capturers[:0]

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
	}

	return cmdErr
}
