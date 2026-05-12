# Support Ops Hub — CI / CLAUDE.md

GitHub Actions ワークフロー（`.github/workflows/`）で作業する際のルール。**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずユーザーに確認すること。

リポジトリ全体のルールは [../CLAUDE.md](../CLAUDE.md) を参照。

---

## ワークフロー一覧

| ファイル | 対象 | トリガパス | 現在のジョブ |
|---------|------|-----------|-----------|
| [workflows/backend.yml](workflows/backend.yml) | Go バックエンド | `backend/**` | `golangci-lint` / `go test -race -cover` |
| [workflows/frontend.yml](workflows/frontend.yml) | Nuxt フロント | `frontend/**` / `docs/api/openapi.yaml` | `pnpm lint` / `pnpm typecheck` / `pnpm test`（各ジョブで `pnpm openapi:types` 実行） |
| [workflows/terraform.yml](workflows/terraform.yml) | Terraform | `infra/**` | `terraform fmt` / `terraform validate`（dev）/ `tflint` / `trivy config` |
| [workflows/quality.yml](workflows/quality.yml) | CI品質・シークレット | 全 PR / `main,develop` への push | `actionlint` / `zizmor` / `trivy secret` |
| [workflows/backend-deps.yml](workflows/backend-deps.yml) | Go 依存脆弱性 | `backend/go.mod` / `backend/go.sum` | `govulncheck` |
| [workflows/frontend-deps.yml](workflows/frontend-deps.yml) | Node 依存脆弱性 | `frontend/package.json` / `frontend/pnpm-lock.yaml` | `pnpm audit --prod --audit-level=high` |
| [workflows/weekly-scan.yml](workflows/weekly-scan.yml) | 定期セキュリティ再検査 | 毎週月曜 03:00 JST / `main` push / 手動実行 | `trivy secret` / `trivy vuln` / `govulncheck` / `pnpm audit` |

---

## ルール

### ルール C1：paths フィルタを必ず付ける

各ワークフローは関係のない変更で動かないよう、`on.push.paths` と `on.pull_request.paths` で対象パスを限定する：

```yaml
on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'
      - '.github/workflows/backend.yml'
  pull_request:
    paths:
      - 'backend/**'
      - '.github/workflows/backend.yml'
```

ワークフロー自身（`.github/workflows/<name>.yml`）も paths に必ず含める（ワークフロー修正時に動かないと困る）。

frontend ワークフローは追加で [`../docs/api/openapi.yaml`](../docs/api/openapi.yaml) も対象にする（型再生成のため）。

### ルール C2：working-directory を明示する

各ジョブは対象ディレクトリ（`backend/` / `frontend/` / `infra/terraform/`）で動作するため、`defaults.run.working-directory` を必ず設定する。

### ルール C3：トリガブランチは main / develop のみ

- `push` トリガは `branches: [main, develop]` 固定
- それ以外のブランチは pull_request 経由でのみ CI を走らせる

### ルール C4：依存キャッシュを使う

| ツール | キャッシュ設定 |
|--------|------------|
| Go | `actions/setup-go@v5` の `cache-dependency-path: backend/go.sum` |
| Node | `actions/setup-node@v4` の `cache: 'pnpm'` + `cache-dependency-path: frontend/pnpm-lock.yaml` |
| Terraform | キャッシュなし（軽量のため） |

### ルール C5：シークレットは GitHub Secrets と OIDC で管理する

- 静的なシークレット（外部 API トークン等）は GitHub Secrets で管理
- AWS 認証は GitHub OIDC + IAM Role（[../infra/terraform/modules/iam/](../infra/terraform/modules/iam/)）を使用。長期 access key を使わない
- ワークフロー内で `echo ${{ secrets.* }}` しない（ログに漏れる）

### ルール C6：apply・push 系の操作は手動承認を挟む

- AWS への `terraform apply` や ECR への image push を将来追加する場合は、`environment:` で承認必須にする
- main ブランチへの直接 push 後の自動 deploy は禁止。明示的な承認後にのみ実行

### ルール C7：バージョン固定

以下は CI のバージョンを明示的に指定する：

| ツール | バージョン |
|--------|----------|
| Go | `1.26` |
| Node.js | `24` |
| pnpm | `10.13.1` |
| Terraform | `1.10.5` |

バージョン更新は別 PR で行い、影響範囲を分離する。固定値はリポジトリ全体で整合させること（[../Makefile](../Makefile)・[../frontend/package.json](../frontend/package.json) の `packageManager`・[../backend/go.mod](../backend/go.mod) 等）。

---

## 今後追加予定のジョブ

| ワークフロー | 追加予定ジョブ | マイルストーン |
|------------|-------------|-------------|
| backend.yml | sqlc・oapi-codegen の生成差分チェック | M2+ |
| frontend.yml | テスト結果のレポート出力（必要時） | M2+ |
| terraform.yml | stg / prod env の validate / 計画的な plan 出力 | M6+ |

追加時は本ドキュメントの「ワークフロー一覧」と「今後追加予定のジョブ」を更新する。

---

## ドキュメント更新義務

- 新しいワークフローを追加した場合は本ファイルの「ワークフロー一覧」を更新する
- ツールバージョンを更新した場合は本ファイルの「バージョン固定」を更新する
