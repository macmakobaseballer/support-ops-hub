# Support Ops Hub — prod environment entry point (M0 skeleton).
# 本番は Multi-AZ・autoscaling を有効化する想定（後続マイルストーンで反映）。

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
      Environment = "prod"
      ManagedBy   = "terraform"
    }
  }
}

variable "region" {
  description = "AWS region for the prod environment."
  type        = string
  default     = "ap-northeast-1"
}
