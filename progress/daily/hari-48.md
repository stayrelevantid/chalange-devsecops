# Hari 48 — Intelligent Alert Routing (Slack)

**📅 Tanggal:** 2026-07-26
**⏱️ Durasi Belajar:** ~90 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Simpan `SLACK_WEBHOOK_URL` di `.env` — gitignored
- [x] Update webhook receiver: parse `.env`, POST ke Slack via `urllib`
- [x] Perbaiki Falco values: `host.k3d.internal` → `host.docker.internal` (DNS fix)
- [x] Test end-to-end: attacker pod → Falco → Falcosidekick → webhook → Slack
- [x] Verifikasi IF node logic: CRITICAL → Slack, non-CRITICAL → log
- [x] Cleanup test pod
- [x] Update tracker & dokumentasi

---

## ✅ Yang Berhasil Dikerjakan

### 1. `.env` + Gitignore
- File `.env` dibuat di `security/n8n-webhook/.env`
- Berisi `SLACK_WEBHOOK_URL` dengan webhook URL Slack
- Ditambahkan ke `.gitignore` — aman dari commit tidak sengaja

### 2. Webhook Receiver Update
- Module-level load `.env` file (tanpa `python-dotenv`, parse manual)
- Fungsi `post_to_slack(alert)` — pakai `urllib.request` (no pip install needed)
- Payload format: `🚨 *Falco CRITICAL Alert*\n*Rule:* ...\n*Pod:* ...`
- `post_to_slack()` dipanggil di IF node logic saat `priority == "critical"`
- Sebelumnya: hanya log "Slack payload prepared"
- Sekarang: beneran POST ke Slack webhook URL

### 3. Falco DNS Fix
- **Masalah:** Falcosidekick gagal resolve `host.k3d.internal` (NXDOMAIN)
- **Penyebab:** k3d di Docker Desktop pakai `host.docker.internal`, bukan `host.k3d.internal`
- **Fix:** `falco-values.yaml` → ubah address ke `http://host.docker.internal:5678/webhook/falco-alert`
- **Hasil:** Helm upgrade → Falcosidekick sukses kirim alert ke webhook receiver

### 4. End-to-End Test Results

| Attack | Priority | IF Node | Webhook Terima? | Slack? |
|--------|----------|---------|-----------------|--------|
| Shell spawned in SecureBank | WARNING | Log only | ✅ | — |
| Sensitive file read in SecureBank | CRITICAL | Slack | ✅ | ✅ HTTP 200 |
| Read sensitive file untrusted | WARNING | Log only | ✅ | — |
| Network tool in SecureBank | NOTICE | Log only | ✅ | — |
| Sensitive file read (2nd) | CRITICAL | Slack | ✅ | ✅ HTTP 200 |

**Total alerts received:** 6
**CRITICAL → Slack:** 2/2 (100%)
**Non-CRITICAL → Log:** 4/4 (100%)

### 5. Discovery: `host.k3d.internal` vs `host.docker.internal`
- Di k3d versi tertentu + Docker Desktop, host gateway menggunakan `host.docker.internal`
- `host.k3d.internal` hanya tersedia di k3d dengan containerd (Linux native)
- Solusi: gunakan `host.docker.internal` untuk kompatibilitas Docker Desktop

---

## 📝 Catatan Teknis

### Webhook receiver: Slack payload structure
```python
slack_payload = {
    "text": f"🚨 *Falco CRITICAL Alert*\n*Rule:* {rule}\n*Pod:* {pod} (ns: {ns})\n*Process:* {proc}\n*Output:* `{output[:200]}`",
}
```

### Falco values fix (before vs after)
```yaml
# Before (DNS NXDOMAIN)
address: "http://host.k3d.internal:5678/webhook/falco-alert"

# After (works with Docker Desktop)
address: "http://host.docker.internal:5678/webhook/falco-alert"
```

### Test commands
```bash
# Start webhook receiver
python3 webhook_receiver.py --port 5678

# Deploy attacker pod (alpine-based, tagged as securebank:attacker)
kubectl apply -f attacker-pod.yaml

# Trigger attacks
kubectl exec attacker-pod -n securebank -- bash -c "echo pwned"
kubectl exec attacker-pod -n securebank -- cat /etc/shadow
kubectl exec attacker-pod -n securebank -- curl -s http://example.com

# Check alerts
tail -f logs/falco-alerts.log
```

---

## 📊 Perubahan File

| File | Status | Description |
|------|--------|-------------|
| `.gitignore` | ✅ Modified | Added `.env` gitignore entry |
| `security/n8n-webhook/webhook_receiver.py` | ✅ Modified | Added Slack POST logic via `urllib`, `.env` parsing |
| `security/n8n-webhook/.env` | ✅ Created | `SLACK_WEBHOOK_URL` (gitignored) |
| `security/falco-values.yaml` | ✅ Modified | `host.k3d.internal` → `host.docker.internal` |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi |
|----------|--------|
| `host.k3d.internal` tidak resolve di Docker Desktop | Ganti ke `host.docker.internal` di `falco-values.yaml` |
| Falcosidekick gagal kirim webhook | DNS fix + Helm upgrade |

---

## 📤 Output Hari Ini

- [x] Slack alert routing untuk Falco CRITICAL alerts
- [x] `.env` file dengan `SLACK_WEBHOOK_URL` (gitignored)
- [x] End-to-end test: attack → Falco → webhook → Slack

---

## 💡 Pelajaran Baru

- K3d di Docker Desktop tidak support `host.k3d.internal`, harus pakai `host.docker.internal`
- `urllib.request` bisa jadi alternatif `requests` tanpa perlu pip install
- Webhook Slack bisa diverifikasi dengan curl langsung tanpa deploy ke cluster

---

## 🔗 Referensi

- [Day 42](hari-42.md) — Webhook receiver awal
- [Day 41](hari-41.md) — Attack simulation & Falco detection
- [Day 37](hari-37.md) — NetworkPolicy (masi blok egress)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | Seru liat alert nyampe ke Slack |
| Pemahaman materi | 5 | IF node logic + webhook chain paham |
| Progres sesuai target | 5 | End-to-end berhasil sekali coba |

---

## ➡️ Rencana Besok

- [ ] **Day 49: AI Remediation Node** — Auto-ringkas SAST finding via LLM → Slack dev

---

*[← Hari 47](hari-47.md) | [Hari 49 →](hari-49.md)*
