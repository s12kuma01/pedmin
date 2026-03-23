// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/s12kuma01/pedmin/internal/model"
)

// TikTokClient fetches video data from the tikwm.com proxy API.
type TikTokClient struct {
	http *http.Client
}

func NewTikTokClient(timeout time.Duration) *TikTokClient {
	return &TikTokClient{
		http: &http.Client{Timeout: timeout},
	}
}

type tikwmResponse struct {
	Code int       `json:"code"`
	Data tikwmData `json:"data"`
}

type tikwmData struct {
	Title        string      `json:"title"`
	Play         string      `json:"play"`
	PlayCount    int         `json:"play_count"`
	DiggCount    int         `json:"digg_count"`
	CommentCount int         `json:"comment_count"`
	ShareCount   int         `json:"share_count"`
	CreateTime   int64       `json:"create_time"`
	OriginCover  string      `json:"origin_cover"`
	Author       tikwmAuthor `json:"author"`
}

type tikwmAuthor struct {
	UniqueID string `json:"unique_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (c *TikTokClient) GetVideo(ctx context.Context, username, videoID string) (*model.TikTokVideo, error) {
	tiktokURL := fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", username, videoID)
	endpoint := "https://www.tikwm.com/api/?url=" + url.QueryEscape(tiktokURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tikwm API returned %d: %s", resp.StatusCode, string(body))
	}

	var tkResp tikwmResponse
	if err := json.Unmarshal(body, &tkResp); err != nil {
		return nil, fmt.Errorf("invalid tikwm response: %w", err)
	}

	if tkResp.Code != 0 {
		return nil, fmt.Errorf("tikwm API error code %d", tkResp.Code)
	}

	data := tkResp.Data
	return &model.TikTokVideo{
		Title: data.Title,
		Author: model.TikTokAuthor{
			UniqueID: data.Author.UniqueID,
			Nickname: data.Author.Nickname,
			Avatar:   data.Author.Avatar,
		},
		CoverURL:     data.OriginCover,
		VideoURL:     data.Play,
		PlayCount:    data.PlayCount,
		LikeCount:    data.DiggCount,
		CommentCount: data.CommentCount,
		ShareCount:   data.ShareCount,
		CreatedAt:    time.Unix(data.CreateTime, 0),
	}, nil
}
