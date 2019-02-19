package capture

import (
	"fmt"
	"sync"

	"github.com/tdimitrov/rpcap/output"
	"github.com/tdimitrov/rpcap/rplog"
)

// Storage is a thread safe container for Capturers.
type Storage struct {
	capturers       map[string]Capturer
	mut             sync.Mutex
	events          CapturerEventChan
	wg              sync.WaitGroup
	handlerFinished chan struct{}
}

// NewStorage creates a Storage instance
func NewStorage() *Storage {
	ret := &Storage{
		make(map[string]Capturer),
		sync.Mutex{},
		make(CapturerEventChan, 1),
		sync.WaitGroup{},
		make(chan struct{}, 1),
	}

	go ret.eventHandler()

	return ret
}

// GetChan returns the channel used for events
func (c *Storage) GetChan() CapturerEventChan {
	return c.events
}

// Add appends new Capturer to the container
func (c *Storage) Add(newCapt Capturer) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	_, exists := c.capturers[newCapt.Name()]
	if exists == true {
		return fmt.Errorf("capturer [%s] already exists", newCapt.Name())
	}

	c.capturers[newCapt.Name()] = newCapt
	c.wg.Add(1)

	return nil
}

// StopAll calls Stop() on each Capturer in the container
func (c *Storage) StopAll() {
	if c == nil {
		//Can be called on a nil pointer
		rplog.Info("capture.Storage: StopAll() called on a nil instance")
		return
	}

	rplog.Info("capture.Storage: Calling Stop for each capturer")

	c.mut.Lock()
	for _, c := range c.capturers {
		if err := c.Stop(); err != nil {
			errMsg := fmt.Sprintf("Can't stop %s. %s", c.Name(), err)
			rplog.Error(errMsg)
			rplog.Feedback(errMsg)
		}
	}
	c.mut.Unlock()

}

// Close terminates the event handler routine
func (c *Storage) Close() {
	rplog.Info("Terminating storage")
	c.StopAll()

	c.wg.Wait()
	close(c.events)

	<-c.handlerFinished

	rplog.Info("Storage terminated")
}

// AddNewOutput adds new Outputer to each capturer.
// Input parameters:
// factFn - lambda function, which accepts event channel as parameter and returns
//			new outputter. The idea is to avoid a dependency between capture and
//			output packages
// targets - slice with the names of all targets, for which the outputer should be
//			started. Empty slice means 'start for each capturer'.
func (c *Storage) AddNewOutput(factFn output.OutputerFactory, targets []string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if len(targets) == 0 {
		// wireshark cmd is called without arguments. Start outputer for each capturer

		//Check if there are running captureres
		if len(c.capturers) == 0 {
			rplog.Error("wireshark called without arguments, but no capturers are running.")
			rplog.Feedback("There are no running capturers. Use start first.\n")
			return
		}

		//Add outputer for each capturer
		for _, capt := range c.capturers {
			err := capt.AddOutputer(factFn)
			if err != nil {
				rplog.Error("Error adding Outputer to capturer %s", capt.Name())
			}
		}
	} else {
		// wireshark cmd is called for specific target(s). Start outputer only for them.
		for _, t := range targets {
			if capt, ok := c.capturers[t]; ok {
				err := capt.AddOutputer(factFn)
				if err != nil {
					rplog.Error("Error adding Outputer to capturer %s", capt.Name())
				}
			} else {
				errMsg := fmt.Sprintf("Target <%s> doesn't exist.\n", t)
				rplog.Feedback(errMsg)
				rplog.Error(errMsg)
			}
		}
	}
}

// Empty returns true if there are no Captureres in the storage
func (c *Storage) Empty() bool {
	c.mut.Lock()
	defer c.mut.Unlock()

	return len(c.capturers) == 0
}

func (c *Storage) eventHandler() {
	defer func() { c.handlerFinished <- struct{}{} }()

	rplog.Info("capture.Storage: Starting eventHandler main loop")
	for e := range c.events {
		rplog.Info("Storage: got an event from %s", e.from)
		c.mut.Lock()
		delete(c.capturers, e.from)
		rplog.Info("Storage: Removed %s", e.from)
		c.wg.Done()
		c.mut.Unlock()
	}

	rplog.Info("capture.Storage: All capturers are stopped. Exited from eventHandler main loop")
}
