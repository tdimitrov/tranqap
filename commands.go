package main

import (
	"github.com/abiosoft/ishell"
	"github.com/tdimitrov/rpcap/capture"
	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/rplog"
)

const (
	cmdErr  = iota
	cmdOk   = iota
	cmdExit = iota
)

var capturers *capture.Storage

func initStorage() {
	capturers = capture.NewStorage()
}

func cmdStart(ctx *ishell.Context) {
	rplog.Info("Called start command")
	// Check if there is a running job
	if capturers.Empty() == false {
		ctx.Println("There is alreaedy a running capture")
		return
	}

	// Get configuration
	config, err := getConfig("config.json")
	if err != nil {
		ctx.Printf("Error loading configuration.\n", err)
		return
	}

	for _, t := range config.Targets {
		c, d, err := getClientConfig(&t)
		if err != nil {
			ctx.Printf("Error parsing client configuration for target <%s>: %s\n", *t.Name, err)
			return
		}

		// Create file output
		f := output.NewFileOutput(*t.Destination, *t.FilePattern, *t.RotationCnt)
		if f == nil {
			ctx.Printf("Can't create File output for target <%s>\n", *t.Name)
			return
		}

		// Create multioutput and attach the file output to it
		m := output.NewMultiOutput(f)
		if m == nil {
			ctx.Printf("Can't create MultiOutput for target <%s>\n.", *t.Name)
			return
		}

		// Create capturer
		sshClient := NewSSHClient(*d, *c)
		capt := capture.NewTcpdump(*t.Name, m, capturers.GetChan(), sshClient)
		if capt == nil {
			ctx.Printf("Error creating Capturer for target <%s>\n", *t.Name)
			return
		}

		if capt.Start() == false {
			ctx.Printf("Error starting Capturer for target <%s>\n", *t.Name)
			return
		}

		capturers.Add(capt)
	}
}

func cmdStop(ctx *ishell.Context) {
	rplog.Info("Called stop command")

	// Check if there is a running job
	if capturers.Empty() == true {
		ctx.Println("There are no running captures.")
		return
	}

	capturers.StopAll()
}

func cmdWireshark(ctx *ishell.Context) {
	rplog.Info("Called wireshark command")

	// Prepare a factory function, which creates Wireshark Outputer
	factFn := func(p output.MOEventChan) output.Outputer {
		return output.NewWsharkOutput(p)
	}

	capturers.AddNewOutput(factFn)
}

func cmdTargets(ctx *ishell.Context) {
	rplog.Info("Called targets command")

	// Get configuration
	config, err := getConfig("config.json")
	if err != nil {
		ctx.Printf("Error loading configuration: %s\n", err)
		return
	}

	for _, t := range config.Targets {
		c, d, err := getClientConfig(&t)
		if err != nil {
			ctx.Printf("Error parsing client configuration for target <%s>: %s\n", *t.Name, err)
			return
		}

		ctx.Printf("=== Running checks for target <%s> ===\n", *t.Name)
		sshClient := NewSSHClient(*d, *c)
		if output, err := checkPermissions(sshClient); err != nil {
			ctx.Printf("%s\n", err)
			return
		} else {
			ctx.Printf("%s\n", output)
		}
	}

	return
}
