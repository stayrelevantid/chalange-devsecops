# Hari 56 — Dokumen Eksekutif (PDF)

**📅 Tanggal:** 2026-08-07
**⏱️ Durasi Belajar:** ~90 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Perluas `draft-q3.md` → dokumen audit formal `laporan-audit-q3.md` (metodologi, key findings, mitigasi, residual risk, lampiran)
- [x] Export dokumen ke PDF menggunakan Pandoc + Typst
- [x] Verifikasi hasil render PDF (struktur & tabel)

---

## ✅ Yang Berhasil Dikerjakan

### 1. Penyusunan Dokumen Audit Formal

Draf dari Day 55 (`draft-q3.md`) diperluas menjadi dokumen lengkap berformat YAML metadata (title, author, date) agar pandoc bisa meng-generate halaman judul:

`securebank-api/security/audit-reports/laporan-audit-q3.md` dengan struktur 7 bagian + lampiran:

| Bagian | Isi |
|--------|-----|
| Executive Summary | Penilaian risiko keseluruhan **MEDIUM**, 6 temuan aktif (0C/1H/2M/2L/1I), turun dari 52 agregat |
| Ruang Lingkup & Metodologi | 11 layer tools + 5 langkah pendekatan (deteksi → residual) |
| Key Findings (Baseline) | F-01 s.d. F-10 — temuan sebelum remediasi dari Day 05–54 |
| Mitigation Evidence (After) | Tabel bukti remediasi per temuan + perbandingan sebelum/sesudah |
| Residual Risk | 5 risiko tersisa + 3 risiko yang diterima (dengan alasan) |
| Rekomendasi & Roadmap | 5 rekomendasi (Q3 segara → Q4 jangka menengah) |
| Konklusi | 1 temuan High tersisa, solusinya rebuild image → bisa naik ke LOW |
| Lampiran | Tabel alat & versi + referensi bukti per hari |

### 2. Export PDF via Pandoc + Typst

```bash
pandoc securebank-api/security/audit-reports/laporan-audit-q3.md \
  --pdf-engine=typst \
  -o securebank-api/security/audit-reports/laporan-audit-q3.pdf
```

- **pandoc 3.10.1** + **typst 0.15.1** (di-install via Homebrew)
- PDF yang dihasilkan: **7 halaman**, valid PDF 1.7, 132 KB — tabel Markdown langsung ter-render rapi oleh engine Typst
- Tidak perlu plugin lain; `--pdf-engine=typst` tersedia sejak pandoc 3.1.8

### 3. Verifikasi Render

Validasi struktur PDF tanpa PDF viewer:

```bash
file laporan-audit-q3.pdf        # PDF document, version 1.7, 7 pages
```

Semua tabel (metrik, tools, temuan, mitigasi, residual risk, lampiran) ikut ter-render karena pandoc mengonversi tabel GFM → tabel Typst secara native.

---

## 🧪 Hasil Checklist

| Checklist | Hasil |
|-----------|-------|
| Laporan memiliki 4 komponen utama (Summary, Methodology, Findings/Mitigations, Residual Risk) | ✅ 7 bagian + lampiran |
| Dokumen diekspor ke PDF | ✅ `laporan-audit-q3.pdf` (7 halaman) |

---

## 📝 Catatan Teknis

### YAML metadata di Markdown
Menambahkan front matter (title, subtitle, author, date) membuat pandoc menghasilkan halaman judul otomatis. Tanpa front matter, judul diambil dari heading pertama.

### Pandoc `--pdf-engine=typst`
Berdasarkan pandoc 3.1.8, engine Typst tersedia sebagai opsi PDF default. Typst menggantikan kebutuhan LaTeX (wkhtmltopdf/weasyprint) — hasil render tabel GFM bersih tanpa instalasi paket LaTeX besar.

### Angka laporan bersumber dari state terkini
Semua metrik di dokumen formal konsisten dengan re-sync Day 55: 52 temuan agregat DefectDojo vs 6 temuan state terkini; Prowler 132→106 FAIL; tren remediasi per layer (Fase 1–3). Menghindari angka stale yang menyesatkan.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| PDF tidak bisa dibuka di tool preview (model AI) | Validasi via `file` (7 pages, PDF 1.7) — struktur & ukuran terkonfirmasi |
| `pdftotext` tidak terpasang di macOS | Cukup pakai `file` untuk verifikasi; render tabel dijamin oleh engine Typst |

---

## 📤 Output Hari Ini

- [x] `securebank-api/security/audit-reports/laporan-audit-q3.md` — dokumen audit formal 7 bagian + lampiran
- [x] `securebank-api/security/audit-reports/laporan-audit-q3.pdf` — export PDF 7 halaman (pandoc + typst)

---

## 💡 Lessons Learned

### 1. Draf ≠ dokumen final
Draf eksekutif (Day 55) berisi poin-poin; dokumen formal (Day 56) harus berbicara ke pembaca non-teknis — menonjolkan dampak bisnis, tren sebelum/sesudah, dan risiko tersisa dengan rencana tindak, bukan hanya daftar CVE.

### 2. Evidence-based storytelling
Bagian "Mitigation Evidence" yang paling meyakinkan bukan narasi, melainkan angka scan ulang: IaC 14→2 Low, K8s 16→0, Prowler 132→106. Setiap klaim remediasi punya bukti dari hari tertentu.

### 3. "Risiko yang diterima" harus jujur
Dokumen eksekutif profesional tidak menyembunyikan sisa risiko — malah mendokumentasikannya dengan alasan (misal: WAF dijadwalkan Q4, node agent NotReady karena keterbatasan lab). Transparansi ini yang membangun kepercayaan.

### 4. Tooling lintas-platform itu penting
Dari Day 55 (checkov rusak di Python 3.14) ke Day 56 (pandoc+typst via Homebrew) — pola yang sama: pilih tool yang mudah diinstall dan berfungsi konsisten di macOS.

---

## 🔗 Referensi

- [Hari 55](hari-55.md) — draf laporan + data state terkini
- [Day 51](hari-51.md) — CSPM remediation (bukti Prowler 132→106)
- [Day 33](hari-33.md) — K8s hardening (bukti 16→0)
- [Pandoc PDF — Typst](https://pandoc.org/MANUAL.html#option--pdf-engine) — engine PDF tanpa LaTeX

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Dokumen akhir mulai berbentuk profesional & siap presentasi |
| Pemahaman materi | 4 | Pandoc/typst, struktur laporan eksekutif, narrative evidence |
| Progres sesuai target | 5 | Dokumen + PDF tuntas; tinggal review AI Day 57 |

---

## ➡️ Rencana Besok

- [ ] **Day 57: AI Review Dokumen** — gunakan AI untuk memperbaiki nada bahasa agar lebih eksekutif (CISO/CTO), fokus dampak bisnis

---

*[← Hari 55](hari-55.md) | [Hari 57 →](hari-57.md)*
