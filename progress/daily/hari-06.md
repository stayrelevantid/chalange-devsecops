# Hari 06 — SCA di Pipeline (Quality Gate)

**📅 Tanggal:** 2026-06-10  
**⏱️ Durasi Belajar:** 1 jam  
**🏷️ Fase:** Fase 1 — Secure SDLC & AppSec  
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Tambah job `sca-scan` ke CI pipeline
- [x] Pipeline gagal/merah karena ada CVE CRITICAL ditemukan
- [x] Report JSON ter-upload sebagai artifact
- [x] `exit-code: 1` berfungsi sebagai quality gate

---

## ✅ Yang Berhasil Dikerjakan

- Tambah job `sca-scan` (Trivy) ke `.github/workflows/ci.yml` — berjalan paralel dengan `build-and-test` dan `secret-scan`
- Menggunakan `aquasecurity/trivy-action@master` dengan konfigurasi:
  - `scan-type: fs` — scan filesystem (go.mod)
  - `scan-ref: 'securebank-api'` — point ke subdirektori proyek
  - `severity: CRITICAL,HIGH` — hanya scan severity tinggi
  - `exit-code: 1` — gagalkan pipeline jika ditemukan
  - `format: json` — output dalam format JSON
- Upload Trivy report sebagai artifact (`if: always()`)
- Pipeline sekarang punya 3 job paralel: Build & Test, Secret Scan (Gitleaks), SCA Scan (Trivy)

---

## 📝 Catatan Teknis

```yaml
sca-scan:
  name: SCA Scan (Trivy)
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4

    - uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: 'securebank-api'
        format: 'json'
        output: 'trivy-sca.json'
        severity: 'CRITICAL,HIGH'
        exit-code: '1'

    - name: Upload Trivy Report
      if: always()
      uses: actions/upload-artifact@v4
      with:
        name: trivy-sca-report
        path: trivy-sca.json
```

### Keputusan desain:
- **`trivy-action@master`** dipilih (bukan binary install manual) karena Aquasecurity menyediakan GitHub Action gratis — beda dengan Gitleaks yang action-nya butuh lisensi berbayar
- **`scan-ref: 'securebank-api'`** — pointing ke subdirektori karena `go.mod` ada di sana, bukan di root repo
- **`.trivyignore` sengaja TIDAK dipakai di CI** — supaya pipeline gagal dan membuktikan quality gate berfungsi. `.trivyignore` hanya dipakai di lokal untuk dokumentasi accepted risk

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Pipeline hijau padahal seharusnya merah (4 CVE ditemukan) | `.trivyignore` di root repo mengecualikan semua 4 CVE — Trivy membaca file ini di CI dan meng-ignore semua temuan. Solusi: hapus `.trivyignore` supaya quality gate berfungsi. File ini hanya dipakai di lokal untuk dokumentasi accepted risk |

---

## 📤 Output Hari Ini

- [x] Job `sca-scan` ditambahkan ke `ci.yml`
- [x] Pipeline expected gagal (RED) karena 4 CVE ditemukan
- [x] Trivy report JSON sebagai artifact
- [x] Quality gate `exit-code: 1` dikonfigurasi

---

## 💡 Pelajaran Baru

- **`exit-code: 1`** adalah mekanisme quality gate — jika Trivy menemukan CVE dengan severity yang ditentukan, pipeline akan gagal. Ini memastikan kode yang rentan tidak masuk ke main branch
- **`trivy-action@master`** gratis dan tidak butuh lisensi — berbeda dengan `gitleaks-action@v2` yang sekarang butuh lisensi berbayar
- **`scan-ref`** harus mengarah ke direktori yang berisi `go.mod` — karena proyek kita ada di subdirektori `securebank-api/`, bukan di root repo
- **`if: always()`** memastikan report tetap di-upload meskipun Trivy menemukan CVE dan pipeline gagal — penting untuk debugging
- **`.trivyignore` hanya untuk lokal** — di CI, kita tidak pakai `.trivyignore` supaya pipeline gagal. File ini dihapus dari repo setelah diketahui menyebabkan false green di CI
- **`.trivyignore` vs CI** — `.trivyignore` menyebabkan semua CVE di-ignore sehingga pipeline hijau padahal ada 4 CVE. Harus dihapus supaya quality gate bekerja

---

## 🔗 Referensi

- [Trivy GitHub Action](https://github.com/aquasecurity/trivy-action)
- [Trivy Documentation](https://trivy.dev/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | CI pipeline makin lengkap |
| Pemahaman materi | 4 | Quality gate konsepnya jelas |
| Progres sesuai target | 5 | Lancar, tidak ada hambatan |

---

## ➡️ Rencana Besok

- [ ] Patch vulnerable dependencies (Day 07 - SCA Remediation)
- [ ] Ganti `dgrijalva/jwt-go` dengan `golang-jwt/jwt/v5`
- [ ] Upgrade gin dan transitive dependencies
- [ ] Trivy scan: 0 CRITICAL/HIGH
- [ ] Pipeline kembali hijau

---

*[← Hari 05](hari-05.md) | [Hari 07 →](hari-07.md)*