package handlers

import (
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

	metas, err := ops.GetFileMetas(path[0], listIfDir)
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, metas)
	}
}
