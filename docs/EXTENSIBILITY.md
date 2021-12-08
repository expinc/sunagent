# Extensibility

SunAgent has an extensible architecture. You could contribute more features by adding operations, adding platforms to be supported, adding protocals to expose the operations.

## Add More Operations

To add operations, first, you add functions in *ops* package with below signature:
```go
func function_name(ctx context.Context, other_argument type_of_other_argument...) (return_value return_type...)
```

Then you add corresponding APIs of the protocals that expose the operation. For the built-in HTTP protocal, you add a request handler in *http/hanlders* package with below signature:
```go
func handler_name(ctx *gin.Context)
```

Finally, register the handler in function ```http.registerHandlers()```. You must group the URL prefix with the built-in middlewares if you want to enable authentication for the API you add (Refer to the existing code for examples).

For other protocals that you add to SunAgent, just expose the operations in the way of your protocals.

If the operations you are going to add could be implemented by system command, leverage grimoire, which is illustrated in [Architecture Specification](ARCHITECTURE_SPECIFICATION.md).

## Support More Platforms

Most of the operations that SunAgent provides are implemented by standard libraries and *shirou/gopsutil*. Usually you don't need to change the code when you compiling SunAgent for the platforms that SunAgent does not announce to support.

If you find some operations fail on some platforms, just fix it by modifying the operations.

If the operations are implemented by system commands that are configured in grimoire, you could just add more spells for the OS families, or add more grimoires for the OS types.

## Expose by More Protocals

*ops* package just includes a set of functions that could be called by other packages. You could add packages of the protocals what ever you want, like gRPC. You route the requests of the protocals to the operations, then fabricate the response by the operation results. In function ```main.main()```, you start the servers of your protocals. Additionally, you may define more configuration options for your protocals in [etc/config.conf](../etc/config.conf).

## Background Job

If you are going to add operations that could be executed in the background, i.e., supporting asynchronous requests, then you could leverage built-in job framework.

First, you create a struct type that consists of struct *ops.jobBase*.

Then, you implement the functions of interface *ops.Job*, except ```getInfo() *JobInfo``` that has already been implemented by *ops.jobBase*.

* ```execute() error```: You implement how to execute the operation in this function. An error should be returned in case of failure. You may change the progress during execution by modifying *ops.jobBase.info.Progress*, which is the percentage in 0~100 of how much the job is accomplished.
* ```cancel()```: You implement the cancel logics in this function. You may leverage *ops.jobBase.ctx* to implement it. Note that **You must set job status as CANCELED by yourself. Otherwise, the job status will end up with SUCCESSFUL or FAILED**. If the operation could not be canceled, just leave it empty.
* ```dispose()```: You release delegated resources, e.g., close files opened by the job, in this function. If nothing to release, just leave it empty.

Finally, you specify how the job is created in function *ops.createJob*. You may initialize *ops.jobBase.params* to specify the parameters necessary to the job.

You may refer to *ops.dummyJob* to see how a job should be implemented.
