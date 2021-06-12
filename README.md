# SunAgent

SunAgent is a lightweight cross-platform remote node operation service.

Clients of SunAgent could get system information, execute scripts, create files, etc., remotely. Clients only need to care about what to do, but care little about how to do, especially how to accomplish the same task on different platforms. For example, installing software "foobar" should do "apt-get install foobar" on debian but "yum install foobar" on redhat. Clients only tell SunAgent to "install a software foobar" regardless of the node should install it by apt-get or by yum. SunAgent encapsulates the actual way to accomplish the installation.

A node infer to some computing resource. It may be a host, a virtual machine, or a container. Different kind of node could do different operations. For example, you cannot call systemctl in docker container without setting docker entrypoint as /usr/sbin/init. SunAgent provides functionalities that work on a host, though some of them may fail in a container.

SunAgent exposes its functionality by common protocals. Currently there is only RESTful APIs.

## Build & Run

TBD

## Configuration

TBD

## API Reference

TBD

## Configurable Shell Command Execution

TBD

## Extensibility

TBD
