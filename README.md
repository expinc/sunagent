# SunAgent

SunAgent is a lightweight cross-platform remote node operation service.

Clients of SunAgent could get system information, execute scripts, create files, etc., remotely. Clients only need to care about what to do, but care little about how to do, especially how to accomplish the same task on different platforms. For example, installing software "foobar" should do "apt-get install foobar" on debian but "yum install foobar" on redhat. Clients only tell SunAgent to "install a software foobar" regardless of the node should install it by apt-get or by yum. SunAgent encapsulates the actual way to accomplish the installation.

A node infer to some computing resource. It may be a host, a virtual machine, or a container. Different kind of node could do different operations. For example, you cannot call systemctl in docker container without setting docker entrypoint as /usr/sbin/init. SunAgent provides functionalities that work on a host, though some of them may fail in a container.

SunAgent exposes its functionality by common protocals. Currently there is only RESTful APIs.

## Build & Run

Run ```make.sh``` to build the executable. The executable and other files necessary to run it will be output to path *gen*.

Run ```sunagentd``` to start the process.

## Configuration

Configuration is defined by *config.conf*.

*GO* section defines behaviors of the process. They are all go related parameters since SunAgent is implemented by go:

* gomaxprocs: The maximum processors used.

*HTTP* section defines the parameters of the HTTP server who exposes RESTful APIs:

* ip: IP address of the HTTP server. This option is introduced because there may be multiple network adaptors with different IP addresses in the node.
* port: Port number of the HTTP server.
* auth: The type of authorization. Valid values are *none*, *basic*.
* user: User of the basic authorization.
* password: Password of the basic authorization.

*LOG* sectino defines behaviors of logging.

* level: Log level. Valid values are *debug*, *info*, *warn*, *error*, *fatal* (from lower to higher).
* filesizelimitmb: Log file size limit by MB. When the log file exceeds the limit, it will be rotated.

## API Reference

TBD

## Configurable Shell Command Execution

TBD

## Extensibility

TBD
