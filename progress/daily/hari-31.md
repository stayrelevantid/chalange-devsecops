# Hari 31 — K8s Cluster + Deploy

**📅 Tanggal:** 2026-07-08  
**⏱️ Durasi Belajar:** ~1 jam  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat k3d cluster lokal
- [x] Buat K8s manifest files (namespace, secret, deployment, service)
- [x] Import distroless image ke k3d cluster
- [x] Deploy SecureBank API ke K8s (2 replicas)
- [x] Test via port-forward — `/health` dan `/balance` merespons

---

## ✅ Yang Berhasil Dikerjakan

- k3d cluster `securebank` created dengan 2 agents (port 9080 → loadbalancer)
- 4 K8s manifest files: `namespace.yaml`, `secret.yaml`, `deployment.yaml`, `service.yaml`
- K8s Secret `securebank-secrets` untuk `JWT_SECRET` (base64 encoded, best practice)
- Image `securebank:v1` (7.97MB distroless) imported ke cluster via `k3d image import`
- 2 replicas Running, `/health` returns `{"status":"healthy"}`, `/balance` returns 401 without JWT + data with JWT

---

## 📝 Catatan Teknis

### K8s Manifest Files
```
securebank-api/k8s/
├── namespace.yaml     # Namespace: securebank
├── secret.yaml        # Secret: securebank-secrets (JWT_SECRET base64)
├── deployment.yaml    # Deployment: 2 replicas, distroless image, envFrom secret
└── service.yaml       # Service: ClusterIP, port 80 → targetPort 8080
```

### Cluster Setup
```bash
k3d cluster create securebank --port "9080:80@loadbalancer" --agents 2
k3d image import securebank:v1 -c securebank
kubectl apply -f k8s/
```

### Deployment Details

| Property | Value |
|----------|-------|
| Image | securebank:v1 (distroless, 7.97MB) |
| Replicas | 2 |
| Container port | 8080 |
| Service | ClusterIP, port 80 → 8080 |
| Secret | JWT_SECRET via K8s Secret (base64) |
| SecurityContext | **Intentionally missing** (Day 32 scan target) |
| Resources | **Intentionally missing** (Day 32 scan target) |

### Port Mapping
- Port 8080 di host dipakai Docker → k3d loadbalancer pakai port **9080**
- `kubectl port-forward svc/securebank-api 9080:80 -n securebank`
- Host:9080 → k3d LB:80 → Service:80 → Pod:8080

### JWT Secret via K8s Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: securebank-secrets
  namespace: securebank
type: Opaque
data:
  JWT_SECRET: Y2ktdGVzdC1zZWNyZXQ=  # base64("ci-test-secret")
```

Deployment inject via `envFrom`:
```yaml
envFrom:
  - secretRef:
      name: securebank-secrets
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Port 8080 di host dipakai Docker | k3d pakai port 9080 untuk loadbalancer |
| `kubectl apply -f k8s/` gagal: namespace not found | Namespace belum ready saat apply — re-apply setelah 2 detik |
| YAML indentation error di deployment.yaml | Fix `protocol: TCP` indentation (align dengan `containerPort`) |
| JWT token generation untuk test endpoint | Buat temporary Go script dengan golang-jwt/jwt untuk generate token |

---

## 📤 Output Hari Ini

- [x] k3d cluster `securebank` running (1 server + 2 agents)
- [x] 4 K8s manifest files di `securebank-api/k8s/`
- [x] 2 pods Running, `/health` + `/balance` tested
- [x] K8s Secret untuk JWT_SECRET (best practice, bukan plain env var)

---

## 💡 Pelajaran Baru

- **K8s Secret untuk sensitive data.** Daripada hardcode JWT_SECRET di deployment.yaml (visible di `kubectl describe`), pakai K8s Secret dengan base64 encoding. Secret tetap decodeable dengan `kubectl decode`, tapi terpisah dari manifest dan bisa di-manage secara terpisah (Day 43: External Secrets Operator).

- **k3d image import untuk local development.** Tidak perlu push ke registry untuk testing local. `k3d image import` inject image ke semua node di cluster. Cepat dan simple untuk dev.

- **Namespace apply order.** `kubectl apply -f k8s/` apply semua file sekaligus, tapi namespace butuh waktu untuk terbentuk sebelum resource lain bisa dibuat. Solusi: re-apply setelah 2 detik, atau gunakan `kubectl apply` dengan namespace dulu, lalu sisanya.

- **Port conflict handling.** Docker Desktop di macOS bind ke port 8080. k3d loadbalancer butuh port yang free. Pilih 9080 sebagai alternatif. Container port (8080) tetap sama — hanya host mapping yang berubah.

---

## 🔗 Referensi

- [k3d documentation](https://k3d.io/)
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [kubectl port-forward](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Fase 3 dimulai, K8s pertama kalinya |
| Pemahaman materi | 4 | k3d, namespace, secret, deployment, service |
| Progres sesuai target | 5 | 2 replicas running, endpoint tested, semua on track |

---

## ➡️ Rencana Besok

- [ ] Hari 32: K8s Misconfiguration Scan — Kubesec/Checkov scan `deployment.yaml`, temukan 5+ findings

---

*[← Hari 30](hari-30.md) | [Hari 32 →](hari-32.md)*