#!/bin/bash
# ==============================================================
# revoke.sh — Auto-Revoke Leaked AWS Credentials
# Day 53: Red Team — Leaked Credentials
#
# Skrip untuk merespons access key IAM yang bocor:
#   1. Deactivate access key (revoke cepat)
#   2. Attach Deny-All inline policy (isolasi user)
#   3. (Optional) List API calls terakhir dari key tsb via CloudTrail
#
# Usage:
#   ./revoke.sh <IAM_USERNAME> <ACCESS_KEY_ID>
#   ./revoke.sh leaked-dev AKIAZ6PEJSWXU5QVF5KY
# ==============================================================
set -euo pipefail

TARGET_USER="${1:?Usage: $0 <IAM_USERNAME> <ACCESS_KEY_ID>}"
ACCESS_KEY_ID="${2:?Usage: $0 <IAM_USERNAME> <ACCESS_KEY_ID>}"
REGION="${AWS_DEFAULT_REGION:-ap-southeast-1}"
DENY_POLICY_NAME="BlockAll-On-Leak"

echo "=== [1/3] CloudTrail — riwayat aktivitas key ${ACCESS_KEY_ID} ==="
aws cloudtrail lookup-events \
  --lookup-attributes AttributeKey=AccessKeyId,AttributeValue="${ACCESS_KEY_ID}" \
  --region "${REGION}" \
  --query 'Events[].[EventTime,EventName,EventSource]' \
  --output table 2>&1 || echo "  (CloudTrail events mungkin belum tersedia)"

echo ""
echo "=== [2/3] Deactivate access key ${ACCESS_KEY_ID} ==="
aws iam update-access-key \
  --user-name "${TARGET_USER}" \
  --access-key-id "${ACCESS_KEY_ID}" \
  --status Inactive
echo "  -> Access key ${ACCESS_KEY_ID} DEACTIVATED"

echo ""
echo "=== [3/3] Attach Deny-All policy ke user ${TARGET_USER} ==="
aws iam put-user-policy \
  --user-name "${TARGET_USER}" \
  --policy-name "${DENY_POLICY_NAME}" \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Deny",
      "Action": "*",
      "Resource": "*"
    }]
  }'
echo "  -> Deny-All policy '${DENY_POLICY_NAME}' attached"

echo ""
echo "=== DONE: User ${TARGET_USER} telah di-revoke. ==="
echo "Verify: aws iam get-access-key-last-used --access-key-id ${ACCESS_KEY_ID}"
