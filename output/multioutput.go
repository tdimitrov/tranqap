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

const pcapHeaderSize = 24 // From the struct above: (32 + 2*16 + 4*32) / 8

const (
	// OutputerDead is generated to the MultiOutput when
	// the Outputer process (e.g. Wireshark) dies
	OutputerDead = iota
)

// MultiOutputEvent represents the structure of the event generated from Outputer
// to MultiOtuput. It has got two parameters:
// from - the address of the Outputer struct in memory. It is used to identify the Outputer
// event - the type of the event. This value should be equal on one of the consts above.
type MultiOutputEvent struct {
	from  Outputer
	event int
}

// MOEventChan is the type of the channel used by MultiOutput for event handling
type MOEventChan chan MultiOutputEvent

// OutputerFactory is a function which creates new Outputer. It receives one parameter
// of type MOEventChan.
// The purpose is tha have a factory function which creates Outputer and passes to it the
// event handling channel of the MultiOutput. This way MultiOutput can create Outputers
// without knowing anything about their creation process.
type OutputerFactory func(MOEventChan) Outputer

// MultiOutput redirects PCAP traffic to multiple outputers, which are saved
// in the members slice.
// It also saves the pcapHeader, received at the start of the capturing, so that
// the header can be reinjected when an outputer is restarted.
type MultiOutput struct {
	members    []Outputer
	membersMut sync.Mutex
	pcapHeader []byte
	events     MOEventChan
}

// NewMultiOutput create new MultiOutput instance. The function receives one or more
// Outputers as input parameters, which are added to the members slice.
func NewMultiOutput(outputers ...Outputer) *MultiOutput {
	ret := &MultiOutput{outputers, sync.Mutex{}, nil, make(MOEventChan, 1)}
	go ret.eventHandler()
	return ret
}

// Write delivers PCAP traffic to all Outputers. It also saves the pcap header.
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

// Close closes all member Outputers
func (mo *MultiOutput) Close() {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()
	for _, o := range mo.members {
		o.Close()
	}
}

// GetEventsChan returns the channel used by the MultiOutput instance for event handling
func (mo *MultiOutput) GetEventsChan() chan MultiOutputEvent {
	return mo.events
}

// AddMember adds new Outputer to the members slice
func (mo *MultiOutput) AddMember(newOutFn OutputerFactory) error {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()

	// Create new member
	newMember := newOutFn(mo.events)
	if newMember == nil {
		return errors.New("Error creating Outputer with factory function")
	}

	// Send the PCAP header
	newMember.Write(mo.pcapHeader)

	// Add to members list
	mo.members = append(mo.members, newMember)
	return nil
}

// eventHandler handles events from member Outputers
// Effectively at the moment this function just removes dead Outputers from
// the members slice
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
