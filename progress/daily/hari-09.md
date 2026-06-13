# Hari 09 — SAST di Pipeline (Quality Gate)

**📅 Tanggal:** 2026-06-13  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Job SAST berjalan di pipeline
- [x] Pipeline gagal karena ada Semgrep findings (quality gate)
- [x] SARIF report ter-integrasi dengan GitHub Security tab

---

## ✅ Yang Berhasil Dikerjakan

- Tambah job `sast-scan` ke `.github/workflows/ci.yml` menggunakan `semgrep/semgrep-action@v1`
- Konfigurasi: `p/golang` + `p/owasp-top-ten` + custom `.semgrep.yml`
- Pipeline sekarang punya 4 job paralel: Build & Test, Secret Scan, SCA Scan, SAST Scan
- SAST Scan menemukan **2 blocking findings** → pipeline gagal (MERAH) sebagai quality gate
- SARIF upload ke GitHub Security tab via `codeql-action/upload-sarif@v3`

---

## 📝 Catatan Teknis

```yaml
sast-scan:
  name: SAST Scan (Semgrep)
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Run Semgrep
      uses: semgrep/semgrep-action@v1
      with:
        config: >-
          p/golang
          p/owasp-top-ten
          .semgrep.yml

    - name: Upload SARIF
      if: always()
      uses: github/codeql-action/upload-sarif@v3
      with:
        sarif_file: semgrep.sarif
```

### Keputusan desain:
- **`semgrep/semgrep-action@v1`** — action resmi dari Semgrep, gratis untuk public repo
- **Config: `p/golang` + `p/owasp-top-ten` + `.semgrep.yml`** — built-in rules + OWASP + custom rules
- **`.semgrep.yml`** di root repo — konfigurasi konsisten dengan `.gitleaks.toml` dan `.trivyignore`
- **`semgrep-action`** otomatis gagalkan pipeline jika ada findings — tidak perlu flag tambahan
- **SARIF upload** ke GitHub Security tab — findings bisa dilihat langsung di GitHub UI

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| — | Tidak ada hambatan — `semgrep-action` bekerja smooth, gratis untuk public repo |

---

## 📤 Output Hari Ini

- [x] Job `sast-scan` ditambahkan ke `ci.yml`
- [x] Pipeline expected gagal (RED) karena 2 Semgrep findings
- [x] Semgrep findings: `no-md5-usage` (ERROR) + `no-http-listen-without-tls` (WARNING)
- [x] SARIF report di-upload ke GitHub Security tab
- [x] Pipeline punya 4 job: Build & Test, Secret Scan, SCA Scan, SAST Scan

---

## 💡 Pelajaran Baru

- **`semgrep-action` gratis untuk public repo** — tidak butuh SEMGREP_APP_TOKEN untuk scan dasar
- **Semgrep otomatis gagalkan pipeline jika ada findings** — berbeda dengan Trivy yang perlu `exit-code: 1`
- **SARIF (Static Analysis Results Interchange Format)** memungkinkan findings ditampilkan di GitHub Security tab — integrasi yang bagus antara SAST tool dan platform
- **Config bisa combine multiple sources** — `p/golang` (built-in), `p/owasp-top-ten` (OWASP), dan `.semgrep.yml` (custom) semua jalan bersamaan
- **Quality gate SAST** berbeda dengan SCA — SAST gagal karena kode yang kita tulis, bukan karena dependensi

---

## 🔗 Referensi

- [Semgrep GitHub Action](https://github.com/semgrep/semgrep-action)
- [Semgrep Rules Registry](https://semgrep.dev/explore)
- [GitHub SARIF Support](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/sarif-support-for-code-scanning)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline makin lengkap — 4 lapis pertahanan |
| Pemahaman materi | 4 | SAST di CI straightforward |
| Progres sesuai target | 5 | Lancar, tanpa hambatan |

---

## ➡️ Rencana Besok

- [ ] Fix MD5 → SHA-256 di `pkg/crypto/hash.go` (Day 10 - SAST Remediation)
- [ ] Pertimbangkan TLS untuk HTTP server
- [ ] Pipeline kembali hijau setelah semua findings di-remediasi

---

*[← Hari 08](hari-08.md) | [Hari 10 →](hari-10.md)*