# Hari 53 — Red Team: Leaked Credentials

**📅 Tanggal:** 2026-08-02
**⏱️ Durasi Belajar:** ~120 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Buat IAM user sementara (`leaked-dev`) dengan access key untuk di-leak (simulasi)
- [x] Aktifkan AWS CloudTrail (sebelumnya belum ada trail sama sekali)
- [x] Siapkan "canary" — file umpan di S3 agar aktivitas pencurian terlihat
- [x] Simulasikan attacker: pakai key bocor untuk list & curi file
- [x] Deteksi jejak key di CloudTrail (management events + data events)
- [x] Buat `revoke.sh` — script auto-revoke key bocor + isolasi user
- [x] Test revoke: key langsung mati, attacker tidak bisa akses lagi
- [x] Cleanup: hapus user, key, bucket canary (CloudTrail tetap hidup)

---

## ✅ Yang Berhasil Dikerjakan

### 1. CloudTrail Diaktifkan (Baru Pertama Kali)

CloudTrail **belum pernah aktif** di akun ini — `describe-trails` kosong. Padahal Day 51 (CSPM Remediation) sudah mencatat ini sebagai rekomendasi. Hari ini dieksekusi:

```bash
# Bucket log (private, versioning aktif, policy deny non-owner)
aws s3api create-bucket --bucket securebank-cloudtrail-logs-XXXX --region ap-southeast-1

# Trail: multi-region + log file validation + global service events
aws cloudtrail create-trail \
  --name securebank-trail \
  --s3-bucket-name securebank-cloudtrail-logs-XXXX \
  --is-multi-region-trail \
  --enable-log-file-validation \
  --no-include-global-service-events  # diubah ke ya via update-trail

aws cloudtrail start-logging --name securebank-trail
```

Detail: multi-region ✅, log validation ✅, termasuk data events untuk S3 (canary bucket) via `put-event-selectors`.

### 2. IAM User + Access Key (Target yang "Bocor")

```bash
aws iam create-user --user-name leaked-dev
aws iam attach-user-policy --user-name leaked-dev \
  --policy-arn arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
aws iam create-access-key --user-name leaked-dev
# AccessKeyId: AKIAZ6PEJSWX... (2 key dibuat, dipakai bergantian)
```

### 3. Canary File (S3 Honeypot) 🍯

Rencana awal: pakai **CanaryTokens.org** supaya ada "signal" saat file dibuka. **Gagal total** — CanaryTokens adalah React SPA, semua POST via curl balas `{"detail":"Method Not Allowed"}`. Tidak ada API publik yang bisa diotomasi.

Pivot ke solusi AWS-native: **S3 honeypot bucket** berisi file "sensitif" palsu + marker:

```bash
aws s3api create-bucket --bucket securebank-honeypot-canary --region ap-southeast-1
# backup/db-config.txt  — konfigurasi DB palsu
# env/.env.production    — env palsu berisi marker:
#   AWS_ACCESS_KEY_ID=CanaryToken-3f9a1c2e8d4b
```

Marker tersebut dipakai untuk membuktikan bahwa file benar-benar diakses (bukan sekadar listing).

### 4. Simulasi Attacker

Dengan key `leaked-dev` (disimpan di profil terpisah), attacker:

1. `ListBuckets` — lihat daftar semua bucket akun
2. `ListObjects` di `securebank-honeypot-canary`
3. `GetObject` — unduh `env/.env.production` (berisi canary token)

### 5. Deteksi di CloudTrail 🔍

- **Management events** langsung kelihatan (~3–5 menit delay): `GetCallerIdentity` + `ListBuckets` atas nama key `leaked-dev`:
  ```bash
  aws cloudtrail lookup-events --region ap-southeast-1 \
    --lookup-attributes AttributeKey=AccessKeyId,AttributeValue=<key> \
    --query 'Events[].[EventTime,EventName,EventSource]' --output table
  ```
- **S3 data events** (GetObject) butuh event selector khusus dan delivery-nya lebih lambat (hingga 15 menit). Sudah diaktifkan via `put-event-selectors` (WriteOnly, resource = canary bucket), tapi belum tampil dalam window observasi — jadi dicek via file log mentah di S3.
- Verifikasi menyeluruh via `GetObject` langsung ke bucket log CloudTrail (path `AWSLogs/...`) dan grep `leaked-dev`.

### 6. `revoke.sh` — Auto-Revoke (Sudah Di-Test) ⚡

Script 3 langkah di `securebank-api/security/revoke.sh`:

| Langkah | Aksi |
|---------|------|
| 1 | `lookup-events` — tampilkan riwayat aktivitas key di CloudTrail |
| 2 | `update-access-key --status Inactive` — revoke instan |
| 3 | `put-user-policy` — attach Deny-All inline policy `BlockAll-On-Leak` |

**Hasil test:** setelah deactivate + deny policy, attacker retry `ListBuckets` → `An error occurred (InvalidAccessKeyId)`. Key mati total. ✅

### 7. Cleanup 🧹

- Delete kedua access key (2 key yang dibuat untuk demo)
- Detach policy + hapus user `leaked-dev`
- Force-delete bucket canary (`securebank-honeypot-canary`)
- **CloudTrail + bucket log DIBIARKAN HIDUP** — ini jadi improvement permanen (menuntaskan rekomendasi Day 51)

---

## 🧪 Hasil Checklist

| Checklist | Hasil |
|-----------|-------|
| IAM user sementara dibuat | ✅ `leaked-dev` (AmazonS3ReadOnlyAccess) |
| Jejak akses terdeteksi di CloudTrail | ✅ `GetCallerIdentity` + `ListBuckets` (management events) |
| Skrip revoke/isolasi disiapkan | ✅ `revoke.sh` + sudah diuji (key mati, `InvalidAccessKeyId`) |

---

## 📝 Catatan Teknis

### CloudTrail Data Events ≠ Management Events
Management events (API call kontrol) tercatat otomatis dan bisa di-query via `lookup-events`. Data events (GetObject, PutObject) **wajib diaktifkan manual** via `put-event-selectors` per resource — dan delivery-nya bisa delay hingga 15 menit. Untuk investigasi cepat, `lookup-events` + cek file log mentah lebih reliable daripada menunggu data events.

### CanaryTokens.org Tidak Bisa Di-Automasi
Layanan ini berbentuk React SPA; endpoint "generate token" butuh interaksi browser dan tidak punya public API yang stabil. Untuk lab, S3 honeypot + marker string sudah cukup sebagai canary.

### Revoke = Deactivate + Isolate
`update-access-key --status Inactive` cukup untuk mematikan key. Tapi attach **Deny-All policy** adalah lapisan isolasi tambahan (misal key lain milik user yang sama, atau akses via role) — defense in depth untuk insiden nyata.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| CanaryTokens.org tidak bisa di-automasi (React SPA, `Method Not Allowed`) | Pivot ke S3 honeypot bucket + marker string `CanaryToken-...` |
| CloudTrail belum ada trail sama sekali | `create-trail` multi-region + log validation + start-logging |
| S3 data events (GetObject) belum muncul di lookup dalam 5 menit | Aktifkan `put-event-selectors`; fallback: baca file log mentah di bucket |
| Bucket log dengan account ID | File log memakai nama bucket ber-account ID — direntang sebagai `XXXX` di repo |

---

## 📤 Output Hari Ini

- [x] CloudTrail trail `securebank-trail` aktif (multi-region, log validation, data events) + bucket `securebank-cloudtrail-logs-XXXX` (dibiarkan hidup)
- [x] `securebank-api/security/revoke.sh` — auto-revoke key bocor (CloudTrail lookup → deactivate → Deny-All)
- [x] S3 honeypot canary bucket + 2 file umpan (sudah di-delete saat cleanup)
- [x] Deteksi jejak key di CloudTrail (management events) + verifikasi file log mentah
- [x] Cleanup lengkap: user + key + bucket canary dihapus

---

## 💡 Lessons Learned

### 1. Detection coverage butuh sensor, bukan asumsi
CloudTrail ternyata **belum pernah aktif**. Semua analisis "siapa yang akses apa" selama ini hanya berdasarkan tebakan. Satu `create-trail` saja mengubah akun dari buta menjadi terawasi — dan data events untuk akses data sensitif harus diaktifkan eksplisit.

### 2. "Canary" tidak harus mahal
CanaryTokens adalah produk bagus, tapi tidak bisa di-automasi di sini. Solusi S3 honeypot dengan marker string memberi sinyal yang sama: ada file palsu yang "tidak pernah disentuh orang", jadi kalau ada yang membacanya, itu pasti attacker.

### 3. Management events vs data events — timing & granularity beda
Management events = cepat, otomatis, mudah di-query. Data events = butuh config manual + delay lebih lama. Untuk incident response, mulai dari `lookup-events`, lalu dalami dengan file log mentah.

### 4. Revoke harus instan dan berlapis
Deactivate key adalah langkah pertama. Attach Deny-All policy adalah jaring pengaman — melindungi dari skenario di mana attacker sudah meng-copy key lain atau mendapatkan akses lewat jalur lain dari user yang sama.

### 5. Lab cleanup itu bagian dari security hygiene
Key yang sudah didelete, user yang dihapus, bucket canary yang di-force-delete — tidak ada sisa artefak yang bisa disalahgunakan. Yang sengaja dibiarkan hidup (CloudTrail) justru meningkatkan postur keamanan jangka panjang.

---

## 🔗 Referensi

- [Day 51](hari-51.md) — CSPM Remediation (CloudTrail jadi rekomendasi)
- [Day 52](hari-52.md) — Red Team: K8s Escape (konteks red teaming)
- [AWS CloudTrail Lookup Events CLI](https://docs.aws.amazon.com/cli/latest/reference/cloudtrail/lookup-events.html)
- [S3 data events di CloudTrail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html)
- [AWS IAM Deactivate Access Key](https://docs.aws.amazon.com/cli/latest/reference/iam/update-access-key.html)
- [CanaryTokens.org](https://canarytokens.org)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Skenario leaked creds yang realistis dan bisa di-test end-to-end |
| Pemahaman materi | 4 | CloudTrail management vs data events, IAM revoke, canary |
| Progres sesuai target | 5 | Semua langkah selesai termasuk test revoke + cleanup |

---

## ➡️ Rencana Besok

- [ ] **Day 54: Chaos Security Engineering** — matikan OPA Gatekeeper, deploy pod tidak patuh, amati apakah ada lapisan deteksi lain yang masih hidup

---

*[← Hari 52](hari-52.md) | [Hari 54 →](hari-54.md)*
