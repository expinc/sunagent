package handlers

import (
	"expinc/sunagent/ops"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNodeInfo(ctx *gin.Context) {
	info := ops.GetNodeInfo(createStandardContext(ctx))
	RespondSuccessfulJson(ctx, http.StatusOK, info)
}

func GetCpuInfo(ctx *gin.Context) {
	info := ops.GetCpuInfo(createStandardContext(ctx))
	RespondSuccessfulJson(ctx, http.StatusOK, info)
}

func GetCpuStat(ctx *gin.Context) {
	perCpuStr, ok := ctx.Request.URL.Query()["perCpu"]
	perCpu := false
	if ok {
		perCpu, _ = strconv.ParseBool(perCpuStr[0])
	}

	stat, err := ops.GetCpuStat(createStandardContext(ctx), perCpu)
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, stat)
	}
}

func GetMemStat(ctx *gin.Context) {
	stat, err := ops.GetMemStat(createStandardContext(ctx))
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, stat)
	}
}

func GetDiskInfo(ctx *gin.Context) {
	info, err := ops.GetDiskInfo(createStandardContext(ctx))
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, info)
	}
}

func GetNetInfo(ctx *gin.Context) {
	info, err := ops.GetNetInfo(createStandardContext(ctx))
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, info)
	}
}
