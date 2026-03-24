// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// GatherSystemInfo collects system information for neofetch-style output.
func GatherSystemInfo() (*model.SystemInfo, error) {
	info := &model.SystemInfo{}

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
	info.GPUInfo = gatherPCIDevices("VGA", "3D controller", "Display controller")

	// NPU
	info.NPUInfo = gatherPCIDevices("Processing accelerators", "accel")

	return info, nil
}

// gatherPCIDevices searches lspci output for lines matching any of the given keywords.
func gatherPCIDevices(keywords ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "lspci").Output()
	if err != nil {
		return "N/A"
	}

	var devices []string
	for _, line := range strings.Split(string(out), "\n") {
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				if idx := strings.Index(line, ": "); idx != -1 {
					devices = append(devices, strings.TrimSpace(line[idx+2:]))
				}
				break
			}
		}
	}
	if len(devices) == 0 {
		return "N/A"
	}
	return strings.Join(devices, "\n")
}
