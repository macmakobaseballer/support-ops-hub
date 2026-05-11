# Support Ops Hub — stg environment entry point (M0 skeleton).

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
      Environment = "stg"
      ManagedBy   = "terraform"
    }
  }
}

variable "region" {
  description = "AWS region for the stg environment."
  type        = string
  default     = "ap-northeast-1"
}
