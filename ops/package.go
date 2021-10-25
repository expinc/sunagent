package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type PackageInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Summary      string `json:"summary"`
}

func getPackageInfoFromDpkg(ctx context.Context, cmdOutput string) (pkgInfo PackageInfo, err error) {
	re := regexp.MustCompile(`Package: .*`)
	entry := re.FindString(cmdOutput)
	value := strings.Split(entry, ":")[1]
	pkgInfo.Name = strings.TrimSpace(value)
	re = regexp.MustCompile(`Version: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Version = strings.TrimSpace(value)
	re = regexp.MustCompile(`Architecture: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Architecture = strings.TrimSpace(value)
	re = regexp.MustCompile(`Description: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Summary = strings.TrimSpace(value)
	return
}

func getPackageInfoFromRpm(ctx context.Context, cmdOutput string) (pkgInfo PackageInfo, err error) {
	re := regexp.MustCompile(`Name\s*: .*`)
	entry := re.FindString(cmdOutput)
	value := strings.Split(entry, ":")[1]
	pkgInfo.Name = strings.TrimSpace(value)
	re = regexp.MustCompile(`Version\s*: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Version = strings.TrimSpace(value)
	re = regexp.MustCompile(`Architecture\s*: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Architecture = strings.TrimSpace(value)
	re = regexp.MustCompile(`Summary\s*: .*`)
	entry = re.FindString(cmdOutput)
	value = strings.Split(entry, ":")[1]
	pkgInfo.Summary = strings.TrimSpace(value)
	return
}

func getPackageInfoFromCmdOutput(ctx context.Context, cmdOutput string) (pkgInfo PackageInfo, err error) {
	switch nodeInfo.OsFamily {
	case "debian":
		pkgInfo, err = getPackageInfoFromDpkg(ctx, cmdOutput)
	case "rhel":
		pkgInfo, err = getPackageInfoFromRpm(ctx, cmdOutput)
	case "suse":
		pkgInfo, err = getPackageInfoFromRpm(ctx, cmdOutput)
	}
	return
}

func GetPackageInfo(ctx context.Context, name string) (pkgInfo PackageInfo, err error) {
	if "linux" != nodeInfo.OsType {
		err = common.NewError(common.ErrorNotImplemented, "Package operations only support Linux systems")
		return
	}

	var output []byte
	output, err = castGrimoireArcane("get-package", name)
	if nil != err {
		log.ErrorCtx(ctx, err)

		// if the command returns non-zero, it indicates that the package is not installed
		_, cmdFailed := err.(*exec.ExitError)
		if cmdFailed {
			msg := fmt.Sprintf("Package \"%s\" is not installed", name)
			err = common.NewError(common.ErrorNotFound, msg)
		}
		return
	}

	outputStr := string(output)
	pkgInfo, err = getPackageInfoFromCmdOutput(ctx, outputStr)
	if nil != err {
		log.ErrorCtx(ctx, err)
	}
	return
}

func InstallPackageByName(ctx context.Context, name string) (pkgInfo PackageInfo, err error) {
	if "linux" != nodeInfo.OsType {
		err = common.NewError(common.ErrorNotImplemented, "Package operations only support Linux systems")
		return
	}

	_, err = GetPackageInfo(ctx, name)
	if nil == err {
		msg := fmt.Sprintf("Package \"%s\" is already installed", name)
		err = common.NewError(common.ErrorUnexpected, msg)
		return
	}

	output, err := castGrimoireArcane("install-package", name)
	if nil != err {
		log.ErrorCtx(ctx, err)

		_, cmdFailed := err.(*exec.ExitError)
		if cmdFailed {
			err = common.NewError(common.ErrorUnexpected, string(output))
		}
		return
	}

	pkgInfo, err = GetPackageInfo(ctx, name)
	return
}
