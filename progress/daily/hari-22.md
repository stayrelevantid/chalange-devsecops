# Hari 22 — IaC Scan di Pipeline (GitHub Actions)

**📅 Tanggal:** 2026-06-26  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat `.github/workflows/infra.yml` — pipeline khusus infrastructure security
- [x] 2 scan jobs berjalan paralel: Checkov (IaC) + Trivy IaC
- [x] Pipeline trigger hanya pada perubahan `securebank-api/terraform/**`
- [x] Pipeline **GAGAL** (RED) karena Terraform masih intentionally misconfigured
- [x] SARIF uploaded ke GitHub Security tab
- [x] Artifacts uploaded (checkov-sarif, trivy-iac-report)

---

## ✅ Yang Berhasil Dikerjakan

- Buat `.github/workflows/infra.yml` dengan 2 parallel jobs:
  - **Checkov (IaC) Scan** — `bridgecrewio/checkov-action@v12`, `soft_fail: false`, output SARIF, upload ke GitHub Security
  - **Trivy IaC Scan** — `aquasecurity/trivy-action@master`, `scan-type: config`, `exit-code: 1`
- Push ke main → pipeline trigger otomatis
- Pipeline run: **FAILED** (expected RED) — 35 detik total
- Checkov: 10 finding annotations di GitHub Actions
- Trivy IaC: exit code 1 (findings detected)
- SARIF uploaded ke GitHub Security tab ✅
- Artifacts: `checkov-sarif` + `trivy-iac-report` ✅

---

## 📝 Catatan Teknis

```yaml
# .github/workflows/infra.yml — key sections

on:
  push:
    paths:
      - 'securebank-api/terraform/**'
      - '.github/workflows/infra.yml'
  pull_request:
    paths:
      - 'securebank-api/terraform/**'
      - '.github/workflows/infra.yml'

jobs:
  checkov-scan:          # Checkov (IaC) — SARIF output + GitHub Security
    uses: bridgecrewio/checkov-action@v12
    with:
      directory: securebank-api/terraform/
      soft_fail: false   # FAIL pipeline if findings

  trivy-iac-scan:        # Trivy IaC — config mode
    uses: aquasecurity/trivy-action@master
    with:
      scan-type: 'config'
      scan-ref: 'securebank-api/terraform/'
      exit-code: '1'     # FAIL pipeline if findings
```

```bash
# Pipeline result
$ gh run view 28216885163

Jobs:
  X Trivy IaC Scan — FAILED (11s)
  X Checkov (IaC) Scan — FAILED (30s)

Checkov annotations (10 findings):
  CKV2_AWS_6   — S3 no public access block
  CKV_AWS_18   — S3 no access logging
  CKV_AWS_145  — S3 no KMS encryption
  CKV2_AWS_5   — SG not attached to resource
  CKV2_AWS_12  — Default SG doesn't restrict traffic
  CKV_AWS_260  — Ingress 0.0.0.0/0 port 80
  CKV_AWS_23   — SG no description
  CKV_AWS_24   — Ingress 0.0.0.0/0 port 22
  CKV_AWS_25   — Ingress 0.0.0.0/0 port 3389
  CKV_AWS_382  — Egress 0.0.0.0/0 all ports

Trivy IaC: exit code 1 (findings detected)

Artifacts: checkov-sarif, trivy-iac-report
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Tutorial pakai tfsec-action | tfsec deprecated. Ganti dengan Trivy IaC (`scan-type: config`) — sudah mengintegrasikan tfsec rule set |
| Path filter harus sesuai struktur repo | Pakai `securebank-api/terraform/**` (bukan `terraform/**` dari tutorial) |
| Node.js 20 deprecation warning | Warning saja, tidak menggagalkan pipeline. Actions tetap jalan (forced to Node.js 24) |
| CodeQL Action v3 deprecation warning | Warning saja. Akan update ke v4 di Day 29 (Pipeline Consolidation) |

---

## 📤 Output Hari Ini

- [x] `.github/workflows/infra.yml` — Infrastructure Security workflow
- [x] Pipeline run: ID 28216885163 — FAILED (expected RED)
- [x] SARIF uploaded ke GitHub Security tab
- [x] Artifacts: checkov-sarif + trivy-iac-report

---

## 💡 Pelajaran Baru

- **Pipeline RED = security gate bekerja.** Tujuan hari ini bukan bikin pipeline hijau, tapi memastikan pipeline **gagal** ketika ada misconfiguration. Kalau pipeline hijau dengan 15 findings, berarti security gate-nya rusak.

- **Path filter bikin pipeline efisien.** Workflow hanya trigger saat ada perubahan di `securebank-api/terraform/**` atau `infra.yml` itu sendiri. Push ke `main.go` atau `Dockerfile` tidak trigger infra pipeline — hemat CI minutes.

- **`soft_fail: false` di Checkov = pipeline gagal.** Default-nya `true` (soft fail = findings dicatat tapi pipeline tetap hijau). Dengan `false`, findings menggagalkan pipeline. Ini yang kita mau untuk security gate.

- **`exit-code: 1` di Trivy = pipeline gagal.** Trivy mengembalikan exit code 1 kalau ada findings. GitHub Actions menganggap exit code non-zero sebagai failure. Tanpa `exit-code: 1`, Trivy mengembalikan 0 walaupun ada findings.

- **SARIF = format standar untuk security findings.** Checkov output SARIF, di-upload ke GitHub Security tab lewat `codeql-action/upload-sarif`. Hasilnya: findings muncul di repo Security tab dengan lokasi file dan baris yang tepat. Trivy juga support SARIF, tapi hari ini kita pakai JSON untuk Trivy.

- **2 tools = 2 perspektif.** Checkov dan Trivy IaC berjalan paralel, keduanya gagal. Tapi finding-nya tidak identik — Checkov nemuin 5 unique findings (Day 21 comparison). Pipeline dengan 2 tools = defense in depth di CI/CD level.

- **GitHub Actions menampilkan annotations otomatis.** 10 Checkov findings muncul sebagai annotations di pipeline run view. Setiap annotation punya check ID, rule name, dan link ke guideline. Tidak perlu buka report JSON untuk lihat findings — langsung di UI.

---

## 🔗 Referensi

- [Checkov GitHub Action](https://github.com/bridgecrewio/checkov-action)
- [Trivy GitHub Action](https://github.com/aquasecurity/trivy-action)
- [SARIF Upload to GitHub Security](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/uploading-a-sarif-file-to-github)
- [Path Filter for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#onpushpull_requestpull_requestpathspaths-ignore)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline RED yang expected = security gate jalan! |
| Pemahaman materi | 5 | Path filter, soft_fail, exit-code, SARIF — semua clear |
| Progres sesuai target | 5 | Day 22 selesai, Fase 2 lanjut 7/15 |

---

## ➡️ Rencana Besok

- [ ] Hari 23: IaC Remediation — fix semua misconfiguration Terraform hingga pipeline hijau

---

*[← Hari 21](hari-21.md) | [Hari 23 →](hari-23.md)*