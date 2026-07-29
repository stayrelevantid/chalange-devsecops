# Hari 50 — CSPM (Prowler)

**📅 Tanggal:** 2026-07-29
**⏱️ Durasi Belajar:** ~90 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Install Prowler
- [x] Run Prowler scan against AWS sandbox (ap-southeast-1)
- [x] Generate HTML + JSON reports
- [x] Analyze findings — kategorisasi by service & severity
- [x] Dokumentasikan baseline untuk Day 51 remediation
- [x] Update tracker & dokumentasi

---

## ✅ Yang Berhasil Dikerjakan

### 1. Prowler Installation
- Syllabus bilang `pip install prowler` — tapi pydantic v1 vs v2 conflict di Python 3.14
- Solution: `brew install prowler` — Prowler v3.11.3 with Python 3.13, no conflicts
- Checkov (needs pydantic v2) dan Prowler (needs pydantic v1) coexist via separate Python envs

### 2. Prowler Scan Results

| Metric | Count |
|--------|-------|
| **Total checks** | 328 |
| **PASS** | 192 (58.5%) |
| **FAIL** | 132 (40.2%) |
| **INFO** | 4 (1.2%) |

### 3. FAIL by Severity

| Severity | Count |
|----------|-------|
| CRITICAL | 4 |
| HIGH | 23 |
| MEDIUM | 85 |
| LOW | 20 |
| **Total FAIL** | **132** |

### 4. FAIL by Service (Top 10)

| Service | Fails | Key Issues |
|---------|-------|------------|
| IAM | 44 | Admin access without MFA, no password policy, access keys not rotated |
| CloudWatch | 23 | No metric filters/alarms for network changes, log groups without KMS |
| S3 | 15 | Block Public Access not configured, no KMS encryption, ACLs enabled |
| EC2 | 10 | NACL all ports open, default SG all ports open, EBS not encrypted |
| VPC | 6 | Flow logs disabled, subnets assign public IP by default |
| CloudFormation | 5 | Potential secrets in stack outputs (CRITICAL), termination protection disabled |
| CloudTrail | 4 | No CloudTrail enabled, no S3 data events |
| Organizations | 3 | AWS Organizations not in use |
| SNS | 3 | SNS topics not KMS encrypted |
| Athena | 2 | Query results not encrypted, config not enforced |

### 5. CRITICAL Findings (4)

All 4 CRITICAL findings are CloudFormation stack outputs containing potential secrets — old Amplify stacks from 2021.

| CheckID | Resource |
|---------|----------|
| cloudformation_stack_outputs_find_secrets | Stack amplify-react-draw-devy-152142 |
| cloudformation_stack_outputs_find_secrets | Stack amplify-authcra-deva-173047-authcognitocf0c6096 |
| cloudformation_stack_outputs_find_secrets | Stack amplify-authcra-deva-173047 |

### 6. Top Findings for Day 51 Remediation

1. **[HIGH] S3 Block Public Access not configured** — `s3_account_level_public_access_blocks`
2. **[HIGH] IAM admin access without MFA** — `iam_administrator_access_with_mfa`
3. **[HIGH] CloudTrail not enabled** — `cloudtrail_multi_region_enabled`
4. **[CRITICAL] Potential secrets in CloudFormation outputs** — `cloudformation_stack_outputs_find_secrets`

---

## 📝 Catatan Teknis

### Prowler Run Commands
```bash
# JSON report
prowler aws -M json -f ap-southeast-1 -o /tmp/prowler-output

# HTML report (separate run for reliability)
prowler aws -M html -f ap-southeast-1 -o /tmp/prowler-output
```

### Prowler vs Checkov Scope
- **Prowler:** Cloud runtime posture (IAM, S3, CloudTrail, VPC, CloudWatch) — live AWS account scan
- **Checkov:** IaC static analysis (Terraform, K8s manifests) — pre-deploy config scan
- **Complementary:** Prowler finds runtime drift, Checkov catches config-level issues

---

## 📊 Perubahan File

| File | Status | Description |
|------|--------|-------------|
| `security/prowler/prowler-report.json` | ✅ Created | Prowler JSON scan report (328 checks) |
| `security/prowler/prowler-report.html` | ✅ Created | Prowler HTML report (visual dashboard) |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi |
|----------|--------|
| pydantic v1/v2 conflict (Python 3.14) | `brew install prowler` — Python 3.13 env, no conflict |
| Prowler scan >2 min timeout | Run JSON and HTML separately, background process |
| First JSON truncated (timeout cut scan) | Re-run with longer wait (150 seconds) |

---

## 📤 Output Hari Ini

- [x] Prowler v3.11.3 installed via Homebrew
- [x] AWS scan complete: 328 checks, 192 PASS, 132 FAIL
- [x] HTML + JSON reports saved to `security/prowler/`
- [x] 4 CRITICAL + 23 HIGH findings identified for Day 51

---

## 💡 Pelajaran Baru

- **Prowler = CSPM scanner.** Scan live AWS account against CIS Benchmarks (300+ checks). Berbeda dengan Checkov (IaC pre-deploy). Prowler catches runtime drift.
- **Feedback loop.** Checkov (IaC) → Prowler (runtime) = full coverage. Checkov passes Terraform, tapi Prowler finds real AWS account issues.
- **Python version matters.** pydantic v1 (Prowler) vs v2 (Checkov) conflict. brew install = separate Python env.

---

## 🔗 Referensi

- [Prowler GitHub](https://github.com/prowler-cloud/prowler)
- [AWS CIS Benchmarks](https://docs.aws.amazon.com/securityhub/latest/userguide/cis-aws-foundations-benchmark.html)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Real AWS scan, 132 findings — eye opening |
| Pemahaman materi | 5 | CIS Benchmarks, CSPM scope clear |
| Progres sesuai target | 5 | Scan done, baseline for Day 51 |

---

## ➡️ Rencana Besok

- [ ] **Day 51: CSPM Remediation** — Fix 3 critical findings from Prowler baseline

---

*[← Hari 49](hari-49.md) | [Hari 51 →](hari-51.md)*