variable "aws_region" {
  default     = "ap-southeast-1"
  description = "Primary AWS region"
}

variable "replication_region" {
  default     = "us-east-1"
  description = "Destination region for S3 cross-region replication"
}

variable "environment" {
  default     = "dev"
  description = "Environment name (dev, staging, prod)"
}

variable "instance_type" {
  default     = "t3.micro"
  description = "EC2 instance type for dummy instance"
}