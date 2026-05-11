# Module: elasticache
# Implemented in a later milestone. Keep this file present so `terraform fmt -recursive` and module discovery work.

variable "env" {
  description = "Environment name (dev/stg/prod)."
  type        = string
}

variable "tags" {
  description = "Common tags to apply to all resources."
  type        = map(string)
  default     = {}
}
