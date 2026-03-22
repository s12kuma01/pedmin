package panel

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (p *Panel) HandleModal(e *events.ModalSubmitInteractionCreate) {
	if !p.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ui.ErrorMessage("このコマンドを使用する権限がありません。"))
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
		_ = e.CreateMessage(ui.ErrorMessage("コマンドを入力してください。"))
		return
	}

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.HTTPClientTimeout)
	defer cancel()

	if err := p.SendConsoleCommand(ctx, identifier, command); err != nil {
		p.logger.Error("failed to send command", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildConsoleError(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildConsoleResult(command))
}
