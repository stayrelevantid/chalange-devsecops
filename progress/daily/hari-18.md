# Hari 18 — Dockerfile Hardening

**📅 Tanggal:** 2026-06-22  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Tambah `COPY --chown=nonroot:nonroot` di Dockerfile
- [x] Buat `docker-compose.yml` dengan security hardening
- [x] Test container dengan `docker compose up -d`
- [x] Verifikasi semua security settings via `docker inspect`

---

## ✅ Yang Berhasil Dikerjakan

- Update Dockerfile: `COPY --chown=nonroot:nonroot` untuk binary dan CA certs
- Buat `docker-compose.yml` dengan 8 layer security hardening:
  1. `user: "65532:65532"` — run as nonroot (UID distroless)
  2. `read_only: true` — root filesystem read-only
  3. `tmpfs: /tmp` — writable /tmp in-memory (64MB)
  4. `security_opt: no-new-privileges:true` — block privilege escalation
  5. `cap_drop: ALL` — drop semua Linux capabilities
  6. `memory: 128M` — memory limit
  7. `cpus: "0.5"` — CPU limit (0.5 core)
  8. `pids: 64` — process count limit
- Test: `docker compose up -d` → `/health` → `{"status":"healthy"}` ✅
- Verifikasi: `docker inspect` menampilkan semua 8 settings

---

## 📝 Catatan Teknis

```bash
# Start container dengan docker compose
$ cd securebank-api
$ docker compose up -d
 Container securebank-api-securebank-1  Started

# Test health endpoint
$ curl -s http://localhost:8080/health
{"status":"healthy"}

# Test auth-protected endpoints (should return 401)
$ curl -s http://localhost:8080/balance
missing authorization header

# Verify all security settings
$ docker inspect --format='User={{.Config.User}} ReadOnly={{.HostConfig.ReadonlyRootfs}} SecurityOpt={{.HostConfig.SecurityOpt}} CapDrop={{.HostConfig.CapDrop}} MemLimit={{.HostConfig.Memory}} NanoCpus={{.HostConfig.NanoCpus}} PidsLimit={{.HostConfig.PidsLimit}} Tmpfs={{.HostConfig.Tmpfs}}' securebank-api-securebank-1

User=65532:65532
ReadOnly=true
SecurityOpt=[no-new-privileges:true]
CapDrop=[ALL]
MemLimit=134217728        ← 128MB
NanoCpus=500000000        ← 0.5 CPU
PidsLimit=64
Tmpfs=map[/tmp:size=64m,mode=1777]

# Dockerfile change
COPY --from=builder --chown=nonroot:nonroot /securebank /securebank
COPY --from=builder --chown=nonroot:nonroot /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Tear down
$ docker compose down
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Docker TLS handshake timeout saat pull `golang:1.26-alpine` | Base image tidak tersedia lokal (sudah di-clean Docker). Gunakan image `securebank:v1` yang sudah di-build dari Day 16. Dockerfile change (`--chown`) akan take effect di CI rebuild |
| Tidak bisa test `--chown` langsung | `--chown` adalah best practice. Image yang ada sudah run as nonroot via `USER nonroot:nonroot`. `--chown` cuma tambahan: pastikan file ownership juga nonroot, bukan root |

---

## 📤 Output Hari Ini

- [x] `securebank-api/docker-compose.yml` — 8 layer security hardening
- [x] `securebank-api/Dockerfile` — `COPY --chown=nonroot:nonroot`
- [x] Container test: health endpoint ✅, auth endpoint ✅
- [x] `docker inspect` verifikasi: 8/8 settings aktif

---

## 💡 Pelajaran Baru

- **`read_only: true` membuat root filesystem tidak bisa ditulis.** Container tidak bisa modify binary, config, atau file sistem. Tapi `/tmp` masih perlu writable, jadi pakai `tmpfs`. Tanpa tmpfs, aplikasi yang butuh write ke /tmp akan crash.

- **`no-new-privileges:true` mencegah privilege escalation.** Bahkan jika attacker menemukan vulnerability (misal SUID binary), mereka tidak bisa naik privilege. Ini kernel-level protection, bukan Docker-specific.

- **`cap_drop: ALL` menghapus semua Linux capabilities.** Default container punya 14 capabilities (CHOWN, DAC_OVERRIDE, NET_BIND_SERVICE, dll). Dropping ALL berarti container cuma bisa apa yang user nonroot bisa lakukan, tanpa special kernel permissions.

- **UID 65532 itu UID nonroot di distroless.** Bukan 1000 atau 1001. `gcr.io/distroless/static-debian12:nonroot` hardcoded UID 65532. Di docker-compose.yml, kita pakai `user: "65532:65532"` (numeric) bukan `user: nonroot` (named) karena lebih reliable.

- **`pids: 64` mencegah fork bomb.** Container hanya bisa spawn 64 process. Kalau attacker coba fork bomb, container akan crash — tapi tidak akan take down host. Ini kernel cgroup feature.

- **HEALTHCHECK di-skip di Dockerfile karena distroless tidak punya shell/wget/curl.** Docker HEALTHCHECK butuh command yang bisa dieksekusi. Distroless hanya punya binary aplikasi. Solusi: K8s liveness probe di Fase 3 (HTTP GET ke `/health`).

---

## 🔗 Referensi

- [Docker Compose Security Hardening](https://docs.docker.com/compose/compose-file/compose-file-v3/#security_opt)
- [Docker Read-Only Filesystem](https://docs.docker.com/engine/reference/run/#read-only-root-filesystem)
- [Linux Capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)
- [Distroless Nonroot UID](https://github.com/GoogleContainerTools/distroless#non-root-images)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Docker compose hardening itu straightforward, 8 layer sekaligus |
| Pemahaman materi | 5 | Tiap setting punya reasoning jelas (kernel-level protection) |
| Progres sesuai target | 5 | Day 18 selesai, Fase 2 lanjut 3/15 |

---

## ➡️ Rencana Besok

- [ ] Hari 19: Image Signing (Cosign) — key pair + signed image di registry

---

*[← Hari 17](hari-17.md) | [Hari 19 →](hari-19.md)*