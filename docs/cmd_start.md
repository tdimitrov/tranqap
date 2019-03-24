> start

Starts packet capturing on target(s). Files are saved to the directory specified with **Destination** parameter in the configuration.

The files are named according to the value specified in **File Pattern** parameter.

On each start, PCAP files for each target are rotated. How many files to be kept is specified with **File Rotation Count** parameter.

The file from the current capture is always named **FILE_PATTERN.pcap**. On the next start it is rotated to FILE_PATTERN.1.pcap.

Here is an example :

```json
{
    "Destination": "PCAPs/local_target_1",
    "File Pattern": "trace",
    "File Rotation Count": 5,
}
```

With this configuration the PCAP files will be saved in a location, relative to the current working directory of the binary - **PCAPs/local_target_1**. The files there will have got the following names:

* trace.pcap
* trace.1.pcap
* trace.2.pcap
* trace.3.pcap
* trace.4.pcap
* trace.5.pcap

On the next start, trace.5.pcap will be deleted, trace.4.pcap will be renamed to trace.5.pcap and so on.