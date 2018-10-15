package output

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
	members    []Outputer
	pcapHeader []byte
}

// NewMultiOutput constructs multiOutput object
func NewMultiOutput(outers ...Outputer) (Outputer, error) {
	return &multiOutput{outers, nil}, nil
}

func (mo *multiOutput) Write(p []byte) (n int, err error) {
	// Save the header
	currHdrLen := len(mo.pcapHeader)
	if currHdrLen < pcapHeaderSize {
		mo.pcapHeader = append(mo.pcapHeader, p[0:pcapHeaderSize-currHdrLen]...)
	}

	// Forward to the capturers
	for _, o := range mo.members {
		o.Write(p)
	}
	return len(p), nil
}

func (mo *multiOutput) Close() {
	for _, o := range mo.members {
		o.Close()
	}
}
