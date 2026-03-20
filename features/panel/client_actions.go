package panel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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
