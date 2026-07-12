# Pro & Kontra — Technology Decisions

> Koleksi pembahasan pro/kontra teknologi yang dipakai di project SecureBank API.
> Setiap section bisa ditambah seiring ada technology decision baru.

---

## Distroless Container Images

Distroless adalah container image yang hanya berisi aplikasi binary + runtime dependencies (CA certs, etc) — tanpa shell, package manager, atau utilities. Dipakai di SecureBank API (`gcr.io/distroless/static-debian12:nonroot`, 7.97MB).

### Pro

| Keuntungan | Detail |
|-----------|--------|
| Attack surface minimal | No shell, no package manager, no utilities. Attacker yang dapat RCE tidak bisa `ls`, `cat`, `curl`, atau pivot |
| Image kecil | 7.97MB vs ~350MB alpine. Push/pull cepat, storage hemat |
| 0 CVE di base image | Tidak ada package yang bisa punya CVE (karena tidak ada package) |
| Immutable by design | Tidak ada util untuk modify filesystem — cocok dengan `readOnlyRootFilesystem: true` |
| Compliance friendly | PCI-DSS, CIS Benchmark suka distroless karena attack surface terukur kecil |

### Kontra

| Masalah | Detail | Mitigation |
|---------|--------|------------|
| **No shell = susah debug** | Tidak bisa `kubectl exec -it pod -- sh` untuk troubleshoot | Pakai **ephemeral debug container** (`kubectl debug`) yang sidecar dengan shell, tanpa restart Pod |
| **No curl/wget** | Tidak bisa test connectivity dari dalam container | Pakai `kubectl debug` atau sidecar dengan tools |
| **No package manager** | Tidak bisa install tools on-the-fly | Semua tools harus dibaked saat build atau via debug container |
| **Log tidak bisa di-tail manual** | Tidak bisa `exec` dan `tail -f /var/log/app.log` | Logging harus ke stdout/stderr (12-factor app) → `kubectl logs` |
| **Crash investigation terbatas** | Tidak bisa inspect filesystem Pod yang crash | Volume mounts + init container untuk pre-crash snapshot, atau `kubectl debug` |
| **Learning curve** | Developer harus build dengan static binary, handle SSL certs, CA bundles | Dokumentasi + CI pipeline yang handle ini |

### Mitigation Terbaik: `kubectl debug`

```bash
# Buat ephemeral container dengan shell di Pod yang sedang running
kubectl debug -it <pod-name> --image=busybox --target=<container-name> -n securebank

# Atau copy Pod dengan debug shell
kubectl debug <pod-name> --image=ubuntu --copy-to=<pod-name>-debug --share-processes
```

Ini bikin container sementara yang share namespace dengan Pod target — bisa inspect filesystem, process, network — **tanpa** merestart Pod atau mengubah distroless image.

### Kesimpulan

Distroless **bukan buruk** — tapi memang menuntut perubahan cara debug:

- Dari "exec ke Pod dan troubleshoot manual" → ke "observability by design" (logging, metrics, tracing, `kubectl debug`)
- Production seharusnya tidak butuh exec ke container untuk debugging — itu adalah anti-pattern
- Debugging di production butuh **ephemeral container**, bukan permanent shell di app image

**Secara security, pro-nya jauh lebih besar dari kontra-nya.** Attack surface minimal + 0 CVE + immutable = sulit untuk attacker. Debugging challenge bisa di-solve dengan `kubectl debug`.

Untuk SecureBank API (financial context), security > convenience.

---

<!-- 
Template untuk section baru:

## Technology Name

Deskripsi singkat teknologi dan kenapa dipakai di project.

### Pro

| Keuntungan | Detail |
|-----------|--------|
| ... | ... |

### Kontra

| Masalah | Detail | Mitigation |
|---------|--------|------------|
| ... | ... | ... |

### Kesimpulan

...

---
-->