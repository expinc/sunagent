package handlers

import (
	"expinc/sunagent/ops"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

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
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, metas)
	}
}

func GetFileContent(ctx *gin.Context) {
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	content, err := ops.GetFileContent(createCancellableContext(ctx), path[0])
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err)
	} else {
		RespondBinary(ctx, http.StatusOK, content)
	}
}

func CreateFile(ctx *gin.Context) {
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	isDirStr, ok := ctx.Request.URL.Query()["isDir"]
	isDir := false
	if ok {
		isDir, _ = strconv.ParseBool(isDirStr[0])
	}

	content, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err)
	}

	meta, err := ops.WriteFile(createCancellableContext(ctx), path[0], content, isDir, false)
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, meta)
	}
}
