# Support Ops Hub — CLAUDE.md

AI エージェントがこのプロジェクトで作業する際のルートルール。**詳細ルールは各領域の `CLAUDE.md` に分割している。** 該当領域で作業する際は、必ずそのディレクトリの `CLAUDE.md` を読むこと。

---

## プロジェクト概要

問い合わせ管理システム（Support Ops Hub）。
運用・保守チームが顧客からの問い合わせを一元管理する Web アプリケーション。

| 項目 | 内容 |
|------|------|
| フロントエンド | Nuxt 3 + TypeScript + Pinia + Tailwind |
| バックエンド | Go（モジュラーモノリス） |
| データベース | MySQL 8.0（AWS RDS） |
| インフラ | AWS（ECS on EC2 / RDS MySQL 8.0 / S3 / CloudFront、Redis 不使用）、Terraform |
| CI | GitHub Actions |
| 要件定義 | [docs/requirements/要件定義書.md](docs/requirements/要件定義書.md) |
| API 仕様 | [docs/api/openapi.yaml](docs/api/openapi.yaml) |
| アーキテクチャ設計 | [docs/api/architecture.md](docs/api/architecture.md) |
| 画面仕様 | [docs/screens/画面仕様書.md](docs/screens/画面仕様書.md) |
| 画面遷移 | [docs/screens/画面遷移図.md](docs/screens/画面遷移図.md) |
| データ定義 | [docs/data/データ定義.md](docs/data/データ定義.md) |
| ER 図 | [docs/data/ER図.md](docs/data/ER図.md) |
| テーブル定義書 | [docs/data/テーブル定義書.md](docs/data/テーブル定義書.md) |

---

## ルールドキュメント索引

| 領域 | ルールドキュメント | 主な内容 |
|------|----------------|---------|
| バックエンド | [backend/CLAUDE.md](backend/CLAUDE.md) | マイクロサービス・アーキテクチャルール（ルール1〜8） / API 設計規約 / Go コード規約 / 実装上の注意事項 |
| フロントエンド | [frontend/CLAUDE.md](frontend/CLAUDE.md) | Nuxt 3 + Pinia 構成 / 画面と SCR の対応 / API クライアント集約 / 状態管理 |
| インフラ | [infra/CLAUDE.md](infra/CLAUDE.md) | Terraform モジュール構成 / 環境分離 / 命名・タグ規約 / シークレット管理 |
| CI | [.github/CLAUDE.md](.github/CLAUDE.md) | ワークフロー構成 / paths フィルタ / OIDC・Secrets 管理 |
| API ユニットテスト | [docs/testing/api-unit-test.md](docs/testing/api-unit-test.md) | Go バックエンドのテスト規約 |
| フロント ユニットテスト | [docs/testing/frontend-unit-test.md](docs/testing/frontend-unit-test.md) | Nuxt フロントのテスト規約 |
| E2E テスト | 未策定 | 後続マイルストーンで策定 |

Claude Code はディレクトリ配下の `CLAUDE.md` を自動でロードするが、ルートで作業する場合や他領域に影響する変更を行う場合は、関連する `CLAUDE.md` を明示的に参照すること。

---

## 横断ルール

以下はどの領域で作業する場合も適用する。

### 横断ルール1：ドキュメント駆動

仕様変更は **必ずドキュメントを先に更新**してからコードを書く。ドキュメントとコードの乖離はバグ扱い。

| 変更内容 | 先に更新するドキュメント |
|---------|-------------------|
| エンドポイントの追加・変更 | [docs/api/openapi.yaml](docs/api/openapi.yaml) |
| 画面の追加・変更 | [docs/screens/画面仕様書.md](docs/screens/画面仕様書.md) |
| データ構造の変更 | [docs/data/データ定義.md](docs/data/データ定義.md) / [docs/data/ER図.md](docs/data/ER図.md) / [docs/data/テーブル定義書.md](docs/data/テーブル定義書.md) |
| サービス境界の変更 | [docs/api/architecture.md](docs/api/architecture.md) |

### 横断ルール2：OpenAPI 変更時の波及

[docs/api/openapi.yaml](docs/api/openapi.yaml) を更新したら、以下も連動して更新する：

- backend：oapi-codegen で `internal/apigen/` を再生成し、コミット
- frontend：`pnpm openapi:types` で `types/api.d.ts` を再生成（gitignore のため CI で自動生成）
- backend / frontend 双方のテストが新スキーマで通ることを確認

### 横断ルール3：新サービス追加手順

新しいバックエンドサービスを追加する場合、以下を **すべて** 更新する：

1. [docs/api/architecture.md](docs/api/architecture.md) のサービスマップ
2. [docs/api/openapi.yaml](docs/api/openapi.yaml) に新タグを追加
3. [backend/CLAUDE.md](backend/CLAUDE.md) の「所有テーブル」表
4. [backend/](backend/) 配下に `cmd/<service>/main.go` を作成
5. Phase 1：既存の単一バイナリ・ECS タスクに新サービスのモジュールをリンクする（[backend/CLAUDE.md](backend/CLAUDE.md) 構成参照）。Phase 2 で ECS タスクを分離する場合は [infra/terraform/envs/*/main.tf](infra/terraform/envs/) で配線

### 横断ルール4：エラー形式の整合

バックエンドが返すエラーレスポンスとフロントエンドの解釈を一致させる：

```json
{
  "code": "TICKET_NOT_FOUND",
  "message": "指定されたチケットは存在しません",
  "details": {}
}
```

- `code` は SCREAMING_SNAKE_CASE
- `message` は日本語ユーザー向けメッセージ
- 詳細は [backend/CLAUDE.md](backend/CLAUDE.md) ルール4 と [frontend/CLAUDE.md](frontend/CLAUDE.md) ルール F3

### 横断ルール5：ハードコンストレイント逸脱は事前確認

各領域の `CLAUDE.md` に書かれたルールは **提案ではなくハードコンストレイント**。逸脱する場合は、コードを書く前に必ずユーザーに確認すること。

### 横断ルール6：開発フローは Issue 起点で必ず PR 経由

**すべての変更は Issue を起点に、feature ブランチ → develop への PR を経由してマージする。**

| ステップ | 内容 |
|---|---|
| 1. Issue を切る | `gh issue create` で起票。タイトルは変更の意図を簡潔に。本文に背景・スコープ・受け入れ条件を書く |
| 2. feature ブランチを切る | `develop` から `feature/<issue#>-<slug>` で分岐。例：`feature/12-m1-migrations`、`feature/15-fix-typo` |
| 3. 実装・コミット | コミットメッセージは日本語で内容を簡潔に。Claude Code が編集した場合は `Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>` トレーラーを付与 |
| 4. PR を作る | base: `develop`、head: `feature/<issue#>-...`。本文に **`Closes #<issue#>`** を必ず含める（マージで Issue が自動クローズされる） |
| 5. CI 確認 → マージ | 全ワークフロー（backend / frontend / terraform）が pass したら **merge commit** でマージ（squash / rebase はデフォルト不可） |
| 6. Issue クローズ確認 | 自動クローズされなかった場合は手動で `gh issue close <issue#>` |

リリース（develop → main）は別フロー：Issue 不要、`gh pr create --base main --head develop` で PR を作り、merge commit でマージする。

**例外：** 本ルール導入のシード commit と、緊急障害対応のみ `develop` への直 push を許可する。後者の場合は事後に Issue を起票して履歴を残すこと。`main` への直 push は禁止。

---

## 開発の進め方

実装は段階的に進める。マイルストーンと現在の進捗は [README.md](README.md) と各領域の `README.md` を参照。

ローカル開発の起動は [Makefile](Makefile) のターゲットを使う：

```bash
make up                # MySQL / LocalStack を起動
make backend-build     # backend をビルド
make front-dev         # frontend を起動（http://localhost:3000）
make tf-validate-dev   # terraform 構文チェック
```
