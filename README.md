# Support Ops Hub

問い合わせ管理システム（Support Ops Hub）。
運用・保守チームが顧客からの問い合わせを一元管理するための Web アプリケーション。

## 構成

| ディレクトリ | 内容 |
|---|---|
| [backend/](backend/) | Go バックエンド（Chi + sqlc + oapi-codegen / モジュラーモノリス） |
| [frontend/](frontend/) | Nuxt 3 フロントエンド |
| [infra/terraform/](infra/terraform/) | AWS インフラ（Terraform） |
| [docs/](docs/) | API・アーキ・テーブル定義・画面仕様 |
| [prototype/](prototype/) | 静的 HTML プロトタイプ（参考） |

仕様詳細は [docs/](docs/) と [CLAUDE.md](CLAUDE.md) を参照。

## ローカル起動

前提：Docker、Go 1.26+、Node 20.19+ または 22.12+（Nuxt 推奨は 24）、pnpm 10+、Terraform 1.10+、`make`（Ubuntu/WSL なら `sudo apt install -y make`）。

```bash
cp .env.example .env

# MySQL / LocalStack を起動
make up

# バックエンド（現時点ではビルドのみ通る空エントリ。M2 以降で実装）
make backend-build

# フロントエンド
make front-install
make front-dev   # http://localhost:3000
```

主要な make ターゲットは `make help` で確認。

## 品質チェック

PR を出す前に各レイヤーの lint / 型 / unit test / セキュリティをローカルで通す。Claude Code からは `/check-backend` などの skill で同じチェックを起動できる。

```bash
# 初回のみ：品質チェックツールを導入（~/.local/bin と ~/go/bin に入る）
make tools-install

# 日常（PR ゲート相当）
make check-backend     # golangci-lint v2 + go test -race -cover
make check-frontend    # ESLint + nuxt typecheck + vitest
make check-infra       # terraform fmt/validate + tflint + trivy config
make check-ci          # actionlint + zizmor (.github/workflows/)
make check-secrets     # trivy fs --scanners secret (リポ全体)
make check-all         # 上記 5 つを順次

# 依存追加・変更時のみ
make check-deps        # govulncheck + pnpm audit
```

CI は **3 階層** で動く：

| 階層 | ワークフロー | 走るタイミング |
|---|---|---|
| PR ゲート（毎 push） | `backend.yml` / `frontend.yml` / `terraform.yml` / `quality.yml` | path-filter で関連レイヤーのみ起動 |
| 依存変更時のみ | `backend-deps.yml` / `frontend-deps.yml` | `go.sum` / `pnpm-lock.yaml` の変更時のみ |
| 週次 + main マージ | `weekly-scan.yml` | スケジュール（毎週月曜 03:00 JST）と main へのマージで全 CVE / 全 secret を再検査 |

ツール選定とポリシーは [/home/mako/.claude/plans/api-reflective-whale.md](.) と [.github/zizmor.yml](.github/zizmor.yml) を参照。

## マイルストーン

実装はマイルストーン単位で進める。各マイルストーンの作業は Issue を起点に、`feature/<issue#>-<slug>` ブランチ → `develop` への PR 経由で取り込む（[CLAUDE.md](CLAUDE.md) 横断ルール 6）。

| M | 名称 | Backend | Frontend | Infra | CI/CD |
|---|---|---|---|---|---|
| ✅ M0 | プロジェクト骨格 | go.mod + 空 cmd | Nuxt 雛形 | tf 骨格 | path-filter のみ |
| M1 | Foundation | 8 テーブル migration、sqlc、oapi-codegen | openapi-typescript、layout、エラー表示、Pinia 基盤 | — | CI に `pnpm openapi:types` 追加 |
| M2 | 横断基盤 + dev-auth Gateway | `internal/{httperr,middleware,domain,auth/jwt,config}`、Gateway リバプロ + dev-auth ミドルウェア | 開発用「ユーザー切替」セレクタ、`composables/useApi.ts` | — | — |
| M3 | Tickets Core | ticket-service CRUD + ステータス遷移 + history | SCR-02/03/04 基本/05/06 | — | — |
| M4 | Comments & Attachments | `/comments`・`/attachments`、LocalStack S3 | SCR-04 強化 | — | — |
| M5 | Masters & Analytics | customer / system / user / analytics サービス | SCR-07/08/09/10/11/12、サイドバー | — | — |
| M6 | E2E & 仕上げ | バグ修正 | Playwright E2E、UI/UX 仕上げ | — | E2E を CI で実行 |
| M7 | 本物 Auth 切替 | auth-service 実装（bcrypt + JWT 発行）、Gateway middleware を JWT 検証に差替 | **SCR-01 ログイン画面**、cookie ベース保存に切替 | — | — |
| M8 | Infra 実装 & dev apply | — | — | state backend、各モジュール resource 定義、`envs/dev/` を apply、スモークテスト | terraform plan/apply ワークフロー（OIDC） |
| M9 | CI/CD 完成 | — | — | — | backend: ECR push → ECS deploy / frontend: build → S3 sync → CloudFront / terraform: plan on PR, apply on merge |

### Auth 戦略について

M2〜M6 は Gateway 側の **dev-auth ミドルウェア**（`X-Dev-User-Email` ヘッダーから user を引いて `X-User-ID`/`X-User-Role` を注入）で開発し、各サービスは [CLAUDE.md](CLAUDE.md) ルール 6 に沿ってヘッダーの出元を意識しない実装にする。M7 で Gateway middleware を本物 JWT 検証に差し替え、SCR-01 ログイン画面を追加する。サービスコードは差替時に一切変更しない。詳細は [docs/api/architecture.md](docs/api/architecture.md) §3.1 を参照。
