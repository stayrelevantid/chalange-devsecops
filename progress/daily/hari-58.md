# Hari 58 — CI/CD Multi-Environment & Release Governance

**📅 Tanggal:** 2026-08-09
**⏱️ Durasi Belajar:** ~180 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Memperbaiki quality gate CI dengan container image scan
- [x] Membuat CD workflow DEV → UAT → PRE-PROD → security approval → PROD
- [x] Menambahkan branch protection, GitHub Environments, dan Kustomize overlays
- [x] Menunda wipe lab CDP agar security evidence tetap tersedia

## ✅ Yang Berhasil Dikerjakan

- `.github/workflows/ci.yml` kini menjalankan build/test, Gitleaks, Trivy SCA, Semgrep, ZAP, Checkov, Trivy IaC, dan image scan.
- `.github/workflows/cd-deploy.yml` memakai `workflow_dispatch`, approval environment, security approval manual, dan post-deployment verification.
- Manifest Kubernetes dipisah menjadi Kustomize base serta overlay `dev`, `uat`, `preprod`, dan `prod`.
- Branch `develop`, `staging`, `uat`, dan `main` dibuat/protected; lima GitHub Environment dikonfigurasi.

## 📝 Catatan Teknis

```bash
kustomize build securebank-api/k8s/overlays/dev
kustomize build securebank-api/k8s/overlays/uat
kustomize build securebank-api/k8s/overlays/preprod
kustomize build securebank-api/k8s/overlays/prod
```

Deployment masih berupa simulasi karena runner GitHub tidak dapat mengakses k3d lokal. Cosign di-skip karena private key belum diberikan; tidak ada key/secret yang dibuat otomatis.

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Runner GitHub tidak dapat mengakses k3d lokal | Render Kustomize, dry-run Kubernetes, dan container health check |
| Cosign membutuhkan private key | Signing di-skip; key/secret tetap harus diberikan manual |

## 📤 Output Hari Ini

- [x] Quality gate CI dengan image scan
- [x] CD workflow multi-environment dengan approval
- [x] Kustomize base dan overlay per environment
- [x] Branch protection dan GitHub Environments

## 💡 Pelajaran Baru

- Quality gate perlu memeriksa image final, bukan hanya source code.
- Approval release harus menjadi kontrol nyata, bukan sekadar log output.
- Deployment simulation harus dijelaskan secara transparan.
- Wipe lab ditunda agar evidence security yang sudah ada tidak hilang.

## 🔗 Referensi

- [Hari 59](hari-59.md) — persiapan ujian CDP berikutnya
- [CI workflow](../../.github/workflows/ci.yml)
- [CD workflow](../../.github/workflows/cd-deploy.yml)

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Release flow lebih realistis |
| Pemahaman materi | 5 | Quality gate, approval, branch governance |
| Progres sesuai target | 5 | CI/CD selesai; persiapan ujian dilanjutkan besok |

## ➡️ Rencana Besok

- [ ] **Day 59:** persiapan CDP exam simulation dan finalisasi scope latihan

*[← Hari 57](hari-57.md) | [Hari 59 →](hari-59.md)*
