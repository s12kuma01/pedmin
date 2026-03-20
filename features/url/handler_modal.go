package url

import (
	"context"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/disgoorg/disgo/events"
)

func (u *URL) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, action, _ := strings.Cut(customID, ":")

	switch action {
	case "shorten_modal":
		u.handleShortenModal(e)
	case "check_modal":
		u.handleCheckModal(e)
	}
}

func (u *URL) handleShortenModal(e *events.ModalSubmitInteractionCreate) {
	rawURL := strings.TrimSpace(e.Data.Text(ModuleID + ":url"))

	if !isValidURL(rawURL) {
		_ = e.DeferUpdateMessage()
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			BuildErrorPanel("有効なURL（http:// または https://）を入力してください。"))
		return
	}

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shortURL, err := u.client.ShortenURL(ctx, rawURL)
	if err != nil {
		u.logger.Error("failed to shorten URL", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			BuildErrorPanel("URL短縮に失敗しました: "+err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		BuildShortenResult(rawURL, shortURL))
}

func (u *URL) handleCheckModal(e *events.ModalSubmitInteractionCreate) {
	rawURL := strings.TrimSpace(e.Data.Text(ModuleID + ":url"))

	if !isValidURL(rawURL) {
		_ = e.DeferUpdateMessage()
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			BuildErrorPanel("有効なURL（http:// または https://）を入力してください。"))
		return
	}

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := u.client.ScanURL(ctx, rawURL)
	if err != nil {
		u.logger.Error("failed to scan URL", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			BuildErrorPanel("URLスキャンに失敗しました: "+err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		BuildCheckResult(rawURL, result))
}

func isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
