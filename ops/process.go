package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
)

type ProcInfo struct {
	Pid            int32     `json:"pid"`
	Name           string    `json:"name"`
	Cmd            string    `json:"cmd"`
	StartTime      time.Time `json:"startTime"`
	ElapsedSeconds int64     `json:"elapsedSeconds"`
	Owner          string    `json:"owner"`
}

func proc2ProcInfo(ctx context.Context, proc *process.Process) (info ProcInfo) {
	// when some of the fields fail to retrieve, just log the error instead of fail the caller
	// because the process existance has already been confirmed
	var err1 error
	info.Pid = proc.Pid
	info.Name, err1 = proc.Name()
	if nil != err1 {
		log.ErrorCtx(ctx, err1)
	}
	info.Cmd, err1 = proc.Cmdline()
	if nil != err1 {
		log.ErrorCtx(ctx, err1)
	}
	startSecTimestamp, err1 := proc.CreateTime()
	if nil != err1 {
		log.ErrorCtx(ctx, err1)
	}
	startSecTimestamp = startSecTimestamp * int64(time.Millisecond)
	info.StartTime = time.Unix(0, startSecTimestamp)
	info.ElapsedSeconds = int64(time.Now().Sub(info.StartTime).Seconds())
	info.Owner, err1 = proc.Username()
	if nil != err1 {
		log.ErrorCtx(ctx, err1)
	}
	return
}

func GetProcInfoByPid(ctx context.Context, pid int32) (info ProcInfo, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Getting process info: pid=%v", pid))
	proc, err := process.NewProcess(pid)
	if nil != err {
		log.ErrorCtx(ctx, err)
		err = common.NewError(common.ErrorNotFound, err.Error())
		return
	}

	info = proc2ProcInfo(ctx, proc)
	return
}

func GetProcInfosByName(ctx context.Context, name string) (infos []ProcInfo, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Getting process info: name=%v", name))
	procs, err := process.Processes()
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}

	for _, proc := range procs {
		procName, err1 := proc.Name()
		if nil != err1 {
			log.ErrorCtx(ctx, err1)
			continue
		}

		if procName == name {
			infos = append(infos, proc2ProcInfo(ctx, proc))
		}
	}

	if 0 == len(infos) {
		err = common.NewError(common.ErrorNotFound, fmt.Sprintf("No process named %s", name))
		log.ErrorCtx(ctx, err)
	}
	return
}

func KillProcByPid(ctx context.Context, pid int32, signal int) (err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Killing process: pid=%v, signal=%v", pid, signal))
	proc, err := process.NewProcess(pid)
	if nil != err {
		log.ErrorCtx(ctx, err)
		err = common.NewError(common.ErrorNotFound, err.Error())
		return
	}

	if 0 == signal {
		err = proc.Kill()
	} else {
		err = proc.SendSignal(syscall.Signal(signal))
	}
	if nil != err {
		log.ErrorCtx(ctx, err)
	}
	return
}
