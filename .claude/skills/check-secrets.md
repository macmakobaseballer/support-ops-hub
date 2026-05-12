---
description: Detect committed secrets (API keys, credentials) in the working tree using trivy.
allowed-tools: Bash
---

# check-secrets

リポジトリ全体に対し、trivy のシークレット検出を実行する。

## 手順

```
trivy fs --scanners secret --severity HIGH,CRITICAL --skip-dirs node_modules,.nuxt,.terraform .
```

ヒットがあれば **直ちに止めてユーザーに報告**。誤検知の場合は `.trivyignore` に追加し、再実行を案内する。

成功なら "no secrets detected ✓"。
