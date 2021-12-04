package handlers

import (
	"expinc/sunagent/ops"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetJobInfo(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		RespondMissingParams(ctx, []string{"id"})
		return
	}

	info, err := ops.GetJobInfo(createStandardContext(ctx), id)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, info)
	} else {
		RespondFailedJson(ctx, http.StatusNotFound, err, nil)
	}
}
