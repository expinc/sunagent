# API Reference

All APIs are exposed by url prefix */api/v1*. For example, to get basic information of SunAgent, you should ```GET http://{ip}:{port}/api/v1/info```.

Trace ID could be specified by request header *traceId*. Otherwise, SunAgent will generate a trace ID for each request.

For a textual response, the response body example is:

```json
{
    "successful": true,
    "status": 200,
    "traceId": "Trace ID",
    "data": "Payload object of the response",
    "error": "error message if error occurs, otherwise empty"
}
```

For a binary response, the response body is the binary content.

Some of the APIs could be called asynchronously. The request will create a corresponding background job, and response the ID and status of the job. You may check the status or cancel the job. For a asynchronous call, your request will get a response payload the same as the API of [getting job status](#get-job-status).

The status of background job is not persistent. If SunAgent is restarted, the status will be lost. However, the ID of the job is unique even if SunAgent is restarted.

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
    - [Execute Script](#execute-script)
* [Package Management](#package-management)
    - [Get Package Information](#get-package-information)
    - [Install Package](#install-package)
    - [Upgrade Package](#upgrade-package)
    - [Uninstall Package](#uninstall-package)
* [Background Job](#background-job)
    - [Get Job Information](#get-job-information)
    - [List Jobs](#list-jobs)
    - [Cancel Job](#cancel-job)
* [Grimoire Management](#grimoire-management)
    - [Get Grimoire](#get-grimoire)
    - [Cast Arcane](#cast-arcane)
    - [Set Arcane](#set-arcane)
    - [Remove Arcane](#remove-arcane)

## SunAgent Management

### Get Information

Get the basic information of SunAgent, etc., version.

Method: ```GET /info```

Status:
- 200 OK: Request succeeded.

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

Status:
- 204 No Content: Request succeeded.

## File Management

### Get File Meta

Get metadata of files. If the specified file is a directory and parameter *list* is *true*, the result will be metadata of all the files under the directory (not recursive).

Response field *size* is calculated by bytes. The format of response field *mode* is platform dependent.

Method: ```GET /file/meta```

Parameters:

* path: Absolute path to the file.
* list: Whether list all the files within the *path* if it is a directory. Valid values: *false* (default), *true*.

Status:
- 200 OK: Request succeeded.
- 404 Not Found: File not found.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Request succeeded.
- 404 Not Found: File not found.
- 500 Internal Server Error: Request failed.

Response: The binary content of the specified file.

### Create File

Create a file. It will fail if the file already exists or the parent directory does not exist.

Method: ```POST /file```

Parameters:

* path: Absolute path to the file.
* isDir: If the file is a directory.

Body: Binary content of the file if it is a regular file.

Status:
- 201 Created: File created.
- 400 Bad Request: Corrupted file content.
- 500 Internal Server Error: File creation failed.

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

Status:
- 200 OK: File overwritten.
- 400 Bad Request: Corrupted file content.
- 500 Internal Server Error: Request failed.

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

### Append File

Append some content to an existing file.

Method: ```POST /file/append```

Parameters:

* path: Absolute path to the file.

Body: Binary content of the file if it is a regular file.

Status:
- 200 OK: File content appended.
- 400 Bad Request: Corrupted file content.
- 404 Not Found: File not found.
- 500 Internal Server Error: Request failed.

### Delete File

Delete a file. If the specified path is a directory, the files under it will also be deleted if parameter *recursive* is true, otherwise it will fail.

Method: ```DELETE /file```

Parameters:

* path: Absolute path to the file.
* recursive: Delete all files under the path if it is a directory.

Status:
- 200 OK: File deleted.
- 404 Not Found: File not found.
- 500 Internal Server Error: Request failed.

## Process Management

### Get Process Information

Get basic information of some processes.

If the path parameter is pid, only information of one process with the specified pid will be returned. If the path parameter is process name, information of a number of processes with the specified process name will be returned.

Method: ```GET /processes/{pid or name}```

Status:
- 200 OK: Request succeeded.
- 404 Not Found: Process not found.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Signal sent to process.
- 404 Not Found: Process not found.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Process is terminated.
- 404 Not Found: Process not found.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Request succeeded.

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

Status:
- 200 OK: Request succeeded.

Response:

```json
{
    "modelName": "model name, e.g., Intel(R) Core(TM) i7-4710MQ CPU @ 2.50GHz",
    "vendorId": "vendor ID, e.g., GenuineIntel",
    "Mhz": "frequency in Mhz",
    "count": 12
}
```

### Get CPU Statistics

Get CPU usage and load.

Method: ```GET /sys/cpus/stats```

Parameters:

* perCpu: *true* for getting usage of each CPU (default). *false* for getting total usage of all CPUs.

Status:
- 200 OK: Request succeeded.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Request succeeded.
- 500 Internal Server Error: Request failed.

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

Status:
- 200 OK: Request succeeded.
- 500 Internal Server Error: Request failed.

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

Get network interface information.

Method: ```GET /sys/net/info```

Status:
- 200 OK: Request succeeded.
- 500 Internal Server Error: Request failed.

Response:

```json
[
    {
        "name": "network adaptor name, e.g., eth0, lo0",
        "maxTransmissionUnit": 65535,
        "hardwareAddress": "MAC address",
        "ipAddresses": [
            "127.0.0.1/8"
        ]
    },
    {
        "name": "network adaptor name, e.g., eth0, lo0",
        "maxTransmissionUnit": 1500,
        "hardwareAddress": "MAC address",
        "ipAddresses": [
            "172.10.10.10/24",
            "fe80::d987:8873:f85f:e621/64"
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
* waitSeconds: Seconds to wait for the script execution to complete. If timeout, the script process will be killed. The default value is 60. Set it as 0 if wait until the script completes execution.
* async: *true* to make this request be handled by a background job. *false* to make this request be handled synchronously as usual. Default is *false*.

Body: script content.

Status:
- 200 OK: Script execution succeeded.
- 202 Accepted: Background job created.
- 400 Bad Request: Invalid parameters or request body.
- 408 Request Timeout: Script execution timeout.
- 500 Internal Server Error: Script execution failed.

Response:

When *separateOutput=false*, it will be like below:

```json
{
    "output": "combined stdout and stderr",
    "exitStatus": 0,
    "error": "error message if error occurs, otherwise empty"
}
```

When *separateOutput=true*, it will be like below:

```json
{
    "stdout": "stdout content",
    "stderr": "stderr content",
    "exitStatus": 0,
    "error": "error message if error occurs, otherwise empty"
}
```

**Note: *exitStatus* will be 0 if timeout.**

## Package Management

**Note: Windows is not supported**

### Get Package Information

Method: ```GET /package/{name}```

Status:
- 200 OK: Request succeeded.
- 404 Not Found: Package not installed.
- 500 Internal Server Error: Request failed.

Response:

```json
{
    "name": "package name",
    "version": "package version",
    "architecture": "package architecture",
    "summary": "package summary"
}
```

### Install Package

Install a package. It will fail if the package has already been installed.

Method: ```POST /package/{name}```

Parameters:

* path: Package file path. This parameter works only if the url parameter *name* is not specified.
* async: *true* to make this request be handled by a background job. *false* to make this request be handled synchronously as usual. Default is *false*.

Status:
- 201 Created: Package installed.
- 202 Accepted: Background job created.
- 400 Bad Request: Invalid parameter.
- 404 Not Found: Archive not found.
- 500 Internal Server Error: Request failed.

Response:

```json
{
    "name": "package name",
    "version": "package version",
    "architecture": "package architecture",
    "summary": "package summary"
}
```

### Upgrade Package

Upgrade a Package. The package will be installed if it is not installed. The package will be upgraded if the installed version is older than the target package. Nothing will do if the package has already been installed.

Method: ```PUT /package/{name}```

Parameters:

* path: Package file path. This parameter works only if the url parameter *name* is not specified.

Status:
- 200 OK: Package upgraded.
- 202 Accepted: Background job created.
- 400 Bad Request: Invalid parameter.
- 404 Not Found: Archive not found.
- 500 Internal Server Error: Request failed.

Response:

```json
{
    "name": "package name",
    "version": "package version",
    "architecture": "package architecture",
    "summary": "package summary"
}
```

### Uninstall Package

Method: ```DELETE /package/{name}```

Status:
- 200 OK: Package uninstalled.
- 404 Not Found: Package not installed.
- 500 Internal Server Error: Request failed.

## Background Job

### Get Job Information

Get the information of a background job.

Method: ```GET /jobs/{ID}```

Status:
- 200 OK: Request succeeded.
- 404 Not Found: Job not found.

Response:

```json
{
    "type": "Type of the job",
    "id": "ID of the job",
    "status": "Status of the job. It could be EXECUTING, SUCCESSFUL, FAILED or CANCELED",
    "beginTime": "Local wall clock time when this job begins executing",
    "endTime": "Local wall clock time when this job ends executing",
    "progress": 100, // the percentage of job progress
    "result": "Result of the job"
}
```

### List Jobs

List all jobs. Some early ended jobs may not listed according to the configuration *CORE.jobCleanThreshold*.

Method: ```GET /jobs```

Status:
- 200 OK: Request succeeded.

Response:
```json
[
    {
        "type": "Type of the job",
        "id": "ID of the job",
        "status": "Status of the job. It could be EXECUTING, SUCCESSFUL, FAILED or CANCELED",
        "beginTime": "Local wall clock time when this job begins executing",
        "endTime": "Local wall clock time when this job ends executing",
        "progress": 100, // the percentage of job progress
        "result": "Result of the job"
    },
    {
        "type": "Type of the job",
        "id": "ID of the job",
        "status": "Status of the job. It could be EXECUTING, SUCCESSFUL, FAILED or CANCELED",
        "beginTime": "Local wall clock time when this job begins executing",
        "endTime": "Local wall clock time when this job ends executing",
        "progress": 100, // the percentage of job progress
        "result": "Result of the job"
    }
]
```

### Cancel Job

Cancel an executing job.

Method: ```POST /jobs/{ID}/cancel```

Status:
- 200 OK: Job cancelled.
- 404 Not Found: Job not found.
- 500 Internal Server Error: Request failed.

Response:

```json
{
    "type": "Type of the job",
    "id": "ID of the job",
    "status": "Status of the job. It could be EXECUTING, SUCCESSFUL, FAILED or CANCELED",
    "beginTime": "Local wall clock time when this job begins executing",
    "endTime": "Local wall clock time when this job ends executing",
    "progress": 100, // the percentage of job progress
    "result": "Result of the job"
}
```

## Grimoire Management

### Get Grimoire

Get the content of the grimoire of specified OS.

You could specify path parameter *osType* as "default" to specify it as the OS that SunAgent is currently running on.

Method: ```GET /grimoires/{osType}```

Status:
- 200 OK: Request succeeded.
- 404 Not Found: Grimoire not found.

Response: The content of the grimoire.

### Cast Arcane

Execute the command specified by an arcane.

You must specify path parameter *osType* as "default" or as the OS that SunAgent is currently running on.

Method: ```POST /grimoires/{osType}/arcanes/{arcaneName}/cast```

Parameters:

* async: *true* to make this request be handled by a background job. *false* to make this request be handled synchronously as usual. Default is *false*.

Body:

```
parameter1
parameter2
parameter3
```

Status:
- 200 OK: Arcane casting succeeded.
- 202 Accepted: Background job created.
- 400 Bad Request: Invalid parameter or corrupted request body.
- 404 Not Found: Arcane not found.
- 406 Not Acceptable: Cannot cast arcane this way.
- 408 Request Timeout: Arcane casting timeout.
- 500 Internal Server Error: Arcane casting failed.

Response:

When *separateOutput=false*, it will be like below:

```json
{
    "output": "combined stdout and stderr",
    "exitStatus": 0,
    "error": "error message if error occurs, otherwise empty"
}
```

When *separateOutput=true*, it will be like below:

```json
{
    "stdout": "stdout content",
    "stderr": "stderr content",
    "exitStatus": 0,
    "error": "error message if error occurs, otherwise empty"
}
```

### Set Arcane

Set an arcane in the grimoire. It will override the existing one if there is already an arcane has the same name as the one that is being set.

You could specify path parameter *osType* as "default" to specify it as the OS that SunAgent is currently running on.

Method: ```PUT /grimoires/{osType}/arcanes/{arcaneName}```

Body:

```yaml
timeout: 60 # timeout in seconds
spells:
  osFamily1:
    args: command arguments for osFamily1   # use {} as parameter place holder, use {{}} as literal {}
  osFamily2:
    args: command arguments for osFamily2   # use {} as parameter place holder, use {{}} as literal {}
  osFamily3:
    args: command arguments for osFamily3   # use {} as parameter place holder, use {{}} as literal {}
```

Status:
- 200 OK: Arcane was set.
- 400 Bad Request: Corrupted request body.
- 500 Internal Server Error: Request failed.

### Remove Arcane

Remove an arcane.

You could specify path parameter *osType* as "default" to specify it as the OS that SunAgent is currently running on.

Method: ```DELETE /grimoires/{osType}/arcanes/{arcaneName}```

Status:
- 200 OK: Arcane removed.
- 404 Not Found: Arcane not found.
- 500 Internal Server Error: Request failed.
