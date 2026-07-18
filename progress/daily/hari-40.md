# Hari 40 — Falco Custom Rules

**📅 Tanggal:** 2026-07-18  
**⏱️ Durasi Belajar:** ~50 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Buat custom rules khusus untuk SecureBank (4 rules)
- [x] Create `falco-values.yaml` untuk reproducible Helm upgrades
- [x] Helm upgrade Falco dengan custom rules
- [x] Verify rules loaded (schema validation ok, 25 + 4 = 29 rules)
- [x] Test trigger — simulate attacker dengan shell di securebank container
- [x] Document k3d tracepoint limitation (rules valid, trigger blocked by k3d kernel)

---

## ✅ Yang Berhasil Dikerjakan

- 4 custom rules dibuat: Shell spawned, Sensitive file read, Network tool, K8s API access
- `falco-values.yaml` dibuat — reproducible Helm config (driver + falcosidekick + customRules)
- Helm upgrade revision 2 sukses, 3 Falco pods rolling update completed
- Schema validation: OK untuk semua 3 pods (`/etc/falco/rules.d/securebank-rules.yaml | schema validation: ok`)
- 29 total rules: 25 default + 4 custom
- Test pod dengan `securebank:test` image (busybox tagged) + `sh -c "cat /etc/passwd"` — Gatekeeper-compliant
- Rules tidak trigger karena k3d tracepoint limitation (no `/sys/kernel/tracing/events/syscalls/`)
- Test pod dan image cleaned up, SecureBank API pods tetap healthy

---

## 📝 Catatan Teknis

### 4 Custom Rules

**1. Shell spawned in SecureBank container (WARNING)**
```yaml
condition: >
  spawned_process and
  container and
  container.image.repository contains "securebank" and
  proc.name in (bash, sh, ash, zsh)
```
Distroless pods tidak punya shell. Kalau shell muncul = sangat suspicious (attacker dapat RCE atau image di-tamper).

**2. Sensitive file read in SecureBank (CRITICAL)**
```yaml
condition: >
  open_read and
  container and
  container.image.repository contains "securebank" and
  (fd.name startswith /etc/shadow or
   fd.name startswith /etc/passwd or
   fd.name startswith /proc/self/environ or
   fd.name startswith /root/.ssh or
   fd.name startswith /var/run/secrets)
```
Tutorial cuma cek 3 path. Aku tambah `/root/.ssh` (SSH keys) dan `/var/run/secrets` (K8s mounted secrets). `/var/run/secrets` relevan karena Day 38 set `automountServiceAccountToken: false` — kalau ada akses ke path ini, attacker mencoba baca K8s secrets.

**3. Network tool in SecureBank container (NOTICE)**
```yaml
condition: >
  spawned_process and
  container and
  container.image.repository contains "securebank" and
  proc.name in (nc, ncat, nmap, wget, curl, dig, nslookup, socat)
```
Tutorial cek 6 tools. Aku tambah `nslookup` (DNS recon) dan `socat` (reverse shell / port forwarding). Network tools di distroless container = potential lateral movement atau data exfiltration.

**4. K8s API access from SecureBank (WARNING) — BONUS RULE**
```yaml
condition: >
  evt.type=connect and
  (fd.typechar=4 or fd.typechar=6) and
  container and
  container.image.repository contains "securebank" and
  k8s_api_server
```
Rule ini tie ke Day 38 RBAC work: `automountServiceAccountToken: false` = tidak ada token di-mount. Kalau Pod contact K8s API server (`kubernetes.default.svc.cluster.local`), attacker mungkin telah menemukan credentials secara lain (misal dari ConfigMap, env var, atau `/var/run/secrets`).

### Helm Values File (Reproducible Upgrades)

```yaml
# securebank-api/security/falco-values.yaml
driver:
  kind: modern_ebpf
falcosidekick:
  enabled: true
  webui:
    enabled: true
customRules:
  securebank-rules.yaml: |-
    <inline rules content>
```

**Keuntungan values.yaml vs `--set` flags:**
- Reproducible — siapa pun bisa `helm upgrade --values falco-values.yaml` dengan hasil sama
- Version controlled — file committed ke repo, perubahan ter-track di git
- Self-documenting — semua config di satu tempat (driver, falcosidekick, rules)
- Helm upgrade tidak lupa flags — Day 39 pakai 3 `--set` flags, kalau upgrade tanpa values.yaml harus re-type semua

### Helm Upgrade

```bash
$ helm upgrade falco falcosecurity/falco --namespace falco --values security/falco-values.yaml
Release "falco" has been upgraded. Happy Helming!
REVISION: 2
```

Rolling update: 3 DaemonSet pods restart satu per satu (~90s total). Falcosidekick + Web UI pods tidak restart (config tidak berubah).

### Verify Rules Loaded

```
=== Files in /etc/falco/ ===
falco.yaml
falco_rules.yaml        # 25 default rules
rules.d/
  securebank-rules.yaml # 4 custom rules

=== Schema validation (all 3 pods) ===
falco-dqg5b:  /etc/falco/rules.d/securebank-rules.yaml | schema validation: ok
falco-hrtlv:  /etc/falco/rules.d/securebank-rules.yaml | schema validation: ok
falco-jnl7w:  /etc/falco/rules.d/securebank-rules.yaml | schema validation: ok

=== Rule count ===
Default rules: 25
Custom rules:  4
Total:         29 rules
```

Custom rules di-mount via ConfigMap ke `/etc/falco/rules.d/` directory. Falco config (`falco.yaml`) sudah set `rules_files` untuk include `rules.d` directory.

### Test Trigger — Simulated Attacker

```bash
# Tag busybox as securebank image (has shell, will match rule condition)
docker pull busybox:latest
docker tag busybox:latest securebank:test
k3d image import securebank:test -c securebank

# Create Gatekeeper-compliant test pod
# image: securebank:test (matches container.image.repository contains "securebank")
# command: sh -c "cat /etc/passwd && sleep 30"
# Triggers: Rule 1 (shell) + Rule 2 (sensitive file read /etc/passwd)
kubectl apply -f test-attacker-pod.yaml
```

Test pod Running di `k3d-securebank-server-0`. Falco pod di node yang sama: `falco-jnl7w`.

### Test Result: Rules Did NOT Trigger (k3d Limitation)

```
=== Falco alerts from server node (last 60s) ===
(no custom or test-attacker alerts found)

=== Tracepoint availability ===
docker exec k3d-securebank-server-0 ls /sys/kernel/tracing/events/syscalls/
→ No such file or directory

docker exec k3d-securebank-server-0 ls /sys/kernel/debug/tracing/events/syscalls/
→ No such file or directory
```

**Kenapa tidak trigger?** k3d nodes (Docker containers dengan linuxkit kernel) tidak expose `/sys/kernel/tracing/events/syscalls/` — tracepoint filesystem tidak mounted. Modern eBPF probe loaded successfully, tapi tidak bisa attach ke syscall tracepoints. `spawned_process` macro butuh `execve` tracepoint, `open_read` butuh `openat` tracepoint — keduanya tidak available.

**Ini k3d limitation, bukan Falco atau rules issue:**
- Rules schema validation: OK (3/3 pods)
- Rule conditions syntactically correct
- `container.image.repository contains "securebank"` match logic benar
- Production dengan real kernel: rules akan trigger
- Day 39 juga ada tracepoint warnings, tapi "Clear Log Activities" alert tetap firing dari old pods (events ter-capture sebelum tracepoint issue manifest)

### Before vs After

| Aspek | Before (Day 39) | After (Day 40) |
|-------|-----------------|----------------|
| Rules | 25 default only | 29 (25 default + 4 custom) |
| Helm config | `--set` flags (not reproducible) | `falco-values.yaml` (committed) |
| App-specific detection | ❌ Generic only | ✅ SecureBank-targeted |
| Rule scope | All containers | `container.image.repository contains "securebank"` |
| Helm revision | 1 | 2 |
| MITRE tags | Default tags | + mitre_execution, mitre_credential_access, mitre_discovery |

### Falco Priority Levels

| Priority | When | Custom Rules Using It |
|----------|------|----------------------|
| EMERGENCY | System unusable | — |
| ALERT | Immediate action needed | — |
| **CRITICAL** | Serious issue | Sensitive file read |
| ERROR | Error condition | — |
| **WARNING** | Warning condition | Shell spawned, K8s API access |
| **NOTICE** | Normal but significant | Network tool |
| INFO | Informational | — |
| DEBUG | Debug-level | — |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| busybox image tidak ada lokal | `docker pull busybox:latest` sebelum tag |
| Custom rules tidak muncul di `falco_rules.yaml` | Custom rules di-mount ke `/etc/falco/rules.d/` directory (separate file), bukan di main rules file. Cek dengan `ls /etc/falco/rules.d/` |
| Test trigger tidak fire | k3d tracepoint limitation — no `/sys/kernel/tracing/events/syscalls/`. Rules valid (schema validation ok), trigger blocked by k3d kernel. Documented sebagai k3d limitation |
| Helm upgrade restart pods | Expected behavior — DaemonSet rolling update. 3 pods restart satu per satu (~90s). Falcosidekick + UI pods tidak restart |

---

## 📤 Output Hari Ini

- [x] `security/falco-rules/securebank-rules.yaml` — 4 custom rules (committed)
- [x] `security/falco-values.yaml` — Helm values file (committed, reproducible)
- [x] Helm upgrade revision 2, 3 Falco pods with custom rules
- [x] Schema validation: OK (3/3 pods)
- [x] 29 total rules (25 default + 4 custom)
- [x] Test pod created and cleaned up
- [x] k3d tracepoint limitation documented

---

## 💡 Pelajaran Baru

- **Custom rules = app-specific detection.** Default rules generic (semua container). Custom rules scope ke app tertentu dengan `container.image.repository contains "securebank"`. Distroless pods tidak punya shell — kalau "Shell spawned" rule trigger, something is very wrong.

- **`rules.d/` directory untuk custom rules.** Falco load rules dari multiple files: `falco_rules.yaml` (default) + `rules.d/` directory (custom). Helm `customRules` ConfigMap di-mount ke `rules.d/`. Tidak perlu modify default rules file.

- **values.yaml > `--set` flags.** Reproducible, version-controlled, self-documenting. Day 39 pakai 3 `--set` flags — kalau upgrade tanpa values.yaml, harus re-type semua. Dengan values.yaml: `helm upgrade --values falco-values.yaml` = satu command.

- **Falco priority levels: EMERGENCY → DEBUG.** Pilih priority berdasarkan severity: CRITICAL untuk credential access, WARNING untuk shell/API access, NOTICE untuk network tools. Priority menentukan alert routing (Day 42 akan route CRITICAL → Slack #security-alerts).

- **k3d tracepoint limitation persists.** Day 39 lihat tracepoint warnings, Day 40 confirm root cause: `/sys/kernel/tracing/events/syscalls/` tidak ada di k3d nodes. Modern eBPF probe loaded, tapi tidak bisa attach ke syscall tracepoints. Production dengan real kernel = full coverage. Rules valid, limitation adalah environment-specific.

- **Rule ties ke RBAC work.** Rule 4 (K8s API access) relevan karena Day 38 set `automountServiceAccountToken: false`. Kalau Pod contact API server tanpa token = suspicious. Custom rules bisa tie ke security decisions dari hari-hari sebelumnya.

---

## 🔗 Referensi

- [Falco Custom Rules](https://falco.org/docs/rules/)
- [Falco Rule Syntax](https://falco.org/docs/rules/conditions/)
- [Falco Priority Levels](https://falco.org/docs/rules/outputs/#priority)
- [Falco Helm Chart Custom Rules](https://github.com/falcosecurity/charts/tree/master/charts/falco#custom-rules)
- [MITRE ATT&CK Framework](https://attack.mitre.org/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Custom rules = app-specific, tapi k3d limitation block test |
| Pemahaman materi | 5 | Rule syntax, priority levels, values.yaml, rules.d directory |
| Progres sesuai target | 5 | 4 rules loaded, schema validation ok, helm upgrade smooth |

---

## ➡️ Rencana Besok

- [ ] Hari 41: Falco Attack Simulation — `kubectl exec` → cek Falco log

---

*[← Hari 39](hari-39.md) | [Hari 41 →](hari-41.md)*
