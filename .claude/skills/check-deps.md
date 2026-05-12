---
description: Dependency CVE scan (govulncheck for Go, pnpm audit for Node). Run when go.sum or pnpm-lock.yaml has changed.
allowed-tools: Bash
---

# check-deps

依存関係の CVE スキャンを実行する。**毎回走らせると遅いため、依存追加・変更時のみ手動で起動する**想定。

## 手順

1. `cd backend && govulncheck ./...`
2. `cd frontend && pnpm audit --prod --audit-level=high`

成功なら "no CVE ✓"。

CI では `go.sum` / `pnpm-lock.yaml` 変更時のみ自動で走る（`.github/workflows/backend-deps.yml`、`.github/workflows/frontend-deps.yml`）。
