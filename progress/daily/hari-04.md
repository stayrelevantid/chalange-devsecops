# Hari 04 — Secret Scan di Pipeline + Remediasi

**📅 Tanggal:** 2026-06-08  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Job `secret-scan` berjalan di pipeline
- [x] Tidak ada hardcoded credential di source code
- [x] `configs/config.go` membaca dari environment variables
- [x] Pipeline hijau (build + test + secret scan)

---

## ✅ Yang Berhasil Dikerjakan

- Menambah job `Secret Scan (Gitleaks)` ke `.github/workflows/ci.yml`
- Job berjalan paralel dengan `Build & Test (Go 1.26)` — konsisten labeling
- Membuat `securebank-api/configs/config.go` — struct `Config` yang baca dari env vars
- Refactor `cmd/api/main.go` — port dibaca dari `configs.Load()` (env var `PORT`, default `8080`)
- Menambah unit test untuk `configs` package (2 test: default values & env var override)
- Semua 8 unit test PASS (6 handler + 2 config)
- Gitleaks lokal: 0 leaks found
- `GITHUB_TOKEN` otomatis tersedia di GitHub Actions — ga perlu setup manual

---

## 📝 Catatan Teknis

```yaml
# .github/workflows/ci.yml — job baru
secret-scan:
  name: Secret Scan (Gitleaks)
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4
      with:
        fetch-depth: 0  # full history untuk scan commit lama

    - name: Run Gitleaks
      uses: gitleaks/gitleaks-action@v2
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

```go
// securebank-api/configs/config.go
type Config struct {
    Port       string
    DBHost     string
    DBPassword string
}

func Load() *Config {
    return &Config{
        Port:       getEnv("PORT", "8080"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPassword: getEnv("DB_PASSWORD", ""),
    }
}
```

**Poin penting:**
- `fetch-depth: 0` di checkout — biar Gitleaks bisa scan seluruh git history, bukan cuma commit terakhir
- `GITHUB_TOKEN` otomatis dibuat oleh GitHub Actions setiap workflow run — ga perlu bikin token manual
- Job `secret-scan` berjalan paralel dengan `build-and-test` — ga saling tunggu, lebih cepat

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `GITHUB_TOKEN` perlu setup manual? | Tidak — otomatis tersedia di setiap GitHub Actions workflow run |
| Port hardcode `:8080` di `main.go` | Dipindah ke env var `PORT` via `configs.Load()`, default 8080 |
| Label job harus konsisten | Format: `Job Name (Tool/Version)` — `Build & Test (Go 1.26)`, `Secret Scan (Gitleaks)` |

---

## 📤 Output Hari Ini

- [x] `.github/workflows/ci.yml` — pipeline dengan 2 job: Build & Test, Secret Scan
- [x] `securebank-api/configs/config.go` — config struct baca dari env vars
- [x] `securebank-api/configs/config_test.go` — 2 unit test untuk config
- [x] `securebank-api/cmd/api/main.go` — refactored, port dari env var
- [x] Gitleaks lokal: 0 leaks found
- [x] Commit: `e337a20 security: integrate gitleaks in CI, refactor config to env vars (Day 04)`

---

## 💡 Pelajaran Baru

**1. `GITHUB_TOKEN` itu otomatis, ga perlu bikin manual.** Setiap workflow run di GitHub Actions udah punya token sendiri. Tinggal pakai `${{ secrets.GITHUB_TOKEN }}`. Expire-nya cuma selama workflow run tersebut — makin aman karena ga ada token permanen yang bisa bocor.

**2. `fetch-depth: 0` itu penting buat Gitleaks.** Tanpa ini, Gitleaks cuma scan commit terakhir. Dengan `fetch-depth: 0`, seluruh git history di-scan, jadi secret yang mungkin masuk di commit lama juga bakal ketangkap.

**3. Hardcoded config itu debt teknis.** Walaupun sekarang cuma port yang dipindah ke env var, pattern `configs.Load()` ini bikin gampang nanti waktu mau nambah database connection string, API key, dan lain-lain. Cukup tambah field di struct dan env var, ga perlu ubah kode bisnis.

**4. Job paralel bikin pipeline lebih cepat.** `secret-scan` ga perlu nunggu `build-and-test` selesai. Keduanya jalan bersamaan. Jadi kalo ada secret leak, pipeline langsung gagal tanpa nunggu build selesai.

---

## 🔗 Referensi

- [Gitleaks GitHub Action](https://github.com/gitleaks/gitleaks-action)
- [GitHub Actions - GITHUB_TOKEN](https://docs.github.com/en/actions/security-guides/automatic-token-authentication)
- [12-Factor App - Config](https://12factor.net/config)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline makin lengkap! |
| Pemahaman materi | 4 | Env vars dan CI integration paham |
| Progres sesuai target | 5 | Sesuai silabus hari ke-4 |

---

## ➡️ Rencana Besok

- [ ] Install Trivy lokal
- [ ] Jalankan `trivy fs .` buat scan dependensi di `go.mod`
- [ ] Identifikasi CVE dari dependensi
- [ ] Simpan report JSON ke `security/`

---

*[← Hari 03](hari-03.md) | [Hari 05 →](hari-05.md)*