# Fase 3: Kubernetes & Runtime Security (Hari 31–45)

> **Proyek:** SecureBank API — Deploy ke K8s, hardening cluster, policy enforcement, runtime monitoring
> **Output Fase:** Cluster K8s yang hardened dengan OPA Gatekeeper, Network Policies, Falco monitoring, dan alerting.

---

## Hari 31: K8s Cluster Setup & Deploy

### Tujuan
Membuat cluster k3d lokal dan deploy SecureBank API.

### Tutorial

**1. Install k3d:**
```bash
brew install k3d
```

**2. Buat cluster:**
```bash
k3d cluster create securebank \
  --port "8080:80@loadbalancer" \
  --agents 2
```

**3. Buat `k8s/deployment.yaml`:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: securebank-api
  namespace: securebank
  labels:
    app: securebank-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: securebank-api
  template:
    metadata:
      labels:
        app: securebank-api
    spec:
      containers:
        - name: api
          image: securebank:v1
          ports:
            - containerPort: 8080
          # Intentionally missing: securityContext, resources
          # Will be fixed on Day 33
```

**4. Buat `k8s/service.yaml`:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: securebank-api
  namespace: securebank
spec:
  selector:
    app: securebank-api
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
```

**5. Deploy:**
```bash
kubectl create namespace securebank
k3d image import securebank:v1 -c securebank
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl get pods -n securebank
```

**6. Port-forward untuk test:**
```bash
kubectl port-forward svc/securebank-api 8080:80 -n securebank &
curl http://localhost:8080/health
```

### Checklist
- [ ] k3d cluster running dengan 2 agents
- [ ] SecureBank API deployed (2 replicas)
- [ ] Service bisa diakses via port-forward
- [ ] `/health` endpoint merespons

---

## Hari 32: K8s Misconfiguration Scanning

### Tujuan
Memindai manifest Kubernetes untuk menemukan konfigurasi berisiko.

### Tutorial

**1. Install Kubesec:**
```bash
brew install kubesec
# atau gunakan online API
```

**2. Scan deployment:**
```bash
# Kubesec (lokal atau API)
kubesec scan k8s/deployment.yaml

# Checkov untuk K8s
checkov -f k8s/deployment.yaml --framework kubernetes

# Trivy config
trivy config k8s/
```

**3. Temuan yang diharapkan:**
- Container running as root
- No resource limits/requests
- No readOnlyRootFilesystem
- No securityContext defined
- allowPrivilegeEscalation not set to false

**4. Dokumentasikan temuan:**
```bash
checkov -f k8s/deployment.yaml --framework kubernetes \
  --output json > security/k8s-scan-report.json
```

### Checklist
- [ ] Kubesec + Checkov + Trivy scan berjalan
- [ ] ≥ 5 misconfiguration ditemukan
- [ ] Report JSON disimpan
- [ ] Temuan dipahami (apa risikonya)

---

## Hari 33: SecurityContext Hardening

### Tujuan
Memperbaiki deployment.yaml dengan security best practices.

### Tutorial

**Update `k8s/deployment.yaml`:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: securebank-api
  namespace: securebank
  labels:
    app: securebank-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: securebank-api
  template:
    metadata:
      labels:
        app: securebank-api
    spec:
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
        runAsGroup: 65534
        fsGroup: 65534
        seccompProfile:
          type: RuntimeDefault
      containers:
        - name: api
          image: securebank:v1
          ports:
            - containerPort: 8080
              protocol: TCP
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          resources:
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 200m
              memory: 128Mi
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 3
            periodSeconds: 5
```

**Re-deploy dan re-scan:**
```bash
kubectl apply -f k8s/deployment.yaml
checkov -f k8s/deployment.yaml --framework kubernetes
# Semua check PASSED
```

### Checklist
- [ ] `readOnlyRootFilesystem: true`
- [ ] `allowPrivilegeEscalation: false`
- [ ] `capabilities.drop: ["ALL"]`
- [ ] `runAsNonRoot: true`
- [ ] Resource limits/requests defined
- [ ] Liveness + readiness probes
- [ ] Checkov: 0 FAILED
- [ ] Pod tetap berjalan normal

---

## Hari 34: OPA Gatekeeper Setup

### Tujuan
Menginstall OPA Gatekeeper sebagai Admission Controller di cluster.

### Tutorial

**1. Add Helm repo dan install:**
```bash
helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
helm repo update

helm install gatekeeper gatekeeper/gatekeeper \
  --namespace gatekeeper-system \
  --create-namespace \
  --set replicas=1
```

**2. Verifikasi:**
```bash
kubectl get pods -n gatekeeper-system
kubectl get crd | grep gatekeeper
# Harus ada: constrainttemplates.templates.gatekeeper.sh
```

### Checklist
- [ ] Gatekeeper pods running
- [ ] CRDs terinstall
- [ ] Webhook aktif (validating admission webhook)

---

## Hari 35: Rego Policy — Wajib Resource Limits

### Tujuan
Menulis policy Rego yang menolak Pod tanpa resource limits.

### Tutorial

**1. Buat `k8s/gatekeeper/constraint-templates/require-limits.yaml`:**
```yaml
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlimits
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredLimits
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredlimits

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          not container.resources.limits.cpu
          msg := sprintf("Container '%v' has no CPU limit", [container.name])
        }

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          not container.resources.limits.memory
          msg := sprintf("Container '%v' has no memory limit", [container.name])
        }
```

**2. Buat `k8s/gatekeeper/constraints/require-limits.yaml`:**
```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sRequiredLimits
metadata:
  name: require-resource-limits
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
      - apiGroups: ["apps"]
        kinds: ["Deployment"]
    namespaces:
      - securebank
```

**3. Apply:**
```bash
kubectl apply -f k8s/gatekeeper/constraint-templates/require-limits.yaml
kubectl apply -f k8s/gatekeeper/constraints/require-limits.yaml
```

### Checklist
- [ ] ConstraintTemplate applied
- [ ] Constraint applied di namespace `securebank`
- [ ] Policy siap diuji (hari 36)

---

## Hari 36: OPA Policy Testing

### Tujuan
Memvalidasi bahwa Gatekeeper menolak Pod tanpa resource limits.

### Tutorial

**1. Buat pod tanpa limits:**
```yaml
# k8s/test-no-limits.yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-no-limits
  namespace: securebank
spec:
  containers:
    - name: nginx
      image: nginx:latest
      # NO resources defined
```

**2. Coba apply:**
```bash
kubectl apply -f k8s/test-no-limits.yaml
# ERROR: admission webhook "validation.gatekeeper.sh" denied the request:
# Container 'nginx' has no CPU limit
```

**3. Verifikasi constraints:**
```bash
kubectl get constraints
kubectl describe k8srequiredlimits require-resource-limits
```

**4. Hapus test file, commit policy:**
```bash
rm k8s/test-no-limits.yaml
git add k8s/gatekeeper/
git commit -m "security: add OPA Gatekeeper policy requiring resource limits"
```

### Checklist
- [ ] Pod tanpa limits DITOLAK oleh Gatekeeper
- [ ] Error message jelas dan informatif
- [ ] Pod dengan limits (deployment.yaml) tetap diterima
- [ ] Policy di-commit ke repo

---

## Hari 37: Network Policies

### Tujuan
Menerapkan Default Deny dan whitelist traffic yang diizinkan.

### Tutorial

**1. Buat `k8s/network-policy.yaml`:**
```yaml
# Default Deny All Ingress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: securebank
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress

---
# Allow Ingress to API only from ingress controller
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-api-ingress
  namespace: securebank
spec:
  podSelector:
    matchLabels:
      app: securebank-api
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
      ports:
        - protocol: TCP
          port: 8080

---
# Allow DNS egress (semua pod butuh DNS)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns-egress
  namespace: securebank
spec:
  podSelector: {}
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector: {}
      ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
```

**2. Apply:**
```bash
kubectl apply -f k8s/network-policy.yaml
```

**3. Test — dari pod lain, coba akses API:**
```bash
# Ini harus GAGAL (denied by network policy)
kubectl run test-curl --rm -it --image=curlimages/curl \
  -n default -- curl http://securebank-api.securebank.svc:80
```

### Checklist
- [ ] Default Deny All diterapkan
- [ ] API hanya bisa diakses dari ingress controller
- [ ] DNS egress diizinkan
- [ ] Cross-namespace access diblokir

---

## Hari 38: RBAC Auditing

### Tujuan
Mengaudit dan memperketat izin akses di cluster.

### Tutorial

**1. Install rakkess:**
```bash
kubectl krew install access-matrix
```

**2. Audit:**
```bash
# Siapa yang bisa create pods?
kubectl who-can create pods -n securebank

# Access matrix untuk service account default
kubectl access-matrix -n securebank --sa default
```

**3. Buat `k8s/rbac.yaml` — service account khusus:**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: securebank-api
  namespace: securebank
  automountServiceAccountToken: false

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: securebank-api-role
  namespace: securebank
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get"]
    resourceNames: ["securebank-config"]
  # Minimal permissions — hanya baca configmap spesifik

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: securebank-api-binding
  namespace: securebank
subjects:
  - kind: ServiceAccount
    name: securebank-api
    namespace: securebank
roleRef:
  kind: Role
  name: securebank-api-role
  apiGroup: rbac.authorization.k8s.io
```

**4. Update deployment untuk menggunakan SA baru:**
```yaml
spec:
  template:
    spec:
      serviceAccountName: securebank-api
```

### Checklist
- [ ] RBAC diaudit dengan rakkess/who-can
- [ ] Service account khusus dibuat (least privilege)
- [ ] Deployment menggunakan SA baru
- [ ] Default SA tidak digunakan

---

## Hari 39: Runtime Security — Falco Setup

### Tujuan
Menginstall Falco untuk monitoring ancaman runtime.

### Tutorial

**1. Install Falco via Helm:**
```bash
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update

helm install falco falcosecurity/falco \
  --namespace falco \
  --create-namespace \
  --set falcosidekick.enabled=true \
  --set falcosidekick.webui.enabled=true
```

**2. Verifikasi:**
```bash
kubectl get pods -n falco
kubectl logs -l app.kubernetes.io/name=falco -n falco --tail=20
```

**3. Lihat default rules yang aktif:**
```bash
kubectl exec -n falco $(kubectl get pod -n falco -l app.kubernetes.io/name=falco -o name | head -1) \
  -- cat /etc/falco/falco_rules.yaml | head -50
```

### Checklist
- [ ] Falco pods running
- [ ] Falcosidekick enabled
- [ ] Default rules aktif dan menghasilkan log
- [ ] Memahami format output Falco

---

## Hari 40: Falco Custom Rules

### Tujuan
Membuat aturan Falco kustom yang relevan untuk SecureBank.

### Tutorial

**1. Buat `security/falco-rules/securebank-rules.yaml`:**
```yaml
- rule: Shell spawned in SecureBank container
  desc: Detect shell execution in securebank containers
  condition: >
    spawned_process and
    container and
    container.image.repository contains "securebank" and
    proc.name in (bash, sh, ash, zsh)
  output: >
    ALERT: Shell opened in SecureBank container
    (user=%user.name container=%container.name
    image=%container.image.repository command=%proc.cmdline
    pod=%k8s.pod.name namespace=%k8s.ns.name)
  priority: WARNING
  tags: [securebank, shell]

- rule: Sensitive file read in SecureBank
  desc: Detect access to sensitive paths
  condition: >
    open_read and
    container and
    container.image.repository contains "securebank" and
    (fd.name startswith /etc/shadow or
     fd.name startswith /etc/passwd or
     fd.name startswith /proc/self/environ)
  output: >
    ALERT: Sensitive file read in SecureBank
    (file=%fd.name user=%user.name container=%container.name
    pod=%k8s.pod.name)
  priority: CRITICAL
  tags: [securebank, filesystem]

- rule: Network tool in SecureBank container
  desc: Detect network reconnaissance tools
  condition: >
    spawned_process and
    container and
    container.image.repository contains "securebank" and
    proc.name in (nc, ncat, nmap, wget, curl, dig)
  output: >
    ALERT: Network tool detected in SecureBank
    (tool=%proc.name user=%user.name pod=%k8s.pod.name)
  priority: NOTICE
  tags: [securebank, network]
```

**2. Apply via Helm upgrade:**
```bash
helm upgrade falco falcosecurity/falco \
  --namespace falco \
  --set-file customRules."securebank-rules\.yaml"=security/falco-rules/securebank-rules.yaml
```

### Checklist
- [ ] 3 custom rules dibuat
- [ ] Rules di-load oleh Falco
- [ ] Rules relevan untuk konteks perbankan

---

## Hari 41: Falco Attack Simulation

### Tujuan
Mensimulasikan serangan dan memastikan Falco mendeteksinya.

### Tutorial

**1. Exec ke container (simulasi attacker):**
```bash
# Ini akan trigger "Shell spawned" rule
kubectl exec -it $(kubectl get pod -n securebank -l app=securebank-api -o name | head -1) \
  -n securebank -- /bin/sh
```

*Catatan: jika menggunakan distroless, shell tidak tersedia — ini membuktikan keamanan distroless! Deploy pod test dengan alpine jika perlu:*

```bash
kubectl run attacker-sim --rm -it \
  -n securebank --image=alpine -- /bin/sh
# Di dalam container:
cat /etc/passwd
wget http://securebank-api/health
```

**2. Cek Falco logs:**
```bash
kubectl logs -l app.kubernetes.io/name=falco -n falco --tail=50 | grep "ALERT"
```

**3. Verifikasi deteksi:**
- [ ] Shell spawned → WARNING
- [ ] Sensitive file read → CRITICAL (jika custom rule aktif)
- [ ] Network tool → NOTICE

**4. Cleanup:**
```bash
kubectl delete pod attacker-sim -n securebank --ignore-not-found
```

### Checklist
- [ ] Simulasi serangan berhasil dilakukan
- [ ] Falco mendeteksi semua aktivitas
- [ ] Log level sesuai (WARNING/CRITICAL)
- [ ] Distroless membuktikan shell tidak tersedia di production container

---

## Hari 42: Automated Alerting — n8n Webhook

### Tujuan
Menghubungkan alert Falco ke sistem automasi n8n.

### Tutorial

**1. Deploy n8n lokal:**
```bash
docker run -d --name n8n \
  -p 5678:5678 \
  -v n8n_data:/home/node/.n8n \
  n8nio/n8n
```

**2. Buka `http://localhost:5678`, buat workflow baru.**

**3. Tambahkan node:**
- **Webhook** trigger (method: POST, path: `/falco-alert`)
- **IF** node: cek `body.priority` == `Critical`
- **True branch**: (akan dihubungkan ke Slack di hari 48)
- **False branch**: log ke file

**4. Konfigurasi Falcosidekick untuk kirim ke n8n:**
```bash
helm upgrade falco falcosecurity/falco \
  --namespace falco \
  --set falcosidekick.config.webhook.address="http://host.k3d.internal:5678/webhook/falco-alert"
```

**5. Test: trigger ulang simulasi hari 41 → cek n8n mendapat data.**

### Checklist
- [ ] n8n berjalan dan workflow dibuat
- [ ] Webhook endpoint siap menerima JSON
- [ ] Falcosidekick mengirim alert ke n8n
- [ ] Data alert masuk ke n8n workflow

---

## Hari 43: K8s Secret Management — External Secrets

### Tujuan
Menarik secrets dari AWS Secrets Manager ke K8s secara aman.

### Tutorial

**1. Install External Secrets Operator:**
```bash
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets \
  --namespace external-secrets \
  --create-namespace
```

**2. Buat SecretStore (contoh AWS):**
```yaml
# k8s/external-secrets/secret-store.yaml
apiVersion: external-secrets.io/v1beta1
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
            name: aws-credentials
            key: access-key-id
          secretAccessKeySecretRef:
            name: aws-credentials
            key: secret-access-key
```

**3. Buat ExternalSecret:**
```yaml
# k8s/external-secrets/db-secret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: securebank-db
  namespace: securebank
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets
    kind: SecretStore
  target:
    name: securebank-db-credentials
    creationPolicy: Owner
  data:
    - secretKey: DB_PASSWORD
      remoteRef:
        key: securebank/db-password
```

**4. Update deployment untuk mount secret:**
```yaml
env:
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: securebank-db-credentials
        key: DB_PASSWORD
```

### Checklist
- [ ] External Secrets Operator terinstall
- [ ] SecretStore mengarah ke AWS Secrets Manager
- [ ] ExternalSecret membuat K8s Secret secara otomatis
- [ ] Deployment membaca secret dari environment variable

---

## Hari 44: AI Threat Modeling pada K8s

### Tujuan
Menggunakan AI untuk menganalisis attack surface cluster.

### Tutorial

**1. Kumpulkan data cluster:**
```bash
kubectl get networkpolicies -A -o yaml > security/threat-model/netpol.yaml
kubectl get roles,rolebindings -A -o yaml > security/threat-model/rbac.yaml
kubectl get pods -A -o yaml > security/threat-model/pods.yaml
```

**2. Prompt ke AI:**
```
Kamu adalah Kubernetes Security Specialist. Analisis konfigurasi cluster berikut:

[paste netpol, rbac, pods yaml]

Identifikasi:
1. Top 3 jalur serangan (attack paths) — misalnya: pod escape, lateral movement, privilege escalation
2. Untuk setiap jalur: prasyarat, langkah eksploitasi, dan mitigasi
3. Skor risiko berdasarkan CVSS framework
4. Rekomendasi hardening tambahan
```

**3. Dokumentasikan di `security/threat-model/k8s-threat-analysis.md`.**

### Checklist
- [ ] Data cluster dikumpulkan
- [ ] AI menganalisis 3 attack paths
- [ ] Mitigasi didokumentasikan
- [ ] Threat model K8s tersimpan

---

## Hari 45: Dokumentasi Fase 3

### Tujuan
Merangkum semua K8s security yang diterapkan.

### Tutorial
Tulis `docs/fase-3-k8s-runtime.md`:
- Arsitektur cluster (diagram)
- SecurityContext checklist
- OPA Gatekeeper policies dan efeknya
- Network Policies topology
- Falco rules dan alert flow
- RBAC audit findings
- Metrics: jumlah policy violations blocked

### Checklist
- [ ] Dokumen lengkap
- [ ] Diagram arsitektur disertakan
- [ ] Commit: `docs: fase 3 kubernetes and runtime security`

---

> ✅ **Selesai Fase 3** — Lanjut ke [Fase 4: Vulnerability Management & Red Teaming](fase-4-vuln-redteam.md)
