// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// RenderComponentPanel converts a panel's components to a Discord container.
func RenderComponentPanel(panel *model.ComponentPanel) discord.ContainerComponent {
	var subs []discord.ContainerSubComponent

	for _, comp := range panel.Components {
		switch comp.Type {
		case model.PanelComponentText:
			subs = append(subs, discord.NewTextDisplay(comp.Content))

		case model.PanelComponentSection:
			var textDisplays []discord.SectionSubComponent
			for _, t := range comp.Texts {
				if t != "" {
					textDisplays = append(textDisplays, discord.NewTextDisplay(t))
				}
			}
			if len(textDisplays) == 0 {
				continue
			}
			section := discord.NewSection(textDisplays...)
			if comp.ThumbnailURL != "" {
				section = section.WithAccessory(discord.NewThumbnail(comp.ThumbnailURL))
			}
			subs = append(subs, section)

		case model.PanelComponentSeparator:
			if comp.Spacing == "large" {
				subs = append(subs, discord.NewLargeSeparator())
			} else {
				subs = append(subs, discord.NewSmallSeparator())
			}

		case model.PanelComponentMedia:
			var items []discord.MediaGalleryItem
			for _, item := range comp.Items {
				mgItem := discord.MediaGalleryItem{
					Media: discord.UnfurledMediaItem{URL: item.URL},
				}
				if item.Description != "" {
					mgItem.Description = item.Description
				}
				items = append(items, mgItem)
			}
			if len(items) > 0 {
				subs = append(subs, discord.NewMediaGallery(items...))
			}

		case model.PanelComponentLinks:
			var buttons []discord.InteractiveComponent
			for _, btn := range comp.Buttons {
				b := discord.NewLinkButton(btn.Label, btn.URL)
				if btn.Emoji != "" {
					b = b.WithEmoji(discord.ComponentEmoji{Name: btn.Emoji})
				}
				buttons = append(buttons, b)
			}
			if len(buttons) > 0 {
				subs = append(subs, discord.NewActionRow(buttons...))
			}
		}
	}

	if len(subs) == 0 {
		subs = append(subs, discord.NewTextDisplay("-# パネルにコンポーネントがありません"))
	}

	return discord.NewContainer(subs...)
}
