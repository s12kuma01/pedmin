// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "time"

// TikTokVideo represents a processed TikTok video for embed display.
type TikTokVideo struct {
	Title        string
	Author       TikTokAuthor
	CoverURL     string
	VideoURL     string
	PlayCount    int
	LikeCount    int
	CommentCount int
	ShareCount   int
	CreatedAt    time.Time
}

// TikTokAuthor holds TikTok author information.
type TikTokAuthor struct {
	UniqueID string
	Nickname string
	Avatar   string
}
