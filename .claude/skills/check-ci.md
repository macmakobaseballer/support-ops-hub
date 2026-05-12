---
description: GitHub Actions workflow quality + security checks (actionlint, zizmor). Use when .github/workflows/*.yml is modified.
allowed-tools: Bash
---

# check-ci

GitHub Actions ワークフロー（`.github/workflows/*.yml`）に対し、以下を順次実行して **失敗のみ簡潔に報告** する。

## 手順

1. `actionlint .github/workflows/*.yml`
2. `zizmor --min-severity=high .github/workflows/`

成功なら "ci OK ✓"。

`zizmor` の medium 以下の警告は情報レベル、CI では止めない（high 以上で fail）。
