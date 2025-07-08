#!/bin/bash

# GitHub Pages用ローカル開発サーバー
# 使用方法: ./serve_pages.sh

set -e

echo "🏮 江戸前エルフ - GitHub Pages ローカルサーバー起動"

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

# docsディレクトリの存在確認
if [ ! -d "docs" ]; then
    log_error "docsディレクトリが見つかりません"
    log_info "まず ./build_pages.sh を実行してビルドしてください"
    exit 1
fi

# 必要ファイルの確認
REQUIRED_FILES=("docs/index.html" "docs/game.wasm" "docs/wasm_exec.js")
for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        log_error "必要なファイルが見つかりません: $file"
        log_info "まず ./build_pages.sh を実行してビルドしてください"
        exit 1
    fi
done

log_success "必要なファイルが揃っています"

# ポート設定
PORT=8080
if [ ! -z "$1" ]; then
    PORT=$1
fi

# サーバーの選択と起動
log_info "適切なHTTPサーバーを検索中..."

cd docs

# Python 3 を最初に試す
if command -v python3 &> /dev/null; then
    log_success "Python3が見つかりました"
    echo ""
    echo "🌐 サーバーを起動中..."
    echo "   URL: http://localhost:$PORT"
    echo "   終了: Ctrl+C"
    echo ""
    
    trap 'echo ""; log_info "サーバーを停止しています..."; exit 0' INT
    
    python3 -m http.server $PORT --bind 127.0.0.1
    
# Python 2 を次に試す
elif command -v python2 &> /dev/null; then
    log_success "Python2が見つかりました"
    echo ""
    echo "🌐 サーバーを起動中..."
    echo "   URL: http://localhost:$PORT"
    echo "   終了: Ctrl+C"
    echo ""
    
    trap 'echo ""; log_info "サーバーを停止しています..."; exit 0' INT
    
    python2 -m SimpleHTTPServer $PORT
    
# Python (バージョン不明) を試す
elif command -v python &> /dev/null; then
    log_success "Pythonが見つかりました"
    PYTHON_VERSION=$(python --version 2>&1 | awk '{print $2}' | cut -d. -f1)
    
    echo ""
    echo "🌐 サーバーを起動中..."
    echo "   URL: http://localhost:$PORT"
    echo "   終了: Ctrl+C"
    echo ""
    
    trap 'echo ""; log_info "サーバーを停止しています..."; exit 0' INT
    
    if [ "$PYTHON_VERSION" = "3" ]; then
        python -m http.server $PORT --bind 127.0.0.1
    else
        python -m SimpleHTTPServer $PORT
    fi
    
# PHP を試す
elif command -v php &> /dev/null; then
    log_success "PHPが見つかりました"
    echo ""
    echo "🌐 サーバーを起動中..."
    echo "   URL: http://localhost:$PORT"
    echo "   終了: Ctrl+C"
    echo ""
    
    trap 'echo ""; log_info "サーバーを停止しています..."; exit 0' INT
    
    php -S 127.0.0.1:$PORT
    
# Node.js を試す
elif command -v npx &> /dev/null; then
    log_success "Node.js (npx)が見つかりました"
    echo ""
    echo "🌐 サーバーを起動中..."
    echo "   URL: http://localhost:$PORT"
    echo "   終了: Ctrl+C"
    echo ""
    
    trap 'echo ""; log_info "サーバーを停止しています..."; exit 0' INT
    
    npx http-server -p $PORT -a 127.0.0.1
    
else
    log_error "適切なHTTPサーバーが見つかりません"
    echo ""
    echo "以下のいずれかをインストールしてください:"
    echo "  • Python 3: sudo apt install python3"
    echo "  • Python 2: sudo apt install python"
    echo "  • PHP: sudo apt install php"
    echo "  • Node.js: sudo apt install nodejs npm"
    echo ""
    echo "または、以下のコマンドで手動でサーバーを起動してください:"
    echo "  cd docs && python3 -m http.server $PORT"
    exit 1
fi