package embedfix

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type EmbedFixSettings struct {
	Platforms map[Platform]bool `json:"platforms"`
}

var AllPlatforms = []struct {
	Key   Platform
	Label string
}{
	{PlatformTwitter, "X / Twitter"},
	{PlatformReddit, "Reddit"},
	{PlatformTikTok, "TikTok"},
}

func defaultSettings() *EmbedFixSettings {
	platforms := make(map[Platform]bool, len(AllPlatforms))
	for _, p := range AllPlatforms {
		platforms[p.Key] = true
	}
	return &EmbedFixSettings{Platforms: platforms}
}

func (s *EmbedFixSettings) IsPlatformEnabled(p Platform) bool {
	enabled, ok := s.Platforms[p]
	if !ok {
		return true // enabled by default
	}
	return enabled
}

func LoadSettings(gs store.GuildStore, guildID snowflake.ID) (*EmbedFixSettings, error) {
	s, err := store.LoadModuleSettings(gs, guildID, ModuleID, defaultSettings)
	if err != nil {
		return nil, err
	}
	if s.Platforms == nil {
		return defaultSettings(), nil
	}
	return s, nil
}

func SaveSettings(gs store.GuildStore, guildID snowflake.ID, settings *EmbedFixSettings) error {
	return store.SaveModuleSettings(gs, guildID, ModuleID, settings)
}
