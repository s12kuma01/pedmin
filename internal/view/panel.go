// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

// PanelServerList builds the server list panel.
func PanelServerList(servers []model.Server) discord.MessageCreate {
	if len(servers) == 0 {
		return discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplay("### 🎮 Game Servers"),
				discord.NewSmallSeparator(),
				discord.NewTextDisplay("サーバーが見つかりませんでした。"),
			),
		)
	}

	options := make([]discord.StringSelectMenuOption, 0, len(servers))
	for _, s := range servers {
		emoji := PanelStatusEmoji(s.Status)
		if s.IsSuspended {
			emoji = "⏸️"
		}
		options = append(options, discord.NewStringSelectMenuOption(s.Name, s.Identifier).
			WithDescription(fmt.Sprintf("%s %s", emoji, s.Status)))
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("### 🎮 Game Servers"),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.PanelModuleID+":select", "サーバーを選択...", options...),
			),
			discord.NewActionRow(
				discord.NewSecondaryButton("🔄 更新", model.PanelModuleID+":refresh_list"),
			),
		),
	)
}

// PanelServerDetail builds the server detail panel.
func PanelServerDetail(server model.Server, res *model.Resources) discord.MessageUpdate {
	state := res.CurrentState
	emoji := PanelStatusEmoji(state)

	// Header
	header := fmt.Sprintf("### %s\n%s %s", server.Name, emoji, state)
	if state == "running" && res.Uptime > 0 {
		header += fmt.Sprintf("  |  Uptime: %s", ui.FormatUptime(time.Duration(res.Uptime)*time.Millisecond))
	}

	// Resource bars
	barLen := 10

	cpuPercent := res.CPUAbsolute
	cpuLimit := server.Limits.CPU
	cpuBar := fmt.Sprintf("**CPU:**  %s %.1f%% / %d%%", ui.BuildBar(cpuPercent/float64(cpuLimit)*100, barLen, false), cpuPercent, cpuLimit)

	memUsed := ui.FormatBytes(uint64(res.MemoryBytes))
	memLimit := PanelFormatMBToHuman(server.Limits.Memory)
	memPercent := 0.0
	if server.Limits.Memory > 0 {
		memPercent = float64(res.MemoryBytes) / (float64(server.Limits.Memory) * 1024 * 1024) * 100
	}
	memBar := fmt.Sprintf("**RAM:**  %s %s / %s", ui.BuildBar(memPercent, barLen, false), memUsed, memLimit)

	diskUsed := ui.FormatBytes(uint64(res.DiskBytes))
	diskLimit := PanelFormatMBToHuman(server.Limits.Disk)
	diskPercent := 0.0
	if server.Limits.Disk > 0 {
		diskPercent = float64(res.DiskBytes) / (float64(server.Limits.Disk) * 1024 * 1024) * 100
	}
	diskBar := fmt.Sprintf("**Disk:** %s %s / %s", ui.BuildBar(diskPercent, barLen, false), diskUsed, diskLimit)

	netLine := fmt.Sprintf("**Net ↑:** %s  **Net ↓:** %s", ui.FormatBytes(uint64(res.NetworkTxBytes)), ui.FormatBytes(uint64(res.NetworkRxBytes)))

	resourceBlock := discord.NewTextDisplay(fmt.Sprintf("%s\n%s\n%s\n%s", cpuBar, memBar, diskBar, netLine))

	// Power buttons
	id := server.Identifier
	isRunning := state == "running"
	isOffline := state == "offline"

	startBtn := discord.NewSuccessButton("▶ 起動", model.PanelModuleID+":power_start:"+id)
	if isRunning {
		startBtn = startBtn.AsDisabled()
	}
	restartBtn := discord.NewPrimaryButton("🔄 再起動", model.PanelModuleID+":power_restart:"+id)
	if isOffline {
		restartBtn = restartBtn.AsDisabled()
	}
	stopBtn := discord.NewDangerButton("⏹ 停止", model.PanelModuleID+":power_stop:"+id)
	if isOffline {
		stopBtn = stopBtn.AsDisabled()
	}

	consoleBtn := discord.NewSecondaryButton("💻 コンソール", model.PanelModuleID+":console:"+id)
	if !isRunning {
		consoleBtn = consoleBtn.AsDisabled()
	}

	backBtn := discord.NewSecondaryButton("← 戻る", model.PanelModuleID+":back")
	refreshBtn := discord.NewSecondaryButton("🔄 更新", model.PanelModuleID+":refresh:"+id)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(header),
			discord.NewSmallSeparator(),
			resourceBlock,
			discord.NewSmallSeparator(),
			discord.NewActionRow(startBtn, restartBtn, stopBtn),
			discord.NewActionRow(consoleBtn, backBtn, refreshBtn),
		),
	})
}
