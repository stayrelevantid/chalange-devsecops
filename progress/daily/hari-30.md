# Hari 30 — Dokumentasi Fase 2

**📅 Tanggal:** 2026-07-06  
**⏱️ Durasi Belajar:** ~1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Tulis retrospective Fase 2 di `docs/fase-2-infra-container.md`
- [x] Isi `progress/retrospektif/fase-2-retrospektif.md`
- [x] Isi `progress/retrospektif/fase-1-retrospektif.md` (backlog)
- [x] Update progress files (tracker, README)

---

## ✅ Yang Berhasil Dikerjakan

- Tulis Fase 2 retrospective lengkap di `docs/fase-2-infra-container.md` (9 sections: arsitektur pipeline, tools, quality gates, security improvements, metrics, lessons learned, recommendations, file structure, pipeline history)
- Isi `fase-1-retrospektif.md` (15/15, Start/Stop/Continue, Top 5 lessons, challenges, skor diri)
- Isi `fase-2-retrospektif.md` (15/15, Start/Stop/Continue, Top 5 lessons, challenges, skor diri)
- Update tracker: Fase 2 15/15 ✅, Total 30/60, Fase Selesai 2/4

---

## 📝 Catatan Teknis

### Fase 2 Retrospective Structure (di docs/)
Mengikuti format Fase 1 retrospective:
1. Arsitektur Pipeline (mermaid diagram 8-job)
2. Tools yang Digunakan (11 tools)
3. Quality Gates (container hardening, signing, IaC, DAST, compliance)
4. Security Improvements (14 fixes Day 16-29)
5. Metrics (image size, CVE count, checkov results, ZAP results, tests)
6. Lessons Learned (9 items)
7. Recommendations for Fase 3 (8 items)
8. Final File Structure
9. Pipeline Execution History

### Key Metrics Documented
- Image size: 350MB → 7.97MB (44x reduction)
- Checkov: 102/0 passed/failed
- Trivy IaC: 0 findings (MEDIUM+)
- ZAP: 0 FAIL, 1 WARN (false positive), 66 PASS
- InSpec: 3/3 controls PASS
- Pipeline: 8 jobs in 1 file, all green
- Unit tests: 25 all PASS

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Fase 1 retrospective template juga kosong | Isi sekalian dengan data dari `docs/fase-1-appsec.md` retrospective |
| Banyak data untuk dirangkum (15 hari, 14 security fixes) | Kelompokkan per kategori: container, IaC, DAST, pipeline |

---

## 📤 Output Hari Ini

- [x] `docs/fase-2-infra-container.md` — retrospective lengkap (~350 lines)
- [x] `progress/retrospektif/fase-1-retrospektif.md` — filled (15/15)
- [x] `progress/retrospektif/fase-2-retrospektif.md` — filled (15/15)
- [x] `progress/daily/hari-30.md` — diisi
- [x] tracker, README — updated (30/60, Fase 2: 15/15)

---

## 💡 Pelajaran Baru

- **Retrospective = forcing function untuk refleksi.** Menulis retrospective bikin aku kembali review semua 15 hari, lihat pola, dan identifikasi apa yang akan dilakukan berbeda di fase berikutnya.

- **Start/Stop/Continue framework itu powerful.** Lebih actionable daripada list pelajaran. "Stop assuming 0 findings = safe" lebih konkret daripada "scanners have limitations."

- **Dokumentasi retrospective adalah knowledge transfer.** Kalau ini proyek tim, retrospective ini akan jadi onboarding document untuk developer baru — mereka tahu kenapa setiap decision dibuat, apa yang work, apa yang tidak.

---

## 🔗 Referensi

- [Fase 1 retrospective](../../docs/fase-1-appsec.md) (line 815-1143)
- [Agile retrospective framework](https://www.atlassian.com/agile/scrum/retrospectives)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Retrospective = closure untuk Fase 2 |
| Pemahaman materi | 5 | Semua 15 hari terdokumentasi |
| Progres sesuai target | 5 | 30/60 separuh journey, Fase 2 complete |

---

## ➡️ Rencana Besok

- [ ] Hari 31: K8s Cluster + Deploy — k3d cluster, `kubectl apply` deployment & service
- [ ] Mulai Fase 3: Kubernetes & Runtime Security

---

*[← Hari 29](hari-29.md) | [Hari 31 →](hari-31.md)*