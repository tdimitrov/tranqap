package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tdimitrov/rpcap/rplog"

	"golang.org/x/crypto/ssh"
)

type config struct {
	Targets []target
}

type target struct {
	Name        *string
	Host        *string
	Port        *int
	User        *string
	Key         *string
	Destination *string
	FilePattern *string `json:"File Pattern"`
	RotationCnt *int    `json:"File Rotation Count"`
	UseSudo     *bool   `json:"Use sudo"`
}

func getConfig(fname string) (config, error) {
	var conf config
	var err error

	confFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		return conf, fmt.Errorf("Error parsing %s: %s", fname, err.Error())
	}

	if len(conf.Targets) == 0 {
		return conf, fmt.Errorf("No targets defined in config")
	}

	return conf, nil
}

func getClientConfig(t *target) (*ssh.ClientConfig, *string, error) {
	var clientConfig ssh.ClientConfig

	clientConfig.Auth = make([]ssh.AuthMethod, 0, 2)

	if t.Name == nil {
		return nil, nil, errors.New("Missing Name in configuration")
	}

	if t.User == nil {
		return nil, nil, fmt.Errorf("Missing user for target <%s> in configuration", *t.Name)
	}

	if t.Key == nil {
		return nil, nil, fmt.Errorf("Missing Key path for target <%s> in configuration", *t.Name)
	}

	if t.Host == nil {
		return nil, nil, fmt.Errorf("Missing Host for target <%s> in configuration", *t.Name)
	}

	if t.Port == nil {
		return nil, nil, fmt.Errorf("Missing port for target <%s> in configuration", *t.Name)
	}

	if t.Destination == nil {
		return nil, nil, fmt.Errorf("Missing destination for target <%s> in configuration", *t.Name)
	}

	if t.FilePattern == nil {
		return nil, nil, fmt.Errorf("Missing File Pattern for target <%s>", *t.Name)
	}

	if t.RotationCnt == nil {
		rplog.Info("File Rotation Count not set for target <%s>. Setting to 10.\n", *t.Name)
		t.RotationCnt = new(int)
		*t.RotationCnt = 10
	}

	if *t.RotationCnt < 0 {
		return nil, nil, fmt.Errorf("Invalid rotation count for target <%s> (%d)", *t.Name, *t.RotationCnt)
	}

	if t.UseSudo == nil {
		t.UseSudo = new(bool)
		*t.UseSudo = false
	}

	dest := fmt.Sprintf("%s:%d", *t.Host, *t.Port)

	clientConfig.User = *t.User
	clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	if t.Key != nil {
		key, err := ioutil.ReadFile(*t.Key)
		if err != nil {
			msg := fmt.Sprintf("unable to read private key: %v", err)
			return nil, nil, errors.New(msg)
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			msg := fmt.Sprintf("unable to parse private key: %v", err)
			return nil, nil, errors.New(msg)
		}

		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}

	return &clientConfig, &dest, nil
}

func generateSampleConfig(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists. Will not overwrite existing config", path)
	}

	name := "Target name. Informational identification only."
	host := "Hostname/IP address of the target."
	port := 22
	login := "SSH login."
	key := "Path to private key, used for authentication."
	dest := "Path to destination dir for the PCAP files."
	pattern := "Filename pattern for each pcap file. Index and file extension will be added to this string."
	rotCnt := 5
	useSudo := true

	t := make([]target, 1, 1)
	t[0] = target{&name, &host, &port, &login, &key, &dest, &pattern, &rotCnt, &useSudo}
	conf := make(map[string][]target)
	conf["targets"] = t

	// And finally create the new file
	confJSON, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return fmt.Errorf("Error serializing sample configuration: %s", err)
	}

	err = ioutil.WriteFile(path, confJSON, 0644)
	if err != nil {
		return fmt.Errorf("Error writing sample configuration to file: %s", err)
	}

	return nil
}
