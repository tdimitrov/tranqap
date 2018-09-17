package output

// Outputer is a destination for PCAPs. Concrete implementations
// can be 'save to file', 'view in wireshark', etc.
type Outputer interface {
	Write(p []byte) (n int, err error)
	Close()
}
