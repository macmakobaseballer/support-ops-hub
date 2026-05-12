---
description: Frontend Nuxt quality checks (ESLint, typecheck, vitest). Use this when changes are made under frontend/ or before opening a PR.
allowed-tools: Bash
---

# check-frontend

`frontend/` 配下に対し、以下を順次実行して **失敗のみ簡潔に報告** する。

## 手順

1. `cd frontend && pnpm lint`
2. `cd frontend && pnpm typecheck`
3. `cd frontend && pnpm test`

すべて成功なら "frontend OK ✓"。途中で失敗してもすべて実行して集約レポート。

依存関係 audit は別 skill `/check-deps` で実行（毎回走らせると遅いため）。
