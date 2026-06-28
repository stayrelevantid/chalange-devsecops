# Network Security — fixes SG + VPC findings

# Security Group — FIXED: restricted ports + descriptions + VPC-only CIDR
resource "aws_security_group" "api" {
  name        = "securebank-api-sg"
  description = "SecureBank API security group — restrict ingress to 443/8080 from VPC"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  ingress {
    description = "API port from VPC"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    description = "HTTPS outbound within VPC only"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  tags = {
    Name = "securebank-api-sg"
  }
}

# Default SG — restrict all traffic (CKV2_AWS_12)
resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.main.id

  # No ingress or egress rules = restrict all traffic
  tags = {
    Name = "securebank-default-sg-restricted"
  }
}

# VPC Flow Log (CKV2_AWS_11 / AWS-0178)
resource "aws_cloudwatch_log_group" "flow_log" {
  name              = "/aws/vpc/securebank-flow-log"
  retention_in_days = 365
  kms_key_id        = aws_kms_key.securebank.arn
}

resource "aws_flow_log" "main" {
  log_destination      = aws_cloudwatch_log_group.flow_log.arn
  log_destination_type = "cloud-watch-logs"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.main.id
  iam_role_arn         = aws_iam_role.flow_log.arn
}