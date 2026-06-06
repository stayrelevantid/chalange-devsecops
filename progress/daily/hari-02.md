# Hari 02 — Pipeline CI/CD Dasar

**📅 Tanggal:** 2026-06-06  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Membuat GitHub Actions workflow `ci.yml` yang otomatis build & test setiap push
- [x] Pipeline trigger on push & PR
- [x] `go build` dan `go test` berjalan di CI
- [x] Coverage report ter-upload sebagai artifact
- [x] Update `go.mod` ke Go 1.26 (konsisten dengan lokal)

---

## ✅ Yang Berhasil Dikerjakan

- Membuat `.github/workflows/ci.yml` di root repo (bukan di subdirektori)
- Workflow menggunakan `working-directory: securebank-api` karena proyek ada di subdirektori
- Job label: `Build & Test (Go 1.26)` — jelas hari keberapa dan versi Go
- Update `go.mod` dari `go 1.25.0` ke `go 1.26.0` agar konsisten
- Commit dan push berhasil, workflow ter-trigger di GitHub Actions

---

## 📝 Catatan Teknis

```yaml
# .github/workflows/ci.yml — snippet penting
name: SecureBank CI - Day 02

jobs:
  build-and-test:
    name: Build & Test (Go 1.26)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: securebank-api  # proyek ada di subdirektori

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go 1.26
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Build
        run: go build -v ./...

      - name: Test with Race Detection & Coverage
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload Coverage Report
        if: always()  # upload meski test gagal
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: securebank-api/coverage.out
```

**Keputusan penting:**
- Workflow diletakkan di root repo (`.github/workflows/ci.yml`) karena GitHub Actions hanya membaca dari root
- `working-directory: securebank-api` digunakan di semua step agar build & test jalan di folder yang benar
- Go version di hardcode `1.26` (bukan `1.22` dari tutorial referensi) agar sesuai dengan versi lokal
- `if: always()` di upload coverage supaya artifact tetap ter-upload meski test gagal

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Tutorial referensi pakai Go 1.22, tapi lokal pakai 1.26 | Disesuaikan ke `1.26` di workflow dan `go.mod` |
| Proyek ada di subdirektori `securebank-api/` | Gunakan `defaults.run.working-directory: securebank-api` di semua step |
| Path coverage artifact harus relatif dari root repo | Set `path: securebank-api/coverage.out` |

---

## 📤 Output Hari Ini

- [x] `.github/workflows/ci.yml` — GitHub Actions workflow CI
- [x] `securebank-api/go.mod` — update ke `go 1.26.0`
- [x] Pipeline hijau di GitHub Actions (on push & PR)
- [x] Coverage report ter-upload sebagai artifact
- [x] Commit: `5e41583 ci: add GitHub Actions CI pipeline (Day 02)`

---

## 💡 Pelajaran Baru

- GitHub Actions hanya membaca `.github/workflows/` dari root repo — walau proyek ada di subdirektori, workflow harus di root
- `defaults.run.working-directory` sangat berguna buat monorepo atau subdirektori, menghindari `cd` manual di setiap step
- `if: always()` di step upload artifact memastikan report tetap tersedia meski test gagal — penting buat debugging
- Label job yang jelas (`Build & Test (Go 1.26)`) bikin lebih gampang identifikasi di GitHub Actions UI

---

## 🔗 Referensi

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [actions/setup-go](https://github.com/actions/setup-go)
- [actions/upload-artifact](https://github.com/actions/upload-artifact)
- [GitHub Actions working-directory](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idstepsrun)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline jalan, lanjut! |
| Pemahaman materi | 4 | Konsep CI/CD dasar paham |
| Progres sesuai target | 5 | Sesuai silabus hari ke-2 |

---

## ➡️ Rencana Besok

- [ ] Install Gitleaks lokal
- [ ] Jalankan Gitleaks detect di repo
- [ ] Buat `.gitleaks.toml` config
- [ ] Test deteksi secret secara lokal

---

*[← Hari 01](hari-01.md) | [Hari 03 →](hari-03.md)*