package main

import (
	"github.com/tdimitrov/rpcap/capture"
	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/shell"
)

const (
	cmdErr  = iota
	cmdOk   = iota
	cmdExit = iota
)

var capturers *capture.Storage

func cmdStart() int {
	// Check if there is a running job
	if capturers != nil {
		logger.Error("There is already a running capture.")
		return cmdErr
	}

	capturers = capture.NewStorage()

	// Get configuration
	config, err := getConfig("config.json")
	if err != nil {
		logger.Error("Error loading configuration. ", err)
		return cmdErr
	}

	for _, t := range config.Targets {
		c, d, err := getClientConfig(&t)
		if err != nil {
			logger.Error("Error parsing client configuration for target <%s>: %s\n", *t.Name, err)
			return cmdErr
		}

		// Create file output
		f := output.NewFileOutput(*t.Destination, *t.FilePattern, *t.RotationCnt)
		if f == nil {
			logger.Error("Can't create File output for target <%s>\n", *t.Name)
			return cmdErr
		}

		// Create multioutput and attach the file output to it
		m := output.NewMultiOutput(f)
		if m == nil {
			logger.Error("Can't create MultiOutput for target <%s\n.", *t.Name)
			return cmdErr
		}

		// Create capturer
		capt := capture.NewTcpdump(*d, c, m, capturers.GetChan())
		if capt == nil {
			logger.Error("Error creating Capturer for target <%s>\n", *t.Name)
			return cmdErr
		}

		if capt.Start() == false {
			logger.Error("Error starting Capturer for target <%s>\n", *t.Name)
			return cmdErr
		}

		capturers.Add(capt)
	}

	return cmdOk
}

func cmdStop() int {
	// Check if there is a running job
	if capturers == nil {
		logger.Error("There are no running captures.")
		return cmdErr
	}

	capturers.StopAll()

	capturers = nil

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

func cmdCheckTargets() int {
	// Get configuration
	config, err := getConfig("config.json")
	if err != nil {
		logger.Error("Error loading configuration. ", err)
		return cmdErr
	}

	for _, t := range config.Targets {
		c, d, err := getClientConfig(&t)
		if err != nil {
			logger.Error("Error parsing client configuration for target <%s>: %s\n", *t.Name, err)
			return cmdErr
		}

		logger.Error("=== Running checks for target <%s> ===\n", *t.Name)
		if shell.CheckPermissions(c, *d) == false {
			return cmdErr
		}
		logger.Error("=========================")
	}

	return cmdOk
}

func processCmd(cmd string) int {
	logger.Info("Executed %s command", cmd)
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

	case "check targets":
		return cmdCheckTargets()

	case "":
		return cmdOk

	default:
		logger.Error("No such command %s", cmd)
		return cmdErr
	}
}
