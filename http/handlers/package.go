package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPackageInfo(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		RespondMissingParams(ctx, []string{"name"})
		return
	}

	pkgInfo, err := ops.GetPackageInfo(createStandardContext(ctx), name)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pkgInfo)
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

func InstallPackage(ctx *gin.Context) {
	// get params
	byFile := false
	nameOrPath, ok := ctx.Params.Get("name")
	if !ok {
		path, ok := ctx.Request.URL.Query()["path"]
		if ok {
			nameOrPath = path[0]
			byFile = true
		} else {
			RespondMissingParams(ctx, []string{"path"})
			return
		}
	}
	async := false
	asyncParams, ok := ctx.Request.URL.Query()["async"]
	if ok {
		async, _ = strconv.ParseBool(asyncParams[0])
	}

	// create background job if async == true
	if async {
		params := make(map[string]interface{})
		params["nameOrPath"] = nameOrPath
		params["byFile"] = byFile
		params["upgradeIfAlreadyInstalled"] = false
		jobInfo, err := ops.StartJob(createStandardContext(ctx), ops.JobTypeInstallPackage, params)
		if nil == err {
			RespondSuccessfulJson(ctx, http.StatusAccepted, jobInfo)
		} else {
			RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		}
		return
	}

	// execute operation
	pkgInfo, err := ops.InstallPackage(createStandardContext(ctx), nameOrPath, byFile, false)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusCreated, pkgInfo)
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

func UpgradePackage(ctx *gin.Context) {
	byFile := false
	nameOrPath, ok := ctx.Params.Get("name")
	if !ok {
		path, ok := ctx.Request.URL.Query()["path"]
		if ok {
			nameOrPath = path[0]
			byFile = true
		} else {
			RespondMissingParams(ctx, []string{"path"})
			return
		}
	}
	async := false
	asyncParams, ok := ctx.Request.URL.Query()["async"]
	if ok {
		async, _ = strconv.ParseBool(asyncParams[0])
	}

	// create background job if async == true
	if async {
		params := make(map[string]interface{})
		params["nameOrPath"] = nameOrPath
		params["byFile"] = byFile
		params["upgradeIfAlreadyInstalled"] = true
		jobInfo, err := ops.StartJob(createStandardContext(ctx), ops.JobTypeInstallPackage, params)
		if nil == err {
			RespondSuccessfulJson(ctx, http.StatusAccepted, jobInfo)
		} else {
			RespondFailedJson(ctx, http.StatusBadRequest, err, nil)
		}
		return
	}

	pkgInfo, err := ops.InstallPackage(createStandardContext(ctx), nameOrPath, byFile, true)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pkgInfo)
	} else {
		RespondFailedJson(ctx, http.StatusInternalServerError, err, nil)
	}
}

func UninstallPackage(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		RespondMissingParams(ctx, []string{"name"})
		return
	}

	err := ops.UninstallPackage(createStandardContext(ctx), name)
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, nil)
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
