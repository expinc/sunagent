package ops

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/host"
)

type NodeInfo struct {
	HostName      string    `json:"hostName"`
	BootTime      time.Time `json:"bootTime"`
	OsType        string    `json:"osType"`        // ex: freebsd, linux
	OsFamily      string    `json:"osFamily"`      // ex: debian, rhel
	OsVersion     string    `json:"osVersion"`     // operating system release version
	KernelVersion string    `json:"kernelVersion"` // operating system kernel version
	CpuArch       string    `json:"cpuArch"`       // ex: x86_64, aarch64
}

var nodeInfo NodeInfo

func init() {
	infoStat, _ := host.Info()
	nodeInfo = NodeInfo{
		HostName:      infoStat.Hostname,
		BootTime:      time.Unix(int64(infoStat.BootTime), 0),
		OsType:        infoStat.OS,
		OsFamily:      infoStat.PlatformFamily,
		OsVersion:     infoStat.PlatformVersion,
		KernelVersion: infoStat.KernelVersion,
		CpuArch:       infoStat.KernelArch,
	}
}

func GetNodeInfo(ctx context.Context) NodeInfo {
	return nodeInfo
}
