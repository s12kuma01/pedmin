// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// ComponentTypeName returns the display name for a component type.
func ComponentTypeName(t model.PanelComponentType) string {
	switch t {
	case model.PanelComponentText:
		return "📝 テキスト"
	case model.PanelComponentSection:
		return "📋 セクション"
	case model.PanelComponentSeparator:
		return "─ セパレータ"
	case model.PanelComponentMedia:
		return "🖼 画像"
	case model.PanelComponentLinks:
		return "🔗 リンクボタン"
	default:
		return string(t)
	}
}

// ComponentSummary returns a short summary for a component.
func ComponentSummary(comp model.PanelComponent) string {
	switch comp.Type {
	case model.PanelComponentText:
		return truncateBuilderStr(comp.Content, 40)
	case model.PanelComponentSection:
		if len(comp.Texts) > 0 {
			return truncateBuilderStr(comp.Texts[0], 40)
		}
		return "(空)"
	case model.PanelComponentSeparator:
		if comp.Spacing == "large" {
			return "大"
		}
		return "小"
	case model.PanelComponentMedia:
		return fmt.Sprintf("%d枚", len(comp.Items))
	case model.PanelComponentLinks:
		return fmt.Sprintf("%d個", len(comp.Buttons))
	default:
		return ""
	}
}

// BuilderErrorContainer builds an error container.
func BuilderErrorContainer(text string) discord.ContainerComponent {
	return discord.NewContainer(discord.NewTextDisplay(text))
}

func truncateBuilderStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
