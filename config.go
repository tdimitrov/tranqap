package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type target struct {
	Host *string
	Port *int
	User *string
	Pass *string
	Key  *string
}

func getTarget(fname string) (target, error) {
	var t target
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

func getClientConfig(t *target) (*ssh.ClientConfig, *string, error) {
	var clientConfig ssh.ClientConfig

	clientConfig.Auth = make([]ssh.AuthMethod, 0, 2)

	if t.User == nil {
		return nil, nil, errors.New("Missing user in configuration")
	}

	if t.Pass == nil && t.Key == nil {
		return nil, nil, errors.New("Missing authentication method in configuration - provide Password or/and private key")
	}

	if t.Host == nil {
		return nil, nil, errors.New("Missing host in configuration")
	}

	if t.Port == nil {
		return nil, nil, errors.New("Missing port in configuration")
	}

	dest := fmt.Sprintf("%s:%d", *t.Host, *t.Port)

	clientConfig.User = *t.User
	clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	if t.Pass != nil {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(*t.Pass))
	}

	if t.Key != nil {
		key, err := ioutil.ReadFile(*t.Key)
		if err != nil {
			msg := fmt.Sprintf("unable to read private key: %v", err)
			return nil, nil, errors.New(msg)
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			msg := fmt.Sprintf("unable to parse private key: %v", err)
			return nil, nil, errors.New(msg)
		}

		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	return &clientConfig, &dest, nil
}
