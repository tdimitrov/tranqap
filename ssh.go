package main

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// SSHClient wraps crypto/ssh library. Destination is set during initialisation
type SSHClient struct {
	dest   string
	config ssh.ClientConfig
	client *ssh.Client
}

// NewSSHClient creates new sshClient instance
func NewSSHClient(dest string, config ssh.ClientConfig) *SSHClient {
	return &SSHClient{dest, config, nil}
}

// IsActive returns true if there is an initialised SSH client
func (c *SSHClient) IsActive() bool {
	return c.client != nil
}

// Connect initialises connection to the destination
func (c *SSHClient) Connect() error {
	var err error
	c.client, err = ssh.Dial("tcp", c.dest, &c.config)
	if err != nil {
		return err
	}

	return nil
}

// Run executes shell command synchronously
func (c *SSHClient) Run(cmd string, stdout io.Writer, stderr io.Writer) error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session: %s", err)
	}

	defer session.Close()

	session.Stdout = stdout
	session.Stderr = stderr

	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("Error running '%s': %s", cmd, err)
	}

	return nil
}
