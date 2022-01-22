package handlers

import (
	"expinc/sunagent/common"
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

	metas, err := ops.GetFileMetas(createStandardContext(ctx), path[0], listIfDir)
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
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

	content, err := ops.GetFileContent(createStandardContext(ctx), path[0])
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
	} else {
		RespondBinary(ctx, http.StatusOK, content)
	}
}

func writeFile(ctx *gin.Context, overwrite bool) {
	// Get parameter "path"
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	// Get parameter "isDir"
	isDirStr, ok := ctx.Request.URL.Query()["isDir"]
	isDir := false
	if ok {
		isDir, _ = strconv.ParseBool(isDirStr[0])
	}

	// Get request body
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(common.FileUploadMaxBytes))
	content, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		return
	}

	// Execute operation & render response
	meta, err := ops.WriteFile(createStandardContext(ctx), path[0], content, isDir, overwrite)
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	} else {
		status := http.StatusOK
		if !overwrite {
			status = http.StatusCreated
		}
		RespondSuccessfulJson(ctx, status, meta)
	}
}

func CreateFile(ctx *gin.Context) {
	writeFile(ctx, false)
}

func OverwriteFile(ctx *gin.Context) {
	writeFile(ctx, true)
}

func AppendFile(ctx *gin.Context) {
	// Get parameters
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	// Read content
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(common.FileUploadMaxBytes))
	content, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		return
	}

	// Render response
	meta, err := ops.AppendFile(createStandardContext(ctx), path[0], content)
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, meta)
	}
}

func DeleteFile(ctx *gin.Context) {
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	recursiveStr, ok := ctx.Request.URL.Query()["recursive"]
	recursive := false
	if ok {
		recursive, _ = strconv.ParseBool(recursiveStr[0])
	}

	err := ops.DeleteFile(createStandardContext(ctx), path[0], recursive)
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
	} else {
		RespondSuccessfulJson(ctx, http.StatusOK, nil)
	}
}
