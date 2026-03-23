// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"time"

	"github.com/s12kuma01/pedmin/internal/client"
	"github.com/s12kuma01/pedmin/internal/model"
)

// PanelService handles game server management operations via the Pelican API.
type PanelService struct {
	pelican *client.PelicanClient
}

// NewPanelService creates a new PanelService.
func NewPanelService(pelican *client.PelicanClient) *PanelService {
	return &PanelService{pelican: pelican}
}

// ListServersWithStatus returns all servers with their current resource state.
func (s *PanelService) ListServersWithStatus(ctx context.Context) ([]model.Server, error) {
	servers, err := s.pelican.ListServers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range servers {
		res, err := s.pelican.GetResources(ctx, servers[i].Identifier)
		if err == nil {
			servers[i].Status = res.CurrentState
		}
	}

	return servers, nil
}

// FindServer looks up a single server by identifier.
func (s *PanelService) FindServer(ctx context.Context, identifier string) (*model.Server, error) {
	servers, err := s.pelican.ListServers(ctx)
	if err != nil {
		return nil, err
	}
	for _, srv := range servers {
		if srv.Identifier == identifier {
			return &srv, nil
		}
	}
	return nil, model.ErrPanelNotFound
}

// GetServerDetail returns a server and its current resources.
func (s *PanelService) GetServerDetail(ctx context.Context, identifier string) (*model.Server, *model.Resources, error) {
	server, err := s.FindServer(ctx, identifier)
	if err != nil {
		return nil, nil, err
	}

	res, err := s.pelican.GetResources(ctx, identifier)
	if err != nil {
		return nil, nil, err
	}

	return server, res, nil
}

// PowerAction sends a power signal and waits briefly for state transition, then returns updated detail.
func (s *PanelService) PowerAction(ctx context.Context, identifier, signal string) (*model.Server, *model.Resources, error) {
	if err := s.pelican.SendPowerAction(ctx, identifier, signal); err != nil {
		return nil, nil, err
	}

	// Wait for state transition
	time.Sleep(2 * time.Second)

	return s.GetServerDetail(ctx, identifier)
}

// SendConsoleCommand sends a console command to a server.
func (s *PanelService) SendConsoleCommand(ctx context.Context, identifier, command string) error {
	return s.pelican.SendCommand(ctx, identifier, command)
}
