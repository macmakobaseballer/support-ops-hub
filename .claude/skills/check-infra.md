---
description: Terraform infrastructure quality + security checks (fmt, validate, tflint, trivy config). Use when changes are made under infra/.
allowed-tools: Bash
---

# check-infra

`infra/terraform/` に対し、以下を順次実行して **失敗のみ簡潔に報告** する。

## 手順

1. `cd infra/terraform && terraform fmt -recursive -check`
2. `cd infra/terraform/envs/dev && terraform init -backend=false -input=false && terraform validate`
3. `cd infra/terraform && tflint --init && tflint --recursive`
4. `trivy config infra/terraform/ --severity HIGH,CRITICAL`

成功なら "infra OK ✓"。

tflint の Warning（`unused_declarations` 等）は exit 0 のまま通る。HIGH/CRITICAL の trivy 検出は必ず止める。
