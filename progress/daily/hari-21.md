# Hari 21 — IaC Scan: Checkov (IaC) vs Trivy IaC

**📅 Tanggal:** 2026-06-25  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Scan Terraform dengan Trivy IaC mode (`trivy config`)
- [x] Simpan Trivy IaC report JSON
- [x] Buat comparison table: Checkov (IaC) vs Trivy IaC
- [x] Analisis perbedaan coverage, severity, dan granularity

---

## ✅ Yang Berhasil Dikerjakan

- Scan Terraform dengan Trivy IaC: **14 findings** (1 CRITICAL, 6 HIGH, 2 MEDIUM, 5 LOW)
- Simpan report: `securebank-api/security/trivy-iac-report.json`
- Buat comparison doc: `securebank-api/security/iac-scan-comparison.md`
- Analisis: Checkov 15 failed/6 passed, Trivy IaC 14 failed/0 passed
- Checkov menemukan 5 unique findings, Trivy IaC 0 unique

---

## 📝 Catatan Teknis

```bash
# Trivy IaC scan
$ trivy config securebank-api/terraform/

Report Summary
  main.tf │ terraform │ 14 misconfigurations

Failures: 14 (UNKNOWN: 0, LOW: 5, MEDIUM: 2, HIGH: 6, CRITICAL: 1)

# Key findings:
# CRITICAL: AWS-0104 — Unrestricted egress to 0.0.0.0/0
# HIGH: AWS-0107 — Unrestricted ingress from 0.0.0.0/0
# HIGH: AWS-0086/0087/0091/0093 — S3 no public access block (4 rules)
# HIGH: AWS-0132 — S3 no KMS encryption
# MEDIUM: AWS-0090 — S3 no versioning
# MEDIUM: AWS-0178 — VPC no flow logging
# LOW: AWS-0089/0094/0099/0124 — S3 logging, SG descriptions

# Save report
$ trivy config securebank-api/terraform/ --format json --output securebank-api/security/trivy-iac-report.json
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| tfsec sudah deprecated | Skip tfsec, pakai Trivy IaC (`trivy config`) saja — sudah mengintegrasikan tfsec's rule set |
| Trivy perlu download checks bundle | Otomatis ter-download saat first run (~234KB, cepat) |

---

## 📤 Output Hari Ini

- [x] `securebank-api/security/trivy-iac-report.json` — Trivy IaC report (14 findings)
- [x] `securebank-api/security/iac-scan-comparison.md` — Comparison table Checkov vs Trivy IaC
- [x] `progress/daily/hari-21.md` — Progress harian

---

## 💡 Pelajaran Baru

- **Trivy IaC punya severity level, Checkov tidak.** Trivy memberi CRITICAL/HIGH/MEDIUM/LOW untuk setiap finding. Checkov cuma Pass/Fail. Severity penting untuk prioritisasi remediasi — CRITICAL di-fix dulu.

- **Checkov lebih granular untuk port-specific rules.** Checkov bikin rule terpisah untuk port 22, 80, 3389. Trivy pakai 1 rule generic (AWS-0107) untuk semua unrestricted ingress. Checkov lebih detail, Trivy lebih ringkas.

- **Trivy lebih granular untuk S3 public access block.** Checkov pakai 1 rule (CKV2_AWS_6). Trivy pecah jadi 4 rules (AWS-0086/0087/0091/0093) — satu untuk setiap aspek public access block (block ACLs, block policies, ignore ACLs, restrict buckets).

- **Checkov menampilkan passed checks, Trivy tidak.** Checkov transparan — tampilkan 6 checks yang passed. Trivy hanya tampilkan failures. Passed checks berguna untuk audit trail dan dokumentasi compliance.

- **Checkov menemukan 5 unique findings yang Trivy tidak.** S3 cross-region replication, event notifications, lifecycle config, SG not attached to resource, default SG restrict traffic. Checkov lebih comprehensive untuk best-practice rules.

- **1 tool (Trivy) = 3 mode scan.** Trivy bisa `fs` (dependency/SCA), `image` (container CVE), dan `config` (IaC). Tidak perlu install tool terpisah. Tapi Checkov tetap value-add untuk coverage yang lebih luas.

---

## 🔗 Referensi

- [Trivy IaC Scanning](https://trivy.dev/latest/docs/scanner/misconfig/)
- [Trivy AWS Checks](https://avd.aquasec.com/)
- [Checkov Documentation](https://www.checkov.io/)
- [tfsec Deprecated → Trivy](https://github.com/aquasecurity/tfsec)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Bandingkan 2 tools, happy dengan insight |
| Pemahaman materi | 5 | Clear: Checkov comprehensive, Trivy severity-based |
| Progres sesuai target | 5 | Day 21 selesai, Fase 2 lanjut 6/15 |

---

## ➡️ Rencana Besok

- [ ] Hari 22: IaC Scan di Pipeline — tambah Checkov + Trivy IaC ke `infra.yml` workflow

---

*[← Hari 20](hari-20.md) | [Hari 22 →](hari-22.md)*