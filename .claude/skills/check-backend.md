---
description: Backend Go quality checks (lint v2, unit test with race+cover). Use this when changes are made under backend/ or before opening a PR.
allowed-tools: Bash
---

# check-backend

`backend/` 配下の Go コードに対し、以下を順次実行して **失敗のみ簡潔に報告** する。

## 手順

1. `cd backend && golangci-lint run ./...`
2. `cd backend && go test -race -cover ./...`

すべて成功なら "backend OK ✓" を表示。途中で失敗してもすべて実行してから集約レポート。

ツール未インストールなら明示する（`make tools-install` の案内）。

依存関係（CVE）チェックは別 skill `/check-deps` で実行（毎回走らせると重いため）。
