package command

import (
	"bytes"
	"context"
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

func logOnExecStart(ctx context.Context, program string, args []string, timeout time.Duration) {
	content := fmt.Sprintf("Executing commad: program=%s, args=%s, timeout=%v", program, strings.Join(args, ", "), timeout)
	log.DebugCtx(ctx, content)
}

func logOnExecFinish(ctx context.Context, program string, args []string, timeout time.Duration, err error, stdOut string, stdErr string) {
	content := fmt.Sprintf("Executed commad: program=%s, args=%s, timeout=%v, err=%v, stdout=%s, stderr=%s",
		program,
		strings.Join(args, ", "),
		timeout,
		err,
		stdOut,
		stdErr,
	)
	log.DebugCtx(ctx, content)
}

// Execute command "program" with arguments "args".
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
// The provided context is used to kill the process (by calling os.Process.Kill)
// if the context becomes done before the command completes on its own.
func CheckCallContext(ctx context.Context, program string, args []string, timeout time.Duration) error {
	logOnExecStart(ctx, program, args, timeout)
	cmd := exec.CommandContext(ctx, program, args...)
	err := cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	logOnExecFinish(ctx, program, args, timeout, err, "<ignored>", "<ignored>")
	return err
}

// Execute command "program" with arguments "args".
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCall(program string, args []string, timeout time.Duration) error {
	return CheckCallContext(context.Background(), program, args, timeout)
}

// Execute command "program" with arguments "args".
// Return combined stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
// The provided context is used to kill the process (by calling os.Process.Kill)
// if the context becomes done before the command completes on its own.
func CheckCombinedOutputContext(ctx context.Context, program string, args []string, timeout time.Duration) (output []byte, err error) {
	logOnExecStart(ctx, program, args, timeout)
	cmd := exec.CommandContext(ctx, program, args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err = cmd.Start()
	if nil == err {
		err = waitForTimeout(cmd, timeout)
	}
	output = buffer.Bytes()
	logOnExecFinish(ctx, program, args, timeout, err, string(output), "<combined with output>")
	return
}

// Execute command "program" with arguments "args".
// Return combined stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckCombinedOutput(program string, args []string, timeout time.Duration) (output []byte, err error) {
	return CheckCombinedOutputContext(context.Background(), program, args, timeout)
}

// Execute command "program" with arguments "args".
// Return separate stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
// The provided context is used to kill the process (by calling os.Process.Kill)
// if the context becomes done before the command completes on its own.
func CheckSeparateOutputContext(ctx context.Context, program string, args []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	logOnExecStart(ctx, program, args, timeout)
	cmd := exec.CommandContext(ctx, program, args...)
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
	logOnExecFinish(ctx, program, args, timeout, err, string(stdout), string(stderr))
	return
}

// Execute command "program" with arguments "args".
// Return separate stdout & stderr if succeeds.
// Stop the command process after "timeout" and return common.Error with code ErrorTimeout.
// If command execution returns non-zero, this function will return a *os/exec.ExitError.
func CheckSeparateOutput(program string, args []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	return CheckSeparateOutputContext(context.Background(), program, args, timeout)
}
