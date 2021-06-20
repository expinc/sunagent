package handlers

import "github.com/gin-gonic/gin"

const TraceIdHeader = "traceId"

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
		TraceId:    context.Value(TraceIdHeader).(string),
		Data:       data,
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}

func respondFailedJson(context *gin.Context, status int, err error) {
	response := &TextualResponse{
		Successful: false,
		Status:     status,
		TraceId:    context.Value(TraceIdHeader).(string),
		Data:       err.Error(),
	}
	context.Set("status", status)
	context.JSON(response.Status, response)
}
