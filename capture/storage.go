package capture

import (
	"fmt"
	"sync"

	"github.com/tdimitrov/rpcap/output"
)

// Storage is a thread safe container for Capturers.
type Storage struct {
	capturers   []Capturer
	mut         sync.Mutex
	events      CapturerEventChan
	handlerDone chan struct{}
}

// NewStorage creates a Storage instance
func NewStorage() *Storage {
	ret := &Storage{[]Capturer{}, sync.Mutex{}, make(CapturerEventChan, 1), make(chan struct{})}
	go ret.eventHandler()

	return ret
}

// GetChan returns the channel used for events
func (c *Storage) GetChan() CapturerEventChan {
	return c.events
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

	close(c.events)

	<-c.handlerDone
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

func (c *Storage) eventHandler() {
	for e := range c.events {
		c.mut.Lock()

		for i, capt := range c.capturers {
			if capt == e.from {
				c.capturers = append(c.capturers[:i], c.capturers[i+1:]...)
				break
			}
		}
		c.mut.Unlock()
	}

	c.handlerDone <- struct{}{}
}
