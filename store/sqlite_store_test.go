package store

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create SQLiteStore: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSQLiteStore_GetDefault(t *testing.T) {
	s := newTestStore(t)

	guildID := snowflake.ID(123456789)
	settings, err := s.Get(guildID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if settings.GuildID != guildID {
		t.Errorf("expected GuildID %d, got %d", guildID, settings.GuildID)
	}
	if len(settings.EnabledModules) != 0 {
		t.Errorf("expected empty EnabledModules, got %v", settings.EnabledModules)
	}
	if len(settings.ModuleSettings) != 0 {
		t.Errorf("expected empty ModuleSettings, got %v", settings.ModuleSettings)
	}
}

func TestSQLiteStore_SetModuleEnabled(t *testing.T) {
	s := newTestStore(t)

	guildID := snowflake.ID(123456789)

	// Initially disabled
	enabled, err := s.IsModuleEnabled(guildID, "player")
	if err != nil {
		t.Fatalf("IsModuleEnabled failed: %v", err)
	}
	if enabled {
		t.Error("expected module to be disabled by default")
	}

	// Enable
	if err := s.SetModuleEnabled(guildID, "player", true); err != nil {
		t.Fatalf("SetModuleEnabled failed: %v", err)
	}

	enabled, err = s.IsModuleEnabled(guildID, "player")
	if err != nil {
		t.Fatalf("IsModuleEnabled failed: %v", err)
	}
	if !enabled {
		t.Error("expected module to be enabled")
	}

	// Disable
	if err := s.SetModuleEnabled(guildID, "player", false); err != nil {
		t.Fatalf("SetModuleEnabled failed: %v", err)
	}

	enabled, err = s.IsModuleEnabled(guildID, "player")
	if err != nil {
		t.Fatalf("IsModuleEnabled failed: %v", err)
	}
	if enabled {
		t.Error("expected module to be disabled")
	}
}

func TestSQLiteStore_SaveAndGet(t *testing.T) {
	s := newTestStore(t)

	guildID := snowflake.ID(987654321)
	settings := &GuildSettings{
		GuildID: guildID,
		EnabledModules: map[string]bool{
			"player":   true,
			"settings": true,
		},
		ModuleSettings: map[string]any{
			"player": map[string]any{
				"default_volume": float64(50),
			},
		},
	}

	if err := s.Save(settings); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := s.Get(guildID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if loaded.GuildID != guildID {
		t.Errorf("expected GuildID %d, got %d", guildID, loaded.GuildID)
	}

	if !loaded.EnabledModules["player"] {
		t.Error("expected player module to be enabled")
	}
	if !loaded.EnabledModules["settings"] {
		t.Error("expected settings module to be enabled")
	}

	playerSettings, ok := loaded.ModuleSettings["player"]
	if !ok {
		t.Fatal("expected player module settings")
	}
	ps, ok := playerSettings.(map[string]any)
	if !ok {
		t.Fatal("expected player settings to be map[string]any")
	}
	if ps["default_volume"] != float64(50) {
		t.Errorf("expected default_volume 50, got %v", ps["default_volume"])
	}
}

func TestSQLiteStore_MultipleGuilds(t *testing.T) {
	s := newTestStore(t)

	guild1 := snowflake.ID(111111111)
	guild2 := snowflake.ID(222222222)

	if err := s.SetModuleEnabled(guild1, "player", true); err != nil {
		t.Fatalf("SetModuleEnabled failed: %v", err)
	}
	if err := s.SetModuleEnabled(guild2, "player", false); err != nil {
		t.Fatalf("SetModuleEnabled failed: %v", err)
	}

	e1, _ := s.IsModuleEnabled(guild1, "player")
	e2, _ := s.IsModuleEnabled(guild2, "player")

	if !e1 {
		t.Error("expected player enabled for guild1")
	}
	if e2 {
		t.Error("expected player disabled for guild2")
	}
}

func TestSQLiteStore_ConcurrentAccess(t *testing.T) {
	s := newTestStore(t)

	guildID := snowflake.ID(333333333)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			enabled := i%2 == 0
			if err := s.SetModuleEnabled(guildID, "player", enabled); err != nil {
				t.Errorf("concurrent SetModuleEnabled failed: %v", err)
			}
			if _, err := s.IsModuleEnabled(guildID, "player"); err != nil {
				t.Errorf("concurrent IsModuleEnabled failed: %v", err)
			}
		}(i)
	}
	wg.Wait()
}

