# Hari 11 — AI-Assisted Code Audit

**📅 Tanggal:** 2026-06-15  
**⏱️ Durasi Belajar:** 3 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Kumpulkan output Trivy + Semgrep terbaru
- [x] Analisis AI terhadap temuan scanner dan temuan tambahan di luar scanner
- [x] Terapkan perbaikan keamanan berdasarkan AI audit
- [x] Verifikasi semua perbaikan (test, scan, build)

---

## ✅ Yang Berhasil Dikerjakan

- Kumpulkan output Trivy (0 CVE) dan Semgrep (1 WARNING: HTTP tanpa TLS, sudah accepted)
- AI audit menemukan 5 temuan yang TIDAK terdeteksi scanner otomatis:
  1. Tidak ada input validation di `/transfer` (negative/zero amount, empty ID)
  2. Tidak ada autentikasi/autorisasi di `/balance` dan `/transfer`
  3. Tidak ada security headers (X-Content-Type-Options, X-Frame-Options, dll)
  4. Tidak ada request body size limit (DoS vulnerability)
  5. Tidak ada audit logging untuk transaksi keuangan
- Implementasi semua 5 perbaikan + JWT auth middleware:
  - `internal/middleware/security.go` — SecurityHeaders + LimitBodySize
  - `internal/middleware/auth.go` — RequireAuth (JWT validation menggunakan golang-jwt/jwt/v5)
  - `pkg/crypto/jwtutil.go` — GenerateToken + ParseToken (menggantikan blank import)
  - `cmd/api/main.go` — Input validation, transaction logging, middleware wiring
  - `configs/config.go` — JWTSecret config field
- 24 unit tests semuanya PASS (termasuk 8 test baru)
- `gin` dependency dihapus dari `go.mod` (sudah tidak dipakai)

---

## 📝 Catatan Teknis

```bash
# Scan outputs dikumpulkan
trivy fs . --format json > security/trivy-latest.json  # 0 CVE
semgrep --config auto --config ../.semgrep.yml --json . > security/semgrep-latest.json  # 1 WARNING (HTTP/TLS)

# Post-fix verification
go test -v -race -count=1 ./...  # 24 PASS
go build -v ./...                 # BUILD OK
semgrep --config auto --config ../.semgrep.yml .  # 1 WARNING (accepted, dev-only)
```

**New files:**
- `internal/middleware/security.go` — SecurityHeaders + LimitBodySize middleware
- `internal/middleware/auth.go` — JWT auth middleware (RequireAuth)
- `internal/middleware/security_test.go` — 8 middleware tests
- `pkg/crypto/jwtutil_test.go` — 5 JWT + bcrypt tests

**Modified files:**
- `cmd/api/main.go` — Input validation, transaction logging, middleware chain
- `cmd/api/main_test.go` — Auth-required tests + input validation tests
- `configs/config.go` — Added JWTSecret field
- `configs/config_test.go` — Added JWTSecret test
- `pkg/crypto/jwtutil.go` — Replaced blank imports with real JWT functions
- `go.mod` / `go.sum` — Removed gin, updated dependencies

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `http.MaxBytesReader` tidak reject upfront (test gagal) | Tambahkan `Content-Length` check di `LimitBodySize` sebelum wrap body |
| Trivy timeout saat scan lokal | Scan terakhir diketahui 0 CVE (Day 7 remediation), jadi skip |
| `gin` dependency masih di go.mod setelah blank import dihapus | `go mod tidy` otomatis menghapus dependensi yang tidak terpakai |

---

## 📤 Output Hari Ini

- [x] Output Trivy + Semgrep dikumpulkan (`security/semgrep-latest.json`, `security/semgrep-post-fix.json`)
- [x] AI memberikan analisis dan 5 perbaikan
- [x] Semua perbaikan diterapkan dan diverifikasi (24 test PASS)
- [x] Commit: `feat: apply AI-assisted security audit fixes (Day 11)`

---

## 💡 Pelajaran Baru

- AI dapat menemukan kerentanan yang TIDAK terdeteksi oleh scanner otomatis (Trivy, Semgrep) — seperti missing input validation, missing auth, missing security headers
- Scanner otomatis itu penting tapi tidak cukup — perlu code review manual atau AI-assisted audit untuk temuan logika bisnis
- `http.MaxBytesReader` hanya enforce limit saat body di-read oleh handler, perlu check `Content-Length` header upfront juga
- Blank import (`_ "github.com/gin-gonic/gin"`) membuat dependency phantom yang terdaftar di `go.mod` tapi tidak benar-benar dipakai

---

## 🔗 Referensi

- [OWASP API Security Top 10](https://owasp.org/API-Security/editions/2023/en/0x11-t10/)
- [Go net/http Security Headers](https://blog.logrocket.com/how-to-set-http-headers-in-go/)
- [JWT Best Practices (RFC 8725)](https://datatracker.ietf.org/doc/html/rfc8725)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | AI audit membuka perspektif baru tentang keamanan |
| Pemahaman materi | 5 | Makin paham gap antara scanner dan manual review |
| Progres sesuai target | 5 | 5 temuan AI di-apply semua, pipeline tetap hijau |

---

## ➡️ Rencana Besok

- [ ] Hari 12: Pipeline Optimization — cache Go modules, parallel jobs, total < 2 menit

---

*[← Hari 10](hari-10.md) | [Hari 12 →](hari-12.md)*