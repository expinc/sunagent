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

func checkOsType() error {
	if "linux" != nodeInfo.OsType {
		return common.NewError(common.ErrorNotImplemented, "Package operations only support Linux systems")
	} else {
		return nil
	}
}

func getPackageArchiveInfo(ctx context.Context, path string) (pkgInfo PackageInfo, err error) {
	output, err := castGrimoireArcane("get-package-archive-info", path)
	if nil == err {
		pkgInfo, err = getPackageInfoFromCmdOutput(ctx, string(output))
	}
	return
}

func GetPackageInfo(ctx context.Context, name string) (pkgInfo PackageInfo, err error) {
	if err = checkOsType(); nil != err {
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

func InstallPackage(ctx context.Context, nameOrPath string, byFile bool) (pkgInfo PackageInfo, err error) {
	if err = checkOsType(); nil != err {
		return
	}

	// get package name
	var name string
	if !byFile {
		name = nameOrPath
	} else {
		var archiveInfo PackageInfo
		archiveInfo, err = getPackageArchiveInfo(ctx, nameOrPath)
		if nil != err {
			return
		}
		name = archiveInfo.Name
	}

	// check if package is already installed
	_, err = GetPackageInfo(ctx, name)
	if nil == err {
		msg := fmt.Sprintf("Package \"%s\" is already installed", name)
		err = common.NewError(common.ErrorUnexpected, msg)
		return
	}

	// install package
	output, err := castGrimoireArcane("install-package", nameOrPath)
	if nil != err {
		log.ErrorCtx(ctx, err)

		_, cmdFailed := err.(*exec.ExitError)
		if cmdFailed {
			err = common.NewError(common.ErrorUnexpected, string(output))
		}
		return
	}

	// return package info
	pkgInfo, err = GetPackageInfo(ctx, name)
	return
}
