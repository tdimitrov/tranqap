package capture

import (
	"sync"

	"golang.org/x/crypto/ssh"
)

// Capturer interface represents a general capturer. There are concrete implementations
// for tcpdump. In the future more can be added, e.g. tshark, dumpcap, etc.
type Capturer interface {
	Start() bool
	Stop() bool
}

// Storage is a thread safe container for Capturers.
type Storage struct {
	capturers []Capturer
	mut       sync.Mutex
}

// Count returns the number of Capturers in the container
func (c *Storage) Count() int {
	c.mut.Lock()
	defer c.mut.Unlock()

	return len(c.capturers)
}

// Add appends new Capturer to the container
func (c *Storage) Add(newCapt Capturer) {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.capturers = append(c.capturers, newCapt)
}

// StopAll calls Stop() on each Capturer in the container
func (c *Storage) StopAll() {
	c.mut.Lock()
	defer c.mut.Unlock()

	for _, c := range c.capturers {
		c.Stop()
	}
}

// Clear removes all Capturers from the container. Don't forget to call StopAll() before it
func (c *Storage) Clear() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.capturers = c.capturers[:0]
}

func connect(dest string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {

	client, err := ssh.Dial("tcp", dest, clientConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}
