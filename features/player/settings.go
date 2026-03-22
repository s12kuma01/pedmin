package player

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type PlayerSettings struct {
	DefaultVolume *int `json:"default_volume"` // nil = use global default
}

func LoadSettings(gs store.GuildStore, guildID snowflake.ID) (*PlayerSettings, error) {
	data, err := gs.GetModuleSettings(guildID, ModuleID)
	if err != nil {
		return nil, err
	}
	var s PlayerSettings
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return &PlayerSettings{}, nil
	}
	return &s, nil
}

func SaveSettings(gs store.GuildStore, guildID snowflake.ID, settings *PlayerSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return gs.SetModuleSettings(guildID, ModuleID, string(data))
}
