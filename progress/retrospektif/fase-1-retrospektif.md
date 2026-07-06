# 🔍 Retrospektif Fase 1 — Secure SDLC & Application Security

**📅 Periode:** Hari 1–15  
**📅 Tanggal Retrospektif:** 2026-07-06  
**⏱️ Total Waktu Belajar:** ~15 hari  
**👤 Ditulis oleh:** Muhammad Indragiri  

---

## 📊 Ringkasan Progres

| Hari | Topik | Status |
|------|-------|--------|
| 01 | Setup Repo & Golang API | ✅ |
| 02 | Pipeline CI/CD Dasar | ✅ |
| 03 | Secret Scanning (Gitleaks) | ✅ |
| 04 | Secret Scan di Pipeline + Remediasi | ✅ |
| 05 | SCA Setup (Trivy FS) | ✅ |
| 06 | SCA di Pipeline (Gate) | ✅ |
| 07 | SCA Remediation | ✅ |
| 08 | SAST Setup (Semgrep) | ✅ |
| 09 | SAST di Pipeline (Gate) | ✅ |
| 10 | SAST Remediation | ✅ |
| 11 | AI-Assisted Code Audit | ✅ |
| 12 | Pipeline Optimization | ✅ |
| 13 | Threat Modeling (STRIDE) | ✅ |
| 14 | Intentional Vuln Test | ✅ |
| 15 | Dokumentasi Fase 1 | ✅ |

**Selesai: 15/15** ✅

---

## 🟢 Start — Hal yang Ingin Mulai Dilakukan

- **Container security**: Dockerfile multi-stage, distroless, image signing (Fase 2)
- **IaC scanning**: Terraform + Checkov/Trivy untuk infrastructure (Fase 2)
- **DAST**: OWASP ZAP untuk dynamic testing (Fase 2)
- **Rate limiter**: Fase 1 threat model menemukan DREAD score 9.6 untuk no rate limiting — perlu mulai implement

---

## 🔴 Stop — Hal yang Perlu Dihentikan

- **Asumsi scanner cukup**: Trivy + Semgrep hanya menemukan pattern, bukan business logic flaw. Stop bergantung 100% pada scanner — harus manual/AI review juga
- **Ignore `.gitignore` sebagai security boundary**: `.gitignore` hanya mencegah accidental commit, `git add -f` bypass. Jangan anggap `.gitignore` = security control

---

## 🟡 Continue — Hal yang Sudah Baik & Perlu Diteruskan

- **Parallel jobs di CI**: 4 jobs paralel lebih cepat dari sequential. Lanjut ke 8 jobs paralel di Fase 2
- **Pipeline caching**: 4 layer cache (Go modules, Gitleaks, Trivy DB, Semgrep) hemat 30-40% per run. Wajib dipertahankan
- **Threat modeling (STRIDE + DREAD)**: Systematic framework untuk identifikasi ancaman. Lanjut di Fase 3 untuk K8s topology
- **Custom rules**: `.semgrep.yml` custom rule (no-md5-usage) menambah coverage. Tambah rules baru untuk pattern yang spesifik project

---

## 🧠 Top 5 Pelajaran Terpenting dari Fase 1

1. **Scanner != Production-Ready** — Trivy menemukan 4 CVE, Semgrep 2 findings, tapi AI audit menemukan 5 kerentanan lain yang keduanya lewati (input validation, auth, security headers, body limit, audit logging). Business logic gap tidak terdeteksi oleh pattern matching.

2. **`.gitignore` bukan security boundary** — Gitleaks default menghormati `.gitignore`, jadi file yang di-`git add -f` tetap di-skip. Selalu validasi bahwa scanner benar-benar men-scan semua tracked files.

3. **Pipeline caching itu penting** — Dari ~30-50 detik per run untuk download dependencies, turun ke ~5-10 detik dengan 4 layer caching. ROI tinggi untuk effort minimal.

4. **JWT Authentication ≠ Authorization** — Day 11 implementasi JWT (authentication: "siapa kamu?"), tapi threat modeling Day 13 menunjukkan authorization ("apa yang boleh kamu lakukan?") belum ada. Siapa saja dengan JWT valid bisa lihat saldo akun siapa saja.

5. **Threat modeling mengisi gap scanner** — STRIDE + DREAD menemukan 18 threats yang tidak ditemukan scanner otomatis. 6 sudah di-mitigasi, 4 partial, 8 belum. Tanpa threat modeling, 8 threat ini mungkin tidak teridentifikasi.

---

## 🏆 Pencapaian Terbaik

Membangun pipeline CI/CD dengan 4 parallel quality gates yang memblokir kode tidak aman — **17 security fixes** dalam 15 hari, dari 0 ke 24 unit tests all passing, 0 CVE. Pipeline GREEN dari Day 14 onward.

Momen "aha!": **Intentional Vulnerability Test (Day 14)** — commit fake secret untuk validasi bahwa Gitleaks CI benar-benar menangkapnya. Hasilnya: pipeline RED, proof bahwa defense layer berfungsi.

---

## 😓 Tantangan Terbesar

| Tantangan | Cara Mengatasi |
|-----------|----------------|
| `go get -u` timeout di slow network | Edit `go.mod` langsung, lalu `go mod tidy` |
| Gitleaks `--no-gitignore` flag tidak ada (akumulasi RED dari Day 14) | Hapus flag, gunakan `[allowlist]` di `.gitleaks.toml` |
| ZAP rule 10049 false positive untuk API | `continue-on-error: true` dengan justification |
| Trivy CVE di jwt-go (deprecated) | `replace` directive di `go.mod` → golang-jwt/v5 |
| Semgrep MD5 finding → refactor ke bcrypt | Custom rule `.semgrep.yml` + `golang.org/x/crypto/bcrypt` |

---

## 🔗 Koneksi ke Fase Berikutnya

- **SecureBank API** (kode Go + 25 tests) → akan di-containerize dengan Dockerfile di Hari 16
- **Pipeline CI/CD** (ci.yml 4 jobs) → akan ditambah job scan image, IaC, DAST di Fase 2
- **Semgrep/Trivy output** (JSON_reports) → akan diintegrasikan ke DefectDojo di Fase 4
- **Threat model** (STRIDE + DREAD) → 8 unmitigated threats carry-over ke Fase 2/3 (rate limiter, RBAC, authorization)
- **Security headers** (5 headers) → akan ditambah 3 lagi di Fase 2 Day 26 (Referrer-Policy, Permissions-Policy, CORP)

---

## 📈 Skor Diri (Jujur!)

| Aspek | Skor (1–10) | Catatan |
|-------|-------------|---------|
| Pemahaman konsep Secret Scan | 8 | Gitleaks config, allowlist, CI integration |
| Pemahaman konsep SCA | 8 | Trivy FS, severity gates, CVE remediation |
| Pemahaman konsep SAST | 7 | Semgrep rules, custom rules, severity levels |
| Kemampuan baca pipeline YAML | 9 | Cache, parallel jobs, matrix, needs, continue-on-error |
| Kemampuan Go (dasar) | 7 | HTTP handlers, middleware, testing, JWT |
| Konsistensi belajar harian | 6 | Beberapa gap antar hari (streak 0) |

---

## 📝 Catatan Bebas

Fase 1 adalah fondasi. Tanpa fondasi ini (kode aman, pipeline berjalan, scanner aktif), Fase 2 tidak akan maksimal. Yang paling berkesan: AI audit di Day 11 yang menemukan 5 kerentanan yang scanner lewati. Membuktikan bahwa tools otomatis perlu, tapi tidak cukup — butuh manual/AI review untuk business logic gap.

Threat modeling (STRIDE + DREAD) di Day 13 membuka mata: ada 18 threats yang belum teridentifikasi. 8 masih unmitigated dan akan carry-over ke Fase 2/3. Ini narasi yang berkelanjutan — tidak ada "finish line" di security.

---

*[← Tracker](../tracker.md) | [Retrospektif Fase 2 →](fase-2-retrospektif.md)*