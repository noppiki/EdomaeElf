# 🏮 EdomaeElf WebGL版 - 参拝客システム

参拝客システム付きの巫女さんゲームをWebブラウザで実行するためのガイドです。

## 🎮 概要

このプロジェクトは、Ebitengine（旧Ebiten）を使用したGo言語ゲームのWebGL/WebAssembly版です。ランダムなタイミングで参拝客が訪れて、賽銭を入れて去っていくシステムを実装しています。

## 🚀 クイックスタート

### 1. ビルド
```bash
./build_webgl.sh
```

### 2. サーバー起動
```bash
./serve_webgl.sh
```

### 3. ブラウザでアクセス
```
http://localhost:8080
```

## 🔧 詳細な手順

### 前提条件
- Go 1.11以上（WebAssemblyサポートが必要）
- WebサーバーまたはPython3/Python2/PHP

### 手動ビルド手順

1. **WebAssemblyファイルの生成**
   ```bash
   GOOS=js GOARCH=wasm go build -o web/main.wasm miko_game_with_worshippers.go
   ```

2. **wasm_exec.jsの取得**
   ```bash
   cp $(go env GOROOT)/misc/wasm/wasm_exec.js web/
   ```

3. **アセットのコピー**
   ```bash
   cp -r assets web/
   ```

4. **Webサーバーの起動**
   ```bash
   cd web && python3 -m http.server 8080
   ```

### 手動サーバー起動オプション

#### Python3
```bash
cd web && python3 -m http.server 8080
```

#### Python2
```bash
cd web && python -m SimpleHTTPServer 8080
```

#### PHP
```bash
cd web && php -S localhost:8080
```

#### Node.js (http-server)
```bash
cd web && npx http-server -p 8080
```

## 🎯 ゲーム機能

### 参拝客システム
- **ランダム出現**: 5秒間隔で30%の確率で出現
- **AI行動**: 賽銭箱に向かって移動 → 参拝 → 退場
- **状態管理**: 接近中・参拝中・退場中の3状態
- **視覚効果**: 色バリエーション・アニメーション・フェードアウト

### 操作方法
- **WASD または 矢印キー**: プレイヤー移動
- **E**: 編集モード切替
- **編集モード時**:
  - **Q/R**: タイルX選択
  - **T/Y**: タイルY選択
  - **左クリック**: タイル配置
  - **Space**: カメラリセット

## 📁 ファイル構成

```
web/
├── index.html      # HTMLテンプレート
├── main.wasm       # WebAssemblyゲーム本体
├── wasm_exec.js    # Go WebAssemblyランタイム
└── assets/         # ゲームアセット
    ├── characters/
    │   └── miko_girl.png
    └── tilemap/
        └── japanese_town_tileset.png
```

## 🐛 トラブルシューティング

### よくある問題

#### 1. "WebAssembly読み込みエラー"
- **原因**: file://でアクセスしている
- **解決**: 必ずHTTPサーバー経由でアクセスしてください

#### 2. "wasm_exec.js が見つかりません"
- **原因**: Go環境の問題
- **解決**: 
  ```bash
  cp $(go env GOROOT)/misc/wasm/wasm_exec.js web/
  ```

#### 3. "アセットが読み込めません"
- **原因**: assetsディレクトリのパス問題
- **解決**: assetsディレクトリをwebディレクトリにコピー

#### 4. "ゲームが重い"
- **原因**: WebAssemblyの性能制限
- **解決**: 
  - ブラウザのタブを1つだけ開く
  - ハードウェアアクセラレーションを有効化
  - 最新のブラウザを使用

### デバッグ方法

1. **ブラウザの開発者ツールを開く** (F12)
2. **コンソールタブでエラーを確認**
3. **ネットワークタブで読み込み状況を確認**

## 🌐 ブラウザ対応

### 対応ブラウザ
- **Chrome**: ✅ 推奨
- **Firefox**: ✅ 対応
- **Safari**: ✅ 対応
- **Edge**: ✅ 対応

### 必要な機能
- WebAssembly対応
- WebGL対応
- HTML5 Canvas対応

## 🚀 デプロイ

### 静的ホスティング
webディレクトリの内容を以下のサービスにデプロイできます：

- **GitHub Pages**
- **Netlify**
- **Vercel**
- **Firebase Hosting**

### 設定例 (GitHub Pages)
```yaml
# .github/workflows/deploy.yml
name: Deploy to GitHub Pages
on:
  push:
    branches: [ main ]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.19
    - run: ./build_webgl.sh
    - uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./web
```

## 📊 パフォーマンス

### 最適化のヒント
- **アセット圧縮**: PNG画像の最適化
- **WASM圧縮**: gzip圧縮の有効化
- **キャッシュ**: 適切なキャッシュヘッダーの設定

### 一般的なパフォーマンス
- **読み込み時間**: 2-5秒（初回）
- **FPS**: 30-60fps（デバイス依存）
- **メモリ使用量**: 20-50MB

## 🔄 更新手順

ゲームを更新する場合：

1. **ソースコードを修正**
2. **再ビルド**
   ```bash
   ./build_webgl.sh
   ```
3. **ブラウザでリロード** (Ctrl+F5)

## 🆘 サポート

### 問題報告
- **Issue**: GitHubのIssueタブで報告
- **ログ**: ブラウザコンソールのエラーログを添付

### 開発者向け
- **Go version**: `go version`
- **Ebitengine version**: `v2.8.8`
- **WebAssembly**: `GOOS=js GOARCH=wasm`

---

## 🎉 楽しいゲーム体験を！

参拝客システムを楽しんでください。新しい参拝客が現れるたびに、賽銭がどんどん増えていきます！

---

**開発者**: EdomaeElf プロジェクト  
**技術**: Go + Ebitengine + WebAssembly  
**バージョン**: 1.0.0