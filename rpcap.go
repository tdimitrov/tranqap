package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	cmdErr  = iota
	cmdOk   = iota
	cmdExit = iota
)

func getTarget(fname string) (Target, error) {
	var t Target
	var err error

	confFile, err := ioutil.ReadFile(fname)
	if err != nil {
		msg := fmt.Sprintf("Error opening %s: %s\n", fname, err.Error())
		return t, errors.New(msg)
	}

	err = json.Unmarshal(confFile, &t)
	if err != nil {
		msg := fmt.Sprintf("Error parsing %s: %s\n", fname, err.Error())
		return t, errors.New(msg)
	}

	return t, nil
}

func connect(t *Target) (*ssh.Client, error) {
	var clientConfig ssh.ClientConfig

	clientConfig.Auth = make([]ssh.AuthMethod, 0, 2)

	if t.User == nil {
		return nil, errors.New("Missing user in configuration")
	}

	if t.Pass == nil && t.Key == nil {
		return nil, errors.New("Missing authentication method in configuration - provide Password or/and private key")
	}

	if t.Host == nil {
		return nil, errors.New("Missing host in configuration")
	}

	if t.Port == nil {
		return nil, errors.New("Missing port in configuration")
	}

	clientConfig.User = *t.User
	clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	if t.Pass != nil {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(*t.Pass))
	}

	if t.Key != nil {
		key, err := ioutil.ReadFile(*t.Key)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}

		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	dest := fmt.Sprintf("%s:%d", *t.Host, *t.Port)

	client, err := ssh.Dial("tcp", dest, &clientConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

func capture() {
	t, err := getTarget("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}

	client, err := connect(&t)
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	//fmt.Println(client.LocalAddr().(*net.TCPAddr).IP)

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Error creating session!")
		return
	}

	defer session.Close()

	writer := pcapwirter{}
	writer.Init("test.fifo")
	defer writer.DeInit()

	session.Stdout = writer

	err = session.Start("sudo tcpdump -U -s0 -w - 'ip and not port 22'")
	if err != nil {
		fmt.Println("Error running command!")
		return
	}

	session.Wait()
}

func processCmd(cmd string) int {
	switch cmd {
	case "exit":
		return cmdExit

	case "quit":
		return cmdExit
	}

	return cmdErr
}

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
