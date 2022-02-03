package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"expinc/sunagent/ops"
	"fmt"
	"io"
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
	// Get parameter
	path, ok := ctx.Request.URL.Query()["path"]
	if !ok {
		RespondMissingParams(ctx, []string{"path"})
		return
	}

	// Get file size
	metas, err := ops.GetFileMetas(createStandardContext(ctx), path[0], false)
	if nil != err {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		RespondFailedJson(ctx, status, err, nil)
		return
	}
	ctx.Header("Content-Length", fmt.Sprint(metas[0].Size))

	// Read from file
	reader, err := ops.NewFileStreamReader(ctx, path[0])
	defer reader.Close()
	if nil != err {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	// FIXME: what should be the good chan size?
	chanStream := make(chan []byte, 10)
	var readErr error
	go func() {
		defer close(chanStream)
		buffer := make([]byte, common.FileTransferSizeLimit)
		for {
			n, readErr := reader.Read(buffer)
			if nil == readErr {
				chanStream <- buffer[:n]
			} else {
				break
			}
		}
		if io.EOF == readErr {
			// EOF means finished reading and no problem
			readErr = nil
		}
	}()

	// Write to response
	var writeErr error
	ctx.Stream(func(responseWriter io.Writer) bool {
		chunk, ok := <-chanStream
		if ok {
			_, writeErr = responseWriter.Write(chunk)
		}
		return ok && nil == writeErr
	})

	// Handle errors
	if nil == readErr && nil == writeErr {
		ctx.Status(http.StatusOK)
	} else {
		err = writeErr
		if nil != readErr {
			err = readErr
		}
		log.ErrorCtx(ctx, err)
		ctx.Status(http.StatusInternalServerError)
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
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(common.FileTransferSizeLimit))
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
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(common.FileTransferSizeLimit))
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
