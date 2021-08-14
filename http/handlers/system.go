package handlers

import (
	"expinc/sunagent/ops"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNodeInfo(ctx *gin.Context) {
	info := ops.GetNodeInfo(createCancellableContext(ctx))
	RespondSuccessfulJson(ctx, http.StatusOK, info)
}

func GetCpuInfo(ctx *gin.Context) {
	info := ops.GetCpuInfo(createCancellableContext(ctx))
	RespondSuccessfulJson(ctx, http.StatusOK, info)
}

func GetCpuStat(ctx *gin.Context) {
	perCpuStr, ok := ctx.Request.URL.Query()["perCpu"]
	perCpu := false
	if ok {
		perCpu, _ = strconv.ParseBool(perCpuStr[0])
	}

	stat, err := ops.GetCpuStat(createCancellableContext(ctx), perCpu)
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, stat)
	}
}
