# infra/terraform

AWS インフラを Terraform で管理する。

## ディレクトリ

```
terraform/
├── versions.tf          required_providers (aws ~> 5.x)
├── modules/             再利用可能なリソースモジュール
│   ├── network/         VPC・サブネット・NAT・SG
│   ├── rds/             RDS for MySQL 8.0
│   ├── elasticache/     ElastiCache for Redis
│   ├── s3/              添付ファイル用バケット
│   ├── ecr/             コンテナレジストリ
│   ├── ecs/             ECS Fargate サービス（gateway/auth/ticket）
│   ├── alb/             ALB + listener rules
│   ├── cloudfront/      フロント配信
│   ├── iam/             ECS task role + GitHub OIDC
│   └── secrets/         Secrets Manager + SSM
└── envs/
    ├── dev/             開発環境
    ├── stg/             ステージング
    └── prod/            本番（Multi-AZ）
```

## 開発

M0 時点では **各モジュールは variables/outputs と空の resource ブロックのみ**。実リソース定義は後続マイルストーンで埋める。

```bash
make tf-fmt              # terraform fmt -recursive -check
make tf-validate-dev     # dev envs で terraform init + validate（リモート state なし）
```

実 apply とリソース作成は backend が縦通し完了（M6）した後の別タスク。
