# Module: rds
#
# Implemented in a later milestone. Keep this file present so `terraform fmt -recursive`
# and module discovery work.
#
# 設計方針（コスト最適化リビジョン）:
#   - engine: MySQL 8.0（変更なし。バージョンはコスト要因ではない）
#   - dev/stg: db.t4g.micro / gp3 20GB / Single-AZ
#   - prod:    db.t4g.small / gp3 50GB / Multi-AZ
#   - バックアップ保持 7 日
#   - Performance Insights / Enhanced Monitoring は prod のみ最低限
#   - parameter group は utf8mb4_0900_ai_ci、time_zone = Asia/Tokyo

variable "env" {
  description = "Environment name (dev/stg/prod)."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to all resources."
  type        = map(string)
  default     = {}
}
