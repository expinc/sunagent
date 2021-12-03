package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"io/ioutil"
	"net/http"
	"strconv"

	"os/exec"

	"github.com/gin-gonic/gin"
)

func ExecScript(ctx *gin.Context) {
	// get parameters
	programParams, ok := ctx.Request.URL.Query()["program"]
	if !ok {
		RespondMissingParams(ctx, []string{"program"})
		return
	}
	program := programParams[0]
	separateOutput := false
	separateOutputParams, ok := ctx.Request.URL.Query()["separateOutput"]
	if ok {
		separateOutput, _ = strconv.ParseBool(separateOutputParams[0])
	}
	waitSeconds := int64(60)
	waitSecondsParams, ok := ctx.Request.URL.Query()["waitSeconds"]
	if ok {
		waitSeconds, _ = strconv.ParseInt(waitSecondsParams[0], 10, 64)
	}
	async := false
	asyncParams, ok := ctx.Request.URL.Query()["async"]
	if ok {
		async, _ = strconv.ParseBool(asyncParams[0])
	}

	// get body
	var result interface{}
	var err error
	script, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
	}

	// create background job if async == true
	if async {
		params := make(map[string]interface{})
		params["program"] = program
		params["script"] = string(script)
		params["waitSeconds"] = waitSeconds
		params["ifSeparateOutput"] = separateOutput
		var jobInfo ops.JobInfo
		jobInfo, err = ops.StartJob(createStandardContext(ctx), ops.JobTypeExecScript, params)
		if nil == err {
			RespondSuccessfulJson(ctx, http.StatusOK, jobInfo)
		} else {
			RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		}
		return
	}

	// execute operation
	if separateOutput {
		result, err = ops.ExecScriptWithSeparateOutput(createStandardContext(ctx), program, string(script), waitSeconds)
	} else {
		result, err = ops.ExecScriptWithCombinedOutput(createStandardContext(ctx), program, string(script), waitSeconds)
	}

	// response
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, result)
	} else {
		// defaultly, should respond InternalServerError and no data
		status := http.StatusInternalServerError
		var data interface{}
		data = nil

		// if the execution timeout, respond RequestTimeout and execution result
		internalErr, ok := err.(common.Error)
		if ok && common.ErrorTimeout == internalErr.Code() {
			status = http.StatusRequestTimeout
			data = result
		}

		// if the execution returns non-zero, respond InternalServerError and execution result
		_, ok = err.(*exec.ExitError)
		if ok {
			data = result
		}

		RespondFailedJson(ctx, status, err, data)
	}
}
