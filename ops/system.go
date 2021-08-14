package ops

import (
	"context"
	"expinc/sunagent/log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
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

type CpuStat struct {
	Usages []float64 `json:"usages"`
	Load1  float64   `json:"load1"`
	Load5  float64   `json:"load5"`
	Load15 float64   `json:"load15"`
}

type MemStat struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
	Free      uint64 `json:"free"`
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

func GetCpuStat(ctx context.Context, perCpu bool) (stat CpuStat, err error) {
	usages, err := cpu.Percent(time.Second, perCpu)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}
	stat.Usages = usages

	loads, err := load.Avg()
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}
	stat.Load1 = loads.Load1
	stat.Load5 = loads.Load5
	stat.Load15 = loads.Load15
	return
}

func GetMemStat(ctx context.Context) (stat MemStat, err error) {
	memStat, err := mem.VirtualMemory()
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}
	stat.Total = memStat.Total
	stat.Available = memStat.Available
	stat.Used = memStat.Used
	stat.Free = memStat.Free
	return
}
