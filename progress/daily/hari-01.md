# Hari 01 — Setup Repo & Golang API

**📅 Tanggal:** 2026-06-05  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Membuat repositori Git dengan `.gitignore` untuk Go
- [x] Membuat 3 endpoint API: `/health`, `/balance`, `/transfer`
- [x] Menulis unit test dan memastikan hijau
- [x] Commit pertama ke repo

---

## ✅ Yang Berhasil Dikerjakan

- Inisialisasi repositori `securebank-api/` sebagai subdirektori di dalam `chalange-devsecops/`
- Membuat struktur folder Go standar: `cmd/api/`, `internal/`, `pkg/`, `configs/`, `security/`
- Menulis `cmd/api/main.go` dengan 3 endpoint fungsional
- Menulis `cmd/api/main_test.go` dengan 6 unit test (semua PASS)
- Verifikasi API server berjalan di `:8080` dan semua endpoint merespon dengan benar
- Commit pertama: `feat: initial SecureBank API with 3 endpoints`

---

## 📝 Catatan Teknis

```bash
# Inisialisasi Go module
go mod init github.com/stayrelevantid/securebank-api

# Jalankan unit test
go test ./cmd/api/ -v
# 6/6 PASS

# Build & jalankan server
go build -o /tmp/securebank-api ./cmd/api/
/tmp/securebank-api

# Test endpoint
curl http://localhost:8080/health
# {"status":"healthy"}

curl "http://localhost:8080/balance?id=ACC001"
# {"id":"ACC001","name":"Alice","balance":10000}

curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"ACC001","to":"ACC002","amount":500}'
# {"status":"success"}
```

**Struktur folder yang dibuat:**
```
securebank-api/
├── cmd/api/
│   ├── main.go          # 3 endpoint API
│   └── main_test.go     # 6 unit test
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── pkg/crypto/
├── configs/
├── security/threat-model/
├── .github/workflows/
├── .gitignore
└── go.mod
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Tidak ada hambatan signifikan | Semua berjalan lancar |

---

## 📤 Output Hari Ini

- [x] `securebank-api/cmd/api/main.go` — 3 endpoint `/health`, `/balance`, `/transfer`
- [x] `securebank-api/cmd/api/main_test.go` — 6 unit test hijau
- [x] `securebank-api/.gitignore` — konfigurasi Git untuk Go
- [x] `securebank-api/go.mod` — Go module terinisialisasi
- [x] Commit: `6ee1d7e feat: initial SecureBank API with 3 endpoints`

---

## 💡 Pelajaran Baru

- Struktur folder Go standar: `cmd/`, `internal/`, `pkg/` memisahkan concerns dengan baik
- `sync.RWMutex` penting untuk concurrent access ke map accounts
- Unit test dengan `httptest.NewRecorder()` dan `httptest.NewRequest()` sangat praktis untuk HTTP handler testing

---

## 🔗 Referensi

- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [net/http package](https://pkg.go.dev/net/http)
- [testing package](https://pkg.go.dev/testing)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Hari pertama, semangat tinggi! |
| Pemahaman materi | 4 | Dasar Go API sudah dipahami |
| Progres sesuai target | 5 | Semua checklist tercapai |

---

## ➡️ Rencana Besok

- [ ] Membuat GitHub Actions workflow `ci.yml` (auto build & test on push)
- [ ] Pipeline trigger on push & PR
- [ ] Coverage report ter-upload sebagai artifact

---

*[← Tracker](../tracker.md) | [Hari 02 →](hari-02.md)*