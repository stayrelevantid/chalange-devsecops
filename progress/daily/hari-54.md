# Hari 54 — Chaos Security Engineering

**📅 Tanggal:** 2026-08-04
**⏱️ Durasi Belajar:** ~90 menit
**🏷️ Fase:** Fase 4 — Vulnerability Management & Red Teaming
**📊 Status:** ✅ Selesai

---

## 🎯 Tujuan Hari Ini

- [x] Uji defense in depth: apakah runtime layer tetap menangkap workload berbahaya jika admission control (OPA Gatekeeper) mati
- [x] Matikan admission control (scale controller-manager → 0)
- [x] Deploy pod non-compliant yang sebelumnya di-deny oleh OPA → buktikan fail-open
- [x] Observasi layer runtime (Falco): yang mana terdeteksi, yang mana lolos
- [x] Cegah/du enemy gap drift detection: OPA audit + script `drift-check.sh`
- [x] Kembalikan sistem, verifikasi kontrol pulih

---

## ✅ Yang Berhasil Dikerjakan

### 1. Baseline (Control Aktif)

Buat `k8s/chaos/test-no-limits.yaml` (pod alpine tanpa requests/limits) → apply di namespace `securebank`:

```
Error from server (Forbidden): admission webhook "validation.gatekeeper.sh" denied:
  [require-resource-limits] Container 'no-limits' has no CPU limit/request
  [require-resource-limits] Container 'no-limits' has no memory limit/request
```

**Bukti:** selama OPA hidup, pod non-compliant diblokir di admission.

### 2. Matikan Lapisan Admission (Chaos)

```bash
kubectl scale deployment gatekeeper-controller-manager -n gatekeeper-system --replicas=0
```

Controller yang serving admission webhook mati. Deployment `gatekeeper-audit` sengaja dibiarkan hidup agar bisa jadi drift detector.

### 3. Fail-Open Terbukti

Admission diekspektasikan "fail-open" (tanpa webhook, request biak). Re-apply manifest + deploy pod privileged:

```bash
kubectl apply -f k8s/chaos/test-no-limits.yaml     # pod/test-no-limits created   (sebelumnya Forbidden!)
kubectl apply -f k8s/chaos/test-privileged.yaml    # pod/test-privileged created
```

Kedua pod Running. **Admission control mati = policy tidak ditegakkan.**

### 4. Layer Runtime (Falco) — defense in depth

| Chaos pod | Falco | Penjelasan |
|-----------|-------|-----------|
| `test-no-limits` (diam, sleep 3600) | **Diam total (0 alert)** | Tidak ada aktivitas syscall aneh → runtime tidak peduli soal "tidak ada limits". Ini gap (policy drift tak terdeteksi runtime) |
| `test-privileged` (privileged) | **CRITICAL** `Privileged container launched` | Meski lolos admission, runtime-capture primitive escalasi. Alert → falcosidekick → webhook → Slack `#security-alerts` (Slack POST 200) |

Bukti key: **walau admission mati, layer runtime (Falco) masih menangkap pod privileged** → eskalasi tetap kelihatan. Tapi workload yang "hanya tidak compliant" (tanpa limits) tidak memicu apa pun di runtime — perlu detector terpisah.

### 5. Drift Detection (2 cara)

**A. OPA Audit (masih hidup meski admission mati):**

```bash
kubectl get k8srequiredlimits.constraints.gatekeeper.sh -o jsonpath='{.items[0].status.totalViolations}'
# 8
kubectl get k8sdisallowedprivileged.constraints.gatekeeper.sh -o jsonpath='{.items[0].status.totalViolations}'
# 1
```

Audit mencatat 8 pelanggaran resource + 1 pelanggaran privileged — **menangkap apa yang admission lewatkan**. Audit jalan periodik (default 60s) dan bukan blocker, tapi memberikan visibilitas.

**B. `drift-check.sh` (script)**

Buat `security/chaos/drift-check.sh` — script kubectl + jq untuk list container yang kurang `requests.cpu/memory` atau `limits.cpu/memory`:

```
=== namespace: securebank ===
DRIFT: pod=test-no-limits container=no-limits req.cpu=MISSING ...
DRIFT: pod=test-privileged container=priv req.cpu=MISSING ...
  ^ 2 container drift ditemukan
RESULT: 2 container drift terdeteksi. ⚠️ ...
```

### 6. Restore & Verifikasi

```bash
kubectl scale deployment gatekeeper-controller-manager -n gatekeeper-system --replicas=1
kubectl delete pod -n securebank test-no-limits test-privileged --grace-period=0 --force
```

Re-apply `test-no-limits.yaml` → **Forbidden lagi** (control pulih). `drift-check.sh` → **0 drift**, semua pod sesuai policy.

---

## 🧪 Hasil Checklist

| Checklist | Hasil |
|-----------|-------|
| OPA Gatekeeper sengaja dimatikan | ✅ controller-manager → 0 |
| Pod non-compliant berhasil di-deploy | ✅ no-limits + privileged (fail-open) |
| Memahami pentingnya layer berlapis (defense in depth) | ✅ Falco tangkap privileged walau admission mati |
| Sistem dikembalikan normal | ✅ control pulih, pod chaos dihapus, 0 drift |

---

## 📝 Catatan Teknis

### Admission ≠ Runtime
Admission control (Gatekeeper) mengecek **pada saat deploy** (admission review). Runtime security (Falco/eBPF) mengecek **saat eksekusi syscall**. Ketika admission mati: workload baru yang buruk lolos; tetapi yang melakukan primitive berbahaya (privileged, setns) tetap muncul di runtime. Yang "hanya buruk secara konfigurasi" (tanpa limits) tidak terlihat di syscall — butuh detector konfigurasi (audit / drift-check).

### Fail-open behavior
`webhookConfig` Gatekeeper dipakai `failurePolicy: Ignore`? Default Gatekeeper memakai `failurePolicy: Ignore` → saat controller mati, klaster **tidak menolak** pod (fail-open), bukan fail-closed. Ini tradeoff lazim: ketersediaan vs keamanan.

### OPA Audit ≠ Admission
Audit berjalan di deployment terpisah (`gatekeeper-audit`) dan secara periodik menandai `total-violations` di status constraint. Itu sebabnya audit tetap mencatat chaos pod — dan itulah detektor drift paling "native" tanpa menambah tool.

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `kubectl get constraints` tidak lagi dikenali (`resource type "constraints"` tidak ada) | Pakai kind penuh: `k8srequiredlimits.constraints.gatekeeper.sh` / `k8sdisallowedprivileged.constraints.gatekeeper.sh` |
| Pod chaos yang sudah di-force-delete tetap Running sebentar | `--grace-period=0 --force` + tunggu; setelah scale balik, state bersih & drift 0 |
| 2 pod gatekeeper stuck `Terminating` 24d (agent nodes NotReady) | Dicatat sebagai temuan saja; bukan blocker scaling |

---

## 📤 Output Hari Ini

- [x] `securebank-api/k8s/chaos/test-no-limits.yaml` — pod non-compliant (tanpa limits)
- [x] `securebank-api/k8s/chaos/test-privileged.yaml` — pod privileged
- [x] `securebank-api/security/chaos/drift-check.sh` — script ubect + jq untuk list pod tanpa limits
- [x] Bukti fail-open + defence-in-depth: Falco menangkap privileged pod → Slack
- [x] Bukti drift detection: OPA audit (9 violation) + `drift-check.sh` (2 drift)
- [x] Restore + verifikasi (Forbidden kembali, 0 drift)

---

## 💡 Lessons Learned

### 1. Admission control adalah satu lapis, bukan segalanya
Saat OPA dimatikan, pod buruk langsung lolos (fail-open). Tapi Falco tetap menandai `privileged`. Artinya klaim "kita aman karena pakai OPA" kurangpas — OPA harus dikombinasi dengan runtime detection. Defense in depth = tidak bergantung pada kegagalan pertama.

### 2. Fail-open itu keputusan desain, bukan bug
Webhook `failurePolicy: Ignore` membuat klaster tetap sehat kala admission error — tapi konsekuensinya policy tidak ditegakkan saat control mati. Saat memilih fail-open, harus ada fallback detector (audit/drift-check/daemon yang resoliant) supaya gap tidak tidak ketahuan.

### 3. Runtime tidak memahami "kebijakan bisnis"
Falco melacak syscall, bukan konfigurasi pod. Pod tanpa limits yang diam tidak pernah ter-trigger — bukan bug Falco, tapi batas area tanggung jawab. "Policy drift" harus diaudit oleh detector konfigurasi.

### 4. Audit itu safety net
Meski admission mati, OPA audit masih menandalah di status constraint. Cek `status.totalViolations` secara berkala = detector drift paling sederhana tanpa tool baru. Integrasikan ke kriteria Prowler/CI kalau serius.

### 5. Latihan chaos menyingkap asumsi
Saya mulus mengira runtime akan melihat "semua yang buruk". Kenyataannya cuma primitive yang terlihat. Chaos test (turn off control) itu cara paling jujur untuk tahu benar luput apa yang attacker/deploy yang salah lakukan.

---

## 🔗 Referensi

- [Day 36](hari-36.md) — OPA testing (pod tanpa limits ditolak pertama kali)
- [Day 52](hari-52.md) — Red Team K8s escape + rule `Privileged container launched`
- [OPA Gatekeeper Audit](https://open-policy-agent.github.io/gatekeeper/website/docs/audit/)
- [Falco event][https://falco.org/docs/events/] — runtime detection
- [K8s Resource Limits](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 5 | E2E: control mati → fail-open → runtime detect → Drift detector |
| Pemahaman materi | 5 | Perbedaan admission vs runtime vs audit, failurePolicy |
| Progres sesuai target | 5 | Semua langkah selesai + sistem pulih |

---

## ➡️ Rencana Besok

- [ ] **Day 55: Laporan Audit** — export DefectDojo → Executive Summary draft

---

*[← Hari 53](hari-53.md) | [Hari 55 →](hari-55.md)*