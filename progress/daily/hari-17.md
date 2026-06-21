# Hari 17 — Container Image Scan

**📅 Tanggal:** 2026-06-21  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Scan distroless image dengan Trivy
- [x] Build naive image (alpine + build tools) untuk perbandingan
- [x] Scan naive image dengan Trivy
- [x] Bandingkan hasil CVE dan ukuran image
- [x] Simpan report JSON di `security/`

---

## ✅ Yang Berhasil Dikerjakan

- Scan `securebank:v1` (distroless): **0 CVE, 7.97MB**
- Build `securebank:naive` (alpine builder stage): **0 CVE, 352MB**
- Simpan kedua report JSON:
  - `security/trivy-image-report.json` (distroless)
  - `security/trivy-naive-image-report.json` (naive)
- Naive image sudah dihapus dari lokal Docker

---

## 📝 Catatan Teknis

```bash
# Distroless image scan
$ trivy image --severity HIGH,CRITICAL securebank:v1
Target                          Type      Vulnerabilities  Secrets
securebank:v1 (debian 12.14)   debian    0               -
securebank                       gobinary  0               -

# Naive image scan
$ trivy image --severity HIGH,CRITICAL securebank:naive
Target                            Type      Vulnerabilities  Secrets
securebank:naive (alpine 3.24.1) alpine    0               -
securebank                        gobinary  0               -
usr/local/go/bin/go               gobinary  0               -
... (9 more Go tool binaries)

# Image size comparison
securebank:v1     7.97MB    ← distroless
securebank:naive  352MB     ← alpine builder
```

**Hasil:** Kedua image 0 CVE, tapi perbedaan **size 44x** (7.97MB vs 352MB). Ini karena naive image membawa seluruh Go toolchain + alpine packages.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Kedua image 0 CVE — kurang dramatis untuk perbandingan | Alpine 3.24.1 sangat baru (rilis 2025). Di dunia nyata, base image yang lebih lama biasanya punya CVE. Size difference 44x masih sangat signifikan |
| Naive image 352MB, bukan 800MB seperti di estimasi | golang:1.26-alpine lebih ringan dari golang:1.26 biasa (tanpa alpine). Estimasi 800MB itu untuk full golang image |

---

## 📤 Output Hari Ini

- [x] `security/trivy-image-report.json` — distroless image scan (0 CVE)
- [x] `security/trivy-naive-image-report.json` — naive image scan (0 CVE)
- [x] Perbandingan: distroless 7.97MB vs naive 352MB (44x lebih kecil)
- [x] Commit: `security: add Trivy image scan reports (Day 17)`

---

## 💡 Pelajaran Baru

- **0 CVE di kedua image itu hasil yang baik.** Alpine 3.24.1 dan distroless/static-debian12:nonroot keduanya image yang well-maintained. Di dunia nyata, kalau pakai base image yang lebih lama (misal alpine 3.18), CVE count biasanya lebih tinggi.

- **CVE count bukan satu-satunya metric keamanan.** Attack surface lebih penting: naive image punya shell, package manager, dan 11 Go tool binaries yang bisa di-exploit. Distroless punya hanya binary aplikasi + CA certs.

- **Size difference 44x itu dramatis.** 7.97MB vs 352MB. Distroless image lebih cepat di-pull, lebih kecil di registry, dan lebih murah di storage. Di scale, ini berarti hemat bandwidth dan biaya cloud.

- **Trivy `image` mode berbeda dari `fs` mode.** `trivy fs` scan filesystem dependencies (go.mod), `trivy image` scan container image layers (OS packages + binary dependencies). Keduanya penting — fs untuk dev-time, image untuk runtime.

- **Multi-stage buildeliminasi clang/chain, bukan sekadar shrink size.** Builder stage dengan Go toolchain itu attack vector. Kalau image naif di-deploy ke production, attacker yang masuk bisa compile code langsung di container.

---

## 🔗 Referensi

- [Trivy Image Scanning](https://aquasecurity.github.io/trivy/latest/docs/target/container_image/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Scan image itu straightforward, perbandingan visualnya impactful |
| Pemahaman materi | 5 | Ngitung CVE vs attack surface = lebih dari sekadar CVE count |
| Progres sesuai target | 5 | Image scan selesai, Fase 2 lanjut |

---

## ➡️ Rencana Besok

- [ ] Hari 18: Dockerfile Hardening — `USER nonroot`, `COPY --chown`, read-only filesystem, health check

---

*[← Hari 16](hari-16.md) | [Hari 18 →](hari-18.md)*