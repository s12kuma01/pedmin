# Module Development Guide

## Architecture Overview

Each module is implemented across multiple layer packages. Handler structs implement the `module.Module` interface and
are registered in `cmd/pedmin/main.go`.

```
internal/
├── handler/myfeature.go       # Handler struct, module.Module implementation
├── service/myfeature.go       # Business logic (optional)
├── view/myfeature.go          # UI builders (optional)
└── model/myfeature.go         # Domain types (optional)
```

Handlers use `package handler`, services use `package service`, views use `package view`, models use `package model`.

## Creating a New Module

### Step 1: Add ModuleID to `internal/model/constants.go`

```go
const MyFeatureModuleID = "myfeature"
```

### Step 2: Create `internal/handler/myfeature.go`

```go
package handler

import (
    "log/slog"

    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
    "github.com/disgoorg/snowflake/v2"
    "github.com/s12kuma01/pedmin/internal/model"
    "github.com/s12kuma01/pedmin/internal/module"
    "github.com/s12kuma01/pedmin/internal/view"
)

type MyFeatureHandler struct {
    logger *slog.Logger
}

func NewMyFeatureHandler(logger *slog.Logger) *MyFeatureHandler {
    return &MyFeatureHandler{logger: logger}
}

func (h *MyFeatureHandler) Info() module.Info {
    return module.Info{
        ID:          model.MyFeatureModuleID,
        Name:        "My Feature",
        Description: "Does something cool",
        AlwaysOn:    false,
    }
}

func (h *MyFeatureHandler) Commands() []discord.ApplicationCommandCreate {
    return []discord.ApplicationCommandCreate{
        discord.SlashCommandCreate{
            Name:        "mycommand",
            Description: "Does the thing",
        },
    }
}

func (h *MyFeatureHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
    ui := view.BuildMyFeatureResponse()
    _ = e.CreateMessage(discord.NewMessageCreateV2(ui))
}

func (h *MyFeatureHandler) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (h *MyFeatureHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (h *MyFeatureHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent { return nil }
```

### Step 3: Create `internal/view/myfeature.go`

```go
package view

import "github.com/disgoorg/disgo/discord"

func BuildMyFeatureResponse() discord.ContainerComponent {
    return discord.NewContainer(
        discord.NewTextDisplay("Hello from my module!"),
    )
}
```

### Step 4: Register in `cmd/pedmin/main.go`

```go
myFeatureHandler := handler.NewMyFeatureHandler(logger)
b.Register(myFeatureHandler)
```

## Adding a Service Layer

For features with business logic, create a service struct:

### `internal/service/myfeature.go`

```go
package service

import "github.com/s12kuma01/pedmin/internal/repository"

type MyFeatureService struct {
    store repository.GuildStore
}

func NewMyFeatureService(store repository.GuildStore) *MyFeatureService {
    return &MyFeatureService{store: store}
}

func (s *MyFeatureService) DoSomething() string {
    return "result"
}
```

Then inject the service into the handler:

```go
// handler
type MyFeatureHandler struct {
    svc    *service.MyFeatureService
    logger *slog.Logger
}
```

## Layer Responsibilities

| Layer | Package | Does | Does NOT |
|-------|---------|------|----------|
| **Handler** | `internal/handler/` | Receive events, call service, build view response | Contain business logic or store access |
| **Service** | `internal/service/` | Business logic, API calls, store operations | Build Discord UI components |
| **View** | `internal/view/` | Build Components V2 UI (pure functions) | Have side effects |
| **Model** | `internal/model/` | Domain types, settings structs, constants | Import handler/service/view |

### Handler → Service → View flow

```go
// handler/myfeature.go: orchestration only
func (h *MyFeatureHandler) HandleComponent(e *events.ComponentInteractionCreate) {
    result := h.svc.DoSomething()                 // service
    ui := view.BuildMyFeatureResult(result)        // view
    _ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

// service/myfeature.go: logic + store/API access
func (s *MyFeatureService) DoSomething() string {
    return "result"
}

// view/myfeature.go: pure function, state in → components out
func BuildMyFeatureResult(result string) discord.ContainerComponent {
    return discord.NewContainer(
        discord.NewTextDisplay(result),
    )
}
```

## Module Interface Reference

| Method | When Called | Purpose |
|--------|------------|---------|
| `Info()` | Registration, routing | Module metadata |
| `Commands()` | `SyncCommands` | Slash command definitions |
| `HandleCommand(e)` | User runs a slash command | Process command |
| `HandleComponent(e)` | User clicks button/select with matching CustomID | Handle interaction |
| `HandleModal(e)` | User submits modal with matching CustomID | Process modal data |
| `SettingsPanel(guildID)` | User views module in /settings | Return settings UI |

### Optional Interfaces

| Interface | Method | Purpose |
|-----------|--------|---------|
| `SettingsSummarizer` | `SettingsSummary(guildID) string` | Brief settings summary in /settings list |
| `VoiceStateListener` | `OnVoiceStateUpdate(guildID, channelID, userID)` | React to voice state changes |

## CustomID Convention

```
{moduleID}:{action}:{extra}
```

- **moduleID**: Must match `Info().ID` — use `model.MyFeatureModuleID`
- **action**: Operation name (e.g., `pause`, `skip`, `toggle`)
- **extra**: Optional data (e.g., target ID, page number)

Examples:

```
player:pause              → Player module, pause action
player:vol_up             → Player module, volume up
settings:toggle:player    → Settings module, toggle player
```

The router splits on the first colon to find the module. The handler parses the rest.

## Settings Integration

Return components from `SettingsPanel()` to show in the /settings detail view:

```go
func (h *MyFeatureHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
    return []discord.LayoutComponent{
        discord.NewTextDisplay("**My Feature Settings**"),
        discord.NewActionRow(
            discord.NewSecondaryButton("Configure", model.MyFeatureModuleID+":configure"),
        ),
    }
}
```

## Per-Module Settings Persistence

Use the generic helpers in `internal/repository/`:

```go
// In model
type MyFeatureSettings struct {
    Option string `json:"option"`
}

// In service
func (s *MyFeatureService) LoadSettings(guildID snowflake.ID) (*model.MyFeatureSettings, error) {
    return repository.LoadModuleSettings(s.store, guildID, model.MyFeatureModuleID, func() *model.MyFeatureSettings {
        return &model.MyFeatureSettings{}
    })
}

func (s *MyFeatureService) SaveSettings(guildID snowflake.ID, settings *model.MyFeatureSettings) error {
    return repository.SaveModuleSettings(s.store, guildID, model.MyFeatureModuleID, settings)
}
```

## Adding Event Listeners

For features that listen to Discord events (not slash commands), create a setup function:

```go
// handler/myfeature_listener.go
package handler

import disgobot "github.com/disgoorg/disgo/bot"

func SetupMyFeatureListeners(client *disgobot.Client, h *MyFeatureHandler) {
    client.AddEventListeners(disgobot.NewListenerFunc(h.onMessageCreate))
}
```

Register in `cmd/pedmin/main.go`:

```go
handler.SetupMyFeatureListeners(b.Client, myFeatureHandler)
```

## Bot Interface Pattern

If your handler needs bot methods (e.g., `IsModuleEnabled`), define a local interface:

```go
// handler/myfeature.go
type MyFeatureBot interface {
    IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type MyFeatureHandler struct {
    bot    MyFeatureBot
    logger *slog.Logger
}
```

The `*bot.Bot` struct satisfies this interface. Pass it in `main.go`:

```go
myFeatureHandler := handler.NewMyFeatureHandler(b, logger)
```
