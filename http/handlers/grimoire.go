package handlers

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetGrimoire(ctx *gin.Context) {
	osType, ok := ctx.Params.Get("osType")
	if !ok {
		RespondMissingParams(ctx, []string{"osType"})
		return
	}

	output, err := ops.GetGrimoireAsYaml(createStandardContext(ctx), osType)
	if nil == err {
		RespondBinary(ctx, http.StatusOK, output)
	} else {
		RespondFailedJson(ctx, http.StatusNotFound, err, nil)
	}
}

func CastArcane(ctx *gin.Context) {
	// get parameters
	osType, ok := ctx.Params.Get("osType")
	if !ok {
		RespondMissingParams(ctx, []string{"osType"})
		return
	}
	if "default" != osType && ops.GetNodeInfo(context.Background()).OsType != osType {
		err := common.NewError(common.ErrorInvalidParameter, "Cannot cast arcane in grimoire of other OS types")
		RespondFailedJson(ctx, http.StatusNotAcceptable, err, nil)
		return
	}
	arcaneName, ok := ctx.Params.Get("arcaneName")
	if !ok {
		RespondMissingParams(ctx, []string{"arcaneName"})
		return
	}

	// get arcane args
	var err error
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		return
	}
	args := strings.Split(string(body), "\n")

	// create background job if async == true
	async := false
	asyncParams, ok := ctx.Request.URL.Query()["async"]
	if ok {
		async, _ = strconv.ParseBool(asyncParams[0])
	}
	if async {
		params := make(map[string]interface{})
		params["arcaneName"] = arcaneName
		params["args"] = args
		jobInfo, err := ops.StartJob(createStandardContext(ctx), ops.JobTypeCastArcane, params)
		if nil == err {
			RespondSuccessfulJson(ctx, http.StatusAccepted, jobInfo)
		} else {
			RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		}
		return
	}

	// execute and get output
	output, err := ops.CastGrimoireArcaneContext(createStandardContext(ctx), arcaneName, args...)
	result := ops.CombinedScriptResult{
		Output: string(output),
	}

	// response
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, result)
	} else {
		result.Error = err.Error()

		// defaultly, should respond InternalServerError and no data
		status := http.StatusInternalServerError
		var data interface{}
		data = nil

		// handle errors of timeout and arcane not found
		internalErr, ok := err.(common.Error)
		if ok {
			if common.ErrorTimeout == internalErr.Code() {
				status = http.StatusRequestTimeout
			} else if common.ErrorNotFound == internalErr.Code() {
				status = http.StatusNotFound
			}
			data = result
		}

		// handle command execution failure
		exitError, ok := err.(*exec.ExitError)
		if ok {
			result.ExitStatus = exitError.ExitCode()
			data = result
		}

		RespondFailedJson(ctx, status, err, data)
	}
}

func SetArcane(ctx *gin.Context) {
	// Get parameters and body
	osType, ok := ctx.Params.Get("osType")
	if !ok {
		RespondMissingParams(ctx, []string{"osType"})
		return
	}
	arcaneName, ok := ctx.Params.Get("arcaneName")
	if !ok {
		RespondMissingParams(ctx, []string{"arcaneName"})
		return
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		return
	}

	// Execute operation and respond
	err = ops.SetGrimoireArcane(createStandardContext(ctx), osType, arcaneName, body)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, nil)
	} else {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
}

func RemoveArcane(ctx *gin.Context) {
	// Get parameters
	osType, ok := ctx.Params.Get("osType")
	if !ok {
		RespondMissingParams(ctx, []string{"osType"})
		return
	}
	arcaneName, ok := ctx.Params.Get("arcaneName")
	if !ok {
		RespondMissingParams(ctx, []string{"arcaneName"})
		return
	}

	// Execute operation and respond
	err := ops.RemoveGrimoireArcane(createStandardContext(ctx), osType, arcaneName)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, nil)
	} else {
		status := http.StatusInternalServerError
		internalErr, ok := err.(common.Error)
		if ok && common.ErrorNotFound == internalErr.Code() {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
	}
}
