#!/usr/bin/env bash
CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
gofmt -w src
go install server \
&& go install server \
&& mkdir -p bin/config \
&& cp ./src/config/client.config.json ./src/config/server.config.json bin/config
export GOPATH="$OLDGOPATH"
echo 'finished'