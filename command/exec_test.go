package command

import (
	"context"
	"expinc/sunagent/common"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/stretchr/testify/assert"
)

var (
	sleepScript         string
	outputScript        string
	failScript          string
	combinedOutput      string
	combinedFail        string
	separateOutputOut   string
	separateOutputErr   string
	separateFailOut     string
	separateFailErr     string
	sleepOrdinaryOutput string
	sleepOutputPattern  string
)

func init() {
	sleepScript = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "dummy-proc.py")
	outputScript = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "output-script.py")
	failScript = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "fail-script.py")

	combinedOutput = "stdout 1\nstderr 1\nstdout 2\nstderr 2\n"
	combinedFail = "start script\nexit with 1\n"
	separateOutputOut = "stdout 1\nstdout 2\n"
	separateOutputErr = "stderr 1\nstderr 2\n"
	sleepOutputPattern = "sleeping %d seconds\n"
	separateFailOut = "start script\n"
	separateFailErr = "exit with 1\n"

	info, _ := host.Info()
	if "windows" == info.OS {
		combinedOutput = strings.ReplaceAll(combinedOutput, "\n", "\r\n")
		combinedFail = strings.ReplaceAll(combinedFail, "\n", "\r\n")
		separateOutputOut = strings.ReplaceAll(separateOutputOut, "\n", "\r\n")
		separateOutputErr = strings.ReplaceAll(separateOutputErr, "\n", "\r\n")
		sleepOutputPattern = strings.ReplaceAll(sleepOutputPattern, "\n", "\r\n")
		separateFailOut = strings.ReplaceAll(separateFailOut, "\n", "\r\n")
		separateFailErr = strings.ReplaceAll(separateFailErr, "\n", "\r\n")
	}
}

func TestCheckCall_Ordinary(t *testing.T) {
	err := CheckCall("python3", []string{sleepScript, "1"}, DefaultTimeout)
	assert.Equal(t, nil, err)
}

func TestCheckCall_Failure(t *testing.T) {
	err := CheckCall("python3", []string{sleepScript, "t"}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
}

func TestCheckCall_Timeout(t *testing.T) {
	err := CheckCall("python3", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
}

func TestCheckCombinedOutput_Ordinary(t *testing.T) {
	output, err := CheckCombinedOutput("python3", []string{outputScript}, DefaultTimeout)
	assert.Equal(t, nil, err)
	assert.Equal(t, combinedOutput, string(output))
}

func TestCheckCombinedOutput_Failure(t *testing.T) {
	output, err := CheckCombinedOutput("python3", []string{failScript}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
	assert.Equal(t, combinedFail, string(output))
}

func TestCheckCombinedOutput_Timeout(t *testing.T) {
	output, err := CheckCombinedOutput("python3", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
	assert.Equal(t, fmt.Sprintf(sleepOutputPattern, 10), string(output))
}

func TestCheckSeparateOutput_Ordinary(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python3", []string{outputScript}, DefaultTimeout)
	assert.Equal(t, nil, err)
	assert.Equal(t, separateOutputOut, string(stdout))
	assert.Equal(t, separateOutputErr, string(stderr))
}

func TestCheckSeparateOutput_Failure(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python3", []string{failScript}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
	assert.Equal(t, separateFailOut, string(stdout))
	assert.Equal(t, separateFailErr, string(stderr))
}

func TestCheckSeparateOutput_Timeout(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python3", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
	assert.Equal(t, fmt.Sprintf(sleepOutputPattern, 10), string(stdout))
	assert.Equal(t, "", string(stderr))
}

func TestCheckCallContext_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var err error
	startTime := time.Now()
	go func(ctx context.Context) {
		err = CheckCallContext(ctx, "python3", []string{sleepScript, "10"}, DefaultTimeout)
	}(ctx)
	cancel()
	endTime := time.Now()

	assert.Equal(t, nil, err)
	assert.Less(t, endTime.Sub(startTime), 5*time.Second)
}

func TestCheckCombinedOutput_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var err error
	startTime := time.Now()
	go func(ctx context.Context) {
		_, err = CheckCombinedOutputContext(ctx, "python3", []string{sleepScript, "10"}, DefaultTimeout)
	}(ctx)
	cancel()
	endTime := time.Now()

	assert.Equal(t, nil, err)
	assert.Less(t, endTime.Sub(startTime), 5*time.Second)
}

func TestCheckSeparateOutput_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var err error
	startTime := time.Now()
	go func(ctx context.Context) {
		_, _, err = CheckSeparateOutputContext(ctx, "python3", []string{sleepScript, "10"}, DefaultTimeout)
	}(ctx)
	cancel()
	endTime := time.Now()

	assert.Equal(t, nil, err)
	assert.Less(t, endTime.Sub(startTime), 5*time.Second)
}
