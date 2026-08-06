# Hari 55 — Pembuatan Laporan Audit

**📅 Tanggal:** 2026-08-06
**⏱️ Durasi Belajar:** ~120 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Re-sync DefectDojo dengan report scan terbaru (data sebelumnya stale dari Day 47)
- [x] Export temuan DefectDojo sebagai evidence + hitung metrik
- [x] Generate report DefectDojo (Executive Summary) via UI
- [x] Buat `security/audit-reports/draft-q3.md` — draf laporan eksekutif dengan metrik nyata

---

## ✅ Yang Berhasil Dikerjakan

### 1. Audit & Regenerasi Report Scanner

Report yang lama ternyata **stale** (sebagian besar dari Juni, sebelum remediasi Day 23/33). Di-regenerate pakai state terkini:

| Report | Sebelum (stale) | Setelah regenerasi | Keterangan |
|--------|-----------------|--------------------|------------|
| `trivy-fs-report.json` | 0 (Jun 11) | 1 Info (GO-2026-5932 openpgp) | Dijalankan ulang `--scanners vuln` agar bersih dari secret scan file gitignored |
| `trivy-image-report.json` | 0 (Jun 21) | **1 High + 1 Medium** (stdlib v1.26.4) | CVE baru di base image |
| `trivy-iac-report.json` | 14 misconfig (Jun 25) | **2 Low** (AWS-0089 S3 logging) | Bukti remediasi Day 23 bekerja |
| `trivy-k8s-post-fix-report.json` | 0 (Jul 10) | 0 (prod manifests) | Di-scope ke manifest produksi (exclude chaos/redteam) |
| `semgrep-report.json` | 1 (Jun 14) | 1 Medium (TLS) | Konsisten |
| `checkov-k8s-post-fix-report.json` | 0 failed | 0 failed | Bersih |
| `checkov-report.json` | 15 failed (Jun 24) | — | CLI checkov gagal di Python 3.14 → di-rename jadi `.stale` |

**Temuan penting regenerasi:**
- IaC: dari **1 Critical + 6 High + 5 Medium → 2 Low** (remediasi S3/SG Day 23 terbukti di scan ulang)
- Image: **CVE-2026-39822 (HIGH, stdlib v1.26.4 → fix 1.26.5)** — action item rebuild image
- go.mod: hanya **GO-2026-5932 (INFO)** — openpgp deprecated

### 2. Re-sync DefectDojo

```bash
bash securebank-api/security/defectdojo/upload-scans.sh
```

Import 6 report ke engagement **Q3 Security Audit** → test baru id **16–21**, total test 21, findings 46 → **52**.

| test | scanner | temuan (state terkini) |
|------|---------|------------------------|
| 16 | Trivy FS | 1 Info |
| 17 | Trivy Image | 1 High + 1 Medium |
| 18 | Trivy IaC | 2 Low |
| 19 | Trivy K8s (post-fix) | 0 |
| 20 | Semgrep | 1 Medium |
| 21 | Checkov K8s | 0 |

### 3. Export & Metrik

```bash
GET /api/v2/findings/  (paginasi) → audit-reports/findings-export.json
```

Aggregat DefectDojo: **52 findings (1 Critical, 7 High, 36 Medium, 7 Low, 1 Info)** — mayoritas dari test lama (pre-fix). State terkini (test 16–21): **6 findings (0 Critical, 1 High, 2 Medium, 2 Low, 1 Info)**.

### 4. Generate Report DefectDojo (UI)

Product **SecureBank API** → Findings → Generate Report → **Executive Summary** (HTML). Output disimpan sebagai evidence lokal; dokumentasi metode ada di catatan ini.

### 5. `draft-q3.md`

Buat `securebank-api/security/audit-reports/draft-q3.md`:
- Executive Summary + metrik ringkas (state terkini: 6 temuan, 0 Critical)
- Tabel tren remediasi (before → after per area)
- Key findings terkini (CVE stdlib HIGH, TLS, S3 logging, openpgp)
- Remediation evidence (S3/SG, K8s hardening, crypto, deps, defense-in-depth)
- Residual risk + metodologi tools

---

## 🧪 Hasil Checklist

| Checklist | Hasil |
|-----------|-------|
| Report DefectDojo berhasil digenerate | ✅ Executive Summary (HTML) + findings-export.json |
| Draf laporan awal berisi metrik keamanan yang jelas | ✅ draft-q3.md (52 aggregate / 6 current, tren, residual risk) |

---

## 📝 Catatan Teknis

### DefectDojo data stale ≠ error
Temuan DefectDojo adalah snapshot saat import. Karena CI `upload-defectdojo` di-skip (`DEFECTDOJO_URL` kosong — DefectDojo lokal), data terakhir Day 47. Re-sync manual = cara refresher-nya. **Duplicate test** muncul (15 → 21 test) karena import membuat test baru per run — dedup DefectDojo tidak menyatukan antar run.

### checkov CLI di Python 3.14
`pip install checkov` sukses tapi entrypoint `checkov` tidak jalan (`No module named checkov.__main__`), image docker `bridgecrewio/checkov` sudah tidak ada di Docker Hub. Workaround: rename report lama jadi `.stale` agar tidak tersaji sebagai state terkini, dan andalkan trivy config untuk coverage IaC (akurat).

### trivy fs dan file gitignored
`trivy fs` (default scanner secret) men-scan **semua file termasuk yang gitignored** (`cosign.key`, `aws-credentials.yaml`, `.env`) → report berisi secret file tersebut. Untuk report SCA yang di-commit, jalankan `--scanners vuln` (vuln only) supaya bersih. CI aman karena checkout GitHub tidak berisi file gitignored.

### Severity `UNKNOWN` dari trivy (GO-2026-5932)
Trivy memakai severity dari vendor lain (GOVULNDB) untuk beberapa vuln; GO-2026-5932 ber-Severity UNKNOWN → tidak menggagalkan gate CI (threshold CRITICAL,HIGH), tapi tetap tercatat.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Checkov CLI tidak jalan (Python 3.14) + image Docker tidak ada | Rename `checkov-report.json` → `.stale`; IaC pakai trivy config (2 Low) |
| `trivy config` tidak menerima multiple target file | Scan dir + `--skip-dirs k8s/chaos` + `--skip-files redteam-pod.yaml` |
| trivy fs menangkap secret file gitignored (cosign.key, .env) | Jalankan ulang `--scanners vuln` (tanpa secret scanner) untuk report yang di-commit |
| Data DefectDojo stale (Day 47) + duplicate test pada re-import | Re-sync manual + catat test id 16-21 sebagai "state terkini" |

---

## 📤 Output Hari Ini

- [x] Report scanner di-regenerate ke state terkini (trivy fs/image/iac/k8s, semgrep, checkov-k8s)
- [x] `checkov-report.json` stale → di-rename `.stale` (tidak tersaji sebagai current)
- [x] DefectDojo re-sync: 6 test baru (16–21), 52 findings total
- [x] `securebank-api/security/audit-reports/findings-export.json` — evidence export
- [x] `securebank-api/security/audit-reports/draft-q3.md` — draf laporan eksekutif
- [x] Temuan baru teridentifikasi: CVE-2026-39822 (HIGH, stdlib 1.26.4)

---

## 💡 Lessons Learned

### 1. Laporan audit sebaiknya tidak dibangun dari data yang bisa saja stale
Snapshot DefectDojo dari bulan lalu menyajikan "1 Critical + 7 High" yang sudah lama ter-remediasi. Re-sync scan sebelum membuat laporan mengubah cerita secara drastis: IaC 14 → 2 Low, K8s 0, deps 0. **Selalu validate data sebelum menyimpulkan.**

### 2. Scan ulang adalah bukti remediasi paling jujur
Daripada mengklaim "sudah diperbaiki", scan ulang menunjukkan angkanya: trivy config turun 14 → 2. Ini bukti yang bisa dibawa ke manajemen — bukan sekadar narasi.

### 3. Scanner lokal ≠ scanner CI (cakupan file)
`trivy fs` lokal men-scan file gitignored (termasuk private key & secret) → berisiko report yang di-commit mengandung secret. Ingat beda cakupan: CI hanya melihat tracked files. Filter `--scanners vuln` / `--skip-dirs` untuk artefak publik.

### 4. Data aggregate ≠ data current
52 findings aggregate di DefectDojo terlihat menakutkan, padahal 46 di antaranya temuan lama yang sudah ter-remediasi. Laporan eksekutif harus memisahkan "aggregate history" dari "current state" supaya tidak menyesatkan.

### 5. Tooling stale itu sendiri temuan
checkov yang tidak bisa jalan di Python 3.14, image Docker yang tak ada lagi — ini reminder bahwa rantai tooling juga perlu pemeliharaan (seperti temuan node NotReady di Day 52).

---

## 🔗 Referensi

- [Day 46](hari-46.md) — DefectDojo setup
- [Day 47](hari-47.md) — DefectDojo API integration + upload-scans.sh
- [Day 51](hari-51.md) — CSPM Remediation (bukti S3/SG fix)
- [Day 33](hari-33.md) — K8s hardening (trivy-k8s-post-fix 0)
- [DefectDojo API](https://defectdojo.github.io/django-DefectDojo/) — findings, import-scan

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Scan ulang membuktikan remediasi; temuan baru image stdlib |
| Pemahaman materi | 4 | Data pipeline DefectDojo, dedup/duplicate, cakupan scanner |
| Progres sesuai target | 5 | Re-sync + export + draft laporan semua tuntas |

---

## ➡️ Rencana Besok

- [ ] **Day 56: Dokumen Eksekutif (PDF)** — perluas draft-q3.md → dokumen formal (metodologi, temuan, mitigasi, sisa risiko) + export PDF (Pandoc)

---

*[← Hari 54](hari-54.md) | [Hari 56 →](hari-56.md)*