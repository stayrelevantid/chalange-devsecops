# 🔍 Retrospektif Fase 2 — Infrastructure as Code & Container Security

**📅 Periode:** Hari 16–30  
**📅 Tanggal Retrospektif:** 2026-07-06  
**⏱️ Total Waktu Belajar:** ~15 hari  
**👤 Ditulis oleh:** Muhammad Indragiri  

---

## 📊 Ringkasan Progres

| Hari | Topik | Status |
|------|-------|--------|
| 16 | Dockerfile Multi-stage | ✅ |
| 17 | Container Image Scan | ✅ |
| 18 | Dockerfile Hardening | ✅ |
| 19 | Image Signing (Cosign) | ✅ |
| 20 | Terraform Setup + IaC Scan (Checkov) | ✅ |
| 21 | IaC Scan (tfsec/Trivy) | ✅ |
| 22 | IaC di Pipeline | ✅ |
| 23 | IaC Remediation | ✅ |
| 24 | DAST Setup (OWASP ZAP) | ✅ |
| 25 | DAST di Pipeline | ✅ |
| 26 | DAST Remediation | ✅ |
| 27 | AI-Assisted IaC Fix | ✅ |
| 28 | Compliance as Code (InSpec) | ✅ |
| 29 | Pipeline Consolidation | ✅ |
| 30 | Dokumentasi Fase 2 | ✅ |

**Selesai: 15/15** ✅

---

## 🟢 Start — Hal yang Ingin Mulai Dilakukan

- **Kubernetes deployment**: Deploy distroless image ke k3d cluster, konfigurasi SecurityContext (Fase 3)
- **Runtime monitoring**: Install Falco untuk deteksi anomali di container runtime (Fase 3)
- **Policy enforcement**: OPA Gatekeeper untuk admission control di K8s (Fase 3)
- **Rate limiter**: Carry-over dari Fase 1 (DREAD 9.6) — mungkin bisa pakai Ingress rate limit atau middleware
- **Vulnerability dashboard**: Integrate semua scan reports ke DefectDojo (Fase 4)

---

## 🔴 Stop — Hal yang Perlu Dihentikan

- **Multiple workflow files**: `infra.yml` terpisah dari `ci.yml` hanya bikin duplikasi dan maintenance overhead. Konsolidasi ke satu pipeline = satu source of truth
- **Asumsi 0 findings = production-ready**: AI review Day 27 menemukan 6 improvements yang scanner lewati (variable validation, IAM Condition, S3 lifecycle). Stop percaya bahwa scanner hijau = aman
- **`gem install` untuk CLI tools**: InSpec 7.x gem hanya ship library, bukan CLI. Selalu cek apakah package punya CLI executable sebelum install

---

## 🟡 Continue — Hal yang Sudah Baik & Perlu Diteruskan

- **Multi-stage Docker build**: golang:alpine → distroless = 44x image reduction. Pattern ini wajib dipertahankan untuk semua container
- **8-layer container hardening**: non-root, read-only, tmpfs, no-new-privileges, cap_drop ALL, resource limits. Best practice tetap dipakai
- **Cosign image signing**: Supply chain security — semua image yang deploy harus di-verify
- **Dual IaC scanner**: Checkov (policy detail) + Trivy (severity gate) complementary, keduanya dipakai
- **Security gate pattern**: `needs` semua scan jobs = single checkpoint untuk Branch Protection. Lanjut dipakai di Fase 3 dengan tambahan K8s scanner
- **`continue-on-error` untuk false positive**: ZAP 10049 tidak boleh block pipeline, tapi tidak boleh di-ignore begitu saja. Document justification

---

## 🧠 Top 5 Pelajaran Terpenting dari Fase 2

1. **Distroless = attack surface minimal** — Image turun dari ~350MB ke 7.97MB (44x). No shell, no package manager, no utilities. Attacker yang dapat RCE tidak bisa pivot. Trade-off: no HEALTHCHECK (no shell) → pakai K8s liveness probe di Fase 3.

2. **Checkov + Trivy IaC complementary** — Checkov untuk policy-as-code detail (102 checks spesifik Terraform), Trivy untuk multi-format severity gate (CRITICAL/HIGH/MEDIUM). Keduanya dipakai bersamaan, bukan dipilih salah satu.

3. **ZAP false positive butuh dokumentasi** — Rule 10049 (Cacheability) adalah false positive untuk API dengan `Cache-Control: no-store`. `continue-on-error: true` solusinya, tapi justification harus dicatat. Tanpa dokumentasi, next developer akan bingung kenapa DAST "fail" tapi pipeline hijau.

4. **Pipeline consolidation = maintainability** — Dari 2 workflow files (ci.yml + infra.yml) jadi 1 (ci.yml). 8 jobs, 1 file, 1 place to debug, 1 source of truth. Security gate pattern bikin "apakah pipeline aman?" jadi pertanyaan yes/no.

5. **InSpec gem ≠ CLI** — `gem install inspec` (v7.x) hanya install library Ruby. Untuk CLI, butuh Chef Workstation installer via `brew install --cask`. InSpec juga butuh `--chef-license accept` di first run. Lesson: selalu cek apakah package punya CLI executable.

---

## 🏆 Pencapaian Terbaik

**Pipeline 8/8 green** dalam satu file workflow. 7 scanner paralel + security gate. Dari 2 file jadi 1.

**Distroless image 7.97MB** — 44x lebih kecil dari alpine-builder naive. Tidak ada shell, tidak ada package manager, attack surface minimal.

**17 security fixes dalam 15 hari**: multi-stage Docker, 8-layer hardening, Cosign signing, Terraform IaC (102/0 Checkov, 0 Trivy), DAST remediasi (8 security headers), AI-assisted IaC fix (6 improvements), pipeline consolidation.

Momen "aha!": **Go binary approach untuk DAST di CI** (Day 25). Dibanding build Docker image (~30s), compile Go binary langsung (~5s) dengan hasil ZAP scan identik. Lebih cepat, lebih simple, hasil sama.

---

## 😓 Tantangan Terbesar

| Tantangan | Cara Mengatasi |
|-----------|----------------|
| Gitleaks `--no-gitignore` RED dari Day 14-22 | Flag tidak pernah ada. Hapus flag, gunakan `[allowlist]` di `.gitleaks.toml` |
| ZAP rule 10049 WARN di CI (pipeline RED) | `continue-on-error: true` + dokumentasi false positive untuk API |
| `go get -u` timeout di slow network | Edit `go.mod` langsung, lalu `go mod tidy` |
| `gem install inspec` tidak ada CLI (v7.x) | `brew install --cask chef/chef/inspec` (pkg installer butuh sudo) |
| Docker `--network host` tidak work di macOS | Pakai `host.docker.internal` untuk ZAP container reach host |
| `.trivyignore` menyebabkan false green di CI | Hapus `.trivyignore` dari repo, fix root cause di Terraform |
| `http.HandleFunc` (global mux) tidak applying middleware ke 404 | Refactor ke `http.NewServeMux()` + `SecurityHeadersHandler(mux)` wrapper |

---

## 🔗 Koneksi ke Fase Berikutnya

- **Docker image (signed + hardened)** → akan di-deploy ke K8s cluster di Hari 31
- **Terraform infra** → menjadi target IaC scan lanjutan & CSPM di Fase 4
- **Pipeline gabungan (8 jobs)** → akan ditambah K8s manifest scanner (kubesec/Checkov K8s) di Fase 3
- **Cosign verification** → K8s admission control bisa verify signature sebelum deploy
- **Threat model carry-over**: Rate limiter (DREAD 9.6), user-scoped authorization (8.4), RBAC (6.6) masih unmitigated
- **InSpec profile** → pattern "compliance as code" akan diteruskan dengan OPA Rego policies di K8s

---

## 📈 Skor Diri (Jujur!)

| Aspek | Skor (1–10) | Catatan |
|-------|-------------|---------|
| Pemahaman Docker & container | 9 | Multi-stage, distroless, 8-layer hardening, Cosign |
| Kemampuan Terraform dasar | 8 | 10 files, variable validation, outputs, KMS, lifecycle |
| Pemahaman IaC scanning | 9 | Checkov 102/0, Trivy IaC, inline skip, SARIF |
| Pemahaman DAST | 8 | ZAP baseline, false positive handling, security headers |
| Kemampuan menulis pipeline kompleks | 9 | 8 jobs paralel, needs, continue-on-error, security gate |
| Konsistensi belajar harian | 7 | Beberapa gap antar hari, tapi 15/15 selesai |

---

## 📝 Catatan Bebas

Fase 2 adalah fase yang paling "hands-on" — banyak tools baru (Docker, Cosign, Terraform, Checkov, ZAP, InSpec) dalam 15 hari. Yang paling challenging: ZAP false positive (10049) yang membuat pipeline RED dari Day 25-26 sampai `continue-on-error` diimplementasi.

Yang paling satisfying: melihat 8/8 jobs green di pipeline consolidation (Day 29). Dari 2 file workflow jadi 1, dari 5 jobs jadi 8, dengan security gate sebagai final checkpoint. Clean, simple, effective.

Distroless image 7.97MB itu game-changer. Tidak hanya soal ukuran, tapi soal attack surface — tidak ada shell, tidak ada package manager, attacker yang dapat RCE tidak bisa apa-apa. Trade-off (no HEALTHCHECK) akan di-solve di Fase 3 dengan K8s liveness probe.

---

*[← Retrospektif Fase 1](fase-1-retrospektif.md) | [Retrospektif Fase 3 →](fase-3-retrospektif.md)*