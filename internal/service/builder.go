// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/repository"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// BuilderService handles component panel business logic.
type BuilderService struct {
	store  repository.GuildStore
	client *disgobot.Client
	logger *slog.Logger
}

// NewBuilderService creates a new BuilderService.
func NewBuilderService(store repository.GuildStore, client *disgobot.Client, logger *slog.Logger) *BuilderService {
	return &BuilderService{
		store:  store,
		client: client,
		logger: logger,
	}
}

// CreatePanel creates a new empty panel.
func (s *BuilderService) CreatePanel(guildID snowflake.ID, name string) (*model.ComponentPanel, error) {
	panel := &model.ComponentPanel{
		GuildID:    guildID,
		Name:       name,
		Components: []model.PanelComponent{},
	}
	if err := s.store.CreatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// GetPanels returns all panels for a guild.
func (s *BuilderService) GetPanels(guildID snowflake.ID) ([]model.ComponentPanel, error) {
	return s.store.GetPanels(guildID)
}

// GetPanel returns a single panel.
func (s *BuilderService) GetPanel(id int64, guildID snowflake.ID) (*model.ComponentPanel, error) {
	return s.store.GetPanel(id, guildID)
}

// DeletePanel deletes a panel.
func (s *BuilderService) DeletePanel(id int64, guildID snowflake.ID) error {
	return s.store.DeletePanel(id, guildID)
}

// CountPanels returns the number of panels for a guild.
func (s *BuilderService) CountPanels(guildID snowflake.ID) (int, error) {
	return s.store.CountPanels(guildID)
}

// RenamePanel renames a panel.
func (s *BuilderService) RenamePanel(id int64, guildID snowflake.ID, newName string) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	panel.Name = newName
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// AddComponent appends a component to a panel.
func (s *BuilderService) AddComponent(id int64, guildID snowflake.ID, comp model.PanelComponent) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if len(panel.Components) >= model.MaxComponentsPerPanel {
		return nil, fmt.Errorf("コンポーネント数が上限(%d)に達しています", model.MaxComponentsPerPanel)
	}
	panel.Components = append(panel.Components, comp)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// RemoveComponent removes a component at the given index.
func (s *BuilderService) RemoveComponent(id int64, guildID snowflake.ID, index int) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(panel.Components) {
		return nil, fmt.Errorf("無効なインデックスです")
	}
	panel.Components = append(panel.Components[:index], panel.Components[index+1:]...)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// MoveComponent moves a component from one index to another.
func (s *BuilderService) MoveComponent(id int64, guildID snowflake.ID, from, to int) (*model.ComponentPanel, error) {
	panel, err := s.store.GetPanel(id, guildID)
	if err != nil {
		return nil, err
	}
	if from < 0 || from >= len(panel.Components) || to < 0 || to >= len(panel.Components) {
		return nil, fmt.Errorf("無効なインデックスです")
	}
	comp := panel.Components[from]
	panel.Components = append(panel.Components[:from], panel.Components[from+1:]...)
	// Insert at 'to' position
	panel.Components = append(panel.Components[:to], append([]model.PanelComponent{comp}, panel.Components[to:]...)...)
	if err := s.store.UpdatePanel(panel); err != nil {
		return nil, err
	}
	return panel, nil
}

// PreviewPanel renders the panel as an ephemeral message.
// Note: uses view.RenderComponentPanel for component rendering.
func (s *BuilderService) PreviewPanel(panel *model.ComponentPanel) discord.MessageCreate {
	container := view.RenderComponentPanel(panel)
	return discord.NewMessageCreateV2(container).WithEphemeral(true)
}

// DeployPanel sends the rendered panel to a channel.
func (s *BuilderService) DeployPanel(panel *model.ComponentPanel, channelID snowflake.ID) error {
	container := view.RenderComponentPanel(panel)
	msg := discord.NewMessageCreateV2(container)
	_, err := s.client.Rest.CreateMessage(channelID, msg)
	return err
}
