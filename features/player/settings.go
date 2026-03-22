package player

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type PlayerSettings struct {
	DefaultVolume *int `json:"default_volume"` // nil = use global default
}

func LoadSettings(gs store.GuildStore, guildID snowflake.ID) (*PlayerSettings, error) {
	return store.LoadModuleSettings(gs, guildID, ModuleID, func() *PlayerSettings {
		return &PlayerSettings{}
	})
}

func SaveSettings(gs store.GuildStore, guildID snowflake.ID, settings *PlayerSettings) error {
	return store.SaveModuleSettings(gs, guildID, ModuleID, settings)
}
