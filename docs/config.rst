
Configuration file
==================

tranqap looks for a file named **config.yaml** in the current working directory. As the name suggests, it's a YAML
formatted file, which contains an array of targets. The general structure is like this:

.. code:: yaml

    targets:
        - name: Target_1
        - name: Target_2
        - name: Target_N

For each target a set of mandatory and optional parameters can be set. 

Mandatory parameters
--------------------

**Name** - Identificator for the target. Should be unique for each target in a configuration file.

**Host** - Hostname/IP address to connect to.

**User** - SSH login.

**Key** - Path to a private key, used for SSH authentication.

**Destination** - Destination directory, where PCAP files should be saved.

**File Pattern** - Base file name for each file. Rotation index and .pcap extension will be added to this value.


Optional parameters
-------------------

**Port** - Port number to connect to. Default value: 22.

**File Rotation count** - How many PCAP files to keep for the target. Default value: 10.

**Use sudo** - true or false. Whether capturer should be invoked with or without sudo. Default value: false.

**Filter port** - Tranqap doesn't include the traffic from the SSH session used to connect to the remote machine. 
The reason is to avoid bloating the PCAP file with irrelevant traffic. However if the target is behind NAT or 
there is a port redirection, the port used for connection might differ from the actual port, on which SSH service 
listens. In that case tranqap will set wrong capture filter for tcpdump and the traffic from the SSH session will 
not be excluded from the capture. This option allows the default filter port to be overridden. If not set, **Port** 
value will be used for the filter. Default value: unset.
