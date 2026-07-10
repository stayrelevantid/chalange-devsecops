# Hari 34 — OPA Gatekeeper Setup

**📅 Tanggal:** 2026-07-11  
**⏱️ Durasi Belajar:** ~45 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install OPA Gatekeeper sebagai Admission Controller di k3d cluster
- [x] Verifikasi pods running, CRDs terinstall, validating webhook aktif
- [x] Pahami arsitektur Gatekeeper (audit vs enforcement)

---

## ✅ Yang Berhasil Dikerjakan

- Helm repo `gatekeeper` ditambahkan
- Gatekeeper installed via Helm ke namespace `gatekeeper-system` (replicas=1 untuk dev)
- 2 pods Running: `gatekeeper-audit` + `gatekeeper-controller-manager`
- 17 CRDs terinstall (utama: `constrainttemplates.templates.gatekeeper.sh`)
- ValidatingWebhookConfiguration aktif: `gatekeeper-validating-webhook-configuration`
- Service `gatekeeper-webhook-service` (ClusterIP, port 443)

---

## 📝 Catatan Teknis

### Install Commands

```bash
# 1. Add Helm repo
helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
helm repo update

# 2. Install Gatekeeper
helm install gatekeeper gatekeeper/gatekeeper \
  --namespace gatekeeper-system \
  --create-namespace \
  --set replicas=1

# 3. Verifikasi
kubectl get pods -n gatekeeper-system
kubectl get crd | grep gatekeeper
kubectl get validatingwebhookconfigurations | grep gatekeeper
```

### Verification Output

**Pods:**
```
NAME                                             READY   STATUS    RESTARTS   AGE
gatekeeper-audit-db77bbf66-4b86j                 1/1     Running   0          86s
gatekeeper-controller-manager-59bcd49475-vdqzb   1/1     Running   0          86s
```

**Resources di gatekeeper-system:**
```
NAME                                                 READY   STATUS
pod/gatekeeper-audit-db77bbf66-4b86j                 1/1     Running
pod/gatekeeper-controller-manager-59bcd49475-vdqzb   1/1     Running
service/gatekeeper-webhook-service                   ClusterIP   10.43.109.224   443/TCP
deployment.apps/gatekeeper-audit                     1/1     1/1     1       86s
deployment.apps/gatekeeper-controller-manager        1/1     1/1     1       86s
```

**CRDs (17 total):**
| CRD | Fungsi |
|-----|--------|
| `constrainttemplates.templates.gatekeeper.sh` | Template untuk custom policy (Rego) |
| `configs.config.gatekeeper.sh` | Konfigurasi Gatekeeper |
| `providers.externaldata.gatekeeper.sh` | External data provider |
| `assign.mutations.gatekeeper.sh` | Mutation: assign values |
| `assignimage.mutations.gatekeeper.sh` | Mutation: modify image |
| `assignmetadata.mutations.gatekeeper.sh` | Mutation: modify metadata |
| `modifyset.mutations.gatekeeper.sh` | Mutation: modify sets |
| `syncsets.syncset.gatekeeper.sh` | Sync data dari cluster |
| `connections.connection.gatekeeper.sh` | External connections |
| `expansiontemplate.expansion.gatekeeper.sh` | Template expansion |
| `*.status.gatekeeper.sh` (6 CRDs) | Status tracking (pod, config, constraint, etc) |

**ValidatingWebhookConfiguration:**
```
gatekeeper-validating-webhook-configuration   2
```
Angka `2` = 2 webhook rules (biasanya: `apiGroups: ["*"]` untuk catch-all).

### Arsitektur Gatekeeper

```
┌──────────────────────────────────────────────────┐
│                  k3d cluster                       │
│                                                    │
│  ┌──────────────┐    ┌──────────────────────┐     │
│  │   API Server  │───→│  ValidatingWebhook   │     │
│  │  (kubectl)    │    │  (Gatekeeper)         │     │
│  └──────────────┘    └───────────┬──────────┘     │
│                                   │                 │
│                    ┌──────────────┴──────────┐     │
│                    │                         │     │
│           ┌────────▼─────────┐  ┌───────────▼───┐ │
│           │  Controller       │  │  Audit        │ │
│           │  Manager          │  │  (periodic    │ │
│           │  (enforcement:    │  │   scan cis    │ │
│           │   block on apply) │  │   violations) │ │
│           └──────────────────┘  └───────────────┘ │
│                                                    │
│  Namespace: gatekeeper-system                      │
└──────────────────────────────────────────────────┘
```

**Dua komponen utama:**

1. **Controller Manager** (enforcement) — Validating Admission Webhook. Setiap kali ada `kubectl apply`, API Server mengirim request ke Gatekeeper webhook. Gatekeeper evaluasi policy Rego. Kalau violate → request ditolak (403). Kalau pass → resource dibuat.

2. **Audit** (monitoring) — Periodically scan semua resource yang sudah ada di cluster. Kalau ada yang violate constraint → report via `status` field (tidak block, hanya audit). Berguna untuk "what if" sebelum enable enforcement.

### Before vs After: Cluster tanpa vs dengan Gatekeeper

| Aspek | Before (Day 1-33) | After (Day 34) |
|-------|-------------------|----------------|
| Admission control | ❌ Tidak ada policy enforcement | ✅ Gatekeeper webhook aktif |
| Policy language | ❌ Tidak ada | ✅ Rego policy (Day 35) |
| Block non-compliant | ❌ Semua resource bisa di-apply | ✅ Resource yang violate akan ditolak |
| Audit existing | ❌ Tidak ada | ✅ Audit scan periodic |
| Namespace | — | `gatekeeper-system` (isolasidal) |
| Pods | — | 2 pods (audit + controller) |
| CRDs | — | 17 CRDs terinstall |
| Webhook | — | `gatekeeper-validating-webhook-configuration` |

### Helm Install Detail

```bash
helm install gatekeeper gatekeeper/gatekeeper \
  --namespace gatekeeper-system \
  --create-namespace \
  --set replicas=1
```

- `--namespace gatekeeper-system` → dedicated namespace (best practice)
- `--create-namespace` → bikin namespace kalau belum ada
- `--set replicas=1` → dev cluster (production: 3 untuk HA)

### Koneksi ke Hari Berikutnya

| Hari | Task | Detail |
|------|------|--------|
| **Day 35** | Rego Policy Writing | Tulis ConstraintTemplate + Constraint (e.g., wajib resource limits) |
| **Day 36** | Policy Testing | Test: apply Pod yang violate → expect ditolak. Apply Pod compliant → expect diterima |

Hari ini install "mesin"-nya. Besok tulis "aturan"-nya. Lusa test "aturan"-nya bekerja.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `sealed-secrets` chart repo 404 di `helm repo update` | Tidak related — repo lama yang sudah tidak di-maintain. Ignore, tidak mengganggu gatekeeper |
| `retool` chart repo 503 | Sama — repo lain yang temporary down. Tidak mengganggu |

---

## 📤 Output Hari Ini

- [x] OPA Gatekeeper installed via Helm di `gatekeeper-system` namespace
- [x] 2 pods Running (audit + controller-manager)
- [x] 17 CRDs terinstall
- [x] ValidatingWebhookConfiguration aktif
- [x] Gatekeeper siap untuk Day 35 (Rego policy writing)

---

## 💡 Pelajaran Baru

- **Admission Controller = "bouncer di pintu masuk cluster".** Tanpa Gatekeeper, siapa saja bisa `kubectl apply` resource apa saja — selama YAML valid, K8s terima. Dengan Gatekeeper, ada policy check di "pintu masuk": kalau violate, request ditolak sebelum resource dibuat.

- **Dua mode: enforcement vs audit.** Controller Manager = enforcement (block di apply). Audit = monitoring (scan existing, report without blocking). Kombinasi keduanya = defense in depth: cegah yang baru, audit yang lama.

- **Rego = policy language OPA.** Bukan YAML, bukan JSON. Bahasa declarative untuk express policy. Besok (Day 35) akan tulis Rego policy pertama: "Pod tanpa resource limits ditolak".

- **CRD = Custom Resource Definition.** Gatekeeper install 17 CRD — artinya K8s sekarang mengerti 17 jenis resource baru. Yang utama: `ConstraintTemplate` (template policy) dan `Constraint` (instance dari template).

---

## 🔗 Referensi

- [OPA Gatekeeper docs](https://open-policy-agent.github.io/gatekeeper/)
- [Gatekeeper Helm chart](https://github.com/open-policy-agent/gatekeeper/tree/master/charts/gatekeeper)
- [Kubernetes Admission Controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Install smooth, tinggal 2 menit |
| Pemahaman materi | 4 | Arsitektur Gatekeeper, enforcement vs audit |
| Progres sesuai target | 5 | Install + verify selesai, siap untuk Day 35 |

---

## ➡️ Rencana Besok

- [ ] Hari 35: Rego Policy — ConstraintTemplate + Constraint untuk wajib resource limits

---

*[← Hari 33](hari-33.md) | [Hari 35 →](hari-35.md)*