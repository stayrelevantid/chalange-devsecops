# Hari 58 — CI/CD Multi-Environment dengan Auto-Deploy

**📅 Tanggal:** 2026-08-09
**⏱️ Durasi Belajar:** ~180 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Memperbaiki quality gate CI dengan container image scan
- [x] Membuat CD otomatis yang ter-trigger setelah CI sukses (`workflow_run`)
- [x] Menambahkan branch protection, GitHub Environments, dan Kustomize overlays
- [x] Validasi manifest tanpa cluster menggunakan `kubeconform`
- [x] Menunda wipe lab CDP agar security evidence tetap tersedia

## ✅ Yang Berhasil Dikerjakan

### 1. Quality Gate CI

`.github/workflows/ci.yml` menjalankan:

| Check | Gate |
|-------|------|
| Build & Test (Go 1.26) | build + race test + coverage wajib lulus |
| Secret Scan (Gitleaks) | secret terdeteksi = gagal |
| SCA Scan (Trivy) | CRITICAL/HIGH = gagal |
| SAST Scan (Semgrep) | finding HIGH = gagal |
| DAST Scan (OWASP ZAP) | evidence, `-I` false positive tidak memblokir |
| IaC Scan (Checkov) | misconfig = gagal |
| IaC Scan (Trivy) | CRITICAL/HIGH/MEDIUM = gagal |
| Container Build & Image Scan | CVE CRITICAL/HIGH pada image = gagal |
| Promotion Policy | jalur PR harus sesuai `feature/fix → develop → staging → main` |
| Security Gate | konsolidasi semua check di atas |

### 2. CD Otomatis

`.github/workflows/cd-deploy.yml` menggunakan `workflow_run` — CD ter-trigger otomatis ketika **SecureBank CI** selesai sukses, tanpa perlu klik `Run workflow` dan tanpa approval environment:

```yaml
on:
  workflow_run:
    workflows: ["SecureBank CI"]
    types: [completed]
    branches: [develop, staging, main]
```

Mapping deployment berdasarkan branch CI yang sukses:

```text
CI develop sukses → Deploy DEV (simulation)
CI staging sukses → Deploy STAGING (simulation)
CI main sukses    → Deploy PROD (simulation) + Post-deployment verification
```

### 3. Validasi Manifest Tanpa Cluster

Runner GitHub-hosted tidak memiliki Kubernetes API server. `kubectl apply --dry-run` akan mencoba konek ke `localhost:8080` dan gagal. Solusi: gunakan **kubeconform** untuk validasi schema manifest:

```bash
kubeconform -strict -ignore-missing-schemas /tmp/dev.yaml
```

### 4. Kustomize Overlays

Manifest reusable di `securebank-api/k8s/base/`, overlay per environment:

```text
dev
staging
prod
```

Perbedaan antar environment: suffix nama, replica count, dan `APP_ENV`.

### 5. Branch Protection & Environments

Branch `develop`, `staging`, `main` dilindungi:

- Required status checks: `Build & Test (Go 1.26)` + `Security Gate`
- Force push diblokir
- Branch deletion diblokir
- Conversation resolution wajib
- Merge commit diizinkan (untuk menjaga ancestry promotion)

Environment GitHub:

```text
dev
staging
prod
```

## 📝 Catatan Teknis

### Base image Go

Dockerfile builder di-update untuk menutup CVE stdlib:

```dockerfile
FROM golang:1.26.6-alpine AS builder
```

Menutup `CVE-2026-39821` dan `CVE-2026-46600`.

### Repository merge settings

- Merge commit: ✅ aktif
- Rebase merge: ❌ nonaktif
- Squash merge: ⚠️ tersedia via API GitHub, tetapi dilarang untuk PR promotion
- Auto-delete branch: ✅ aktif

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Runner GitHub tidak dapat mengakses k3d lokal | Deployment simulation: render Kustomize + `kubeconform` + container health check |
| `kubectl apply --dry-run` mencoba `localhost:8080` | Ganti ke `kubeconform -strict` (tanpa cluster) |
| CVE stdlib di image | Bump builder ke `golang:1.26.6-alpine` |
| Squash merge membuat ancestry divergen | Merge commit wajib untuk `develop → staging` dan `staging → main` |
| Cosign butuh private key | Signing di-skip; key/secret tidak dibuat otomatis |

## 📤 Output Hari Ini

- [x] Quality gate CI 10 check
- [x] CD otomatis `workflow_run` tanpa approval
- [x] `kubeconform` validation
- [x] Kustomize base + overlay `dev/staging/prod`
- [x] Branch protection + 3 GitHub Environments
- [x] Go builder `1.26.6`

## 💡 Pelajaran Baru

- Quality gate harus memeriksa artefak final (image), bukan hanya source code.
- Auto-deploy yang aman = CD ter-trigger hanya jika CI sukses, dan deployment target dipilih dari branch pemicu.
- Approval tidak selalu harus environment gate; untuk solo project, approval bisa berupa PR checks + merge commit.
- Validasi manifest tidak memerlukan cluster bila memakai schema validator seperti kubeconform.
- Merge commit untuk promotion menjaga ancestry tetap dapat ditelusuri.

## 🔗 Referensi

- [Branching & Merge Policy](../../docs/branching-and-merge-policy.md)
- [Hari 59](hari-59.md) — persiapan CDP exam simulation berikutnya
- [CI workflow](../../.github/workflows/ci.yml)
- [CD workflow](../../.github/workflows/cd-deploy.yml)

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline akhir terasa seperti rilis sungguhan |
| Pemahaman materi | 5 | workflow_run, kubeconform, promotion, merge commit |
| Progres sesuai target | 5 | CI/CD tuntas dan diuji end-to-end |

## ➡️ Rencana Besok

- [ ] **Day 59:** persiapan CDP exam simulation dan finalisasi scope latihan

*[← Hari 57](hari-57.md) | [Hari 59 →](hari-59.md)*
