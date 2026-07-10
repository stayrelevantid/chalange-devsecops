# Hari 33 — SecurityContext Hardening

**📅 Tanggal:** 2026-07-10  
**⏱️ Durasi Belajar:** ~1 jam  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Fix semua 20+ findings dari Day 32
- [x] Update deployment.yaml dengan SecurityContext hardening
- [x] Re-deploy ke k3d cluster
- [x] Re-scan dengan 3 scanner — target 0 failed
- [x] Pod tetap running dengan hardened config

---

## ✅ Yang Berhasil Dikerjakan

- Deployment.yaml hardened: 26 lines → 65 lines (39 lines security added)
- Pod-level: automountServiceAccountToken, runAsNonRoot, runAsUser 65532, runAsGroup, fsGroup, seccompProfile RuntimeDefault
- Container-level: allowPrivilegeEscalation false, readOnlyRootFilesystem true, capabilities drop ALL
- Resources: requests (cpu 50m, mem 64Mi) + limits (cpu 200m, mem 128Mi)
- Probes: liveness + readiness (httpGet /health:8080)
- emptyDir volume for /tmp (preventive, best practice with readOnlyRootFilesystem)
- imagePullPolicy: IfNotPresent (local k3d, no registry)
- 2 pods Running, /health endpoint tested OK

### Scanner Results: Before vs After

| Scanner | Before (Day 32) | After (Day 33) | Status |
|---------|-----------------|----------------|--------|
| Kubesec | Score 0 (14 advise) | **Score 11** (3 advise remain) | ✅ 11/14 passed |
| Checkov | 69 passed, **20 failed** | **85 passed, 0 failed** (4 skipped) | ✅ 0 failed |
| Trivy | **16 findings** (3H/3M/10L) | **0 findings** | ✅ Clean |

---

## 📝 Catatan Teknis

### BEFORE: deployment.yaml (Day 31 — Insecure)

```yaml
# 26 lines — no security whatsoever
spec:
  template:
    spec:
      containers:
        - name: api
          image: securebank:v1
          ports:
            - containerPort: 8080
              protocol: TCP
          envFrom:
            - secretRef:
                name: securebank-secrets
```

**Masalah:** No securityContext, no resources, no probes, no automountServiceAccountToken, no imagePullPolicy, no volumes.

### AFTER: deployment.yaml (Day 33 — Hardened)

```yaml
# 65 lines — 39 lines of security added
spec:
  template:
    spec:
      automountServiceAccountToken: false           # NEW
      securityContext:                               # NEW — Pod level
        runAsNonRoot: true
        runAsUser: 65532                             # distroless nonroot UID
        runAsGroup: 65532
        fsGroup: 65532
        seccompProfile:
          type: RuntimeDefault
      volumes:                                       # NEW — emptyDir for /tmp
        - name: tmp
          emptyDir: {}
      containers:
        - name: api
          image: securebank:v1
          imagePullPolicy: IfNotPresent              # NEW
          ports:
            - containerPort: 8080
              protocol: TCP
          envFrom:
            - secretRef:
                name: securebank-secrets
          volumeMounts:                              # NEW — mount tmp
            - name: tmp
              mountPath: /tmp
          securityContext:                           # NEW — Container level
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          resources:                                 # NEW — Resource limits
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 200m
              memory: 128Mi
          livenessProbe:                             # NEW — Health check
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:                            # NEW — Traffic readiness
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 3
            periodSeconds: 5
```

### Line-by-Line Changes Summary

| # | Property Added | Level | Scanner Finding Fixed | Day-32 Severity |
|---|---------------|-------|----------------------|-----------------|
| 1 | `automountServiceAccountToken: false` | Pod | Kubesec, CKV_K8S_38 | MEDIUM |
| 2 | `runAsNonRoot: true` | Pod | Kubesec, CKV_K8S_23, KSV-0012 | MEDIUM |
| 3 | `runAsUser: 65532` | Pod | Kubesec, CKV_K8S_40, KSV-0020 | LOW |
| 4 | `runAsGroup: 65532` | Pod | Kubesec, KSV-0021 | LOW |
| 5 | `fsGroup: 65532` | Pod | best practice | — |
| 6 | `seccompProfile: RuntimeDefault` | Pod | Kubesec, CKV_K8S_31, KSV-0030/0104 | MEDIUM |
| 7 | `allowPrivilegeEscalation: false` | Container | CKV_K8S_20, KSV-0001 | MEDIUM |
| 8 | `readOnlyRootFilesystem: true` | Container | Kubesec, CKV_K8S_22, KSV-0014 | HIGH |
| 9 | `capabilities.drop: ["ALL"]` | Container | Kubesec, CKV_K8S_28/37, KSV-0003/0004/0106 | LOW |
| 10 | `resources.requests` (cpu/mem) | Container | Kubesec, CKV_K8S_10/12, KSV-0015/0016 | LOW |
| 11 | `resources.limits` (cpu/mem) | Container | Kubesec, CKV_K8S_11/13, KSV-0011/0018 | LOW |
| 12 | `livenessProbe` | Container | CKV_K8S_8 | MEDIUM |
| 13 | `readinessProbe` | Container | CKV_K8S_9 | MEDIUM |
| 14 | `imagePullPolicy: IfNotPresent` | Container | CKV_K8S_15 (skipped) | LOW |
| 15 | `emptyDir` volume for /tmp | Pod | best practice with readOnlyRootFilesystem | — |
| 16 | `volumeMounts` for /tmp | Container | supports readOnlyRootFilesystem | — |

### UID: 65532 vs 65534

Tutorial pakai `runAsUser: 65534` (nobody). Tapi distroless `gcr.io/distroless/static-debian12:nonroot` pakai UID **65532**. Kalau pakai 65534, Pod akan crash karena UID mismatch. Selalu cek image User dengan `docker inspect <image> --format '{{.Config.User}}'`.

### Findings yang Di-accept (with justification + skip-check)

| Finding | Scanner | Skip ID | Reason | Remediation Day |
|---------|---------|---------|--------|-----------------|
| Image pull policy not Always | Checkov | CKV_K8S_15 | Local k3d (no registry), IfNotPresent = correct | Saat pakai registry |
| Image should use digest | Checkov | CKV_K8S_43 | Dev environment, tag v1 fine | Production |
| Secrets as env vars | Checkov | CKV_K8S_35 | envFrom secretRef, will migrate | Day 43 (External Secrets) |
| No NetworkPolicy | Checkov | CKV2_K8S_6 | Dedicated day | Day 37 |
| AppArmor | Kubesec | — | k3d local, kernel-level | Production |
| serviceAccountName | Kubesec | — | Dedicated day | Day 38 (RBAC) |

### Report Files (Post-Fix)
- `security/kubesec-post-fix-report.json` — score 11
- `security/checkov-k8s-post-fix-report.json` — 85/0 (4 skipped)
- `security/trivy-k8s-post-fix-report.json` — 0 findings

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| Port 9080 masih dipakai dari port-forward hari lalu | `pkill -f port-forward` sebelum port-forward baru |
| Rolling update: pod lama Error, pod baru ContainerCreating | Tunggu ~15 detik, k3d rolling update selesai otomatis |
| Tutorial pakai runAsUser 65534 | Cek distroless UID via `docker inspect` → 65532, sesuaikan |
| readOnlyRootFilesystem bisa crash app | Tambah emptyDir volume untuk /tmp sebagai preventive best practice |

---

## 📤 Output Hari Ini

- [x] deployment.yaml hardened (26 → 65 lines, 16 security properties added)
- [x] 2 pods Running dengan hardened config
- [x] Kubesec: 0 → 11 score (11/14 passed)
- [x] Checkov: 20 failed → 0 failed (85 passed, 4 skipped)
- [x] Trivy: 16 findings → 0 findings (clean)
- [x] 3 post-fix report JSONs saved

---

## 💡 Pelajaran Baru

- **`readOnlyRootFilesystem: true` + emptyDir = best practice.** Immutable root filesystem mencegah attacker menulis binary/config. Tapi app mungkin butuh write ke /tmp, jadi mount emptyDir. Preventive lebih baik daripada debug crash-loop nanti.

- **UID harus match dengan image.** Tutorial pakai 65534, tapi distroless `nonroot` pakai 65532. Kalau mismatch, Pod crash dengan `exec user process caused: setgroups: operation not permitted`. Selalu cek `docker inspect <image> --format '{{.Config.User}}'`.

- **3 scanner = 3 blind spot.** Setelah hardening, Trivy langsung 0. Checkov 0 (dengan 4 skip). Kubesec masih ada 3 advise (AppArmor, ServiceAccountName, SeccompAny annotation) yang akan di-remediasi di hari berikutnya. Tanpa 3 scanner, beberapa finding mungkin terlewat.

- **`imagePullPolicy: IfNotPresent` untuk local k3d.** Kalau pakai `Always` dengan local image (tanpa registry), K8s akan cari image di registry yang tidak ada → ImagePullBackOff. `IfNotPresent` = pakai image yang sudah ada di node. Nanti saat pakai registry, ganti ke `Always`.

---

## 🔗 Referensi

- [Kubernetes SecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
- [Pod Security Standards: Restricted](https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted)
- [Checkov K8s skip checks](https://www.checkov.io/5.Policy%20Index/kubernetes.html)
- [readOnlyRootFilesystem + emptyDir best practice](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | 20+ → 0 findings, sangat satisfying |
| Pemahaman materi | 5 | SecurityContext pod vs container, UID match, emptyDir |
| Progres sesuai target | 5 | 3 scanner all green (or near-green), pods running |

---

## ➡️ Rencana Besok

- [ ] Hari 34: RBAC Setup — ServiceAccount, Role, RoleBinding untuk least privilege

---

*[← Hari 32](hari-32.md) | [Hari 34 →](hari-34.md)*