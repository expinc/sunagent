# SunAgent

SunAgent is a lightweight cross-platform remote node operation service.

Clients of SunAgent could get system information, execute scripts, create files, etc., remotely. Clients only need to care about what to do, but care little about how to do, especially how to accomplish the same task on different platforms. For example, installing software "foobar" should do "apt-get install foobar" on debian but "yum install foobar" on redhat. Clients only tell SunAgent to "install a software foobar" regardless of the node should install it by apt-get or by yum. SunAgent encapsulates the actual way to accomplish the installation.

A node infer to some computing resource. It may be a host, a virtual machine, or a container. Different kind of node could do different operations. For example, you cannot call systemctl in docker container without setting docker entrypoint as /usr/sbin/init. SunAgent provides functionalities that work on a host, though some of them may fail in a container.

SunAgent exposes its functionality by common protocals. Currently there is only HTTP APIs.

## Build & Run

Run below command to build the executable. The executable and other files necessary to run it will be output to path *gen*.
```
# Linux:
./make.sh

# Windows:
make.bat
```

Run below command to start the process.
```
# Linux:
./gen/sunagentd --config=gen/config.conf

# Windows:
gen\sunagentd.exe --config=gen\config.conf
```

## Configuration

Configuration is defined by *config.conf*.

*GO* section defines behaviors of the process. They are all go related parameters since SunAgent is implemented by go:

* gomaxprocs: The maximum processors used.

*HTTP* section defines the parameters of the HTTP server who exposes HTTP APIs:

* ip: IP address of the HTTP server. This option is introduced because there may be multiple network adaptors with different IP addresses in the node.
* port: Port number of the HTTP server.
* auth: The type of authorization. Valid values are *none*, *basic*.
* user: User of the basic authorization.
* password: Password of the basic authorization.

*LOG* sectino defines behaviors of logging.

* level: Log level. Valid values are *debug*, *info*, *warn*, *error*, *fatal* (from lower to higher).
* filelimitmb: Log file size limit by MB. When the log file exceeds the limit, it will be rotated.

## API Reference

Please refer to [here](docs/API_REFERENCE.md).

## Test

To run tests, use *test.py*:
* Run ```python3 test.py -h``` to see the usage.
* Run ```python3 test.py -t unit``` to run unit test.
* Run ```python3 test.py -t func``` to run functionality test.
* Run ```python3 test.py``` to run all tests.

Note:
* *test.py* requires python 3.8 and packages illustrated in *requirements.txt*.
* Your python 3 must use "python3" as the python executable name (either the binary file or a symbolic link to it) instead of "python". This restriction is to avoid disrupting RPM, which depends on python 2.

## Extensibility

TBD
