package handlers

import (
	"context"
	"expinc/sunagent/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// create standard context from gin.Context to support request cancellation like Done()
// functions called by gin handlers should use context this function returns instead of gin.Context
func createStandardContext(ginCtx *gin.Context) context.Context {
	if nil != ginCtx {
		traceId := ginCtx.Value(common.TraceIdContextKey)
		return context.WithValue(context.Background(), common.TraceIdContextKey, traceId)
	} else {
		return context.Background()
	}
}

type JsonResponse struct {
	Successful bool        `json:"successful"`
	Status     int         `json:"status"`
	TraceId    string      `json:"traceId"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
}

func RespondSuccessfulJson(context *gin.Context, status int, data interface{}) {
	response := &JsonResponse{
		Successful: true,
		Status:     status,
		TraceId:    context.Value(common.TraceIdContextKey).(string),
		Data:       data,
		Error:      "",
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}

func RespondFailedJson(context *gin.Context, status int, err error, data interface{}) {
	response := &JsonResponse{
		Successful: false,
		Status:     status,
		TraceId:    context.Value(common.TraceIdContextKey).(string),
		Data:       data,
		Error:      err.Error(),
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}

func RespondMissingParams(context *gin.Context, params []string) {
	err := common.NewError(common.ErrorInvalidParameter, "Missing parameter: "+strings.Join(params, ", "))
	RespondFailedJson(context, http.StatusBadRequest, err, nil)
}

func RespondBinary(context *gin.Context, status int, content []byte) {
	context.Set("status", status)
	context.Data(status, "application/octet-stream", content)
}
