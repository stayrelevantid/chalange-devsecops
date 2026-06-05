# Fase 1: Secure SDLC & Application Security (Hari 1–15)

> **Proyek:** SecureBank API — REST API perbankan sederhana (Golang)
> **Output Fase:** Pipeline CI/CD lengkap dengan Secret Scan + SCA + SAST yang memblokir kode tidak aman.

---

## Hari 1: Setup Repositori & Aplikasi SecureBank API

### Tujuan
Membuat fondasi proyek: repositori Git, struktur folder Golang, dan 3 endpoint API fungsional.

### Tutorial

**1. Inisialisasi repositori:**
```bash
mkdir securebank-api && cd securebank-api
git init
go mod init github.com/<username>/securebank-api
```

**2. Buat struktur folder:**
```bash
mkdir -p cmd/api internal/{handler,middleware,model,repository,service} pkg/crypto configs docs security
```

**3. Tulis `cmd/api/main.go`:**
```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
)

type Account struct {
    ID      string  `json:"id"`
    Name    string  `json:"name"`
    Balance float64 `json:"balance"`
}

type TransferReq struct {
    From   string  `json:"from"`
    To     string  `json:"to"`
    Amount float64 `json:"amount"`
}

var (
    accounts = map[string]*Account{
        "ACC001": {ID: "ACC001", Name: "Alice", Balance: 10000},
        "ACC002": {ID: "ACC002", Name: "Bob", Balance: 5000},
    }
    mu sync.RWMutex
)

func getBalance(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    mu.RLock()
    acc, ok := accounts[id]
    mu.RUnlock()
    if !ok {
        http.Error(w, "account not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(acc)
}

func transfer(w http.ResponseWriter, r *http.Request) {
    var req TransferReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    mu.Lock()
    defer mu.Unlock()
    from, ok1 := accounts[req.From]
    to, ok2 := accounts[req.To]
    if !ok1 || !ok2 {
        http.Error(w, "account not found", http.StatusNotFound)
        return
    }
    if from.Balance < req.Amount {
        http.Error(w, "insufficient balance", http.StatusBadRequest)
        return
    }
    from.Balance -= req.Amount
    to.Balance += req.Amount
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
    http.HandleFunc("/health", healthCheck)
    http.HandleFunc("/balance", getBalance)
    http.HandleFunc("/transfer", transfer)
    log.Println("SecureBank API running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**4. Tulis unit test `cmd/api/main_test.go`:**
```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthCheck(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    healthCheck(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestGetBalanceNotFound(t *testing.T) {
    req := httptest.NewRequest("GET", "/balance?id=INVALID", nil)
    w := httptest.NewRecorder()
    getBalance(w, req)
    if w.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", w.Code)
    }
}
```

**5. Jalankan dan verifikasi:**
```bash
go run cmd/api/main.go &
curl http://localhost:8080/health
curl "http://localhost:8080/balance?id=ACC001"
go test ./cmd/api/ -v
```

### Checklist
- [ ] Repo Git terinisialisasi dengan `.gitignore` untuk Go
- [ ] 3 endpoint berfungsi: `/health`, `/balance`, `/transfer`
- [ ] Unit test berjalan hijau
- [ ] Commit pertama: `feat: initial SecureBank API with 3 endpoints`

---

## Hari 2: Build Pipeline CI/CD Dasar

### Tujuan
Membuat GitHub Actions workflow yang otomatis build dan test setiap push.

### Tutorial

**1. Buat `.github/workflows/ci.yml`:**
```yaml
name: SecureBank CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload Coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.out
```

**2. Push dan verifikasi:**
```bash
git add .
git commit -m "ci: add basic build and test pipeline"
git push origin main
```

**3. Cek GitHub → Actions tab → pipeline harus hijau.**

### Checklist
- [ ] Pipeline trigger on push & PR
- [ ] `go build` dan `go test` berjalan di CI
- [ ] Coverage report ter-upload sebagai artifact

---

## Hari 3: Secret Scanning (Gitleaks) — Lokal

### Tujuan
Mendeteksi kredensial yang bocor di kode menggunakan Gitleaks.

### Tutorial

**1. Install Gitleaks:**
```bash
brew install gitleaks    # macOS
# atau download binary dari https://github.com/gitleaks/gitleaks/releases
```

**2. Sengaja tambahkan secret untuk uji coba (di branch terpisah):**
```bash
git checkout -b test/secret-leak
```

Buat file `configs/database.go`:
```go
package configs

// INTENTIONAL — ini akan dihapus setelah pengujian
const (
    DBPassword = "SuperSecret123!"
    AWSKey     = "AKIAIOSFODNN7EXAMPLE"
    AWSSecret  = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
)
```

**3. Jalankan Gitleaks:**
```bash
gitleaks detect --source . -v
```

**4. Buat config `.gitleaks.toml`:**
```toml
title = "SecureBank Gitleaks Config"

[allowlist]
  paths = [
    '''go\.sum''',
    '''vendor/''',
  ]
```

**5. Jalankan ulang:**
```bash
gitleaks detect --source . -v --config .gitleaks.toml
```

**6. Bersihkan: hapus file test, kembali ke main:**
```bash
git checkout main
git branch -D test/secret-leak
```

### Checklist
- [ ] Gitleaks terinstall dan berjalan lokal
- [ ] Berhasil mendeteksi AWS key dan password hardcoded
- [ ] `.gitleaks.toml` dikonfigurasi untuk exclude false positive
- [ ] Branch test dihapus (secret tidak masuk main)

---

## Hari 4: Integrasi Secret Scan di Pipeline + Remediasi

### Tujuan
Menjalankan Gitleaks otomatis di CI/CD dan memindahkan semua kredensial ke environment variables.

### Tutorial

**1. Tambahkan job ke `.github/workflows/ci.yml`:**
```yaml
  secret-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # full history untuk scan commit lama

      - name: Run Gitleaks
        uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**2. Refactor kode — pindahkan config ke env vars:**

Buat `configs/config.go`:
```go
package configs

import "os"

type Config struct {
    DBHost     string
    DBPassword string
    Port       string
}

func Load() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        Port:       getEnv("PORT", "8080"),
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
```

**3. Tambahkan secrets di GitHub → Settings → Secrets:**
- `DB_PASSWORD`
- Nilai credential lain yang dibutuhkan

**4. Commit dan push:**
```bash
git add .
git commit -m "security: integrate gitleaks in pipeline, move secrets to env vars"
git push
```

### Checklist
- [ ] Job `secret-scan` berjalan di pipeline
- [ ] Tidak ada hardcoded credential di source code
- [ ] `configs/config.go` membaca dari environment variables
- [ ] Pipeline hijau

---

## Hari 5: SCA Setup — Trivy Filesystem Scan

### Tujuan
Memindai dependensi Go (`go.mod`/`go.sum`) untuk menemukan pustaka dengan CVE.

### Tutorial

**1. Install Trivy:**
```bash
brew install trivy    # macOS
```

**2. Tambahkan beberapa dependensi (termasuk yang rentan untuk latihan):**
```bash
go get github.com/gin-gonic/gin@v1.9.0
go get github.com/dgrijalva/jwt-go@v3.2.0  # known vulnerable
```

**3. Jalankan scan filesystem:**
```bash
trivy fs . --severity HIGH,CRITICAL
```

**4. Buat `.trivyignore` untuk false positive:**
```
# Accepted risks — documented
# CVE-XXXX-YYYY — not applicable because ...
```

**5. Simpan output sebagai referensi:**
```bash
trivy fs . --format json --output security/trivy-fs-report.json
```

### Checklist
- [ ] Trivy terinstall dan berjalan lokal
- [ ] Scan menemukan CVE dari `jwt-go` (library deprecated)
- [ ] `.trivyignore` dibuat
- [ ] Report JSON tersimpan di `security/`

---

## Hari 6: SCA Pipeline Integration (Quality Gate)

### Tujuan
Trivy di pipeline gagalkan build jika ada CVE CRITICAL.

### Tutorial

**Tambahkan job ke `.github/workflows/ci.yml`:**
```yaml
  sca-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Trivy SCA
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'json'
          output: 'trivy-sca.json'
          severity: 'CRITICAL,HIGH'
          exit-code: '1'  # gagalkan pipeline jika ditemukan

      - name: Upload Trivy Report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: trivy-sca-report
          path: trivy-sca.json
```

### Checklist
- [ ] Pipeline gagal/merah karena CVE CRITICAL ditemukan
- [ ] Report JSON ter-upload sebagai artifact
- [ ] `--exit-code 1` berfungsi sebagai quality gate

---

## Hari 7: SCA Remediation — Patch Dependensi

### Tujuan
Memperbaiki pustaka rentan yang ditemukan Trivy, pipeline kembali hijau.

### Tutorial

**1. Ganti library yang deprecated:**
```bash
# Ganti jwt-go (deprecated) dengan golang-jwt
go get -u github.com/golang-jwt/jwt/v5
```

**2. Update semua dependensi:**
```bash
go get -u ./...
go mod tidy
```

**3. Verifikasi lokal:**
```bash
trivy fs . --severity CRITICAL
# harus 0 findings
```

**4. Commit dan push:**
```bash
git add go.mod go.sum
git commit -m "fix: patch vulnerable dependencies (jwt-go → golang-jwt)"
git push
```

### Checklist
- [ ] `jwt-go` diganti dengan `golang-jwt/jwt/v5`
- [ ] `trivy fs .` lokal: 0 CRITICAL
- [ ] Pipeline CI kembali hijau
- [ ] `go mod tidy` bersih

---

## Hari 8: SAST Setup — Semgrep Lokal

### Tujuan
Menganalisis kode Golang secara statis untuk menemukan kerentanan (SQL injection, insecure crypto, dll).

### Tutorial

**1. Install Semgrep:**
```bash
pip3 install semgrep
# atau
brew install semgrep
```

**2. Sengaja tambahkan kode rentan untuk latihan:**

Buat `pkg/crypto/hash.go`:
```go
package crypto

import (
    "crypto/md5"   // INSECURE — intentional for training
    "encoding/hex"
)

// HashPassword uses MD5 — this is intentionally insecure for SAST demo
func HashPassword(password string) string {
    h := md5.Sum([]byte(password))
    return hex.EncodeToString(h[:])
}
```

**3. Jalankan Semgrep:**
```bash
semgrep --config auto . --lang go
semgrep --config "p/golang" .
```

**4. Buat custom rule `.semgrep.yml`:**
```yaml
rules:
  - id: no-md5-usage
    patterns:
      - pattern: crypto/md5
    message: "MD5 is cryptographically broken. Use SHA-256 or bcrypt."
    languages: [go]
    severity: ERROR
```

### Checklist
- [ ] Semgrep terinstall dan berjalan
- [ ] Menemukan penggunaan MD5 yang tidak aman
- [ ] Custom rule `.semgrep.yml` dibuat
- [ ] Output dipahami (severity, file, line number)

---

## Hari 9: SAST Pipeline Integration (Quality Gate)

### Tujuan
Semgrep di pipeline memblokir PR jika ditemukan kerentanan HIGH.

### Tutorial

**Tambahkan job ke `.github/workflows/ci.yml`:**
```yaml
  sast-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Semgrep
        uses: semgrep/semgrep-action@v1
        with:
          config: >-
            p/golang
            p/owasp-top-ten
          generateSarif: "1"

      - name: Upload SARIF
        if: always()
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: semgrep.sarif
```

### Checklist
- [ ] Job SAST berjalan di pipeline
- [ ] PR di-blokir jika ada finding HIGH+
- [ ] SARIF report terintegrasi dengan GitHub Security tab

---

## Hari 10: SAST Remediation — Fix Kode

### Tujuan
Memperbaiki semua temuan SAST: ganti MD5 dengan bcrypt, parameterized query.

### Tutorial

**1. Fix `pkg/crypto/hash.go`:**
```go
package crypto

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword uses bcrypt — industry standard
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword verifies a password against its hash
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**2. Tambahkan dependency:**
```bash
go get golang.org/x/crypto/bcrypt
```

**3. Verifikasi:**
```bash
semgrep --config auto . --lang go
# 0 findings
go test ./... -v
```

### Checklist
- [ ] MD5 diganti dengan bcrypt
- [ ] Semgrep lokal: 0 HIGH findings
- [ ] Unit test tetap hijau
- [ ] Pipeline hijau

---

## Hari 11: AI-Assisted Code Audit

### Tujuan
Menggunakan AI untuk menganalisis output scanner dan menghasilkan perbaikan berkualitas.

### Tutorial

**1. Ambil output dari scan sebelumnya:**
```bash
trivy fs . --format json > security/trivy-latest.json
semgrep --config auto . --json > security/semgrep-latest.json
```

**2. Berikan ke AI dengan prompt terstruktur:**
```
Kamu adalah security engineer senior. Analisis laporan keamanan berikut dari aplikasi Golang:

[paste JSON output]

Untuk setiap temuan:
1. Jelaskan risiko nyata dalam konteks REST API perbankan
2. Berikan kode perbaikan yang mengikuti OWASP best practice
3. Sebutkan dampak jika tidak diperbaiki
```

**3. Terapkan perbaikan yang valid, commit:**
```bash
git add .
git commit -m "security: apply AI-assisted audit fixes"
```

### Checklist
- [ ] Output Trivy + Semgrep dikumpulkan
- [ ] AI memberikan analisis dan perbaikan
- [ ] Perbaikan diterapkan dan diverifikasi
- [ ] Commit tercatat

---

## Hari 12: Pipeline Optimization & Caching

### Tujuan
Pipeline total berjalan di bawah 2 menit dengan caching dan parallel jobs.

### Tutorial

**Update `.github/workflows/ci.yml` — tambahkan caching:**
```yaml
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true  # auto-cache Go modules

      - name: Build
        run: go build -v ./...

      - name: Test with Race Detection
        run: go test -v -race -count=1 ./...
```

**Jalankan jobs secara paralel (bukan sequential):**
```yaml
jobs:
  build-and-test:
    # ...
  secret-scan:
    # ... (tidak depends on build)
  sca-scan:
    # ... (tidak depends on build)
  sast-scan:
    # ... (tidak depends on build)
```

### Checklist
- [ ] Go module cache aktif
- [ ] Security scan jobs berjalan paralel
- [ ] Total pipeline < 2 menit
- [ ] Semua jobs tetap hijau

---

## Hari 13: Threat Modeling (STRIDE)

### Tujuan
Mengidentifikasi ancaman terhadap SecureBank API menggunakan metodologi STRIDE.

### Tutorial

**1. Gambar arsitektur di Draw.io atau Mermaid:**

Buat `security/threat-model/architecture.md`:
```markdown
# SecureBank API — Threat Model

## Arsitektur
Client → [API Gateway] → [SecureBank API] → [Database]

## Analisis STRIDE

| Threat | Target | Deskripsi | Mitigasi |
|--------|--------|-----------|----------|
| **S**poofing | Auth endpoint | Attacker berpura-pura jadi user lain | JWT + rate limiting |
| **T**ampering | /transfer | Modifikasi amount di request body | Input validation + HMAC |
| **R**epudiation | Transaction log | User menyangkal melakukan transfer | Audit logging + timestamp |
| **I**nfo Disclosure | /balance | Data balance user lain bocor | Authorization check per user |
| **D**enial of Service | All endpoints | Flood request ke API | Rate limiter + WAF |
| **E**levation of Privilege | Admin endpoint | User biasa akses admin API | RBAC + middleware auth |
```

**2. Prioritaskan: urutkan berdasarkan risk score (Impact × Likelihood).**

### Checklist
- [ ] Diagram arsitektur dibuat
- [ ] 6 kategori STRIDE dianalisis
- [ ] Mitigasi diidentifikasi per ancaman
- [ ] Dokumen disimpan di `security/threat-model/`

---

## Hari 14: Intentional Vulnerability Test

### Tujuan
Memvalidasi bahwa pipeline benar-benar menangkap secret leak.

### Tutorial

**1. Buat branch test:**
```bash
git checkout -b test/intentional-leak
```

**2. Tambahkan fake secret:**
```bash
echo 'STRIPE_SECRET_KEY=sk_live_abc123verysecretkey' > configs/.env.production
git add . && git commit -m "test: intentional secret for pipeline validation"
git push origin test/intentional-leak
```

**3. Buka PR → pipeline HARUS gagal di job `secret-scan`.**

**4. Verifikasi error message, lalu tutup PR dan hapus branch:**
```bash
git checkout main
git branch -D test/intentional-leak
git push origin --delete test/intentional-leak
```

### Checklist
- [ ] Pipeline gagal mendeteksi secret leak
- [ ] Error message jelas menunjukkan file dan line
- [ ] Branch test dihapus sepenuhnya
- [ ] Validasi bahwa defense layer bekerja

---

## Hari 15: Dokumentasi Fase 1

### Tujuan
Merangkum semua yang dipelajari dalam dokumen teknis terstruktur.

### Tutorial

Tulis `docs/fase-1-appsec.md` dengan struktur:

```markdown
# Fase 1: Application Security — Dokumentasi

## 1. Arsitektur Pipeline
[Diagram pipeline: Push → Build → Test → Secret Scan → SCA → SAST]

## 2. Tools yang Digunakan
| Tool | Fungsi | Config File |
|------|--------|-------------|
| Gitleaks | Secret Scanning | .gitleaks.toml |
| Trivy | SCA (dependency scan) | .trivyignore |
| Semgrep | SAST (static analysis) | .semgrep.yml |

## 3. Quality Gates
- Secret Scan: BLOCK jika ada credential terdeteksi
- SCA: BLOCK jika ada CVE CRITICAL
- SAST: BLOCK jika ada finding HIGH+

## 4. Lessons Learned
- [Tuliskan 3-5 pelajaran penting]

## 5. Metrics
- Jumlah vulnerabilities ditemukan: X
- Jumlah vulnerabilities diperbaiki: Y
- Waktu pipeline rata-rata: Z menit
```

### Checklist
- [ ] Dokumen teknis Fase 1 lengkap
- [ ] Mencakup arsitektur, tools, quality gates
- [ ] Metrics tercatat
- [ ] Commit: `docs: fase 1 application security documentation`

---

> ✅ **Selesai Fase 1** — Lanjut ke [Fase 2: Infrastructure as Code & Container Security](fase-2-infra-container.md)
