# Outputs — visibility into created resources (best practice)

output "vpc_id" {
  description = "ID of the SecureBank VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "CIDR block of the SecureBank VPC"
  value       = aws_vpc.main.cidr_block
}

output "logs_bucket_name" {
  description = "Name of the S3 logs bucket"
  value       = aws_s3_bucket.logs.id
}

output "logs_bucket_arn" {
  description = "ARN of the S3 logs bucket"
  value       = aws_s3_bucket.logs.arn
}

output "access_logs_bucket_name" {
  description = "Name of the S3 access logs bucket"
  value       = aws_s3_bucket.access_logs.id
}

output "replication_dest_bucket_arn" {
  description = "ARN of the S3 replication destination bucket"
  value       = aws_s3_bucket.replication_dest.arn
}

output "kms_key_arn" {
  description = "ARN of the SecureBank KMS key"
  value       = aws_kms_key.securebank.arn
}

output "kms_key_alias" {
  description = "Alias of the SecureBank KMS key"
  value       = aws_kms_alias.securebank.name
}

output "api_security_group_id" {
  description = "ID of the API security group"
  value       = aws_security_group.api.id
}

output "app_subnet_id" {
  description = "ID of the app subnet"
  value       = aws_subnet.app.id
}

output "flow_log_group_name" {
  description = "CloudWatch log group name for VPC flow logs"
  value       = aws_cloudwatch_log_group.flow_log.name
}

output "sns_topic_arn" {
  description = "ARN of the SNS topic for S3 event notifications"
  value       = aws_sns_topic.s3_events.arn
}