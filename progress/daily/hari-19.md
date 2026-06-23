# Hari 19 — Image Signing (Cosign)

**📅 Tanggal:** 2026-06-23  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install Cosign binary
- [x] Generate key pair (cosign.key + cosign.pub)
- [x] Add `cosign.key` ke `.gitignore` (private key tidak di-commit)
- [x] Push image ke GHCR (GitHub Container Registry)
- [x] Sign image dengan Cosign
- [x] Verify signature dengan public key

---

## ✅ Yang Berhasil Dikerjakan

- Install Cosign v3.1.1 (下载 dari GitHub releases, binary `darwin-arm64`)
- Generate key pair: `cosign.key` (private) + `cosign.pub` (public)
- `cosign.key` ditambahkan ke `.gitignore` — tidak akan di-commit
- Tag image: `securebank:v1` → `ghcr.io/stayrelevantid/securebank:v1`
- Push image ke GHCR: berhasil, digest `sha256:df0ecb33...`
- Sign image: `cosign sign --key cosign.key ghcr.io/stayrelevantid/securebank:v1`
- Verify signature: cosign verify — 3/3 checks lulus (claims, transparency log, public key)

---

## 📝 Catatan Teknis

```bash
# Install Cosign (binary langsung dari GitHub releases)
$ curl -sL -o /opt/homebrew/bin/cosign \
  "https://github.com/sigstore/cosign/releases/download/v3.1.1/cosign-darwin-arm64"
$ chmod +x /opt/homebrew/bin/cosign
$ cosign version
cosign: A tool for Container Signing, Verification and Storage in an OCI registry.
GitVersion: v3.1.1
Platform: darwin/arm64

# Generate key pair (passwordless untuk non-interactive CI)
$ COSIGN_PASSWORD="" cosign generate-key-pair
Private key written to cosign.key
Public key written to cosign.pub

# Tag image untuk GHCR
$ docker tag securebank:v1 ghcr.io/stayrelevantid/securebank:v1

# Push image ke GHCR
$ docker push ghcr.io/stayrelevantid/securebank:v1
v1: digest: sha256:df0ecb33e4efc30ed24b6c45404dee21c1814cd8c6acad66c87f4f66efc9faa3

# Sign image
$ COSIGN_PASSWORD="" cosign sign --key cosign.key ghcr.io/stayrelevantid/securebank:v1
Signing artifact...
Pushing signature to: ghcr.io/stayrelevantid/securebank

# Verify signature
$ COSIGN_PASSWORD="" cosign verify --key cosign.pub ghcr.io/stayrelevantid/securebank:v1

Verification for ghcr.io/stayrelevantid/securebank:v1 --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The signatures were verified against the specified public key

[{"critical":{"identity":{"docker-reference":"ghcr.io/stayrelevantid/securebank:v1"},
"image":{"docker-manifest-digest":"sha256:df0ecb33..."},"type":"https://sigstore.dev/cosign/sign/v1"},
"optional":{}}]
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Homebrew bisa timeout saat install Cosign | Download binary langsung dari GitHub releases ke `/opt/homebrew/bin/` — lebih cepat dan reliable |
| Cosign minta password saat sign/verify | Set `COSIGN_PASSWORD=""` env var untuk non-interactive mode (key dibuat tanpa password) |
| Cosign warning: "uses a tag, not a digest" | Expected. Untuk production seharusnya pakai digest (`@sha256:...`), tapi untuk challenge ini tag `v1` sudah cukup |
| Butuh GitHub PAT untuk push ke GHCR | Buat PAT dengan scope `write:packages` di GitHub Settings → Developer settings → Personal access tokens |

---

## 📤 Output Hari Ini

- [x] `securebank-api/cosign.pub` — public key (di-commit, aman untuk share)
- [x] `securebank-api/cosign.key` — private key (gitignored, TIDAK di-commit)
- [x] `securebank-api/.gitignore` — ditambah `cosign.key`
- [x] Image `ghcr.io/stayrelevantid/securebank:v1` pushed ke GHCR
- [x] Image ditandatangani dengan Cosign
- [x] Signature terverifikasi (3/3 checks lulus)

---

## 💡 Pelajaran Baru

- **Image signing itu bukan tentang enkripsi, tapi tentang integritas.** Cosign tidak mengubah image. Dia membuat signature terpisah yang disimpan di registry sebagai OCI artifact. Verifikasi memastikan image yang di-pull persis sama dengan image yang di-sign.

- **Private key tidak pernah meninggalkan environment trusted.** `cosign.key` di `.gitignore`, tidak di-commit ke repo. Di CI, private key di-store sebagai GitHub Secret. Public key (`cosign.pub`) aman di-commit — siapapun bisa verifikasi, tapi cuma yang punya private key yang bisa sign.

- **Transparency log itu bukti signature tercatat.** Cosign mengirim signature ke Rekor (transparency log) sebagai audit trail. Verifikasi mencek bahwa signature ada di transparency log. Ini mencegah penyangkalan — signer tidak bisa bilang "saya tidak pernah sign image itu."

- **Tag vs digest di Cosign.** Cosign warning saat pakai tag (`:v1`) karena tag bisa di-repoint ke image berbeda. Best practice: pakai digest (`@sha256:...`) untuk signing. Tapi untuk challenge ini, tag `v1` sudah cukup karena kita kontrol full lifecycle.

- **GHCR gratis untuk public repo.** GitHub Container Registry tidak butuh subscription untuk public images. Pakai GitHub PAT dengan scope `write:packages` untuk push.

---

## 🔗 Referensi

- [Cosign Documentation](https://docs.sigstore.dev/cosign/)
- [Sigstore Transparency Log (Rekor)](https://docs.sigstore.dev/rekor/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packagesRegistry/working-with-the-container-registry)
- [Cosign Key Management](https://docs.sigstore.dev/cosign/key_management/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Cosign straightforward, sign+verify dalam 2 command |
| Pemahaman materi | 5 | Konsep signing: private sign, public verify. Mirip GPG |
| Progres sesuai target | 5 | Day 19 selesai, Fase 2 lanjut 4/15 |

---

## ➡️ Rencana Besok

- [ ] Hari 20: Terraform Setup + IaC Scan (Checkov) — `terraform/main.tf` (VPC+S3) + checkov scan

---

*[← Hari 18](hari-18.md) | [Hari 20 →](hari-20.md)*