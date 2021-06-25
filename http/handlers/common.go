package handlers

import (
	"context"
	"expinc/sunagent/common"

	"github.com/gin-gonic/gin"
)

// create standard context from gin.Context to support request cancellation like Done()
// functions called by handlers should use context this function returns instead of gin.Context
func createCancellableContext(ginCtx *gin.Context) context.Context {
	if nil != ginCtx {
		traceId := ginCtx.Value(common.TraceIdContextKey)
		return context.WithValue(context.Background(), common.TraceIdContextKey, traceId)
	} else {
		return context.Background()
	}
}

type TextualResponse struct {
	Successful bool        `json:"successful"`
	Status     int         `json:"status"`
	TraceId    string      `json:"traceId"`
	Data       interface{} `json:"data"`
}

func respondSuccessfulJson(context *gin.Context, status int, data interface{}) {
	response := &TextualResponse{
		Successful: true,
		Status:     status,
		TraceId:    context.Value(common.TraceIdContextKey).(string),
		Data:       data,
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}

func RespondFailedJson(context *gin.Context, status int, err error) {
	response := &TextualResponse{
		Successful: false,
		Status:     status,
		TraceId:    context.Value(common.TraceIdContextKey).(string),
		Data:       err.Error(),
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}
