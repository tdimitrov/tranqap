#!/bin/bash

if [ -f /usr/sbin/tcpdump ];
then
    mv /usr/sbin/tcpdump /usr/sbin/__tcpdump
fi