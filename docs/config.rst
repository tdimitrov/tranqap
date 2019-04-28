
Configuration file
==================

rpcap looks for a file named **json.config** in the current working directory. As the name suggests, it's a JSON
formatted file, which contains an array of targets. The general structure is like this:

.. code:: json

    {
        "targets" : [
            {
            }
        ]
    }

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

**Filter port** - If the target is behind NAT or there is a port redirection, the port used for connection might 
differ from the actual port, on which SSH service listens. In that case rpcap will set wrong capture filter for 
tcpdump and the traffic for the SSH session will not be excluded in the capture. This option allows the default 
filter port to be overridden. If not set, **Port** value will be used for the filter. Default value: unset.