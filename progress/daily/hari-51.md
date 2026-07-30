# Hari 51 — CSPM Remediation

**📅 Tanggal:** 2026-07-30
**⏱️ Durasi Belajar:** ~60 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Fix 1: S3 Block Public Access — aktifkan di account level
- [x] Fix 2: Default Security Group — revoke ingress yang terbuka ke internet
- [x] Re-scan Prowler — verify FAIL → PASS
- [x] Dokumentasi CloudTrail recommendations (skip eksekusi — lab challenge)
- [x] Dokumentasi IAM Password Policy recommendations (skip eksekusi — lab challenge)
- [x] Update tracker & dokumentasi

---

## ✅ Yang Berhasil Dikerjakan

### Fix 1: S3 Block Public Access (HIGH → PASS)

**Sebelum:**
- Account-level Block Public Access tidak dikonfigurasi
- 15 S3 findings (public access, no KMS, ACLs, no versioning, dll)

**Command:**
```bash
aws s3control put-public-access-block \
  --account-id <REDACTED> \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true \
  --region ap-southeast-1
```

**Sesudah:**
- `s3_account_level_public_access_blocks`: PASS ✅
- **Semua 15 S3 findings resolved** — 0 S3 FAIL di post-fix scan
- 1 account-level setting fix = 15 findings resolved (cascade effect)

### Fix 2: Default Security Group Revoke Ingress (HIGH → PASS)

**Sebelum:**
- Default SG `sg-bb99b6c8` (VPC `vpc-1d617f7a`): ingress ALL dari 0.0.0.0/0 + ::/0

**Command:**
```bash
aws ec2 revoke-security-group-ingress \
  --group-id sg-bb99b6c8 \
  --region ap-southeast-1 \
  --ip-permissions 'IpProtocol=-1,IpRanges=[{CidrIp=0.0.0.0/0}],Ipv6Ranges=[{CidrIpv6=::/0}]'
```

**Sesudah:**
- Default SG ingress: 0 rules ✅
- `ec2_securitygroup_allow_ingress_from_internet_to_any_port` (default SG): PASS ✅
- Note: Custom SG `sg-0b909f8f4332eecf9` (SSH HTTP HTTPS ICMP) masih FAIL — diluar scope hari ini

---

## 📊 Before vs After

| Metric | Day 50 (Before) | Day 51 (After) | Delta |
|--------|-----------------|----------------|-------|
| Total checks | 328 | 278 | -50 (resource changes) |
| PASS | 192 | 168 | -24 |
| **FAIL** | **132** | **106** | **-26** ✅ |
| CRITICAL | 4 | 1 | -3 |
| HIGH | 23 | 17 | -6 |
| MEDIUM | 85 | 70 | -15 |
| LOW | 20 | 18 | -2 |

### Key Wins

| Fix | Checks Fixed | Findings Resolved |
|-----|-------------|-------------------|
| S3 Block Public Access | 1 account-level → 15 bucket-level | 15 S3 FAIL → 0 |
| Default SG revoke ingress | 1 SG | 1 HIGH → PASS |
| **Total** | | **26 FAILs resolved** |

---

## 📝 Catatan Teknis

### Cascade Effect: S3 Block Public Access
1 command di account level fixed 15 findings. Kenapa? Karena Block Public Access di account level override semua bucket-level settings. Prowler check `s3_account_level_public_access_blocks` PASS, sekaligus bucket-level checks (public read/list/write, bucket policy, ACL) juga PASS karena di-block di account level.

### Post-Fix: Remaining FAILs

| Severity | Count | Top Issues |
|----------|-------|------------|
| CRITICAL | 1 | Root account pakai virtual MFA (bukan hardware MFA) |
| HIGH | 17 | IAM admin tanpa MFA, custom SG open, SNS tidak encrypted, ACM cert akan expired |
| MEDIUM | 70 | CloudWatch alarms, CloudTrail, VPC flow logs, config recorder |
| LOW | 18 | CloudTrail data events, backup, macie, support plan |

---

## 📋 Recommendations (Tidak Dieksekusi — Lab Challenge)

### CloudTrail Recommendations

CloudTrail = audit log untuk semua API calls di AWS. CIS Benchmark mewajibkan minimal 1 multi-region trail.

**Best practice setup:**
1. **Create dedicated S3 bucket** untuk CloudTrail logs (globally unique name)
2. **Bucket policy** granting `cloudtrail.amazonaws.com` write access (`s3:PutObject` + `s3:GetBucketAcl`)
3. **Multi-region trail** (`--is-multi-region-trail`) → record events dari semua region
4. **Log file validation** → SHA256 hash chain untuk detect tampering
5. **KMS encryption** untuk log files at rest
6. **CloudWatch Logs integration** → real-time alerting untuk suspicious API calls
7. **S3 lifecycle policy** → auto-archive ke Glacier setelah 90 hari (save cost)

**CIS Benchmark checks:**
- `cloudtrail_multi_region_enabled` (HIGH)
- `cloudtrail_multi_region_enabled_logging_management_events` (LOW)
- `cloudtrail_s3_dataevents_read_enabled` (LOW)
- `cloudtrail_s3_dataevents_write_enabled` (LOW)

**Cost estimate:** ~$0.05/bulan (management events free, bayar S3 storage)

**Kenapa skip di lab:** CloudTrail butuh S3 bucket + bucket policy + trail creation — 4 steps. Cost minimal tapi tidak esensial untuk belajar CSPM concept. Diprioritaskan di production environment.

### IAM Password Policy Recommendations

Password policy = aturan password untuk semua IAM users yang punya console access.

**Best practice setup (1 command):**
```bash
aws iam update-account-password-policy \
  --minimum-password-length 14 \
  --require-symbols \
  --require-numbers \
  --require-uppercase-characters \
  --require-lowercase-characters \
  --max-password-age 90 \
  --password-reuse-prevention 24
```

**CIS Benchmark checks (5 MEDIUM → PASS):**
- `iam_password_policy_minimum_length_14`
- `iam_password_policy_lowercase`
- `iam_password_policy_uppercase`
- `iam_password_policy_expires_passwords_within_90_days_or_less`
- `iam_password_policy_reused_24`

**Kenapa skip di lab:** Akun sandbox punya 8 IAM users, mayoritas API-only (no console login). Password policy hanya affect users dengan console access. Cost = FREE tapi impact minimal di lab environment.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi |
|----------|--------|
| `--public-access-block-configuration` parameter salah | 4 field: BlockPublicAcls, IgnorePublicAcls, BlockPublicPolicy, RestrictPublicBuckets |
| zsh glob error pada `--ip-permissions` | Quote parameter dengan single quotes |

---

## 📤 Output Hari Ini

- [x] Fix 1: S3 Block Public Access → 15 S3 findings resolved (cascade effect)
- [x] Fix 2: Default SG revoke ingress → 1 HIGH finding resolved
- [x] Re-scan Prowler: 132 FAIL → 106 FAIL (26 resolved)
- [x] CloudTrail recommendations documented
- [x] IAM Password Policy recommendations documented

---

## 💡 Pelajaran Baru

- **Cascade effect.** 1 account-level S3 setting fix = 15 bucket-level findings resolved. Account-level config > individual resource config. Always check if there's an account-level setting yang bisa fix multiple findings sekaligus.
- **CSPM remediation = prioritization.** 132 findings, tidak bisa fix semua sehari. Pilih yang HIGH severity + high impact (cascade). S3 Block Public Access = 1 command, 15 findings. ROI tertinggi.
- **Lab vs Production.** CloudTrail dan IAM password policy penting untuk production, tapi di lab challenge cukup didokumentasikan sebagai recommendations. Production butuh audit trail + password policy, lab cukup paham konsepnya.

---

## 🔗 Referensi

- [Day 50](hari-50.md) — CSPM baseline scan (132 FAIL)
- [AWS S3 Block Public Access](https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-block-public-access.html)
- [AWS CloudTrail Best Practices](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/best-practices-for-security.html)
- [AWS IAM Password Policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_passwords_account-policy.html)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | 1 command = 15 findings, cascade effect very satisfying |
| Pemahaman materi | 5 | Account-level vs bucket-level, CSPM prioritization |
| Progres sesuai target | 5 | 26 FAILs resolved, recommendations documented |

---

## ➡️ Rencana Besok

- [ ] **Day 52: Red Team — K8s Escape** — Privileged pod → host filesystem → Falco/OPA deteksi

---

*[← Hari 50](hari-50.md) | [Hari 52 →](hari-52.md)*
