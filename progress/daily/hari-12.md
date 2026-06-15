# Hari 12 — Pipeline Optimization

**📅 Tanggal:** 2026-06-15  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Tambah Go module caching di Build & Test job
- [x] Tambah Gitleaks binary caching di Secret Scan job
- [x] Tambah Trivy DB caching di SCA Scan job
- [x] Tambah Semgrep rules caching di SAST Scan job
- [x] Verifikasi YAML syntax benar

---

## ✅ Yang Berhasil Dikerjakan

- Optimasi 4 job di CI pipeline dengan caching:
  1. **Build & Test** — `cache: true` di `actions/setup-go@v5` + `-count=1` di test
  2. **Secret Scan (Gitleaks)** — Cache binary di `/usr/local/bin/gitleaks` dengan key `gitleaks-8.30.1`
  3. **SCA Scan (Trivy)** — Cache DB di `~/.cache/trivy` + pre-download DB step
  4. **SAST Scan (Semgrep)** — Cache rules di `~/.semgrep` dengan key baseret pada `.semgrep.yml` hash
- Semua 4 job tetap paralel (no dependencies between them)
- Semua caching menggunakan `actions/cache@v4`

---

## 📝 Catatan Teknis

```yaml
# Build & Test: Go module cache
- uses: actions/setup-go@v5
  with:
    go-version: '1.26'
    cache: true  # auto-caches ~/go/pkg/mod and ~/.cache/go-build

# Secret Scan: Gitleaks binary cache
- uses: actions/cache@v4
  with:
    path: /usr/local/bin/gitleaks
    key: gitleaks-8.30.1
- if: steps.cache-gitleaks.outputs.cache-hit != 'true'
  # hanya download jika belum cache

# SCA Scan: Trivy DB cache
- uses: actions/cache@v4
  with:
    path: ~/.cache/trivy
    key: trivy-db-${{ github.run_id }}
    restore-keys: trivy-db-
- run: trivy image --download-db-only 2>/dev/null || true

# SAST Scan: Semgrep rules cache
- uses: actions/cache@v4
  with:
    path: ~/.semgrep
    key: semgrep-rules-${{ hashFiles('.semgrep.yml') }}
```

**Perubahan lain:**
- `go test` ditambah `-count=1` untuk bypass local cache di CI

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Gitleaks binary harus di-cache per versi | Gunakan `key: gitleaks-8.30.1` supaya cache invalid saat upgrade versi |
| Trivy DB perlu di-download sebelum scan | Tambah step `trivy image --download-db-only` sebelum scan, dengan `|| true` supaya gak fail kalau DB sudah cached |
| Semgrep rules gak perlu di-re-download tiap run | Cache `~/.semgrep` dengan key berdasarkan hash `.semgrep.yml` supaya cache invalid kalau custom rules berubah |

---

## 📤 Output Hari Ini

- [x] CI pipeline dioptimasi dengan 4 layer caching
- [x] Semua job tetap paralel dan independen
- [x] Commit: `ci: optimize pipeline with Go module cache, Gitleaks cache, Trivy DB cache, Semgrep rules cache (Day 12)`

---

## 💡 Pelajaran Baru

- **`actions/setup-go@v5` sudah built-in Go module caching.** Cukup `cache: true`, gak perlu manual `actions/cache` buat Go modules.
- **Cache key per versi itu penting buat tools.** Gitleaks key `gitleaks-8.30.1` berarti kalau upgrade ke 8.31.0, cache otomatis invalid dan binary baru di-download.
- **Semgrep action punya internal caching juga**, tapi caching explicit di `~/.semgrep` memberikan kontrol lebih.
- **`-count=1` di `go test`** mencegah Go menggunakan cached test results dari run sebelumnya, memastikan test selalu fresh di CI.
- **Trivy DB download** bisa memakan waktu 5-10 detik tiap run. Dengan caching, DB yang sudah di-download bisa di-reuse.

---

## 🔗 Referensi

- [actions/setup-go caching](https://github.com/actions/setup-go#caching-dependency-files)
- [actions/cache documentation](https://github.com/actions/cache)
- [Trivy DB caching in CI](https://aquasecurity.github.io/trivy/latest/docs/advanced/cache/)
- [Semgrep in CI optimization](https://semgrep.dev/docs/semgrep-ci/sample-ci-configs/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Optimization day — lebih teknis tapi penting |
| Pemahaman materi | 5 | Paham mekanisme GitHub Actions caching |
| Progres sesuai target | 5 | Caching pada 4 job selesai |

---

## ➡️ Rencana Besok

- [ ] Hari 13: Threat Modeling (STRIDE) — diagram arsitektur + tabel ancaman di `threat-model/`

---

*[← Hari 11](hari-11.md) | [Hari 13 →](hari-13.md)*