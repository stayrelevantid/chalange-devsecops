# IaC Scan Comparison: Checkov (IaC) vs Trivy IaC

**📅 Tanggal:** 2026-06-25  
**📋 Target:** `securebank-api/terraform/main.tf` (3 resources: VPC, S3, Security Group)

---

## Summary

| Metric | Checkov (IaC) | Trivy IaC |
|-------|--------------|-----------|
| **Total Failed** | 15 | 14 |
| **Total Passed** | 6 | 0 |
| **CRITICAL** | 0 | 1 |
| **HIGH** | 0 | 6 |
| **MEDIUM** | 0 | 2 |
| **LOW** | 0 | 5 |
| **Severity System** | Pass/Fail only | CVSS-like (CRITICAL→LOW) |

---

## Findings Comparison

### S3 Bucket (`aws_s3_bucket.logs`)

| # | Finding | Checkov (IaC) | Trivy IaC | Match? |
|---|---------|--------------|-----------|--------|
| 1 | No public access block | ✅ CKV2_AWS_6 | ✅ AWS-0086 (HIGH) | ✅ Match |
| 2 | No public access block (block public ACLs) | ✅ CKV2_AWS_6 | ✅ AWS-0091 (HIGH) | ⚠️ Partial (Trivy splits into 4 rules) |
| 3 | No public access block (block public policies) | ✅ CKV2_AWS_6 | ✅ AWS-0087 (HIGH) | ⚠️ Partial |
| 4 | No public access block (restrict public buckets) | ✅ CKV2_AWS_6 | ✅ AWS-0093 (HIGH) | ⚠️ Partial |
| 5 | No public access block (general) | ✅ CKV2_AWS_6 | ✅ AWS-0094 (LOW) | ⚠️ Partial |
| 6 | No access logging | ✅ CKV_AWS_18 | ✅ AWS-0089 (LOW) | ✅ Match |
| 7 | No versioning | ✅ CKV_AWS_21 | ✅ AWS-0090 (MEDIUM) | ✅ Match |
| 8 | No KMS encryption | ✅ CKV_AWS_145 | ✅ AWS-0132 (HIGH) | ✅ Match |
| 9 | No cross-region replication | ✅ CKV_AWS_144 | ❌ | Checkov only |
| 10 | No event notifications | ✅ CKV2_AWS_62 | ❌ | Checkov only |
| 11 | No lifecycle configuration | ✅ CKV2_AWS_61 | ❌ | Checkov only |

### Security Group (`aws_security_group.api`)

| # | Finding | Checkov (IaC) | Trivy IaC | Match? |
|---|---------|--------------|-----------|--------|
| 12 | Ingress 0.0.0.0/0 to port 22 (SSH) | ✅ CKV_AWS_24 | ✅ AWS-0107 (HIGH) | ✅ Match (Trivy broader) |
| 13 | Ingress 0.0.0.0/0 to port 80 (HTTP) | ✅ CKV_AWS_260 | ✅ AWS-0107 (HIGH) | ✅ Match (Trivy broader) |
| 14 | Ingress 0.0.0.0/0 to port 3389 (RDP) | ✅ CKV_AWS_25 | ✅ AWS-0107 (HIGH) | ✅ Match (Trivy broader) |
| 15 | Egress 0.0.0.0/0 to all ports | ✅ CKV_AWS_382 | ✅ AWS-0104 (CRITICAL) | ✅ Match |
| 16 | No description on SG | ✅ CKV_AWS_23 | ✅ AWS-0099 (LOW) | ✅ Match |
| 17 | No description on ingress rule | ✅ CKV_AWS_23 | ✅ AWS-0124 (LOW) | ✅ Match (Trivy splits rules) |
| 18 | No description on egress rule | ✅ CKV_AWS_23 | ✅ AWS-0124 (LOW) | ✅ Match (Trivy splits rules) |
| 19 | SG not attached to resource | ✅ CKV2_AWS_5 | ❌ | Checkov only |

### VPC (`aws_vpc.main`)

| # | Finding | Checkov (IaC) | Trivy IaC | Match? |
|---|---------|--------------|-----------|--------|
| 20 | No VPC flow logging | ✅ CKV2_AWS_11 | ✅ AWS-0178 (MEDIUM) | ✅ Match |
| 21 | Default SG doesn't restrict traffic | ✅ CKV2_AWS_12 | ❌ | Checkov only |

---

## Unique Findings

### Checkov (IaC) Only (4 findings)
- CKV_AWS_144 — S3 no cross-region replication
- CKV2_AWS_62 — S3 no event notifications
- CKV2_AWS_61 — S3 no lifecycle configuration
- CKV2_AWS_5 — SG not attached to any resource
- CKV2_AWS_12 — Default SG of VPC doesn't restrict traffic

### Trivy IaC Only (0 findings)
Trivy tidak menemukan finding yang unik. Semua temuan Trivy juga ditemukan oleh Checkov.

---

## Key Differences

### 1. Severity Classification
- **Checkov:** Hanya Pass/Fail, tidak ada severity level
- **Trivy IaC:** Ada severity (CRITICAL, HIGH, MEDIUM, LOW) — berguna untuk prioritisasi

### 2. Rule Granularity
- **Checkov:** 1 rule untuk "no public access block" (CKV2_AWS_6)
- **Trivy IaC:** 4 rules untuk hal yang sama (AWS-0086, AWS-0087, AWS-0091, AWS-0093, AWS-0094) — lebih detail

### 3. Port-Specific vs Generic Rules
- **Checkov:** Buat rule terpisah untuk setiap port (22, 80, 3389)
- **Trivy IaC:** 1 rule generic (AWS-0107) untuk semua unrestricted ingress

### 4. Coverage
- **Checkov:** 15 failed + 6 passed = 21 total checks
- **Trivy IaC:** 14 failed + 0 passed = 14 total checks
- **Checkov lebih comprehensive** dalam hal total rules

### 5. Passed Checks
- **Checkov:** Menampilkan 6 passed checks (transparan)
- **Trivy IaC:** Hanya tampilkan failures, tidak tampilkan passed

---

## Recommendation

| Use Case | Recommended Tool |
|----------|-----------------|
| Comprehensive scanning | Checkov (21 rules vs 14) |
| Severity-based prioritization | Trivy IaC (CRITICAL/HIGH/MEDIUM/LOW) |
| CI/CD pipeline (GitHub Actions) | Trivy IaC (sudah pakai Trivy untuk SCA, 1 tool 2 mode) |
| Compliance reporting | Checkov (passed + failed transparan) |
| Production gate | Keduanya (defense in depth — 2 vendor) |

### Final Verdict
Pakai **keduanya** di pipeline. Checkov untuk coverage, Trivy IaC untuk severity-based triage. Overlap findings = confirmatory, unique findings = complementary.