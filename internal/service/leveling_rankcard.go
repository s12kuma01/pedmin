// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/pkg/rankcard"
)

// GenerateRankCard generates a rank card PNG image for a user.
func (s *LevelingService) GenerateRankCard(ctx context.Context, guildID, userID snowflake.ID) ([]byte, error) {
	ux, err := s.store.GetUserXP(guildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user xp: %w", err)
	}

	rank, err := s.store.GetUserRank(guildID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rank: %w", err)
	}

	user, err := s.client.Rest.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	avatarURL := user.EffectiveAvatarURL(discord.WithSize(256), discord.WithFormat(discord.FileFormatPNG))
	avatarData, err := fetchImage(ctx, avatarURL)
	if err != nil {
		s.logger.Error("failed to fetch avatar, using nil", slog.Any("error", err))
		avatarData = nil
	}

	return rankcard.Generate(rankcard.Data{
		Username:  user.EffectiveName(),
		AvatarPNG: avatarData,
		Level:     ux.Level,
		CurrentXP: ux.CurrentXP(),
		NeededXP:  ux.NeededXP(),
		TotalXP:   ux.TotalXP,
		Rank:      rank,
	})
}

func fetchImage(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("avatar fetch returned %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
