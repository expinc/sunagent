package handlers

import (
	"expinc/sunagent/ops"
	"net/http"

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
