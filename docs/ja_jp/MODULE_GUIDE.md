# モジュール開発ガイド

## アーキテクチャ概要

各モジュールは複数のレイヤーパッケージにまたがって実装されます。Handler 構造体は `module.Module` インターフェースを実装し、`cmd/pedmin/main.go` で登録されます。

```
internal/
├── handler/myfeature.go       # Handler 構造体、module.Module の実装
├── service/myfeature.go       # ビジネスロジック（任意）
├── view/myfeature.go          # UI ビルダー（任意）
└── model/myfeature.go         # ドメイン型（任意）
```

Handler は `package handler`、Service は `package service`、View は `package view`、Model は `package model` を使用します。

## 新しいモジュールの作成

### ステップ 1: `internal/model/constants.go` に ModuleID を追加

```go
const MyFeatureModuleID = "myfeature"
```

### ステップ 2: `internal/handler/myfeature.go` を作成

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

### ステップ 3: `internal/view/myfeature.go` を作成

```go
package view

import "github.com/disgoorg/disgo/discord"

func BuildMyFeatureResponse() discord.ContainerComponent {
    return discord.NewContainer(
        discord.NewTextDisplay("Hello from my module!"),
    )
}
```

### ステップ 4: `cmd/pedmin/main.go` で登録

```go
myFeatureHandler := handler.NewMyFeatureHandler(logger)
b.Register(myFeatureHandler)
```

## Service レイヤーの追加

ビジネスロジックを含む機能には、Service 構造体を作成します。

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

その後、Service を Handler に注入します。

```go
type MyFeatureHandler struct {
    svc    *service.MyFeatureService
    logger *slog.Logger
}
```

## レイヤーの責務

| レイヤー | パッケージ | 責務 | 責務外 |
|---------|-----------|------|--------|
| **Handler** | `internal/handler/` | イベントの受信、Service の呼び出し、View レスポンスの構築 | ビジネスロジックやストアへのアクセス |
| **Service** | `internal/service/` | ビジネスロジック、API 呼び出し、ストア操作 | Discord UI コンポーネントの構築 |
| **View** | `internal/view/` | Components V2 UI の構築（純粋関数） | 副作用の発生 |
| **Model** | `internal/model/` | ドメイン型、設定構造体、定数 | handler/service/view のインポート |

### Handler → Service → View のフロー

```go
// handler/myfeature.go: オーケストレーションのみ
func (h *MyFeatureHandler) HandleComponent(e *events.ComponentInteractionCreate) {
    result := h.svc.DoSomething()
    ui := view.BuildMyFeatureResult(result)
    _ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

// service/myfeature.go: ロジック + ストア/API アクセス
func (s *MyFeatureService) DoSomething() string {
    return "result"
}

// view/myfeature.go: 純粋関数、状態入力 → コンポーネント出力
func BuildMyFeatureResult(result string) discord.ContainerComponent {
    return discord.NewContainer(
        discord.NewTextDisplay(result),
    )
}
```

## Module インターフェースリファレンス

| メソッド | 呼び出しタイミング | 目的 |
|---------|------------------|------|
| `Info()` | 登録時、ルーティング時 | モジュールのメタデータ |
| `Commands()` | `SyncCommands` 実行時 | スラッシュコマンドの定義 |
| `HandleCommand(e)` | ユーザーがスラッシュコマンドを実行した時 | コマンドの処理 |
| `HandleComponent(e)` | ユーザーが一致する CustomID のボタン/セレクトをクリックした時 | インタラクションの処理 |
| `HandleModal(e)` | ユーザーが一致する CustomID のモーダルを送信した時 | モーダルデータの処理 |
| `SettingsPanel(guildID)` | ユーザーが /settings でモジュールを表示した時 | 設定 UI の返却 |

### オプションインターフェース

| インターフェース | メソッド | 目的 |
|----------------|---------|------|
| `SettingsSummarizer` | `SettingsSummary(guildID) string` | /settings 一覧での簡潔な設定概要 |
| `VoiceStateListener` | `OnVoiceStateUpdate(guildID, channelID, userID)` | ボイスステート変更への対応 |

## CustomID 規約

```
{moduleID}:{action}:{extra}
```

- **moduleID**: `Info().ID` と一致する必要がある -- `model.MyFeatureModuleID` を使用
- **action**: 操作名
- **extra**: オプションのデータ

ルーターは最初のコロンで分割してモジュールを特定します。Handler が残りの部分を解析します。

## 設定画面の統合

`SettingsPanel()` からコンポーネントを返すと、/settings の詳細ビューに表示されます。

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

## モジュール別設定の永続化

`internal/repository/` の汎用ヘルパーを使用します。

```go
func (s *MyFeatureService) LoadSettings(guildID snowflake.ID) (*model.MyFeatureSettings, error) {
    return repository.LoadModuleSettings(s.store, guildID, model.MyFeatureModuleID, func() *model.MyFeatureSettings {
        return &model.MyFeatureSettings{}
    })
}

func (s *MyFeatureService) SaveSettings(guildID snowflake.ID, settings *model.MyFeatureSettings) error {
    return repository.SaveModuleSettings(s.store, guildID, model.MyFeatureModuleID, settings)
}
```

## イベントリスナーの追加

```go
// handler/myfeature_listener.go
func SetupMyFeatureListeners(client *disgobot.Client, h *MyFeatureHandler) {
    client.AddEventListeners(disgobot.NewListenerFunc(h.onMessageCreate))
}
```

`cmd/pedmin/main.go` で登録します。

```go
handler.SetupMyFeatureListeners(b.Client, myFeatureHandler)
```

## Bot インターフェースパターン

Handler が Bot のメソッドを必要とする場合、ローカルインターフェースを定義します。

```go
type MyFeatureBot interface {
    IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type MyFeatureHandler struct {
    bot    MyFeatureBot
    logger *slog.Logger
}
```
