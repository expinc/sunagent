package command

import (
	"bytes"
	"os/exec"
	"time"
)

const (
	DefaultTimeout = 60 * time.Second
	NoTimeout      = (1<<63 - 1) * time.Nanosecond
)

// Execute command "program" with arguments "args".
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCall(program string, args []string, timeout time.Duration) error {
	cmd := exec.Command(program, args...)
	err := cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	return err
}

// Execute command "program" with arguments "args".
// Return combined stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCombinedOutput(program string, args []string, timeout time.Duration) (output []byte, err error) {
	cmd := exec.Command(program, args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err = cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	output = buffer.Bytes()
	return
}

// Execute command "program" with arguments "args".
// Return separate stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckSeparateOutput(program string, args []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
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
	return
}
