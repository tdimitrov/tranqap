package main

import (
	"errors"
	"sync"

	"github.com/tdimitrov/rpcap/output"
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

type multiOutput struct {
	members    []output.Outputer
	membersMut sync.Mutex
	pcapHeader []byte
}

// NewMultiOutput constructs multiOutput object
func newMultiOutput(outers ...output.Outputer) (*multiOutput, error) {
	return &multiOutput{outers, sync.Mutex{}, nil}, nil
}

func (mo *multiOutput) Write(p []byte) (n int, err error) {
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

func (mo *multiOutput) Close() {
	mo.membersMut.Lock()
	for _, o := range mo.members {
		o.Close()
	}
	mo.membersMut.Unlock()
}

func (mo *multiOutput) AddOutputer(newOut output.Outputer) error {
	mo.membersMut.Lock()
	defer mo.membersMut.Unlock()

	for i := range mo.members {
		if mo.members[i] == newOut {
			return errors.New("Outputer already added")
		}
	}

	mo.members = append(mo.members, newOut)
	return nil
}

func (mo *multiOutput) RemoveOutputer(member output.Outputer) error {
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
