package command

import (
	"bytes"
	"expinc/sunagent/log"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	DefaultTimeout = 60 * time.Second
	NoTimeout      = (1<<63 - 1) * time.Nanosecond
)

func logOnExecStart(program string, args []string, timeout time.Duration) {
	content := fmt.Sprintf("Executing commad: program=%s, args=%s, timeout=%v", program, strings.Join(args, ", "), timeout)
	log.Debug(content)
}

func logOnExecFinish(program string, args []string, timeout time.Duration, err error, stdOut string, stdErr string) {
	content := fmt.Sprintf("Executed commad: program=%s, args=%s, timeout=%v, err=%v, stdout=%s, stderr=%s",
		program,
		strings.Join(args, ", "),
		timeout,
		err,
		stdOut,
		stdErr,
	)
	log.Debug(content)
}

// Execute command "program" with arguments "args".
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCall(program string, args []string, timeout time.Duration) error {
	logOnExecStart(program, args, timeout)
	cmd := exec.Command(program, args...)
	err := cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	logOnExecFinish(program, args, timeout, err, "<ignored>", "<ignored>")
	return err
}

// Execute command "program" with arguments "args".
// Return combined stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCombinedOutput(program string, args []string, timeout time.Duration) (output []byte, err error) {
	logOnExecStart(program, args, timeout)
	cmd := exec.Command(program, args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err = cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	output = buffer.Bytes()
	logOnExecFinish(program, args, timeout, err, string(output), "<combined with output>")
	return
}

// Execute command "program" with arguments "args".
// Return separate stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckSeparateOutput(program string, args []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	logOnExecStart(program, args, timeout)
	cmd := exec.Command(program, args...)
	var bufferOut bytes.Buffer
	var bufferErr bytes.Buffer
	cmd.Stdout = &bufferOut
	cmd.Stderr = &bufferErr

	err = cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	stdout = bufferOut.Bytes()
	stderr = bufferErr.Bytes()
	logOnExecFinish(program, args, timeout, err, string(stdout), string(stderr))
	return
}
