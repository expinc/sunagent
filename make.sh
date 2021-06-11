#!/usr/bin/env sh
go build -o gen/sunagent cmd/sunagent.go
cp config.conf gen/
