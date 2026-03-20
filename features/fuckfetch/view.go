package fuckfetch

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildFuckfetchUI(info *SystemInfo) discord.ContainerComponent {
	title := discord.NewTextDisplay("### 🖥️ fuckfetch")

	osBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**OS:** %s (%s)\n**Kernel:** %s\n**Uptime:** %s",
		info.OS,
		info.Platform,
		info.KernelVersion,
		formatUptime(info.Uptime),
	))

	cpuBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**CPU:** %s (%dC/%dT)\n**Usage:** %s",
		info.CPUModel,
		info.CPUCores,
		info.CPUThreads,
		buildBar(info.CPUUsage),
	))

	memBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**RAM:** %s / %s %s\n**Swap:** %s / %s %s",
		formatBytes(info.MemUsed),
		formatBytes(info.MemTotal),
		buildBar(info.MemUsage),
		formatBytes(info.SwapUsed),
		formatBytes(info.SwapTotal),
		buildBar(info.SwapUsage),
	))

	diskNetBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**Disk (/):** %s / %s %s\n**Net ↑:** %s  **Net ↓:** %s",
		formatBytes(info.DiskUsed),
		formatBytes(info.DiskTotal),
		buildBar(info.DiskUsage),
		formatBytes(info.NetBytesSent),
		formatBytes(info.NetBytesRecv),
	))

	gpuNpuBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**GPU:** %s\n**NPU:** %s",
		info.GPUInfo,
		info.NPUInfo,
	))

	return discord.NewContainer(
		title,
		discord.NewSmallSeparator(),
		osBlock,
		discord.NewSmallSeparator(),
		cpuBlock,
		memBlock,
		discord.NewSmallSeparator(),
		diskNetBlock,
		gpuNpuBlock,
	)
}
