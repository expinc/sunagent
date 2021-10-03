package handlers

import (
	"expinc/sunagent/ops"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPackageInfo(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		RespondMissingParams(ctx, []string{"name"})
		return
	}

	pkgInfo, err := ops.GetPackageInfo(createCancellableContext(ctx), name)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pkgInfo)
	} else {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
}
