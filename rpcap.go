package main

import (
	"flag"
	"fmt"

	"github.com/tdimitrov/rpcap/rplog"

	"github.com/abiosoft/ishell"
)

func main() {
	var configFile = flag.String("c", "config.json", "config file to use")
	var logFile = flag.String("l", "", "path to log file")

	flag.Parse()

	// Create shell
	shell := ishell.New()
	shell.SetPrompt("rpcap> ")

	// Create logger
	if len(*logFile) > 0 {
		printCb := func(f string, a ...interface{}) { shell.Printf(f, a...) }
		if err := rplog.Init(*logFile, printCb); err != nil {
			fmt.Printf("Error initialising logger: %s\nLog file won't be generated", err)
		}
	}

	// Initialise capturers storage
	initStorage()

	rplog.Info("Program started.")

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		c.Stop()
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "start file capturing",
		Func: func(ctx *ishell.Context) { cmdStart(ctx, *configFile) },
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "stop file capturing",
		Func: cmdStop,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "wireshark",
		Help: "fork wireshark for each capture",
		Func: cmdWireshark,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "targets",
		Help: "show information about loaded targets",
		Func: func(ctx *ishell.Context) { cmdTargets(ctx, *configFile) },
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "create empty config file",
		Func: func(ctx *ishell.Context) { cmdCreateConfig(ctx, *configFile) },
	})

	shell.Run()
	capturers.Close()
	rplog.Close()
}
