# S3 Bucket Hardening — fixes 7 Checkov findings + 7 Trivy findings

# 1. Versioning (CKV_AWS_21 / AWS-0090)
resource "aws_s3_bucket_versioning" "logs" {
  bucket = aws_s3_bucket.logs.id

  versioning_configuration {
    status = "Enabled"
  }
}

# 2. KMS Encryption (CKV_AWS_145 / AWS-0132)
resource "aws_s3_bucket_server_side_encryption_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.securebank.arn
    }
  }
}

# Also encrypt access logs bucket with KMS
resource "aws_s3_bucket_server_side_encryption_configuration" "access_logs" {
  bucket = aws_s3_bucket.access_logs.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.securebank.arn
    }
  }
}

# Also encrypt replication destination bucket with KMS
resource "aws_s3_bucket_server_side_encryption_configuration" "replication_dest" {
  bucket = aws_s3_bucket.replication_dest.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.securebank.arn
    }
  }
}

# Versioning for access logs bucket
resource "aws_s3_bucket_versioning" "access_logs" {
  bucket = aws_s3_bucket.access_logs.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Versioning for replication destination (required for replication)
resource "aws_s3_bucket_versioning" "replication_dest" {
  bucket = aws_s3_bucket.replication_dest.id

  versioning_configuration {
    status = "Enabled"
  }
}

# 3. Public Access Block (CKV2_AWS_6 / AWS-0086, AWS-0087, AWS-0091, AWS-0093, AWS-0094)
resource "aws_s3_bucket_public_access_block" "logs" {
  bucket                  = aws_s3_bucket.logs.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "access_logs" {
  bucket                  = aws_s3_bucket.access_logs.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "replication_dest" {
  bucket                  = aws_s3_bucket.replication_dest.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# 4. Access Logging (CKV_AWS_18 / AWS-0089)
resource "aws_s3_bucket_logging" "logs" {
  bucket        = aws_s3_bucket.logs.id
  target_bucket = aws_s3_bucket.access_logs.id
  target_prefix = "log/securebank/"
}

# 5. Cross-Region Replication (CKV_AWS_144)
resource "aws_s3_bucket_replication_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id
  role   = aws_iam_role.replication.arn

  rule {
    id     = "replicate-all"
    status = "Enabled"

    destination {
      bucket        = aws_s3_bucket.replication_dest.arn
      storage_class = "STANDARD"
    }

    filter {
      prefix = ""
    }
  }

  depends_on = [aws_s3_bucket_versioning.logs]
}

# 6. Event Notifications (CKV2_AWS_62)
resource "aws_s3_bucket_notification" "logs" {
  bucket = aws_s3_bucket.logs.id

  topic {
    topic_arn     = aws_sns_topic.s3_events.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "log/"
  }

  depends_on = [aws_sns_topic_policy.s3_events]
}

# 7. Lifecycle Configuration (CKV2_AWS_61)
resource "aws_s3_bucket_lifecycle_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id

  rule {
    id     = "transition-to-glacier"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 90
      storage_class = "GLACIER"
    }

    expiration {
      days = 365
    }
  }
}

# Lifecycle for access logs bucket
resource "aws_s3_bucket_lifecycle_configuration" "access_logs" {
  bucket = aws_s3_bucket.access_logs.id

  rule {
    id     = "expire-access-logs"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }

    expiration {
      days = 90
    }
  }
}

# Lifecycle for replication destination bucket
resource "aws_s3_bucket_lifecycle_configuration" "replication_dest" {
  bucket = aws_s3_bucket.replication_dest.id

  rule {
    id     = "replication-dest-lifecycle"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }

    expiration {
      days = 365
    }
  }
}