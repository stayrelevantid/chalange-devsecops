# Hari 03 — Secret Scanning (Gitleaks)

**📅 Tanggal:** 2026-06-07  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Gitleaks terinstall dan berjalan lokal
- [x] Berhasil mendeteksi secret di repository (finding dari docs tutorial)
- [x] `.gitleaks.toml` dikonfigurasi untuk exclude false positive
- [x] Branch test dihapus (secret tidak masuk main)

---

## ✅ Yang Berhasil Dikerjakan

- Install Gitleaks v8.21.2 (download binary dari GitHub releases)
- Membuat branch `test/secret-leak` dengan file `configs/database.go` berisi fake secret (DBPassword, AWSKey, AWSSecret)
- Menjalankan `gitleaks detect --source . -v` — menemukan 1 finding dari `docs/fase-1-appsec.md` (contoh Stripe key di dokumentasi tutorial)
- Gitleaks tidak melaporkan fake AWS key karena `AKIAIOSFODNN7EXAMPLE` adalah AWS example key yang sudah dikenali sebagai false positive oleh Gitleaks
- Membuat `.gitleaks.toml` di root repo yang meng-exclude `go.sum`, `vendor/`, dan `docs/`
- Scan ulang dengan config: **0 leaks found** ✅
- Menyimpan report kosong ke `securebank-api/security/gitleaks-report.json`
- Branch test dihapus, fake secret tidak masuk ke main

---

## 📝 Catatan Teknis

```bash
# Install Gitleaks (macOS ARM64)
curl -L https://github.com/gitleaks/gitleaks/releases/download/v8.21.2/gitleaks_8.21.2_darwin_arm64.tar.gz -o /tmp/gitleaks.tar.gz
tar -xzf /tmp/gitleaks.tar.gz -C /tmp/
cp /tmp/gitleaks /opt/homebrew/bin/gitleaks
gitleaks version  # v8.21.2

# Scan tanpa config — menemukan 1 finding di docs
gitleaks detect --source . -v --no-banner
# Finding: stripe-access-token di docs/fase-1-appsec.md:745

# Scan dengan config — 0 leaks found
gitleaks detect --source . -v --no-banner --config .gitleaks.toml
# no leaks found ✅

# Simpan report
gitleaks detect --source . --config .gitleaks.toml \
  --report-format json \
  --report-path securebank-api/security/gitleaks-report.json
```

**Konfigurasi `.gitleaks.toml`:**
```toml
title = "SecureBank Gitleaks Config"

[allowlist]
  paths = [
    '''go\.sum''',
    '''vendor/''',
    '''docs/''',
  ]
```

**Penting:** `docs/` di-exclude karena berisi contoh secret di dokumentasi tutorial (bukan secret asli).

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `brew install gitleaks` timeout karena koneksi lambat | Download binary manual dari GitHub releases |
| Fake AWS key (`AKIAIOSFODNN7EXAMPLE`) tidak terdeteksi | Gitleaks mengenali ini sebagai example key (bukan secret asli) — ini behavior yang benar! |
| Gitleaks menemukan 1 finding di `docs/fase-1-appsec.md` | Itu contoh Stripe key di dokumentasi, di-exclude lewat `.gitleaks.toml` |
| `.gitleaks.toml` harus di root repo (bukan di subdirektori) | Ditaruh di root agar scan dari root juga membaca config-nya |

---

## 📤 Output Hari Ini

- [x] Gitleaks v8.21.2 terinstall di `/opt/homebrew/bin/gitleaks`
- [x] `.gitleaks.toml` di root repo — config untuk exclude false positive
- [x] `securebank-api/security/gitleaks-report.json` — report bersih (0 finding)
- [x] Branch test/secret-leak sudah dihapus, fake secret tidak masuk main
- [x] Commit: `9bf4464 security: add gitleaks config to exclude docs and false positives`

---

## 💡 Pelajaran Baru

**1. Gitleaks itu cukup pintar.** AWS example key seperti `AKIAIOSFODNN7EXAMPLE` tidak dilaporkan karena Gitleaks tahu itu bukan kunci nyata. Ini mengurangi false positive secara signifikan.

**2. Dokumentasi bisa jadi sumber false positive.** File `docs/fase-1-appsec.md` berisi contoh secret buat tutorial — Gitleaks akan menandai ini sebagai kebocoran nyata kalau tidak di-exclude. Makanya `.gitleaks.toml` penting banget.

**3. Scan git history vs scan file biasa beda.** `gitleaks detect` secara default scan git history (semua commit). Kalau mau scan file di disk tanpa git, tambahkan flag `--no-git`.

**4. Branch terpisah itu strategi yang aman.** Dengan bikin `test/secret-leak` branch, kita bisa menguji Gitleaks tanpa risiko secret masuk ke main. Setelah selesai, branch langsung dihapus.

---

## 🔗 Referensi

- [Gitleaks GitHub](https://github.com/gitleaks/gitleaks)
- [Gitleaks Configuration](https://github.com/gitleaks/gitleaks#configuration)
- [AWS Example Keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Seru lihat tool security langsung kerja |
| Pemahaman materi | 4 | Gitleaks cukup mudah dipahami |
| Progres sesuai target | 5 | Sesuai silabus hari ke-3 |

---

## ➡️ Rencana Besok

- [ ] Integrasi Gitleaks ke GitHub Actions pipeline
- [ ] Job `secret-scan` di CI gagalkan build jika ada secret
- [ ] Pindahkan kredensial ke environment variables
- [ ] Pipeline hijau setelah remediasi

---

*[← Hari 02](hari-02.md) | [Hari 04 →](hari-04.md)*