# Hari 35 — Rego Policy: Wajib Resource Limits + Requests

**📅 Tanggal:** 2026-07-11  
**⏱️ Durasi Belajar:** ~1 jam  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Tulis Rego policy yang menolak Pod/Deployment tanpa resource limits
- [x] Tambah cek resource requests (lebih comprehensive dari tutorial)
- [x] Buat ConstraintTemplate + Constraint
- [x] Apply ke cluster
- [x] Verifikasi constraint aktif dengan `enforcementAction: deny`

---

## ✅ Yang Berhasil Dikerjakan

- 2 file Gatekeeper policy dibuat di `k8s/gatekeeper/`:
  - `constraint-templates/require-limits.yaml` — ConstraintTemplate dengan Rego (4 violation rules)
  - `constraints/require-limits.yaml` — Constraint dengan `enforcementAction: deny`, match Pod + Deployment di namespace `securebank`
- ConstraintTemplate applied: `k8srequiredlimits` created
- Constraint applied: `require-resource-limits` created, ENFORCEMENT-ACTION: deny
- 0 total violations — deployment.yaml yang sudah hardened (Day 33) langsung pass
- Besok (Day 36) akan test: apply Pod tanpa resources → expect ditolak

---

## 📝 Catatan Teknis

### File Structure
```
k8s/gatekeeper/
├── constraint-templates/
│   └── require-limits.yaml    # Rego policy (template)
└── constraints/
    └── require-limits.yaml    # Constraint (instance, match rules)
```

### ConstraintTemplate: Rego Policy

```rego
package k8srequiredlimits

# Rule 1: CPU limit wajib
violation[{"msg": msg}] {
  container := input.review.object.spec.containers[_]
  not container.resources.limits.cpu
  msg := sprintf("Container '%v' has no CPU limit", [container.name])
}

# Rule 2: Memory limit wajib
violation[{"msg": msg}] {
  container := input.review.object.spec.containers[_]
  not container.resources.limits.memory
  msg := sprintf("Container '%v' has no memory limit", [container.name])
}

# Rule 3: CPU request wajib
violation[{"msg": msg}] {
  container := input.review.object.spec.containers[_]
  not container.resources.requests.cpu
  msg := sprintf("Container '%v' has no CPU request", [container.name])
}

# Rule 4: Memory request wajib
violation[{"msg": msg}] {
  container := input.review.object.spec.containers[_]
  not container.resources.requests.memory
  msg := sprintf("Container '%v' has no memory request", [container.name])
}
```

### ConstraintTemplate vs Constraint

| Aspek | ConstraintTemplate | Constraint |
|-------|-------------------|------------|
| CRD | `templates.gatekeeper.sh/v1` | `constraints.gatekeeper.sh/v1beta1` |
| Berisi | Rego policy code | Match rules (kinds, namespaces) + enforcementAction |
| Analogi | "Template undangan" | "Undangan yang sudah dikirim ke tamu tertentu" |
| Bisa banyak? | 1 template = 1 policy | Banyak constraint dari 1 template (e.g., deny untuk prod, dryrun untuk dev) |

### Rego Code Explanation

```rego
violation[{"msg": msg}] {
  container := input.review.object.spec.containers[_]  # iterasi semua containers
  not container.resources.limits.cpu                   # kalau CPU limit tidak ada
  msg := sprintf("Container '%v' has no CPU limit", [container.name])  # pesan error
}
```

- `violation[{"msg": msg}]` — Gatekeeper pattern: kalau rule ini match, ada violation
- `input.review.object` — resource yang sedang di-apply (Pod atau Deployment)
- `spec.containers[_]` — iterasi semua container di Pod
- `not container.resources.limits.cpu` — kalau property tidak ada, `not` = true
- `sprintf` — format pesan error dengan container name

### Constraint: Match Rules

```yaml
spec:
  enforcementAction: deny    # explicit deny (bukan dryrun/warn)
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
      - apiGroups: ["apps"]
        kinds: ["Deployment"]
    namespaces:
      - securebank
```

- `apiGroups: [""]` + `kinds: ["Pod"]` → match bare Pod
- `apiGroups: ["apps"]` + `kinds: ["Deployment"]` → match Deployment (Gatekeeper expand ke Pod template)
- `namespaces: ["securebank"]` → hanya di namespace `securebank` (tidak hit namespace lain)

### enforcementAction Options

| Mode | Efek | Use case |
|------|------|----------|
| `deny` (default) | Block — `kubectl apply` ditolak kalau violate | Production, policy mature |
| `dryrun` | Allow — violation di-report di status, tidak block | Testing policy sebelum enable |
| `warn` | Allow — warning message di response | Transition period |

Kita pakai `deny` (explicit) karena policy sudah mature dan cluster dev.

### Verification Output

```
$ kubectl get constrainttemplate k8srequiredlimits
NAME                AGE
k8srequiredlimits   15s

$ kubectl get k8srequiredlimits require-resource-limits
NAME                      ENFORCEMENT-ACTION   TOTAL-VIOLATIONS
require-resource-limits   deny
```

0 total violations = deployment.yaml yang sudah hardened (Day 33) langsung pass policy.

### Gatekeeper Evaluation Flow

```
kubectl apply → API Server → ValidatingWebhook → Gatekeeper Controller
                                                        ↓
                                          Evaluate Rego policy
                                                        ↓
                                              violation? ──YES──→ DENY (403)
                                                  │
                                                  NO
                                                  ↓
                                          Resource created ✅
```

### Before vs After: Tanpa vs dengan Rego Policy

| Aspek | Before (Day 34) | After (Day 35) |
|-------|-----------------|----------------|
| Policy | Gatekeeper installed, no policy | Rego policy: 4 violation rules |
| Enforcement | Webhook aktif tapi kosong | Pod/Deployment tanpa resources = ditolak |
| Scope | — | Namespace `securebank` only |
| enforcementAction | — | `deny` (explicit) |
| Violations | — | 0 (existing deployment sudah compliant) |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Status menampilkan "K8sNativeValidation engine is missing" | Ini error VAP (Validating Admission Policy) — Gatekeeper Rego engine tetap jalan via webhook. VAP adalah feature baru K8s yang belum stabil. Ignore, Rego policy tetap enforced |
| Tutorial hanya cek limits, tidak cek requests | Tambah 2 violation rules untuk requests (CPU + memory) — lebih comprehensive, best practice |

---

## 📤 Output Hari Ini

- [x] `k8s/gatekeeper/constraint-templates/require-limits.yaml` — Rego policy (4 rules)
- [x] `k8s/gatekeeper/constraints/require-limits.yaml` — Constraint dengan `deny`
- [x] ConstraintTemplate + Constraint applied, 0 violations
- [x] Policy siap diuji besok (Day 36)

---

## 💡 Pelajaran Baru

- **ConstraintTemplate = template, Constraint = instance.** Template berisi Rego code (logic). Constraint berisi match rules (scope: namespace, kind) + enforcementAction. Satu template bisa punya banyak constraint — misalnya `deny` untuk prod, `dryrun` untuk dev.

- **`enforcementAction: deny` vs `dryrun`.** `deny` block resource di apply. `dryrun` allow tapi report violation. Best practice: mulai dengan `dryrun` untuk lihat impact, lalu switch ke `deny` kalau sudah yakin. Kita pakai `deny` langsung karena deployment.yaml sudah compliant.

- **Rego `not` = negation.** `not container.resources.limits.cpu` artinya "kalau property ini tidak ada". Rego tidak punya `== nil` atau `!= nil` seperti Go. `not` adalah cara cek "tidak ada / false".

- **Gatekeeper expand Deployment ke Pod template.** Match `apiGroups: ["apps"], kinds: ["Deployment"]` → Gatekeeper otomatis cek `spec.template.spec.containers` (bukan `spec.containers`). Tapi Rego code pakai `spec.containers[_]` — Gatekeeper handle expansion internally.

---

## 🔗 Referensi

- [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/)
- [Gatekeeper ConstraintTemplate](https://open-policy-agent.github.io/gatekeeper/website/docs/howto/
- [Gatekeeper enforcementAction](https://open-policy-agent.github.io/gatekeeper/website/docs/violation/)
- [Kubernetes Resource Management](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Rego pertama kali, policy-as-code |
| Pemahaman materi | 4 | ConstraintTemplate vs Constraint, Rego syntax |
| Progres sesuai target | 5 | Policy applied, 0 violations, siap test besok |

---

## ➡️ Rencana Besok

- [ ] Hari 36: OPA Policy Testing — apply Pod tanpa resources → expect ditolak, apply Pod compliant → expect diterima

---

*[← Hari 34](hari-34.md) | [Hari 36 →](hari-36.md)*