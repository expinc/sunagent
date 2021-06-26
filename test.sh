#!/usr/bin/env sh
mkdir -p gen
go test -v -coverpkg=./... ./... -coverprofile gen/coverprofile > gen/test.log
go tool cover -html ./gen/coverprofile -o gen/test.html
