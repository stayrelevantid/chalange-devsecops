# Hari 59 — Persiapan CDP Exam Simulation

**📅 Tanggal:** 2026-08-16
**⏱️ Durasi Belajar:** ~120 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Menyusun strategi, checklist skill, dan cheatsheet untuk simulasi ujian CDP
- [x] Mengintegrasikan 15 DevSecOps Best Practices dari Practical DevSecOps
- [x] Membuat `docs/cdp-exam-guide.md` sebagai panduan reusable
- [x] Menetapkan aturan ujian (3 jam, boleh Google, dilarang menyalin dari `main`)

## ✅ Yang Berhasil Dikerjakan

Dokumentasi persiapan ujian CDP selesai dalam satu file: [`docs/cdp-exam-guide.md`](../../docs/cdp-exam-guide.md).

Isinya:

| Bagian | Konten |
|--------|--------|
| Format ujian | Scope 3 jam, aturan, branch `cdp-exam` terpisah |
| Strategi alokasi waktu | 180 menit: perencanaan 15m → CI 45m → secret+SCA 25m → SAST 20m → Dockerfile 25m → K8s 30m → verifikasi 20m |
| 15 Best Practices | Shift Left, Automation, Continuous Testing, Risk Management, Integrate Security Tools, Collaboration, Secure Coding Standards, Access Controls, Monitoring, Training, **Policy as Code**, Threat Modeling, Incident Response, Immutable Infrastructure, Observability — dipetakan ke tindakan ujian |
| Checklist skill | GitHub Actions, Gitleaks, Trivy, Semgrep, Dockerfile multi-stage, K8s SecurityContext, bash |
| Cheatsheet | Snippet siap pakai untuk workflow, gitleaks, trivy, semgrep, Dockerfile, deployment/service aman, kubeconform |
| Tips & jebakan | `kubectl apply --dry-run` butuh cluster → kubeconform; `fetch-depth: 0` untuk gitleaks; pin golang patch; `labels` vs `commonLabels`; cache key; konsistensi artifact path |
| Persiapan & evaluasi | Checklist sebelum ujian, strategi verifikasi tanpa mencontek, refleksi pasca ujian |

## 📝 Catatan Teknis

Sumber best practices: [Practical DevSecOps — Top 15 DevSecOps Best Practices](https://www.practical-devsecops.com/devsecops-best-practices/).

Pola yang divalidasi dari 60 hari sebelumnya:

```text
golang:1.26.6-alpine   # builder pin versi patch (tutup CVE stdlib)
kubeconform -strict    # validasi manifest TANPA cluster
labels: (bukan commonLabels)  # Kustomize modern
fetch-depth: 0         # agar Gitleaks scan seluruh history
```

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Runner GitHub tidak punya cluster K8s | Validasi manifest pakai `kubeconform -strict`, bukan `kubectl apply --dry-run` |
| CVE stdlib muncul di image | Pin `golang:1.26.6-alpine` di builder |
| Squash merge membuat ancestry divergen | Merge commit wajib untuk promotion `develop → staging → main` |

## 📤 Output Hari Ini

- [x] `docs/cdp-exam-guide.md` — panduan persiapan ujian lengkap
- [x] Best practices terintegrasi ke strategi ujian
- [x] Cheatsheet sintaks siap pakai
- [x] Checklist persiapan & evaluasi diri

## 💡 Pelajaran Baru

- Persiapan ujian yang baik = strategi waktu + checklist skill + cheatsheet + daftar jebakan.
- Best practices tidak cukup sebagai daftar — harus dipetakan ke tindakan konkret di pipeline.
- Jebakan teknis (kubeconform, golang patch, labels) lebih berharga dicatat daripada dihafal ulang saat ujian.

## 🔗 Referensi

- [Panduan Ujian CDP](../../docs/cdp-exam-guide.md)
- [Hari 58](hari-58.md) — CI/CD multi-environment (dasar materi ujian)
- [Practical DevSecOps Best Practices](https://www.practical-devsecops.com/devsecops-best-practices/)

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Persiapan terasa terarah dan praktis |
| Pemahaman materi | 4 | Materi ujian tersusun rapi dalam satu referensi |
| Progres sesuai target | 5 | Guide + best practices + cheatsheet selesai |

## ➡️ Rencana Besok

- [ ] **Day 60:** Project Showcase — diagram arsitektur E2E + publikasi final

*[← Hari 58](hari-58.md) | [Hari 60 →](hari-60.md)*
