package main

import (
	"fmt"

	"github.com/tdimitrov/rpcap/rplog"

	"github.com/abiosoft/ishell"
)

func main() {
	shell := ishell.New()
	shell.SetPrompt("rpcap> ")

	if err := rplog.Init("rpcap.log"); err != nil {
		fmt.Printf("Error initialising logger: %s\nLog file won't be generated", err)
	}

	rplog.Info("Program started.")

	shell.EOF(func(c *ishell.Context) {
		capturers.StopAll()
		c.Stop()
	})

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		capturers.StopAll()
		c.Stop()
	})

	shell.AddCmd(&ishell.Cmd{Name: "start", Help: "start file capturing", Func: cmdStart})
	shell.AddCmd(&ishell.Cmd{Name: "stop", Help: "stop file capturing", Func: cmdStop})
	shell.AddCmd(&ishell.Cmd{Name: "wireshark", Help: "fork wireshark for each capture", Func: cmdWireshark})
	shell.AddCmd(&ishell.Cmd{Name: "targets", Help: "show information about loaded targets", Func: cmdTargets})

	shell.Run()
}
