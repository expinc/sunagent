package command

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"expinc/sunagent/common"

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
	sleepScript = path.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "dummy-proc.py")
	outputScript = path.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "output-script.py")
	failScript = path.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "fail-script.py")

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
	err := CheckCall("python", []string{sleepScript, "1"}, DefaultTimeout)
	assert.Equal(t, nil, err)
}

func TestCheckCall_Failure(t *testing.T) {
	err := CheckCall("python", []string{sleepScript, "t"}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
}

func TestCheckCall_Timeout(t *testing.T) {
	err := CheckCall("python", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
}

func TestCheckCombinedOutput_Ordinary(t *testing.T) {
	output, err := CheckCombinedOutput("python", []string{outputScript}, DefaultTimeout)
	assert.Equal(t, nil, err)
	assert.Equal(t, combinedOutput, string(output))
}

func TestCheckCombinedOutput_Failure(t *testing.T) {
	output, err := CheckCombinedOutput("python", []string{failScript}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
	assert.Equal(t, combinedFail, string(output))
}

func TestCheckCombinedOutput_Timeout(t *testing.T) {
	output, err := CheckCombinedOutput("python", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
	assert.Equal(t, fmt.Sprintf(sleepOutputPattern, 10), string(output))
}

func TestCheckSeparateOutput_Ordinary(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python", []string{outputScript}, DefaultTimeout)
	assert.Equal(t, nil, err)
	assert.Equal(t, separateOutputOut, string(stdout))
	assert.Equal(t, separateOutputErr, string(stderr))
}

func TestCheckSeparateOutput_Failure(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python", []string{failScript}, DefaultTimeout)
	assert.IsType(t, &exec.ExitError{}, err)
	assert.Equal(t, separateFailOut, string(stdout))
	assert.Equal(t, separateFailErr, string(stderr))
}

func TestCheckSeparateOutput_Timeout(t *testing.T) {
	stdout, stderr, err := CheckSeparateOutput("python", []string{sleepScript, "10"}, 1*time.Second)
	assert.IsType(t, common.Error{}, err)
	assert.Equal(t, common.ErrorTimeout, err.(common.Error).Code())
	assert.Equal(t, fmt.Sprintf(sleepOutputPattern, 10), string(stdout))
	assert.Equal(t, "", string(stderr))
}