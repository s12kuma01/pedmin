package store

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
)

// LoadModuleSettings loads and unmarshals module-specific settings from the store.
// If the data is missing or invalid, defaultFn provides the fallback value.
func LoadModuleSettings[T any](gs GuildStore, guildID snowflake.ID, moduleID string, defaultFn func() *T) (*T, error) {
	data, err := gs.GetModuleSettings(guildID, moduleID)
	if err != nil {
		return nil, err
	}
	var s T
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return defaultFn(), nil
	}
	return &s, nil
}

// SaveModuleSettings marshals and persists module-specific settings to the store.
func SaveModuleSettings[T any](gs GuildStore, guildID snowflake.ID, moduleID string, settings *T) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return gs.SetModuleSettings(guildID, moduleID, string(data))
}
