# Hari 47 — DefectDojo API Integration

**📅 Tanggal:** 2026-07-23  
**⏱️ Durasi Belajar:** ~50 menit  
**🏷️ Fase:** Fase 4 — Vuln Mgmt & Red Team  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Create upload-scans.sh — automated upload semua scanner reports ke DefectDojo
- [x] Test upload-scans.sh — 7 reports successful, 15 tests imported, 46 findings
- [x] Add CI job `upload-defectdojo` ke ci.yml (conditional skip kalau DEFECTDOJO_URL secret tidak ada)
- [x] Fix Day 46 terminology mapping (Product Type → Organization, Product → Asset)
- [x] Document local vs CI approach

---

## ✅ Yang Berhasil Dikerjakan

- `upload-scans.sh` created + tested: 7 scanner reports uploaded, 15 tests imported, 46 findings
- CI job `upload-defectdojo` added ke ci.yml (job 9, conditional pada `vars.DEFECTDOJO_URL`)
- Day 46 daily note updated dengan terminology mapping (API v2 vs UI v3.x)
- ZAP report butuh XML format (JSON ditolak DefectDojo) — documented
- Product lifecycle "active" field invalid di API v2 — pakai default (omit field)

---

## 📝 Catatan Teknis

### upload-scans.sh

Script yang automate upload 7 scanner reports ke DefectDojo via `import-scan` API endpoint:

```bash
#!/usr/bin/env bash
# Usage: bash upload-scans.sh
# Env vars: DEFECTDOJO_URL, DEFECTDOJO_USER, DEFECTDOJO_PASS,
#           PRODUCT_NAME, ENGAGEMENT_NAME

# 1. Get API token via /api/v2/api-token-auth/
# 2. Upload 7 reports via /api/v2/import-scan/ (Trivy x4, Semgrep, Checkov x2)
# 3. Print summary: tests count, findings count
```

### Test Results (upload-scans.sh)

```
=== Uploading scanner reports ===
  OK: Trivy Scan -> test_id=9 (trivy-fs-report.json)
  OK: Trivy Scan -> test_id=10 (trivy-image-report.json)
  OK: Trivy Scan -> test_id=11 (trivy-iac-report.json)
  OK: Trivy Scan -> test_id=12 (trivy-k8s-post-fix-report.json)
  OK: Semgrep JSON Report -> test_id=13 (semgrep-report.json)
  OK: Checkov Scan -> test_id=14 (checkov-report.json)
  OK: Checkov Scan -> test_id=15 (checkov-k8s-post-fix-report.json)

=== Summary ===
Tests imported: 15
Findings: 46
```

### Scanner Compatibility dengan DefectDojo

| Scanner | Report Format | scan_type | Compatible | Findings Contributed |
|---------|---------------|-----------|------------|---------------------|
| Trivy FS | JSON | Trivy Scan | ✅ | 0 (post-fix, clean) |
| Trivy Image | JSON | Trivy Scan | ✅ | 0 (distroless, 0 CVE) |
| Trivy IaC | JSON | Trivy Scan | ✅ | 0 (post-fix, clean) |
| Trivy K8s | JSON | Trivy Scan | ✅ | 0 (post-fix, clean) |
| Semgrep | JSON | Semgrep JSON Report | ✅ | 1 (Medium: TLS audit) |
| Checkov IaC | JSON | Checkov Scan | ✅ | 0 (102 passed, 0 failed) |
| Checkov K8s | JSON | Checkov Scan | ✅ | 15 (Medium: SG/S3/KMS findings) |
| ZAP | JSON | ZAP Scan | ❌ | Wrong format — needs XML |

**Total: 7/8 scanner reports berhasil import, 16 findings centralized.**

### CI Job: upload-defectdojo

```yaml
  upload-defectdojo:
    name: DefectDojo Upload
    needs: [sca-scan, sast-scan, iac-checkov, iac-trivy]
    runs-on: ubuntu-latest
    if: ${{ vars.DEFECTDOJO_URL != '' }}  # Skip kalau secret tidak ada (local)

    steps:
      - Download artifacts (trivy-sca, semgrep, checkov)
      - Upload Trivy to DefectDojo (POST /api/v2/import-scan/)
      - Upload Semgrep to DefectDojo
      - Upload Checkov to DefectDojo
```

### Conditional Skip Logic

```yaml
if: ${{ vars.DEFECTDOJO_URL != '' }}
```

- **Local (current):** `DEFECTDOJO_URL` secret tidak di-set → job skip → local pakai `upload-scans.sh` manual
- **Production:** deploy DefectDojo di cloud → set `DEFECTDOJO_URL` + `DEFECTDOJO_API_KEY` sebagai GitHub repository variables/secrets → CI auto-upload

### CI vs Local Split

| Aspek | Local (current) | CI (production) |
|-------|-----------------|-----------------|
| DefectDojo URL | `http://localhost:8088` | `https://defectdojo.company.com` |
| Trigger | Manual `bash upload-scans.sh` | Automatic di setiap push/PR |
| Reports | All 7 local reports | Artifact dari CI jobs (trivy, semgrep, checkov) |
| Speed | Immediate (~3s) | After CI scan jobs complete (~5 min) |
| Networking | localhost → localhost | GitHub Actions → cloud DefectDojo |

### DefectDojo API: import-scan Endpoint

```bash
curl -X POST "http://localhost:8088/api/v2/import-scan/" \
  -H "Authorization: Token <API_TOKEN>" \
  -F "scan_type=Trivy Scan" \
  -F "product_name=SecureBank API" \
  -F "engagement_name=Q3 Security Audit" \
  -F "file=@trivy-sca.json"
```

**Key fields:**
- `scan_type`: DefectDojo supported scanner name (exact match required)
- `product_name`: Product/Asset name dari hierarchy
- `engagement_name`: Engagement name dari hierarchy
- `file`: scanner report file (multipart form upload)

### Terminology Mapping (Day 46 Fix)

| API v2 (curl commands) | UI v3.x (browser) | Our Value |
|-------------------------|---------------------|-----------|
| Product Type | Organization | `Fintech` |
| Product | Asset | `SecureBank API` |
| Engagement | Engagement (unchanged) | `Q3 Security Audit` |

Tutorial di `docs/fase-4-vuln-redteam.md` masih pakai terminology lama ("Product Type", "Product"). UI v3.x sudah rename ke "Organization" dan "Asset". API v2 tetap pakai nama lama — jadi curl commands work dengan `product_types/` dan `products/` endpoints.

### ZAP Import Failure

```
Error: ['Internal error: Wrong file format, please use xml.']
```

DefectDojo ZAP import butuh **XML format**, bukan JSON. ZAP scan reports kita di `security/zap-report.json` adalah JSON format. Fix options:
- Re-run ZAP dengan `--xml` output flag
- Atau skip ZAP import (1 false positive only — rule 10049)

Documented sebagai known limitation.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| DefectDojo local tidak reachable dari GitHub Actions CI | Conditional job `if: vars.DEFECTDOJO_URL != ''` — skip kalau local, run kalau production |
| ZAP import gagal — "Wrong file format, please use xml" | DefectDojo butuh XML for ZAP. JSON ditolak. Skip ZAP import (1 false positive only) |
| Product lifecycle "active" invalid choice | Hapus `lifecycle` field dari API payload — DefectDojo assign default |
| DefectDojo v3.x terminology berbeda dari tutorial | Product Type → Organization, Product → Asset. API v2 tetap pakai nama lama. Documented di Day 46 fix |
| upload-scans.sh path resolving | Pakai `$(dirname "$0")/../` untuk resolve relative path dari script location |

---

## 📤 Output Hari Ini

- [x] `security/defectdojo/upload-scans.sh` — automated upload script (7 reports, committed)
- [x] CI job `upload-defectdojo` added ke ci.yml (job 9, conditional skip)
- [x] Day 46 daily note updated dengan terminology mapping
- [x] 15 tests imported, 46 findings centralized di DefectDojo dashboard
- [x] 7/8 scanner reports successfully imported (ZAP needs XML)
- [x] Local vs CI approach documented

---

## 💡 Pelajaran Baru

- **Single pane of glass achieved.** 7 scanner reports dari 4 tools (Trivy, Semgrep, Checkov, ZAP-skip) centralized di 1 DefectDojo dashboard. 46 findings dari 2 scanners (Semgrep + Checkov K8s). Visibility > siloed reports.

- **DefectDojo import-scan API.** Endpoint `/api/v2/import-scan/` terima multipart form upload dengan `scan_type`, `product_name`, `engagement_name`, `file`. `scan_type` butuh exact match dengan DefectDojo supported scanner list.

- **Conditional CI jobs.** `if: ${{ vars.DEFECTDOJO_URL != '' }}` = job skip kalau variable tidak di-set. Berguna untuk local vs production split — local pakai manual script, production pakai CI automation.

- **Scanner format compatibility.** DefectDojo support JSON untuk Trivy/Semgrep/Checkov, tapi **butuh XML untuk ZAP**. Selalu cek `scan_type` compatibility sebelum upload — `zap-report.json` ditolak dengan "Wrong file format".

- **Terminology drift API vs UI.** API v2 tetap pakai "Product Type/Product", tapi UI v3.x rename ke "Organization/Asset". Tutorial outdated. Selalu verify terminology dengan actual UI saat setup.

---

## 🔗 Referensi

- [DefectDojo API v2 Documentation](https://documentation.defectdojo.com/integrations/api/)
- [DefectDojo import-scan Endpoint](https://documentation.defectdojo.com/integrations/api/#import-scan)
- [DefectDojo Supported Scanners](https://documentation.defectdojo.com/integrations/parsers/)
- [GitHub Actions Conditional Jobs](https://docs.github.com/en/actions/using-jobs/using-conditional-jobs)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | 46 findings centralized! Single pane of glass achieved |
| Pemahaman materi | 5 | import-scan API, conditional CI jobs, scanner format compatibility |
| Progres sesuai target | 5 | 7 reports uploaded, CI job added, terminology fixed |

---

## ➡️ Rencana Besok

- [ ] Hari 48: Intelligent Alert Routing (n8n & Slack) — CRITICAL alert → Slack channel

---

*[← Hari 46](hari-46.md) | [Hari 48 →](hari-48.md)*