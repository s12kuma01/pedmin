# Lavalink 統合ガイド

## セットアップ

### Docker Compose

Lavalink はボットと並行して Docker コンテナとして実行される。

```yaml
services:
  lavalink:
    image: ghcr.io/lavalink-devs/lavalink:4-alpine
    volumes:
      - ./lavalink/application.yml:/opt/Lavalink/application.yml
```

### 設定 (`lavalink/application.yml`)

主要な設定項目:

- `server.port`: WebSocket/REST ポート (デフォルト: 2333)
- `lavalink.server.password`: 認証パスワード
- `lavalink.server.sources`: 音声ソースの有効化/無効化
- `lavalink.plugins`: プラグインの依存関係とリポジトリ

### 環境変数

| 変数                | デフォルト値      | 説明                       |
|---------------------|-------------------|----------------------------|
| `LAVALINK_HOST`     | `lavalink:2333`   | Lavalink のアドレス        |
| `LAVALINK_PASSWORD` | `youshallnotpass`  | Lavalink の認証パスワード  |

## disgolink API

### ノードへの接続

```go
link := disgolink.New(botAppID)
_, err := link.AddNode(ctx, disgolink.NodeConfig{
    Name:     "main",
    Address:  "lavalink:2333",
    Password: "youshallnotpass",
})
```

### トラックの読み込み

```go
node := link.BestNode()

// ハンドラコールバックを使用
node.LoadTracksHandler(ctx, "ytsearch:query", disgolink.NewResultHandler(
    func(track lavalink.Track) { /* 単一トラック */ },
    func(playlist lavalink.Playlist) { /* プレイリスト読み込み */ },
    func(tracks []lavalink.Track) { /* 検索結果 */ },
    func() { /* 一致なし */ },
    func(err error) { /* 読み込み失敗 */ },
))

// または直接結果を取得
result, err := node.LoadTracks(ctx, "ytsearch:query")
```

### プレイヤー操作

```go
player := link.Player(guildID)

// トラックを再生
player.Update(ctx, lavalink.WithTrack(track))

// 一時停止/再開
player.Update(ctx, lavalink.WithPaused(true))

// 音量 (0-200)
player.Update(ctx, lavalink.WithVolume(50))

// シーク
player.Update(ctx, lavalink.WithPosition(30000)) // 30秒

// 停止 (トラックをクリア)
player.Update(ctx, lavalink.WithNullTrack())

// プレイヤーを破棄
player.Destroy(ctx)
link.RemovePlayer(guildID)
```

### プレイヤー状態の取得

```go
player.Track()    // *lavalink.Track (再生中でなければ nil)
player.Paused()   // bool
player.Position() // lavalink.Duration (ミリ秒)
player.Volume()   // int (0-200)
```

## イベントリスナー

disgolink クライアントにリスナーを登録する。

```go
link.AddListeners(
    disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.TrackStartEvent) {
        // トラックの再生が開始された
    }),
    disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.TrackEndEvent) {
        // トラックが終了した - event.Reason を確認
        if event.Reason == lavalink.TrackEndReasonFinished {
            // 次のトラックを再生
        }
    }),
)
```

### イベントの種類

| イベント               | 発生タイミング                     |
|------------------------|------------------------------------|
| `TrackStartEvent`      | トラックの再生が開始された時       |
| `TrackEndEvent`        | トラックが終了/失敗/置換された時   |
| `TrackExceptionEvent`  | 再生エラーが発生した時             |
| `TrackStuckEvent`      | トラックがスタックした時 (音声フレームなし) |
| `WebSocketClosedEvent` | 音声 WebSocket が閉じられた時      |

### TrackEndReason の値

| 理由                       | 意味               | 対処                  |
|----------------------------|--------------------|-----------------------|
| `TrackEndReasonFinished`   | 正常終了           | 次のトラックを再生    |
| `TrackEndReasonLoadFailed` | 読み込み失敗       | 次を再生 / 通知       |
| `TrackEndReasonStopped`    | 手動で停止された   | 何もしない            |
| `TrackEndReasonReplaced`   | 別のトラックが開始 | 何もしない            |
| `TrackEndReasonCleanup`    | プレイヤーが破棄   | 何もしない            |

## 音声接続

Lavalink が動作するには、音声状態/サーバー更新を転送する必要がある。

```go
// ボットの初期化時
bot.WithEventListenerFunc(func(e *events.GuildVoiceStateUpdate) {
    if e.VoiceState.UserID != client.ApplicationID {
        return
    }
    link.OnVoiceStateUpdate(ctx, e.VoiceState.GuildID, e.VoiceState.ChannelID, e.VoiceState.SessionID)
})

bot.WithEventListenerFunc(func(e *events.VoiceServerUpdate) {
    link.OnVoiceServerUpdate(ctx, e.GuildID, e.Token, *e.Endpoint)
})
```

ボイスチャンネルに参加するには:

```go
client.UpdateVoiceState(ctx, guildID, &channelID, false, true) // selfMute=false, selfDeaf=true
```

## プラグイン

### lavasearch-plugin

複数のソースにまたがる拡張検索を提供する。

```yaml
plugins:
  lavasearch:
    sources:
      - youtube
      - soundcloud
```

### lavalyrics-plugin

再生中のトラックの歌詞を提供する。

```yaml
plugins:
  lavalyrics:
    sources:
      - youtube
```

## トラブルシューティング

| 問題                         | 解決方法                                                     |
|------------------------------|--------------------------------------------------------------|
| `no lavalink node available` | Lavalink が起動中でアクセス可能か確認する                    |
| 接続拒否                     | 環境変数のホスト/ポートが Lavalink の設定と一致しているか確認 |
| 音声が出ない                 | ボットがデフンされているか、音声イベントが転送されているか確認 |
| トラック読み込み失敗         | Lavalink のログでソース固有のエラーを確認する                |
| WebSocket closed (4014)      | ボットに `GuildVoiceStates` Intent が必要                    |
