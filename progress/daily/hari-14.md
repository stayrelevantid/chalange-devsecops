# Hari 14 — Intentional Vulnerability Test

**📅 Tanggal:** 2026-06-18  
**⏱️ Durasi Belajar:** 2 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Commit fake secret ke branch terpisah
- [x] Verifikasi bahwa Gitleaks detect secret tersebut
- [x] Identifikasi gap dalam pipeline secret scanning
- [x] Fix pipeline dan re-verify
- [x] Hapus branch test dan tutup PR

---

## ✅ Yang Berhasil Dikerjakan

- Membuat branch `test/intentional-leak` dengan file `configs/.env.production` berisi:
  - `AWS_ACCESS_KEY_ID=AKIA3E4Y5Z6X7W8Q9R0T`
  - `AWS_SECRET_ACCESS_KEY=7hGjKmNpQrStUvWxYz2B4C6D8E0F1G2H3J4K5L6M`
  - `DATABASE_URL=postgresql://prodadmin:Xk9mP2vR7nLqW4jF@db.securebank.internal:5432/securebank_prod`
- Gitleaks lokal **berhasil deteksi** 2 findings saat scan file langsung
- Gitleaks di CI **GAGAL deteksi** — ini gap kritis yang ditemukan!
- Root cause: Gitleaks menghormati `.gitignore`, sehingga file yang di-gitignore (meskip di-force-add) tidak di-scan
- Fix: tambahkan `--no-gitignore` flag ke Gitleaks CI command
- PR #1 ditutup dan branch dihapus

---

## 📝 Catatan Teknis

```bash
# Gitleaks detect secret secara lokal (scan langsung file)
$ gitleaks detect --source securebank-api/configs/.env.production -v --no-git
2 findings: aws-access-token, generic-api-key ✅

# Gitleaks GAGAL detect saat scan directory
$ gitleaks detect --source . -v --config .gitleaks.toml
no leaks found ❌

# Penyebab: .gitignore memfilter file .env.production
# File ada di git (force-add), tapi Gitleaks skip karena di-gitignore

# Fix: tambah --no-git-ignore flag
$ gitleaks detect --source . -v --config .gitleaks.toml --no-gitignore
```

**PR:** https://github.com/stayrelevantid/chalange-devsecops/pull/1 (closed)

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `.gitignore` mencegah accidental commit | Ini pertahanan pertama yang bekerja — `git add` biasa ditolak, butuh `git add -f` |
| Gitleaks tidak scan file yang di-`.gitignore` | Tambahkan `--no-gitignore` flag ke CI command |
| `aws_access_key_id` dengan pattern `AKIA...` dianggap example | Gunakan pattern yang lebih realistis (random digits bukan sequential) |
| CI gak trigger di branch `test/intentional-leak` | Buat PR ke main supaya pipeline jalan di PR context |

---

## 📤 Output Hari Ini

- [x] Fake secret committed dan di-push ke branch terpisah
- [x] Gitleaks lokal berhasil detect 2 findings
- [x] Gitleaks CI gagal detect (gap ditemukan!)
- [x] Fix: `--no-gitignore` flag ditambahkan ke CI
- [x] PR #1 ditutup dan branch dihapus
- [x] Pipeline fix di-push ke main

---

## 💡 Pelajaran Baru

- **`.gitignore` adalah pertahanan pertama, bukan satu-satunya.** File `.env` gak bisa di-commit secara normal (blocked by `.gitignore`), tapi bisa di-force-add (`git add -f`). Pipeline secret scanning harus tetap menangkap kasus ini.

- **Gitleaks `--no-gitignore` itu WAJIB di CI.** Tanpa flag ini, Gitleaks skip semua file yang ada di `.gitignore`. Ini berarti jika seseorang `git add -f .env.production`, Gitleaks gak detect. Flag ini bikin Gitleaks scan SEMUA file yang tracked di git, regardless of `.gitignore`.

- **Intentional vulnerability test itu penting.** Tanpa test ini, kita gak bakal tahu bahwa Gitleaks CI gak detect file yang di-gitignore. Ini gap yang gak terlihat saat setup awal.

- **PR context vs branch push context beda di CI.** GitHub Actions bisa trigger di branch push (hanya `main`, `develop`) dan PR events (ke `main`). Untuk test security gate, PR context lebih realistis karena men-scan diff branch terhadap base.

- **AWS example keys (`AKIAIOSFODNN7EXAMPLE`) di-flag oleh Gitleaks sebagai known example keys.** Kunci AWS yang valid harus lebih random, bukan pattern sequential.

---

## 🔗 Referensi

- [Gitleaks Documentation](https://gitleaks.io/)
- [Gitleaks --no-gitignore flag](https://github.com/gitleaks/gitleaks#configuration)
- [OWASP - Secret Management](https://owasp.org/www-community/vulnerabilities/Use_of_hard-coded_cryptographic_key)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Menemukan gap nyata di pipeline itu sangat memuaskan |
| Pemahaman materi | 5 | Paham sekarang kenapa .gitignore bukan security boundary |
| Progres sesuai target | 5 | Intentional test berhasil, gap ditemukan, fix diterapkan |

---

## ➡️ Rencana Besok

- [ ] Hari 15: Dokumentasi Fase 1 — rangkum cara kerja SAST, SCA, Secret Scan

---

*[← Hari 13](hari-13.md) | [Hari 15 →](hari-15.md)*