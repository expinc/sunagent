package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
	"os/exec"
	"regexp"
)

type PackageInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Summary      string `json:"summary"`
}

func getPackageInfoFromDpkg(ctx context.Context, cmdOutput string) (pkgInfo PackageInfo, err error) {
	re := regexp.MustCompile(`Package: .*`)
	pkgInfo.Name = re.FindString(cmdOutput)[9:]
	re = regexp.MustCompile(`Version: .*`)
	pkgInfo.Version = re.FindString(cmdOutput)[9:]
	re = regexp.MustCompile(`Architecture: .*`)
	pkgInfo.Architecture = re.FindString(cmdOutput)[14:]
	re = regexp.MustCompile(`Description: .*`)
	pkgInfo.Summary = re.FindString(cmdOutput)[13:]
	return
}

func getPackageInfoFromRpm(ctx context.Context, cmdOutput string) (pkgInfo PackageInfo, err error) {
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
