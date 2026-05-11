# ==========================================================================
# Support Ops Hub — dev environment entry point
#
# M0: provider のみ宣言し、モジュールは未呼び出し。
# 後続マイルストーンで modules/* のリソース定義を埋めながら、ここから呼び出す。
# ==========================================================================

terraform {
  required_version = ">= 1.10.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project     = "support-ops-hub"
      Environment = "dev"
      ManagedBy   = "terraform"
    }
  }
}

variable "region" {
  description = "AWS region for the dev environment."
  type        = string
  default     = "ap-northeast-1"
}
