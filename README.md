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

# バックエンド（M0 時点ではビルドのみ通る空エントリ）
make backend-build

# フロントエンド
make front-install
make front-dev   # http://localhost:3000
```

主要な make ターゲットは `make help` で確認。

## マイルストーン

実装はマイルストーン単位で進める。

- **M0**（現在）：プロジェクト骨格
- M1：DB マイグレーション・sqlc・oapi-codegen 連携
- M2：横断基盤 + auth-service
- M3：API Gateway 最小
- M4：ticket-service コア
- M5：コメント & 添付（S3）
- M6：縦通し E2E 検証
