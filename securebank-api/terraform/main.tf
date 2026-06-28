# ⚠️ DO NOT RUN `terraform apply`
# This Terraform config is for STATIC ANALYSIS ONLY (Checkov + Trivy IaC)
# No AWS resources are created. No cost is incurred.
# To destroy (if accidentally applied): see terraform/DESTROY.md

terraform {
  required_version = ">= 1.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "securebank"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

# VPC — fixed: added flow logging + default SG restriction (in network.tf)
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "securebank-vpc"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# S3 Bucket — base resource (hardening in s3.tf)
resource "aws_s3_bucket" "logs" {
  bucket = "securebank-logs-${var.environment}"

  lifecycle {
    prevent_destroy = true
  }
}

# S3 Access Log destination bucket (for S3 access logging)
resource "aws_s3_bucket" "access_logs" {
  #checkov:skip=CKV2_AWS_62: Access log destination bucket does not need event notifications
  #checkov:skip=CKV_AWS_144: Access log destination does not need cross-region replication (it IS the log destination)
  # Note: Trivy AWS-0089 (LOW) accepted as risk — log destination bucket does not need its own access logging

  bucket   = "securebank-access-logs-${var.environment}"
  provider = aws

  lifecycle {
    prevent_destroy = true
  }
}

# S3 Replication destination bucket (different region)
# Note: Checkov cannot associate standalone resources with alias provider buckets.
# Hardening (encryption, versioning, lifecycle, public access block) is configured in s3.tf.
resource "aws_s3_bucket" "replication_dest" {
  #checkov:skip=CKV2_AWS_62: Replication destination does not need event notifications
  #checkov:skip=CKV_AWS_18: Replication destination does not need access logging (it IS the log destination)
  #checkov:skip=CKV2_AWS_61: Lifecycle configured in s3.tf
  #checkov:skip=CKV2_AWS_6: Public access block configured in s3.tf
  #checkov:skip=CKV_AWS_145: KMS encryption configured in s3.tf
  #checkov:skip=CKV_AWS_21: Versioning configured in s3.tf
  #checkov:skip=CKV_AWS_144: This IS the replication destination — no replication needed
  #checkov:skip=CKV_AWS_19: Encryption configured in s3.tf
  # Note: Trivy AWS-0089 (LOW) accepted as risk — replication destination does not need access logging

  bucket   = "securebank-replication-dest-${var.environment}"
  provider = aws.replication_region

  lifecycle {
    prevent_destroy = true
  }
}

# Alias provider for replication region
provider "aws" {
  alias  = "replication_region"
  region = var.replication_region
}