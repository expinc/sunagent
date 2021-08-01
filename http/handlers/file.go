package handlers

import (
	"expinc/sunagent/ops"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetFileMeta(ctx *gin.Context) {
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	_, listIfDir := ctx.Request.URL.Query()["list"]

	metas, err := ops.GetFileMetas(createCancellableContext(ctx), path[0], listIfDir)
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, metas)
	}
}

func GetFileContent(ctx *gin.Context) {
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	content, err := ops.GetFileContent(createCancellableContext(ctx), path[0])
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err)
	} else {
		RespondBinary(ctx, http.StatusOK, content)
	}
}
