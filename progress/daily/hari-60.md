# Hari 60 — Project Showcase

**📅 Tanggal:** 2026-08-19
**⏱️ Durasi Belajar:** ~3 jam
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Membuat diagram arsitektur end-to-end (Draw.io + PNG)
- [x] Publikasi pencapaian: blog post showcase skill + draft LinkedIn
- [x] AWS cleanup — hapus resource lab agar biaya kembali nol
- [x] Teardown resource lokal (k3d, DefectDojo, n8n webhook receiver)
- [x] Menutup challenge: update tracker, retrospektif, dan promotion ke `main`

---

## ✅ Yang Berhasil Dikerjakan

### 1. Diagram Arsitektur End-to-End

File dibuat di [`docs/architecture/day-60-architecture.drawio`](../../docs/architecture/day-60-architecture.drawio) dan di-render ke PNG (scale 2) dengan draw.io desktop headless.

Lapisannya:

| Band | Isi |
|------|-----|
| 1. Source & CI/CD | Developer → GitHub (branch protection) → SecureBank CI: Build & Test, Gitleaks, Trivy SCA, Semgrep, ZAP DAST, Checkov+Trivy IaC, Image Build & Scan, Promotion Policy, Security Gate, kubeconform |
| 2. Artifacts, Deploy & Runtime | GHCR → Auto-CD (`workflow_run`) → Kustomize (dev/staging/prod) → k3d (OPA Gatekeeper, Falco, NetworkPolicy/RBAC/ESO, Pods) |
| 3. Monitoring, Cloud & Output | Falco → n8n → Slack; AWS lab (CloudTrail, Prowler, IAM); Output (Audit PDF, Blog + LinkedIn) |

### 2. Publikasi Showcase

- Draft blog `blogpost.md` (gitignored): **"Hari 60: Membangun Arsitektur DevSecOps End-to-End dalam 60 Hari"** — framing pencapaian "berhasil menguasai skill DevSecOps" dengan section **Keterampilan yang Dikuasai**: AppSec, Container & IaC, K8s & Runtime, Cloud Security & Red Team, Vuln Mgmt & Automation, CI/CD & Engineering Practices.
- Draft LinkedIn post dengan link repo.

### 3. AWS Cleanup (cost → $0)

Resource yang dihapus (semua di `ap-southeast-1`):

| Resource | Detail |
|----------|--------|
| Secrets Manager `securebank/jwt-secret` | `delete-secret --force-delete-without-recovery` — hemat $0.40/bulan |
| CloudTrail `securebank-trail` | `stop-logging` + `delete-trail` |
| S3 bucket `securebank-cloudtrail-logs-683915449775` | Versioning enabled → hapus **18.098 versions** + delete markers, lalu `rb --force` |

Yang **sengaja dipertahankan** (cost Rp0, hasil remediasi keamanan Day 51): S3 Block Public Access account-level, default Security Group revoke, IAM password policy.

### 4. Teardown Local

- `k3d cluster delete securebank` — hapus hanya klaster challenge (klaster project lain tidak disentuh).
- `docker-compose down` di `securebank-api/` dan `securebank-api/security/defectdojo/`.
- Stop webhook receiver Python (PID di port 5678).

---

## 📝 Catatan Teknis

```bash
# Render diagram headless (draw.io desktop)
brew install --cask drawio
/opt/homebrew/bin/drawio --export --format png --scale 2 \
  docs/architecture/day-60-architecture.drawio

# AWS cleanup — versioned S3 bucket
while true; do
  v=$(aws s3api list-object-versions --bucket <BUCKET>)
  [ count == 0 ] && break
  aws s3api delete-objects --bucket <BUCKET> --delete "$(gen_json)"
done
aws s3 rb s3://<BUCKET> --force
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| S3 bucket log CloudTrail punya versioning → `rb --force` gagal `BucketNotEmpty` | Hapus semua 18.098 object versions + delete markers via `list-object-versions` → `delete-objects` berulang, lalu `rb --force` |
| draw.io CLI belum tersedia | `brew install --cask drawio` menyediakan wrapper CLI `/opt/homebrew/bin/drawio` untuk export headless |
| Docker daemon sempat mati | Diaktifkan ulang manual sebelum teardown; verifikasi `k3d cluster list` sebelum delete |
| Repo ternyata sudah PUBLIC | Tidak perlu ubah visibility, cukup verifikasi `gh repo view` + pastikan tidak ada rahasia |

---

## 📤 Output Hari Ini

- [x] `docs/architecture/day-60-architecture.drawio` + `.png` — diagram E2E
- [x] `blogpost.md` (gitignored) — artikel showcase + LinkedIn post
- [x] AWS resource challenge dihapus (secret, trail, bucket) — cost kembali nol
- [x] k3d `securebank` + DefectDojo + webhook receiver dimatikan
- [x] `progress/tracker.md` — **Fase 4: 15/15, Total 60/60**, Hari Aktif 56
- [x] Retrospektif Fase 4 terisi + README/progress README ter-update

---

## 💡 Pelajaran Baru

1. **S3 bucket versioning menyembunyikan biaya & kegagalan delete** — CloudTrail log bucket meng-versioning otomatis; cleanup wajib `delete-objects` per version + delete marker.
2. **Framing hasil = menguasai skill, bukan sekadar selesai** — blog showcase lebih efektif saat dikelompokkan per kemampuan (SAST/SCA/DAST, IaC, K8s runtime, cloud security, red teaming) dengan bukti output konkret.
3. **Cleanup adalah bagian dari profesionalisme** — dokumentasi teardown (AWS + lokal) menutup siklus challenge dengan rapi dan tanpa biaya residu.
4. **Diagram end-to-end membantu bercerita** — satu gambar menjelaskan 60 hari kerja: dari `git push` sampai laporan audit dan alert runtime.

---

## 🔗 Referensi

- [Diagram arsitektur](../../docs/architecture/day-60-architecture.drawio)
- [Fase 4 — Implementation Notes](../../docs/fase-4-vuln-redteam.md)
- [Tracking utama](../../progress/tracker.md)
- [Retrospektif Fase 4](../../progress/retrospektif/fase-4-retrospektif.md)
- Blog showcase: `blogpost.md` (gitignored)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Menutup 60 hari dengan rasa puas |
| Pemahaman materi | 5 | Semua fase terhubung dalam satu arsitektur utuh |
| Progres sesuai target | 5 | 60/60 selesai, cost AWS nol, repo publik |

---

## ➡️ Rencana Besok

- [ ] Simulasi ujian CDP (opsional, bisa kapan pun) memakai `docs/cdp-exam-guide.md`
- [ ] Publikasikan blog + LinkedIn post (draft sudah siap di `blogpost.md`)

*[← Hari 59](hari-59.md) | [Tracker →](../tracker.md)*