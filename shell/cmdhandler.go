package shell

// CmdHandler processes the output of a specific shell command
type CmdHandler interface {
	Write(p []byte) (n int, err error)
}
