package fuckfetch

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type SystemInfo struct {
	// OS
	OS, Platform, KernelVersion string
	Uptime                                time.Duration
	// CPU
	CPUModel   string
	CPUCores   int // physical
	CPUThreads int // logical
	CPUUsage   float64
	// Memory
	MemTotal, MemUsed, MemAvailable uint64
	MemUsage                        float64
	// Swap
	SwapTotal, SwapUsed uint64
	SwapUsage           float64
	// Disk
	DiskTotal, DiskUsed, DiskFree uint64
	DiskUsage                     float64
	// Network
	NetBytesSent, NetBytesRecv uint64
	// GPU / NPU (optional)
	GPUInfo string
	NPUInfo string
}

func GatherSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}

	// Host
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("host info: %w", err)
	}
	info.OS = hostInfo.OS
	info.Platform = hostInfo.Platform + " " + hostInfo.PlatformVersion
	info.KernelVersion = hostInfo.KernelVersion
	info.Uptime = time.Duration(hostInfo.Uptime) * time.Second

	// CPU info
	cpuInfos, err := cpu.Info()
	if err == nil && len(cpuInfos) > 0 {
		info.CPUModel = cpuInfos[0].ModelName
		info.CPUCores = int(cpuInfos[0].Cores)
	}
	logicalCount, err := cpu.Counts(true)
	if err == nil {
		info.CPUThreads = logicalCount
	}
	physicalCount, err := cpu.Counts(false)
	if err == nil && physicalCount > 0 {
		info.CPUCores = physicalCount
	}

	// CPU usage (500ms sample)
	percents, err := cpu.Percent(500*time.Millisecond, false)
	if err == nil && len(percents) > 0 {
		info.CPUUsage = percents[0]
	}

	// Memory
	vmem, err := mem.VirtualMemory()
	if err == nil {
		info.MemTotal = vmem.Total
		info.MemUsed = vmem.Used
		info.MemAvailable = vmem.Available
		info.MemUsage = vmem.UsedPercent
	}

	// Swap
	swapMem, err := mem.SwapMemory()
	if err == nil {
		info.SwapTotal = swapMem.Total
		info.SwapUsed = swapMem.Used
		info.SwapUsage = swapMem.UsedPercent
	}

	// Disk
	diskUsage, err := disk.Usage("/")
	if err == nil {
		info.DiskTotal = diskUsage.Total
		info.DiskUsed = diskUsage.Used
		info.DiskFree = diskUsage.Free
		info.DiskUsage = diskUsage.UsedPercent
	}

	// Network
	counters, err := net.IOCounters(false)
	if err == nil && len(counters) > 0 {
		info.NetBytesSent = counters[0].BytesSent
		info.NetBytesRecv = counters[0].BytesRecv
	}

	// GPU
	info.GPUInfo = gatherGPUInfo()

	// NPU
	info.NPUInfo = gatherNPUInfo()

	return info, nil
}

func gatherGPUInfo() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "lspci").Output()
	if err != nil {
		return "N/A"
	}

	var gpus []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "VGA") || strings.Contains(line, "3D controller") || strings.Contains(line, "Display controller") {
			if idx := strings.Index(line, ": "); idx != -1 {
				gpus = append(gpus, strings.TrimSpace(line[idx+2:]))
			}
		}
	}
	if len(gpus) == 0 {
		return "N/A"
	}
	return strings.Join(gpus, "\n")
}

func gatherNPUInfo() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "lspci").Output()
	if err != nil {
		return "N/A"
	}

	var npus []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "Processing accelerators") || strings.Contains(line, "accel") {
			if idx := strings.Index(line, ": "); idx != -1 {
				npus = append(npus, strings.TrimSpace(line[idx+2:]))
			}
		}
	}
	if len(npus) == 0 {
		return "N/A"
	}
	return strings.Join(npus, "\n")
}
