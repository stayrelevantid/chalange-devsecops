# Hari 16 — Dockerfile Multi-stage Build

**📅 Tanggal:** 2026-06-20  
**⏱️ Durasi Belajar:** 1.5 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat Dockerfile multi-stage build (build di alpine, run di distroless)
- [x] Binary di-strip dengan `-ldflags="-w -s"`
- [x] Image < 15MB
- [x] Aplikasi berjalan normal di container
- [x] `.dockerignore` dikonfigurasi

---

## ✅ Yang Berhasil Dikerjakan

- Membuat `securebank-api/Dockerfile` dengan multi-stage build:
  - Stage 1 (builder): `golang:1.26-alpine` — build binary dengan `CGO_ENABLED=0` dan `-ldflags="-w -s"`
  - Stage 2 (runtime): `gcr.io/distroless/static-debian12:nonroot` — copy binary + CA certs only
- Membuat `securebank-api/.dockerignore` — exclude `.git`, `.github`, `docs`, `security`, `*.md`, test files
- Build image: `docker build -t securebank:v1 .`
- Test container: `curl http://localhost:8080/health` → `{"status":"healthy"}`
- Image size: **7.97MB** (target < 15MB ✅)

---

## 📝 Catatan Teknis

```dockerfile
# Stage 1: Build
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /securebank ./cmd/api

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /securebank /securebank
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/securebank"]
```

```bash
# Build & test
docker build -t securebank:v1 .
docker run -d -p 8080:8080 -e JWT_SECRET=test-secret securebank:v1
curl http://localhost:8080/health
# {"status":"healthy"}

# Image size
docker images securebank:v1
# securebank:v1   7.97MB
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `golang:1.26-alpine` image belum cached | First build memakan waktu ~30s buat download, subsequent build akan lebih cepat karena Docker layer cache |
| Distroless image tidak punya shell | Tidak bisa `docker exec -it` untuk debugging — ini security feature, bukan bug |
| `JWT_SECRET` perlu di-pass saat run | `docker run -e JWT_SECRET=test-secret` — di production akan pakai K8s Secrets (Fase 3) |

---

## 📤 Output Hari Ini

- [x] `securebank-api/Dockerfile` — multi-stage build (alpine → distroless)
- [x] `securebank-api/.dockerignore` — exclude unnecessary files
- [x] Image size 7.97MB (< 15MB target ✅)
- [x] Container test passed: `/health` returns `{"status":"healthy"}`
- [x] Commit: `feat: add multi-stage Dockerfile with distroless runtime (Day 16)`

---

## 💡 Pelajaran Baru

- **Multi-stage build itu fundamental untuk container security.** Stage 1 (builder) bawa toolchain lengkap (~300MB), tapi hanya binary final yang di-copy ke stage 2 (runtime). Attack surface drastis berkurang.

- **Distroless image tidak punya shell.** Ini feature, bukan bug. Tanpa shell, attacker yang masuk ke container tidak bisa menjalankan command. Tapi artinya debugging harus pakai `docker cp` atau `kubectl exec` dengan debug container.

- **`-ldflags="-w -s"` menghapus DWARF debug info dan symbol table.** Ini mengurangi binary size ~20-30% dan juga mengurangi informasi yang bisa di-extract oleh attacker.

- **`CGO_ENABLED=0` menghasilkan statically-linked binary.** Ini penting karena distroless tidak punya glibc. Tanpa flag ini, binary akan dinamically link ke glibc dan crash di distroless.

- **7.97MB itu sangat kecil.** Bandingkan: alpine base image ~7MB, golang:alpine ~300MB. Distroless + stripped binary = attack surface minimal.

---

## 🔗 Referensi

- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Go Binary Optimization](https://pkg.go.dev/cmd/link)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Fase 2 dimulai! Container security itu exciting |
| Pemahaman materi | 5 | Multi-stage build dan distroless jadi clear |
| Progres sesuai target | 5 | Image 7.97MB, jauh di bawah target 15MB |

---

## ➡️ Rencana Besok

- [ ] Hari 17: Container Image Scan — `trivy image securebank:v1` dan bandingkan dengan naive image

---

*[← Hari 15](hari-15.md) | [Hari 17 →](hari-17.md)*