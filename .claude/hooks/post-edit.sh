#!/bin/bash
# PostToolUse hook: Edit/Write 後に自動チェックを実行する

input=$(cat)

# ファイルパスを抽出
file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty' 2>/dev/null)
if [ -z "$file_path" ]; then
    file_path=$(echo "$input" | grep -o '"file_path":"[^"]*"' | head -1 | sed 's/"file_path":"//;s/"//')
fi

[ -z "$file_path" ] && exit 0

root="$(git -C "$(dirname "$0")" rev-parse --show-toplevel 2>/dev/null || echo "/home/makoj/dev/support-ops-hub")"

# openapi.yaml 変更 → oapi-codegen 再生成
if [[ "$file_path" == *"openapi.yaml" ]]; then
    echo "openapi.yaml changed — running make oapi-codegen..."
    make -C "$root" oapi-codegen && echo "oapi-codegen: done"
fi

# backend の .go ファイル変更 → go vet
if [[ "$file_path" == *.go ]] && [[ "$file_path" == *"backend/"* ]]; then
    echo "Running go vet..."
    cd "$root/backend" && go vet ./... && echo "go vet: passed"
fi
