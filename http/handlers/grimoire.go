package handlers

import (
	"expinc/sunagent/ops"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGrimoire(ctx *gin.Context) {
	osType, ok := ctx.Params.Get("osType")
	if !ok {
		RespondMissingParams(ctx, []string{"osType"})
		return
	}

	output, err := ops.GetGrimoireAsYaml(createStandardContext(ctx), osType)
	if nil == err {
		RespondBinary(ctx, http.StatusOK, output)
	} else {
		RespondFailedJson(ctx, http.StatusNotFound, err, nil)
	}
}
