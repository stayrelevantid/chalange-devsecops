# Hari 38 — RBAC Auditing

**📅 Tanggal:** 2026-07-15  
**⏱️ Durasi Belajar:** ~55 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install krew + plugins (access-matrix, who-can)
- [x] Audit RBAC state — siapa yang bisa create pods, get secrets?
- [x] Buat dedicated ServiceAccount dengan least privilege
- [x] Buat Role + RoleBinding (get configmap only)
- [x] Update deployment.yaml — `serviceAccountName: securebank-api`
- [x] Verify pods Running dengan SA baru
- [x] Re-audit — confirm least privilege
- [x] Re-run Kubesec — ServiceAccountName advise resolved (score 11 → 12)

---

## ✅ Yang Berhasil Dikerjakan

- krew v0.5.0 installed (darwin_arm64), 2 plugins: access-matrix + who-can
- Audit before: pods pakai `default` SA (shared, no accountability), 0 Role/RoleBinding
- `k8s/rbac.yaml` dibuat: ServiceAccount + Role + RoleBinding
- deployment.yaml updated: `serviceAccountName: securebank-api`
- Rolling update sukses: 2 pods Running dengan SA `securebank-api`
- API tetap healthy (HTTP 200 via port-forward)
- `kubectl auth can-i` confirms least privilege: get securebank-config = YES, everything else = NO
- Kubesec score 11 → 12 (ServiceAccountName advise resolved)
- Checkov: 101 passed, 0 failed (resources 10 → 13)

---

## 📝 Catatan Teknis

### Tooling: krew + rakkess

Tutorial pakai `rakkess` binary. Tapi repo `rajat-saxena/rakkess` di GitHub sudah **404 (deleted/archived)**. Rakkess sekarang menjadi `access-matrix` plugin di krew. Solusi: install krew, kemudian `kubectl krew install access-matrix who-can`.

```bash
# Install krew v0.5.0
cd /tmp && curl -sL -o krew.tar.gz "https://github.com/kubernetes-sigs/krew/releases/download/v0.5.0/krew-darwin_arm64.tar.gz"
tar -xzf krew.tar.gz && ./krew-darwin_arm64 install krew
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"

# Install plugins
kubectl krew install access-matrix who-can
```

### Audit Before Fix

**Current state:**
- Pods pakai `default` SA (shared across all pods without serviceAccountName)
- `default` SA has 0 permissions (access-matrix shows all ✖)
- No Role or RoleBinding in `securebank` namespace
- `automountServiceAccountToken: false` already set (Day 33)

**who-can create pods -n securebank:**
```
No subjects found with permissions to create pods assigned through RoleBindings
CLUSTERROLEBINDING: cluster-admin, helm-traefik, system controllers...
(No securebank SA listed)
```

**who-can get secrets -n securebank:**
```
No subjects found with permissions to get secrets assigned through RoleBindings
CLUSTERROLEBINDING: cluster-admin, gatekeeper-admin, traefik...
(No securebank SA listed)
```

**Problem:** `default` SA is shared — no accountability. Kalau Pod compromise, attacker pakai `default` SA identity. Tidak bisa audit siapa yang akses apa.

### RBAC Fix: 3 Resources

**1. ServiceAccount** — dedicated SA untuk SecureBank API:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: securebank-api
  namespace: securebank
automountServiceAccountToken: false  # no API token mounted
```

**2. Role** — least privilege: `get` on one specific configmap only:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: securebank-api-role
  namespace: securebank
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get"]
    resourceNames: ["securebank-config"]
```

`resourceNames` = restrict ke configmap spesifik. SA tidak bisa `get` configmap lain, tidak bisa `list`/`create`/`update`/`delete` configmap apa pun.

**3. RoleBinding** — bind SA to Role:
```yaml
subjects:
  - kind: ServiceAccount
    name: securebank-api
    namespace: securebank
roleRef:
  kind: Role
  name: securebank-api-role
```

### Deployment Update

```yaml
# Before (Day 33):
spec:
  automountServiceAccountToken: false
  # no serviceAccountName — uses default SA

# After (Day 38):
spec:
  serviceAccountName: securebank-api  # dedicated SA
  automountServiceAccountToken: false
```

### Verify: Least Privilege Confirmed

```bash
$ kubectl auth can-i get configmaps/securebank-config -n securebank --as=system:serviceaccount:securebank:securebank-api
yes  # ✅ can get specific configmap

$ kubectl auth can-i get configmaps/other-config -n securebank --as=system:serviceaccount:securebank:securebank-api
no   # ✅ cannot get other configmaps

$ kubectl auth can-i list pods -n securebank --as=system:serviceaccount:securebank:securebank-api
no   # ✅ cannot list pods

$ kubectl auth can-i create pods -n securebank --as=system:serviceaccount:securebank:securebank-api
no   # ✅ cannot create pods

$ kubectl auth can-i get secrets -n securebank --as=system:serviceaccount:securebank:securebank-api
no   # ✅ cannot get secrets
```

Perfect least privilege: hanya `get` pada `securebank-config`, nothing else.

### Scanner Results: Before vs After

| Scanner | Before (Day 37) | After (Day 38) |
|---------|-----------------|----------------|
| Kubesec score | 11 | 12 (+1 ServiceAccountName) |
| Kubesec advise | 3 (AppArmor, SA, Seccomp) | 2 (AppArmor, Seccomp) |
| Checkov passed | 88 | 101 (+13 for SA/Role/RoleBinding) |
| Checkov failed | 0 | 0 |
| Checkov resources | 10 | 13 |
| Trivy | 0 findings | 0 findings |

### who-can vs access-matrix vs auth can-i

| Tool | What it shows | Best for |
|------|--------------|----------|
| `kubectl who-can <verb> <resource>` | Who can perform a verb on a resource type | Finding overprivileged subjects |
| `kubectl access-matrix --sa <name>` | Matrix of all resources × verbs for a SA | Overview of SA permissions |
| `kubectl auth can-i <verb> <resource> --as=<sa>` | Yes/No for specific action | Verifying least privilege |

`kubectl auth can-i` paling precise untuk verify least privilege — bisa test exact resource + resourceName.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| rakkess repo 404 (deleted/archived) | Rakkess sekarang menjadi `access-matrix` krew plugin. Install via `kubectl krew install access-matrix` |
| krew not installed | Download krew v0.5.0 dari GitHub releases (darwin_arm64), self-install |
| `who-can get configmaps` tidak list securebank-api SA | who-can checks verb on resource TYPE (all configmaps), bukan specific resourceName. Our Role pakai `resourceNames: ["securebank-config"]` — restrict ke configmap spesifik. `kubectl auth can-i` confirm permission works |
| access-matrix shows all ✖ for securebank-api SA | access-matrix shows LIST/CREATE/UPDATE/DELETE, bukan GET. Our Role hanya allow `get` (not shown in those columns). `kubectl auth can-i` confirm GET works |

---

## 📤 Output Hari Ini

- [x] krew v0.5.0 + 2 plugins (access-matrix, who-can) installed
- [x] `k8s/rbac.yaml` — ServiceAccount + Role + RoleBinding (least privilege)
- [x] `k8s/deployment.yaml` updated — `serviceAccountName: securebank-api`
- [x] Pods Running dengan dedicated SA, API healthy
- [x] Least privilege verified: get securebank-config = YES, all else = NO
- [x] Kubesec score 11 → 12 (ServiceAccountName resolved)
- [x] `security/kubesec-rbac-report.json` — post-RBAC Kubesec report
- [x] `security/checkov-k8s-rbac-report.json` — post-RBAC Checkov report (101/0)

---

## 💡 Pelajaran Baru

- **`default` SA = shared, no accountability.** Semua Pod tanpa `serviceAccountName` pakai `default` SA. Kalau Pod compromise, attacker pakai `default` identity — tidak bisa audit siapa yang akses apa. Dedicated SA = explicit, auditable.

- **Role + resourceNames = true least privilege.** Role tanpa `resourceNames` allow verb pada ALL resources of that type. Dengan `resourceNames: ["securebank-config"]`, SA hanya bisa `get` configmap itu saja — tidak bisa `get` configmap lain, tidak bisa `list`/`create`/`update`/`delete`.

- **`automountServiceAccountToken: false` + dedicated SA = defense in depth.** Token tidak di-mount (Pod tidak bisa talk to K8s API), TAPI Role tetap dibuat untuk "if needed" — kalau future feature butuh akses API, tinggal set `automountServiceAccountToken: true` dan SA sudah punya least privilege Role.

- **rakkess → access-matrix krew plugin.** rakkess repo di GitHub sudah 404. Maintainer memindahkan ke krew sebagai `access-matrix` plugin. Selalu cek status repo sebelum install — kalau 404, cari fork atau plugin replacement.

- **3 RBAC audit tools, 3 use cases.** `who-can` = cari siapa yang bisa verb X. `access-matrix` = overview SA permissions. `auth can-i` = verify exact permission (most precise).

---

## 🔗 Referensi

- [Kubernetes RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [kubectl krew plugin manager](https://krew.sigs.k8s.io/)
- [access-matrix (rakkess)](https://github.com/corneliusweig/rakkess)
- [kubectl who-can](https://github.com/aquasecurity/kubectl-who-can)
- [Least privilege pattern](https://kubernetes.io/docs/concepts/security/rbac-good-practices/#least-privilege)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | rakkess 404 discovery, least privilege verified |
| Pemahaman materi | 5 | SA/Role/RoleBinding, resourceNames, 3 audit tools |
| Progres sesuai target | 5 | Kubesec 12, Checkov 101/0, pods Running |

---

## ➡️ Rencana Besok

- [ ] Hari 39: Falco Setup — Helm install Falco, verifikasi logs berjalan

---

*[← Hari 37](hari-37.md) | [Hari 39 →](hari-39.md)*
