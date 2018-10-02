package shell

// StderrToDevNull is a bash snippet which redirects stderr output to /dev/null
const StderrToDevNull = " 2> /dev/null "

// RunInBackground is a bash snippet which runs the previous command in background
const RunInBackground = " & "

// PidPrefix is the string, which is put in front of the PID, when it is transmitted over stderr
const PidPrefix = "MY_PID_IS:"

// CmdHandlePid returns a Bash one-liner, which does the following:
// 1. Saves the PID of the last command executed in a Bash variable. This is supposed to be the capture command
// 2. Echoes the PID to stderr, so that it can be saved by the Capturer. Stderr is used, because PCAP data is
//		transmitted over stdout
// 3. Waits the PID to finish, so that the session remains active until stop command is sent from rpcap shell
func CmdHandlePid() string {
	return "RPCAP_MY_PID=$! ; echo " + PidPrefix + " $RPCAP_MY_PID >&2 ; wait $RPCAP_MY_PID"
}
