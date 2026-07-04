variable "aws_region" {
  description = "Primary AWS region"
  type        = string
  default     = "ap-southeast-1"

  validation {
    condition     = can(regex("^(us|ap|eu|ca|me|sa)-(east|west|north|south|central|southeast|northeast)-[0-9]+$", var.aws_region))
    error_message = "aws_region must be a valid AWS region format (e.g., ap-southeast-1)."
  }
}

variable "replication_region" {
  description = "Destination region for S3 cross-region replication"
  type        = string
  default     = "us-east-1"

  validation {
    condition     = can(regex("^(us|ap|eu|ca|me|sa)-(east|west|north|south|central|southeast|northeast)-[0-9]+$", var.replication_region))
    error_message = "replication_region must be a valid AWS region format (e.g., us-east-1)."
  }
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "environment must be one of: dev, staging, prod."
  }
}

variable "instance_type" {
  description = "EC2 instance type for dummy instance"
  type        = string
  default     = "t3.micro"

  validation {
    condition     = can(regex("^(t3|t2|m5|m6|c5|c6|r5|r6)\\.(nano|micro|small|medium|large|xlarge|2xlarge|4xlarge)$", var.instance_type))
    error_message = "instance_type must be a valid AWS instance type (e.g., t3.micro, m5.large)."
  }
}

variable "kms_deletion_window" {
  description = "KMS key deletion window in days (7-30, use 30 for production safety)"
  type        = number
  default     = 30

  validation {
    condition     = var.kms_deletion_window >= 7 && var.kms_deletion_window <= 30
    error_message = "kms_deletion_window must be between 7 and 30 days (AWS limit)."
  }
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 365

  validation {
    condition     = contains([1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653], var.log_retention_days)
    error_message = "log_retention_days must be a valid CloudWatch retention period."
  }
}