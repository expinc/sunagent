package ops

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
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

type CpuInfo struct {
	ModelName string  `json:"modelName"`
	VendorId  string  `json:"vendorId"`
	Mhz       float64 `json:"mhz"`
	Count     int32   `json:"count"`
}

var (
	nodeInfo NodeInfo
	cpuInfo  CpuInfo
)

func init() {
	// init node info
	nodeInfoStat, _ := host.Info()
	nodeInfo = NodeInfo{
		HostName:      nodeInfoStat.Hostname,
		BootTime:      time.Unix(int64(nodeInfoStat.BootTime), 0),
		OsType:        nodeInfoStat.OS,
		OsFamily:      nodeInfoStat.PlatformFamily,
		OsVersion:     nodeInfoStat.PlatformVersion,
		KernelVersion: nodeInfoStat.KernelVersion,
		CpuArch:       nodeInfoStat.KernelArch,
	}

	// init CPU info
	cpuInfoStat, _ := cpu.Info()
	count := int32(0)
	for _, oneCpu := range cpuInfoStat {
		count += oneCpu.Cores
	}
	cpuInfo = CpuInfo{
		ModelName: cpuInfoStat[0].ModelName,
		VendorId:  cpuInfoStat[0].VendorID,
		Mhz:       cpuInfoStat[0].Mhz,
		Count:     count,
	}
}

func GetNodeInfo(ctx context.Context) NodeInfo {
	return nodeInfo
}

func GetCpuInfo(ctx context.Context) CpuInfo {
	return cpuInfo
}
