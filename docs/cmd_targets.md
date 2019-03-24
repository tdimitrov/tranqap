> targets

This command has got two main purposes:

- To list all targets in the configuration file.
- To check if PCAP files can be collected there.

Here is a sample output from the command. It is similar to the one in [Quick start](quickstart.md) section, but here more details will be provided about each line:

```
rpcap> targets
=== Running checks for target <Local target 1> ===
Check if tcpdump is installed: Yes
Check if sudo is installed: NO
Check if tcpdump can be run with sudo, without password: NO
Check if tcpdump has got cap_net_admin capabilities: NO
Check if tcpdump has got cap_net_raw+eip capabilities: NO
User is member of the binary's group: Yes
```

**Check if tcpdump is installed:** Yes/NO.

rpcap uses tcpdump to collect traffic. This check verifies if tcpdump command is available on the target.

**Check if sudo is installed:** Yes/NO

Usually only privileged users can run tcpdump. One way to achieve this is with sudo. This check verifies if there is sudo installed on target.

**Check if tcpdump can be run with sudo, without password:** Yes/NO

rpcap can't provide a password for sudo. For this reason if rpcap should be started via sudo, it should be configured to execute tcpdump command without asking for a password. This line checks if sudo tcpdump requires a password.

**Check if tcpdump has got cap_net_admin capabilities:** Yes/NO

**Check if tcpdump has got cap_net_raw+eip capabilities:** Yes/NO

**User is member of the binary's group:** Yes/NO

These three lines are connected. The target might be configured to allow an unprivileged user to collect PCAPs with tcpdump. This is possible when the binary has got cap_net_admin and cap_net_raw+eip capabilities enabled. Additionally the user, executing tcpdump, should either be owner of the binary or be member of the owner group. As the owner of tcpdump is usually root, the second condition is checked.
