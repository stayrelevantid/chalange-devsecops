# 🔍 Retrospektif Fase 3 — Kubernetes & Runtime Security

**📅 Periode:** Hari 31–45  
**📅 Tanggal Retrospektif:** 2026-07-22  
**⏱️ Total Waktu Belajar:** ~15 hari  
**👤 Ditulis oleh:** Muhammad Indragiri  

---

## 📊 Ringkasan Progres

| Hari | Topik | Status |
|------|-------|--------|
| 31 | K8s Cluster + Deploy | ✅ |
| 32 | K8s Misconfiguration Scan | ✅ |
| 33 | SecurityContext Hardening | ✅ |
| 34 | OPA Gatekeeper Setup | ✅ |
| 35 | Rego Policy Writing | ✅ |
| 36 | OPA Policy Testing | ✅ |
| 37 | Network Policies | ✅ |
| 38 | RBAC Auditing | ✅ |
| 39 | Falco Setup | ✅ |
| 40 | Falco Custom Rules | ✅ |
| 41 | Falco Attack Simulation | ✅ |
| 42 | Alerting Webhook (n8n) | ✅ |
| 43 | K8s Secret Management | ✅ |
| 44 | AI Threat Modeling K8s | ✅ |
| 45 | Dokumentasi Fase 3 | ✅ |

**Selesai: 15/15** ✅

---

## 🟢 Start — Hal yang Ingin Mulai Dilakukan

- **DefectDojo integration**: Semua scan results (Trivy, Semgrep, Checkov, Falco) masuk ke vulnerability dashboard untuk centralized tracking (Fase 4, Day 47)
- **Slack alerting**: Connect Falcosidekick webhook ke Slack #security-alerts — CRITICAL alerts route langsung ke team channel (Day 48)
- **Chaos security engineering**: Simulasi serangan ke cluster yang sudah hardened — test apakah defense in depth benar-benar hold (Day 54)
- **Vulnerability scanning pipeline**: Automated scan semua K8s manifests + container images + Terraform, schedule harian (Fase 4)
- **IRSA migration planning**: Untuk production, migrate dari AWS credentials di K8s Secret ke IRSA (EKS only)

---

## 🔴 Stop — Hal yang Perlu Dihentikan

- **Asumsi k3d limitation = total failure**: Day 39-40 bilang "rules can't trigger", Day 41 correction: PARTIAL. Stop making assumptions tanpa testing — always verify dengan real test
- **Secrets in git**: Base64 di K8s Secret committed ke git = bad practice. Stop committing secrets, gunakan ESO atau Sealed Secrets
- **n8n image pull on slow network**: 400MB image timed out 2x. Stop forcing tool yang tidak cooperate — pivot ke alternative yang achieve same objectives
- **Ignoring eBPF warm-up time**: Day 40 test tidak fire karena probe belum ready. Stop testing immediately after helm upgrade — wait ~25 min for eBPF to fully attach

---

## 🟡 Continue — Hal yang Sudah Baik & Perlu Diteruskan

- **Defense in depth pattern**: 6 layers (Gatekeeper + distroless + NetworkPolicy + RBAC + Falco + ESO). Setiap layer independent, tidak ada single point of failure. Lanjut dipakai dan di-extend di Fase 4
- **Default deny + whitelist**: NetworkPolicy pattern ini proven works (Day 37, 41). Apply ke semua namespace baru
- **Secret Reference Pattern**: Git = references, AWS = values, K8s = runtime cache. Pattern ini wajib untuk semua secrets
- **values.yaml untuk Helm**: Reproducible, version-controlled, self-documenting. Semua Helm install pakai values.yaml
- **3 scanner comparison**: Kubesec + Checkov + Trivy — each unique checks, complementary. Lanjut dipakai untuk K8s audit
- **AI-assisted security analysis**: Day 44 threat modeling dengan AI successful. Lanjut pakai AI untuk review, threat modeling, dan recommendations
- **Pragmatic pivoting**: n8n timeout → Python receiver. Tutorial adalah guide, bukan gospel. Kalau tool tidak cooperate, pivot ke alternative

---

## 🧠 Top 5 Pelajaran Terpenting dari Fase 3

1. **Defense in depth works** — Day 41 proven: distroless blocks shell, Falco detects, NetworkPolicy blocks egress, Gatekeeper prevents bad pods. 6 layers, 5/6 tested in real attack simulation. Tidak ada single point of failure.

2. **Distroless = prevention > detection** — 4 shell exec attempts all fail di distroless pod. Attacker dengan RCE tidak bisa exec shell untuk pivot. Ini security win paling konkret di Fase 3. Falco detection penting, tapi prevention lebih baik.

3. **Secret Reference Pattern + GitHub Push Protection** — Git cuma berisi referensi (nama), nilai di AWS. GitHub Push Protection actually blocked real AWS keys di Checkov report. Real security layer in action, bukan theoretical.

4. **k3d tracepoint limitation is PARTIAL** — Initial assumption "total failure" WRONG. execve + openat work, connect missing. 3/4 custom rules fired. eBPF warm-up time ~25 min. Always test, jangan assume.

5. **AI sebagai threat modeler** — glm-5.2 analyze cluster config, identify 3 attack paths dengan CVSS + MITRE ATT&CK. AI actually read YAML dan produce analysis. Tie ke ai-assistant-brainstorm.md — skill threat-modeler adalah real capability.

---

## 🏆 Pencapaian Terbaik

**Day 41: Falco Attack Simulation — 6 alerts fired, defense in depth proven**

Ini moment paling satisfying di Fase 3. Setelah 11 hari (Day 31-40) setup hardening — SecurityContext, Gatekeeper, NetworkPolicy, RBAC, Falco, custom rules — akhirnya test dengan real attack simulation. Hasilnya melampaui ekspektasi:

- Distroless blocks 4 shell exec attempts (prevention)
- 6 Falco alerts fired: 1 WARNING + 3 CRITICAL + 2 NOTICE (detection)
- NetworkPolicy blocks both wget attempts (network isolation)
- Falcosidekick forwards 5/6 alerts to WebUI (alert forwarding)
- k3d limitation corrected: PARTIAL, bukan total (3/4 rules work)

Semua 5 defense layers proven working together. Ini bukan theoretical — ini actually worked.

---

## 😓 Tantangan Terbesar

| Tantangan | Cara Mengatasi |
|-----------|----------------|
| k3d tracepoint limitation (Day 39-40) | Document sebagai "total failure", Day 41 correction: PARTIAL. execve/openat work, connect missing. eBPF warm-up ~25 min |
| n8n Docker image pull timed out (Day 42) | Pivot ke Python webhook receiver. Same objectives, no Docker pull, committed to repo. Pragmatic solution |
| GitHub Push Protection blocked push (Day 43) | Checkov report mengandung AWS keys. Remove report dari commit, gitignore. Push Protection = real security win |
| ESO CRD v1beta1 not found (Day 43) | Chart terbaru serve CRD sebagai v1 (GA). Update apiVersion, wait for CRD established |
| Falcosidekick UI pod ImagePullBackOff (Day 39) | Docker Hub network timeout. Delete pod, Kubernetes recreate → berhasil pull di retry |
| rakkess repo 404 (Day 38) | Repo deleted/archived. Rakkess sekarang jadi access-matrix krew plugin. Install via krew |

---

## 🔗 Koneksi ke Fase Berikutnya

> *Apa hasil dari Fase 3 yang akan langsung digunakan di Fase 4?*

- **Falco alerts + webhook receiver** → connect ke Slack #security-alerts di Day 48 (CRITICAL routing)
- **K8s cluster (hardened)** → menjadi target red team simulation di Day 52
- **OPA Gatekeeper policies** → menjadi bagian dari chaos security engineering di Day 54
- **Semua scan results** (Trivy, Checkov, Kubesec, Falco) → masuk ke DefectDojo di Day 47
- **AI threat model** (3 attack paths) → menjadi scenario untuk red team simulation
- **ESO + AWS Secrets Manager** → secret management foundation untuk production deployment
- **Defense in depth (6 layers)** → diuji kembali dengan chaos engineering di Fase 4

---

## 📈 Skor Diri (Jujur!)

| Aspek | Skor (1–10) | Catatan |
|-------|-------------|---------|
| Pemahaman K8s fundamental | 8 | k3d, manifests, namespace, deployment, service. Solid basics |
| Kemampuan menulis Rego policy | 7 | 4 violation rules, ConstraintTemplate. Masih basic, butuh practice lebih |
| Pemahaman Falco & runtime security | 8 | Install, custom rules, attack sim, webhook. eBPF tracepoint limitation understood |
| Kemampuan setup RBAC | 8 | Dedicated SA, Role with resourceNames, 3 audit tools. Least privilege verified |
| Pemahaman Network Policy | 9 | Default deny + whitelist, k3s flannel enforcement, kubelet bypass. Deep understanding |
| Konsistensi belajar harian | 8 | 15/15 selesai, tapi beberapa hari ada gap (Day 39→40, Day 43→44) |

---

## 📝 Catatan Bebas

Fase 3 ini paling intense di challenge ini. Banyak tools yang install (Gatekeeper, Falco, ESO, krew), banyak yang first time (Rego, eBPF, Falcosidekick, External Secrets). Tapi paling satisfying juga — terutama Day 41 attack simulation yang actually worked.

Discovery terbesar: k3d tracepoint limitation is PARTIAL, bukan total. Day 39-40 saya bilang "rules can't trigger", Day 41 correction. Lesson: **always test, jangan assume.** Kalau tidak test attack simulation, saya akan terus percaya limitation = total failure.

Discovery terbesar #2: GitHub Push Protection actually works. Ini bukan something yang saya planning untuk test — itu terjadi naturally saat commit Checkov report yang mengandung AWS keys. Real security layer in action.

Pivot terbaik: n8n Docker pull timed out → Python webhook receiver. Pragmatic DevSecOps = problem solving. Tutorial adalah guide, bukan gospel.

AI threat modeling (Day 44) = unique opportunity. Saya (glm-5.2) adalah AI yang menganalisis cluster config. Tie ke ai-assistant-brainstorm.md (Day 39) — skill threat-modeler adalah real capability, bukan hypothetical.

---

*[← Retrospektif Fase 2](fase-2-retrospektif.md) | [Retrospektif Fase 4 →](fase-4-retrospektif.md)*
