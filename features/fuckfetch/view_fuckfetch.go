package fuckfetch

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildFuckfetchUI(info *SystemInfo) discord.ContainerComponent {
	// 1. Title
	title := discord.NewTextDisplay("### 🖥️ fuckfetch")

	// 2. OS block
	osBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**Hostname:** %s\n**OS:** %s (%s)\n**Kernel:** %s\n**Uptime:** %s",
		info.Hostname,
		info.OS,
		info.Platform,
		info.KernelVersion,
		formatUptime(info.Uptime),
	))

	// 3. CPU block
	cpuBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**CPU:** %s (%dC/%dT)\n**Usage:** %s",
		info.CPUModel,
		info.CPUCores,
		info.CPUThreads,
		buildBar(info.CPUUsage),
	))

	// 4. Memory block
	memBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**RAM:** %s / %s %s\n**Swap:** %s / %s %s",
		formatBytes(info.MemUsed),
		formatBytes(info.MemTotal),
		buildBar(info.MemUsage),
		formatBytes(info.SwapUsed),
		formatBytes(info.SwapTotal),
		buildBar(info.SwapUsage),
	))

	// 5. Disk + Network block
	diskNetBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**Disk (/):** %s / %s %s\n**Net ↑:** %s  **Net ↓:** %s",
		formatBytes(info.DiskUsed),
		formatBytes(info.DiskTotal),
		buildBar(info.DiskUsage),
		formatBytes(info.NetBytesSent),
		formatBytes(info.NetBytesRecv),
	))

	// 6. GPU / NPU block
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
