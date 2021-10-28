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
	byFile := false
	nameOrPath, ok := ctx.Params.Get("name")
	if !ok {
		path, ok := ctx.Request.URL.Query()["path"]
		if ok {
			nameOrPath = path[0]
			byFile = true
		} else {
			RespondMissingParams(ctx, []string{"path"})
			return
		}
	}

	pkgInfo, err := ops.InstallPackage(createCancellableContext(ctx), nameOrPath, byFile)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pkgInfo)
	} else {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
}
