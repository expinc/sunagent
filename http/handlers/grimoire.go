package handlers

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	}
	args := strings.Split(string(body), "\n")

	// execute and get output
	output, err := ops.CastGrimoireArcaneContext(createStandardContext(ctx), arcaneName, args...)
	result := ops.CombinedScriptResult{
		Output: string(output),
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
		exitError, ok := err.(*exec.ExitError)
		if ok {
			result.ExitStatus = exitError.ExitCode()
			data = result
		}

		RespondFailedJson(ctx, status, err, data)
	}
}
