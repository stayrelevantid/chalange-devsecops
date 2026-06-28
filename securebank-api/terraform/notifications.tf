# SNS Topic for S3 Event Notifications (CKV2_AWS_62)

resource "aws_sns_topic" "s3_events" {
  name              = "securebank-s3-events"
  kms_master_key_id = aws_kms_key.securebank.arn
}

resource "aws_sns_topic_policy" "s3_events" {
  arn = aws_sns_topic.s3_events.arn

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action   = "sns:Publish"
        Resource = aws_sns_topic.s3_events.arn
        Condition = {
          ArnLike = {
            "aws:SourceArn" = aws_s3_bucket.logs.arn
          }
        }
      }
    ]
  })
}