# Fase 4: Vulnerability Management & Red Teaming (Hari 46–60)

> **Proyek:** SecureBank API — Visibilitas ancaman terpusat, respons otomatis, dan simulasi serangan (Red Teaming)
> **Output Fase:** Dashboard kerentanan, workflow remediase dengan AI, eksekusi Red Teaming, dan laporan audit profesional.

---

## Hari 46: Vulnerability Management Dashboard (DefectDojo)

### Tujuan
Menggelar DefectDojo secara lokal sebagai Single Pane of Glass untuk semua temuan keamanan.

### Tutorial

**1. Clone dan jalankan DefectDojo:**
```bash
git clone https://github.com/DefectDojo/django-DefectDojo.git
cd django-DefectDojo
# Gunakan profile docker-compose untuk setup ringan
./dc-build.sh
docker-compose --profile postgres-redis up -d
```

**2. Dapatkan kredensial admin:**
```bash
docker-compose logs initializer | grep "Admin password"
```

**3. Login ke Dashboard:**
- Buka `http://localhost:8080`
- Login dengan username `admin` dan password yang didapat.

**4. Setup awal di DefectDojo:**
- Buat Product Type: `Fintech`
- Buat Product: `SecureBank API`
- Buat Engagement: `Q3 Security Audit`

### Checklist
- [ ] DefectDojo berjalan di lokal.
- [ ] Berhasil login sebagai Admin.
- [ ] Product `SecureBank API` telah dibuat di dashboard.

---

## Hari 47: DefectDojo API Integration

### Tujuan
Mengotomatiskan pengunggahan laporan (Trivy, Semgrep) dari CI/CD pipeline ke DefectDojo.

### Tutorial

**1. Dapatkan API Key DefectDojo:**
- Di DefectDojo, klik profil kanan atas -> `API v2 Key`.
- Salin key tersebut.

**2. Tambahkan API Key ke GitHub Secrets:**
- Tambahkan secret `DEFECTDOJO_API_KEY` ke repositori GitHub.

**3. Update `.github/workflows/security-scan.yml`:**
Tambahkan job baru di akhir pipeline (butuh jobs sebelumnya selesai):
```yaml
  upload-defectdojo:
    needs: [sca-scan, sast-scan]
    runs-on: ubuntu-latest
    if: always()
    steps:
      - uses: actions/checkout@v4
      - name: Download all artifacts
        uses: actions/download-artifact@v4
      - name: Upload Trivy to DefectDojo
        run: |
          curl -X POST "http://<defectdojo-url>/api/v2/import-scan/" \
            -H "Authorization: Token ${{ secrets.DEFECTDOJO_API_KEY }}" \
            -F "scan_type=Trivy Scan" \
            -F "product_name=SecureBank API" \
            -F "engagement_name=CI/CD Pipeline" \
            -F "file=@trivy-sca.json"
      - name: Upload Semgrep to DefectDojo
        run: |
          curl -X POST "http://<defectdojo-url>/api/v2/import-scan/" \
            -H "Authorization: Token ${{ secrets.DEFECTDOJO_API_KEY }}" \
            -F "scan_type=Semgrep JSON Report" \
            -F "product_name=SecureBank API" \
            -F "engagement_name=CI/CD Pipeline" \
            -F "file=@semgrep.sarif"
```
*(Catatan: Sesuaikan URL DefectDojo jika dijalankan di VM/Cloud, atau gunakan ngrok/tunnel untuk akses lokal).*

### Checklist
- [ ] API Key DefectDojo tersimpan di GitHub Secrets.
- [ ] Pipeline berhasil mengirim data ke DefectDojo (cek HTTP 201 Created).
- [ ] Temuan Trivy dan Semgrep muncul di dashboard DefectDojo.

---

## Hari 48: Intelligent Alert Routing (n8n & Slack)

### Tujuan
Mengirimkan notifikasi darurat (CRITICAL severity) ke Slack via n8n.

### Tutorial

**1. Siapkan Slack Workspace:**
- Buat channel `#security-alerts`.
- Buat Slack App, aktifkan `Incoming Webhooks`, salin Webhook URL.

**2. Update Workflow n8n (Lanjutan Hari 42):**
- Buka workflow webhook Falco di `http://localhost:5678`.
- Pada cabang **True** dari node IF (Priority == Critical), tambahkan node **HTTP Request** atau **Slack**.
- Jika menggunakan Slack node:
  - Hubungkan kredensial webhook.
  - Set text: `🚨 CRITICAL ALERT DETECTED 🚨\nRule: {{ $json.rule }}\nPod: {{ $json.output_fields["k8s.pod.name"] }}\nMessage: {{ $json.output }}`

**3. Uji Workflow:**
- Jalankan ulang trigger simulasi Falco dari Hari 41 yang memicu Critical Alert.
- Periksa channel Slack, pastikan pesan masuk.

### Checklist
- [ ] Slack Incoming Webhook terkonfigurasi.
- [ ] Node Slack ditambahkan di n8n untuk severity Critical.
- [ ] Pesan alert berhasil masuk ke channel Slack.

---

## Hari 49: AI Remediation Node di n8n

### Tujuan
Menggunakan LLM untuk memberikan ringkasan perbaikan instan ke developer saat ada temuan SAST baru.

### Tutorial

**1. Buat Webhook baru di n8n untuk SAST:**
- Buat workflow baru. Trigger: Webhook (POST `/sast-alert`).
- Asumsikan inputnya adalah temuan SAST (JSON dari DefectDojo Webhook atau custom script di CI).

**2. Tambahkan node OpenAI / LLM di n8n:**
- Masukkan API Key (OpenAI / Anthropic).
- Prompt:
  ```text
  Berikut adalah temuan kerentanan statis pada kode Golang:
  Judul: {{$json.title}}
  Deskripsi: {{$json.description}}
  File: {{$json.file_path}}

  Buat ringkasan maksimal 3 kalimat tentang apa bahayanya dan berikan contoh perbaikan kodenya.
  ```

**3. Kirim ke Slack:**
- Tambahkan node Slack setelah LLM.
- Kirim pesan hasil generate LLM ke channel `#dev-alerts`.

### Checklist
- [ ] Workflow n8n baru untuk SAST remediation dibuat.
- [ ] Node LLM terhubung dan menghasilkan rekomendasi.
- [ ] Rekomendasi dikirim secara otomatis ke Slack.

---

## Hari 50: Cloud Security Posture Management (CSPM)

### Tujuan
Memindai lingkungan Cloud (AWS/GCP) untuk kepatuhan CIS Benchmarks.

### Tutorial

**1. Siapkan Cloud Sandbox:**
- Pastikan Anda memiliki akun AWS (Free Tier) dan telah mengatur AWS CLI (`aws configure`).

**2. Install Prowler (AWS):**
```bash
pip install prowler
```

**3. Jalankan Prowler Scan:**
```bash
prowler aws
```
Atau output ke HTML/JSON:
```bash
prowler aws -M html json -f ap-southeast-1
```

**4. Analisis Temuan:**
- Buka file HTML hasil Prowler (di folder `output/`).
- Perhatikan warna merah (FAIL) pada bagian IAM, S3, dan Networking.

### Checklist
- [ ] Prowler terinstall dan berhasil dijalankan di akun sandbox.
- [ ] Laporan HTML / JSON dihasilkan.
- [ ] Mengidentifikasi temuan kritikal (misal: MFA belum aktif).

---

## Hari 51: CSPM Remediation

### Tujuan
Memperbaiki temuan kritis dari Prowler di lingkungan cloud.

### Tutorial

**1. Analisis dan Perbaiki (Pilih 3 Temuan Kritis), Contoh:**
- **Temuan 1: IAM Root Account tidak pakai MFA.**
  - **Remediasi:** Login ke AWS Console -> IAM -> Aktifkan MFA untuk root user.
- **Temuan 2: S3 Bucket Publik.**
  - **Remediasi:** Edit Terraform `aws_s3_bucket_public_access_block` atau perbaiki langsung via Console.
- **Temuan 3: Default VPC digunakan / Security Group terbuka.**
  - **Remediasi:** Hapus default VPC jika tidak dipakai, atau pastikan SG default tidak mengizinkan Ingress/Egress (sesuai best practice).

**2. Re-scan dengan Prowler:**
```bash
prowler aws -c aws_iam_1_2,aws_s3_2_1  # (contoh spesifik check ID)
```

### Checklist
- [ ] Memahami setidaknya 3 temuan CIS Benchmarks.
- [ ] Remediasi berhasil diterapkan di AWS.
- [ ] Scan ulang membuktikan masalah (FAIL) berubah menjadi (PASS).

---

## Hari 52: Red Teaming — Scenario 1 (K8s Escape)

### Tujuan
Simulasi penyerang yang berhasil masuk ke pod dan mencoba kabur (escape) ke node Kubernetes.

### Tutorial

**1. Deploy "Bad Pod" (Privileged):**
Buat `k8s/redteam-pod.yaml`:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: attacker-pod
  namespace: default
spec:
  containers:
  - name: attacker
    image: alpine
    command: ["sleep", "3600"]
    securityContext:
      privileged: true # Ini harusnya diblokir OPA di namespace securebank
```
*(Deploy di default namespace untuk lolos dari OPA, anggap attacker masuk dari celah aplikasi lain).*
```bash
kubectl apply -f k8s/redteam-pod.yaml
```

**2. Coba eksploitasi:**
```bash
kubectl exec -it attacker-pod -- /bin/sh
# Di dalam pod:
nsenter -t 1 -m -u -n -i sh
# Sekarang Anda berada di shell host (node K8s)
cat /etc/kubernetes/kubelet.conf
```

**3. Cek Observability:**
- Lihat log Falco / n8n / Slack: Apakah kejadian `nsenter` atau `Sensitive file read` terdeteksi?

### Checklist
- [ ] Attacker pod (privileged) berhasil di-deploy.
- [ ] Eksploitasi nsenter berhasil menembus ke node.
- [ ] Falco mencatat aktivitas tidak wajar ini.

### Implementation Notes (Hari 52 — Hasil Praktek)
1. **`privileged: true` saja TIDAK cukup** untuk `nsenter -t 1` mencapai host — pod masih punya PID namespace sendiri (`nsenter -t 1` hanya masuk ke PID 1 pod). Tambah **`hostPID: true`** supaya PID 1 host terlihat. Kombinasi `privileged + hostPID` = escape vector klasik.
2. **Path kubelet di k3s beda dari kubeadm.** `/etc/kubernetes/kubelet.conf` tidak ada; k3s pakai `/var/lib/rancher/k3s/agent/kubelet.kubeconfig` + `client-kubelet.crt/key`.
3. **Cek ruleset Falco yang benar-benar ter-load.** Cluster tanpa internet → falcoctl gagal download full rules → hanya 26 rules minimal aktif, tidak ada rule privileged/setns. Tambahkan rule custom: `Privileged container launched`, `Container escape via setns`, `Read kubelet credentials post-escape`.
4. **Hati-hati FP:** rule kubelet jangan pakai prefix `/var/lib/rancher/k3s/agent` (containerd snapshot juga di sana) — match nama file persis (`endswith client-kubelet.crt`, dll).
5. **Gap OPA:** pastikan constraint `privileged: true` + `hostPID` benar-benar ada (default setup hanya require resource limits). Lihat `k8s/gatekeeper/constraints/deny-privileged.yaml`.

---

## Hari 53: Red Teaming — Scenario 2 (Leaked Credentials)

### Tujuan
Simulasi kunci akses AWS bocor dan respons otomatis pencabutannya.

### Tutorial

**1. Buat User Sementara & Leak (Sengaja):**
- Buat IAM User `leaked-dev` dengan akses S3ReadOnly.
- Buat Access Key.

**2. Trigger Alert:**
- Gunakan key tersebut dari IP/mesin asing, atau commit key tersebut ke GitHub publik sebentar (GitHub akan otomatis mendeteksi dan mengirim email).
- Atau gunakan layanan CanaryTokens.

**3. Observasi CloudTrail:**
- Buka AWS CloudTrail, perhatikan kapan key tersebut digunakan (`ConsoleLogin` atau `ListBuckets`).

**4. Script Revoke Otomatis (Bash):**
Buat skrip `revoke.sh` yang bisa dipicu:
```bash
#!/bin/bash
TARGET_USER=$1
aws iam update-access-key --access-key-id $2 --status Inactive --user-name $TARGET_USER
aws iam put-user-policy --user-name $TARGET_USER --policy-name BlockAll --policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Deny","Action":"*","Resource":"*"}]}'
echo "User $TARGET_USER revoked."
```

### Checklist
- [ ] IAM User sementara dibuat.
- [ ] Jejak akses terdeteksi di CloudTrail.
- [ ] Skrip revoke/isolasi telah disiapkan.

### Implementation Notes (Hari 53 — Hasil Praktek)
1. **CloudTrail kemungkinan belum aktif di akun Anda.** Cek dulu: `aws cloudtrail describe-trails`. Jika kosong, buat trail multi-region + log file validation + `start-logging`, dan tambahkan bucket log versi `securebank-cloudtrail-logs-XXXX` (private, deny non-owner).
2. **Management events cepat, data events lambat.** `GetCallerIdentity`/`ListBuckets` langsung muncul di `cloudtrail lookup-events` (~3–5 menit). S3 data events (GetObject) wajib diaktifkan manual via `put-event-selectors` dan bisa delay hingga 15 menit — fallback: baca file log mentah di bucket (path `AWSLogs/...`).
3. **CanaryTokens.org tidak bisa di-automasi** (React SPA, semua POST balas `Method Not Allowed`). Ganti dengan **S3 honeypot bucket** berisi file umpan `backup/db-config.txt` + `env/.env.production` yang memuat marker string `CanaryToken-...`. Marker = bukti file benar-benar diakses.
4. **Script revoke siap dipakai:** `securebank-api/security/revoke.sh` — alur: (1) tampilkan riwayat key via CloudTrail, (2) deactivate key (`update-access-key --status Inactive`), (3) attach Deny-All inline policy `BlockAll-On-Leak`. Setelah dijalankan, key yang bocor langsung balas `InvalidAccessKeyId`.
5. **Cleanup wajib:** hapus access key, detach policy, delete user, force-delete bucket canary. CloudTrail + bucket log sengaja dibiarkan hidup sebagai improvement permanen.

---

## Hari 54: Chaos Security Engineering

### Tujuan
Menguji ketahanan sistem peringatan dengan mematikan lapisan keamanan.

### Tutorial

**1. Matikan Lapisan Kunci (Gatekeeper):**
```bash
kubectl scale deployment gatekeeper-controller-manager -n gatekeeper-system --replicas=0
```

**2. Lakukan Instalasi Berbahaya:**
- Deploy pod tanpa resource limits di namespace `securebank` (yang sebelumnya ditolak di Hari 36).
```bash
kubectl apply -f k8s/test-no-limits.yaml
# Harusnya berhasil (Gatekeeper mati)
```

**3. Verifikasi Log Audit / Drift:**
- Jika ada tools drift detection atau CSPM K8s yang jalan berkala, apakah ada alert bahwa pod ini tidak patuh?
- Falco mungkin mencatat pod start up tanpa policy yang benar jika ada rule khusus.

**4. Kembalikan Sistem:**
```bash
kubectl scale deployment gatekeeper-controller-manager -n gatekeeper-system --replicas=1
kubectl delete -f k8s/test-no-limits.yaml
```

### Checklist
- [ ] OPA Gatekeeper sengaja dimatikan.
- [ ] Pod tidak patuh berhasil di-deploy.
- [ ] Memahami pentingnya layer monitoring berlapis (defense in depth).
- [ ] Sistem dikembalikan ke kondisi normal.

### Implementation Notes (Hari 54 — Hasil Praktek)
1. **Admission ≠ Runtime.** Gatekeeper mengecek saat deploy, Falco mengecek syscall saat runtime. Dengan `scale deployment gatekeeper-controller-manager --replicas=0`, admission mati dan klaster menjadi **fail-open** — pod non-compliant langsung diterima (`test-no-limits`, `test-privileged` created).
2. **Runtime hanya melihat primitive, bukan kebijakan.** Pod tanpa limits yang diam → **0 alert Falco**. Pod `privileged: true` → rule `Privileged container launched` fire **CRITICAL** → falcosidekick → webhook → Slack. Runtime tidak "mengerti" resource limits — itu domain detector konfigurasi.
3. **OPA audit tetap hidup** di deployment terpisah dan mencatat pelanggaran di `status.totalViolations` (8 untuk resource limits + 1 untuk privileged). Cek via kind penuh: `kubectl get k8srequiredlimits.constraints.gatekeeper.sh`.
4. **Drift detector buatan sendiri:** `securebank-api/security/chaos/drift-check.sh` — script kubectl+jq untuk list container tanpa `requests/limits` per namespace; exit code 1 jika ada drift (bisa dipakai di CI/cron).
5. **Restore + verifikasi:** scale controller-manager balik ke 1, hapus pod chaos (`--grace-period=0 --force`), re-apply manifest non-compliant harus kembali **Forbidden**. Manifes chaos ada di `k8s/chaos/`.

---

## Hari 55: Pembuatan Laporan Audit

### Tujuan
Mengekspor data dari DefectDojo menjadi laporan profesional.

### Tutorial

**1. Generate Report di DefectDojo:**
- Masuk ke dashboard DefectDojo, navigasi ke Product `SecureBank API`.
- Klik `Findings` -> `Generate Report` (pilih format Executive Summary / PDF jika plugin tersedia, atau HTML).
- Pastikan mencakup statistik SAST, SCA, dan DAST.

**2. Buat Draf Struktur Laporan:**
Buat file `security/audit-reports/draft-q3.md`:
```markdown
# Executive Summary
Aplikasi SecureBank API telah diuji menggunakan metodologi DevSecOps (SAST, SCA, DAST, IaC Scan).
Jumlah kerentanan kritis awal: 12
Jumlah kerentanan kritis tersisa: 0
Tingkat risiko saat ini: LOW
...
```

### Checklist
- [ ] Report DefectDojo berhasil digenerate.
- [ ] Draf laporan awal berisi metrik keamanan yang jelas.

### Implementation Notes (Hari 55 — Hasil Praktek)
1. **Data DefectDojo cepat stale.** Karena DefectDojo lokal, job CI `upload-defectdojo` di-skip (variabel `DEFECTDOJO_URL` kosong) → data terakhir dari import manual Day 47. **Re-sync sebelum laporan**: `bash securebank-api/security/defectdojo/upload-scans.sh` — import menambah test baru per run (dedup tidak menyatukan antar run), catat id test baru sebagai "state terkini".
2. **Regenerasi report scanner** sebelum re-sync supaya angka mencerminkan kenyataan: `trivy fs --scanners vuln`, `trivy config` (terraform & k8s), `trivy image`, `semgrep --config "p/golang"`. Hasil: IaC turun 14 → 2 Low (bukti remediasi Day 23), k8s-post-fix 0, tapi image stdlib muncul **CVE-2026-39822 HIGH** (rebuild image = action item baru).
3. **trivy fs men-scan file gitignored** (cosign.key, aws-credentials.yaml, .env) → report lokal bisa berisi secret. Untuk report yang di-commit, gunakan `--scanners vuln` atau `--skip-dirs`. CI aman karena checkout GitHub tidak memuat file gitignored.
4. **Scope K8s scan ke manifest produksi**: `trivy config --skip-dirs k8s/chaos --skip-files k8s/redteam-pod.yaml k8s/` — supaya report "post-fix" tidak terkotori manifest red team/chaos.
5. **checkov CLI rusak di Python 3.14** (`No module named checkov.__main__`) dan image `bridgecrewio/checkov` sudah tidak ada → report `checkov-report.json` lama di-rename `.stale`; coverage IaC dialihkan ke trivy config.
6. **Draf laporan**: `securebank-api/security/audit-reports/draft-q3.md` — pisahkan metrik *aggregate* (52 findings DefectDojo) dari *state terkini* (6 findings), plus tabel tren sebelum/sesudah dan residual risk.

---

## Hari 56: Penyusunan Dokumen Eksekutif

### Tujuan
Membuat laporan keamanan akhir yang siap dipresentasikan ke manajemen (PDF/Word).

### Tutorial

**1. Lengkapi Dokumen:**
Tambahkan bagian ke dalam laporan:
- **Metodologi:** Penjelasan tools yang dipakai (Semgrep, Trivy, ZAP, Checkov).
- **Key Findings (Before):** Screenshot atau daftar temuan kritis (MD5 usage, S3 Public).
- **Mitigation Evidence (After):** Bukti implementasi (bcrypt, S3 Encrypted, K8s Network Policy).
- **Residual Risk:** Risiko yang masih diterima (misal: "Aplikasi belum memiliki WAF, akan diimplementasi di Q4").

**2. Format dan Export:**
- Konversi `draft-q3.md` ke PDF menggunakan Pandoc atau alat markdown-to-pdf lainnya.

### Checklist
- [ ] Laporan memiliki 4 komponen utama (Summary, Methodology, Findings/Mitigations, Residual Risk).
- [ ] Dokumen diekspor ke PDF.

---

## Hari 57: Review Dokumen Bersama AI

### Tujuan
Memanfaatkan AI untuk menyempurnakan nada bahasa (tone) agar sesuai dengan bahasa eksekutif (CISO/CTO).

### Tutorial

**1. Siapkan Prompt:**
```text
Saya memiliki draf laporan audit keamanan teknis ini. Tolong perbaiki nada bahasanya agar lebih profesional, tidak terlalu teknikal, berfokus pada dampak bisnis (business impact), dan cocok dibaca oleh Chief Information Security Officer (CISO).

[Copy-paste isi draft-q3.md]
```

**2. Terapkan Saran:**
- AI mungkin menyarankan mengganti "Kami mengganti MD5 dengan bcrypt" menjadi "Mengimplementasikan standar kriptografi modern (bcrypt) untuk mitigasi risiko kebocoran data nasabah".
- Update dokumen berdasarkan saran AI.

### Checklist
- [ ] Draf laporan di-review oleh LLM.
- [ ] Nada bahasa berhasil ditingkatkan menjadi lebih manajerial/profesional.

---

## Hari 58: CDP Exam Simulation (Lab Setup)

### Tujuan
Membersihkan environment (Wipe) dan mempersiapkan skenario mirip ujian sertifikasi Certified DevSecOps Professional.

### Tutorial

**1. Hapus Semua Konfigurasi CI/CD & Infra:**
*(Sebaiknya buat branch baru `cdp-exam` dari branch awal aplikasi atau backup yang sudah ada).*
```bash
git checkout -b cdp-exam
rm -rf .github/workflows/*
rm -rf terraform/*
rm -rf k8s/*
rm -rf security/*
git commit -am "chore: wipe config for exam sim"
```

**2. Siapkan "Raw Application":**
Pastikan Anda hanya memiliki kode sumber Go di `cmd/api/main.go` yang mengandung kerentanan (seperti versi di Hari 8).

### Checklist
- [ ] Branch ujian dibuat.
- [ ] Semua konfigurasi CI/CD, K8s, dan Terraform dihapus.
- [ ] Hanya tersisa kode aplikasi mentah.

---

## Hari 59: CDP Exam Simulation (Execution - 3 Jam)

### Tujuan
Menguji kemampuan mengingat dan mengimplementasikan pipeline DevSecOps dari nol tanpa menyalin kode sebelumnya.

### Tutorial

**1. Set Stopwatch 3 Jam.**

**2. Target yang Harus Diselesaikan:**
- Tulis `.github/workflows/ci.yml` dari nol.
- Tambahkan Secret Scan (Gitleaks).
- Tambahkan SCA (Trivy).
- Tambahkan SAST (Semgrep).
- Buat Dockerfile multi-stage.
- Tulis manifest K8s (deployment & service) dengan SecurityContext yang aman.

**3. Aturan:**
- Boleh mencari di Google/Dokumentasi resmi.
- Dilarang mencontek dari branch `main` repositori Anda sendiri.

### Checklist
- [ ] Simulasi selesai dilakukan dalam batas waktu.
- [ ] Pipeline berhasil jalan dan berwarna hijau.
- [ ] Evaluasi diri: Bagian mana yang masih sering lupa syntax-nya?

---

## Hari 60: Project Showcase

### Tujuan
Merayakan keberhasilan dengan publikasi arsitektur "60 Days DevSecOps Mastery".

### Tutorial

**1. Buat Diagram End-to-End:**
- Gunakan Draw.io atau Excalidraw.
- Gambarkan Developer -> Git Push -> GitHub Actions (Gitleaks, Trivy, Semgrep) -> Docker Build (Cosign) -> K8s k3d -> OPA Gatekeeper -> Falco -> n8n -> Slack.

**2. Publikasikan Pencapaian:**
- Tulis artikel di LinkedIn, Medium, atau blog pribadi.
- Judul rekomendasi: *"Bagaimana Saya Membangun Arsitektur DevSecOps End-to-End dalam 60 Hari"*
- Sertakan link repositori GitHub (pastikan repositori sudah publik dan bersih dari rahasia nyata).

**3. Penutup:**
- Gabungkan kembali semua branch ke main (jika perlu).
- Matikan/hancurkan resource lokal/cloud yang memakan biaya.
```bash
k3d cluster delete securebank
docker-compose down
```

### Checklist
- [ ] Diagram End-to-End dibuat.
- [ ] Artikel publikasi dirilis/dipersiapkan.
- [ ] Repositori di-public-kan (opsional).
- [ ] Resource dimatikan.

---

> 🎉 **Selamat! Anda telah menyelesaikan tantangan 60 Hari DevSecOps Mastery.** Anda kini siap mengimplementasikan kultur keamanan di proyek nyata dan menghadapi sertifikasi industri.
