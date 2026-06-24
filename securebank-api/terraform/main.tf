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
}

# VPC — intentionally misconfigured for training
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name        = "securebank-vpc"
    Environment = var.environment
  }
}

# S3 Bucket — intentionally insecure for training
resource "aws_s3_bucket" "logs" {
  bucket = "securebank-logs-${var.environment}"
  # Missing: versioning, encryption, public access block
}

# Security Group — intentionally too open
resource "aws_security_group" "api" {
  name   = "securebank-api-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # INTENTIONAL: open to world
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}