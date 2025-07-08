## Protocols

### 作業完了時の通知ルール

#### 1. 必須通知
- **claude codeの作業が完了する際には、作業内容を必ずアラートダイアログで通知すること**
- タグ: `ultrathink`

#### 2. 通知方法
作業完了時は以下のコマンドを使用してアラートダイアログを表示：
```bash
osascript -e 'display alert "Claude Code - 作業完了" message "実行した作業内容の詳細をここに記載" as informational'
```

#### 3. 通知内容に含めるべき情報
- 実行したタスクの概要
- 作成・変更したファイルのリスト
- 重要な変更点や注意事項
- 次のステップ（ある場合）

#### 4. 通知タイミング
以下の場合に必ず通知を送信：
- ユーザーから依頼された主要タスクが完了した時
- 複数のファイルを作成・編集した後
- ビルドやテストの実行が完了した時
- エラーが発生して作業を中断する時

#### 5. 通知の例
```bash
# 例1: 機能追加完了
osascript -e 'display alert "Claude Code - 作業完了" message "Ebitengineのセットアップが完了しました。
作成ファイル:
- main.go (基本サンプル)
- main_advanced.go (高度な機能)
- main_sprite.go (スプライト操作)
- notify.go (通知ヘルパー)

実行: go run main.go で動作確認できます。" as informational'

# 例2: エラー発生時
osascript -e 'display alert "Claude Code - エラー" message "ビルドエラーが発生しました。
エラー内容: undefined variable
該当ファイル: main.go:25
修正が必要です。" as critical'
```

### プロジェクト固有の設定

#### Go開発環境
- **言語**: Go
- **フレームワーク**: Ebitengine v2.8.8
- **プロジェクト名**: EdomaeElf

#### 開発ガイドライン
1. **コーディング規約**
   - Go標準のフォーマットに従う（`go fmt`）
   - エラーハンドリングは必ず行う
   - 変数名は意味のある名前を使用

2. **ファイル構成**
   - ゲームロジック: `main.go`, `game.go`
   - グラフィックス関連: `graphics/`
   - 音声関連: `audio/`
   - ユーティリティ: `utils/`

3. **テスト**
   - 単体テストファイルは `*_test.go` の形式
   - `go test ./...` で全テスト実行