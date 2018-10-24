package output

import (
	"errors"
	"sync"
)

// PCAP header is a struct like this:
//
// typedef struct pcap_hdr_s {
// 	guint32 magic_number;   /* magic number */
// 	guint16 version_major;  /* major version number */
// 	guint16 version_minor;  /* minor version number */
// 	gint32  thiszone;       /* GMT to local correction */
// 	guint32 sigfigs;        /* accuracy of timestamps */
// 	guint32 snaplen;        /* max length of captured packets, in octets */
// 	guint32 network;        /* data link type */
// } pcap_hdr_t;
//
// Source: https://wiki.wireshark.org/Development/LibpcapFileFormat
//
// It should be present in the beginning of each PCAP file/stream.
// It's multiOutput's job to save the header for the stream and to put it in
// the beginning of each new stream.

const pcapHeaderSize = 32 + 2*16 + 4*32

const (
	OutputerDead = iota
)

type MultiOutputEvent struct {
	from  Outputer
	event int
}

type MOEventChan chan MultiOutputEvent
type OutputerFactory func(MOEventChan) Outputer

type MultiOutput struct {
	members    []Outputer
	membersMut sync.Mutex
	pcapHeader []byte
	events     MOEventChan
}

func NewMultiOutput(outputers ...Outputer) *MultiOutput {
	ret := &MultiOutput{outputers, sync.Mutex{}, nil, make(MOEventChan, 1)}
	go ret.eventHandler()
	return ret
}

func (mo *MultiOutput) Write(p []byte) (n int, err error) {
	// Save the header
	currHdrLen := len(mo.pcapHeader)
	if currHdrLen < pcapHeaderSize {
		mo.pcapHeader = append(mo.pcapHeader, p[0:pcapHeaderSize-currHdrLen]...)
	}

	// Forward to the capturers
	mo.membersMut.Lock()
	for _, o := range mo.members {
		o.Write(p)
	}
	mo.membersMut.Unlock()

	return len(p), nil
}

func (mo *MultiOutput) Close() {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()
	for _, o := range mo.members {
		o.Close()
	}
}

func (mo *MultiOutput) GetEventsChan() chan MultiOutputEvent {
	return mo.events
}

func (mo *MultiOutput) AddMember(newOutFn OutputerFactory) error {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()

	newMember := newOutFn(mo.events)
	if newMember == nil {
		return errors.New("Error creating Outputer with factory function")
	}

	mo.members = append(mo.members, newMember)
	return nil
}

func (mo *MultiOutput) eventHandler() {
	for event := range mo.events {
		mo.membersMut.Lock()

		for i, c := range mo.members {
			if c == event.from {
				mo.members = append(mo.members[:i], mo.members[i+1:]...)
				break
			}
		}

		mo.membersMut.Unlock()
	}

}

/*
func (mo *multiOutput) RemoveOutputer(member Outputer) error {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()

	for i := range mo.members {
		if mo.members[i] == member {
			mo.members = append(mo.members[:i], mo.members[i+1:]...)
			return nil
		}
	}

	return errors.New("Outputer not found")
}
*/

/*
	go func() {
		<-onWsharkExit
		w.Close()
		if o.RemoveOutputer(w) != nil {
			fmt.Println("Error removing wireshark outputer from multioutput")
		} else {
			fmt.Println("Wireshark closed. Removing outputer.")
		}
	}()
*/

/*
	w, onWsharkExit, err := output.NewWsharkOutput()
	if err != nil {
		fmt.Println("Can't create Wireshark output.", err)
		return cmdErr
	}
*/
