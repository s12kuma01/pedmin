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
	pad := "        " // 8 spaces for label indent

	body := fmt.Sprintf("```\n"+
		"%-8s%s (%s)\n"+
		"%-8s%s\n"+
		"%-8s%s\n"+
		"%-8s%s (%dC/%dT)\n"+
		"%s%s %5.1f%%\n"+
		"%-8s%s\n"+
		"%-8s%s / %s\n"+
		"%s%s %5.1f%%\n"+
		"%-8s%s / %s\n"+
		"%s%s %5.1f%%\n"+
		"%-8s↑ %s  ↓ %s\n"+
		"%-8s%s\n"+
		"```",
		"OS", info.OS, info.Platform,
		"Kernel", info.KernelVersion,
		"Uptime", ui.FormatUptime(info.Uptime),
		"CPU", info.CPUModel, info.CPUCores, info.CPUThreads,
		pad, ui.BuildBarRaw(info.CPUUsage, barLen), info.CPUUsage,
		"GPU", info.GPUInfo,
		"RAM", ui.FormatBytes(info.MemUsed), ui.FormatBytes(info.MemTotal),
		pad, ui.BuildBarRaw(info.MemUsage, barLen), info.MemUsage,
		"Disk", ui.FormatBytes(info.DiskUsed), ui.FormatBytes(info.DiskTotal),
		pad, ui.BuildBarRaw(info.DiskUsage, barLen), info.DiskUsage,
		"Network", ui.FormatBytes(info.NetBytesSent), ui.FormatBytes(info.NetBytesRecv),
		"NPU", info.NPUInfo,
	)

	return discord.NewContainer(
		discord.NewTextDisplay("### 🖥️ fuckfetch"),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	)
}
