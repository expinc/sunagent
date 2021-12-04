package handlers

import (
	"expinc/sunagent/common"
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

func GetAllJobInfo(ctx *gin.Context) {
	result := ops.ListJobInfo(createStandardContext(ctx))
	RespondSuccessfulJson(ctx, http.StatusOK, result)
}

func CancelJob(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		RespondMissingParams(ctx, []string{"id"})
		return
	}

	info, err := ops.CancelJob(createStandardContext(ctx), id)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, info)
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
