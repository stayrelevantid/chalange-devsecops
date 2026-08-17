# Panduan Persiapan Ujian CDP (Certified DevSecOps Professional)

Dokumen ini merangkum strategi, skill, cheatsheet, dan tips-trik untuk menghadapi simulasi ujian **Certified DevSecOps Professional (CDP)** yang meminta Anda membangun pipeline DevSecOps dari nol dalam 3 jam.

Dibangun dari pengalaman 60 hari membangun SecureBank API dan diintegrasikan dengan 15 best practices dari [Practical DevSecOps — DevSecOps Best Practices](https://www.practical-devsecops.com/devsecops-best-practices/).

---

## 1. Format Ujian (Simulasi 3 Jam)

- Tulis `.github/workflows/ci.yml` dari nol.
- Tambahkan Secret Scan (Gitleaks).
- Tambahkan SCA (Trivy).
- Tambahkan SAST (Semgrep).
- Buat Dockerfile multi-stage.
- Tulis manifest K8s (deployment & service) dengan SecurityContext yang aman.

### Aturan
- Boleh mencari di Google / dokumentasi resmi.
- Dilarang menyalin dari branch `main` repositori Anda sendiri.
- Gunakan branch terpisah (mis. `cdp-exam`) agar tidak merusak repo utama.

---

## 2. Strategi Alokasi Waktu (180 Menit)

| Blok | Durasi | Isi |
|------|--------|-----|
| Perencanaan & struktur | 15 menit | Petakan file yang harus dibuat, urutan pengerjaan, siapkan environment |
| CI pipeline + build/test | 45 menit | Workflow skeleton, build Go, test |
| Secret scan + SCA | 25 menit | Gitleaks, Trivy dependency |
| SAST | 20 menit | Semgrep |
| Dockerfile multi-stage | 25 menit | Build hardened, non-root |
| K8s manifests aman | 30 menit | Deployment + Service, SecurityContext, probes |
| Verifikasi & buffer | 20 menit | Test run, perbaiki error, pastikan pipeline hijau |

> Prinsip: **jangan optimalkan dulu**. Bangun versi bekerja, verifikasi hijau, baru polish.

---

## 3. 15 DevSecOps Best Practices (Practical DevSecOps)

Berikut 15 best practices yang dipetakan ke tindakan konkret saat ujian:

| # | Best Practice | Tindakan Konkret di Ujian |
|---|---------------|---------------------------|
| 1 | **Shift Left** | Mulai scan sejak commit pertama, bukan di akhir. Letakkan security check lebih awal di workflow. |
| 2 | **Adopt Automation** | Semua scan berjalan otomatis via GitHub Actions, tanpa klik manual. |
| 3 | **Continuous Testing** | Setiap push/PR memicu test + scan; bukan hanya saat release. |
| 4 | **Prioritize Risk Management** | Atur severity gate: CRITICAL/HIGH memblokir, MEDIUM/LOW ditoleransi. |
| 5 | **Integrate Security Tools** | Gitleaks + Trivy + Semgrep dalam satu pipeline yang konsisten. |
| 6 | **Collaborate Across Teams** | Dokumentasikan keputusan di repo (policy as docs), gunakan PR sebagai review. |
| 7 | **Implement Secure Coding Standards** | Dockerfile: non-root, `CGO_ENABLED=0`, least privilege; K8s: drop ALL capabilities. |
| 8 | **Enforce Access Controls** | Secrets tidak pernah di-commit; gunakan GitHub Secrets / environment. |
| 9 | **Monitor for Threats** | DAST (ZAP) dan image scan sebagai lapisan deteksi. |
| 10 | **Provide Security Training** | Catat "lessons learned" pasca ujian untuk perbaikan berkelanjutan. |
| 11 | **Embrace Policy as Code** | Branch protection + required status checks sebagai kebijakan yang bisa diverifikasi. |
| 12 | **Utilize Threat Modeling** | Pikirkan ancaman per komponen (API, container, K8s, cloud) sebelum menulis config. |
| 13 | **Expand Incident Response** | Pastikan pipeline bisa memblokir release berbahaya; siapkan cara rollback. |
| 14 | **Leverage Immutable Infrastructure** | Image di-build sekali dan di-tag; jangan mutasi container saat runtime. |
| 15 | **Enhance Security Observability** | Upload report sebagai artifacts; verifikasi hasil scan terlihat dan bisa diaudit. |

---

## 4. Checklist Skill yang Harus Dikuasai

### GitHub Actions
- Syntax workflow: `name`, `on`, `jobs`, `runs-on`, `steps`, `uses`, `run`, `with`.
- Context: `github.sha`, `github.ref`, `github.event_name`, `secrets.*`, `needs.*.outputs`.
- Dependency antar job dengan `needs`.
- Artifacts & cache (`actions/upload-artifact@v4`, `actions/cache@v4`).
- `workflow_dispatch` dan `workflow_run`.

### Gitleaks
- Instal & jalankan di runner: `gitleaks detect --source . --report-path ... --report-format json`.
- `fetch-depth: 0` agar scan seluruh history.
- Custom rules & `.gitleaks.toml`.

### Trivy
- SCA: `trivy fs --scanners vuln --exit-code 1 --severity CRITICAL,HIGH .`.
- IaC: `trivy config --exit-code 1 --severity CRITICAL,HIGH,MEDIUM .`.
- Image: `trivy image --exit-code 1 --severity CRITICAL,HIGH <image>`.
- Update DB sebelum scan: `trivy image --download-db-only`.

### Semgrep
- Rulesets: `--config p/golang`, `--config p/owasp-top-ten`.
- Custom rules `.semgrep.yml`.
- Exit-code untuk memblokir temuan HIGH.

### Dockerfile Multi-Stage
- Builder stage: `FROM golang:1.26.6-alpine AS builder`.
- Runtime stage: `FROM gcr.io/distroless/static-debian12` atau `alpine`.
- `CGO_ENABLED=0`, `GOOS=linux`, `USER 65532:65532`.
- `.dockerignore` untuk mengecualikan file sensitif.

### Kubernetes
- Deployment: image, ports, liveness/readiness probes, resource requests/limits.
- Service: selector, ports, type.
- SecurityContext: `runAsNonRoot: true`, `runAsUser`, `allowPrivilegeEscalation: false`, `capabilities.drop: ["ALL"]`, `readOnlyRootFilesystem: true`, `seccompProfile`.

### Bash
- `set -euo pipefail`.
- Loops, conditional, `$GITHUB_OUTPUT`, env substitution.

---

## 5. Cheatsheet Sintaks (Komentar dalam Bahasa Inggris)

### GitHub Actions skeleton

```yaml
name: SecureBank CI

on:
  push:
    branches: [main, develop, staging]
  pull_request:
    branches: [main, develop, staging]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build
        run: go build ./...
      - name: Test with race detection
        run: go test -race -cover ./...
```

### Gitleaks

```yaml
      - name: Run Gitleaks
        uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Atau manual:

```bash
# Scan repo including full history
gitleaks detect --source . --report-path gitleaks-report.json --report-format json
# Fail on any finding
gitleaks detect --source . --no-banner
```

### Trivy SCA

```yaml
      - name: Run Trivy SCA
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: fs
          scan-ref: .
          severity: CRITICAL,HIGH
          exit-code: '1'
          format: json
          output: trivy-sca.json
```

### Semgrep

```bash
# Scan with community rulesets
semgrep scan --config p/golang --config p/owasp-top-ten --json > semgrep.json
# Fail only on high severity
semgrep scan --config p/golang --error --severity ERROR --json > semgrep.json
```

### Dockerfile multi-stage hardened

```dockerfile
# Build stage
FROM golang:1.26.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o securebank .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/securebank /usr/local/bin/securebank
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["securebank"]
```

### Kubernetes Deployment aman

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: securebank-api
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
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
        runAsGroup: 65532
        seccompProfile:
          type: RuntimeDefault
      containers:
        - name: api
          image: securebank:latest
          ports:
            - containerPort: 8080
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
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
```

### Kubernetes Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: securebank-api
spec:
  selector:
    app: securebank-api
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
```

### Validasi manifest tanpa cluster (kubeconform)

```bash
# Install kubeconform
curl -sSLo /tmp/kubeconform.tar.gz https://github.com/yannh/kubeconform/releases/download/v0.7.0/kubeconform-linux-amd64.tar.gz
tar -xzf /tmp/kubeconform.tar.gz -C /tmp kubeconform
sudo mv /tmp/kubeconform /usr/local/bin/

# Validate against Kubernetes schema (no cluster required)
kubeconform -strict -ignore-missing-schemas manifest.yaml
```

---

## 6. Tips, Trik, dan Jebakan Umum

### Jebakan yang Sering Muncul
1. **`kubectl apply --dry-run` butuh cluster.** Di GitHub-hosted runner tidak ada Kubernetes API server → akan gagal konek ke `localhost:8080`. Gunakan `kubeconform -strict` sebagai gantinya.
2. **Gitleaks tidak melihat history** jika `fetch-depth` default (1). Gunakan `fetch-depth: 0` agar seluruh commit ter-scan.
3. **Trivy image di runner butuh image sudah ada.** Build image dulu (`docker build`), lalu scan tag lokal, atau login ke registry sebelum `trivy image`.
4. **Trivy DB tidak ter-update.** Jalankan `trivy image --download-db-only` atau gunakan `--db-update` agar data CVE terbaru.
5. **`commonLabels` deprecated di Kustomize.** Gunakan format `labels:` (list of pairs) agar tidak ada warning.
6. **CVE di Go stdlib.** Selalu pin versi patch terbaru (mis. `golang:1.26.6-alpine`) untuk menutup CVE `CRITICAL/HIGH` di standard library.
7. **Artifact path Semgrep tidak konsisten** dengan upload path workflow → upload gagal. Pastikan nama file report konsisten antara job scan dan job upload.
8. **Cache key yang salah** membuat dependency di-download ulang setiap run dan scan jadi lambat. Gunakan hash file lockfile sebagai cache key.

### Trik Efisiensi
- **Siapkan snippet sebelum ujian** — cheatsheet di atas bisa dihafalkan dalam bentuk template.
- **Gunakan satu workflow terlebih dahulu** sampai hijau, baru expand. Jangan tulis 10 job sekaligus tanpa tes.
- **Test lokal sebelum push**: `go build`, `gitleaks detect` di mesin sendiri agar error cepat terlihat.
- **Gunakan `continue-on-error` hanya untuk evidence** (mis. DAST), bukan untuk gate utama.
- **Verifikasi log run**, bukan hanya status hijau — pastikan tiap scan benar-benar men-scan target yang dimaksud.

---

## 7. Checklist Persiapan Sebelum Ujian

- [ ] Laptop & environment stabil (git, docker, Go terinstall).
- [ ] GitHub token dengan akses repo & Actions.
- [ ] Snapshot/reference "answer key" disimpan terpisah (tidak dibuka saat ujian).
- [ ] Branch `cdp-exam` sudah disiapkan (atau repo bersih untuk latihan).
- [ ] Stopwatch 3 jam siap.
- [ ] Aturan dicatat: boleh Google, dilarang salin dari `main`.
- [ ] Buffer 20 menit untuk verifikasi sudah dialokasikan.

---

## 8. Strategi Verifikasi (Tanpa Mencontek Jawaban)

- **Ulangi formula dari ingatan**, bukan dari file.
- **Gunakan dokumentasi resmi** (GitHub Actions docs, Trivy docs, Semgrep docs).
- **Jalankan scan lokal** dan bandingkan hasilnya dengan harapan.
- **Pastikan severity gate masuk akal**: `CRITICAL,HIGH` untuk SCA/image, `--error` untuk SAST.
- **Test workflow di repo latihan** (repo kosong/sandbox) sebelum deploy ke repo utama.

---

## 9. Evaluasi Diri Pasca Ujian

- Bagian mana yang paling sering lupa syntax-nya?
- Script/step mana yang paling lama dan bisa dipercepat?
- Apakah severity gate sudah sesuai dan tidak terlalu longgar?
- Apakah semua artifact report bisa diakses dan diaudit?
- Apa 3 perbaikan teratas untuk ujian berikutnya?

---

*Panduan ini terkait dengan [Branching & Merge Policy](branching-and-merge-policy.md) dan catatan harian [Hari 59](../progress/daily/hari-59.md).*
