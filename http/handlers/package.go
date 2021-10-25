package handlers

import (
	"expinc/sunagent/common"
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
		status := http.StatusInternalServerError
		if internalErr, ok := err.(common.Error); ok {
			if common.ErrorNotFound == internalErr.Code() {
				status = http.StatusNotFound
			}
		}
		RespondFailedJson(ctx, status, err, nil)
	}
}

func InstallPackage(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		RespondMissingParams(ctx, []string{"name"})
		return
	}

	pkgInfo, err := ops.InstallPackageByName(createCancellableContext(ctx), name)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pkgInfo)
	} else {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
}
