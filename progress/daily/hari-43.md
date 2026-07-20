# Hari 43 — K8s Secret Management (External Secrets Operator)

**📅 Tanggal:** 2026-07-20  
**⏱️ Durasi Belajar:** ~55 menit  
**🏷️ Fase:** Fase 3 — K8s & Runtime Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Create AWS Secrets Manager secret (securebank/jwt-secret)
- [x] Install External Secrets Operator via Helm
- [x] Create SecretStore + ExternalSecret + aws-credentials manifests
- [x] Gitignore aws-credentials.yaml (AWS keys tidak di-commit)
- [x] Apply manifests, verify ESO sync (SecretSynced)
- [x] Update deployment.yaml — envFrom securebank-jwt (ESO-managed)
- [x] Redeploy, verify API works (health 200, balance 401)
- [x] Re-run Checkov — CKV_K8S_35 status documented

---

## ✅ Yang Berhasil Dikerjakan

- AWS Secrets Manager secret `securebank/jwt-secret` created (ap-southeast-1, $0.40/month)
- ESO Helm install: 3 pods Running (external-secrets, cert-controller, webhook)
- 24 CRDs installed (SecretStore, ExternalSecret, ClusterSecretStore, dll)
- SecretStore `aws-secrets` — Valid, Ready, ReadWrite capabilities
- ExternalSecret `securebank-jwt` — SecretSynced, K8s Secret auto-created in 20s
- `aws-credentials.yaml` gitignored — AWS keys tidak masuk git history
- deployment.yaml updated: envFrom `securebank-secrets` → `securebank-jwt` (ESO-managed)
- Rolling update sukses, 2 pods Running, API healthy (HTTP 200)
- /balance returns 401 without JWT — JWT middleware active = JWT_SECRET loaded from ESO secret
- Checkov: 102 passed, 0 failed (with CKV_K8S_35 skip — trade-off documented)

---

## 📝 Catatan Teknis

### Architecture: Secret Reference Pattern

```
Git Repo (committed)              AWS (encrypted)              K8s Cluster (runtime)
┌──────────────────────┐         ┌─────────────────┐         ┌──────────────────────┐
│ secret-store.yaml    │         │ Secrets Manager │         │ aws-credentials      │
│  (ref: aws-creds)    │────────►│ securebank/     │◄────────│  (AWS keys)          │
│                      │  ESO    │ jwt-secret      │  reads  │                      │
│ jwt-secret.yaml      │  reads  │ ($0.40/month)   │         │ SecretStore          │
│  (ref: AWS secret)   │         └─────────────────┘         │  aws-secrets          │
│                      │                                     │                      │
│ deployment.yaml      │                                     │ ExternalSecret       │
│  (envFrom: jwt)      │                                     │  securebank-jwt       │
└──────────────────────┘                                     │  ↓ auto-creates      │
                                                             │ K8s Secret           │
                                                             │  securebank-jwt       │
                                                             │  ↓ envFrom           │
                                                             │ Pod (JWT_SECRET env) │
└──────────────────────┘                                     └──────────────────────┘
   Git: references only                                         Runtime: actual values
```

### AWS Secrets Manager Secret

```bash
$ aws secretsmanager create-secret \
  --name securebank/jwt-secret \
  --secret-string "$(openssl rand -hex 32)" \
  --region ap-southeast-1

ARN: arn:aws:secretsmanager:ap-southeast-1:683915449775:secret:securebank/jwt-secret-RTgjr9
```

### ESO Helm Install

```bash
$ helm install external-secrets external-secrets/external-secrets \
  --namespace external-secrets --create-namespace

3 pods Running:
- external-secrets-678fb7bf9b-8bw77 (main controller)
- external-secrets-cert-controller-6cf45d9f6-fnwlg (cert management)
- external-secrets-webhook-548487f58b-475mg (admission webhook)

24 CRDs installed (SecretStore, ExternalSecret, ClusterSecretStore, dll)
```

### 3 Manifest Files

**1. aws-credentials.yaml (GITIGNORED)**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-credentials
  namespace: securebank
type: Opaque
stringData:
  access-key-id: "AKIAZ6PE..."      # Real AWS keys
  secret-access-key: "hehAoIu6..."  # NOT committed to git
```

**2. secret-store.yaml (COMMITTED)**
```yaml
apiVersion: external-secrets.io/v1
kind: SecretStore
metadata:
  name: aws-secrets
  namespace: securebank
spec:
  provider:
    aws:
      service: SecretsManager
      region: ap-southeast-1
      auth:
        secretRef:
          accessKeyIDSecretRef:
            name: aws-credentials     # Reference only, no values
            key: access-key-id
          secretAccessKeySecretRef:
            name: aws-credentials
            key: secret-access-key
```

**3. jwt-secret.yaml (COMMITTED)**
```yaml
apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata:
  name: securebank-jwt
  namespace: securebank
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets
    kind: SecretStore
  target:
    name: securebank-jwt
    creationPolicy: Owner
  data:
    - secretKey: JWT_SECRET
      remoteRef:
        key: securebank/jwt-secret   # AWS secret name
```

### ESO Sync Verification

```bash
$ kubectl get externalsecret -n securebank
NAME             STORETYPE     STORE         REFRESH INTERVAL   STATUS         READY   LAST SYNC
securebank-jwt   SecretStore   aws-secrets   1h                 SecretSynced   True    20s

$ kubectl get secretstore -n securebank
NAME          AGE   STATUS   CAPABILITIES   READY
aws-secrets   20s   Valid    ReadWrite      True

$ kubectl get secret securebank-jwt -n securebank
NAME             TYPE     DATA   AGE
securebank-jwt   Opaque   1      20s
# JWT_SECRET: (len=88) — base64 encoded, auto-created by ESO
```

### Deployment Update

```yaml
# Before (Day 31-42):
envFrom:
  - secretRef:
      name: securebank-secrets  # Manual K8s Secret (base64 in git)

# After (Day 43):
envFrom:
  - secretRef:
      name: securebank-jwt  # ESO-managed (synced from AWS)
```

### API Verification

```bash
$ curl http://localhost:9086/health
{"status":"healthy"}  # HTTP 200 — pods running with ESO secret

$ curl http://localhost:9087/balance
missing authorization header  # HTTP 401 — JWT middleware active
```

HTTP 401 dari /balance membuktikan JWT middleware aktif = JWT_SECRET berhasil di-load dari ESO-managed secret.

### Checkov: CKV_K8S_35 Trade-off

```bash
$ checkov -d k8s/ --framework kubernetes --check CKV_K8S_35
FAILED for resource: Deployment.securebank.securebank-api
Check: CKV_K8S_35: "Prefer using secrets as files over secrets as environment variables"
```

**CKV_K8S_35 masih FAILED** — ESO fixes **source** (AWS vs git), bukan **mount method** (env vs file).

| Aspek | Before (Day 31) | After (Day 43) |
|-------|-----------------|----------------|
| Secret source | Git (base64 in secret.yaml) | AWS Secrets Manager (encrypted) |
| Git exposure | ❌ Secret in git history | ✅ No secret in git |
| Auto-rotation | ❌ Manual | ✅ refreshInterval: 1h |
| Mount method | envFrom | envFrom (unchanged) |
| CKV_K8S_35 | FAILED | FAILED (same — env mount) |

**Kenapa CKV_K8S_35 masih failed?** Checkov flags `envFrom secretRef` — prefer secrets as **files** (volumeMounts). Tapi app pakai `os.Getenv("JWT_SECRET")`, butuh env. Fix full: mount as file + modify app to read file. Beyond Day 43 scope.

**With skip: 102 passed, 0 failed.** Trade-off documented: ESO improves source security, CKV_K8S_35 skip remains for env mount.

### CRD Version: v1 (not v1beta1)

Tutorial pakai `apiVersion: external-secrets.io/v1beta1`. Tapi ESO chart terbaru (2026) serve CRD sebagai `v1` (GA). Fix: update `v1beta1` → `v1` di secret-store.yaml dan jwt-secret.yaml.

### Cost

| Item | Cost | Notes |
|------|------|-------|
| AWS Secrets Manager secret | $0.40/month | 1 secret, ap-southeast-1 |
| ESO pods (3 pods) | ~100MB RAM | k3d local, no cloud cost |
| Total | $0.40/month | Acceptable for learning |

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| CRD `v1beta1` not found | ESO chart terbaru serve CRD sebagai `v1` (GA). Update apiVersion dari `v1beta1` → `v1` |
| CRD not established immediately | `kubectl wait --for=condition=Established crd/secretstores.external-secrets.io` — wait ~30s after Helm install |
| CKV_K8S_35 still FAILED | ESO fixes source (AWS), bukan mount method (env). App butuh `os.Getenv()`. Skip documented: env mount trade-off |
| `aws-credentials.yaml` contains real AWS keys | Gitignore file — `securebank-api/k8s/external-secrets/aws-credentials.yaml` di .gitignore |
| /auth endpoint 404 | App tidak punya /auth endpoint. JWT generated externally. /balance 401 = JWT middleware works = JWT_SECRET loaded |

---

## 📤 Output Hari Ini

- [x] AWS Secrets Manager secret `securebank/jwt-secret` created
- [x] ESO Helm install (3 pods, 24 CRDs)
- [x] `k8s/external-secrets/secret-store.yaml` (committed)
- [x] `k8s/external-secrets/jwt-secret.yaml` (committed)
- [x] `k8s/external-secrets/aws-credentials.yaml` (gitignored)
- [x] `.gitignore` updated
- [x] `k8s/deployment.yaml` updated — envFrom `securebank-jwt`
- [x] ESO sync verified: SecretSynced, K8s Secret auto-created
- [x] API verified: health 200, balance 401 (JWT middleware active)
- [x] `security/checkov-k8s-eso-report.json` — Checkov report (102/0 with skip)
- [x] CKV_K8S_35 trade-off documented

---

## 💡 Pelajaran Baru

- **Secret Reference Pattern.** Git repo cuma berisi referensi (nama-nama), bukan nilai secret. SecretStore references `aws-credentials` (nama), ExternalSecret references `securebank/jwt-secret` (nama AWS). Kalau repo leak, attacker tahu nama-nama tapi tidak tahu nilainya.

- **ESO = decouple secrets from git.** Source of truth di AWS Secrets Manager (encrypted, auditable, auto-rotation). K8s Secret cuma cache yang auto-sync. `refreshInterval: 1h` = kalau AWS secret di-rotate, K8s auto-sync within 1 jam.

- **CKV_K8S_35 = mount method, not source.** ESO fixes source (AWS vs git), tapi CKV_K8S_35 flags mount method (env vs file). Full fix: mount as file + modify app to read file. Trade-off: app butuh `os.Getenv()`, env mount tetap dipakai.

- **CRD version evolution.** Tutorial pakai `v1beta1`, chart terbaru serve `v1` (GA). Selalu cek `kubectl api-resources --api-group=external-secrets.io` untuk version yang available.

- **Tutorial reality for non-EKS.** AWS keys masih di K8s Secret (`aws-credentials`). IRSA (IAM Roles for Service Accounts) = EKS-only, tidak available di k3d. Workaround: AWS keys di gitignored K8s Secret. Production di EKS: pakai IRSA, no keys in cluster.

- **Gitignore = security boundary for secrets.** `aws-credentials.yaml` di .gitignore = AWS keys tidak masuk git history. File ada di disk (untuk kubectl apply), tapi tidak di-commit. Pattern ini = "local-only secrets".

---

## 🔗 Referensi

- [External Secrets Operator](https://external-secrets.io/)
- [ESO AWS Secrets Manager Provider](https://external-secrets.io/latest/provider/aws/secrets-manager/)
- [ESO Helm Chart](https://github.com/external-secrets/external-secrets)
- [AWS Secrets Manager Pricing](https://aws.amazon.com/secrets-manager/pricing/)
- [Checkov CKV_K8S_35](https://docs.prismacloud.io/en/enterprise-edition/policy-reference/kubernetes-policies/kubernetes-policy-index/bc-k8s-33)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | ESO sync works, AWS integration real |
| Pemahaman materi | 5 | Secret Reference Pattern, ESO architecture, CKV trade-off |
| Progres sesuai target | 5 | SecretSynced, API healthy, 102/0 Checkov |

---

## ➡️ Rencana Besok

- [ ] Hari 44: AI Threat Modeling pada K8s — topologi + RBAC → AI analisis attack path

---

*[← Hari 42](hari-42.md) | [Hari 44 →](hari-44.md)*
