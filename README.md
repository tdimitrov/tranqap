[![CircleCI](https://circleci.com/gh/tdimitrov/rpcap/tree/master.svg?style=svg)](https://circleci.com/gh/tdimitrov/rpcap/tree/master)

# rpcap

Remote network packet capturing tool which automates the generation of PCAP files from one or more remote machines.

rpcap automates things like logging on the remote mahcine, executing a packet capturer (tcpdump), transferring the PCAP file to the local mahcine and executing a graphical tool (wireshark), which dislays the traffic in real time.

The tool runs on Linux and doesn't require the installation of any specific software (besides the packet capturer itself) on the target. In theory there are no restrictions about the target machine as long as it has got SSH server and tcpdump installed.

## The problem
Network traffic capturing is usually a simple and straigh forward process, but it can become a bit tedious if it needs to be performed on a remote machine. And if there are more than one remote machine, a lot of manual work is required.

For single machine I used to run a simple bash script which executes tcpdump over ssh and pipes the PCAP to a local wireshark process. Essentially the script did this:

```bash
ssh "$@" tcpdump -U -s0 -w - "ip and not port 22" | wireshark -k -i - > /dev/null 2>&1 &
```

For quite a lot of time this worked just fine for my needs, but the script has got its drawbacks:
  * It's up to the user of the script to make sure that the SSH user has got sufficient rights to run tcpdump on the target machine.
  * If wireshark is closed during capture, the only way to resume capture is to restart the whole trace.
  * Running the script for multiple machines is a bit tedious. The captures can be started with a script, but one should be careful killing them.
  * The script needs to be modified in order to work on different targets. E.g. on one machine tcpdump can be run as regular user, on another it requires sudo, etc. So it is hard to run multiple captures on machines, requiring different authentication/permissions.
  
## What is rpcap and how it solves the problem
rpcap is the bash script above on a very strong steroids. It's main purpose is to implement the functionality of the script in more intelligent and robust way. But in the same time rpcap should be as lean and simple as possible. 

rpcap is driven by three main principles:
1. Minimal dependencies and requirements for the target machine(s).
2. Do not take over the role of another tool(s), just glue them toghether.
3. Be simple, minimalistic and user friendly.

What rpcap does:
1. Reads a list of remote targets and configuration parameters from a JSON file. Each config file represents an environment/project/task.
2. A connection to each target is established over SSH. 
3. A capturer (e.g. tcpdump) is run on each target.
4. The resulting PCAP file is saved locally on disk and optionally it can be watched in real time with GUI packet analyser (e.g. wireshark).
5. Handles starting and stopping of the capturer process and also the SSH connections.
6. Checks if the user, set in the config file, has sufficient permissions to run the capturer.

# Quick start guide
For the moment no prebuilt binaries are provided so the only installation option is from source.

## Installation from source
rpcap is written in Go, so the Go distribution should be installed. Instructions for the installation can be found [here](https://golang.org/doc/install). 

After that:
```bash
go get  github.com/tdimitrov/rpcap
go install  github.com/tdimitrov/rpcap
```
At this point rpcap should be installer in `$GOPATH/bin`. This path should be added to system path.

## Configuration and usage
rpcap looks for a file named config.json in the current working dir. A sample file can be found [here](samples/config.json).

Most of the fields are self-explanatory, but anyway:
* **targets** is an array of target machines. For each target:
  * **Name** is a string, which distinguishes the target. No need to be unique, it is used in the error messages and log files to identify the target.
  * **Host**, **Port**, **User**, **Key** - used for the SSH connection. At the moment only Public key authentication is supported. For now Password authentication won't be supported for security reasons.
  * **Destination** is a directory name, where PCAP files for the target will be saved.
  * **File pattern** is a filename pattern for each PCAP file. A rotation counter and file extension will be appended to the pattern. E.g. PATTERN.1.pcap
  * **File rotation count** determines how many files to keep on rotation.
  * **Use sudo** is a bool variable specifying if tcpdump should be run with sudo or not.

The binary should be executed in the directory, where `config.json` is located. The application uses a very basic shell as a UI. The following commands are supported:
* **targets** - Lists all targets and prints details about binary's permissions/capabilities. This is an indication if tcpdump can be run without root.
```
rpcap> targets
=== Running checks for target <Local target 1> ===
Check if tcpdump is installed: Yes
Check if tcpdump has got cap_net_admin capabilities: NO
Check if tcpdump has got cap_net_raw+eip capabilities: NO
User is member of the binary's group: Yes

=== Running checks for target <Local target 2> ===
Check if tcpdump is installed: Yes
Check if tcpdump has got cap_net_admin capabilities: NO
Check if tcpdump has got cap_net_raw+eip capabilities: NO
User is member of the binary's group: Yes

rpcap>  
```
* **start** - Executes the capturer on all remote targets. A PCAP file is saved to the destination directory and the old files are rotated.
* **stop** - Stops the capturer on all targets.
* **wireshark** - Starts a wireshark instance for each target. All traffic, which is saved on disk is also piped to wireshark, so it can be viewed in real  time.

rpcap accepts two command line parameters:
* -l - path to log file
* -c - path to config file
