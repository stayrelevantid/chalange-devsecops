# Hari 25 — DAST di Pipeline (GitHub Actions)

**📅 Tanggal:** 2026-07-02  
**⏱️ Durasi Belajar:** ~1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Tambah job DAST scan ke `ci.yml` dengan OWASP ZAP
- [x] Build Go binary langsung (tanpa Docker build) untuk efisiensi
- [x] ZAP baseline scan berjalan otomatis di CI
- [x] Report ZAP ter-upload sebagai artifact
- [x] Pipeline RED (expected) — ZAP nemu WARN 10049

---

## ✅ Yang Berhasil Dikerjakan

- Tambah job `dast-scan` ke `.github/workflows/ci.yml` (job ke-5)
- Pendekatan: Go binary langsung (bukan Docker build) — lebih ringan dan cepat
- Job `dast-scan` di-trigger setelah `build-and-test` selesai (`needs: [build-and-test]`)
- ZAP action: `zaproxy/action-baseline@v0.13.0` dengan target `http://localhost:8080`
- `cmd_options: '-I'` untuk ignore INFO findings
- Artifact `zap-dast-report` ter-upload (13KB)
- Pipeline run: 2 menit 16 detik
- Hasil: **0 FAIL, 1 WARN (10049), 66 PASS** — sama dengan local scan Day 24
- Pipeline RED (expected) — ZAP action return exit code 1 karena ada WARN

---

## 📝 Catatan Teknis

### DAST Job di ci.yml
```yaml
dast-scan:
  name: DAST Scan (OWASP ZAP)
  needs: [build-and-test]
  runs-on: ubuntu-latest
  defaults:
    run:
      working-directory: securebank-api
  steps:
    - checkout
    - setup Go 1.26
    - go build -o securebank ./cmd/api
    - run: JWT_SECRET=ci-test-secret PORT=8080 ./securebank &
    - wait for health (curl retry 15x, sleep 2s)
    - zaproxy/action-baseline@v0.13.0 (target: localhost:8080, cmd_options: -I)
    - upload artifact: zap-dast-report
```

### Pipeline Result
```
Build & Test (Go 1.26):      success ✅
Secret Scan (Gitleaks):      success ✅
SCA Scan (Trivy):            success ✅
SAST Scan (Semgrep):         success ✅
DAST Scan (OWASP ZAP):       failure ❌ (WARN 10049)
```

### ZAP Scan Output di CI
```
FAIL-NEW: 0  FAIL-INPROG: 0  WARN-NEW: 1  WARN-INPROG: 0  INFO: 0  IGNORE: 0  PASS: 66
WARN-NEW: Storable and Cacheable Content [10049] x 2
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Docker build di CI berat (~30s build + container lifecycle) | Pakai Go binary langsung: `go build -o securebank ./cmd/api`, run as background process. Lebih ringan, lebih cepat |
| App butuh beberapa detik untuk start di CI | Wait loop: `curl -sf http://localhost:8080/health` max 15x dengan sleep 2s |
| ZAP action default fail pada WARN | Biarkan fail — ini security gate. Fix finding di Day 26 |
| Background process cleanup | GitHub Actions auto-cleanup orphan processes di akhir job (terminate pid securebank) |

---

## 📤 Output Hari Ini

- [x] `.github/workflows/ci.yml` — tambah job `dast-scan` (job ke-5)
- [x] Pipeline run: ID 28572291850 — DAST job failure (expected)
- [x] Artifact: `zap-dast-report` (13KB, HTML report)
- [x] Pipeline total: 2 menit 16 detik
- [x] Commit: `5482b9d`

---

## 💡 Pelajaran Baru

- **Go binary approach lebih efisien dari Docker untuk DAST di CI.** Build binary cuma ~5 detik vs Docker build ~30 detik + container lifecycle. DAST scan HTTP response — jadi hasilnya sama whether app jalan sebagai binary atau container. Docker tetap penting untuk production, tapi untuk CI DAST scan, binary cukup.

- **ZAP action baseline default fail pada WARN.** `zaproxy/action-baseline@v0.13.0` return exit code 1 kalau ada WARN atau FAIL. Ini berbeda dari local run yang cuma fail pada FAIL. Di CI, WARN juga dianggap failure — security gate yang lebih ketat.

- **`cmd_options: '-I'` untuk ignore INFO findings.** ZAP baseline scan punya 3 level: FAIL, WARN, INFO. `-I` skip INFO findings. Tanpa `-I`, ZAP action akan fail pada INFO juga — terlalu noisy.

- **`needs: [build-and-test]` untuk sequential dependency.** DAST scan butuh konfirmasi dulu bahwa build dan test lolos. Kalau build gagal, DAST scan nggak jalan — hemat CI minutes. Job lain (Gitleaks, Trivy, Semgrep) tetap parallel — nggak perlu nunggu build.

- **GitHub Actions auto-cleanup orphan processes.** Log menunjukkan: "Terminate orphan process: pid (3962) (securebank)". Setup `./securebank &` (background) otomatis di-kill di akhir job. Nggak perlu manual `kill`.

---

## 🔗 Referensi

- [zaproxy/action-baseline GitHub Action](https://github.com/zaproxy/action-baseline)
- [GitHub Actions services vs background process](https://docs.github.com/en/actions/using-containerized-services/about-service-containers)
- [ZAP Docker baseline scan options](https://www.zaproxy.org/docs/docker/baseline-scan/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | DAST di pipeline jalan! 5 jobs sekarang |
| Pemahaman materi | 5 | Paham trade-off Docker vs binary di CI, ZAP exit code behavior |
| Progres sesuai target | 5 | Pipeline RED expected, fix di Day 26 |

---

## ➡️ Rencana Besok

- [ ] Hari 26: DAST Remediation — fix WARN 10049 (Storable and Cacheable Content), tambah security headers tambahan, pipeline hijau

---

*[← Hari 24](hari-24.md) | [Hari 26 →](hari-26.md)*