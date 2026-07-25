#!/usr/bin/env bash
set -euo pipefail

DEFECTDOJO_URL="${DEFECTDOJO_URL:-http://localhost:8088}"
DEFECTDOJO_USER="${DEFECTDOJO_USER:-admin}"
DEFECTDOJO_PASS="${DEFECTDOJO_PASS:-DefectDojo60DaysChallenge!2026}"
PRODUCT_NAME="${PRODUCT_NAME:-SecureBank API}"
ENGAGEMENT_NAME="${ENGAGEMENT_NAME:-Q3 Security Audit}"
SECURITY_DIR="${SECURITY_DIR:-$(dirname "$0")/../}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "=== DefectDojo Upload Script ==="
echo "URL: $DEFECTDOJO_URL"
echo "Product: $PRODUCT_NAME"
echo "Engagement: $ENGAGEMENT_NAME"
echo ""

# Get API token
echo "Getting API token..."
TOKEN=$(curl -s -X POST "$DEFECTDOJO_URL/api/v2/api-token-auth/" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$DEFECTDOJO_USER\",\"password\":\"$DEFECTDOJO_PASS\"}" \
  | python3 -c "import json,sys; print(json.load(sys.stdin).get('token',''))" 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
  echo "ERROR: Failed to get API token"
  exit 1
fi
echo "Token obtained: ${TOKEN:0:12}..."
echo ""

# Upload function
upload_scan() {
  local report_file="$1"
  local scan_type="$2"

  if [ ! -f "$report_file" ]; then
    echo "  SKIP: $report_file (not found)"
    return 0
  fi

  local response
  response=$(curl -s -X POST "$DEFECTDOJO_URL/api/v2/import-scan/" \
    -H "Authorization: Token $TOKEN" \
    -F "scan_type=$scan_type" \
    -F "product_name=$PRODUCT_NAME" \
    -F "engagement_name=$ENGAGEMENT_NAME" \
    -F "file=@$report_file" 2>&1)

  local test_id
  test_id=$(echo "$response" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('test_id','?'))" 2>/dev/null || echo "?")

  if [ "$test_id" != "?" ] && [ -n "$test_id" ]; then
    echo "  OK: $scan_type -> test_id=$test_id ($report_file)"
  else
    local error
    error=$(echo "$response" | python3 -c "import json,sys; d=json.load(sys.stdin); print(json.dumps(d)[:150])" 2>/dev/null || echo "unknown error")
    echo "  FAIL: $scan_type ($report_file) -> $error"
  fi
}

echo "=== Uploading scanner reports ==="
upload_scan "$SECURITY_DIR/trivy-fs-report.json"        "Trivy Scan"
upload_scan "$SECURITY_DIR/trivy-image-report.json"     "Trivy Scan"
upload_scan "$SECURITY_DIR/trivy-iac-report.json"       "Trivy Scan"
upload_scan "$SECURITY_DIR/trivy-k8s-post-fix-report.json" "Trivy Scan"
upload_scan "$SECURITY_DIR/semgrep-report.json"         "Semgrep JSON Report"
upload_scan "$SECURITY_DIR/checkov-report.json"         "Checkov Scan"
upload_scan "$SECURITY_DIR/checkov-k8s-post-fix-report.json" "Checkov Scan"
echo ""

# Summary
echo "=== Summary ==="
TESTS=$(curl -s "$DEFECTDOJO_URL/api/v2/tests/" -H "Authorization: Token $TOKEN" \
  | python3 -c "import json,sys; print(json.load(sys.stdin).get('count',0))" 2>/dev/null || echo "?")
FINDINGS=$(curl -s "$DEFECTDOJO_URL/api/v2/findings/" -H "Authorization: Token $TOKEN" \
  | python3 -c "import json,sys; print(json.load(sys.stdin).get('count',0))" 2>/dev/null || echo "?")
echo "Tests imported: $TESTS"
echo "Findings: $FINDINGS"
echo ""
echo "Dashboard: $DEFECTDOJO_URL"