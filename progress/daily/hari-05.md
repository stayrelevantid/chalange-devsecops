# Hari 05 — SCA Setup (Trivy FS)

**📅 Tanggal:** 2026-06-09  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Install Trivy dan jalankan scan filesystem
- [x] Tambahkan dependensi rentan (gin v1.9.0, jwt-go v3.2.0)
- [x] Scan menemukan CVE dari library yang rentan
- [x] Buat `.trivyignore` untuk dokumentasi accepted risks
- [x] Simpan report JSON di `security/`

---

## ✅ Yang Berhasil Dikerjakan

- Install Trivy v0.71.0 via Homebrew (setelah beberapa kali gagal download binary langsung karena timeout)
- Tambahkan dependensi rentan: `github.com/gin-gonic/gin@v1.9.0` dan `github.com/dgrijalva/jwt-go@v3.2.0+incompatible`
- Buat `pkg/crypto/jwtutil.go` dengan blank import agar `go mod tidy` tidak menghapus dependensi
- Jalankan `trivy fs . --severity HIGH,CRITICAL` — menemukan **4 CVE** (1 CRITICAL, 3 HIGH)
- Buat `.trivyignore` di root repo untuk accepted risks
- Simpan report JSON di `securebank-api/security/trivy-fs-report.json`
- Semua 8 unit test masih hijau setelah penambahan dependensi

---

## 📝 Catatan Teknis

```bash
# Install Trivy
brew install trivy

# Tambah dependensi rentan
go get github.com/gin-gonic/gin@v1.9.0
go get github.com/dgrijalva/jwt-go@v3.2.0+incompatible

# Scan filesystem
trivy fs . --severity HIGH,CRITICAL

# Simpan report JSON
trivy fs . --severity HIGH,CRITICAL --format json --output security/trivy-fs-report.json
```

### Temuan Trivy:
| Library | CVE | Severity | Fixed Version |
|---------|-----|----------|---------------|
| dgrijalva/jwt-go | CVE-2020-26160 | HIGH | — (deprecated) |
| golang.org/x/crypto | CVE-2024-45337 | CRITICAL | 0.31.0 |
| golang.org/x/crypto | CVE-2025-22869 | HIGH | 0.35.0 |
| golang.org/x/net | CVE-2023-39325 | HIGH | 0.17.0 |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Download binary Trivy v0.58.0 gagal (404) | Versi tidak ada, gunakan v0.71.0 (latest) |
| Binary v0.71.0 timeout saat download (45MB, koneksi lambat) | Gunakan `brew install trivy` — berhasil |
| `go mod tidy` menghapus dependensi yang tidak diimport | Buat `pkg/crypto/jwtutil.go` dengan blank import |

---

## 📤 Output Hari Ini

- [x] Trivy v0.71.0 terinstall
- [x] `go.mod` + `go.sum` berisi gin v1.9.0 dan jwt-go v3.2.0
- [x] Scan menemukan 4 CVE (1 CRITICAL, 3 HIGH)
- [x] `.trivyignore` di root repo
- [x] `securebank-api/security/trivy-fs-report.json` (41KB)
- [x] `pkg/crypto/jwtutil.go` — blank import untuk dependensi

---

## 💡 Pelajaran Baru

- **Blank import** (`_ "package"`) bisa menjaga dependensi tetap di `go.mod` meskipun tidak langsung digunakan di kode
- **Trivy** mengenal Go modules dan otomatis mengidentifikasi CVE dari `go.mod`/`go.sum`
- **jwt-go** sudah deprecated dan tidak ada fixed version — harus migrasi ke `golang-jwt/jwt/v5` (akan dilakukan di Day 07)
- **SCA (Software Composition Analysis)** berbeda dengan SAST: SCA memindai dependensi pihak ketiga, SAST memindai kode sendiri

---

## 🔗 Referensi

- [Trivy Documentation](https://trivy.dev/)
- [CVE-2020-26160 — jwt-go](https://github.com/advisories/GHSA-w73w-5m7g-f7qc)
- [CVE-2024-45337 — golang.org/x/crypto](https://nvd.nist.gov/vuln/detail/CVE-2024-45337)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Trivy jadi tool yang sangat berguna |
| Pemahaman materi | 4 | SCA konsepnya straightforward |
| Progres sesuai target | 4 | Sesuai rencana meskipun ada hambatan koneksi |

---

## ➡️ Rencana Besok

- [ ] Integrasikan Trivy ke CI pipeline (Day 06 - SCA Pipeline)
- [ ] Tambah job `sca-scan` di `.github/workflows/ci.yml`
- [ ] Quality gate: pipeline gagal jika ada CVE CRITICAL

---

*[← Hari 04](hari-04.md) | [Hari 06 →](hari-06.md)*