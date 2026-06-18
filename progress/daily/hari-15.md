# Hari 15 — Dokumentasi Fase 1

**📅 Tanggal:** 2026-06-18  
**⏱️ Durasi Belajar:** 2 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Rangkum arsitektur pipeline CI/CD
- [x] Dokumentasi semua tools, config, dan versi
- [x] Definisikan quality gates (setiap layer)
- [x] Tabel security improvements (14 hari, 17 fixes)
- [x] Threat model summary (STRIDE + DREAD highlights)
- [x] Metrics (vulns found, fixed, tests, deps)
- [x] Lessons learned (7 pelajaran)
- [x] Recommendations for Fase 2 (8 item)

---

## ✅ Yang Berhasil Dikerjakan

- Menambahkan **Dokumentasi Retrospektif Fase 1** di `docs/fase-1-appsec.md` (9 bagian):
  1. Arsitektur Pipeline (Mermaid diagram + tabel 4 parallel jobs)
  2. Tools yang Digunakan (5 tools + 2 Go dependencies)
  3. Quality Gates (4 gates: Secret Scan, SCA, SAST, Build & Test)
  4. Security Improvements Applied (17 fixes dari Day 03-14)
  5. Threat Model Summary (3 Critical, 4 High, dari Day 13)
  6. Metrics (4 CVE found/fixed, 2 SAST findings, 24 tests, 0 current CVE)
  7. Lessons Learned (7 pelajaran)
  8. Recommendations for Fase 2 (8 rekomendasi)
  9. Final File Structure (pohon direktori lengkap)
  10. Pipeline Execution History (timeline Day 02-14)

---

## 📝 Catatan Teknis

```
# Dokumentasi ditambahkan ke file yang sudah ada (bukan file baru)
docs/fase-1-appsec.md  —  appended retrospective setelah line 811

# Struktur retrospektif:
# 1. Arsitektur Pipeline (Mermaid diagram)
# 2. Tools yang Digunakan (tabel)
# 3. Quality Gates (4 gates dengan detail)
# 4. Security Improvements Applied (17 fixes)
# 5. Threat Model Summary (STRIDE + DREAD)
# 6. Metrics
# 7. Lessons Learned (7 pelajaran)
# 8. Recommendations for Fase 2
# 9. Final File Structure
# 10. Pipeline Execution History
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| File tutorial sudah 811 baris, menambah di akhir | Append retrospektif setelah baris terakhir, dipisahkan oleh `---` |
| Banyak data yang perlu dirangkum (14 hari) | Gunakan tabel untuk ringkasan, link ke daily logs untuk detail |

---

## 📤 Output Hari Ini

- [x] Dokumentasi retrospektif Fase 1 di `docs/fase-1-appsec.md`
- [x] 10 bagian dokumentasi lengkap
- [x] Commit: `docs: add Fase 1 retrospective documentation (Day 15)`

---

## 💡 Pelajaran Baru

- **Dokumentasi itu investasi, bukan biaya.** Tanpa retrospektif, knowledge dari 14 hari hilang. Dengan dokumentasi, tim baru bisa onboarding lebih cepat.
- **Metrics itu penting untuk storytelling.** "4 CVE found, 4 fixed, 0 remaining" lebih impactful daripada "we fixed some CVEs".
- **Threat model bukan satu kali.** Ini living document yang harus di-update setiap fase baru. Fase 2 (container security) akan menambah attack surface baru.
- **Lessons learned harus actionable.** Setiap pelajaran harus terhubung ke rekomendasi spesifik. "Gitleaks skip .gitignore" → "tambah --no-gitignore flag".

---

## 🔗 Referensi

- [OWASP Secure SDLC](https://owasp.org/www-project-secure-sdlc/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Gitleaks Configuration](https://gitleaks.io/)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [Semgrep Documentation](https://semgrep.dev/docs/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Dokumentasi kurang exciting tapi sangat penting |
| Pemahaman materi | 5 | Full picture setelah 14 hari kerja keras |
| Progres sesuai target | 5 | Fase 1 complete! 15/15 hari |

---

## ➡️ Rencana Besok

- [ ] Hari 16: Dockerfile Multi-stage — Build di `golang:alpine`, run di `distroless`

---

*[← Hari 14](hari-14.md) | [Hari 16 →](hari-16.md)*