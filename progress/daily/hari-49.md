# Hari 49 — AI Remediation Node

**📅 Tanggal:** 2026-07-27
**⏱️ Durasi Belajar:** ~75 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Tambah endpoint `/webhook/sast-alert` untuk scanner findings (Semgrep/Checkov/Trivy)
- [x] Normalize 3 format scanner report ke unified finding schema
- [x] Integrate LLM via OpenCode Zen free model (`deepseek-v4-flash-free`)
- [x] Auto-post AI remediation ke Slack `#security-alerts`
- [x] Priority routing: CRITICAL/HIGH/MEDIUM → LLM + Slack, WARNING/NOTICE/LOW → log
- [x] Test end-to-end dengan 3 scanner findings
- [x] Update tracker & dokumentasi

---

## ✅ Yang Berhasil Dikerjakan

### 1. Unified Finding Format
Normalize 3 scanner output ke schema seragam:
```json
{
  "id": "CKV_AWS_260",
  "title": "Ensure no security groups allow ingress from 0.0.0.0:0 to port 80",
  "severity": "MEDIUM",
  "file": "terraform/main.tf",
  "description": "Check CKV_AWS_260 failed on resource aws_security_group.api",
  "references": ["https://docs.prismacloud.io/..."]
}
```

### 2. LLM Integration — OpenCode Zen (FREE)
- **Model:** `deepseek-v4-flash-free` — $0.00 cost
- **Endpoint:** `https://opencode.ai/zen/v1/chat/completions` (OpenAI-compatible)
- **Auth:** Bearer token
- **Prompt:** Bahasa Indonesia — ringkasan bahaya (2 kalimat) + contoh kode perbaikan
- **Response:** 496-636 chars per finding — actionable advice

### 3. Priority Routing (IF Node Logic)
```
Severity CRITICAL/HIGH/MEDIUM → call_llm() → post_to_slack_remediation()
Severity WARNING/NOTICE/LOW    → log only (skip LLM, save cost)
```

### 4. Test Results — 3 Scanners

| Test | Scanner | Severity | Route | LLM? | Slack? |
|------|---------|----------|-------|------|--------|
| 1 | Semgrep | WARNING | Log only | — | — |
| 2 | Checkov | MEDIUM | LLM + Slack | ✅ 496 chars | ✅ HTTP 200 |
| 3 | Trivy | CRITICAL | LLM + Slack | ✅ 636 chars | ✅ HTTP 200 |

### 5. Cloudflare User-Agent Fix
- **Masalah:** `urllib.request` default User-Agent (`Python-urllib/3.x`) diblokir Cloudflare → HTTP 403
- **Fix:** Tambah `User-Agent: SecureBank-WebhookReceiver/1.0` di request headers
- **Root cause:** Cloudflare WAF block bot-like User-Agent strings

---

## 📝 Catatan Teknis

### Architecture Flow
```
Scanner finding JSON → POST /webhook/sast-alert
  → normalize_finding(raw, scanner)
  → IF severity in (CRITICAL, HIGH, MEDIUM):
      → call_llm(finding) via OpenCode Zen API
      → post_to_slack_remediation(scanner, finding, ai_analysis)
  → ELSE: log only
```

### Scanner Format Mapping

| Scanner | Source Array | Key Fields |
|---------|-------------|-----------|
| Semgrep | `results[]` | `check_id`, `path`, `start.line`, `extra.message`, `extra.severity`, `extra.metadata.references` |
| Checkov | `results.failed_checks[]` | `check_id`, `check_name`, `file_path`, `resource`, `guideline`, `severity` |
| Trivy | `Results[].Vulnerabilities[]` | `VulnerabilityID`, `PkgName`, `Severity`, `Title`, `Description`, `References` |

### Test Commands
```bash
# Start webhook receiver
python3 webhook_receiver.py --port 5678

# Test 1: Semgrep WARNING (log only)
curl -X POST http://localhost:5678/webhook/sast-alert \
  -d '{"scanner":"semgrep","findings":[{"check_id":"go.lang...use-tls","path":"cmd/api/main.go","start":{"line":80},"extra":{"message":"HTTP without TLS","severity":"WARNING","metadata":{"references":["https://..."]}}}]}'

# Test 2: Checkov MEDIUM (LLM + Slack)
curl -X POST http://localhost:5678/webhook/sast-alert \
  -d '{"scanner":"checkov","findings":[{"check_id":"CKV_AWS_260","check_name":"Ensure no security groups allow ingress from 0.0.0.0:0 to port 80","file_path":"terraform/main.tf","resource":"aws_security_group.api","guideline":"https://...","severity":"MEDIUM"}]}'

# Test 3: Trivy CRITICAL (LLM + Slack)
curl -X POST http://localhost:5678/webhook/sast-alert \
  -d '{"scanner":"trivy","findings":[{"VulnerabilityID":"CVE-2024-1234","PkgName":"golang.org/x/crypto","Severity":"CRITICAL","Title":"Timing attack vulnerability","Description":"...","References":["https://..."]}]}'
```

---

## 📊 Perubahan File

| File | Status | Description |
|------|--------|-------------|
| `security/n8n-webhook/webhook_receiver.py` | ✅ Modified | +endpoint +normalize_finding +call_llm +post_to_slack_remediation |
| `security/n8n-webhook/.env` | ✅ Modified | +`OPENCODE_ZEN_API_KEY` (gitignored) |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi |
|----------|--------|
| Cloudflare block urllib request (HTTP 403) | Tambah `User-Agent` header di request |
| LLM response `content` empty, `reasoning_content` filled | `max_tokens` terlalu kecil (10), naikkan ke 512 |

---

## 📤 Output Hari Ini

- [x] Endpoint `/webhook/sast-alert` untuk general scanner findings
- [x] LLM integration via OpenCode Zen free model ($0.00)
- [x] 3 scanner format normalization (Semgrep, Checkov, Trivy)
- [x] AI remediation auto-posted ke Slack `#security-alerts`
- [x] End-to-end test: 3 findings, 2 routed to LLM+Slack, 1 log only

---

## 💡 Pelajaran Baru

- **OpenCode Zen free models** — 7 model gratis, `deepseek-v4-flash-free` paling bagus untuk reasoning, no privacy concerns
- **Cloudflare WAF blocks `Python-urllib` User-Agent** — selalu tambah custom User-Agent di urllib requests
- **OpenAI-compatible API** — Zen endpoint sama formatnya dengan OpenAI, tinggal ganti base URL + model name
- **Scanner format normalization** — setiap scanner punya JSON schema beda, unified format simplifies downstream processing

---

## 🔗 Referensi

- [Day 48](hari-48.md) — Falco alert routing ke Slack
- [Day 44](hari-44.md) — AI threat modeling dengan LLM
- [OpenCode Zen docs](https://opencode.ai/docs/zen/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | LLM free model pertama kali, langsung work |
| Pemahaman materi | 5 | Normalize + LLM + routing paham full |
| Progres sesuai target | 5 | 3/3 test pass, clean execution |

---

## ➡️ Rencana Besok

- [ ] **Day 50: CSPM (Prowler/ScoutSuite)** — Scan AWS sandbox vs CIS Benchmarks

---

*[← Hari 48](hari-48.md) | [Hari 50 →](hari-50.md)*
