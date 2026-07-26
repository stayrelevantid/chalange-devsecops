# ✅ Master Tracker — 60 Days DevSecOps Challenge

> Centang `[x]` setelah menyelesaikan setiap hari. Update kolom Status dan tanggal selesai.

---

## 🔐 Fase 1: Secure SDLC & Application Security (Hari 1–15)

| Hari | Topik | Output Konkret | Status | Selesai |
|------|-------|----------------|--------|---------|
| [Hari 01](daily/hari-01.md) | Setup Repo & Golang API | `main.go` dengan 3 endpoint + unit test | ✅ | 2026-06-05 |
| [Hari 02](daily/hari-02.md) | Pipeline CI/CD Dasar | `ci.yml` — auto build & test on push | ✅ | 2026-06-06 |
| [Hari 03](daily/hari-03.md) | Secret Scanning (Gitleaks) | Gitleaks lokal + `.gitleaks.toml` dikonfigurasi | ✅ | 2026-06-07 |
| [Hari 04](daily/hari-04.md) | Secret Scan di Pipeline + Remediasi | Job Gitleaks di CI, kredensial ke env vars | ✅ | 2026-06-08 |
| [Hari 05](daily/hari-05.md) | SCA Setup (Trivy FS) | `trivy fs .` menemukan CVE di `go.mod` | ✅ | 2026-06-09 |
| [Hari 06](daily/hari-06.md) | SCA di Pipeline (Gate) | Job Trivy gagalkan build jika CVE CRITICAL | ✅ | 2026-06-10 |
| [Hari 07](daily/hari-07.md) | SCA Remediation | `go get -u` patch dependensi, pipeline hijau | ✅ | 2026-06-11 |
| [Hari 08](daily/hari-08.md) | SAST Setup (Semgrep) | Semgrep lokal menemukan insecure crypto & SQL injection | ✅ | 2026-06-12 |
| [Hari 09](daily/hari-09.md) | SAST di Pipeline (Gate) | Job Semgrep blokir PR jika severity HIGH | ✅ | 2026-06-13 |
| [Hari 10](daily/hari-10.md) | SAST Remediation | Fix md5→bcrypt, pipeline hijau | ✅ | 2026-06-14 |
| [Hari 11](daily/hari-11.md) | AI-Assisted Code Audit | Output Trivy/Semgrep → AI → perbaikan diterapkan | ✅ | 2026-06-15 |
| [Hari 12](daily/hari-12.md) | Pipeline Optimization | Cache Go modules, parallel jobs, total < 2 menit | ✅ | 2026-06-15 |
| [Hari 13](daily/hari-13.md) | Threat Modeling (STRIDE) | Diagram arsitektur + tabel ancaman di `threat-model/` | ✅ | 2026-06-15 |
| [Hari 14](daily/hari-14.md) | Intentional Vuln Test | Commit fake secret → Gitleaks gagalkan build → revert | ✅ | 2026-06-18 |
| [Hari 15](daily/hari-15.md) | Dokumentasi Fase 1 | `fase-1-appsec.md` lengkap dengan cara kerja SAST, SCA, Secret Scan | ✅ | 2026-06-18 |

**Progres Fase 1: 15/15** ✅

---

## 🐳 Fase 2: Infrastructure as Code & Container Security (Hari 16–30)

| Hari | Topik | Output Konkret | Status | Selesai |
|------|-------|----------------|--------|---------|
| [Hari 16](daily/hari-16.md) | Dockerfile Multi-stage | Build di `golang:alpine`, run di `distroless` | ✅ | 2026-06-20 |
| [Hari 17](daily/hari-17.md) | Container Image Scan | `trivy image securebank:v1` — daftar CVE base image | ✅ | 2026-06-21 |
| [Hari 18](daily/hari-18.md) | Dockerfile Hardening | `USER nonroot`, `COPY --chown`, read-only filesystem | ✅ | 2026-06-22 |
| [Hari 19](daily/hari-19.md) | Image Signing (Cosign) | Key pair + signed image di registry | ✅ | 2026-06-23 |
| [Hari 20](daily/hari-20.md) | Terraform Setup + IaC Scan (Checkov) | `terraform/main.tf` (VPC+S3) + `checkov -d .` | ✅ | 2026-06-24 |
| [Hari 21](daily/hari-21.md) | IaC Scan (tfsec/Trivy) | Scan ulang dengan tfsec, bandingkan temuan | ✅ | 2026-06-25 |
| [Hari 22](daily/hari-22.md) | IaC di Pipeline | Job Checkov gagalkan build jika SG terbuka 0.0.0.0/0 | ✅ | 2026-06-26 |
| [Hari 23](daily/hari-23.md) | IaC Remediation | Enkripsi S3, perketat SG rules, pipeline hijau | ✅ | 2026-06-28 |
| [Hari 24](daily/hari-24.md) | DAST Setup (OWASP ZAP) | ZAP Baseline Scan ke `localhost:8080` | ✅ | 2026-07-02 |
| [Hari 25](daily/hari-25.md) | DAST di Pipeline | Job ZAP scan staging setelah deploy | ✅ | 2026-07-02 |
| [Hari 26](daily/hari-26.md) | DAST Remediation | Tambah security headers di Go middleware | ✅ | 2026-07-03 |
| [Hari 27](daily/hari-27.md) | AI-Assisted IaC Fix | Terraform rentan → AI → versi aman diterapkan | ✅ | 2026-07-04 |
| [Hari 28](daily/hari-28.md) | Compliance as Code (InSpec) | Profile InSpec: verifikasi port 22 tertutup | ✅ | 2026-07-05 |
| [Hari 29](daily/hari-29.md) | Pipeline Consolidation | Satu YAML: Secret+SAST+SCA+Image+IaC scan paralel | ✅ | 2026-07-06 |
| [Hari 30](daily/hari-30.md) | Dokumentasi Fase 2 | `fase-2-infra-container.md` lengkap | ✅ | 2026-07-06 |

**Progres Fase 2: 15/15** ✅

---

## ☸️ Fase 3: Kubernetes & Runtime Security (Hari 31–45)

| Hari | Topik | Output Konkret | Status | Selesai |
|------|-------|----------------|--------|---------|
| [Hari 31](daily/hari-31.md) | K8s Cluster + Deploy | k3d cluster + `kubectl apply` deployment & service | ✅ | 2026-07-08 |
| [Hari 32](daily/hari-32.md) | K8s Misconfiguration Scan | Kubesec/Checkov scan `deployment.yaml` | ✅ | 2026-07-09 |
| [Hari 33](daily/hari-33.md) | SecurityContext Hardening | `readOnlyRootFilesystem`, `allowPrivilegeEscalation: false` | ✅ | 2026-07-10 |
| [Hari 34](daily/hari-34.md) | OPA Gatekeeper Setup | Helm install Gatekeeper di cluster | ✅ | 2026-07-11 |
| [Hari 35](daily/hari-35.md) | Rego Policy Writing | Policy: tolak Pod tanpa resource limits + requests | ✅ | 2026-07-11 |
| [Hari 36](daily/hari-36.md) | OPA Policy Testing | Deploy Pod tanpa limits → "Denied by Gatekeeper" | ✅ | 2026-07-13 |
| [Hari 37](daily/hari-37.md) | Network Policies | Default Deny All + whitelist Ingress traffic | ✅ | 2026-07-14 |
| [Hari 38](daily/hari-38.md) | RBAC Auditing | `kubectl who-can` audit, hapus overprivileged SA | ✅ | 2026-07-15 |
| [Hari 39](daily/hari-39.md) | Falco Setup | Helm install Falco, verifikasi logs berjalan | ✅ | 2026-07-16 |
| [Hari 40](daily/hari-40.md) | Falco Custom Rules | Rule: alert jika `bash`/`sh` dijalankan di container | ✅ | 2026-07-18 |
| [Hari 41](daily/hari-41.md) | Falco Attack Simulation | `kubectl exec` → cek Falco log Notice/Warning | ✅ | 2026-07-18 |
| [Hari 42](daily/hari-42.md) | Alerting Webhook (n8n) | n8n workflow: webhook trigger untuk alert Falco | ✅ | 2026-07-19 |
| [Hari 43](daily/hari-43.md) | K8s Secret Management | External Secrets Operator → AWS Secrets Manager | ✅ | 2026-07-20 |
| [Hari 44](daily/hari-44.md) | AI Threat Modeling K8s | Topologi + RBAC → AI analisis 3 attack path | ✅ | 2026-07-21 |
| [Hari 45](daily/hari-45.md) | Dokumentasi Fase 3 | `fase-3-k8s-runtime.md` lengkap | ✅ | 2026-07-22 |

**Progres Fase 3: 15/15** ✅ SELESAI

---

## 🔴 Fase 4: Vulnerability Management & Red Teaming (Hari 46–60)

| Hari | Topik | Output Konkret | Status | Selesai |
|------|-------|----------------|--------|---------|
| [Hari 46](daily/hari-46.md) | DefectDojo Setup | Docker Compose up, login dashboard | ✅ | 2026-07-22 |
| [Hari 47](daily/hari-47.md) | DefectDojo API Integration | Pipeline upload Trivy+Semgrep JSON ke DefectDojo | ✅ | 2026-07-23 |
| [Hari 48](daily/hari-48.md) | Alert Routing (Slack) | CRITICAL alert → Slack `#security-alerts` via webhook receiver | ✅ | 2026-07-26 |
| [Hari 49](daily/hari-49.md) | AI Remediation Node | n8n + LLM: auto-ringkas SAST finding → Slack dev | ⬜ | — |
| [Hari 50](daily/hari-50.md) | CSPM (Prowler/ScoutSuite) | Scan AWS/GCP sandbox vs CIS Benchmarks | ⬜ | — |
| [Hari 51](daily/hari-51.md) | CSPM Remediation | Fix 3 temuan kritis (MFA, S3 bucket, stale keys) | ⬜ | — |
| [Hari 52](daily/hari-52.md) | Red Team: K8s Escape | Privileged pod → host filesystem → Falco/OPA deteksi | ⬜ | — |
| [Hari 53](daily/hari-53.md) | Red Team: Leaked Creds | IAM user sementara → CloudTrail → auto-revoke script | ⬜ | — |
| [Hari 54](daily/hari-54.md) | Chaos Security Engineering | Kill OPA webhook → deploy vulnerable app → alarm? | ⬜ | — |
| [Hari 55](daily/hari-55.md) | Laporan Audit | Export DefectDojo → Executive Summary draft | ⬜ | — |
| [Hari 56](daily/hari-56.md) | Dokumen Eksekutif (PDF) | PDF formal: metodologi, temuan, mitigasi, sisa risiko | ⬜ | — |
| [Hari 57](daily/hari-57.md) | AI Review Dokumen | AI periksa struktur & nada bahasa laporan audit | ⬜ | — |
| [Hari 58](daily/hari-58.md) | CDP Exam Sim: Lab Setup | Wipe semua config, siapkan app mentah | ⬜ | — |
| [Hari 59](daily/hari-59.md) | CDP Exam Sim: Execution | 3 jam: rebuild semua pipeline dari nol | ⬜ | — |
| [Hari 60](daily/hari-60.md) | Project Showcase | Diagram arsitektur E2E + publikasi final | ⬜ | — |

**Progres Fase 4: 3/15**

---

## 📊 Summary

| Metrik | Nilai |
|--------|-------|
| Total Hari Selesai | 48 / 60 |
| Fase Selesai | 3 / 4 |
| Hari Aktif (ada catatan) | 47 |
| Streak Hari Berturut-turut | 0 |

---

*Legend: ⬜ Belum dimulai · 🔄 Sedang dikerjakan · ✅ Selesai · ⏭️ Dilewati*
