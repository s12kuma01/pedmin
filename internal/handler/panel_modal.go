package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *PanelHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	if !h.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ui.ErrorMessage("このコマンドを使用する権限がありません。"))
		return
	}

	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, identifier, _ := strings.Cut(rest, ":")

	if action != "console_modal" {
		return
	}

	command := strings.TrimSpace(e.Data.Text(model.PanelModuleID + ":cmd"))
	if command == "" {
		_ = e.CreateMessage(ui.ErrorMessage("コマンドを入力してください。"))
		return
	}

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.HTTPClientTimeout)
	defer cancel()

	if err := h.service.SendConsoleCommand(ctx, identifier, command); err != nil {
		h.logger.Error("failed to send command", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelConsoleError(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelConsoleResult(command))
}
