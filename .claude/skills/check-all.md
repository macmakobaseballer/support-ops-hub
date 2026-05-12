---
description: Run all PR-gate quality checks (backend, frontend, infra, ci, secrets) sequentially. Use before opening a PR.
allowed-tools: Bash, Skill
---

# check-all

PR 提出前に必須の全チェックを順次実行する。CVE スキャン（依存変更時のみ）は含まない。

## 手順

以下の skill を順番に呼び出し、**各結果サマリを集約**して報告：

1. `check-backend`
2. `check-frontend`
3. `check-infra`
4. `check-ci`
5. `check-secrets`

すべて成功なら "All PR-gate quality checks passed ✓" を表示。

依存関係（go.sum / pnpm-lock.yaml）を変更した場合は別途 `/check-deps` を実行する。
