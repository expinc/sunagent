#!/usr/bin/env sh
docker build -f test/docker/Dockerfile.centos7.amd64 -t sunagent-centos-amd64:03 .
docker build -f test/docker/Dockerfile.debian9.amd64 -t sunagent-debian-amd64:03 .
docker build -f test/docker/Dockerfile.opensuse15.amd64 -t sunagent-suse-amd64:03 .
