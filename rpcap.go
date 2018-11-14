package main

import (
	"fmt"

	"github.com/tdimitrov/rpcap/rplog"

	"github.com/abiosoft/ishell"
)

func main() {
	// Create shell
	shell := ishell.New()
	shell.SetPrompt("rpcap> ")

	// Create logger
	printCb := func(f string, a ...interface{}) { shell.Printf(f, a...) }
	if err := rplog.Init("rpcap.log", printCb); err != nil {
		fmt.Printf("Error initialising logger: %s\nLog file won't be generated", err)
	}

	// Initialise capturers storage
	initStorage()

	rplog.Info("Program started.")

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		c.Stop()
	})

	shell.AddCmd(&ishell.Cmd{Name: "start", Help: "start file capturing", Func: cmdStart})
	shell.AddCmd(&ishell.Cmd{Name: "stop", Help: "stop file capturing", Func: cmdStop})
	shell.AddCmd(&ishell.Cmd{Name: "wireshark", Help: "fork wireshark for each capture", Func: cmdWireshark})
	shell.AddCmd(&ishell.Cmd{Name: "targets", Help: "show information about loaded targets", Func: cmdTargets})

	shell.Run()
	capturers.Close()
	rplog.Close()
}
