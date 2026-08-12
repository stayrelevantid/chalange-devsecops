# Hari 59 — CDP Exam Sim: Execution

**📅 Tanggal:** 2026-08-09
**⏱️ Durasi Belajar:** ~180 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Memperbaiki quality gate CI dengan container image scan
- [x] Menyusun CD workflow manual dengan alur DEV → UAT → PRE-PROD → PROD
- [x] Menambahkan approval GitHub Environment dan security approval manual terpisah
- [x] Memisahkan konfigurasi Kubernetes dengan Kustomize base + overlay
- [x] Membuat branch `develop`, `staging`, dan `uat` untuk promotion flow
- [x] Mengaktifkan branch protection pada seluruh branch promotion

---

## ✅ Yang Berhasil Dikerjakan

### 1. Quality Gate CI

Workflow `.github/workflows/ci.yml` kini memiliki blocking checks:

| Check | Gate |
|-------|------|
| Go build/test | Build, race detection, coverage wajib lulus |
| Gitleaks | Secret terdeteksi = gagal |
| Trivy SCA | CRITICAL/HIGH = gagal |
| Semgrep | SAST job wajib lulus |
| Checkov/Trivy IaC | Misconfiguration sesuai threshold = gagal |
| Trivy image | CRITICAL/HIGH pada image = gagal |
| Security Gate | Single required check untuk promotion |

ZAP tetap dijalankan sebagai DAST evidence dengan false-positive policy yang sudah ada.

### 2. CD Multi-Environment

`.github/workflows/cd-deploy.yml` memakai `workflow_dispatch` dan alur berurutan:

```text
build image + Trivy scan
  -> DEV approval + deployment simulation
  -> UAT approval + manifest validation
  -> PRE-PROD approval + manifest validation
  -> SECURITY APPROVAL manual terpisah
  -> PROD approval + manifest validation
  -> post-deployment verification
```

Deployment menggunakan simulasi karena GitHub-hosted runner tidak dapat mengakses k3d lokal: render Kustomize, `kubectl apply --dry-run`, pull image, lalu health check container di DEV.

### 3. Branch Protection & Environments

Branch `develop`, `staging`, `uat`, dan `main` sudah dibuat dan diproteksi dengan PR review, required status checks, linear history, conversation resolution, serta block force-push/deletion. Environment `dev`, `uat`, `preprod`, `security-approval`, dan `prod` memiliki required reviewer `yogi-indragiri`; `prod` memiliki wait timer 5 menit.

### 4. Kustomize

`securebank-api/k8s/base/` menyimpan manifest reusable. Overlay `dev`, `uat`, `preprod`, dan `prod` menetapkan suffix nama, replica count, dan `APP_ENV` tanpa menyalin seluruh manifest.

---

## 📝 Catatan Teknis

```bash
for env in dev uat preprod prod; do
  kustomize build "securebank-api/k8s/overlays/$env"
done

python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/ci.yml"))'
python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/cd-deploy.yml"))'
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| GitHub-hosted runner tidak bisa akses k3d lokal | Gunakan deployment simulation: render/validate manifest + container health check |
| Cosign signing memerlukan private key manual | Skip signing; tidak membuat key atau secret baru |
| Day 58 lab wipe berisiko menghapus security evidence | Day 58 ditandai skipped; konfigurasi utama dipertahankan |

---

## 📤 Output Hari Ini

- [x] `.github/workflows/ci.yml` — image build + Trivy blocking scan + consolidated security gate
- [x] `.github/workflows/cd-deploy.yml` — promotion pipeline dengan approval bertingkat
- [x] `securebank-api/k8s/base/` + `overlays/{dev,uat,preprod,prod}/`
- [x] `securebank-api/Dockerfile` — Go builder `1.26.5`
- [x] GitHub Environments + branch protection rules

---

## 💡 Pelajaran Baru

- Quality gate harus eksplisit membedakan check yang memblokir rilis dan evidence yang hanya informatif.
- Environment approval menjaga pemisahan tugas antara build, validasi, dan release.
- Kustomize overlay mengurangi copy-paste dan membuat perbedaan environment terlihat jelas.
- Tidak semua deployment bisa benar-benar dijalankan dari CI hosted; simulasi yang jujur lebih baik daripada klaim deployment palsu.

---

## 🔗 Referensi

- [Day 57](hari-57.md) — AI review executive summary
- [Day 58](hari-58.md) — lab setup dilewati
- [CI workflow](../../.github/workflows/ci.yml)
- [CD workflow](../../.github/workflows/cd-deploy.yml)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline mulai menyerupai release process nyata |
| Pemahaman materi | 5 | Quality gate, environment approval, dan promotion flow |
| Progres sesuai target | 5 | CI/CD, branches, environments, dan protection selesai |

---

## ➡️ Rencana Besok

- [ ] **Day 60: Project Showcase** — diagram arsitektur E2E, publikasi, dan cleanup resource lab

---

*[← Hari 58](hari-58.md) | [Hari 60 →](hari-60.md)*
