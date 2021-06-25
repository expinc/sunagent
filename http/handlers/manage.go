package handlers

import (
	"expinc/sunagent/ops"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfo(ctx *gin.Context) {
	info := ops.GetInfo(createCancellableContext(ctx))
	respondSuccessfulJson(ctx, http.StatusOK, info)
}
