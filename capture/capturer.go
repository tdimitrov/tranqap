package capture

import (
	"github.com/tdimitrov/rpcap/output"
	"golang.org/x/crypto/ssh"
)

// Capturer interface represents a general capturer. There are concrete implementations
// for tcpdump. In the future more can be added, e.g. tshark, dumpcap, etc.
type Capturer interface {
	Start() bool
	Stop() bool
	AddOutputer(newOutputer output.OutputerFactory) error
}

func connect(dest string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {

	client, err := ssh.Dial("tcp", dest, clientConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}
