# API Reference

All APIs are exposed by url prefix */api/v1*. For example, to get basic information of SunAgent, you should ```GET http://{ip}:{port}/api/v1/info```.

Trace ID could be specified by request header *traceId*. Otherwise, SunAgent will generate a trace ID for each request.

For a successful textual response, the response body example is:

```json
{
    "successful": true,
    "status": 200,
    "traceId": "Trace ID",
    "data": "Payload object of the response"
}
```

For a failed textual response, the response body example is:

```json
{
    "successful": false,
    "status": 500,
    "traceId": "Trace ID",
    "data": "Error message"
}
```

For a binary response, the response body is the binary content.

## Table of Content

* [SunAgent Management](#sunagent-management)
    - [Get Information](#get-information)
    - [Terminate](#terminate)
* [File Management](#file-management)
    - [Get File Meta](#get-file-meta)
    - [Get File Content](#get-file-content)
    - [Create File](#create-file)
    - [Overwrite File](#overwrite-file)
    - [Delete File](#delete-file)
* [Process Management](#process-management)
    - [Get Process Information](#get-process-information)
    - [Kill Process](#kill-process)
* [System Information](#system-information)
    - [Get Node Information](#get-node-information)
    - [Get CPU Information](#get-cpu-information)
    - [Get CPU Statistics](#get-cpu-statistics)
    - [Get Memory Statistics](#get-memory-statistics)
    - [Get Disk Information](#get-disk-information)
    - [Get Network Information](#get-network-information)
* [Script Execution](#script-execution)
    - [Execute script](#execute-script)

## SunAgent Management

### Get Information

Get the basic information of SunAgent, etc., version.

Method: ```GET /info```

Response:

```json
{
    "Version": "x.y.z"
}
```

### Terminate

Stop SunAgent process.

Method: ```POST /terminate```

Parameters:

* waitSec: Waiting seconds for gracefully stopping (default is 3). If it is exceeded, SunAgent will be stopped forcibly.

## File Management

### Get File Meta

Get metadata of files. If the specified file is a directory and parameter *list* is *true*, the result will be metadata of all the files under the directory (not recursive).

Response field *size* is calculated by bytes. The format of response field *mode* is platform dependent.

Method: ```GET /fileMeta```

Parameters:

* path: Absolute path to the file.
* list: Whether list all the files within the *path* if it is a directory. Valid values: *false* (default), *true*.

Response:

```json
[
    {
        "name": "filename.txt",
        "size": 65536,
        "lastModifiedDate": "YYYY-MM-DDThh:mm:ss",
        "owner": "owner",
        "mode": "-rwxrwxrwx"
    },
    {
        "name": "dirname",
        "size": 4096,
        "lastModifiedDate": "YYYY-MM-DDThh:mm:ss",
        "owner": "owner",
        "mode": "-rwxrwxrwx"
    }
]
```

### Get File Content

Get the binary content of the specified file.

It may fail to get the content of some system files, etc., device files in linux.

Method: ```GET /file```

Parameters:

* path: Absolute path to the file.

Response: The binary content of the specified file.

### Create File

Create a file. It will fail if the file already exists or the parent directory does not exist.

Method: ```POST /file```

Parameters:

* path: Absolute path to the file.
* isDir: If the file is a directory.

Body: Binary content of the file if it is a regular file.

Response:

```json
{
    "name": "filename.txt",
    "size": 65536,
    "lastModifiedDate": "YYYY-MM-DDThh:mm:ss",
    "owner": "owner",
    "mode": "-rwxrwxrwx"
}
```

### Overwrite File

Change the content of a file. Create it if not exists.

Method: ```PUT /file```

Parameters:

* path: Absolute path to the file.
* directory: If the file is a directory.

Body: Binary content of the file if it is a regular file.

Response:

```json
{
    "name": "filename.txt",
    "size": 65536,
    "lastModifiedDate": "YYYY-MM-DDThh:mm:ss",
    "owner": "owner",
    "mode": "-rwxrwxrwx"
}
```

### Delete File

Delete a file. If the specified path is a directory, the files under it will also be deleted if parameter *recursive* is true, otherwise it will fail.

Method: ```DELETE /file```

Parameters:

* path: Absolute path to the file.
* recursive: Delete all files under the path if it is a directory.

## Process Management

### Get Process Information

Get basic information of some processes.

If the path parameter is pid, only information of one process with the specified pid will be returned. If the path parameter is process name, information of a number of processes with the specified process name will be returned.

Method: ```GET /processes/{pid or name}```

Response:

```json
[
    {
        "pid": 1024,
        "name": "ProcessName",
        "cmd": "StartCommand",
        "startTime": "YYYY-MM-DDThh:mm:ss",
        "elapsedSeconds": 3600,
        "owner": "owner"
    }
]
```

### Kill Process

Terminate a process, or send a signal to a process (linux).

If the path parameter is pid, only the process with the specified pid will be killed. If the path parameter is process name, all processes with the specified process name will be killed.

Pids of the killed processes will be returned.

Method: ```POST /processes/{pid or name}/kill```

Parameters:

* signal: Signal to send to the process. The default value is SIGTERM. (For linux)

Response:

```json
[
    1024,
    2048
]
```

### Terminate Process

Terminate a process. For windows, it is equivalent to kill process. For linux, it is equivalent to kill process by signal *SIGKILL*.

Method: ```POST /processes/{pid or name}/terminate```

Response:

```json
[
    1024,
    2048
]
```

## System Information

### Get Node Information

Get basic information of the node.

Method: ```GET /sys/info```

Response:

```json
{
    "hostName": "host name",
    "bootTime": "boot time",
    "osType": "operating system type, e.g., linux, windows",
    "osFamily": "operating system family, e.g., debian, rhel",
    "osVersion": "operating system release version",
    "kernelVersion": "operating system kernel version",
    "cpuArch": "CPU architecture, e.g., x86_64, aarch64"
}
```

### Get CPU Information

Get basic informatino of CPUs. The response field *count* is the number of physical threads of all CPUs.

Method: ```GET /sys/cpus/info```

Response:

```json
{
    "count": 12,
    "vendorId": "vendor ID, e.g., GenuineIntel",
    "model": "model",
    "modelName": "model name, e.g., Intel(R) Core(TM) i7-4710MQ CPU @ 2.50GHz",
    "Mhz": "frequency in Mhz"
}
```

### Get CPU Statistics

Get CPU usage and load.

Method: ```GET /sys/cpus/stats```

Parameters:

* perCpu: *true* for getting usage of each CPU (default). *false* for getting total usage of all CPUs.

Response:

```json
{
    "usages" : [
        12.1,
        50.9
    ],
    "load1": 5.1,
    "load5": 2.3,
    "load15": 1.4
}
```

### Get Memory Statistics

Get memory statistics. The values are count by bytes.

The response field *free* is the kernel's notion of free memory, RAM chips whose bits nobody cares about the value of right now. For a human consumable number, *available* is what you really want.

Method: ```GET /sys/mem/stats```

Response:

```json
{
    "total": 8000000000,
    "available": 3000000000,
    "used": 4000000000,
    "free": 4000000000
}
```

### Get Disk Information

Get disk information. The values are count by bytes.

Method: ```GET /sys/disks/stats```

Response:

```json
[
    {
        "device": "dev/sda",
        "mountPoint": "/",
        "fileSystem": "ext4",
        "total": 100000000000,
        "free": 20000000000,
        "used": 800000000000
    },
    {
        "device": "dev/sdb",
        "mountPoint": "/home",
        "fileSystem": "ext4",
        "total": 100000000000,
        "free": 20000000000,
        "used": 800000000000
    }
]
```

### Get Network Information

Method: ```GET /sys/net/info```

Response:

```json
[
    {
        "name": "network adaptor name, e.g., eth0, lo0",
        "mtu": 65535,
        "hardwareAddress": "MAC address",
        "addresses": [
            "127.0.0.1"
        ]
    },
    {
        "name": "network adaptor name, e.g., eth0, lo0",
        "mtu": 1500,
        "hardwareAddress": "MAC address",
        "addresses": [
            "172.10.10.10",
            "192.168.0.3"
        ]
    }
]
```

## Script Execution

### Execute script

Method: ```POST /script/execute```

Parameters:

* program: The program to execute the script, e.g., bash, python.
* separateOutput: *false* to return all output together (default). *true* to return stdout and stderr separately.
* waitSeconds: Seconds to wait for the script execution to complete. The default value is 60.
* killIfOvertime: *false* to do nothing if the script execution is overtime (default). *true* to kill the script execution process if overtime.

Body: script content.

Response:

When *separateOutput=false*, it will be all output content.

When *separateOutput=true*, it will be like below:

```json
{
    "stdout": "stdout content",
    "stderr": "stderr content"
}
```
