# Module: ecs
#
# Implemented in a later milestone. Keep this file present so `terraform fmt -recursive`
# and module discovery work.
#
# 設計方針（コスト最適化リビジョン）:
#   - launch type: EC2 (NOT Fargate)
#   - ECS cluster + capacity provider + Auto Scaling Group で EC2 を管理
#   - AMI: Amazon ECS-optimized Amazon Linux 2023 (ARM/Graviton)
#   - インスタンス: dev/stg は t4g.small × 1、prod は t4g.medium × 2 (Multi-AZ)
#   - Phase 1 は単一バイナリ・1 タスク・1 コンテナ構成
#     (../../../../backend/CLAUDE.md および ../../README.md を参照)
#   - bridge network、ALB target group へは dynamic port で登録

variable "env" {
  description = "Environment name (dev/stg/prod)."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to all resources."
  type        = map(string)
  default     = {}
}
