package panel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized (invalid API key)")
	ErrNotFound     = errors.New("server not found")
	ErrRateLimited  = errors.New("rate limited")
)

type PelicanClient struct {
	http    *http.Client
	baseURL string
	apiKey  string
}

func NewPelicanClient(baseURL, apiKey string) *PelicanClient {
	return &PelicanClient{
		http:    &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}

type Server struct {
	Identifier  string
	Name        string
	Description string
	Status      string // "running", "starting", "stopping", "offline"
	Node        string
	Limits      ServerLimits
	IsSuspended bool
}

type ServerLimits struct {
	Memory int // MB
	Disk   int // MB
	CPU    int // percent (e.g. 400 = 4 cores)
}

type Resources struct {
	CurrentState   string
	MemoryBytes    int64
	CPUAbsolute    float64
	DiskBytes      int64
	NetworkRxBytes int64
	NetworkTxBytes int64
	Uptime         int64 // milliseconds
}

func (c *PelicanClient) do(ctx context.Context, method, path string, body string) (*http.Response, error) {
	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	var req *http.Request
	var err error
	if reqBody != nil {
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "Application/vnd.pterodactyl.v1+json")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		switch resp.StatusCode {
		case 401, 403:
			return nil, fmt.Errorf("unauthorized (status %d): %s", resp.StatusCode, string(body))
		case 404:
			return nil, ErrNotFound
		case 429:
			return nil, ErrRateLimited
		default:
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
	}

	return resp, nil
}

// ListServers returns the first page of servers.
func (c *PelicanClient) ListServers(ctx context.Context) ([]Server, error) {
	resp, err := c.do(ctx, http.MethodGet, "/api/client", "")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Data []struct {
			Attributes struct {
				Identifier  string `json:"identifier"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Status      *string `json:"status"`
				Node        string `json:"node"`
				IsSuspended bool   `json:"is_suspended"`
				Limits      struct {
					Memory int `json:"memory"`
					Disk   int `json:"disk"`
					CPU    int `json:"cpu"`
				} `json:"limits"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode server list: %w", err)
	}

	servers := make([]Server, len(result.Data))
	for i, d := range result.Data {
		status := "offline"
		if d.Attributes.Status != nil {
			status = *d.Attributes.Status
		}
		servers[i] = Server{
			Identifier:  d.Attributes.Identifier,
			Name:        d.Attributes.Name,
			Description: d.Attributes.Description,
			Status:      status,
			Node:        d.Attributes.Node,
			IsSuspended: d.Attributes.IsSuspended,
			Limits: ServerLimits{
				Memory: d.Attributes.Limits.Memory,
				Disk:   d.Attributes.Limits.Disk,
				CPU:    d.Attributes.Limits.CPU,
			},
		}
	}

	return servers, nil
}

// GetResources returns current resource usage for a server.
func (c *PelicanClient) GetResources(ctx context.Context, identifier string) (*Resources, error) {
	resp, err := c.do(ctx, http.MethodGet, "/api/client/servers/"+identifier+"/resources", "")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Attributes struct {
			CurrentState string `json:"current_state"`
			Resources    struct {
				MemoryBytes      int64   `json:"memory_bytes"`
				CPUAbsolute      float64 `json:"cpu_absolute"`
				DiskBytes        int64   `json:"disk_bytes"`
				NetworkRxBytes   int64   `json:"network_rx_bytes"`
				NetworkTxBytes   int64   `json:"network_tx_bytes"`
				Uptime           int64   `json:"uptime"`
			} `json:"resources"`
		} `json:"attributes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode resources: %w", err)
	}

	return &Resources{
		CurrentState:   result.Attributes.CurrentState,
		MemoryBytes:    result.Attributes.Resources.MemoryBytes,
		CPUAbsolute:    result.Attributes.Resources.CPUAbsolute,
		DiskBytes:      result.Attributes.Resources.DiskBytes,
		NetworkRxBytes: result.Attributes.Resources.NetworkRxBytes,
		NetworkTxBytes: result.Attributes.Resources.NetworkTxBytes,
		Uptime:         result.Attributes.Resources.Uptime,
	}, nil
}

// SendPowerAction sends a power signal (start, stop, restart, kill).
func (c *PelicanClient) SendPowerAction(ctx context.Context, identifier, signal string) error {
	body := fmt.Sprintf(`{"signal":"%s"}`, signal)
	resp, err := c.do(ctx, http.MethodPost, "/api/client/servers/"+identifier+"/power", body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// SendCommand sends a console command to a server.
func (c *PelicanClient) SendCommand(ctx context.Context, identifier, command string) error {
	body := fmt.Sprintf(`{"command":%s}`, jsonString(command))
	resp, err := c.do(ctx, http.MethodPost, "/api/client/servers/"+identifier+"/command", body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
