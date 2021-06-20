package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Info struct {
	Version string `json:"version"`
}

func GetInfo(context *gin.Context) {
	info := Info{
		Version: "1.0.0",
	}
	respondSuccessfulJson(context, http.StatusOK, info)
}
