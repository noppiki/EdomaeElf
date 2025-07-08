#!/bin/bash

# WebGL版EdomaeElfサーバー起動スクリプト
# 参拝客システム付きゲームをブラウザで実行するためのローカルサーバー

set -e

echo "🏮 EdomaeElf WebGL版サーバーを起動します..."

# プロジェクトのルートディレクトリに移動
cd "$(dirname "$0")"

# webディレクトリの存在確認
if [ ! -d "web" ]; then
    echo "❌ webディレクトリが見つかりません。"
    echo "まずビルドスクリプトを実行してください:"
    echo "  ./build_webgl.sh"
    exit 1
fi

# main.wasmファイルの存在確認
if [ ! -f "web/main.wasm" ]; then
    echo "❌ main.wasmファイルが見つかりません。"
    echo "まずビルドスクリプトを実行してください:"
    echo "  ./build_webgl.sh"
    exit 1
fi

# webディレクトリに移動
cd web

# 利用可能なポートを確認
PORT=8080
if command -v lsof >/dev/null 2>&1; then
    if lsof -i :$PORT >/dev/null 2>&1; then
        echo "⚠️  ポート8080は既に使用中です。ポート8081を使用します。"
        PORT=8081
    fi
fi

echo "🌐 ローカルサーバーを起動中..."
echo "📍 URL: http://localhost:$PORT"
echo "🎮 ゲームを開始するには、上記URLをブラウザで開いてください。"
echo ""
echo "⚠️  注意事項:"
echo "   - サーバーを停止するには Ctrl+C を押してください"
echo "   - ゲームが重い場合は、ブラウザのタブを1つだけ開いてください"
echo "   - 開発者ツール(F12)でコンソールを確認できます"
echo ""
echo "🚀 サーバー起動中..."

# Python3が利用できるかチェック
if command -v python3 >/dev/null 2>&1; then
    echo "Python3サーバーを使用"
    python3 -m http.server $PORT
elif command -v python >/dev/null 2>&1; then
    echo "Python2サーバーを使用"
    python -m SimpleHTTPServer $PORT
elif command -v php >/dev/null 2>&1; then
    echo "PHPサーバーを使用"
    php -S localhost:$PORT
else
    echo "❌ Python3, Python2, PHPのいずれも見つかりません。"
    echo "以下のいずれかをインストールしてください:"
    echo "  - Python3: python3 -m http.server $PORT"
    echo "  - Python2: python -m SimpleHTTPServer $PORT"
    echo "  - PHP: php -S localhost:$PORT"
    echo ""
    echo "または、お好みのWebサーバーでwebディレクトリを配信してください。"
    exit 1
fi