package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProcInfo(ctx *gin.Context) {
	pidOrNameStr, ok := ctx.Params.Get("pidOrName")
	if !ok {
		RespondMissingParams(ctx, []string{"pidOrName"})
		return
	}

	pid, err := strconv.ParseInt(pidOrNameStr, 10, 32)
	var infos []ops.ProcInfo
	if nil == err {
		var info ops.ProcInfo
		info, err = ops.GetProcInfoByPid(createCancellableContext(ctx), int32(pid))
		if nil == err {
			infos = append(infos, info)
		}
	} else {
		infos, err = ops.GetProcInfosByName(createCancellableContext(ctx), pidOrNameStr)
	}

	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, infos)
	} else {
		status := http.StatusInternalServerError
		if internalErr, ok := err.(common.Error); ok {
			if common.ErrorNotFound == internalErr.Code() {
				status = http.StatusNotFound
			}
		}
		RespondFailedJson(ctx, status, err)
	}
}
