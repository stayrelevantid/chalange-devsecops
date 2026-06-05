# Fase 2: Infrastructure as Code & Container Security (Hari 16–30)

> **Proyek:** SecureBank API — Kontainerisasi, hardening Docker, Terraform infra, dan DAST
> **Output Fase:** Image Docker yang hardened & signed, Terraform yang aman, pipeline terpadu dengan semua scan paralel.

---

## Hari 16: Dockerfile Multi-stage Build

### Tujuan
Membuat Docker image SecureBank API sekecil dan seaman mungkin menggunakan multi-stage build.

### Tutorial

**1. Buat `Dockerfile`:**
```dockerfile
# === Stage 1: Build ===
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /securebank ./cmd/api

# === Stage 2: Runtime ===
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /securebank /securebank
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/securebank"]
```

**2. Buat `.dockerignore`:**
```
.git
.github
docs
security
terraform
k8s
*.md
```

**3. Build dan test:**
```bash
docker build -t securebank:v1 .
docker run -p 8080:8080 securebank:v1
curl http://localhost:8080/health
```

**4. Cek ukuran image:**
```bash
docker images securebank:v1
# Target: < 15MB karena distroless + stripped binary
```

### Checklist
- [ ] Multi-stage build: golang:alpine → distroless
- [ ] Binary di-strip (`-ldflags="-w -s"`)
- [ ] Image < 15MB
- [ ] Aplikasi berjalan normal di container

---

## Hari 17: Container Image Scanning

### Tujuan
Memindai Docker image untuk CVE pada base image dan layer dependencies.

### Tutorial

**1. Scan image dengan Trivy:**
```bash
trivy image securebank:v1
trivy image --severity HIGH,CRITICAL securebank:v1
```

**2. Bandingkan dengan image non-distroless:**
```bash
# Build versi "naif" untuk perbandingan
docker build -t securebank:naive --target builder .
trivy image securebank:naive
# Jauh lebih banyak CVE karena alpine + go toolchain
```

**3. Export report:**
```bash
trivy image securebank:v1 --format json --output security/trivy-image-report.json
```

### Checklist
- [ ] Scan image distroless: minimal/zero CVE
- [ ] Scan image builder: banyak CVE (bukti pentingnya multi-stage)
- [ ] Report JSON disimpan di `security/`
- [ ] Memahami perbedaan OS-level vs app-level CVE

---

## Hari 18: Dockerfile Hardening

### Tujuan
Menerapkan best practice keamanan container: non-root, read-only, no new privileges.

### Tutorial

**1. Verifikasi Dockerfile sudah memiliki:**
```dockerfile
USER nonroot:nonroot  # ✅ sudah dari hari 16
```

**2. Tambahkan health check di Dockerfile:**
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD ["/securebank", "-health-check"] || exit 1
```

**3. Buat `docker-compose.yml` dengan security options:**
```yaml
version: '3.8'

services:
  securebank:
    build: .
    image: securebank:v1
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=${DB_PASSWORD}
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=10m
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 128M
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: securebank
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    security_opt:
      - no-new-privileges:true

volumes:
  pgdata:
```

**4. Test:**
```bash
DB_PASSWORD=localdev123 docker compose up -d
curl http://localhost:8080/health
```

### Checklist
- [ ] Container berjalan sebagai non-root
- [ ] `read_only: true` + tmpfs untuk /tmp
- [ ] `no-new-privileges` aktif
- [ ] Resource limits (CPU + Memory) ditentukan
- [ ] Health check berjalan

---

## Hari 19: Image Signing dengan Cosign

### Tujuan
Menandatangani image Docker untuk menjamin integritas dan provenance.

### Tutorial

**1. Install Cosign:**
```bash
brew install cosign
```

**2. Generate key pair:**
```bash
cosign generate-key-pair
# Menghasilkan cosign.key (private) dan cosign.pub (public)
```

**3. Push image ke registry (contoh: GitHub Container Registry):**
```bash
docker tag securebank:v1 ghcr.io/<username>/securebank:v1
docker push ghcr.io/<username>/securebank:v1
```

**4. Sign image:**
```bash
cosign sign --key cosign.key ghcr.io/<username>/securebank:v1
```

**5. Verifikasi signature:**
```bash
cosign verify --key cosign.pub ghcr.io/<username>/securebank:v1
```

**6. Tambahkan `cosign.key` ke `.gitignore`:**
```
cosign.key
```

### Checklist
- [ ] Key pair di-generate
- [ ] Image di-push ke registry
- [ ] Image ditandatangani dengan Cosign
- [ ] Signature bisa diverifikasi
- [ ] Private key TIDAK di-commit (ada di `.gitignore`)

---

## Hari 20: Terraform Setup + IaC Scanning (Checkov)

### Tujuan
Menulis Terraform dasar (VPC + S3) dan memindai misconfiguration dengan Checkov.

### Tutorial

**1. Buat `terraform/main.tf`:**
```hcl
terraform {
  required_version = ">= 1.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC — intentionally misconfigured for training
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name        = "securebank-vpc"
    Environment = var.environment
  }
}

# S3 Bucket — intentionally insecure for training
resource "aws_s3_bucket" "logs" {
  bucket = "securebank-logs-${var.environment}"
  # Missing: versioning, encryption, public access block
}

# Security Group — intentionally too open
resource "aws_security_group" "api" {
  name   = "securebank-api-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # INTENTIONAL: open to world
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

**2. Buat `terraform/variables.tf`:**
```hcl
variable "aws_region" {
  default = "ap-southeast-1"
}

variable "environment" {
  default = "dev"
}
```

**3. Install dan jalankan Checkov:**
```bash
pip3 install checkov
checkov -d terraform/
```

**4. Perhatikan temuan: S3 tanpa enkripsi, SG terbuka, dll.**

### Checklist
- [ ] Terraform files dibuat di `terraform/`
- [ ] Checkov menemukan ≥ 5 misconfiguration
- [ ] Output dipahami (check ID, resource, guideline)
- [ ] Belum diperbaiki (itu tugas hari 23)

---

## Hari 21: IaC Scanning dengan tfsec & Trivy

### Tujuan
Membandingkan multiple IaC scanner untuk coverage yang lebih baik.

### Tutorial

**1. Install tfsec:**
```bash
brew install tfsec
```

**2. Scan Terraform:**
```bash
tfsec terraform/
```

**3. Scan juga dengan Trivy (mode IaC):**
```bash
trivy config terraform/
```

**4. Bandingkan:**

Buat `security/iac-scan-comparison.md`:
```markdown
| Finding | Checkov | tfsec | Trivy |
|---------|---------|-------|-------|
| S3 no encryption | ✅ | ✅ | ✅ |
| SG open to world | ✅ | ✅ | ✅ |
| S3 no versioning | ✅ | ✅ | ❌ |
| ... | | | |
```

### Checklist
- [ ] tfsec dan Trivy config scan berjalan
- [ ] Perbandingan 3 tools didokumentasikan
- [ ] Memahami kelebihan masing-masing scanner

---

## Hari 22: IaC Scan di Pipeline

### Tujuan
Pipeline khusus infra yang menggagalkan build jika Terraform tidak aman.

### Tutorial

**Buat `.github/workflows/infra.yml`:**
```yaml
name: Infrastructure Security

on:
  push:
    paths:
      - 'terraform/**'
  pull_request:
    paths:
      - 'terraform/**'

jobs:
  iac-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Checkov
        uses: bridgecrewio/checkov-action@v12
        with:
          directory: terraform/
          framework: terraform
          soft_fail: false  # gagalkan pipeline
          output_format: sarif
          output_file_path: checkov.sarif

      - name: Run tfsec
        uses: aquasecurity/tfsec-action@v1.0.0
        with:
          working_directory: terraform/
          soft_fail: false

      - name: Upload SARIF
        if: always()
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: checkov.sarif
```

### Checklist
- [ ] Pipeline trigger hanya pada perubahan `terraform/`
- [ ] Checkov + tfsec berjalan
- [ ] Pipeline gagal karena SG terbuka 0.0.0.0/0
- [ ] SARIF ter-upload ke GitHub Security

---

## Hari 23: IaC Remediation

### Tujuan
Memperbaiki semua misconfiguration Terraform hingga pipeline hijau.

### Tutorial

**Update `terraform/main.tf`:**
```hcl
# S3 Bucket — FIXED
resource "aws_s3_bucket" "logs" {
  bucket = "securebank-logs-${var.environment}"
}

resource "aws_s3_bucket_versioning" "logs" {
  bucket = aws_s3_bucket.logs.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "logs" {
  bucket                  = aws_s3_bucket.logs.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Security Group — FIXED
resource "aws_security_group" "api" {
  name   = "securebank-api-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  ingress {
    description = "API port from VPC"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS outbound only"
  }
}
```

**Verifikasi:**
```bash
checkov -d terraform/
tfsec terraform/
# Semua check PASSED
```

### Checklist
- [ ] S3: versioning + encryption + public access block
- [ ] SG: hanya port 443 & 8080, hanya dari VPC CIDR
- [ ] Checkov + tfsec: 0 FAILED
- [ ] Pipeline infra hijau

---

## Hari 24: DAST Setup — OWASP ZAP

### Tujuan
Menyerang aplikasi SecureBank API yang sedang berjalan untuk menemukan kerentanan runtime.

### Tutorial

**1. Jalankan SecureBank lokal:**
```bash
docker compose up -d
```

**2. Jalankan ZAP Baseline Scan:**
```bash
docker run --rm --network host \
  ghcr.io/zaproxy/zaproxy:stable \
  zap-baseline.py \
  -t http://localhost:8080 \
  -r zap-report.html \
  -J zap-report.json
```

**3. Buka `zap-report.html` di browser.**

**4. Catat temuan umum:**
- Missing `Content-Security-Policy` header
- Missing `X-Content-Type-Options` header
- Missing `Strict-Transport-Security` header

### Checklist
- [ ] ZAP berjalan dan memindai aplikasi
- [ ] Report HTML dan JSON dihasilkan
- [ ] Temuan security headers dicatat
- [ ] Report disimpan di `security/`

---

## Hari 25: DAST Pipeline Integration

### Tujuan
ZAP berjalan otomatis di pipeline setelah deploy ke staging.

### Tutorial

**Tambahkan job ke `.github/workflows/ci.yml`:**
```yaml
  dast-scan:
    needs: [build-and-test]
    runs-on: ubuntu-latest
    services:
      securebank:
        image: securebank:v1
        ports:
          - 8080:8080
    steps:
      - name: Wait for app
        run: |
          for i in $(seq 1 30); do
            curl -sf http://localhost:8080/health && break
            sleep 2
          done

      - name: ZAP Baseline Scan
        uses: zaproxy/action-baseline@v0.12.0
        with:
          target: 'http://localhost:8080'
          allow_issue_writing: false

      - name: Upload ZAP Report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: zap-report
          path: report_html.html
```

### Checklist
- [ ] ZAP berjalan setelah app di-deploy di CI
- [ ] Report ter-upload sebagai artifact
- [ ] Temuan tercatat

---

## Hari 26: DAST Remediation — Security Headers

### Tujuan
Menambahkan security headers di middleware Golang untuk memperbaiki temuan ZAP.

### Tutorial

**Buat `internal/middleware/security.go`:**
```go
package middleware

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
        w.Header().Set("X-XSS-Protection", "0")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
        next.ServeHTTP(w, r)
    })
}
```

**Update `cmd/api/main.go` untuk menggunakan middleware:**
```go
mux := http.NewServeMux()
mux.HandleFunc("/health", healthCheck)
mux.HandleFunc("/balance", getBalance)
mux.HandleFunc("/transfer", transfer)

handler := middleware.SecurityHeaders(mux)
log.Fatal(http.ListenAndServe(":8080", handler))
```

**Verifikasi:**
```bash
curl -I http://localhost:8080/health
# Harus muncul semua security headers
```

### Checklist
- [ ] Middleware security headers dibuat
- [ ] Semua 7 headers ditambahkan
- [ ] ZAP re-scan: temuan header berkurang
- [ ] Unit test untuk middleware

---

## Hari 27: AI-Assisted IaC Automation

### Tujuan
Menggunakan AI untuk memperbaiki Terraform yang kompleks sesuai standar industri.

### Tutorial

**1. Ambil Terraform block yang masih belum optimal.**

**2. Prompt ke AI:**
```
Kamu adalah AWS Solutions Architect. Perbaiki Terraform berikut agar memenuhi:
- CIS AWS Foundations Benchmark v2.0
- AWS Well-Architected Framework (Security Pillar)
- Prinsip least privilege

[paste terraform code]
```

**3. Terapkan perbaikan, verifikasi dengan Checkov.**

### Checklist
- [ ] AI memberikan perbaikan berbasis standar industri
- [ ] Perbaikan diterapkan dan diverifikasi
- [ ] Checkov pass 100%

---

## Hari 28: Compliance as Code — Chef InSpec

### Tujuan
Menulis test keamanan infrastruktur yang bisa dijalankan berulang.

### Tutorial

**1. Install InSpec:**
```bash
brew install --cask chef/chef/inspec
```

**2. Buat profile `security/inspec-profiles/securebank/`:**
```bash
inspec init profile security/inspec-profiles/securebank
```

**3. Tulis test di `controls/network.rb`:**
```ruby
control 'ssh-port-closed' do
  impact 1.0
  title 'SSH port should be closed'
  desc 'Port 22 should not be open to public internet'

  describe port(22) do
    it { should_not be_listening }
  end
end

control 'api-port-listening' do
  impact 0.7
  title 'API port should be listening'

  describe port(8080) do
    it { should be_listening }
  end
end

control 'no-root-process' do
  impact 1.0
  title 'Application should not run as root'

  describe processes('securebank') do
    its('users') { should_not include 'root' }
  end
end
```

**4. Jalankan:**
```bash
inspec exec security/inspec-profiles/securebank/
```

### Checklist
- [ ] InSpec profile dibuat
- [ ] 3 controls: SSH closed, API listening, non-root
- [ ] Test berjalan dan hasilnya dipahami

---

## Hari 29: Pipeline Consolidation — Full Security Pipeline

### Tujuan
Menggabungkan SEMUA scan ke satu pipeline yang berjalan paralel.

### Tutorial

**Update `.github/workflows/security-scan.yml`:**
```yaml
name: Full Security Scan

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  secret-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: gitleaks/gitleaks-action@v2
        env: { GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}" }

  sca-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: fs
          scan-ref: '.'
          exit-code: '1'
          severity: CRITICAL

  sast-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: semgrep/semgrep-action@v1
        with:
          config: "p/golang p/owasp-top-ten"

  image-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: docker build -t securebank:ci .
      - uses: aquasecurity/trivy-action@master
        with:
          image-ref: securebank:ci
          exit-code: '1'
          severity: CRITICAL

  iac-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: bridgecrewio/checkov-action@v12
        with:
          directory: terraform/
          soft_fail: false

  # Gate — semua scan harus lulus
  security-gate:
    needs: [secret-scan, sca-scan, sast-scan, image-scan, iac-scan]
    runs-on: ubuntu-latest
    steps:
      - run: echo "✅ All security scans passed"
```

### Checklist
- [ ] 5 scan jobs berjalan paralel
- [ ] `security-gate` job sebagai final check
- [ ] Total waktu < 3 menit (paralel)
- [ ] Pipeline hijau = kode & infra aman

---

## Hari 30: Dokumentasi Fase 2

### Tujuan
Merangkum Container Security dan IaC Security.

### Tutorial
Tulis `docs/fase-2-infra-container.md`:
- Perbandingan image size: naive vs distroless
- Checklist hardening Dockerfile (12 item)
- Perbandingan IaC scanner (Checkov vs tfsec vs Trivy)
- DAST findings dan remediasi
- Screenshot pipeline paralel

### Checklist
- [ ] Dokumen lengkap dengan metrics
- [ ] Commit: `docs: fase 2 infrastructure and container security`

---

> ✅ **Selesai Fase 2** — Lanjut ke [Fase 3: Kubernetes & Runtime Security](fase-3-k8s-runtime.md)
