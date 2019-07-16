/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"io"
	"net"

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

// GetRemoteIP returns the IP address of the target which the SSH connection uses.
// It is useful for the cases when a hostname is specified in the configuration.
// For such situatuons the exact IP address is needed for the capture filter.
func (c *SSHClient) GetRemoteIP() *string {
	if c.client == nil {
		return nil
	}

	if addr := c.client.RemoteAddr(); addr != nil {
		ret := c.client.LocalAddr().(*net.TCPAddr).IP.String()
		return &ret
	}

	return nil
}

// GetRemotePort returns the port number of the SSH target
func (c *SSHClient) GetRemotePort() *int {
	if c.client == nil {
		return nil
	}

	if addr := c.client.RemoteAddr(); addr != nil {
		return &c.client.LocalAddr().(*net.TCPAddr).Port
	}

	return nil
}
