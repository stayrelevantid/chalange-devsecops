# Draft Laporan Audit Q3 — SecureBank API

**Status:** Draf (Hari 55) — akan diperluas jadi dokumen eksekutif (PDF) di Hari 56
**Sumber data:** DefectDojo (engagement *Q3 Security Audit*), re-sync scan 2026-08-06
**Metodologi:** SAST (Semgrep), SCA (Trivy), Container Image (Trivy), DAST (OWASP ZAP), IaC (Trivy config, Checkov), K8s (Trivy config, Checkov), CSPM (Prowler)

---

## Executive Summary

Aplikasi SecureBank API telah diuji menggunakan metodologi DevSecOps menyeluruh: SAST, SCA, DAST, Image Scan, IaC Scan, dan K8s Misconfiguration Scan. Semua hasil diagregasi di DefectDojo untuk memudahkan pemantauan.

**Snapshot saat ini (re-sync scan terbaru):** 6 temuan tersisa — 1 High, 2 Medium, 2 Low, 1 Info. Tidak ada temuan **Critical** pada state terkini.

### Metrik Ringkas

| Metrik | Nilai |
|--------|-------|
| Total temuan (aggregat DefectDojo) | 52 |
| Temuan state terkini (re-sync 06-08) | 6 |
| Critical (state terkini) | 0 |
| High (state terkini) | 1 |
| Medium (state terkini) | 2 |
| Low (state terkini) | 2 |
| Info (state terkini) | 1 |

### Tren Remediasi

| Area | Sebelum (baseline) | Sekarang (state terkini) |
|------|--------------------|--------------------------|
| IaC — S3 & Security Group | 1 Critical + 6 High (AWS-0086/0087/0091/0093/0107, dll.) | 2 Low (AWS-0089 S3 Bucket Logging) |
| K8s manifest (`deployment.yaml`) | 16 misconfig (KSV-0001, KSV-0012, ...) | 0 |
| Dependensi Go (`go.mod`) | 4 CVE (1 CRITICAL, 3 HIGH) — Day 05 | 0 Critical/High (1 Info: openpgp) |
| Container image | 0 (baseline distroless) | 1 High + 1 Medium (stdlib v1.26.4) |
| SAST (Semgrep) | 3 temuan (MD5, SQLi, ...) | 1 Medium (HTTP tanpa TLS) |

---

## Key Findings (State Terkini)

1. **[HIGH] CVE-2026-39822 — Stdlib v1.26.4** (image `securebank:v1`)
   Image di-build dengan Go 1.26.4; ada kerentanan stdlib yang sudah diperbaiki di 1.26.5. Aksi: rebuild image dengan Go terbaru.
2. **[MEDIUM] CVE-2026-42505 — Stdlib v1.26.4** (image) — ikut ter-remediasi saat rebuild image.
3. **[MEDIUM] HTTP tanpa TLS** (`cmd/api/main.go:153`) — API berjalan plain HTTP. Aksi: pasang `ListenAndServeTLS` + sertifikat.
4. **[LOW] AWS-0089 — S3 Bucket Logging** (2 bucket) — aktifkan bucket logging untuk audit trail akses S3.
5. **[INFO] GO-2026-5932 — `golang.org/x/crypto` openpgp** — paket openpgp di x/crypto tidak di-maintain; tidak dipakai langsung oleh aplikasi, tapi ter-scan. Aksi: singkirkan dependency yang menyeretnya bila memungkinkan.

---

## Remediation Evidence (Telah Dilakukan)

- **S3 Block Public Access + SG tightening** (Day 51): hasil Prowler turun 132 → 106 FAIL; scan ulang trivy config hanya menyisakan 2 Low.
- **K8s hardening** (Day 33): `readOnlyRootFilesystem`, `allowPrivilegeEscalation: false`, drop capabilities → trivy-k8s-post-fix **0 findings**, checkov-k8s-post-fix **0 failed**.
- **Crypto & SQL injection** (Day 10): MD5 → bcrypt, query parameterized → semgrep tersisa hanya TLS.
- **Dependensi** (Day 07): `go get -u` patch CVE CRITICAL/HIGH di go.mod.
- **Defense in depth** (Day 52–54): gap OPA ditutup (deny-privileged), Falco 7 rules + alerting ke Slack, drift-check.sh.

---

## Residual Risk

| Risiko Tersisa | Tingkat | Catatan |
|----------------|---------|---------|
| Image stdlib CVE (Go 1.26.4) | HIGH | Rebuild image wajib sebelum rilis production |
| API tanpa TLS | MEDIUM | Butuh sertifikat + terminasi TLS |
| S3 bucket logging mati | LOW | Perlu enable logging + integrasi CloudTrail |
| openpgp dependency | INFO | Deprecated, tidak dipakai langsung |
| Red Team artifacts di repo (attacker/chaos manifests) | — | Manifest demo jangan pernah di-deploy ke production |

---

## Metodologi & Tools

| Layer | Tool | Cakupan |
|-------|------|---------|
| SAST | Semgrep (`p/golang`) | Kode Go (auth, crypto, SQL) |
| SCA | Trivy FS | Dependensi Go (go.mod) |
| Container | Trivy Image | Base image + binary dependencies |
| DAST | OWASP ZAP | Endpoint HTTP di `localhost:8080` |
| IaC | Trivy config, Checkov | Terraform (S3, SG, VPC) |
| K8s | Trivy config (KSV), Checkov | Manifest deployment/service |
| CSPM | Prowler | AWS CIS Benchmark (328 checks) |
| Runtime | Falco | syscall events + alerting |

---

*Draf ini akan dilengkapi menjadi dokumen eksekutif final (PDF) di Hari 56.*
