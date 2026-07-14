# Hari 37 — Network Policies

**📅 Tanggal:** 2026-07-14  
**⏱️ Durasi Belajar:** ~50 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat `k8s/network-policy.yaml` — Default Deny All + whitelist
- [x] Apply NetworkPolicy ke cluster
- [x] Test cross-namespace access (dari default namespace)
- [x] Test same-namespace access (dari securebank namespace)
- [x] Verify API pods masih healthy
- [x] Re-run Checkov — CKV2_K8S_6 sekarang PASS (remove skip)
- [x] Cleanup test pods

---

## ✅ Yang Berhasil Dikerjakan

- 3 NetworkPolicy dibuat dan di-apply: default-deny-all, allow-api-ingress, allow-dns-egress
- NetworkPolicy di-enforce oleh k3s flannel (SURPRISE — awalnya expect tidak enforce)
- Cross-namespace access: BLOCKED (HTTP 000 dari default namespace)
- Same-namespace access: BLOCKED (HTTP 000 dari securebank namespace — karena hanya kube-system yang di-whitelist)
- API pods tetap Running 1/1 — kubelet probes bypass NetworkPolicy
- Checkov CKV2_K8S_6: SKIP → PASS (88 passed, 0 failed)
- Test pods cleaned up

---

## 📝 Catatan Teknis

### NetworkPolicy Design: Default Deny + Whitelist

3 policy dalam 1 file:

**1. default-deny-all** — Deny ALL ingress + egress ke semua Pod di namespace `securebank`. `podSelector: {}` = semua Pod. `policyTypes: [Ingress, Egress]` = blok kedua arah. Ini adalah baseline — semua traffic diblokir kecuali yang di-whitelist.

**2. allow-api-ingress** — Allow ingress ke Pod `app: securebank-api` hanya dari namespace `kube-system` (traefik ingress controller), port 8080. Pod dari namespace lain (default, dst) tidak bisa akses API.

**3. allow-dns-egress** — Allow egress DNS (port 53 UDP+TCP) untuk semua Pod di `securebank`. DNS dibutuhkan untuk service discovery (resolve `securebank-api.securebank.svc`).

### Apply & Verify

```bash
$ kubectl apply -f k8s/network-policy.yaml
networkpolicy.networking.k8s.io/default-deny-all created
networkpolicy.networking.k8s.io/allow-api-ingress created
networkpolicy.networking.k8s.io/allow-dns-egress created

$ kubectl get networkpolicy -n securebank
NAME                POD-SELECTOR         AGE
allow-api-ingress   app=securebank-api   3s
allow-dns-egress    <none>               3s
default-deny-all    <none>               3s
```

### Test 1: Cross-Namespace Access (default → securebank) — BLOCKED

```bash
$ kubectl apply -f test-curl-default.yaml  # Pod di namespace default
$ kubectl logs test-curl-default -n default
HTTP 000  # curl exit code 7 = connection refused
```

Pod di `default` namespace tidak bisa akses `securebank-api.securebank.svc:80`. NetworkPolicy blok ingress dari namespace `default` — hanya `kube-system` yang di-whitelist.

### Test 2: Same-Namespace Access (securebank → securebank) — BLOCKED

```bash
$ kubectl apply -f test-curl-ns.yaml  # Pod di namespace securebank
$ kubectl logs test-curl-ns -n securebank
HTTP 000  # curl exit code 7 = connection refused
```

Pod di `securebank` namespace juga tidak bisa akses API! Karena `allow-api-ingress` hanya whitelist dari `kube-system`, bukan dari `securebank` sendiri. Ini design choice: Pod API tidak perlu saling bicara (independent replicas), jadi tidak perlu allow same-namespace traffic.

### Test 3: API Pods Tetap Healthy

```bash
$ kubectl get pods -n securebank -l app=securebank-api
NAME                              READY   STATUS    RESTARTS   AGE
securebank-api-54c957948c-b9bkr   1/1     Running   0          4d2h
securebank-api-54c957948c-swlr4   1/1     Running   0          4d2h
```

Pods tetap Running 1/1. Liveness dan readiness probes tetap work karena kubelet mengakses Pod langsung via node (bypass NetworkPolicy). NetworkPolicy hanya memfilter traffic yang melalui Kubernetes network (CNI).

### Test 4: API via Port-Forward — WORKS

```bash
$ kubectl port-forward svc/securebank-api -n securebank 9081:80 &
$ curl http://localhost:9081/health
{"status":"healthy"}  # HTTP 200
```

Port-forward bypass NetworkPolicy karena traffic masuk melalui API Server (SSH tunnel), bukan melalui CNI network.

### Discovery: k3s Flannel ENFORCE NetworkPolicy

Awalnya expect flannel (k3d default CNI) tidak enforce NetworkPolicy. Tapi ternyata **k3s flannel SUPPORT NetworkPolicy**. k3s menyertakan NetworkPolicy controller yang bekerja dengan flannel. Test membuktikan: cross-namespace access = BLOCKED, same-namespace access = BLOCKED. NetworkPolicy benar-benar di-enforce, bukan hanya object ada di cluster.

### Checkov: CKV2_K8S_6 SKIP → PASS

Sebelum Day 37, Checkov CKV2_K8S_6 (NetworkPolicy) di-skip dengan `--skip-check` karena belum ada NetworkPolicy. Sekarang NetworkPolicy sudah ada:

```bash
$ checkov -d k8s/ --framework kubernetes --skip-check CKV_K8S_15,CKV_K8S_35,CKV_K8S_43

Summary: 88 passed, 0 failed, 0 skipped
CKV2_K8S_6: PASSED on Pod.securebank.securebank-api.app-securebank-api
```

Skip CKV2_K8S_6 dihapus dari command. Sisa 3 skip: CKV_K8S_15 (imagePullPolicy), CKV_K8S_35 (secret as env), CKV_K8S_43 (image digest).

| Metrik | Before (Day 33) | After (Day 37) |
|--------|-----------------|----------------|
| Checkov passed | 85 | 88 |
| Checkov failed | 0 | 0 |
| Checkov skipped | 4 (CLI) | 3 (CLI) |
| CKV2_K8S_6 | SKIP | PASS |
| NetworkPolicy | ❌ None | ✅ 3 policies |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Expect flannel tidak enforce NetworkPolicy | k3s flannel SUPPORT NetworkPolicy. Test membuktikan: traffic benar-benar diblokir. Tidak perlu recreate cluster dengan Calico |
| Test pod ditolak Gatekeeper (no resources) | Buat YAML dengan resources block (requests + limits) sesuai policy Day 35 |
| Cross-namespace curl HTTP 000 | Expected! NetworkPolicy blok ingress dari non-kube-system namespace. Test membuktikan enforcement bekerja |
| Same-namespace curl juga HTTP 000 | Design choice: allow-api-ingress hanya whitelist kube-system, bukan securebank sendiri. Pod API tidak perlu saling bicara |
| localhost:9080 returns 404 | k3d load balancer → traefik → no Ingress resource. Gunakan port-forward untuk test langsung |

---

## 📤 Output Hari Ini

- [x] `k8s/network-policy.yaml` — 3 NetworkPolicy (default-deny-all, allow-api-ingress, allow-dns-egress)
- [x] NetworkPolicy applied dan di-enforce oleh k3s flannel
- [x] Cross-namespace access BLOCKED (verified)
- [x] API pods tetap healthy (kubelet probes bypass NetworkPolicy)
- [x] Checkov CKV2_K8S_6: SKIP → PASS (88/0, 3 skip remaining)
- [x] `security/checkov-k8s-netpol-report.json` — post-NetworkPolicy report
- [x] Test pods cleaned up

---

## 💡 Pelajaran Baru

- **Default deny + whitelist = pola terbaik.** Mulai dengan "block everything", kemudian allow hanya yang dibutuhkan. Lebih aman daripada "allow all + deny specific". Kalau lupa deny sesuatu, tidak ada celah — karena default-nya sudah deny.

- **k3s flannel SUPPORT NetworkPolicy.** Berbeda dari flannel standalone yang tidak enforce NetworkPolicy. k3s menyertakan NetworkPolicy controller yang bekerja dengan flannel. Test membuktikan: traffic benar-benar diblokir, bukan hanya object ada.

- **kubelet probes bypass NetworkPolicy.** Liveness dan readiness probe diakses oleh kubelet langsung via node, bukan melalui CNI network. NetworkPolicy tidak memfilter traffic ini. Pods tetap healthy meski default-deny-all.

- **Port-forward bypass NetworkPolicy.** Traffic dari `kubectl port-forward` masuk melalui API Server (SSH tunnel), bukan CNI network. Berguna untuk debugging tanpa harus bikin hole di NetworkPolicy.

- **NetworkPolicy = namespace-scoped.** Policy di namespace `securebank` hanya memfilter Pod di namespace tersebut. Pod di namespace `default` tidak terkena policy `securebank`. Tapi ingress ke `securebank` dari `default` diblokir karena `allow-api-ingress` hanya whitelist `kube-system`.

---

## 🔗 Referensi

- [Kubernetes Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [NetworkPolicy default-deny pattern](https://kubernetes.io/docs/concepts/services-networking/network-policies/#default-deny-all-inggress-and-all-egress-traffic)
- [k3s NetworkPolicy support](https://docs.k3s.io/networking)
- [Checkov CKV2_K8S_6](https://docs.prismacloud.io/en/enterprise-edition/policy-reference/kubernetes-policies/kubernetes-policy-index/ensure-that-the-default-namespace-has-a-network-policy-defined)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Surprise: k3s flannel enforce NetworkPolicy — lebih baik dari expect |
| Pemahaman materi | 5 | Default-deny pattern, namespace-scoped, kubelet bypass, port-forward bypass |
| Progres sesuai target | 5 | 3 policies, enforcement verified, Checkov CKV2_K8S_6 PASS |

---

## ➡️ Rencana Besok

- [ ] Hari 38: RBAC Auditing — dedicated ServiceAccount with least privilege

---

*[← Hari 36](hari-36.md) | [Hari 38 →](hari-38.md)*
