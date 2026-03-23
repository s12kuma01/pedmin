// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *URLHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	switch action {
	case "shorten":
		_ = e.Modal(discord.ModalCreate{
			CustomID: model.URLModuleID + ":shorten_modal",
			Title:    "URL短縮",
			Components: []discord.LayoutComponent{
				discord.NewLabel("URL",
					discord.NewShortTextInput(model.URLModuleID+":url").
						WithRequired(true).
						WithPlaceholder("https://example.com/long/path"),
				),
			},
		})

	case "check":
		_ = e.Modal(discord.ModalCreate{
			CustomID: model.URLModuleID + ":check_modal",
			Title:    "URLチェッカー",
			Components: []discord.LayoutComponent{
				discord.NewLabel("URL",
					discord.NewShortTextInput(model.URLModuleID+":url").
						WithRequired(true).
						WithPlaceholder("https://example.com"),
				),
			},
		})

	case "back":
		hasXGD := h.cfg.XGDAPIKey != ""
		hasVT := h.cfg.VTAPIKey != ""

		msg := view.BuildURLMainPanel(hasXGD, hasVT)
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2(msg.Components))
	}
}
