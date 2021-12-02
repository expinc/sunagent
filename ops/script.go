package ops

import (
	"context"
	"expinc/sunagent/command"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CombinedScriptResult struct {
	Output     string `json:"output"`
	ExitStatus int    `json:"exitStatus"`
	Error      string `json:"error"`
}

type SeparateScriptResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitStatus int    `json:"exitStatus"`
	Error      string `json:"error"`
}

const scriptLogLimit = 32

func execScript(ctx context.Context, program string, script string, waitSeconds int64, separateOutput bool) (result interface{}, err error) {
	// check parameters
	if "" == strings.TrimSpace(program) {
		err = common.NewError(common.ErrorInvalidParameter, "Parameter \"program\" should be non-empty")
		log.ErrorCtx(ctx, err)
		return
	}
	if "" == strings.TrimSpace(script) {
		err = common.NewError(common.ErrorInvalidParameter, "Parameter \"script\" should be non-empty")
		log.ErrorCtx(ctx, err)
		return
	}
	if 0 > waitSeconds {
		err = common.NewError(common.ErrorInvalidParameter, "Parameter \"waitSeconds\" should be non-negative")
		log.ErrorCtx(ctx, err)
		return
	}

	// put script file
	err = os.MkdirAll(common.TmpFolder, fs.ModeDir)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}
	scriptFilePath := filepath.Join(common.TmpFolder, uuid.New().String())
	params := []string{scriptFilePath}
	// windows cmd must have parameter /C to execute the .bat file as script
	// parameter /Q is also added to avoid printing command lines
	if ("cmd" == program || "cmd.exe" == program) && "windows" == nodeInfo.OsType {
		scriptFilePath += ".bat"
		params = []string{"/Q", "/C", scriptFilePath}
	}
	err = ioutil.WriteFile(scriptFilePath, []byte(script), fs.ModePerm)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}
	defer func() {
		err := os.Remove(scriptFilePath)
		if nil != err {
			log.ErrorCtx(ctx, err)
		}
	}()

	// execute script
	timeout := time.Duration(waitSeconds) * time.Second
	if 0 == waitSeconds {
		timeout = command.NoTimeout
	}
	if separateOutput {
		var stdout []byte
		var stderr []byte
		stdout, stderr, err = command.CheckSeparateOutput(program, params, timeout)
		separateResult := SeparateScriptResult{
			Stdout: string(stdout),
			Stderr: string(stderr),
		}
		if nil != err {
			separateResult.Error = err.Error()
			exitError, ok := err.(*exec.ExitError)
			if ok {
				separateResult.ExitStatus = exitError.ExitCode()
			}
		}
		result = separateResult
	} else {
		var output []byte
		output, err = command.CheckCombinedOutput(program, params, timeout)
		combinedResult := CombinedScriptResult{
			Output: string(output),
		}
		if nil != err {
			combinedResult.Error = err.Error()
			exitError, ok := err.(*exec.ExitError)
			if ok {
				combinedResult.ExitStatus = exitError.ExitCode()
			}
		}
		result = combinedResult
	}
	return result, err
}

func getScriptPreview(script string) string {
	var scriptSnippet string
	if len(script) <= scriptLogLimit {
		scriptSnippet = script
	} else {
		scriptSnippet = script[:scriptLogLimit] + "..."
	}
	return scriptSnippet
}

func ExecScriptWithCombinedOutput(ctx context.Context, program string, script string, waitSeconds int64) (result CombinedScriptResult, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Executing script to get combined output: program=%v, script=%q, waitSeconds=%v",
		program,
		getScriptPreview(script),
		waitSeconds,
	))
	var combinedResult interface{}
	combinedResult, err = execScript(ctx, program, script, waitSeconds, false)
	result, _ = combinedResult.(CombinedScriptResult)
	if nil != err {
		log.ErrorCtx(ctx, err)
	}
	return
}

func ExecScriptWithSeparateOutput(ctx context.Context, program string, script string, waitSeconds int64) (result SeparateScriptResult, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Executing script to get separate output: program=%v, script=%q, waitSeconds=%v",
		program,
		getScriptPreview(script),
		waitSeconds,
	))
	var separateResult interface{}
	separateResult, err = execScript(ctx, program, script, waitSeconds, true)
	result, _ = separateResult.(SeparateScriptResult)
	if nil != err {
		log.ErrorCtx(ctx, err)
	}
	return
}
