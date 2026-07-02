# Hari 24 — DAST Setup (OWASP ZAP)

**📅 Tanggal:** 2026-07-02  
**⏱️ Durasi Belajar:** ~2 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Jalankan SecureBank API di Docker container
- [x] Run OWASP ZAP Baseline Scan via Docker
- [x] Generate HTML + JSON report
- [x] Review temuan ZAP

---

## ✅ Yang Berhasil Dikerjakan

- Pull ZAP Docker image `ghcr.io/zaproxy/zaproxy:stable` (~300MB, butuh waktu lama karena network lambat)
- Build + start SecureBank API via `docker compose up -d --build`
- ZAP Baseline Scan via Docker dengan target `http://host.docker.internal:8080`
- Report HTML + JSON + YAML disimpan di `security/`
- Hasil: **0 FAIL, 1 WARN, 66 PASS**

---

## 📝 Catatan Teknis

### Command ZAP Baseline Scan via Docker
```bash
docker run --rm \
  -v "$(pwd)/securebank-api/security:/zap/wrk" \
  ghcr.io/zaproxy/zaproxy:stable \
  zap-baseline.py \
  -t http://host.docker.internal:8080 \
  -r zap-report.html \
  -J zap-report.json \
  -I
```

### Hasil ZAP Scan
```
FAIL-NEW: 0  FAIL-INPROG: 0  WARN-NEW: 1  WARN-INPROG: 0  INFO: 0  IGNORE: 0  PASS: 66
```

### ZAP Findings Detail

| ID | Name | Risk | Status |
|----|------|------|--------|
| 10049 | Storable and Cacheable Content | Informational (Medium) | WARN-NEW |
| 10021 | X-Content-Type-Options Header Missing | — | PASS ✅ |
| 10035 | Strict-Transport-Security Header | — | PASS ✅ |
| 10038 | Content Security Policy (CSP) Header Not Set | — | PASS ✅ |
| 10063 | Permissions Policy Header Not Set | — | PASS ✅ |
| 10020 | Anti-clickjacking Header | — | PASS ✅ |
| 10010 | Cookie No HttpOnly Flag | — | PASS ✅ |
| 10011 | Cookie Without Secure Flag | — | PASS ✅ |

### Root Cause WARN Finding
- **10049 — Storable and Cacheable Content**: 404 responses dari `/` dan `/robots.txt` tidak punya `Cache-Control` header. Karena 404 adalah response default Go's `http.HandleFunc` untuk path yang tidak didefinisikan, response ini tidak melalui middleware `SecurityHeaders` yang set `Cache-Control: no-store`.
- Fix akan dilakukan di Day 26 (DAST Remediation)

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Docker pull ZAP image lambat (2 layer retry terus) | Tunggu dengan background process. Layer `ccce74329ddc` dan `e01c689efb3b` butuh retry berkali-kali tapi akhirnya selesai |
| `docker compose up` dari root repo gagal (`no configuration file provided`) | Pakai `docker compose -f <full-path>/docker-compose.yml --project-directory <full-path> up -d --build` |
| macOS Docker Desktop: `--network host` tidak berfungsi | Pakai `host.docker.internal:8080` sebagai target URL untuk ZAP |
| ZAP local install (crossplatform zip) butuh `zap_common.py` + Python deps | Skip local install, pakai Docker image yang sudah lengkap |

---

## 📤 Output Hari Ini

- [x] `security/zap-report.html` — ZAP report HTML
- [x] `security/zap-report.json` — ZAP report JSON
- [x] `security/zap.yaml` — ZAP automation framework config
- [x] ZAP scan: 0 FAIL, 1 WARN, 66 PASS
- [x] Commit: `c8241a4`

---

## 💡 Pelajaran Baru

- **DAST vs SAST beda pendekatan.** SAST (Semgrep) scan source code tanpa menjalankan aplikasi. DAST (ZAP) scan aplikasi yang sedang berjalan — menemukan kerentanan runtime yang SAST tidak bisa lihat, seperti missing HTTP headers di response.

- **ZAP Baseline Scan = passive scan.** Baseline scan tidak melakukan active attack (SQL injection, XSS injection). Hanya spidering + passive scanning (analisis response). Untuk active scan, gunakan `zap-full-scan.py` atau `zap-api-scan.py`.

- **`host.docker.internal` di macOS Docker Desktop.** Beda dengan Linux yang pakai `--network host`, macOS Docker Desktop tidak support `--network host`. Pakai `host.docker.internal` untuk akses host dari dalam container.

- **Middleware hanya apply ke route yang didefinisikan.** Go's `http.HandleFunc` hanya apply middleware ke handler yang explicit didaftarkan (`/health`, `/balance`, `/transfer`). Path lain (seperti `/` atau `/robots.txt`) dapat default 404 response tanpa middleware — ZAP menemukan ini.

- **66 PASS dari 67 checks = 98.5% pass rate.** Security headers yang sudah diimplementasi di Day 15 (SecurityHeaders middleware) terbukti efektif. ZAP mengkonfirmasi: HSTS, CSP, X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Permissions Policy — semua PASS.

---

## 🔗 Referensi

- [OWASP ZAP Docker images](https://hub.docker.com/r/softwaresecurityproject/zaproxy)
- [ZAP Baseline Scan documentation](https://www.zaproxy.org/docs/docker/baseline-scan/)
- [ZAP GitHub Actions integration](https://github.com/zaproxy/action-baseline)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | DAST pertama! ZAP via Docker akhirnya jalan |
| Pemahaman materi | 5 | Paham bedanya baseline scan vs full scan, DAST vs SAST |
| Progres sesuai target | 5 | 66 pass, 0 fail — security headers ternyata bekerja |

---

## ➡️ Rencana Besok

- [ ] Hari 25: DAST di Pipeline — ZAP scan otomatis di GitHub Actions CI

---

*[← Hari 23](hari-23.md) | [Hari 25 →](hari-25.md)*