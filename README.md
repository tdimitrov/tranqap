# tranqap

Remote network packet capturing tool which automates the generation of PCAP files from one or more remote machines.

tranqap automates things like logging on the remote machine, executing a packet capturer (tcpdump), transferring the PCAP file to the local mahcine and executing a graphical tool (wireshark), which displays the traffic in real time.

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

## What is tranqap and how it solves the problem
tranqap is the bash script above on a very strong steroids. It's main purpose is to implement the functionality of the script in more intelligent and robust way. But in the same time tranqap should be as lean and simple as possible.

tranqap is driven by three main principles:
1. Minimal dependencies and requirements for the target machine(s).
2. Do not take over the role of another tool(s), just glue them toghether.
3. Be simple, minimalistic and user friendly.

What tranqap does:
1. Reads a list of remote targets and configuration parameters from a YAML file. Each config file represents an environment/project/task.
2. A connection to each target is established over SSH.
3. A capturer (e.g. tcpdump) is run on each target.
4. The resulting PCAP file is saved locally on disk and optionally it can be watched in real time with GUI packet analyser (e.g. wireshark).
5. Handles starting and stopping of the capturer process and also the SSH connections.
6. Checks if the user, set in the config file, has sufficient permissions to run the capturer.

# Learn more
[Documentation](https://tranqap.readthedocs.io)

[Downloads](https://github.com/tdimitrov/tranqap/releases)

[Short demo on YouTube](https://youtu.be/yNIUNll4SaE)
