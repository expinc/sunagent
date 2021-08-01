package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"net/http"

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
		internalError, ok := err.(common.Error)
		if ok && common.ErrorNotFound == internalError.Code() {
			status = http.StatusNotFound
		}

		RespondFailedJson(ctx, status, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, metas)
	}
}
