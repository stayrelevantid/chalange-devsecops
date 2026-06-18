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

---

# Fase 1: Dokumentasi Retrospektif

> **Tanggal:** 2026-06-18  
> **Durasi:** Hari 1–14 (14 hari aktif)  
> **Status:** ✅ Selesai  

---

## 1. Arsitektur Pipeline

```mermaid
graph LR
    Push[git push / PR] --> Build[Build & Test]
    Push --> Secret[Secret Scan<br/>Gitleaks]
    Push --> SCA[SCA Scan<br/>Trivy]
    Push --> SAST[SAST Scan<br/>Semgrep]

    Build --> Check1{Tests Pass?}
    Secret --> Check2{Leaks Found?}
    SCA --> Check3{CVE Found?}
    SAST --> Check4{Findings?}

    Check1 -->|Yes| OK[✅ Pipeline Green]
    Check2 -->|No| OK
    Check3 -->|No CVE| OK
    Check4 -->|No ERROR| OK

    Check1 -->|No| FAIL[❌ Pipeline Red]
    Check2 -->|Yes| FAIL
    Check3 -->|CRITICAL/HIGH CVE| FAIL
    Check4 -->|ERROR finding| FAIL
```

**4 parallel jobs, semuanya harus pass untuk pipeline hijau:**

| Job | Tool | Purpose | Quality Gate |
|-----|------|---------|-------------|
| Build & Test | Go 1.26 | Compile dan unit test dengan race detection | `go test` exit code 0 |
| Secret Scan | Gitleaks 8.30.1 | Deteksi credential/API key yang bocor di git history | Any leak = block |
| SCA Scan | Trivy | Scan dependensi Go untuk CVE CRITICAL/HIGH | Any CRITICAL/HIGH CVE = block |
| SAST Scan | Semgrep | Analisis statis kode Go untuk insecure pattern | Any ERROR finding = block |

---

## 2. Tools yang Digunakan

| Tool | Fungsi | Config File | Versi | Sumber |
|------|--------|-------------|-------|--------|
| Go | Bahasa pemrograman API | `go.mod` | 1.26.0 | golang.org |
| Gitleaks | Secret scanning (git history) | `.gitleaks.toml` | 8.30.1 | Binary install via wget |
| Trivy | SCA — Software Composition Analysis | (tidak ada) | latest (action) | `aquasecurity/trivy-action@master` |
| Semgrep | SAST — Static Application Security Testing | `.semgrep.yml` | latest (action) | `semgrep/semgrep-action@v1` |
| GitHub Actions | CI/CD pipeline | `.github/workflows/ci.yml` | N/A | github.com |

**Dependencies Go:**

| Dependency | Versi | Fungsi |
|------------|-------|--------|
| `golang-jwt/jwt/v5` | 5.3.1 | JWT token generation & validation |
| `golang.org/x/crypto` | 0.53.0 | bcrypt password hashing |

---

## 3. Quality Gates

### 3.1 Secret Scan — Gitleaks

```yaml
# .github/workflows/ci.yml (secret-scan job)
- name: Run Gitleaks
  run: gitleaks detect --source . -v --config .gitleaks.toml --no-gitignore
```

**Gate:** Jika ada leak yang terdeteksi → pipeline gagal.

**Important:** Flag `--no-gitignore` ditambahkan di Day 14 setelah ditemukan bahwa Gitleaks default menghormati `.gitignore`, sehingga file yang di-force-add (`git add -f`) tetap di-skip.

**Allowlist** (`.gitleaks.toml`):
- `go.sum` — dependency checksums, bukan secret
- `vendor/` — third-party code
- `docs/`, `progress/` — dokumentasi
- `security/` — scan reports (JSON)
- `.gitleaks.toml`, `.gitignore`, `.semgrep.yml` — config files

### 3.2 SCA Scan — Trivy

```yaml
# .github/workflows/ci.yml (sca-scan job)
severity: 'CRITICAL,HIGH'
exit-code: '1'
```

**Gate:** Jika ada CVE dengan severity CRITICAL atau HIGH → pipeline gagal.

**State:** Setelah Day 7 remediation (jwt-go → golang-jwt, upgrade gin), 0 CVE ditemukan.

### 3.3 SAST Scan — Semgrep

```yaml
# .github/workflows/ci.yml (sast-scan job)
config: p/golang p/owasp-top-ten .semgrep.yml
```

**Gate:** Jika ada finding dengan severity ERROR → pipeline gagal (Semgrep default behavior).

**Custom rules** (`.semgrep.yml`):
- `no-md5-usage` — ERROR jika ada `md5.Sum(...)` di kode Go

**Accepted WARNING:** `use-tls` (HTTP tanpa TLS) — accepted karena environment development, production akan pakai reverse proxy.

### 3.4 Build & Test

```yaml
# .github/workflows/ci.yml (build-and-test job)
- name: Test with Race Detection & Coverage
  run: go test -v -race -count=1 -coverprofile=coverage.out ./...
```

**Gate:** Jika ada test yang gagal atau race condition terdeteksi → pipeline gagal.

**24 unit tests** mencakup:
- Handler tests (health, balance, transfer)
- Auth middleware tests (valid/invalid/missing JWT)
- Security header tests
- Body size limit tests
- Input validation tests (negative, zero, empty, NaN)
- Config tests (default values, env vars)
- JWT utility tests (generate, parse, expired, wrong secret)
- bcrypt tests (hash, check)

---

## 4. Security Improvements Applied

| Day | Finding | Fix | Status |
|-----|---------|-----|--------|
| 03 | Hardcoded credentials (DB_PASSWORD, AWS_KEY) | Refactored to env vars (`configs/config.go`) | ✅ Fixed |
| 04 | Gitleaks needed CI integration | Binary install via wget in CI | ✅ Fixed |
| 05 | 4 CVE in dependencies (1 CRITICAL, 3 HIGH) | Identified via `trivy fs` | ✅ Identified |
| 07 | jwt-go v3.2.0 (deprecated, CVE) → golang-jwt v5.3.1 | `replace` in go.mod | ✅ Fixed |
| 07 | gin v1.9.0 (CVE) → v1.12.0 | Upgrade | ✅ Fixed |
| 07 | x/crypto, x/net (CVE) → patched | Upgrade | ✅ Fixed |
| 08 | MD5 password hashing (insecure crypto) | Replaced with bcrypt | ✅ Fixed |
| 09 | HTTP without TLS | Accepted as WARNING (dev env) | ⚠️ Accepted |
| 11 | No input validation on `/transfer` | Validate negative/zero/NaN/empty amounts | ✅ Fixed |
| 11 | No authentication on `/balance` and `/transfer` | JWT Bearer auth middleware | ✅ Fixed |
| 11 | No security headers | SecurityHeaders middleware (5 headers) | ✅ Fixed |
| 11 | No body size limit | LimitBodySize middleware (1KB/4KB) | ✅ Fixed |
| 11 | No transaction logging | `log.Printf` for transfers | ✅ Fixed |
| 12 | No Go module caching in CI | `cache: true` in setup-go | ✅ Fixed |
| 12 | No Gitleaks binary caching | `actions/cache@v4` with version key | ✅ Fixed |
| 12 | No Trivy DB caching | Cache `~/.cache/trivy` with restore-keys | ✅ Fixed |
| 12 | No Semgrep rules caching | Cache `~/.semgrep` with hash-based key | ✅ Fixed |
| 14 | Gitleaks skips `.gitignore` files | Added `--no-gitignore` flag | ✅ Fixed |

---

## 5. Threat Model Summary (STRIDE + DREAD)

> Full analysis: `securebank-api/security/threat-model/architecture.md`

### Critical Findings (DREAD Score ≥ 8)

| # | Threat | DREAD Score | Category |
|---|--------|-------------|----------|
| 1 | No rate limiting | 9.6 | Denial of Service |
| 2 | Balance of any account (no user-scoped authz) | 8.4 | Info Disclosure |
| 3 | Account ID manipulation (transfer from any account) | 8.0 | Tampering |

### High Findings (DREAD Score 6-7)

| # | Threat | DREAD Score | Category |
|---|--------|-------------|----------|
| 4 | Hardcoded accounts (no RBAC) | 7.8 | Elevation of Privilege |
| 5 | No connection limit | 7.2 | Denial of Service |
| 6 | No RBAC | 6.6 | Elevation of Privilege |
| 7 | Default JWT secret fallback | 6.2 | Spoofing |
| 8 | No persistent audit log | 6.0 | Repudiation |

### Already Mitigated (6 findings)

- ✅ JWT token forgery → HMAC-SHA256 signing
- ✅ Input validation (negative/zero/NaN amounts)
- ✅ JWT payload manipulation → signature verification
- ✅ Large request body → LimitBodySize middleware
- ✅ Error message leakage → generic JSON errors
- ✅ Security headers → 5 headers implemented

---

## 6. Metrics

| Metrik | Nilai |
|--------|-------|
| Total vulnerabilities ditemukan (SCA) | 4 CVE (1 CRITICAL, 3 HIGH) |
| Total vulnerabilities diperbaiki (SCA) | 4/4 (100%) |
| SAST findings ditemukan | 2 (1 ERROR: MD5, 1 WARNING: no-TLS) |
| SAST findings diperbaiki | 1/2 ERROR fixed (MD5→bcrypt), 1 WARNING accepted |
| Secret leaks ditemukan | 2 (AWS key, DB password) |
| Secret leaks diperbaiki | 2/2 (moved to env vars) |
| AI audit findings ditemukan | 5 (input validation, auth, security headers, body limit, logging) |
| AI audit findings diperbaiki | 5/5 (100%) |
| Threat model findings | 18 (6 ✅ mitigasi, 4 ⚠️ partial, 8 ❌ belum) |
| Pipeline jobs | 4 (semua paralel) |
| Pipeline caching | 4 layers (Go modules, Gitleaks, Trivy DB, Semgrep) |
| Unit tests | 24 (semua PASS) |
| Go version | 1.26.0 |
| External dependencies | 2 direct (golang-jwt v5.3.1, x/crypto v0.53.0) |
| Trivy CVE count (current) | 0 |

---

## 7. Lessons Learned

### 7.1 `.gitignore` Bukan Security Boundary

`.gitignore` mencegah accidental commit, tapi `git add -f` bisa bypass. Gitleaks default menghormati `.gitignore` — jadi file yang di-force-add tetap di-skip. **Flag `--no-gitignore` wajib di CI pipeline** agar semua file yang tracked di git di-scan, regardless of `.gitignore`.

### 7.2 Scanner Otomatis Tidak Cukup

Trivy menemukan 4 CVE. Semgrep menemukan 2 findings. Tapi AI audit di Day 11 menemukan 5 kerentanan yang keduanya lewati: missing input validation, missing auth, missing security headers, missing body size limit, missing audit logging. **Business logic gaps tidak bisa dideteksi oleh pattern matching.**

### 7.3 Pipeline Caching Itu Penting

Dari Day 12, kita belajar bahwa pipeline yang tanpa cache menghabiskan ~30-50 detik per run hanya untuk download. Dengan 4 layer caching (Go modules, Gitleaks, Trivy DB, Semgrep), estimasi saving 30-40% per run.

### 7.4 JWT Authentication ≠ Authorization

Day 11 implementasi JWT auth (authentication — "siapa kamu?"), tapi threat modeling di Day 13 menunjukkan bahwa authorization ("apa yang boleh kamu lakukan?") masih belum ada. Siapa saja dengan JWT valid bisa lihat saldo akun siapa saja dan transfer dari akun siapa saja.

### 7.5 Threat Modeling Mengisi Gap Antara Scanner dan Manual Review

STRIDE + DREAD memberikan framework systematic untuk mengidentifikasi ancaman. Dari 18 threats yang ditemukan, 6 sudah di-mitigasi, 4 partial, dan 8 belum. Tanpa threat modeling, 8 threat ini mungkin tidak teridentifikasi sampai terlambat.

### 7.6 `math.IsInf` dan `math.IsNaN` Penting untuk Financial API

Input validation untuk financial endpoint tidak cukup hanya cek `amount > 0`. `NaN` dan `Inf` adalah float64 values yang valid di JSON tapi bisa menyebabkan behavior yang tidak diinginkan di arithmetic operations.

### 7.7 Intentional Vulnerability Test Memvalidasi Defense Layer

Tanpa test di Day 14, kita tidak akan tahu bahwa Gitleaks CI tidak menangkap file yang di-gitignore. Test ini membuktikan bahwa false sense of security lebih berbahaya daripada tidak punya security sama sekali.

---

## 8. Recommendations for Fase 2

| # | Rekomendasi | Prioritas | DREAD Score | Target |
|---|-------------|-----------|-------------|--------|
| 1 | **Rate limiter middleware** (100 req/min/IP) | 🔴 Critical | 9.6 | Fase 2 |
| 2 | **User-scoped authorization** (JWT sub claim) | 🔴 Critical | 8.4 | Fase 2 |
| 3 | **Persistent audit logging** (structured JSON to file) | 🟠 High | 6.0 | Fase 2 |
| 4 | **RBAC** (admin/user roles) | 🟠 High | 6.6 | Fase 2 |
| 5 | **JWT_SECRET mandatory** (fail if empty) | 🟠 High | 6.2 | Fase 2 |
| 6 | **Connection limits** (http.Server tuning) | 🟠 High | 7.2 | Fase 2 |
| 7 | **Container security** (Dockerfile hardening) | — | — | Fase 2 |
| 8 | **DAST** (OWASP ZAP baseline scan) | — | — | Fase 2 |

---

## 9. Final File Structure

```
securebank-api/
├── cmd/api/
│   ├── main.go                      # Entry point: 3 endpoints + middleware chain
│   └── main_test.go                 # 14 handler + validation + auth tests
├── configs/
│   ├── config.go                    # Config struct: PORT, DB_HOST, DB_PASSWORD, JWT_SECRET
│   └── config_test.go               # 2 config tests
├── internal/middleware/
│   ├── auth.go                      # RequireAuth: JWT Bearer token validation
│   ├── security.go                  # SecurityHeaders + LimitBodySize middleware
│   └── security_test.go             # 8 middleware tests
├── pkg/crypto/
│   ├── hash.go                      # bcrypt HashPassword + CheckPassword
│   ├── jwtutil.go                   # GenerateToken + ParseToken (golang-jwt/jwt/v5)
│   └── jwtutil_test.go              # 5 JWT + bcrypt tests
├── security/
│   ├── gitleaks-report.json         # Gitleaks scan report (Day 03)
│   ├── semgrep-report.json           # Semgrep scan report (Day 08)
│   ├── semgrep-latest.json           # Semgrep latest report (Day 11)
│   ├── semgrep-post-fix.json         # Semgrep post-fix report (Day 11)
│   ├── trivy-fs-report.json          # Trivy FS scan report (Day 05)
│   └── threat-model/
│       └── architecture.md           # STRIDE + DREAD threat model (Day 13)
├── go.mod                            # Module: github.com/stayrelevantid/securebank-api
├── go.sum
└── .gitignore

Root repo files:
├── .github/workflows/ci.yml          # CI pipeline: 4 parallel jobs + caching
├── .gitleaks.toml                    # Gitleaks allowlist config
├── .semgrep.yml                      # Custom Semgrep rule: no-md5-usage
├── blogpost.md                       # Blog post drafts (gitignored)
├── docs/
│   ├── fase-1-appsec.md              # Tutorial + retrospective documentation
│   ├── fase-2-infra-container.md     # Phase 2 tutorial
│   ├── fase-3-k8s-runtime.md        # Phase 3 tutorial
│   ├── fase-4-vuln-redteam.md        # Phase 4 tutorial
│   └── istilah-asing.md             # Glossary
├── progress/
│   ├── README.md                     # Progress overview
│   ├── tracker.md                    # Master tracker
│   └── daily/
│       ├── hari-01.md through hari-14.md
│       └── hari-15.md
├── sylabus.md                        # 60-day curriculum
└── README.md                         # Project overview
```

---

## 10. Pipeline Execution History

| Day | Event | Pipeline Status |
|-----|-------|-----------------|
| 02 | Basic CI (build + test) | ✅ Green |
| 04 | Added Gitleaks job | ✅ Green (after 3 fixes) |
| 06 | Added Trivy SCA job | ❌ Red (4 CVE found) |
| 07 | SCA remediation | ✅ Green (0 CVE) |
| 09 | Added Semgrep SAST job | ❌ Red (2 findings) |
| 10 | SAST remediation (MD5→bcrypt) | ✅ Green (1 WARNING accepted) |
| 11 | AI audit security fixes | ✅ Green (24 tests pass) |
| 12 | Pipeline optimization (caching) | ✅ Green (4 cache layers) |
| 14 | Gitleaks `--no-gitignore` fix | ✅ Green |

---

*Retrospektif ini ditulis pada Hari 15 sebagai penutup Fase 1.*
