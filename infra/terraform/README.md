# infra/terraform

AWS インフラを Terraform で管理する。**コスト最適化リビジョン**を適用。

## 設計方針サマリ

| 観点 | 採用 | 不採用にしたもの |
|---|---|---|
| コンテナ実行 | **ECS on EC2** (capacity provider + ASG) | Fargate |
| KVS | **使わない**（JWT ステートレス + in-process キャッシュ） | ElastiCache for Redis |
| RDB | **RDS for MySQL 8.0**（db.t4g.micro〜small） | Aurora / Serverless |
| ネットワーク | **NAT なし**、アプリ EC2 は public subnet | NAT Gateway / NAT インスタンス |
| S3 egress | **Gateway 型 VPC Endpoint**（無料） | NAT 経由 |

詳細・トレードオフは [`docs/api/architecture.md`](../../docs/api/architecture.md) を参照。

## ディレクトリ

```
terraform/
├── versions.tf          required_providers (aws ~> 5.x)
├── modules/             再利用可能なリソースモジュール
│   ├── network/         VPC・subnet (public/private)・SG・S3 VPC endpoint
│   ├── rds/             RDS for MySQL 8.0
│   ├── s3/              添付ファイル用バケット
│   ├── ecr/             コンテナレジストリ
│   ├── ecs/             ECS on EC2（cluster + ASG + capacity provider）
│   ├── alb/             ALB + listener rules
│   ├── cloudfront/      フロント配信
│   ├── iam/             ECS task role + GitHub OIDC
│   └── secrets/         Secrets Manager + SSM
└── envs/
    ├── dev/             1 EC2、Single-AZ RDS
    ├── stg/             1 EC2、Single-AZ RDS
    └── prod/            2 EC2 (Multi-AZ)、Multi-AZ RDS
```

## 開発

M0 〜 M7 時点では **各モジュールは variables/outputs と方針コメントのみ**。実リソース定義は M8（Infra 実装 & dev apply）で埋める。

```bash
make tf-fmt              # terraform fmt -recursive -check
make tf-validate-dev     # dev envs で terraform init + validate（リモート state なし）
```

実 apply とリソース作成は backend / frontend が E2E まで通った後の **M8** で実施する。
