#!/usr/bin/env sh
mkdir -p gen
go mod tidy
go build -o gen/sunagentd cmd/sunagentd/sunagentd.go

cp -r etc/* gen/
