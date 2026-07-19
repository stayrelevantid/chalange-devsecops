# Hari 42 — Automated Alerting (n8n Webhook)

**📅 Tanggal:** 2026-07-19  
**⏱️ Durasi Belajar:** ~50 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Deploy webhook receiver untuk terima alert Falco
- [x] Create webhook + IF node logic (CRITICAL → Slack, others → log)
- [x] Konfigurasi Falcosidekick webhook output
- [x] Helm upgrade Falco dengan webhook config
- [x] Verify Falcosidekick: Enabled Outputs [Webhook WebUI]
- [x] Test end-to-end: attack → Falco → Falcosidekick → webhook receiver
- [x] Cleanup test pod + image

---

## ✅ Yang Berhasil Dikerjakan

- n8n Docker image pull timed out (slow network, ~400MB) → pivot ke Python webhook receiver
- `webhook_receiver.py` dibuat — simulates n8n webhook + IF node logic
- Falcosidekick webhook output enabled via `falco-values.yaml` update
- Helm upgrade revision 3 — Falcosidekick pods restarted, DaemonSet unchanged (no eBPF warm-up needed)
- `Enabled Outputs: [Webhook WebUI]` — webhook aktif
- End-to-end test: 3 alerts received from real attack (CRITICAL + WARNING + CRITICAL)
- Falcosidekick: 3 Webhook POST OK (200) — alerts forwarded successfully
- IF logic: CRITICAL → "route to Slack" path, WARNING → "log only" path
- Test pod + image cleaned up

---

## 📝 Catatan Teknis

### Pivot: n8n → Python Webhook Receiver

n8n Docker image (`n8nio/n8n`, ~400MB) pull timed out 2x (5 min each). Slow network. Pivot ke Python webhook receiver yang achieves same learning objectives:

| n8n Feature | Python Equivalent |
|-------------|------------------|
| Webhook trigger (POST /falco-alert) | `BaseHTTPRequestHandler.do_POST()` |
| IF node (priority == Critical) | `if priority.lower() == "critical":` |
| True branch → Slack (Day 48) | `logger.info("route to Slack #security-alerts")` |
| False branch → log | `logger.info("log only")` |
| Workflow execution history | `falco-alerts.log` file |

**Keuntungan pivot:**
- No Docker image pull needed (Python3 built-in)
- Committed to repo (reproducible)
- IF logic transparent (code vs n8n GUI)
- Easy to extend Day 48 (add `requests.post(slack_url)`)

**Trade-off:**
- Tidak ada GUI workflow builder
- Tidak ada n8n ecosystem integrations
- Tapi learning objectives tercapai: webhook trigger + IF routing + end-to-end pipeline

### Webhook Receiver Design

```python
# IF node logic — CRITICAL routes to Slack (Day 48), others just log
if priority.lower() == "critical":
    logger.info("IF TRUE: CRITICAL — route to Slack #security-alerts (Day 48)")
    # Day 48 will add: requests.post(slack_webhook_url, json=slack_payload)
else:
    logger.info("IF FALSE: {priority} — log only (below CRITICAL threshold)")
```

Endpoints:
- `POST /webhook/falco-alert` — receives Falco alerts
- `GET /health` — health check
- `GET /alerts` — view last 20 alerts

### Falcosidekick Webhook Config

```yaml
# falco-values.yaml (updated)
falcosidekick:
  enabled: true
  webui:
    enabled: true
  config:
    webhook:
      address: "http://host.k3d.internal:5678/webhook/falco-alert"
      method: "POST"
```

`host.k3d.internal` = k3d DNS yang resolve ke host machine (192.168.65.254). Falcosidekick pod di cluster bisa reach webhook receiver di host melalui DNS ini.

### Helm Upgrade — Revision 3

```bash
$ helm upgrade falco falcosecurity/falco --namespace falco --values security/falco-values.yaml
REVISION: 3
```

**Yang restart:** Falcosidekick pods (config berubah — webhook added)
**Yang TIDAK restart:** Falco DaemonSet pods (config tidak berubah — driver + rules sama)

Ini penting: tidak ada eBPF probe warm-up (Day 41 lesson). Alert bisa trigger immediately.

### Verify: Webhook Output Enabled

```
$ kubectl logs -l app.kubernetes.io/name=falcosidekick -n falco --tail=10
[INFO] : Enabled Outputs: [Webhook WebUI]
[INFO] : Falcosidekick is up and listening on :2801
```

Dari `[WebUI]` (Day 39) → `[Webhook WebUI]` (Day 42). Webhook output aktif.

### End-to-End Test

```bash
# 1. Start webhook receiver
python3 webhook_receiver.py --port 5678

# 2. Deploy attacker pod (alpine tagged securebank:attacker)
kubectl apply -f test-attacker-e2e.yaml
# command: sh -c "cat /etc/passwd && sleep 15"

# 3. Check webhook receiver logs
```

**Webhook receiver logs (real attack alerts):**
```
ALERT RECEIVED | priority=Critical rule=Sensitive file read in SecureBank pod=attacker-e2e ns=securebank
  -> IF TRUE: CRITICAL alert — route to Slack #security-alerts (Day 48)
  -> Slack payload prepared: 'Sensitive file read in SecureBank' in pod 'attacker-e2e'

ALERT RECEIVED | priority=Warning rule=Shell spawned in SecureBank container pod=attacker-e2e ns=securebank
  -> IF FALSE: Warning alert — log only (below CRITICAL threshold)

ALERT RECEIVED | priority=Critical rule=Sensitive file read in SecureBank pod=attacker-e2e ns=securebank
  -> IF TRUE: CRITICAL alert — route to Slack #security-alerts (Day 48)
```

**Falcosidekick logs:**
```
[INFO] : Webhook - POST OK (200)   # 3 alerts forwarded to webhook
[INFO] : WebUI - POST OK (200)     # alerts also to WebUI dashboard
```

### Alert Pipeline: Full Chain

```
Attacker (cat /etc/passwd)
  → Falco eBPF probe (syscall capture: openat /etc/passwd)
  → Falco engine (rule match: Sensitive file read, CRITICAL)
  → Falcosidekick (receive alert, forward to outputs)
  → Webhook (POST http://host.k3d.internal:5678/webhook/falco-alert)
  → Python webhook receiver (IF logic: CRITICAL → Slack path)
  → [Day 48: Slack #security-alerts notification]
```

### Before vs After

| Aspek | Before (Day 41) | After (Day 42) |
|-------|-----------------|----------------|
| Falcosidekick outputs | WebUI only | WebUI + Webhook |
| Alert routing | Dashboard only | IF logic (CRITICAL → Slack) |
| Automated response | ❌ None | ✅ Webhook receiver with routing |
| Helm revision | 2 | 3 |
| Pipeline | Falco → WebUI | Falco → Falcosidekick → Webhook → IF routing |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| n8n Docker image pull timed out (2x, 5 min each) | Pivot ke Python webhook receiver — same learning objectives, no Docker pull needed, committed to repo |
| n8n Helm chart repo 404 (8k8s.github.io) | Tidak dipakai — Python receiver tidak butuh Helm chart |
| Falcosidekick config webhook key | `falcosidekick.config.webhook.address` + `method: POST` — verified via `helm show values` |
| `host.k3d.internal` connectivity | Verified: resolves to 192.168.65.254, Falcosidekick pod can reach port 5678 |
| WebUI rate limit error | Falcosidekick WebUI has "exceeding post rate limit (500)" — 3 alerts in rapid succession. Webhook tidak ada rate limit issue |

---

## 📤 Output Hari Ini

- [x] `security/n8n-webhook/webhook_receiver.py` — Python webhook receiver (committed)
- [x] `security/n8n-webhook/logs/falco-alerts.log` — alert audit log
- [x] `security/falco-values.yaml` updated — webhook config added
- [x] Helm revision 3 — Falcosidekick webhook output enabled
- [x] End-to-end pipeline tested: 3 alerts received, IF routing works
- [x] Falcosidekick: 3 Webhook POST OK (200)

---

## 💡 Pelajaran Baru

- **DevSecOps = problem solving, bukan dogmatic tool following.** n8n image pull timed out → pivot ke Python webhook receiver yang achieves same objectives. Tutorial adalah guide, bukan gospel. Pragmatic solution > stuck on tooling.

- **Falcosidekick multi-output.** WebUI (dashboard, human review) + Webhook (automation, programmatic routing). Multiple outputs = defense in depth untuk alerting. WebUI untuk visual inspection, Webhook untuk automated response.

- **`host.k3d.internal` = k3d DNS untuk host access.** Pods di k3d cluster bisa reach host services melalui `host.k3d.internal` (resolves to 192.168.65.254). Berguna untuk local dev — webhook receiver di host, Falcosidekick di cluster.

- **Helm upgrade selective restart.** Falcosidekick pods restart (config berubah), Falco DaemonSet pods tidak restart (config sama). Artinya: tidak ada eBPF probe warm-up (Day 41 lesson). Alert bisa trigger immediately after upgrade.

- **IF node = conditional alert routing.** CRITICAL → Slack (immediate notification), WARNING/NOTICE → log (periodic review). Priority-based routing = alert fatigue prevention. Day 48 akan implement Slack webhook untuk CRITICAL path.

---

## 🔗 Referensi

- [Falcosidekick Webhook Output](https://github.com/falcosecurity/falcosidekick#webhook)
- [Falcosidekick Helm Chart Config](https://github.com/falcosecurity/charts/blob/master/charts/falcosidekick/values.yaml)
- [n8n Workflow Automation](https://n8n.io/)
- [k3d host.k3d.internal DNS](https://k3d.io/v5.6.0/usage/kubeconfig/#using-k3d-with-external-datastores)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | End-to-end pipeline works! Pivot dari n8n ke Python = pragmatic |
| Pemahaman materi | 5 | Webhook config, IF routing, host.k3d.internal, selective restart |
| Progres sesuai target | 5 | 3 alerts received via webhook, IF logic verified, pipeline complete |

---

## ➡️ Rencana Besok

- [ ] Hari 43: K8s Secret Management — External Secrets Operator → AWS Secrets Manager

---

*[← Hari 41](hari-41.md) | [Hari 43 →](hari-43.md)*
