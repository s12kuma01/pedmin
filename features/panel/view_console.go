package panel

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildConsoleModal(identifier string) discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: ModuleID + ":console_modal:" + identifier,
		Title:    "コンソールコマンド",
		Components: []discord.LayoutComponent{
			discord.NewLabel("コマンド",
				discord.NewShortTextInput(ModuleID+":cmd").
					WithRequired(true).
					WithPlaceholder("say hello"),
			),
		},
	}
}

func BuildConsoleResult(command string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("コマンドを送信しました: `%s`", command)),
		),
	})
}

func BuildConsoleError(errMsg string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("コマンド送信に失敗しました:\n%s", errMsg)),
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
