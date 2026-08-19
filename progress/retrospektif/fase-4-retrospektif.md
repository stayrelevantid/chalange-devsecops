# 🔍 Retrospektif Fase 4 — Vulnerability Management & Red Teaming

**📅 Periode:** Hari 46–60  
**📅 Tanggal Retrospektif:** 2026-08-19  
**⏱️ Total Waktu Belajar:** ~60 hari  
**👤 Ditulis oleh:** Muhammad Indragiri  

---

## 📊 Ringkasan Progres

| Hari | Topik | Status |
|------|-------|--------|
| 46 | DefectDojo Setup | ✅ |
| 47 | DefectDojo API Integration | ✅ |
| 48 | Alert Routing (n8n) | ✅ |
| 49 | AI Remediation Node | ✅ |
| 50 | CSPM (Prowler/ScoutSuite) | ✅ |
| 51 | CSPM Remediation | ✅ |
| 52 | Red Team: K8s Escape | ✅ |
| 53 | Red Team: Leaked Creds | ✅ |
| 54 | Chaos Security Engineering | ✅ |
| 55 | Laporan Audit | ✅ |
| 56 | Dokumen Eksekutif (PDF) | ✅ |
| 57 | AI Review Dokumen | ✅ |
| 58 | CDP Exam Sim: Lab Setup (CI/CD rework) | ✅ |
| 59 | CDP Exam Sim: Preparation Guide | ✅ |
| 60 | Project Showcase | ✅ |

**Selesai: 15/15** ✅

---

## 🟢 Start — Hal yang Ingin Mulai Dilakukan

- **CloudTrail + bucket log** (Day 53) sebagai improvement permanen — menjawab rekomendasi audit Day 51 (log semua akses, termasuk data events S3).
- **Panduan ujian CDP** (Day 59) — strategi 180 menit + cheatsheet + jebakan sebagai template reusable untuk sertifikasi.
- **Showcase akhir** (Day 60) — diagram arsitektur E2E + blog yang menonjolkan "keterampilan yang dikuasai", bukan hanya daftar tool.
- **Disiplin cleanup** — setiap eksperimen cloud/container ditutup dengan dokumentasi teardown (AWS cost → nol).

---

## 🔴 Stop — Hal yang Perlu Dihentikan

- **Bergantung pada satu detector**: Day 54 membuktikan Falco tidak mendeteksi pod tanpa limits yang diam — runtime guard ≠ policy enforcement. Stop berharap satu tool menjawab semua lapisan.
- **Membiarkan layanan "tidak terpakai" tetap jalan**: CloudTrail + secret dikira "gratis", padahal secret $0.40/bln dan bucket versioning menyimpan 18.098 versions. Stop menumpuk residu resource di akun percobaan.
- **Automating bahan bukan inti**: CanaryTokens.org tidak bisa di-automasi — pivot ke cara AWS-native lebih cepat daripada memaksakan tool eksternal.

---

## 🟡 Continue — Hal yang Sudah Baik & Perlu Diteruskan

- **DefectDojo sebagai single pane of glass**: semua scan (Trivy, Semgrep, dll.) terkumpul di satu dashboard → laporan audit lebih cepat disusun.
- **Alerting berlapis**: Falco → n8n → Slack #security-alerts + AI remediation node (LLM auto-ringkas perbaikan) — menutup jeda menunggu developer.
- **CSPM berbasis evidence**: Prowler 328 checks, remediasi account-level (S3 Block Public Access) membuktikan 1 setting = 15 findings fixed.
- **Red team terstruktur**: K8s escape (nsenter) dan leaked credentials dijalankan seperti exercise nyata: serang → deteksi → patch (Falco rules + OPA gap) → verifikasi.
- **CI/CD multi-environment yang stabil**: quality gate 10 check + auto-CD `workflow_run` + merge-commit policy — terbukti stabil sampai hari terakhir.

---

## 🧠 Top 5 Pelajaran Terpenting dari Fase 4

1. **[Detector tak saling menggantikan]** — Falco perhatikan syscall, OPA/Checkov perhatikan konfigurasi. Gap antara keduanya hanya ketahuan lewat chaos engineering (kill webhook → deploy vulnerable app).
2. **[Account-level setting = leverage tinggi]** — 1 config S3 Block Public Access di level account menyelesaikan 15 bucket-level findings sekaligus.
3. **[S3 versioning menyembunyikan residu]** — bucket log CloudTrail mengaktifkan versioning; cleanup wajib hapus versions + delete markers, bukan sekadar `rb --force`.
4. **[Merge commit menjaga traceability]** — squash membuat ancestry divergen dan memicu konflik berulang di promotion; merge commit membuat `develop → staging → main` bisa diaudit.
5. **[Dokumen adalah output profesional]** — laporan audit + PDF eksekutif + panduan ujian membuktikan kemampuan; framework AI review membantu struktur & nada bahasa.

---

## 🏆 Pencapaian Terbaik

- **Laporan audit Q3 resmi** — 6 findings (dari trivy/semgrep/kubesec/checkov/prowler), 0 critical, 1 high (CVE stdlib), lengkap dengan PDF eksekutif siap presentasi.
- **CI/CD 10 quality gate + auto-CD** — smoke test E2E: `feature → develop → staging → main` otomatis men-deploy dev/staging/prod tanpa intervensi.
- **Red team multi-skenario sukses** — K8s escape via nsenter, leaked creds auto-revoke dalam 3 langkah, chaos engineering mengungkap gap monitoring.
- **Kesiapan CDP** — panduan 180 menit lengkap dengan cheatsheet dan 15 best practices terpetakan ke tindakan nyata.
- **Project showcase & cleanup total** — diagram E2E, blog skill showcase, AWS cost kembali nol, klaster lokal dimatikan.

---

## 😓 Tantangan Terbesar

| Tantangan | Cara Mengatasi |
|-----------|----------------|
| Falco di k3d offline → hanya 26 rules, rule escape tidak ada | Diagnosa falcoctl artifact download gagal; tambah custom rules manual + verifikasi di log |
| S3 data events CloudTrail lambat (hingga 15 menit) | Fallback baca file log mentah di bucket + `put-event-selectors` |
| CanaryTokens.org tidak bisa di-automasi (React SPA) | Pivot ke S3 honeypot bucket + marker string |
| Bucket log ber-account ID nyata tidak boleh bocor ke repo | Direntang sebagai `XXXX` di dokumentasi |
| Runner GitHub tanpa cluster → `kubectl apply --dry-run` gagal | Kubeconform `-strict -ignore-missing-schemas` |
| AWS bucket versioning → delete gagal `BucketNotEmpty` | Loop `list-object-versions` → `delete-objects` + delete markers, baru `rb --force` |

---

## 🎯 Kesiapan CDP (Certified DevSecOps Professional)

> *Setelah 60 hari, seberapa siap untuk mengambil sertifikasi?*

| Area | Kesiapan (1–10) | Gap yang Perlu Diisi |
|------|----------------|----------------------|
| Application Security (SAST/SCA/Secret Scan) | 9 | Menghafal sintaks tanpa melihat repo |
| Container Security | 8 | Distroless vs scratch trade-off |
| Kubernetes Security | 8 | Policy Rego yang lebih kompleks |
| IaC Security | 7 | Terraform plan/apply end-to-end dari nol |
| Vulnerability Management | 9 | Automasi DefectDojo lintas tool |
| Red Team Fundamentals | 7 | Eksploitasi yang lebih dalam |
| Laporan & Dokumentasi Audit | 9 | Presentasi lisan di depan stakeholder |

**Keputusan:** [ ] Siap ambil CDP · [x] Perlu review area tertentu dulu (IaC Terraform, Rego, red team) · [ ] Perlu latihan lebih

---

## 🏁 Refleksi 60 Hari Penuh

> *Apa yang paling berubah dalam cara berpikir tentang keamanan software setelah 60 hari ini?*

Pergeseran terbesar: **dari "cari CVE" ke "bangun proses yang mencegah dan merespons"**. Di awal, keamanan = menjalankan scanner. Di akhir, keamanan = pipeline quality gate yang gagal build, policy yang menolak deployment, runtime yang mengirim alert ke Slack, dan laporan yang bisa dipertanggungjawabkan. 60 hari mengajarkan bahwa keamanan adalah sistem yang harus *selalu dites* — bahkan terhadap dirinya sendiri (chaos engineering, red team).

---

## 📈 Skor Diri Final (Jujur!)

| Aspek | Awal (1–10) | Akhir (1–10) | Delta |
|-------|------------|-------------|-------|
| Pemahaman DevSecOps secara keseluruhan | 3 | 8 | +5 |
| Kemampuan praktis (tools & commands) | 2 | 8 | +6 |
| Kemampuan membaca & menulis pipeline | 2 | 8 | +6 |
| Kemampuan audit & dokumentasi | 3 | 8 | +5 |
| Konsistensi belajar | 5 | 9 | +4 |

---

## 🚀 Langkah Selanjutnya

> *Apa yang akan dilakukan setelah challenge 60 hari ini selesai?*

- [ ] Menjalankan simulasi CDP 3 jam memakai `docs/cdp-exam-guide.md` untuk menutup gap (Terraform, Rego, red team).
- [ ] Publikasikan blog + LinkedIn showcase dari draft `blogpost.md`.
- [ ] Menerapkan quality gate + auto-CD pattern ini ke project kerja nyata.

---

## 📝 Catatan Bebas

Challenge selesai 60/60, Fase 4 ditutup dengan **cost AWS nol** dan klaster lokal dimatikan. Artefak lengkap: diagram arsitektur E2E (`docs/architecture/`), laporan audit PDF, panduan CDP, 60 catatan harian, dan retrospektif per fase. Semua evidence tersedia di repositori PUBLIC: [github.com/stayrelevantid/chalange-devsecops](https://github.com/stayrelevantid/chalange-devsecops).

---

*[← Retrospektif Fase 3](fase-3-retrospektif.md) | [Tracker →](../tracker.md)*