package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tdimitrov/rpcap/rplog"

	"github.com/abiosoft/ishell"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "rpcap [global flags] [subcommand [subcommand flags]]\n\n")
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSubcommands:\n")
		fmt.Fprintf(os.Stderr, "init - creates sample configuration file. Works with -c.\n")
		fmt.Fprintf(os.Stderr, "\tE.g. \"%s -c config.json init\" - ", os.Args[0])
		fmt.Fprintf(os.Stderr, "creates sample config named config.json in current working directory.\n")
	}

	var configFile = flag.String("c", "config.json", "config file to use")
	var logFile = flag.String("l", "", "path to log file")

	flag.Parse()

	if len(flag.Args()) > 0 {
		//subcommand
		if len(flag.Args()) == 1 && flag.Arg(0) == "init" {
			//init cmd
			rplog.Info("Called config command")
			err := generateSampleConfig(*configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating sample config: %s\n", err.Error())
			} else {
				fmt.Printf("Saved sample configuration to %s\n", *configFile)
			}
			return
		}

		//bad cmd
		fmt.Fprintf(os.Stderr, "Bad subcommand: %v\n", flag.Args())
		flag.Usage()
		return
	}

	// Get configuration
	config, err := getConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration. %s\n", err)
		fmt.Fprintf(os.Stderr, "Run %s init or provide path to a configuration file with -c.\n", os.Args[0])
		return
	}
	targetsList := config.getTargetsList()

	// Create shell
	shell := ishell.New()
	shell.SetPrompt("rpcap> ")

	// Create logger
	printCb := func(f string, a ...interface{}) { shell.Printf(f, a...) }
	if err := rplog.Init(*logFile, printCb); err != nil {
		fmt.Printf("Error initialising logger: %s\nLog file won't be generated", err)
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
		Func: func(ctx *ishell.Context) { cmdStart(ctx, config) },
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
		Completer: func([]string) []string {
			return targetsList
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "targets",
		Help: "show information about loaded targets",
		Func: func(ctx *ishell.Context) { cmdTargets(ctx, config) },
	})

	shell.Run()
	capturers.Close()
	rplog.Close()
}
