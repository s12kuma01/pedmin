// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// PanelConsoleModal builds the console command input modal.
func PanelConsoleModal(identifier string) discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: model.PanelModuleID + ":console_modal:" + identifier,
		Title:    "コンソールコマンド",
		Components: []discord.LayoutComponent{
			discord.NewLabel("コマンド",
				discord.NewShortTextInput(model.PanelModuleID+":cmd").
					WithRequired(true).
					WithPlaceholder("say hello"),
			),
		},
	}
}

// PanelConsoleResult builds the console command success message.
func PanelConsoleResult(command string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("コマンドを送信しました: `%s`", command)),
		),
	})
}

// PanelConsoleError builds the console command error message.
func PanelConsoleError(errMsg string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("コマンド送信に失敗しました:\n%s", errMsg)),
		),
	})
}

// PanelErrorPanel builds a generic error panel with a back button.
func PanelErrorPanel(errMsg string) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### ❌ エラー\n%s", errMsg)),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", model.PanelModuleID+":back"),
			),
		),
	})
}
