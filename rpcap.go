package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("rpcap> ")
		command, _ := reader.ReadString('\n')

		if processCmd(strings.TrimSuffix(command, "\n")) == cmdExit {
			break
		}
	}
}
