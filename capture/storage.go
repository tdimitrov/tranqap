package capture

import (
	"fmt"
	"sync"

	"github.com/tdimitrov/rpcap/output"
)

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

// AddNewOutput adds new Outputer to each capturer
func (c *Storage) AddNewOutput(factFn output.OutputerFactory) {
	c.mut.Lock()
	defer c.mut.Unlock()

	for _, c := range c.capturers {
		err := c.AddOutputer(factFn)
		if err != nil {
			fmt.Println("Error adding Outputer")
		}
	}
}
