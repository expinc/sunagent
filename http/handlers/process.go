package handlers

import (
	"expinc/sunagent/common"
	"expinc/sunagent/ops"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func GetProcInfo(ctx *gin.Context) {
	// Get params
	pidOrNameStr, ok := ctx.Params.Get("pidOrName")
	if !ok {
		RespondMissingParams(ctx, []string{"pidOrName"})
		return
	}

	// Execute operation
	pid, err := strconv.ParseInt(pidOrNameStr, 10, 32)
	var infos []ops.ProcInfo
	if nil == err {
		var info ops.ProcInfo
		info, err = ops.GetProcInfoByPid(createStandardContext(ctx), int32(pid))
		if nil == err {
			infos = append(infos, info)
		}
	} else {
		infos, err = ops.GetProcInfosByName(createStandardContext(ctx), pidOrNameStr)
	}

	// Render response
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, infos)
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

func KillProc(ctx *gin.Context) {
	// Get pid or name
	pidOrNameStr, ok := ctx.Params.Get("pidOrName")
	if !ok {
		RespondMissingParams(ctx, []string{"pidOrName"})
		return
	}

	// Get signal
	signal := int64(0)
	if "linux" == runtime.GOOS {
		signal = int64(syscall.SIGTERM)
	}
	signalStr, ok := ctx.Request.URL.Query()["signal"]
	if ok {
		signal, _ = strconv.ParseInt(signalStr[0], 10, 32)
	}

	// Execute operation
	pid, err := strconv.ParseInt(pidOrNameStr, 10, 32)
	var pids []int64
	if nil == err {
		err = ops.KillProcByPid(createStandardContext(ctx), int32(pid), int(signal))
		if nil == err {
			pids = append(pids, pid)
		}
	} else {
		var infos []ops.ProcInfo
		infos, err = ops.GetProcInfosByName(createStandardContext(ctx), pidOrNameStr)
		for _, info := range infos {
			pid = int64(info.Pid)
			err = ops.KillProcByPid(createStandardContext(ctx), int32(pid), int(signal))
			if nil == err {
				pids = append(pids, pid)
			} else {
				errMsg := fmt.Sprintf("Failed to kill process %d because %s. Killed: %s", pid, err.Error(), fmt.Sprint(pids))
				err = common.NewError(common.ErrorUnexpected, errMsg)
				break
			}
		}
	}

	// Render response
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pids)
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

func TermProc(ctx *gin.Context) {
	// Get pid or name
	pidOrNameStr, ok := ctx.Params.Get("pidOrName")
	if !ok {
		RespondMissingParams(ctx, []string{"pidOrName"})
		return
	}

	// Get signal
	signal := int64(0)
	if "linux" == runtime.GOOS {
		signal = int64(syscall.SIGTERM)
	}

	// Execute operation
	pid, err := strconv.ParseInt(pidOrNameStr, 10, 32)
	var pids []int64
	if nil == err {
		err = ops.KillProcByPid(createStandardContext(ctx), int32(pid), int(signal))
		if nil == err {
			pids = append(pids, pid)
		}
	} else {
		var infos []ops.ProcInfo
		infos, err = ops.GetProcInfosByName(createStandardContext(ctx), pidOrNameStr)
		for _, info := range infos {
			pid = int64(info.Pid)
			err = ops.KillProcByPid(createStandardContext(ctx), int32(pid), int(signal))
			if nil == err {
				pids = append(pids, pid)
			} else {
				errMsg := fmt.Sprintf("Failed to kill process %d because %s. Killed: %s", pid, err.Error(), fmt.Sprint(pids))
				err = common.NewError(common.ErrorUnexpected, errMsg)
				break
			}
		}
	}

	// Render response
	if nil == err {
		RespondSuccessfulJson(ctx, http.StatusOK, pids)
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
