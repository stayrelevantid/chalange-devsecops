# Hari 46 — DefectDojo Setup

**📅 Tanggal:** 2026-07-22  
**⏱️ Durasi Belajar:** ~45 menit  
**🏷️ Fase:** Fase 4 — Vuln Mgmt & Red Team  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Deploy DefectDojo via Docker Compose (standalone, no clone needed)
- [x] Verifikasi DefectDojo accessible (HTTP 200)
- [x] Get admin password dari initializer logs
- [x] Document setup — Product Type, Product, Engagement plan

---

## ✅ Yang Berhasil Dikerjakan

- `docker-compose.yml` standalone dibuat (no clone 358MB repo, no build from source)
- 7 containers Running: nginx, uwsgi, celerybeat, celeryworker, initializer, postgres, valkey
- DefectDojo accessible di `http://localhost:8088` (HTTP 302 → login page HTTP 200)
- Admin password didapat dari initializer logs
- Images sudah di-pull user sebelumnya (4 images, ~1.3GB total)
- Port 8088 dipakai (port 8080 occupied oleh local API)

---

## 📝 Catatan Teknis

### Approach: Standalone Docker Compose (NO CLONE)

Tutorial minta `git clone https://github.com/DefectDojo/django-DefectDojo.git` (358MB) + `./dc-build.sh` (build from source). Slow network kita sering timeout (n8n 400MB timed out Day 42).

**Pivot:** Download hanya `docker-compose.yml` content via curl, create standalone version:
- Hapus `build:` sections (images sudah di-pull, tidak perlu build)
- Hapus `@sha256:...` image digests (pakai tag saja)
- Change port 8080 → 8088 (`published: 8088`)
- Create dummy `./docker/extra_settings/` directory untuk bind mount

**Hasil:** No 358MB clone, no build from source, images sudah ready. `docker compose up -d` = langsung jalan.

### Docker Compose Services (7 containers)

| Service | Image | Port | Function |
|---------|-------|------|----------|
| nginx | defectdojo/defectdojo-nginx:latest | **8088** → 8080 | Reverse proxy |
| uwsgi | defectdojo/defectdojo-django:latest | 3031 (internal) | Django app server |
| celerybeat | defectdojo/defectdojo-django:latest | — | Task scheduler |
| celeryworker | defectdojo/defectdojo-django:latest | — | Async task processing |
| initializer | defectdojo/defectdojo-django:latest | — | DB migration + admin user (exits after done) |
| postgres | postgres:18.4-alpine | 5432 (internal) | Database |
| valkey | valkey/valkey:9.1.0-alpine | 6379 (internal) | Cache (Redis-compatible) |

### Docker Images (pre-pulled by user)

| Image | Size |
|-------|------|
| defectdojo/defectdojo-django:latest | 709MB |
| defectdojo/defectdojo-nginx:latest | 280MB |
| postgres:18.4-alpine | 298MB |
| valkey/valkey:9.1.0-alpine | 45MB |
| **Total** | **~1.3GB** |

### Startup Sequence

```
1. postgres + valkey start (database + cache)
2. initializer starts → waits for postgres → runs DB migrations → creates admin user → exits
3. uwsgi + celerybeat + celeryworker start (wait for initializer to complete)
4. nginx starts (reverse proxy to uwsgi)
5. DefectDojo accessible di http://localhost:8088
```

### Admin Credentials

```
URL: http://localhost:8088
Username: admin
Password: 6BBFxQtWdmFyl74dHACetD
```

Password di-generate random oleh initializer. Setiap fresh setup = password berbeda.

### DefectDojo Setup Plan (UI — untuk user)

Setup berikut dilakukan via browser di `http://localhost:8088`:

1. **Login** — `admin` + password dari logs
2. **Product Type** — `Fintech` (Settings → Product Type → Add)
3. **Product** — `SecureBank API` (Product Type: Fintech)
4. **Engagement** — `Q3 Security Audit` (Product: SecureBank API)

Setelah ini, Day 47 akan upload scan reports (Trivy, Semgrep) ke DefectDojo via API.

### Port Conflict Resolution

```yaml
# Original (port 8080 — conflict dengan local API):
ports:
  - target: 8080
    published: 8080

# Fixed (port 8088):
ports:
  - target: 8080
    published: 8088
```

### File Location

```
securebank-api/security/defectdojo/
├── docker-compose.yml          # Standalone DefectDojo compose
└── docker/
    └── extra_settings/
        └── .gitkeep            # Dummy dir for bind mount
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Git clone 358MB (slow network timeout) | Tidak clone! Download hanya docker-compose.yml content, create standalone version |
| `./dc-build.sh` builds from source (slow) | Hapus `build:` sections, pakai pre-built images dari Docker Hub |
| Port 8080 conflict dengan local API | Change `published: 8088` di nginx ports |
| `./docker/extra_settings` bind mount butuh directory | Create dummy directory dengan `.gitkeep` |
| Docker image pull (1.3GB total) | User bantu pre-pull 4 images sebelum eksekusi |

---

## 📤 Output Hari Ini

- [x] `security/defectdojo/docker-compose.yml` — standalone DefectDojo compose (committed)
- [x] `security/defectdojo/docker/extra_settings/.gitkeep` — dummy dir (committed)
- [x] 7 containers Running (6 active + 1 initializer exited after DB setup)
- [x] DefectDojo accessible di http://localhost:8088 (HTTP 200)
- [x] Admin password obtained: `6BBFxQtWdmFyl74dHACetD`
- [x] Setup plan documented (Product Type → Product → Engagement)

---

## 💡 Pelajaran Baru

- **No clone needed untuk Docker Compose setups.** Tutorial minta clone 358MB repo, tapi sebenarnya cuma butuh `docker-compose.yml` file. Download content via curl, create standalone version. Pragmatic DevSecOps = problem solving (lagi).

- **Pre-built images > build from source.** `./dc-build.sh` builds images from Dockerfile (slow, butuh source code). Pre-built images dari Docker Hub = tinggal pull dan run. Hapus `build:` sections, pakai `image:` saja.

- **Port conflict resolution.** Docker Compose `published` port = host port. `target` = container port. Change `published: 8088` untuk avoid conflict, `target: 8080` tetap (container internal).

- **DefectDojo = single pane of glass.** Vulnerability management platform yang bisa import scan results dari multiple tools (Trivy, Semgrep, Checkov, ZAP, Falco). Day 47 akan upload scan reports via API. Hierarchy: Product Type → Product → Engagement → Test → Findings.

- **Initializer pattern.** DefectDojo pakai `initializer` container yang run-once (DB migration + admin user), kemudian exit. Container lain `depends_on: initializer: condition: service_completed_successfully` — wait sampai initializer selesai sebelum start.

---

## 🔗 Referensi

- [DefectDojo Documentation](https://documentation.defectdojo.com/)
- [DefectDojo GitHub](https://github.com/DefectDojo/django-DefectDojo)
- [DefectDojo Docker Compose](https://github.com/DefectDojo/django-DefectDojo#docker-compose)
- [DefectDojo API v2](https://documentation.defectdojo.com/getting_started/api/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Fase 4 dimulai! DefectDojo jalan tanpa clone |
| Pemahaman materi | 4 | Docker Compose multi-container, initializer pattern |
| Progres sesuai target | 5 | 7 containers Running, HTTP 200, admin password obtained |

---

## ➡️ Rencana Besok

- [ ] Hari 47: DefectDojo API Integration — upload Trivy + Semgrep JSON ke DefectDojo via API

---

*[← Hari 45](hari-45.md) | [Hari 47 →](hari-47.md)*
