# Hari 20 — Terraform Setup + IaC Scan (Checkov)

**📅 Tanggal:** 2026-06-24  
**⏱️ Durasi Belajar:** 2 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install Terraform v1.15.6 (via tfenv)
- [x] Install Checkov v3.3.2 (via pip3)
- [x] Buat `terraform/main.tf` dengan intentional misconfiguration
- [x] Buat `terraform/variables.tf`
- [x] Run Checkov scan — 15 failed checks, 6 passed checks
- [x] Simpan report JSON

---

## ✅ Yang Berhasil Dikerjakan

- Install Terraform v1.15.6 (tfenv auto-install)
- Install Checkov v3.3.2 (`pip3 install --user --break-system-packages checkov`)
- Buat `securebank-api/terraform/main.tf` — VPC + S3 (insecure) + Security Group (open)
- Buat `securebank-api/terraform/variables.tf` — `aws_region` + `environment`
- Run `checkov -d securebank-api/terraform/` → **6 passed, 15 failed**
- Simpan report: `securebank-api/security/checkov-report.json`

---

## 📝 Catatan Teknis

```bash
# Install Terraform (tfenv auto-install latest stable)
$ tfenv use latest
Terraform v1.15.6

# Install Checkov
$ pip3 install --user --break-system-packages checkov
$ checkov --version
3.3.2

# Checkov scan
$ checkov -d securebank-api/terraform/

Passed checks: 6, Failed checks: 15, Skipped checks: 0

# Failed checks breakdown:
# Security Group (8 failures):
  CKV_AWS_260  - Ingress 0.0.0.0/0 to port 80
  CKV_AWS_24   - Ingress 0.0.0.0/0 to port 22
  CKV_AWS_25   - Ingress 0.0.0.0/0 to port 3389
  CKV_AWS_23   - No description on SG rules
  CKV_AWS_382  - Egress 0.0.0.0/0 to all ports
  CKV2_AWS_5   - SG not attached to any resource
  CKV2_AWS_12  - Default SG of VPC doesn't restrict traffic

# S3 Bucket (7 failures):
  CKV2_AWS_6   - No public access block
  CKV2_AWS_62  - No event notifications
  CKV2_AWS_61  - No lifecycle configuration
  CKV_AWS_18   - No access logging
  CKV_AWS_145  - No KMS encryption
  CKV_AWS_144  - No cross-region replication
  CKV_AWS_21   - No versioning

# Save report
$ checkov -d securebank-api/terraform/ -o json --output-file-path securebank-api/security/
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `terraform init` timeout download AWS provider (~500MB) | Checkov static analyze HCL, tidak butuh provider. Skip `terraform init` |
| Python `externally-managed-environment` error | Pakai `--break-system-packages` flag |
| Checkov binary tidak di PATH | Tambahkan `~/Library/Python/3.14/bin` ke PATH |
| User minta Terraform v3.2.2 | Terraform tidak punya v3.x. Latest: v1.15.6. Lanjut pakai v1.15.6 |

---

## 📤 Output Hari Ini

- [x] `securebank-api/terraform/main.tf` — VPC + S3 (insecure) + SG (open)
- [x] `securebank-api/terraform/variables.tf` — variables
- [x] `securebank-api/security/checkov-report.json` — Checkov report (15 failed, 6 passed)
- [x] Commit: Day 20

---

## 💡 Pelajaran Baru

- **Checkov tidak butuh `terraform init`.** Checkov static analyze HCL files, tidak perlu download provider. Ini berguna saat network lambat atau tidak ada AWS credentials.

- **15 failed checks dari 3 resources.** Satu Security Group yang terbuka ke 0.0.0.0/0 untuk semua port (0-65535) akan trigger multiple Checkov rules — port 22, 80, 3389, dll. Satu bucket S3 tanpa encryption/versioning/logging juga trigger banyak rules.

- **VPC tanpa flow logging langsung kena.** Checkov mendeteksi `aws_vpc.main` tanpa `aws_flow_log` resource terkait. Best practice: enable VPC flow logging untuk audit network traffic.

- **Default Security Group harus restrict all.** Checkov mendeteksi bahwa VPC default security group tidak dikonfigurasi untuk restrict all traffic. AWS create default SG otomatis saat VPC dibuat, dan harus di-hardcoded untuk restrict.

- **IaC scan lebih cepat dari manual review.** Checkov scan 3 file Terraform dalam hitungan detik dan menemukan 15 issues. Manual review akan ambil waktu lama dan mudah ada yang terlewat.

- **`CKV_AWS_19` PASSED tapi `CKV_AWS_145` FAILED.** Checkov punya 2 rules encryption: satu cek SSE (default encryption), satu cek KMS (customer-managed key). Bucket tanpa keduanya pass SSE check (karena S3 default encryption) tapi fail KMS check.

---

## 🔗 Referensi

- [Checkov Documentation](https://www.checkov.io/)
- [Checkov AWS Policies](https://docs.prismacloud.io/en/enterprise-edition/policy-reference/aws-policies)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [tfenv](https://github.com/tfutils/tfenv)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | First time pakai Terraform + Checkov, langsung dapet 15 findings |
| Pemahaman materi | 5 | Konsep IaC scanning jelas — static analysis HCL files |
| Progres sesuai target | 5 | Day 20 selesai, Fase 2 lanjut 5/15 |

---

## ➡️ Rencana Besok

- [ ] Hari 21: IaC Scan (tfsec/Trivy) — scan ulang dengan tfsec, bandingkan temuan dengan Checkov

---

*[← Hari 19](hari-19.md) | [Hari 21 →](hari-21.md)*