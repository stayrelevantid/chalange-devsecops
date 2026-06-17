# Hari 13 — Threat Modeling (STRIDE)

**📅 Tanggal:** 2026-06-15  
**⏱️ Durasi Belajar:** 2 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat diagram arsitektur SecureBank API (Mermaid)
- [x] Analisis STRIDE — 6 kategori ancaman
- [x] DREAD risk assessment — scoring per ancaman
- [x] Mitigation roadmap — prioritas top 3
- [x] Attack tree untuk transfer endpoint

---

## ✅ Yang Berhasil Dikerjakan

- Membuat `security/threat-model/architecture.md` dengan:
  - Diagram komponen (Mermaid) — Client → LB → API → Middleware → Data Store
  - Trust boundary diagram — Internet Zone, DMZ, Application Zone
  - 2 sequence diagram — Authentication Flow dan Transfer Flow
  - 18 threat entries di STRIDE table
  - 12 threat entries di DREAD scoring
  - Attack tree untuk transfer endpoint
  - Mitigation roadmap (top 3 + fase target)

---

## 📝 Catatan Teknis

**STRIDE Analysis Summary:**

| Kategori | Jumlah Temuan |
|----------|---------------|
| Spoofing | 3 (1 ✅, 1 ⚠️, 1 ⚠️) |
| Tampering | 4 (3 ✅, 1 ❌) |
| Repudiation | 2 (1 ⚠️, 1 ❌) |
| Info Disclosure | 4 (1 ✅, 2 ⚠️, 1 ❌) |
| Denial of Service | 3 (1 ✅, 2 ❌) |
| Elevation of Privilege | 2 (2 ❌) |

**DREAD Priority Summary:**

| Priority | Count |
|----------|-------|
| 🔴 Critical | 3 |
| 🟠 High | 4 |
| 🟡 Medium | 4 |
| 🟢 Low | 0 |
| ✅ Mitigated | 6 |

**Top 3 Mitigasi:**
1. Rate limiter middleware (DREAD Score: 9.6)
2. User-scoped authorization (DREAD Score: 8.4)
3. Persistent audit logging (DREAD Score: 6.0)

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Mermaid diagram di Markdown gak render di semua viewer | Diagram ditulis dalam format Mermaid code block — render di GitHub, VS Code, dan Mermaid Live Editor |
| DREAD scoring subjektif | Gunakan skala 0-10 per kategori, rata-ratakan untuk score akhir — lebih objektif daripada "gut feeling" |
| Beberapa threat overlap (no auth ↔ no user scoping) | Pisahkan jadi threat terpisah: spoofing (auth) vs tampering (authorization) |

---

## 📤 Output Hari Ini

- [x] Diagram arsitektur Mermaid (komponen + trust boundary)
- [x] Sequence diagram (auth flow + transfer flow)
- [x] STRIDE analysis table — 18 threats with mitigation status
- [x] DREAD risk assessment — 12 scored threats
- [x] Attack tree untuk transfer endpoint
- [x] Mitigation roadmap (top 3 + fase target)
- [x] Commit: `docs: add STRIDE+DREAD threat model for SecureBank API (Day 13)`

---

## 💡 Pelajaran Baru

- **STRIDE itu systematic, bukan creative.** Setiap komponen dan data flow diuji terhadap 6 kategori ancaman. Gak perlu mikir "apa serangan yang mungkin?" — cukup ikuti framework.
- **DREAD scoring bikin prioritas objektif.** Daripada bilang "ini bahaya," kita punya angka. No rate limiting (9.6) jelas lebih kritis daripada JWT_SECRET di env var (4.0).
- **Attack tree membantu visualisasi attack path.** Dari attack tree transfer endpoint, jelas bahwa path termudah adalah: valid JWT → query balance siapa saja → transfer dari akun siapa saja. Ini langsung tunjukin kenapa authorization penting.
- **Mitigation status itu penting.** Threat model bukan cuma list serangan — tapi juga document status mitigasi. 6 dari 18 threat udah ✅ Implemented dari Day 11.
- **Threat model itu living document.** Ini versi 1.0. Setiap kali ada fitur baru atau lapisan keamanan baru di fase berikutnya, threat model harus di-update.

---

## 🔗 Referensi

- [OWASP Threat Modeling](https://owasp.org/www-community/Threat_Modeling)
- [STRIDE Methodology — Microsoft](https://learn.microsoft.com/en-us/azure/security/develop/threat-modeling-tool-threats)
- [DREAD Risk Assessment](https://owasp.org/www-community/attacks/DREAD)
- [Mermaid Diagram Syntax](https://mermaid.js.org/intro/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Threat modeling sangat relevan buat keamanan nyata |
| Pemahaman materi | 5 | STRIDE + DREAD bikin analisa jadi systematic |
| Progres sesuai target | 5 | 18 threats, 12 DREAD scores, attack tree — lengkap |

---

## ➡️ Rencana Besok

- [ ] Hari 14: Intentional Vuln Test — commit fake secret → Gitleaks gagalkan build → revert

---

*[← Hari 12](hari-12.md) | [Hari 14 →](hari-14.md)*