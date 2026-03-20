package panel

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildServerList(servers []Server) discord.MessageCreate {
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
		emoji := statusEmoji(s.Status)
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
				discord.NewStringSelectMenu(ModuleID+":select", "サーバーを選択...", options...),
			),
			discord.NewActionRow(
				discord.NewSecondaryButton("🔄 更新", ModuleID+":refresh_list"),
			),
		),
	)
}

func BuildServerDetail(server Server, res *Resources) discord.MessageUpdate {
	state := res.CurrentState
	emoji := statusEmoji(state)

	// Header
	header := fmt.Sprintf("### %s\n%s %s", server.Name, emoji, state)
	if state == "running" && res.Uptime > 0 {
		header += fmt.Sprintf("  |  Uptime: %s", formatUptime(res.Uptime))
	}

	// Resource bars
	barLen := 10

	cpuPercent := res.CPUAbsolute
	cpuLimit := server.Limits.CPU
	cpuBar := fmt.Sprintf("**CPU:**  %s %.1f%% / %d%%", buildBar(cpuPercent/float64(cpuLimit)*100, barLen), cpuPercent, cpuLimit)

	memUsed := formatBytes(res.MemoryBytes)
	memLimit := formatMBToHuman(server.Limits.Memory)
	memPercent := 0.0
	if server.Limits.Memory > 0 {
		memPercent = float64(res.MemoryBytes) / (float64(server.Limits.Memory) * 1024 * 1024) * 100
	}
	memBar := fmt.Sprintf("**RAM:**  %s %s / %s", buildBar(memPercent, barLen), memUsed, memLimit)

	diskUsed := formatBytes(res.DiskBytes)
	diskLimit := formatMBToHuman(server.Limits.Disk)
	diskPercent := 0.0
	if server.Limits.Disk > 0 {
		diskPercent = float64(res.DiskBytes) / (float64(server.Limits.Disk) * 1024 * 1024) * 100
	}
	diskBar := fmt.Sprintf("**Disk:** %s %s / %s", buildBar(diskPercent, barLen), diskUsed, diskLimit)

	netLine := fmt.Sprintf("**Net ↑:** %s  **Net ↓:** %s", formatBytes(res.NetworkTxBytes), formatBytes(res.NetworkRxBytes))

	resourceBlock := discord.NewTextDisplay(fmt.Sprintf("%s\n%s\n%s\n%s", cpuBar, memBar, diskBar, netLine))

	// Power buttons
	id := server.Identifier
	isRunning := state == "running"
	isOffline := state == "offline"

	startBtn := discord.NewSuccessButton("▶ 起動", ModuleID+":power_start:"+id)
	if isRunning {
		startBtn = startBtn.AsDisabled()
	}
	restartBtn := discord.NewPrimaryButton("🔄 再起動", ModuleID+":power_restart:"+id)
	if isOffline {
		restartBtn = restartBtn.AsDisabled()
	}
	stopBtn := discord.NewDangerButton("⏹ 停止", ModuleID+":power_stop:"+id)
	if isOffline {
		stopBtn = stopBtn.AsDisabled()
	}

	consoleBtn := discord.NewSecondaryButton("💻 コンソール", ModuleID+":console:"+id)
	if !isRunning {
		consoleBtn = consoleBtn.AsDisabled()
	}

	backBtn := discord.NewSecondaryButton("← 戻る", ModuleID+":back")
	refreshBtn := discord.NewSecondaryButton("🔄 更新", ModuleID+":refresh:"+id)

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

func BuildErrorPanel(errMsg string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### ❌ エラー\n%s", errMsg)),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", ModuleID+":back"),
			),
		),
	})
}
