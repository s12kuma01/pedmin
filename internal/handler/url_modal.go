package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *URLHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, action, _ := strings.Cut(customID, ":")

	switch action {
	case "shorten_modal":
		h.handleShortenModal(e)
	case "check_modal":
		h.handleCheckModal(e)
	}
}

func (h *URLHandler) handleShortenModal(e *events.ModalSubmitInteractionCreate) {
	rawURL := strings.TrimSpace(e.Data.Text(model.URLModuleID + ":url"))

	if !service.URLIsValid(rawURL) {
		_ = e.DeferUpdateMessage()
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			view.BuildURLErrorPanel("有効なURL（http:// または https://）を入力してください。"))
		return
	}

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.ShortenTimeout)
	defer cancel()

	shortURL, err := h.service.Shorten(ctx, rawURL)
	if err != nil {
		h.logger.Error("failed to shorten URL", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			view.BuildURLErrorPanel("URL短縮に失敗しました: "+err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		view.BuildURLShortenResult(rawURL, shortURL))
}

func (h *URLHandler) handleCheckModal(e *events.ModalSubmitInteractionCreate) {
	rawURL := strings.TrimSpace(e.Data.Text(model.URLModuleID + ":url"))

	if !service.URLIsValid(rawURL) {
		_ = e.DeferUpdateMessage()
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			view.BuildURLErrorPanel("有効なURL（http:// または https://）を入力してください。"))
		return
	}

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.ScanTimeout)
	defer cancel()

	result, err := h.service.Check(ctx, rawURL)
	if err != nil {
		h.logger.Error("failed to scan URL", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			view.BuildURLErrorPanel("URLスキャンに失敗しました: "+err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		view.BuildURLCheckResult(rawURL, result))
}
