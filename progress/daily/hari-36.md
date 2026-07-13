# Hari 36 — OPA Policy Testing

**📅 Tanggal:** 2026-07-13  
**⏱️ Durasi Belajar:** ~45 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Test Pod tanpa resources → expect DITOLAK oleh Gatekeeper
- [x] Test Pod dengan resources → expect DITERIMA
- [x] Test Pod dengan partial resources (limits only, no requests) → expect DITERIMA (K8s auto-fill requests)
- [x] Verify deployment.yaml tetap OK
- [x] Cleanup test pods

---

## ✅ Yang Berhasil Dikerjakan

- 3 test cases dijalankan, semua sesuai ekspektasi (dengan 1 discovery menarik)
- Test 1: Pod tanpa resources → **DITOLAK** (4 violation messages: no CPU/memory limit + request)
- Test 2: Pod dengan full resources → **DITERIMA**
- Test 3: Pod dengan partial resources (limits only) → **DITERIMA** (K8s auto-fill requests = limits)
- Deployment.yaml re-apply → **unchanged** (no error, Gatekeeper pass)
- Test pods cleaned up, temp files deleted

---

## 📝 Catatan Teknis

### Test Results Summary

| # | Test | Resources | Expected | Actual | Result |
|---|------|-----------|----------|--------|--------|
| 1 | test-no-limits | ❌ none | DENIED | DENIED | ✅ Match |
| 2 | test-with-limits | ✅ full (limits + requests) | ACCEPTED | ACCEPTED | ✅ Match |
| 3 | test-partial-limits | ⚠️ limits only | DENIED | ACCEPTED | ⚠️ K8s auto-fill |
| 4 | deployment.yaml | ✅ full (Day 33) | ACCEPTED | ACCEPTED (unchanged) | ✅ Match |

### Test 1: Pod Without Resources — DENIED

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-no-limits
  namespace: securebank
spec:
  containers:
    - name: nginx
      image: nginx:latest
      # NO resources defined
```

```bash
$ kubectl apply -f test-no-limits.yaml
Error from server (Forbidden): error when creating "test-no-limits.yaml":
admission webhook "validation.gatekeeper.sh" denied the request:
[require-resource-limits] Container 'nginx' has no CPU limit
[require-resource-limits] Container 'nginx' has no CPU request
[require-resource-limits] Container 'nginx' has no memory limit
[require-resource-limits] Container 'nginx' has no memory request
```

**Analysis: 4 violation messages** — sesuai dengan 4 violation rules di Rego policy. Gatekeeper mengevaluasi semua rules dan mengumpulkan semua violations dalam satu response.

### Test 2: Pod With Full Resources — ACCEPTED

```yaml
resources:
  requests:
    cpu: 50m
    memory: 64Mi
  limits:
    cpu: 200m
    memory: 128Mi
```

```bash
$ kubectl apply -f test-with-limits.yaml
pod/test-with-limits created
```

Policy pass — Pod dengan full resources (limits + requests) diterima.

### Test 3: Pod With Partial Resources — ACCEPTED (Surprise!)

```yaml
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  # NO requests defined
```

```bash
$ kubectl apply -f test-partial-limits.yaml
pod/test-partial-limits created  # DITERIMA, padahal expect DITOLAK!
```

**Kenapa diterima?** Kubernetes otomatis mengisi `requests` sama dengan `limits` jika `requests` tidak ditentukan. Setelah Pod dibuat, `kubectl get pod -o yaml` menunjukkan:

```yaml
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 200m      # K8s auto-fill dari limits
    memory: 128Mi   # K8s auto-fill dari limits
```

Gatekeeper mengevaluasi **admission request** dari API Server — dan API Server sudah auto-fill requests sebelum webhook dipanggil. Jadi dari POV Gatekeeper, requests sudah ada.

**Discovery:** Ini bukan bug di policy — ini Kubernetes default behavior. Gatekeeper webhook dipanggil setelah K8s defaulting. Untuk mencegah auto-fill, perlu `ValidatingAdmissionPolicy` (K8s 1.30+) yang dipanggil sebelum defaulting.

### Test 4: Deployment.yaml — ACCEPTED

```bash
$ kubectl apply -f k8s/deployment.yaml
deployment.apps/securebank-api unchanged
```

Deployment.yaml yang sudah hardened (Day 33) dengan full resources tetap pass. Tidak ada error.

### Gatekeeper Evaluation Flow (Confirmed)

```
Test 1: kubectl apply (no resources)
  → API Server → K8s defaulting (no requests to fill) → Gatekeeper webhook
  → 4 violations → DENY ✅

Test 2: kubectl apply (full resources)
  → API Server → K8s defaulting (nothing to fill) → Gatekeeper webhook
  → 0 violations → ALLOW ✅

Test 3: kubectl apply (limits only)
  → API Server → K8s defaulting (auto-fill requests = limits) → Gatekeeper webhook
  → 0 violations (requests already filled) → ALLOW ✅ (surprise but correct)
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Test 3 expect DENIED tapi DITERIMA | K8s auto-fill requests = limits sebelum webhook dipanggil. Bukan bug — Kubernetes default behavior. Documented di progress |
| Tutorial tidak mention auto-fill behavior | Discovery saat testing. Penting untuk dipahami: Gatekeeper webhook dipanggil SETELAH K8s defaulting |

---

## 📤 Output Hari Ini

- [x] 3 test cases executed, results documented
- [x] Gatekeeper policy confirmed working (deny without resources)
- [x] Discovery: K8s auto-fill requests from limits
- [x] Test pods cleaned up

---

## 💡 Pelajaran Baru

- **Gatekeeper webhook dipanggil SETELAH K8s defaulting.** Kubernetes auto-fill `requests = limits` kalau `requests` tidak ditentukan. Gatekeeper melihat Pod yang sudah di-default, bukan YAML original. Untuk mencegah auto-fill, perlu `ValidatingAdmissionPolicy` (K8s 1.30+) yang dipanggil sebelum defaulting.

- **4 violation messages dalam satu response.** Gatekeeper mengumpulkan semua violations dari semua rules, bukan stop di rule pertama. Developer dapat lihat semua masalah sekaligus — lebih efficient daripada fix satu-satu.

- **Testing = discovery.** Kalau tidak test Test 3 (partial resources), kita tidak akan tahu bahwa K8s auto-fill requests dari limits. Testing membuktikan assumption — dan kadang assumption salah.

- **Policy enforcement works.** Pod tanpa resources = ditolak sebelum masuk cluster. Ini prevention, bukan detection. Scanner (Checkov/Trivy) melaporkan, Gatekeeper mencegah.

---

## 🔗 Referensi

- [Gatekeeper Policy Testing](https://open-policy-agent.github.io/gatekeeper/website/docs/howto/)
- [Kubernetes Resource Defaulting](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits)
- [Admission Webhook Order](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Testing = discovery, K8s auto-fill behavior |
| Pemahaman materi | 5 | Gatekeeper webhook order, K8s defaulting, Rego evaluation |
| Progres sesuai target | 5 | 3 tests, policy confirmed working |

---

## ➡️ Rencana Besok

- [ ] Hari 37: Network Policies — default deny all + whitelist traffic

---

*[← Hari 35](hari-35.md) | [Hari 37 →](hari-37.md)*