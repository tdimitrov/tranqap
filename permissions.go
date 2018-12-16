package main

import (
	"fmt"
	"strings"
)

type stdOutWriter struct {
	output *strings.Builder
}

func (b stdOutWriter) Write(p []byte) (n int, err error) {
	b.output.Write(p)
	return len(p), nil
}

func cmdPermissions() string {
	cmd :=
		`
	rpcap_permissions() {
		printf "Check if tcpdump is installed: "
		TCPDUMP_BIN=$(command -v tcpdump)
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if sudo is installed: "
		SUDO_BIN=$(command -v sudo)
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if tcpdump can be run with sudo, without password: "
		sudo -n tcpdump --version > /dev/null 2>&1
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if tcpdump has got cap_net_admin capabilities: "
		getcap $TCPDUMP_BIN | grep cap_net_admin > /dev/null
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if tcpdump has got cap_net_raw+eip capabilities: "
		getcap $TCPDUMP_BIN | grep 'cap_net_raw+eip' > /dev/null
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		TCPDUMP_USER=$(stat -c '%U' $TCPDUMP_BIN)
		TCPDUMP_GROUP=$(stat -c '%G' $TCPDUMP_BIN)

		if [[ "$(groups)" =~ "$TCPDUMP_GROUP" ]]
		then
			echo "User is member of the binary's group: Yes"
		elif [ "${USER}" == "${TCPDUMP_USER}" ]
		then
			echo "User is owner of the binary: Yes"
		else
			echo "User is owner of the binary OR member of the group of the binary: NO"
		fi
	}
	rpcap_permissions
	`
	return cmd
}

// checkPermissions executes a bash function, which checks if tcpdump can be run on a target machine
func checkPermissions(trans *SSHClient) (string, error) {
	if err := trans.Connect(); err != nil {
		return "", fmt.Errorf("Error connecting: %s", err)
	}

	out := stdOutWriter{&strings.Builder{}}

	if err := trans.Run(cmdPermissions(), out, out); err != nil {
		return "", fmt.Errorf("Error running permissions command: %s", err)
	}

	return out.output.String(), nil
}
