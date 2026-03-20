package panel

import (
	"context"
	"time"
)

// ListServersWithStatus returns all servers with their current resource state.
func (p *Panel) ListServersWithStatus(ctx context.Context) ([]Server, error) {
	servers, err := p.pelican.ListServers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range servers {
		res, err := p.pelican.GetResources(ctx, servers[i].Identifier)
		if err == nil {
			servers[i].Status = res.CurrentState
		}
	}

	return servers, nil
}

// FindServer looks up a single server by identifier.
func (p *Panel) FindServer(ctx context.Context, identifier string) (*Server, error) {
	servers, err := p.pelican.ListServers(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		if s.Identifier == identifier {
			return &s, nil
		}
	}
	return nil, ErrNotFound
}

// GetServerDetail returns a server and its current resources.
func (p *Panel) GetServerDetail(ctx context.Context, identifier string) (*Server, *Resources, error) {
	server, err := p.FindServer(ctx, identifier)
	if err != nil {
		return nil, nil, err
	}

	res, err := p.pelican.GetResources(ctx, identifier)
	if err != nil {
		return nil, nil, err
	}

	return server, res, nil
}

// PowerAction sends a power signal and waits briefly for state transition, then returns updated detail.
func (p *Panel) PowerAction(ctx context.Context, identifier, signal string) (*Server, *Resources, error) {
	if err := p.pelican.SendPowerAction(ctx, identifier, signal); err != nil {
		return nil, nil, err
	}

	// Wait for state transition
	time.Sleep(2 * time.Second)

	return p.GetServerDetail(ctx, identifier)
}

// SendConsoleCommand sends a console command to a server.
func (p *Panel) SendConsoleCommand(ctx context.Context, identifier, command string) error {
	return p.pelican.SendCommand(ctx, identifier, command)
}
