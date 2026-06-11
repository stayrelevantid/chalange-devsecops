# Hari 07 — SCA Remediation

**📅 Tanggal:** 2026-06-11  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] `jwt-go` diganti dengan `golang-jwt/jwt/v5`
- [x] `trivy fs .` lokal: 0 CRITICAL/HIGH
- [x] Pipeline CI kembali hijau
- [x] `go mod tidy` bersih

---

## ✅ Yang Berhasil Dikerjakan

- Ganti `dgrijalva/jwt-go v3.2.0+incompatible` → `golang-jwt/jwt/v5 v5.3.1` (library yang masih maintained)
- Upgrade `gin` dari v1.9.0 → v1.12.0 (membawa transitive dependencies yang lebih baru)
- Upgrade `x/crypto` dari v0.5.0 → v0.48.0 (memperbaiki CVE-2024-45337 dan CVE-2025-22869)
- Upgrade `x/net` dari v0.7.0 → v0.51.0 (memperbaiki CVE-2023-39325)
- Trivy scan lokal: **0 vulnerabilities** — semua CVE ter-remediasi
- Semua 8 unit test masih hijau
- Pipeline CI kembali **hijau** — SCA Scan (Trivy) menemukan 0 CVE

---

## 📝 Catatan Teknis

```bash
# Ganti jwt-go (deprecated) dengan golang-jwt/jwt/v5
go get github.com/golang-jwt/jwt/v5@latest

# Upgrade gin ke versi terbaru
go get github.com/gin-gonic/gin@v1.12.0

# Upgrade semua transitive dependencies
go get -u ./...
go mod tidy

# Verifikasi lokal — harus 0CVE
trivy fs . --severity CRITICAL,HIGH --scanners vuln
```

### Versi sebelum dan sesudah:

| Library | Sebelum | Sesudah |
|---------|---------|---------|
| dgrijalva/jwt-go | v3.2.0+incompatible | **DIHAPUS** (deprecated) |
| golang-jwt/jwt/v5 | — | v5.3.1 (baru) |
| gin | v1.9.0 | v1.12.0 |
| x/crypto | v0.5.0 | v0.48.0 |
| x/net | v0.7.0 | v0.51.0 |
| x/sys | v0.5.0 | v0.41.0 |
| x/text | v0.7.0 | v0.34.0 |

### Perubahan kode:
```go
// pkg/crypto/jwtutil.go — SEBELUM
import (
    _ "github.com/dgrijalva/jwt-go"
    _ "github.com/gin-gonic/gin"
)

// pkg/crypto/jwtutil.go — SESUDAH
import (
    _ "github.com/gin-gonic/gin"
    _ "github.com/golang-jwt/jwt/v5"
)
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `go get -u` timeout karena koneksi lambat | Edit `go.mod` langsung untuk bump versi gin ke v1.12.0, lalu `go mod tidy` |
| `x/crypto` dan `x/net` tidak ikut ter-upgrade saat `go get -u gin` | Manual edit `go.mod` untuk versi minimum yang fix CVE, `go mod tidy` resolve ke versi terbaru yang kompatibel |

---

## 📤 Output Hari Ini

- [x] `dgrijalva/jwt-go` diganti `golang-jwt/jwt/v5 v5.3.1`
- [x] `gin` upgraded ke v1.12.0
- [x] Semua transitive dependencies upgraded
- [x] Trivy scan: **0 vulnerabilities** (dari 4 CVE sebelumnya)
- [x] `go.mod` dan `go.sum` bersih
- [x] 8/8 unit test hijau
- [x] Pipeline CI **hijau** (semua 3 job pass)
- [x] `security/trivy-fs-report.json` diupdate (0 findings)

---

## 💡 Pelajaran Baru

- **Library deprecated harus diganti total, bukan di-upgrade** — `jwt-go` tidak punya fixed version untuk CVE-2020-26160. Satu-satunya jalan adalah migrasi ke `golang-jwt/jwt/v5`
- **Transitive dependency ikut ter-upgrade saat library utama di-upgrade** — gin v1.12.0 membawa `x/crypto` v0.48.0 dan `x/net` v0.51.0 yang sudah memperbaiki CVE
- **`go mod tidy` resolve versi minimum ke versi terbaru kompatibel** — setelah manual set gin v1.12.0, `go mod tidy` mengupgrade semua transitive deps ke versi yang sesuai
- **Network timeout bisa diakali dengan manual edit `go.mod`** — daripada menunggu `go get -u` yang timeout, lebih efektif edit versi langsung lalu `go mod tidy`
- **Remediasi SCA itu about menghapus dan mengganti, bukan cuma patch** — jwt-go dihapus total dan diganti library baru yang berbeda package path-nya

---

## 🔗 Referensi

- [golang-jwt/jwt — GitHub](https://github.com/golang-jwt/jwt)
- [CVE-2020-26160 — jwt-go access restriction bypass](https://github.com/advisories/GHSA-w73w-5m7g-f7qc)
- [CVE-2024-45337 — golang.org/x/crypto auth bypass](https://nvd.nist.gov/vuln/detail/CVE-2024-45337)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Pipeline hijau lagi, sangat memuaskan! |
| Pemahaman materi | 4 | Remediasi SCA konsepnya clear |
| Progres sesuai target | 5 | Sesuai rencana, network issue bisa diakali |

---

## ➡️ Rencana Besok

- [ ] SAST Setup dengan Semgrep (Day 08)
- [ ] Install Semgrep dan scan kode lokal
- [ ] Temukan insecure crypto dan SQL injection pattern

---

*[← Hari 06](hari-06.md) | [Hari 08 →](hari-08.md)*