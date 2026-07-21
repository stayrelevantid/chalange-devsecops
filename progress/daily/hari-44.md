# Hari 44 — AI Threat Modeling K8s

**📅 Tanggal:** 2026-07-21  
**⏱️ Durasi Belajar:** ~60 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Kumpulkan data cluster (netpol, rbac, pods, deploy-svc)
- [x] AI analisis 3 attack paths dengan CVSS + MITRE ATT&CK
- [x] Dokumentasikan mitigasi existing (Day 31-43) + rekomendasi tambahan
- [x] Bandingkan threat model Fase 1 (app-level) vs Fase 3 (cluster-level)
- [x] Simpan threat analysis di `security/threat-model/k8s-threat-analysis.md`

---

## ✅ Yang Berhasil Dikerjakan

- Cluster data collected: netpol.yaml (2.5K), rbac.yaml (20K), pods.yaml (160K), deploy-svc.yaml (68K)
- AI analysis completed oleh glm-5.2 (opencode-go) — 3 attack paths identified
- Threat analysis document created: `k8s-threat-analysis.md` (~450 lines)
- 3 attack paths dengan CVSS scoring + MITRE ATT&CK mapping:
  1. **Container Escape via Privileged Pod** — CVSS 9.8 (Critical), MITRE T1611
  2. **Lateral Movement via NetworkPolicy Bypass** — CVSS 7.5 (High), MITRE T1021.004
  3. **Secret Exfiltration via ESO Compromise** — CVSS 8.6 (High), MITRE T1552.001
- Defense in depth validated: setiap attack path punya multiple mitigations dari Day 31-43
- 9 recommendations prioritized (3 high, 3 medium, 3 low priority)
- Fase 1 vs Fase 3 threat model comparison documented

---

## 📝 Catatan Teknis

### Cluster Data Collection

```bash
kubectl get networkpolicy -A -o yaml > security/threat-model/netpol.yaml
kubectl get role,rolebinding -A -o yaml > security/threat-model/rbac.yaml
kubectl get pods -A -o yaml > security/threat-model/pods.yaml
kubectl get deploy,svc -A -o yaml > security/threat-model/deploy-svc.yaml
```

**Files created:**
- `netpol.yaml` (2.5K) — 3 NetworkPolicies di namespace `securebank`
- `rbac.yaml` (20K) — Roles/RoleBindings untuk ESO, Falco, Gatekeeper, securebank-api
- `pods.yaml` (160K) — 15+ pods across 5 namespaces
- `deploy-svc.yaml` (68K) — Deployments + Services di semua namespaces

### AI Analysis: 3 Attack Paths

#### Attack Path 1: Container Escape via Privileged Pod

**MITRE ATT&CK:** T1611 — Escape to Host  
**CVSS v3.1:** 9.8 (Critical)  
**Vector:** `CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H`

**Exploit scenario:**
1. Attacker create privileged pod dengan `hostPath: /` mount
2. Pod dapat akses ke host filesystem (`/etc/shadow`, `/root/.ssh/`)
3. Attacker exec: `kubectl exec -it attacker-pod -- chroot /host bash`
4. Attacker sekarang di host node, bukan di container
5. Lateral movement ke node lain atau akses cluster-wide resources

**Mitigasi existing:**
- ✅ Gatekeeper (Day 34-36) — Block pods tanpa resource limits
- ✅ SecurityContext (Day 33) — `allowPrivilegeEscalation: false`, `runAsNonRoot: true`
- ✅ NetworkPolicy (Day 37) — Limit lateral movement kalau attacker berhasil escape
- ✅ Falco (Day 39-41) — Detect "release_agent File Container Escapes" + "Linux Kernel Module Injection"

**Recommendations:**
- Pod Security Standards (PSS) — Enable `restricted` profile di namespace level
- Gatekeeper policy tambahan — Block `privileged: true` containers
- K8s audit logs — Detect `pods/create` dengan `privileged: true`

#### Attack Path 2: Lateral Movement via NetworkPolicy Bypass

**MITRE ATT&CK:** T1021.004 — Remote Services: SSH  
**CVSS v3.1:** 7.5 (High)  
**Vector:** `CVSS:3.1/AV:N/AC:H/PR:L/UI:N/S:U/C:H/I:H/A:H`

**Exploit scenario:**
1. Attacker compromise `securebank-api` pod (RCE via vulnerable dependency)
2. Attacker coba reach `external-secrets` namespace (steal AWS credentials)
3. NetworkPolicy `default-deny-all` block semua egress kecuali DNS (port 53)
4. Attacker coba DNS tunneling (encode data di DNS queries)
5. Atau attacker coba exploit DNS service (CoreDNS vulnerability)

**Mitigasi existing:**
- ✅ NetworkPolicy default-deny (Day 37) — Block semua egress kecuali DNS
- ⚠️ NetworkPolicy allow-dns-egress (Day 37) — DNS allowed, tapi DNS tunneling detectable
- ✅ Distroless image (Day 18) — Tidak ada shell, attacker tidak bisa exec `curl`, `nslookup`
- ✅ Falco network tool detection (Day 40) — Detect `curl`/`wget`/`dig` execution
- ✅ RBAC least privilege (Day 38) — ServiceAccount hanya bisa `get configmap`

**Recommendations:**
- DNS logging — Enable CoreDNS query logging untuk detect DNS tunneling
- Egress to specific DNS server — Restrict DNS egress ke CoreDNS IP saja (10.43.0.10/32)
- Falco DNS tunneling detection — Custom rule untuk detect high-frequency DNS queries

#### Attack Path 3: Secret Exfiltration via ESO Compromise

**MITRE ATT&CK:** T1552.001 — Unsecured Credentials: Credentials In Files  
**CVSS v3.1:** 8.6 (High)  
**Vector:** `CVSS:3.1/AV:N/AC:L/PR:L/UI:N/S:C/C:H/I:N/A:N`

**Exploit scenario:**
1. Attacker compromise ESO pod (RCE via vulnerable dependency atau supply chain attack)
2. Attacker extract AWS credentials dari pod environment (stored di K8s Secret `aws-credentials`)
3. Attacker gunakan AWS credentials untuk akses AWS Secrets Manager langsung (bypass K8s RBAC)
4. Attacker steal semua secrets di AWS Secrets Manager (tidak hanya `securebank/jwt-secret`)
5. Attacker bisa rotate secrets di AWS, cause denial of service, atau exfiltrate data

**Mitigasi existing:**
- ✅ ESO di namespace terpisah (Day 43) — ESO pods di `external-secrets`, bukan `securebank`
- ✅ NetworkPolicy default-deny (Day 37) — Block egress dari `securebank` ke `external-secrets`
- ✅ RBAC least privilege (Day 38) — `securebank-api` ServiceAccount tidak bisa `get secrets`
- ✅ Falco runtime detection (Day 39-41) — Detect shell execution di ESO pods
- ⚠️ AWS credentials di K8s Secret (Day 43) — Kalau ESO pod compromised, credentials exposed

**Recommendations:**
- **IRSA (IAM Roles for Service Accounts)** — Gunakan IRSA instead of AWS credentials di K8s Secret (EKS only)
- ESO pod SecurityContext hardening — Apply same hardening as `securebank-api`
- AWS Secrets Manager resource-based policy — Restrict access ke specific IP/VPC
- Falco custom rule untuk ESO namespace — Detect shell execution di ESO pods

### Defense in Depth Validation

Setiap attack path punya **multiple mitigations** dari hari-hari sebelumnya:

| Attack Path | Mitigations |
|-------------|-------------|
| **Container Escape** | Gatekeeper (admission) + SecurityContext (runtime) + NetworkPolicy (network) + Falco (detection) |
| **Lateral Movement** | NetworkPolicy (network) + Distroless (runtime) + Falco (detection) + RBAC (access) |
| **Secret Exfiltration** | ESO namespace isolation (network) + RBAC (access) + NetworkPolicy (network) + Falco (detection) |

**Conclusion:** Defense in depth bekerja — tidak ada single point of failure. Setiap layer mitigasi mengurangi risk, dan multiple layers membuat attack sangat sulit.

### Threat Model Comparison: Fase 1 vs Fase 3

| Aspek | Fase 1 (App-Level) | Fase 3 (Cluster-Level) |
|-------|-------------------|------------------------|
| **Scope** | SecureBank API (single app) | K8s cluster (multiple namespaces, pods, services) |
| **Attack Surface** | HTTP endpoints, JWT, in-memory data | Pods, containers, network, RBAC, secrets |
| **Threat Actors** | External attackers, malicious users | Compromised pods, insider threats, supply chain |
| **Attack Vectors** | SQL injection, XSS, JWT bypass, IDOR | Container escape, lateral movement, RBAC abuse, secret exfiltration |
| **Mitigations** | Input validation, JWT auth, rate limiting, security headers | NetworkPolicy, RBAC, Gatekeeper, Falco, ESO, distroless |
| **Detection** | Application logs, WAF | Falco runtime detection, K8s audit logs |
| **Impact** | Data breach, account takeover | Cluster compromise, host escape, secret exfiltration |

### Recommendations Summary

#### High Priority (Immediate)

1. **Pod Security Standards (PSS)** — Enable `restricted` profile di semua namespaces (Attack Path 1)
2. **ESO pod SecurityContext hardening** — Apply same hardening as `securebank-api` (Attack Path 3)
3. **Falco custom rule untuk ESO namespace** — Detect shell execution di ESO pods (Attack Path 3)

#### Medium Priority (Next Sprint)

4. **DNS logging** — Enable CoreDNS query logging untuk detect DNS tunneling (Attack Path 2)
5. **Egress to specific DNS server** — Restrict DNS egress ke CoreDNS IP saja (Attack Path 2)
6. **AWS Secrets Manager resource-based policy** — Restrict access ke specific VPC/IP (Attack Path 3)

#### Low Priority (Future)

7. **IRSA migration** — Migrate dari AWS credentials di K8s Secret ke IRSA (Attack Path 3, EKS only)
8. **K8s audit logs** — Enable audit logs untuk detect `pods/create` dengan `privileged: true` (Attack Path 1)
9. **Falco DNS tunneling detection** — Custom rule untuk detect high-frequency DNS queries (Attack Path 2)

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Cluster data files besar (pods.yaml 160K) | Committed ke git sebagai evidence, tapi tidak di-read full. AI analysis fokus ke key parts (netpol, rbac, pod spec) |
| IRSA tidak available di k3d | Documented sebagai EKS-only feature. Untuk production di EKS, migrate ke IRSA |
| ESO pods bukan distroless | Documented sebagai risk. Recommendation: apply SecurityContext hardening (Attack Path 3) |
| AWS credentials di K8s Secret | Documented sebagai critical gap. Recommendation: IRSA migration untuk production |

---

## 📤 Output Hari Ini

- [x] `security/threat-model/netpol.yaml` — NetworkPolicy cluster data (2.5K)
- [x] `security/threat-model/rbac.yaml` — RBAC cluster data (20K)
- [x] `security/threat-model/pods.yaml` — Pods cluster data (160K)
- [x] `security/threat-model/deploy-svc.yaml` — Deployments + Services data (68K)
- [x] `security/threat-model/k8s-threat-analysis.md` — AI threat analysis (450 lines)
- [x] 3 attack paths dengan CVSS scoring + MITRE ATT&CK mapping
- [x] 9 recommendations prioritized (3 high, 3 medium, 3 low)
- [x] Defense in depth validation
- [x] Fase 1 vs Fase 3 threat model comparison

---

## 💡 Pelajaran Baru

- **AI sebagai threat modeler** — glm-5.2 (model yang powered session ini) menganalisis cluster config dan identify attack paths. Ini bukan theory — AI actually read the YAML dan produce analysis. Tie ke `docs/ai-assistant-brainstorm.md` (Day 39) di mana skill `threat-modeler` adalah salah satu 7 skills yang dirancang.

- **Cluster-level threats vs app-level threats** — Fase 1 fokus ke app-level (SQL injection, XSS, JWT bypass). Fase 3 fokus ke cluster-level (container escape, lateral movement, RBAC abuse, secret exfiltration). Different attack surface, different mitigations.

- **Defense in depth validation** — Setiap attack path punya multiple mitigations dari hari-hari sebelumnya. Tidak ada single point of failure. Setiap layer mitigasi mengurangi risk, dan multiple layers membuat attack sangat sulit.

- **CVSS scoring untuk prioritas remediation** — CVSS v3.1 memberikan standardized scoring untuk compare risks. Attack Path 1 (9.8 Critical) > Attack Path 3 (8.6 High) > Attack Path 2 (7.5 High). Prioritize high-priority recommendations first.

- **MITRE ATT&CK mapping untuk threat intelligence** — Setiap attack path mapped ke MITRE ATT&CK technique (T1611, T1021.004, T1552.001). Ini standard framework untuk threat intelligence sharing dan detection rule creation.

- **Recommendations prioritization** — 9 recommendations dibagi ke 3 priority levels (high, medium, low). High priority = immediate action (PSS, ESO hardening, Falco rules). Medium priority = next sprint (DNS logging, egress restriction, AWS policy). Low priority = future (IRSA, audit logs, DNS tunneling detection).

---

## 🔗 Referensi

- [MITRE ATT&CK Framework](https://attack.mitre.org/)
- [CVSS v3.1 Specification](https://www.first.org/cvss/v3.1/specification-document)
- [Kubernetes Security Best Practices](https://kubernetes.io/docs/concepts/security/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [External Secrets Operator](https://external-secrets.io/)
- [Falco Runtime Security](https://falco.org/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | AI threat modeling = unique opportunity, defense in depth validated |
| Pemahaman materi | 5 | 3 attack paths, CVSS scoring, MITRE ATT&CK, recommendations |
| Progres sesuai target | 5 | Threat analysis complete, 9 recommendations, Fase 1 vs 3 comparison |

---

## ➡️ Rencana Besok

- [ ] Hari 45: Dokumentasi Fase 3 — Rangkum semua K8s security yang diterapkan

---

*[← Hari 43](hari-43.md) | [Hari 45 →](hari-45.md)*
