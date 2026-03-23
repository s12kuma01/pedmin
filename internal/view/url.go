package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// BuildURLMainPanel builds the main URL tools panel.
func BuildURLMainPanel(hasXGD, hasVT bool) discord.MessageCreate {
	shortenBtn := discord.NewPrimaryButton("🔗 URL短縮", model.URLModuleID+":shorten")
	if !hasXGD {
		shortenBtn = shortenBtn.AsDisabled()
	}

	checkBtn := discord.NewPrimaryButton("🛡️ URLチェッカー", model.URLModuleID+":check")
	if !hasVT {
		checkBtn = checkBtn.AsDisabled()
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("### 🔗 URLツール"),
			discord.NewSmallSeparator(),
			discord.NewActionRow(shortenBtn, checkBtn),
		),
	).WithEphemeral(true)
}

// BuildURLShortenResult builds the URL shorten result panel.
func BuildURLShortenResult(originalURL, shortURL string) discord.MessageUpdate {
	text := fmt.Sprintf("### 🔗 URL短縮\n**元URL:** %s\n**短縮URL:** %s", originalURL, shortURL)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(text),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", model.URLModuleID+":back"),
			),
		),
	})
}

// BuildURLCheckResult builds the URL check result panel.
func BuildURLCheckResult(rawURL string, result *model.VTResult) discord.MessageUpdate {
	var verdict string
	switch {
	case result.Malicious > 0:
		verdict = "🚨 危険"
	case result.Suspicious > 0:
		verdict = "⚠️ 注意"
	default:
		verdict = "✅ 安全"
	}

	text := fmt.Sprintf("### 🛡️ URLチェッカー\n**URL:** %s\n\n%s\n🟢 Harmless: %d\n🔴 Malicious: %d\n🟡 Suspicious: %d\n⚪ Undetected: %d",
		rawURL, verdict,
		result.Harmless, result.Malicious, result.Suspicious, result.Undetected)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(text),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", model.URLModuleID+":back"),
			),
		),
	})
}

// BuildURLErrorPanel builds the URL error panel.
func BuildURLErrorPanel(errMsg string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### ❌ エラー\n%s", errMsg)),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", model.URLModuleID+":back"),
			),
		),
	})
}
