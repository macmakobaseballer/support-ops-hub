# Module: network
#
# Implemented in a later milestone. Keep this file present so `terraform fmt -recursive`
# and module discovery work.
#
# 設計方針（コスト最適化リビジョン）:
#   - VPC: /16、AZ × 2（prod のみ実質 2 AZ 利用、dev/stg は 2 AZ 用意するが片方は休眠）
#   - public subnet  × 2  ← ALB と アプリ EC2 を配置
#   - private subnet × 2  ← RDS のみ配置（DB Subnet Group）
#   - **NAT Gateway / NAT インスタンスは作らない**
#   - アプリ EC2 は public subnet、Security Group で ALB からの受信のみ許可、
#     アウトバウンドは Internet Gateway 直通
#   - S3 へのアクセスは Gateway 型 VPC Endpoint（無料）を 1 つ作成して egress を抑制

variable "env" {
  description = "Environment name (dev/stg/prod)."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to all resources."
  type        = map(string)
  default     = {}
}
