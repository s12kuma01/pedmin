package panel

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (p *Panel) HandleModal(e *events.ModalSubmitInteractionCreate) {
	if !p.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ephemeralError("このコマンドを使用する権限がありません。"))
		return
	}

	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, identifier, _ := strings.Cut(rest, ":")

	if action != "console_modal" {
		return
	}

	command := strings.TrimSpace(e.Data.Text(ModuleID + ":cmd"))
	if command == "" {
		_ = e.CreateMessage(ephemeralError("コマンドを入力してください。"))
		return
	}

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.pelican.SendCommand(ctx, identifier, command); err != nil {
		p.logger.Error("failed to send command", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(fmt.Sprintf("コマンド送信に失敗しました:\n%s", err.Error())),
			),
		}))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("コマンドを送信しました: `%s`", command)),
		),
	}))
}
