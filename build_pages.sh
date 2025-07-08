#!/bin/bash

# GitHub Pages用のWebGL/WASMビルドスクリプト
# 使用方法: ./build_pages.sh

set -e

echo "🏮 江戸前エルフ - GitHub Pages用ビルド開始"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 関数: ログ出力
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Go環境の確認
log_info "Go環境を確認中..."
if ! command -v go &> /dev/null; then
    log_error "Goがインストールされていません"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
log_info "Go バージョン: $GO_VERSION"

# WebAssembly環境の設定
export GOOS=js
export GOARCH=wasm

log_info "WebAssembly環境設定: GOOS=$GOOS, GOARCH=$GOARCH"

# docsディレクトリの作成
log_info "docsディレクトリを作成中..."
mkdir -p docs

# 既存ファイルのクリーンアップ
log_info "既存のビルドファイルをクリーンアップ中..."
rm -f docs/game.wasm docs/wasm_exec.js

# 依存関係のダウンロード
log_info "Go依存関係をダウンロード中..."
go mod download

# WebAssemblyのビルド
log_info "WebAssemblyファイルをビルド中..."
go build -ldflags="-s -w" -o docs/game.wasm miko_game_with_worshippers.go

if [ ! -f "docs/game.wasm" ]; then
    log_error "WebAssemblyファイルのビルドに失敗しました"
    exit 1
fi

log_success "WebAssemblyファイルのビルドが完了しました"

# wasm_exec.jsのコピー
log_info "wasm_exec.jsをコピー中..."
GOROOT_VAR=$(go env GOROOT)

# 複数の場所を試す（Go versionによって場所が異なる）
WASM_EXEC_PATHS=(
    "$GOROOT_VAR/misc/wasm/wasm_exec.js"
    "$GOROOT_VAR/lib/wasm/wasm_exec.js"
)

WASM_EXEC_PATH=""
for path in "${WASM_EXEC_PATHS[@]}"; do
    if [ -f "$path" ]; then
        WASM_EXEC_PATH="$path"
        break
    fi
done

if [ -z "$WASM_EXEC_PATH" ]; then
    log_error "wasm_exec.jsが見つかりません。以下の場所を確認しました:"
    for path in "${WASM_EXEC_PATHS[@]}"; do
        echo "  - $path"
    done
    exit 1
fi

log_info "wasm_exec.jsが見つかりました: $WASM_EXEC_PATH"
cp "$WASM_EXEC_PATH" docs/
log_success "wasm_exec.jsのコピーが完了しました"

# アセットファイルのコピー
if [ -d "assets" ]; then
    log_info "アセットファイルをコピー中..."
    cp -r assets docs/
    log_success "アセットファイルのコピーが完了しました"
else
    log_warning "assetsディレクトリが見つかりません"
fi

# ビルド結果の確認
log_info "ビルド結果を確認中..."
echo ""
echo "📁 docsディレクトリの内容:"
ls -la docs/

echo ""
echo "📊 ファイルサイズ:"
if [ -f "docs/game.wasm" ]; then
    WASM_SIZE=$(du -h docs/game.wasm | cut -f1)
    echo "  📦 game.wasm: $WASM_SIZE"
fi

if [ -f "docs/wasm_exec.js" ]; then
    EXEC_SIZE=$(du -h docs/wasm_exec.js | cut -f1)
    echo "  📄 wasm_exec.js: $EXEC_SIZE"
fi

if [ -f "docs/index.html" ]; then
    HTML_SIZE=$(du -h docs/index.html | cut -f1)
    echo "  🌐 index.html: $EXEC_SIZE"
fi

echo ""
log_success "GitHub Pages用ビルドが完了しました!"

echo ""
echo "🚀 次のステップ:"
echo "  1. GitHub Pagesを有効にしてください (Settings > Pages > Source: GitHub Actions)"
echo "  2. mainブランチにプッシュすると自動デプロイされます"
echo "  3. ローカルでテストする場合: ./serve_pages.sh"

echo ""
echo "🌐 GitHub Pagesが有効になったら以下のURLでアクセスできます:"
echo "  https://noppiki.github.io/EdomaeElf/"