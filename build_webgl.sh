#!/bin/bash

# WebGL/WASM ビルドスクリプト for EdomaeElf
# 参拝客システム付きゲームのWebGL版をビルドします

set -e

echo "🏮 EdomaeElf WebGL/WASM ビルドを開始します..."

# プロジェクトのルートディレクトリに移動
cd "$(dirname "$0")"

# webディレクトリを作成（存在しない場合）
mkdir -p web

# Go version確認
GO_VERSION=$(go version)
echo "Go version: $GO_VERSION"

# WebAssemblyサポートを確認
echo "WebAssembly support確認中..."
if ! go env GOOS | grep -q "js"; then
    echo "⚠️  WebAssemblyサポートが見つかりません。Go 1.11以上が必要です。"
fi

# wasm_exec.jsをコピー
WASM_EXEC_JS=$(go env GOROOT)/misc/wasm/wasm_exec.js
WASM_EXEC_JS_ALT=$(go env GOROOT)/lib/wasm/wasm_exec.js

if [ -f "$WASM_EXEC_JS" ]; then
    echo "📋 wasm_exec.js をコピー中... (misc/wasm/)"
    cp "$WASM_EXEC_JS" web/
elif [ -f "$WASM_EXEC_JS_ALT" ]; then
    echo "📋 wasm_exec.js をコピー中... (lib/wasm/)"
    cp "$WASM_EXEC_JS_ALT" web/
else
    echo "❌ wasm_exec.js が見つかりません:"
    echo "  $WASM_EXEC_JS"
    echo "  $WASM_EXEC_JS_ALT"
    echo "手動でwasm_exec.jsをwebディレクトリにコピーしてください。"
    exit 1
fi

# assetsディレクトリをwebディレクトリにコピー
echo "🎨 アセットをコピー中..."
cp -r assets web/

# WebGL/WASM版をビルド
echo "🔧 WebGL/WASM版をビルド中..."
GOOS=js GOARCH=wasm go build -o web/main.wasm miko_game_with_worshippers.go

# ビルドの成功を確認
if [ -f "web/main.wasm" ]; then
    echo "✅ WebGL/WASMビルドが完了しました！"
    echo ""
    echo "📁 生成されたファイル:"
    echo "   - web/main.wasm      (メインゲーム)"
    echo "   - web/index.html     (HTMLテンプレート)"
    echo "   - web/wasm_exec.js   (WebAssemblyランタイム)"
    echo "   - web/assets/        (ゲームアセット)"
    echo ""
    echo "🚀 実行方法:"
    echo "   1. Webサーバーを起動:"
    echo "      cd web && python3 -m http.server 8080"
    echo "   2. ブラウザで以下にアクセス:"
    echo "      http://localhost:8080"
    echo ""
    echo "⚠️  注意: WebAssemblyファイルは必ずHTTPサーバー経由で開いてください。"
    echo "   file://でのアクセスは動作しません。"
else
    echo "❌ ビルドに失敗しました。"
    exit 1
fi

echo "🎉 WebGL対応完了！"