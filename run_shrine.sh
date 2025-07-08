#!/bin/bash
# 神社マップを実行するスクリプト
# IMK警告を抑制するための環境変数を設定

export GODEBUG=asyncpreemptoff=1
go run shrine_map.go