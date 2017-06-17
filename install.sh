#!/usr/bin/env bash
CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
gofmt -w src
go install server \
&& go install client \
&& mkdir -p ./config \
&& cp ./src/config/client.config.json.example ./config/client.config.json \
&& cp ./src/config/server.config.json.example ./config/server.config.json
export GOPATH="$OLDGOPATH"
echo 'finished'