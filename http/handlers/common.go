package handlers

import (
	"expinc/sunagent/common"

	"github.com/gin-gonic/gin"
)

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

func respondFailedJson(context *gin.Context, status int, err error) {
	response := &TextualResponse{
		Successful: false,
		Status:     status,
		TraceId:    context.Value(common.TraceIdContextKey).(string),
		Data:       err.Error(),
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}
