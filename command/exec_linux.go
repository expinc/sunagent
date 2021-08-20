package command

import (
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"os/exec"
	"syscall"
	"time"
)

const gracefulKillTime = 5 * time.Second

func waitForTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	// kill process after timeout
	killFunc := func() {
		err := cmd.Process.Signal(syscall.SIGTERM)
		if nil != err {
			log.Error(err)
			return
		}

		forciblyKillFunc := func() {
			err := cmd.Process.Signal(syscall.SIGKILL)
			if nil != err {
				log.Error(err)
				return
			}
		}
		time.AfterFunc(gracefulKillTime, forciblyKillFunc)
	}
	killTimer := time.AfterFunc(timeout, killFunc)

	// wait for command finish execution and stop killTimer
	err := cmd.Wait()
	notTimeout := killTimer.Stop()
	if nil != err {
		// if killFunc has already been called, it should be timeout
		if !notTimeout {
			err = common.NewError(common.ErrorTimeout, "Command execution timeout")
		}
	}

	return err
}
