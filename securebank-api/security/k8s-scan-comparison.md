# K8s Misconfiguration Scan Comparison — Day 32

## Scan Target
`securebank-api/k8s/deployment.yaml` — intentionally insecure (no securityContext, no resources, no probes)

## Scanner Results Summary

| Scanner | Total Checks | Passed | Failed | Unique Findings |
|---------|-------------|--------|--------|-----------------|
| Kubesec v2.14.2 | 14 advise | 0 | 14 | AppArmor, ServiceAccountName, Seccomp, automountSAT, RunAsGroup, RunAsNonRoot, RunAsUser, LimitsCPU, LimitsMemory, RequestsCPU, RequestsMemory, CapDropAny, CapDropAll, ReadOnlyRootFS |
| Checkov 3.3.2 | 89 | 69 | 20 | CKV_K8S_8-43 (liveness, readiness, resources, securityContext, cap drop, readOnlyFS, runAsNonRoot, seccomp, automountSAT, imagePullPolicy, image digest, secret as env, network policy) |
| Trivy 0.71.0 | 100 | 84 | 16 | KSV-0001-0118 (allowPrivilegeEscalation, cap drop, resources, runAsNonRoot, runAsUser, runAsGroup, readOnlyFS, seccomp, default security context) |

## Findings by Category

### SecurityContext (HIGH — must fix Day 33)
| Finding | Kubesec | Checkov | Trivy |
|---------|---------|---------|--------|
| No securityContext defined | ✅ | ✅ CKV_K8S_29/30 | ✅ KSV-0118 (HIGH) |
| allowPrivilegeEscalation not false | ❌ | ✅ CKV_K8S_20 | ✅ KSV-0001 (MEDIUM) |
| runAsNonRoot not true | ✅ | ✅ CKV_K8S_23 | ✅ KSV-0012 (MEDIUM) |
| No ReadOnlyRootFilesystem | ✅ | ✅ CKV_K8S_22 | ✅ KSV-0014 (HIGH) |
| Capabilities not dropped (ALL) | ✅ | ✅ CKV_K8S_28/37 | ✅ KSV-0003/0004/0106 (LOW) |
| No seccompProfile | ✅ | ✅ CKV_K8S_31 | ✅ KSV-0030/0104 (MEDIUM/LOW) |
| runAsUser > 10000 | ✅ | ✅ CKV_K8S_40 | ✅ KSV-0020 (LOW) |
| runAsGroup > 10000 | ✅ | ❌ | ✅ KSV-0021 (LOW) |
| AppArmor annotation | ✅ | ❌ | ❌ |

### Resources (MEDIUM — must fix Day 33)
| Finding | Kubesec | Checkov | Trivy |
|---------|---------|---------|--------|
| No CPU limits | ✅ | ✅ CKV_K8S_11 | ✅ KSV-0011 (LOW) |
| No CPU requests | ✅ | ✅ CKV_K8S_10 | ✅ KSV-0015 (LOW) |
| No Memory limits | ✅ | ✅ CKV_K8S_13 | ✅ KSV-0018 (LOW) |
| No Memory requests | ✅ | ✅ CKV_K8S_12 | ✅ KSV-0016 (LOW) |

### Probes (MEDIUM — must fix Day 33)
| Finding | Kubesec | Checkov | Trivy |
|---------|---------|---------|--------|
| No liveness probe | ❌ | ✅ CKV_K8S_8 | ❌ |
| No readiness probe | ❌ | ✅ CKV_K8S_9 | ❌ |

### Service Account (MEDIUM — fix Day 33/38)
| Finding | Kubesec | Checkov | Trivy |
|---------|---------|---------|--------|
| No serviceAccountName | ✅ | ❌ | ❌ |
| automountServiceAccountToken not false | ✅ | ✅ CKV_K8S_38 | ❌ |

### Other (LOW — fix Day 33 or accept)
| Finding | Kubesec | Checkov | Trivy |
|---------|---------|---------|--------|
| Image pull policy not Always | ❌ | ✅ CKV_K8S_15 | ❌ |
| Image should use digest | ❌ | ✅ CKV_K8S_43 | ❌ |
| Secrets as env vars (prefer files) | ❌ | ✅ CKV_K8S_35 | ❌ |
| No NetworkPolicy | ❌ | ✅ CKV2_K8S_6 | ❌ |

## Scanner POV Comparison

### Kubesec POV: Scoring System
- Approach: **point-based scoring** — each advisory has points (1-3), total score = sum of passed
- Output: JSON with `score` and `advise` array
- Strength: Simple, clear — "score of 0" = максимально insecure
- Weakness: No severity levels (HIGH/MEDIUM/LOW), no CIS benchmark mapping
- Unique: AppArmor annotation check (neither Checkov nor Trivy check this)

### Checkov POV: Policy-as-Code
- Approach: **CIS Kubernetes Benchmark** — 89 checks covering cluster + workload
- Output: CLI, JSON, SARIF — rich output with policy URLs
- Strength: Most comprehensive (89 checks), CIS mapping, SARIF for GitHub Security
- Weakness: Some checks are cluster-level (kubelet, API server) not applicable to deployment.yaml
- Unique: Liveness/readiness probe check, image digest check, NetworkPolicy check, secret-as-env-vars check

### Trivy POV: Severity-Based
- Approach: **KSV (Kubernetes Security Verification)** — 100 checks with severity
- Output: CLI, JSON — categorized by severity (HIGH/MEDIUM/LOW)
- Strength: Clear severity levels, link to AVD docs, code snippet with line numbers
- Weakness: No liveness/readiness probe check, no NetworkPolicy check
- Unique: KSV-0118 "using default security context" as HIGH — comprehensive catch-all

## Conclusion

All 3 scanners complement each other:
- **Kubesec**: Quick scoring (0 = bad), AppArmor check
- **Checkov**: Most comprehensive, CIS mapping, probes + NetworkPolicy
- **Trivy**: Severity-based, best for CI gate (HIGH/MEDIUM threshold)

For Day 33 fix: use all 3 reports to ensure 0 findings across all scanners.