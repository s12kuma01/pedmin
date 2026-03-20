# Module Development Guide

## Architecture Overview

Each module follows the **Feature Module Pattern**: a single Go package containing handler, service, and view layers separated by file.

```
features/myfeature/
├── module.go              # Module interface, struct, constructor
├── handler.go             # HandleCommand / HandleComponent / HandleModal
├── service.go             # Business logic (optional, for complex modules)
└── view.go                # UI builders (optional, if module has UI)
```

All files share `package myfeature`. No sub-packages needed.

## Creating a New Module

### Step 1: Create `features/myfeature/module.go`

```go
package myfeature

import (
    "log/slog"

    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
    "github.com/disgoorg/snowflake/v2"
    "github.com/s12kuma01/pedmin/module"
)

const ModuleID = "myfeature"

type MyFeature struct {
    logger *slog.Logger
}

func New(logger *slog.Logger) *MyFeature {
    return &MyFeature{logger: logger}
}

func (m *MyFeature) Info() module.Info {
    return module.Info{
        ID:          ModuleID,
        Name:        "My Feature",
        Description: "Does something cool",
        AlwaysOn:    false,
    }
}

func (m *MyFeature) Commands() []discord.ApplicationCommandCreate {
    return []discord.ApplicationCommandCreate{
        discord.SlashCommandCreate{
            Name:        "mycommand",
            Description: "Does the thing",
        },
    }
}

func (m *MyFeature) HandleModal(_ *events.ModalSubmitInteractionCreate) {}
func (m *MyFeature) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent { return nil }
func (m *MyFeature) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
```

### Step 2: Create `features/myfeature/handler.go`

```go
package myfeature

import (
    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
)

func (m *MyFeature) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
    _ = e.CreateMessage(discord.NewMessageCreateV2(
        discord.NewContainer(
            discord.NewTextDisplay("Hello from my module!"),
        ),
    ))
}

func (m *MyFeature) HandleComponent(e *events.ComponentInteractionCreate) {
    // Handle button/select interactions
}
```

### Step 3: Register in `main.go`

```go
myModule := myfeature.New(logger)
b.Register(myModule)
```

## Layer Responsibilities

| Layer | File pattern | Does | Does NOT |
|-------|-------------|------|----------|
| **Module** | `module.go` | Define struct, Info, Commands, stubs | Contain logic |
| **Handler** | `handler*.go` | Parse interaction, call service, build response | Contain business logic, API calls, or store access |
| **Service** | `service*.go` | Business logic, API calls, store operations | Build Discord UI components |
| **View** | `view*.go` | Build Components V2 UI (pure functions: data in → components out) | Have side effects, call APIs, or access store |
| **Client** | `client.go` | External API HTTP wrappers | Contain business logic |
| **Domain** | `queue.go` etc. | Data structures, types | Import Discord packages |

### Layer Rules

1. **Handlers are dispatchers only.** A handler method should: defer the interaction → call service → call view → respond. No business logic, no direct API calls, no store access.
2. **Services own all logic.** Any operation that involves validation, state changes, external API calls, or store reads/writes belongs in `service*.go`.
3. **Views are pure functions.** They take data as input and return Discord components as output. No side effects.
4. **Clients are HTTP wrappers.** They translate Go method calls to HTTP requests and responses. No business logic.

### Handler → Service → View flow

```go
// handler.go: orchestration only
func (m *MyFeature) HandleComponent(e *events.ComponentInteractionCreate) {
    result := m.doSomething(guildID)       // service.go
    ui := buildResultUI(result)            // view.go
    _ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

// service.go: logic + store/API access
func (m *MyFeature) doSomething(guildID snowflake.ID) string {
    return "result"
}

// view.go: pure function, state in → components out
func buildResultUI(result string) discord.ContainerComponent {
    return discord.NewContainer(
        discord.NewTextDisplay(result),
    )
}
```

## Module Interface Reference

| Method | When Called | Purpose |
|--------|-----------|---------|
| `Info()` | Registration, routing | Module metadata |
| `Commands()` | `SyncCommands` | Slash command definitions |
| `HandleCommand(e)` | User runs a slash command | Process command |
| `HandleComponent(e)` | User clicks button/select with matching CustomID | Handle interaction |
| `HandleModal(e)` | User submits modal with matching CustomID | Process modal data |
| `SettingsPanel(guildID)` | User views module in /settings | Return settings UI |
| `HandleSettingsComponent(e)` | User interacts in settings panel | Handle settings interaction |

## CustomID Convention

```
{moduleID}:{action}:{extra}
```

- **moduleID**: Must match `Info().ID`
- **action**: Operation name (e.g., `pause`, `skip`, `toggle`)
- **extra**: Optional data (e.g., target ID, page number)

Examples:
```
player:pause              → Player module, pause action
player:vol_up             → Player module, volume up
settings:toggle:player    → Settings module, toggle player
```

The router splits on the first colon to find the module. The module parses the rest.

## Settings Integration

Return components from `SettingsPanel()` to show in the /settings detail view:

```go
func (m *MyFeature) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
    return []discord.LayoutComponent{
        discord.NewTextDisplay("**My Feature Settings**"),
        discord.NewActionRow(
            discord.NewSecondaryButton("Configure", ModuleID+":configure"),
        ),
    }
}
```
