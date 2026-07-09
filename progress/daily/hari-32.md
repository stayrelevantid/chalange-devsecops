# Hari 32 — K8s Misconfiguration Scan

**📅 Tanggal:** 2026-07-09  
**⏱️ Durasi Belajar:** ~1 jam  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install Kubesec
- [x] Scan deployment.yaml dengan 3 scanner (Kubesec, Checkov, Trivy)
- [x] Temukan 5+ misconfiguration
- [x] Simpan reports JSON
- [x] Bandingkan POV 3 scanner

---

## ✅ Yang Berhasil Dikerjakan

- Install Kubesec v2.14.2 dari GitHub releases (darwin_arm64)
- Scan dengan 3 scanner — semua menemukan misconfiguration pada deployment.yaml yang intentionally insecure
- Simpan 3 report JSON ke `security/` + comparison doc `k8s-scan-comparison.md`
- Total findings: Kubesec 14 advise, Checkov 20 failed, Trivy 16 failed (3 HIGH, 3 MEDIUM, 10 LOW)

---

## 📝 Catatan Teknis

### 3 Scanner Results

| Scanner | Checks | Passed | Failed |
|---------|--------|--------|--------|
| Kubesec v2.14.2 | 14 | 0 | 14 (score 0) |
| Checkov 3.3.2 | 89 | 69 | 20 |
| Trivy 0.71.0 | 100 | 84 | 16 (3 HIGH, 3 MEDIUM, 10 LOW) |

### Key Findings (all 3 scanners agree)

| # | Finding | Severity |
|---|---------|----------|
| 1 | No securityContext (pod + container) | HIGH |
| 2 | allowPrivilegeEscalation not false | MEDIUM |
| 3 | runAsNonRoot not set | MEDIUM |
| 4 | No readOnlyRootFilesystem | HIGH |
| 5 | Capabilities not dropped (ALL) | LOW |
| 6 | No seccompProfile | MEDIUM |
| 7 | No resource limits/requests | LOW |
| 8 | runAsUser/runAsGroup > 10000 | LOW |

### Unique Findings per Scanner

| Scanner | Unique Finding |
|---------|---------------|
| Kubesec | AppArmor annotation |
| Checkov | Liveness/readiness probe, image digest, NetworkPolicy, secret-as-env-vars |
| Trivy | KSV-0118 "default security context" as HIGH catch-all |

### Scanner POV

- **Kubesec**: Point-based scoring (score 0 = worst). Simple, fast, AppArmor check
- **Checkov**: CIS Benchmark mapping, 89 checks, SARIF output. Most comprehensive
- **Trivy**: Severity-based (HIGH/MEDIUM/LOW), AVD docs links. Best for CI gate

### Report Files
- `security/kubesec-report.json`
- `security/checkov-k8s-report.json`
- `security/trivy-k8s-report.json`
- `security/k8s-scan-comparison.md`

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Kubesec tidak ada di Homebrew (brew search empty) | Download binary langsung dari GitHub releases (`kubesec_darwin_arm64.tar.gz`) |
| `wget` tidak tersedia di macOS | Pakai `curl -sL` sebagai alternatif |
| Checkov `-o` flag expects format, not file path | Pakai `--output-file-path` untuk direktori output |
| Checkov saves as `results_json.json` not specified filename | Rename after run |

---

## 📤 Output Hari Ini

- [x] Kubesec v2.14.2 installed
- [x] 3 JSON reports di `security/`
- [x] Comparison doc `k8s-scan-comparison.md`
- [x] 20+ misconfiguration teridentifikasi (Day 33 target: 0)

---

## 💡 Pelajaran Baru

- **3 scanner = 3 POV yang saling melengkapi.** Kubesec punya AppArmor, Checkov punya probes + NetworkPolicy, Trivy punya severity gate. Pakai semua untuk coverage maksimal.

- **Kubesec scoring system unik.** Score 0 artinya semua check gagal — tidak ada satupun security best practice yang diterapkan. Simple tapi efektif untuk komunikasi ke stakeholder.

- **Trivy KSV-0118 adalah catch-all finding.** "Using default security context" turun ke HIGH — ini menangkap semua yang tidak explicitly set securityContext, bukan per-property. Checkov lebih granular (CKV_K8S_29 pod level, CKV_K8S_30 container level).

- **Checkov mengecek probe dan NetworkPolicy.** Yang lain tidak. Liveness/readiness probe dianggap security karena Pod yang tidak sehat bisa mengakibatkan traffic routing issue. NetworkPolicy karena Pod tanpa policy terbuka ke semua traffic.

---

## 🔗 Referensi

- [Kubesec](https://github.com/controlplaneio/kubesec)
- [Checkov K8s policies](https://docs.prismacloud.io/en/enterprise-edition/policy-reference/kubernetes-policies)
- [Trivy KSV checks](https://avd.aquasec.com/misconfig/ksv/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | 3 scanner comparison = belajar 3 POV |
| Pemahaman materi | 5 | K8s SecurityContext, CIS benchmark, KSV checks |
| Progres sesuai target | 5 | 20+ findings, semua report tersimpan |

---

## ➡️ Rencana Besok

- [ ] Hari 33: SecurityContext Hardening — fix semua findings, re-scan, target 0 failed

---

*[← Hari 31](hari-31.md) | [Hari 33 →](hari-33.md)*