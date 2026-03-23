# データ永続化ガイド

## GuildStore インターフェース

`internal/repository/repository.go` で定義されています。データ型は `internal/model/guild.go` にあります。

```go
type GuildStore interface {
    // ギルド設定
    Get(guildID snowflake.ID) (*model.GuildSettings, error)
    Save(settings *model.GuildSettings) error
    IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
    SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
    GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
    SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error

    // チケット
    CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
    GetTicketByChannel(channelID snowflake.ID) (*model.Ticket, error)
    CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
    DeleteTicket(channelID snowflake.ID) error

    // RSS
    CreateRSSFeed(feed *model.RSSFeed) error
    DeleteRSSFeed(id int64, guildID snowflake.ID) error
    GetRSSFeeds(guildID snowflake.ID) ([]model.RSSFeed, error)
    GetAllRSSFeeds() ([]model.RSSFeed, error)
    CountRSSFeeds(guildID snowflake.ID) (int, error)
    IsItemSeen(feedID int64, itemHash string) (bool, error)
    MarkItemsSeen(feedID int64, itemHashes []string) error
    PruneSeenItems(olderThan time.Time) error

    // ライフサイクル
    Close() error
}
```

すべてのメソッドは並行利用に対して安全です。

## データ型

`internal/model/guild.go` で定義されています。

### GuildSettings

```go
type GuildSettings struct {
    GuildID        snowflake.ID          `json:"guild_id"`
    EnabledModules map[string]bool       `json:"enabled_modules"`
    ModuleSettings map[string]any        `json:"module_settings"`
}
```

| フィールド | 型 | 説明 |
|-----------|------|------|
| `GuildID` | `snowflake.ID` | Discord ギルド ID |
| `EnabledModules` | `map[string]bool` | モジュール ID → 有効状態 |
| `ModuleSettings` | `map[string]any` | モジュール ID → 任意の設定データ（`Get`/`Save` の内部使用） |

デフォルト: レコードが存在しない場合、すべてのモジュールが無効の空マップを返します。

### Ticket

```go
type Ticket struct {
    GuildID   snowflake.ID
    Number    int
    ChannelID snowflake.ID
    UserID    snowflake.ID
    Subject   string
    CreatedAt time.Time
    ClosedAt  *time.Time
    ClosedBy  *snowflake.ID
}
```

### RSSFeed

```go
type RSSFeed struct {
    ID        int64
    GuildID   snowflake.ID
    URL       string
    ChannelID snowflake.ID
    Title     string
    AddedAt   time.Time
}
```

## SQLiteStore 実装

`internal/repository/sqlite.go` にあります。

### データベースの場所

```
{DB_PATH}  (デフォルト: {DATA_DIR}/pedmin.db)
```

`config.toml` の `storage.db_path` 設定、または `DB_PATH` 環境変数で上書きできます。

### スキーマ

SQL マイグレーションファイルは `migrations/` に格納され、`embed.FS` 経由で読み込まれます。

### 設定

- **WAL モード**: 複数の読み取りと単一の書き込みを同時実行可能
- **busy_timeout=5000**: ロック競合時に 5 秒間待機
- **純粋な Go ドライバー**: `modernc.org/sqlite`（CGO 不要）

### パフォーマンス

| 操作 | クエリ | 計算量 |
|------|--------|--------|
| `IsModuleEnabled` | 主キーによる単一行 SELECT | O(1) |
| `SetModuleEnabled` | 単一 UPSERT | O(1) |
| `Get` | 2 クエリ（モジュール + 設定） | O(n モジュール) |
| `Save` | UPSERT を含むトランザクション | O(n モジュール) |
| `GetTicketByChannel` | 主キーによる単一行 SELECT | O(1) |
| `GetRSSFeeds` | guild_id による SELECT | O(n フィード) |
| `IsItemSeen` | 主キーによる単一行 SELECT | O(1) |

### スキーママイグレーション

`schema_migrations` テーブルで管理されます。SQL ファイルは `migrations/` に格納され、起動時に `embed.FS` 経由で読み込まれます。新しいマイグレーションは `migrations/NNN_description.sql` として追加してください。バージョン番号はファイル名のプレフィックスから解析されます。`repository.NewSQLiteStore()` の実行時に自動的に適用されます。

## 新しいストア実装の追加

1. `repository.GuildStore` インターフェース（`Close() error` を含む）を実装する
2. コンストラクタ関数を追加する
3. `cmd/pedmin/main.go` の初期化処理を差し替える

## ModuleSettings の使い方

`internal/repository/module_settings.go` の汎用ヘルパーを使用します。

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

汎用ヘルパーは JSON のマーシャリング/アンマーシャリングを自動的に処理し、デフォルトファクトリによるフォールバックを提供します。
