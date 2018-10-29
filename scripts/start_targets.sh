#!/bin/sh

docker run -d -p 127.0.0.1:2222:22 --name target1 --rm tdimitrov/rpcap-target
docker run -d -p 127.0.0.1:2223:22 --name target2 --rm tdimitrov/rpcap-target