# 🏮 江戸前エルフ - GitHub Pages対応

GitHub Pagesでランダム参拝客システムを動作させるための設定ガイドです。

## 🌐 ライブデモ

設定完了後、以下のURLでゲームをプレイできます：
**https://noppiki.github.io/EdomaeElf/**

## 🚀 クイックスタート

### 1. GitHub Pagesの有効化

1. GitHubリポジトリの **Settings** タブに移動
2. 左メニューから **Pages** を選択
3. **Source** で **GitHub Actions** を選択
4. 保存

### 2. 自動デプロイ

- `main`ブランチにプッシュすると自動的にビルド＆デプロイされます
- GitHub Actionsタブでビルドの進行状況を確認できます

### 3. ローカル開発

```bash
# ビルド
./build_pages.sh

# ローカルサーバー起動
./serve_pages.sh

# ブラウザでアクセス
open http://localhost:8080
```

## 📁 ファイル構成

```
EdomaeElf/
├── docs/                          # GitHub Pages用ファイル
│   ├── index.html                # ゲーム用HTMLテンプレート
│   ├── .nojekyll                 # Jekyll無効化
│   ├── _config.yml               # GitHub Pages設定
│   ├── _headers                  # MIME type設定
│   ├── game.wasm                 # ビルド後のWebAssemblyファイル
│   ├── wasm_exec.js              # Goランタイム
│   └── assets/                   # ゲームアセット（コピー）
├── .github/workflows/
│   └── github-pages.yml          # 自動デプロイワークフロー
├── build_pages.sh                # ローカルビルドスクリプト
└── serve_pages.sh                # ローカル開発サーバー
```

## 🎮 ゲームの特徴

### 参拝客システム
- **ランダム出現**: 5秒間隔で30%の確率で参拝客が出現
- **自動AI行動**: 賽銭箱に向かって移動 → 参拝 → 退場
- **3つの状態管理**: 
  - 🚶‍♂️ **Approaching**: 賽銭箱に移動中
  - 🙏 **Offering**: 賽銭箱で参拝中（2秒）
  - 👋 **Leaving**: 画面外へ退場中

### 視覚効果
- 5種類の色でバリエーション豊富な参拝客
- 参拝時のバウンスアニメーション
- 退場時のフェードアウト効果
- リアルタイム統計表示

### 操作方法
- **WASD / 矢印キー**: プレイヤー移動
- **E**: 編集モード切替
- **マウス**: 編集モード時の操作

## 🛠️ 技術仕様

### 使用技術
- **言語**: Go 1.21+
- **フレームワーク**: Ebitengine v2.8.8
- **ターゲット**: WebAssembly (GOOS=js, GOARCH=wasm)
- **デプロイ**: GitHub Pages + GitHub Actions

### ブラウザ対応
- ✅ Chrome 57+
- ✅ Firefox 52+
- ✅ Safari 11+
- ✅ Edge 16+

### パフォーマンス
- **初回読み込み**: ~18MB (WebAssemblyファイル)
- **実行速度**: ネイティブレベル
- **メモリ使用量**: ~50MB

## 🔧 ローカル開発

### 前提条件
```bash
# Go環境の確認
go version  # 1.21以上

# リポジトリのクローン
git clone https://github.com/noppiki/EdomaeElf.git
cd EdomaeElf
```

### ビルドコマンド
```bash
# GitHub Pages用ビルド
./build_pages.sh

# カスタムビルド（詳細制御）
export GOOS=js
export GOARCH=wasm
go build -ldflags="-s -w" -o docs/game.wasm miko_game_with_worshippers.go
```

### 開発サーバー
```bash
# デフォルト（ポート8080）
./serve_pages.sh

# カスタムポート
./serve_pages.sh 3000
```

## 🚀 デプロイメント

### 自動デプロイ（推奨）
1. `main`ブランチにコミット＆プッシュ
2. GitHub Actionsが自動的にビルド
3. GitHub Pagesに自動デプロイ
4. 数分後にサイトが更新

### 手動デプロイ
```bash
# ローカルでビルド
./build_pages.sh

# docsディレクトリをコミット
git add docs/
git commit -m "feat: update GitHub Pages build"
git push origin main
```

## 📊 モニタリング

### GitHub Actions
- **ビルドログ**: Actions タブでビルドの詳細を確認
- **デプロイ状況**: Pages タブでデプロイ状況を確認
- **エラー通知**: 失敗時にメール通知

### アクセス解析
- GitHub Insights でページビューを確認可能
- ブラウザの開発者ツールでパフォーマンス測定

## 🐛 トラブルシューティング

### よくある問題

#### 1. WASMファイルが読み込めない
```
原因: MIME typeの設定不備
解決: docs/_headers ファイルの確認
```

#### 2. ゲームが表示されない
```
原因: JavaScriptエラー
解決: ブラウザの開発者ツールでエラー確認
```

#### 3. ビルドが失敗する
```
原因: Go環境やパス設定
解決: go env の確認、Go 1.21+の使用
```

#### 4. GitHub Pagesが更新されない
```
原因: Actions権限やPages設定
解決: Settings > Actions/Pages の設定確認
```

### デバッグ方法
```bash
# ローカルビルドテスト
./build_pages.sh

# ファイル存在確認
ls -la docs/

# サーバー起動テスト
./serve_pages.sh
```

## 📞 サポート

### リンク
- [Issue報告](https://github.com/noppiki/EdomaeElf/issues)
- [GitHub Pages公式ドキュメント](https://docs.github.com/pages)
- [Ebitengine公式サイト](https://ebitengine.org/)

### 貢献
プルリクエストや機能提案を歓迎します！

---

**🎮 GitHub Pagesで楽しい参拝体験を！**