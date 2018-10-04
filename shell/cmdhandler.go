package shell

// StderrToDevNull is a bash snippet which redirects stderr output to /dev/null
const StderrToDevNull = " 2> /dev/null "

// RunInBackground is a bash snippet which runs the previous command in background
const RunInBackground = " & "

// CmdHandler processes the output of a specific shell command
type CmdHandler interface {
	Write(p []byte) (n int, err error)
}
