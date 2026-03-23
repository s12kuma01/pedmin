// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/view"
	"github.com/s12kuma01/pedmin/pkg/deepl"
)

// TranslatorService handles message translation logic.
type TranslatorService struct {
	restClient rest.Rest
	deeplClient *deepl.TranslateClient
	logger      *slog.Logger
}

// NewTranslatorService creates a new TranslatorService.
func NewTranslatorService(restClient rest.Rest, deeplClient *deepl.TranslateClient, logger *slog.Logger) *TranslatorService {
	return &TranslatorService{
		restClient:  restClient,
		deeplClient: deeplClient,
		logger:      logger,
	}
}

// IsAvailable reports whether the DeepL client is configured.
func (s *TranslatorService) IsAvailable() bool {
	return s.deeplClient.IsAvailable()
}

// ProcessTranslation fetches a message, translates it, and sends the result.
func (s *TranslatorService) ProcessTranslation(ctx context.Context, channelID, messageID snowflake.ID, targetLang string) {
	msg, err := s.restClient.GetMessage(channelID, messageID)
	if err != nil {
		s.logger.Warn("failed to fetch message for translation",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
		return
	}

	if msg.Author.Bot || msg.Content == "" {
		return
	}

	result, err := s.deeplClient.Translate(ctx, msg.Content, targetLang)
	if err != nil {
		s.logger.Warn("failed to translate message",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
		return
	}

	embed := view.BuildTranslationEmbed(result.TranslatedText, result.DetectedLanguage, targetLang, msg.Author.ID, messageID)
	if _, err := s.restClient.CreateMessage(channelID, embed); err != nil {
		s.logger.Warn("failed to send translation",
			slog.String("message_id", messageID.String()),
			slog.Any("error", err),
		)
	}
}
