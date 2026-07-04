# Hari 27 — AI-Assisted IaC Fix

**📅 Tanggal:** 2026-07-04  
**⏱️ Durasi Belajar:** ~1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] AI review Terraform terhadap CIS AWS Foundations Benchmark + Well-Architected Framework
- [x] Tambah variable validation (region format, environment enum, instance type)
- [x] Buat `outputs.tf` untuk visibility resource attributes
- [x] Tighten IAM assume role policies dengan `Condition: aws:SourceAccount`
- [x] Tambah S3 `noncurrent_version_transitions` di lifecycle rules
- [x] Parameterize KMS deletion window + CloudWatch retention
- [x] Checkov + Trivy IaC tetap 0 findings

---

## ✅ Yang Berhasil Dikerjakan

- **`variables.tf`** — 6 variables dengan validation blocks:
  - `aws_region` — regex validation untuk AWS region format
  - `replication_region` — same regex validation
  - `environment` — enum validation (`dev`, `staging`, `prod`)
  - `instance_type` — regex validation untuk valid AWS instance types
  - `kms_deletion_window` — range validation (7-30 days, AWS limit)
  - `log_retention_days` — enum validation untuk CloudWatch retention values

- **`outputs.tf`** — NEW file, 12 outputs:
  - VPC ID, VPC CIDR, S3 bucket names+ARNs, KMS key ARN+alias, SG ID, subnet ID, flow log group name, SNS topic ARN

- **`iam.tf`** — Added `Condition: aws:SourceAccount` ke 3 trust policies:
  - Flow log role: hanya akun yang sama yang bisa assume
  - EC2 role: same
  - Replication role: same

- **`s3.tf`** — Added `noncurrent_version_transition` + `noncurrent_version_expiration` ke logs bucket lifecycle (move old versions to STANDARD_IA → GLACIER → expire after 365 days)

- **`kms.tf`** — Parameterize `deletion_window_in_days` ke `var.kms_deletion_window`

- **`network.tf`** — Parameterize `retention_in_days` ke `var.log_retention_days`

- **`compute.tf`** — Added `prevent_destroy` lifecycle ke subnet

- Checkov: **102 passed, 0 failed, 10 skipped** ✅
- Trivy IaC: **0 findings (MEDIUM+)** ✅
- Both pipelines GREEN

---

## 📝 Catatan Teknis

### Variable Validation Example
```hcl
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "environment must be one of: dev, staging, prod."
  }
}
```

### IAM Trust Policy Tightening
```hcl
Condition = {
  StringEquals = {
    "aws:SourceAccount" = data.aws_caller_identity.current.account_id
  }
}
```
Tanpa Condition, siapa saja dengan `sts:AssumeRole` bisa assume role walau dari account lain. Dengan `aws:SourceAccount`, hanya service dari account yang sama yang bisa assume.

### S3 Non-Current Version Lifecycle
```hcl
noncurrent_version_transition {
  noncurrent_days = 30
  storage_class   = "STANDARD_IA"
}
noncurrent_version_transition {
  noncurrent_days = 90
  storage_class   = "GLACIER"
}
noncurrent_version_expiration {
  noncurrent_days = 365
}
```
Versioning aktif = setiap update bikin version baru. Tanpa non-current lifecycle, versi lama menumpuk tanpa batas.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| KMS deletion window awalnya diset 90 → AWS limit max 30 | Set ke 30 (AWS max) dan bikin variable dengan validation range 7-30 |
| Checkov `prevent_destroy` belum ada di subnet | Tambah `lifecycle { prevent_destroy = true }` ke `aws_subnet.app` |

---

## 📤 Output Hari Ini

- [x] `variables.tf` — 6 variables, semua dengan validation blocks
- [x] `outputs.tf` — 12 outputs (NEW file)
- [x] `iam.tf` — 3 trust policies tightened dengan `aws:SourceAccount`
- [x] `s3.tf` — noncurrent_version_transitions + expiration di logs bucket
- [x] `kms.tf` — parameterize deletion window
- [x] `network.tf` — parameterize log retention
- [x] `compute.tf` — prevent_destroy di subnet
- [x] Checkov: 102 passed, 0 failed
- [x] Trivy IaC: 0 findings (MEDIUM+)
- [x] Pipeline Infrastructure Security: SUCCESS (29s)
- [x] Pipeline SecureBank CI: SUCCESS (4m28s)
- [x] Commit: `fcd29f1`

---

## 💡 Pelajaran Baru

- **Variable validation = shift-left untuk configuration.** Tanpa validation, typo di `environment = "devv"` akan propagate ke semua resource tags. Dengan validation, error ketangkap saat `terraform plan` — sebelum resource dibuat.

- **KMS `deletion_window_in_days` max 30.** Awalnya ngira bisa 90 hari untuk production safety. Ternyata AWS limit adalah 7-30. Documented di [AWS KMS API](https://docs.aws.amazon.com/kms/latest/APIReference/API_CreateKey.html).

- **`aws:SourceAccount` Condition = confusion matrix prevention.** Tanpa Condition, trust policy `Principal: { Service: "s3.amazonaws.com" }` ngijinin S3 dari ANY account assume role. Dengan `aws:SourceAccount`, hanya S3 dari account yang sama yang bisa.

- **`noncurrent_version_transition` = lifecycle untuk versioned objects.** Versioning aktif bikin setiap overwrite produce version baru. Tanpa non-current lifecycle, versi lama menumpuk permanent. `noncurrent_version_expiration` juga penting untuk cost control.

- **`outputs.tf` bukan boilerplate — itu documentation.** Output yang jelas bikin Terraform config self-documenting. `terraform output` setelah apply langsung kasih VPC ID, bucket ARN, dll. Tanpa output, harus cek console atau `terraform show`.

---

## 🔗 Referensi

- [CIS AWS Foundations Benchmark v2.0](https://docs.aws.amazon.com/security-hub/latest/userguide/cis-aws-foundations-benchmark.html)
- [AWS Well-Architected Security Pillar](https://docs.aws.amazon.com/wellarchitected/latest/security-pillar/)
- [Terraform Variable Validation](https://developer.hashicorp.com/terraform/language/values/variables#custom-validation-rules)
- [AWS KMS deletion window](https://docs.aws.amazon.com/kms/latest/developerguide/deleting-keys.html)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | AI-assisted review, banyak improvement yang scanner nggak nemuin |
| Pemahaman materi | 5 | Variable validation, IAM Condition, S3 non-current lifecycle |
| Progres sesuai target | 5 | Pipeline green, 0 findings, todos clear |

---

## ➡️ Rencana Besok

- [ ] Hari 28: Compliance as Code (InSpec) — profile untuk verifikasi port 22 tertutup, API listening, non-root process

---

*[← Hari 26](hari-26.md) | [Hari 28 →](hari-28.md)*