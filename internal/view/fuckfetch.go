// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// BuildFuckfetchOutput builds the neofetch-style system info display.
func BuildFuckfetchOutput(info *model.SystemInfo) discord.ContainerComponent {
	const barLen = 20

	systemBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**OS** %s (%s)\n**Kernel** %s\n**Uptime** %s",
		info.OS, info.Platform,
		info.KernelVersion,
		ui.FormatUptime(info.Uptime),
	))

	cpuBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**CPU** %s (%dC/%dT)\n%s %.1f%%\n**GPU** %s",
		info.CPUModel, info.CPUCores, info.CPUThreads,
		ui.BuildBar(info.CPUUsage, barLen, false), info.CPUUsage,
		info.GPUInfo,
	))

	memBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**RAM** %s / %s\n%s %.1f%%\n**Disk** %s / %s\n%s %.1f%%",
		ui.FormatBytes(info.MemUsed), ui.FormatBytes(info.MemTotal),
		ui.BuildBar(info.MemUsage, barLen, false), info.MemUsage,
		ui.FormatBytes(info.DiskUsed), ui.FormatBytes(info.DiskTotal),
		ui.BuildBar(info.DiskUsage, barLen, false), info.DiskUsage,
	))

	netBlock := discord.NewTextDisplay(fmt.Sprintf(
		"**Network** ↑ %s ↓ %s\n**NPU** %s",
		ui.FormatBytes(info.NetBytesSent), ui.FormatBytes(info.NetBytesRecv),
		info.NPUInfo,
	))

	return discord.NewContainer(
		discord.NewTextDisplay("### 🖥️ fuckfetch"),
		discord.NewLargeSeparator(),
		systemBlock,
		discord.NewSmallSeparator(),
		cpuBlock,
		discord.NewSmallSeparator(),
		memBlock,
		discord.NewSmallSeparator(),
		netBlock,
	)
}
