# 📊 Progress Tracker — 60 Days DevSecOps Mastery

> Folder ini berisi catatan progres harian dan retrospektif dari challenge **SecureBank API** DevSecOps selama 60 hari.

---

## 📁 Struktur Folder

```
progress/
├── README.md               ← File ini (ringkasan & navigasi)
├── tracker.md              ← Master checklist semua 60 hari
├── daily/                  ← Catatan harian (1 file per hari)
│   ├── hari-01.md
│   ├── hari-02.md
│   └── ...
└── retrospektif/           ← Retrospektif per fase
    ├── fase-1-retrospektif.md   (Hari 1–15)
    ├── fase-2-retrospektif.md   (Hari 16–30)
    ├── fase-3-retrospektif.md   (Hari 31–45)
    └── fase-4-retrospektif.md   (Hari 46–60)
```

---

## 🗺️ Panduan Penggunaan

### Cara mengisi progress harian:
1. Buka file `daily/hari-XX.md` sesuai hari yang sedang dikerjakan
2. Isi semua section: **Tasks**, **Catatan**, **Hambatan**, **Output**
3. Update `tracker.md` — centang hari yang sudah selesai

### Cara mengisi retrospektif:
- Setelah menyelesaikan **setiap fase** (15 hari), isi file retrospektif di folder `retrospektif/`
- Format mengikuti template **Start / Stop / Continue** + learning highlights

---

## 📈 Dashboard Progres

| Fase | Rentang | Status | Selesai |
|------|---------|--------|---------|
| 🔐 Fase 1 — Secure SDLC & AppSec | Hari 1–15 | ✅ Selesai | 15/15 |
| 🐳 Fase 2 — IaC & Container Security | Hari 16–30 | 🔄 Berjalan | 7/15 |
| ☸️ Fase 3 — K8s & Runtime Security | Hari 31–45 | ⏳ Menunggu | 0/15 |
| 🔴 Fase 4 — Vuln Mgmt & Red Team | Hari 46–60 | ⏳ Menunggu | 0/15 |
| **Total** | **60 Hari** | | **22/60** |

> Update tabel ini setiap kali satu fase selesai penuh.

---

## 🏆 Target & Milestone

- [ ] **Hari 15** — Fase 1 selesai, pipeline CI/CD dengan Secret/SAST/SCA scan berjalan
- [ ] **Hari 30** — Fase 2 selesai, Docker hardened + Terraform IaC + DAST live
- [ ] **Hari 45** — Fase 3 selesai, K8s cluster dengan OPA + Falco + RBAC
- [ ] **Hari 60** — Fase 4 selesai, laporan audit + CDP exam simulation done

---

*Mulai: 2026-06-05 | Target Selesai: 2026-08-03*
