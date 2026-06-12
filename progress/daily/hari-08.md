# Hari 08 — SAST Setup (Semgrep)

**📅 Tanggal:** 2026-06-12  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Semgrep terinstall dan berjalan
- [x] Menemukan penggunaan MD5 yang tidak aman
- [x] Custom rule `.semgrep.yml` dibuat
- [x] Output dipahami (severity, file, line number)

---

## ✅ Yang Berhasil Dikerjakan

- Install Semgrep v1.165.0 via Homebrew
- Buat `pkg/crypto/hash.go` — fungsi `HashPassword` yang menggunakan MD5 (sengaja insecure untuk latihan SAST)
- Jalankan `semgrep --config "p/golang" .` — menemukan 2 findings:
  - `use-of-md5` di `pkg/crypto/hash.go:11` — MD5 adalah weak hash (WARNING)
  - `use-tls` di `cmd/api/main.go:80` — HTTP server tanpa TLS (WARNING)
- Buat `.semgrep.yml` di root repo dengan 2 custom rules:
  - `no-md5-usage` (ERROR) — deteksi penggunaan MD5
  - `no-http-listen-without-tls` (WARNING) — deteksi HTTP server tanpa TLS
- Custom rules menemukan 2 findings (sama dengan built-in rules)
- Simpan report JSON di `security/semgrep-report.json`
- Update `.gitleaks.toml` — exclude `.semgrep.yml` dari secret scanning
- Semua 8 unit test masih hijau

---

## 📝 Catatan Teknis

```bash
# Install Semgrep
brew install semgrep

# Scan dengan built-in Go rules
semgrep --config "p/golang" .

# Scan dengan custom rules
semgrep --config .semgrep.yml .

# Simpan report JSON
semgrep --config "p/golang" . --json --output security/semgrep-report.json
```

### Semgrep findings:

| Rule | File | Line | Severity | Pesan |
|------|------|------|----------|-------|
| `use-of-md5` | `pkg/crypto/hash.go` | 11 | WARNING | MD5 is cryptographically broken |
| `use-tls` | `cmd/api/main.go` | 80 | WARNING | HTTP server tanpa TLS |

### Custom rule `.semgrep.yml`:

```yaml
rules:
  - id: no-md5-usage
    patterns:
      - pattern: md5.Sum(...)
    message: "MD5 is cryptographically broken. Use SHA-256 or bcrypt."
    languages: [go]
    severity: ERROR

  - id: no-http-listen-without-tls
    patterns:
      - pattern: http.ListenAndServe(...)
    message: "HTTP server without TLS. Use http.ListenAndServeTLS()."
    languages: [go]
    severity: WARNING
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Custom rule `pattern: crypto/md5` tidak match | Ganti pattern ke `md5.Sum(...)` karena Semgrep match berdasarkan pemanggilan fungsi, bukan import path |
| `.semgrep.yml` harus di-root repo | File scan dari subdirektori `securebank-api/` menggunakan `--config ../.semgrep.yml` |

---

## 📤 Output Hari Ini

- [x] Semgrep v1.165.0 terinstall
- [x] `pkg/crypto/hash.go` — kode insecure (MD5) untuk latihan SAST
- [x] `.semgrep.yml` — 2 custom rules di root repo
- [x] `security/semgrep-report.json` — report hasil scan (2 findings)
- [x] 8/8 unit test masih hijau
- [x] `.gitleaks.toml` diupdate — exclude `.semgrep.yml`

---

## 💡 Pelajaran Baru

- **SAST itu beda sama SCA** — SCA scan dependensi pihak ketiga, SAST scan kode yang kita tulis sendiri. Keduanya saling melengkapi
- **Semgrep bisa pakai built-in rules dan custom rules** — `p/golang` untuk aturan bawaan Go, `.semgrep.yml` untuk aturan khusus project
- **Custom rule pattern harus match pemanggilan fungsi** — `crypto/md5` (import path) tidak match, tapi `md5.Sum(...)` (pemanggilan) match. Ini beda sama SCA yang scan go.mod
- **Severity di custom rule bisa diatur** — kita set MD5 sebagai ERROR (lebih ketat) dan HTTP tanpa TLS sebagai WARNING
- **Semgrep bisa detect insecure pattern yang gak terpikir** — HTTP tanpa TLS di `main.go` gak pernah kita pikirkan sebagai vulnerability, tapi Semgrep langsung flag

---

## 🔗 Referensi

- [Semgrep Documentation](https://semgrep.dev/docs/)
- [Semgrep Go Rules](https://semgrep.dev/p/golang)
- [CWE-328: Use of Weak Hash](https://cwe.mitre.org/data/definitions/328.html)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | SAST itu topik yang menarik |
| Pemahaman materi | 4 | Custom rule pattern agak tricky |
| Progres sesuai target | 4 | Sesuai rencana |

---

## ➡️ Rencana Besok

- [ ] Integrasikan Semgrep ke CI pipeline (Day 09 - SAST Pipeline)
- [ ] Tambah job `sast-scan` di `.github/workflows/ci.yml`
- [ ] Quality gate: pipeline gagal jika severity ERROR

---

*[← Hari 07](hari-07.md) | [Hari 09 →](hari-09.md)*