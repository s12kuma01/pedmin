package fuckfetch

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/ui"
)

func BuildFuckfetchUI(info *SystemInfo) discord.ContainerComponent {
	title := discord.NewTextDisplay("### 🖥️ fuckfetch")

	osBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**OS:** %s (%s)\n**Kernel:** %s\n**Uptime:** %s",
		info.OS,
		info.Platform,
		info.KernelVersion,
		ui.FormatUptime(info.Uptime),
	))

	cpuBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**CPU:** %s (%dC/%dT)\n**Usage:** %s",
		info.CPUModel,
		info.CPUCores,
		info.CPUThreads,
		ui.BuildBar(info.CPUUsage, 20, true),
	))

	memBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**RAM:** %s / %s %s\n**Swap:** %s / %s %s",
		ui.FormatBytes(info.MemUsed),
		ui.FormatBytes(info.MemTotal),
		ui.BuildBar(info.MemUsage, 20, true),
		ui.FormatBytes(info.SwapUsed),
		ui.FormatBytes(info.SwapTotal),
		ui.BuildBar(info.SwapUsage, 20, true),
	))

	diskNetBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**Disk (/):** %s / %s %s\n**Net ↑:** %s  **Net ↓:** %s",
		ui.FormatBytes(info.DiskUsed),
		ui.FormatBytes(info.DiskTotal),
		ui.BuildBar(info.DiskUsage, 20, true),
		ui.FormatBytes(info.NetBytesSent),
		ui.FormatBytes(info.NetBytesRecv),
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
