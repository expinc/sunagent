#!/usr/bin/env sh
mkdir -p gen
go build -o gen/sunagentd cmd/sunagentd/sunagentd.go
cp config.conf gen/
