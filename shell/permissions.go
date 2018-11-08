package shell

import (
	"fmt"

	"github.com/tdimitrov/rpcap/rplog"

	"golang.org/x/crypto/ssh"
)

type stdOutWriter struct {
}

func (b stdOutWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func cmdPermissions() string {
	cmd :=
		`
	rpcap_permissions() {
		printf "Check if tcpdump is installed: "
		BIN=$(command -v tcpdump)
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if tcpdump has got cap_net_admin capabilities: "
		getcap $BIN | grep cap_net_admin > /dev/null
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		printf "Check if tcpdump has got cap_net_raw+eip capabilities: "
		getcap $BIN | grep 'cap_net_raw+eip' > /dev/null
		if [ $? -ne 0 ]
		then
			echo "NO"
		else
			echo "Yes"
		fi

		BIN_USER=$(stat -c '%U' $BIN)
		BIN_GROUP=$(stat -c '%G' $BIN)

		if [[ "$(groups)" =~ "$BIN_GROUP" ]]
		then
			echo "User is member of the binary's group: Yes"
		elif [ "${USER}" == "${BIN_USER}" ]
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

// CheckPermissions executes a bash function, which checks if tcpdump can be run on a target machine
func CheckPermissions(c *ssh.ClientConfig, dest string) bool {
	client, err := ssh.Dial("tcp", dest, c)
	if err != nil {
		rplog.Error("Error connecting: %s\n", err)
		return false
	}

	sess, err := client.NewSession()
	if err != nil {
		rplog.Error("Error creating session!")
		return false
	}

	defer sess.Close()

	sess.Stdout = stdOutWriter{}
	sess.Stderr = stdOutWriter{}

	err = sess.Start(cmdPermissions())
	if err != nil {
		rplog.Error("Error running permissions command")
		return false
	}

	sess.Wait()

	return true
}
