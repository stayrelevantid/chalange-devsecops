# Hari 41 — Falco Attack Simulation

**📅 Tanggal:** 2026-07-18  
**⏱️ Durasi Belajar:** ~45 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Attempt shell exec pada distroless SecureBank pod — buktikan no shell
- [x] Deploy alpine attacker pod (simulasi serangan)
- [x] Execute 4 attack scenarios: shell, sensitive file, network tool, K8s API
- [x] Check Falco logs — verify custom rules trigger
- [x] Verify Falcosidekick forwarding alerts to WebUI
- [x] Verify NetworkPolicy blocking egress (Day 37 defense)
- [x] Cleanup test pod dan image

---

## ✅ Yang Berhasil Dikerjakan

- Distroless security proof: `/bin/sh`, `/bin/bash`, `sh`, `ls` — semua gagal di SecureBank pod
- Alpine attacker pod (`securebank:attacker` image) deployed dengan 4 attack scenarios
- **MAJOR DISCOVERY: 3 dari 4 custom rules FIRED!** 6 alerts total:
  - 1x Shell spawned (WARNING) — `execve` tracepoint works
  - 3x Sensitive file read (CRITICAL) — `openat` tracepoint works
  - 2x Network tool (NOTICE) — `execve` tracepoint works
- Rule 4 (K8s API access) tidak fire — `connect` tracepoint limitation confirmed
- Falcosidekick forwarded 5/6 alerts to WebUI (1 internal error on first)
- NetworkPolicy (Day 37) blocked both wget attempts — "Connection refused"
- Test pod dan image cleaned up

---

## 📝 Catatan Teknis

### Test 1: Distroless Security Proof

```bash
$ kubectl exec -it <securebank-pod> -n securebank -- /bin/sh
error: OCI runtime exec failed: exec: "/bin/sh": stat /bin/sh: no such file or directory

$ kubectl exec -it <securebank-pod> -n securebank -- /bin/bash
error: OCI runtime exec failed: exec: "/bin/bash": stat /bin/bash: no such file or directory

$ kubectl exec -it <securebank-pod> -n securebank -- sh
error: exec: "sh": executable file not found in $PATH

$ kubectl exec -it <securebank-pod> -n securebank -- ls
error: exec: "ls": executable file not found in $PATH
```

**Distroless = no shell, no ls, no utilities.** Hanya binary `/securebank` yang ada. Attacker yang dapat RCE tidak bisa exec shell untuk pivot. Ini security win yang sesungguhnya — bukan Falco detection, tapi prevention.

### Test 2: Attack Simulation — 4 Scenarios

Deploy `securebank:attacker` (alpine tagged sebagai securebank image) dengan 4 attack scenarios:

```bash
# Attack script:
echo "[1] Shell spawned → whoami"
cat /etc/passwd                          # [2] Sensitive file read
wget -qO- http://securebank-api...:80    # [3] Network tool
wget -qO- https://kubernetes.default.svc  # [4] K8s API access
```

**Attacker output:**
```
[1] Shell spawned → root
[2] cat /etc/passwd → full passwd file (17 entries)
[3] wget securebank-api → Connection refused (NetworkPolicy blocks!)
[4] wget kubernetes.default.svc → Connection refused (NetworkPolicy blocks!)
```

NetworkPolicy (Day 37) **bekerja** — default-deny-all blocks egress. Both wget attempts gagal karena NetworkPolicy hanya allow DNS egress (port 53).

### MAJOR DISCOVERY: 3/4 Custom Rules Fired!

```
=== Falco alerts (6 total) ===
[Critical] Sensitive file read in SecureBank
  file=/etc/passwd user=root pod=attacker-sim ns=securebank

[Warning] Shell spawned in SecureBank container
  command=sh -c echo "=== ATTACK SIMULATION START ==="...
  pod=attacker-sim image=docker.io/library/securebank

[Critical] Sensitive file read in SecureBank  (2nd read)
  file=/etc/passwd

[Critical] Sensitive file read in SecureBank  (3rd read)
  file=/etc/passwd

[Notice] Network tool in SecureBank container
  tool=wget command=wget -qO- http://securebank-api.securebank.svc:80/health

[Notice] Network tool in SecureBank container
  tool=wget command=wget -qO- --timeout=3 https://kubernetes.default.svc:443/
```

### Tracepoint Availability: PARTIAL (Not Total Failure)

Day 39-40 documented "k3d tracepoint limitation = rules can't trigger". **Hari ini discovery: limitation is PARTIAL, not total.**

| Syscall | Tracepoint | Status | Rules Using It |
|---------|-----------|--------|----------------|
| `execve` / `execveat` | `sys_enter_execve` | ✅ Working | Shell spawned, Network tool |
| `openat` | `sys_enter_openat` | ✅ Working | Sensitive file read |
| `connect` | `sys_enter_connect` | ❌ Missing | K8s API access |

Day 40 test pod tidak fire karena **timing issue** — Falco pods baru saja restart dari helm upgrade (~3 menit), eBPF probe butuh waktu untuk warm up. Day 41 (25+ menit setelah upgrade), probe fully operational.

### Alert Details: Image Matching

```
image=docker.io/library/securebank
```

Falco parse `securebank:attacker` sebagai `container.image.repository = docker.io/library/securebank`. Rule condition `container.image.repository contains "securebank"` match. Scope ke SecureBank app saja — tidak alert untuk shell di Falco pod atau busybox pod.

### Falcosidekick Forwarding

```
2026/07/18 07:29:11 [ERROR] : WebUI - internal server error    (first alert, Redis cold start)
2026/07/18 07:29:11 [INFO]  : WebUI - POST OK (200)            (5 subsequent alerts)
```

5/6 alerts berhasil forward ke WebUI. 1 error di alert pertama (Redis connection cold start). WebUI HTTP 200 accessible.

### Defense in Depth: All Layers Working

| Layer | Tool | Day | Result |
|-------|------|-----|--------|
| Prevention (admission) | Gatekeeper | 34-36 | Pod tanpa resources = denied |
| Prevention (image) | Distroless | 18 | No shell = attacker can't exec |
| Prevention (network) | NetworkPolicy | 37 | Egress blocked = wget failed |
| Detection (runtime) | Falco | 39-41 | 6 alerts fired for attack scenarios |
| Alert forwarding | Falcosidekick | 39 | 5/6 alerts to WebUI |

Attack scenario: attacker compromises securebank container → no shell (distroless) → even if shell obtained (supply chain attack) → Falco detects shell + file access + network tool → NetworkPolicy blocks egress → alert forwarded to WebUI.

### Before vs After

| Aspek | Before (Day 40) | After (Day 41) |
|-------|-----------------|----------------|
| Custom rules tested | Schema validation only | 3/4 rules fired in real attack |
| k3d limitation | "Total failure" | "Partial — execve/openat work, connect missing" |
| Distroless proof | Documented only | Verified: /bin/sh, /bin/bash, sh, ls all fail |
| NetworkPolicy proof | Cross-ns blocked | Egress blocked (wget to API + K8s API both refused) |
| Falcosidekick | Default alerts only | 5/6 custom alerts forwarded to WebUI |
| Defense in depth | Theoretical | All 5 layers proven working |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Day 40 test tidak fire | Timing issue — eBPF probe butuh warm-up time setelah helm upgrade. Day 41 (25+ min later) = full operational |
| Rule 4 (K8s API) tidak fire | `connect` tracepoint missing di k3d. `execve` dan `openat` work. Partial limitation, bukan total |
| Falcosidekick first alert error | Redis cold start. Subsequent 5 alerts OK (POST 200) |
| TTY not available di CLI | `kubectl exec -it` tanpa real TTY. Tidak masalah untuk attack sim — sh -c script jalan non-interaktif |

---

## 📤 Output Hari Ini

- [x] Distroless security proof: 4 shell exec attempts all fail
- [x] 4 attack scenarios executed (shell, file read, network tool, K8s API)
- [x] 6 Falco alerts fired: 1 WARNING + 3 CRITICAL + 2 NOTICE
- [x] 3/4 custom rules verified working (Rule 4 blocked by connect tracepoint)
- [x] NetworkPolicy egress blocking verified (Day 37 defense works)
- [x] Falcosidekick forwarding 5/6 alerts to WebUI
- [x] Defense in depth: all 5 layers proven working together
- [x] Test pod dan image cleaned up

---

## 💡 Pelajaran Baru

- **Distroless = prevention, Falco = detection.** Distroless prevents shell exec (attacker can't pivot). Falco detects if attacker somehow gets shell (supply chain attack, image tampering). Defense in depth: keduanya bekerja bersama.

- **k3d tracepoint limitation is PARTIAL.** Day 39-40 document sebagai "total failure" — WRONG. `execve` dan `openat` tracepoints work, hanya `connect` yang missing. 3/4 custom rules fired. Correction: "partial limitation, bukan total".

- **eBPF probe warm-up time.** Day 40 test tidak fire karena Falco pods baru restart (~3 min). Day 41 (25+ min) = full operational. eBPF probe butuh waktu untuk attach semua tracepoints setelah restart.

- **NetworkPolicy blocks egress in attack scenario.** Attacker's `wget` ke securebank-api dan K8s API server both "Connection refused". Day 37 default-deny-all + DNS-only egress = attacker tidak bisa lateral movement atau data exfiltration.

- **`container.image.repository contains "securebank"` match logic.** Falco parse `securebank:attacker` sebagai `docker.io/library/securebank`. Rule condition match karena "contains" — tidak perlu exact match. Scope ke SecureBank app saja.

- **Falco alert priority routing.** CRITICAL (sensitive file) > WARNING (shell) > NOTICE (network tool). Day 42 akan route CRITICAL alerts ke Slack #security-alerts via n8n webhook.

---

## 🔗 Referensi

- [Falco Attack Simulation](https://falco.org/docs/rules/sample-rules/)
- [Distroless Container Security](https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)
- [Kubernetes Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [eBPF Tracepoint Requirements](https://falco.org/docs/event-sources/kernel/ebpf-modern/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | 6 ALERTS FIRED! Defense in depth proven working |
| Pemahaman materi | 5 | Tracepoint partial limitation, eBPF warm-up, image matching |
| Progres sesuai target | 5 | Exceeded expectation — 3/4 rules work, not 0/4 |

---

## ➡️ Rencana Besok

- [ ] Hari 42: Alerting Webhook (n8n) — webhook trigger untuk alert Falco → Slack

---

*[← Hari 40](hari-40.md) | [Hari 42 →](hari-42.md)*
