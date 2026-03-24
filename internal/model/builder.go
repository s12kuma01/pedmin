// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// ComponentPanel is a saved panel of user-created components.
type ComponentPanel struct {
	ID         int64
	GuildID    snowflake.ID
	Name       string
	Components []PanelComponent
	CreatedAt  time.Time
}

// PanelComponent is a single component in a panel.
type PanelComponent struct {
	Type         PanelComponentType `json:"type"`
	Content      string             `json:"content,omitempty"`
	Texts        []string           `json:"texts,omitempty"`
	ThumbnailURL string             `json:"thumbnail_url,omitempty"`
	Spacing      string             `json:"spacing,omitempty"`
	Divider      bool               `json:"divider,omitempty"`
	Items        []PanelMediaItem   `json:"items,omitempty"`
	Buttons      []PanelLinkButton  `json:"buttons,omitempty"`
}

// PanelComponentType discriminates the component kind.
type PanelComponentType string

const (
	PanelComponentText      PanelComponentType = "text"
	PanelComponentSection   PanelComponentType = "section"
	PanelComponentSeparator PanelComponentType = "separator"
	PanelComponentMedia     PanelComponentType = "media"
	PanelComponentLinks     PanelComponentType = "links"
)

// PanelMediaItem holds a media gallery entry.
type PanelMediaItem struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// PanelLinkButton holds a link button entry.
type PanelLinkButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Emoji string `json:"emoji,omitempty"`
}
