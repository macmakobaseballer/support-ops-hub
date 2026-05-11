# Support Ops Hub — infra / CLAUDE.md

AWS インフラ（Terraform）で作業する際のルール。**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずユーザーに確認すること。

リポジトリ全体のルールは [../CLAUDE.md](../CLAUDE.md) を参照。

---

## 構成

| 項目 | 内容 |
|------|------|
| ツール | Terraform 1.10.0 以上（CI では 1.10.5 固定） |
| プロバイダ | AWS provider `~> 5.0` |
| 構成パターン | modules / envs 分離（再利用モジュール + 環境別エントリポイント） |

ディレクトリ構成は [terraform/README.md](terraform/README.md) を参照。

---

## アーキテクチャ概要

3 環境（dev / stg / prod）にそれぞれ同じモジュールをデプロイする。**prod のみ Multi-AZ** とする。

```
envs/{dev,stg,prod}/main.tf
  ↓ 呼び出し
modules/
  network        VPC・サブネット・NAT・SG
  rds            RDS for MySQL 8.0
  elasticache    ElastiCache for Redis
  s3             添付ファイル用バケット
  ecr            コンテナレジストリ
  ecs            ECS Fargate
  alb            ALB + listener rules
  cloudfront     フロント配信
  iam            ECS task role + GitHub OIDC role
  secrets        Secrets Manager / SSM
```

**Phase 1 のデプロイ：** バックエンドは単一バイナリにリンクされ、1 つの ECS タスク・1 コンテナで起動する（[../backend/CLAUDE.md](../backend/CLAUDE.md) 構成参照）。`modules/ecs` は Phase 1 では task definition を 1 つだけ作成する。Phase 2 でサービスごとに task definition を分離する際に拡張する。

---

## ルール

### ルール I1：環境差分は envs/ 配下にだけ持つ

- リソース定義は必ず `modules/` 配下に作る
- 環境固有の値（インスタンスサイズ、Multi-AZ 設定、ドメイン名等）は `envs/{env}/main.tf` で module 呼び出しの引数として渡す
- `modules/` 配下に環境別の条件分岐を埋め込まない（`var.environment == "prod"` のような分岐は禁止）

### ルール I2：命名・タグ規約

全リソースに以下のタグを必ず付与する。**実装は `provider "aws" { default_tags { ... } }` で一括設定する**（[terraform/envs/dev/main.tf](terraform/envs/dev/main.tf) を参照）。リソースごとに個別 `tags` を書かない。

```hcl
provider "aws" {
  default_tags {
    tags = {
      Project     = "support-ops-hub"
      Environment = "dev"  # / "stg" / "prod"
      ManagedBy   = "terraform"
    }
  }
}
```

リソース名のプレフィックス：`{project}-{env}-{component}-...`
例：`support-ops-hub-dev-ecs-gateway`

### ルール I3：シークレットは tfvars にも tf にも書かない

- DB パスワード・JWT 秘密鍵・外部 API キー等は **Secrets Manager / SSM Parameter Store** に保存する
- Terraform は「Secrets Manager のリソースを作成し、`ignore_changes = [secret_string]` で値の変更を無視する」だけにとどめる
- ローカル動作確認のため `secret_string` の初期値が必要な場合はダミー値を入れ、本物は手動 or 別の経路でセットする

### ルール I4：ステート管理

- ステートファイルは S3 + DynamoDB ロックで管理する（実装は M6 以降）
- 現時点（M0 〜 M5）では `terraform init -backend=false` で validate のみ実施。実 apply はステートバックエンド整備後
- バックエンド設定は `envs/{env}/backend.tf` に分離する

### ルール I5：apply の権限と経路

- 開発者ローカルから `terraform apply` しない
- CI（GitHub Actions）+ GitHub OIDC role（`modules/iam` で定義）経由で apply する
- ローカルでは `terraform plan` で差分確認まで

### ルール I6：destroy・破壊的変更は要承認

- `terraform destroy` は禁止（環境ごと作り直す場合のみ別途承認の上で実施）
- リソースの置換が発生する変更（`-/+` の差分）は plan で確認し、PR に明記する
- データを持つリソース（RDS / S3 / ElastiCache）の削除や置換は必ず事前に承認を得る

---

## モジュール作成規約

### ファイル構成

各モジュールは以下のファイル構成を基本とする：

```
modules/<name>/
├── main.tf        ← リソース定義
├── variables.tf   ← 入力変数
├── outputs.tf     ← 出力値
└── versions.tf    ← required_providers（必要な場合）
```

### 変数・出力の規約

- すべての `variable` に `type` と `description` を必須にする
- すべての `output` に `description` を必須にする
- 機密値の output は `sensitive = true` を付ける

### モジュール間の依存

- モジュール間の参照は **envs 層で配線する**（モジュール内で他モジュールを呼び出さない）
- 例：`module.ecs` は `module.network.vpc_id` を引数で受け取る。`module.network` を内部で呼び出さない

---

## CI

- `terraform fmt` / `terraform validate` は [`.github/workflows/terraform.yml`](../.github/workflows/terraform.yml) で実行される
- 詳細は [../.github/CLAUDE.md](../.github/CLAUDE.md) を参照

---

## マイルストーン

M0 時点では各モジュールは骨格のみ（variables/outputs と空 resource）。実リソース定義は backend 縦通し（M6）後に整備する。
