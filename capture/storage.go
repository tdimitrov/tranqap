package capture

import (
	"fmt"
	"sync"

	"github.com/tdimitrov/rpcap/rplog"

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
	ret := &Storage{[]Capturer{}, sync.Mutex{}, make(CapturerEventChan, 1), make(chan struct{}, 1)}
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
	if c == nil {
		//Can be called on a nil pointer
		rplog.Info("capture.Storage: StopAll() called on a nil instance")
		return
	}

	rplog.Info("capture.Storage: Stopping all capturers")

	c.mut.Lock()
	for _, c := range c.capturers {
		c.Stop()
	}
	c.mut.Unlock()

	rplog.Info("capture.Storage: Waiting for confirmation")
	<-c.handlerDone

	rplog.Info("capture.Storage: All capturers are now stopped")
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
	rplog.Info("capture.Storage: Starting eventHandler main loop")
	for e := range c.events {
		rplog.Info("capture.Storage: eventHandler got an event from a capturer")
		c.mut.Lock()
		for i, capt := range c.capturers {
			if capt == e.from {
				rplog.Info("capture.Storage: Removed capturer")
				c.capturers = append(c.capturers[:i], c.capturers[i+1:]...)
				break
			}
		}
		capturersCount := len(c.capturers)
		c.mut.Unlock()

		if capturersCount == 0 {
			break
		}
	}

	rplog.Info("capture.Storage: All capturers are stopped. Exited from eventHandler main loop")
	c.handlerDone <- struct{}{}
}
