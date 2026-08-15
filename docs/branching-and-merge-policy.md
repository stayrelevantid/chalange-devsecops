# Branching & Merge Policy

Dokumen ini menjelaskan cara menjaga history repository tetap rapi dan promotion antar-environment tidak mudah konflik.

## Branch Model

```text
feature/* ─┐
           ├─> develop ──> staging ──> main
fix/* ─────┘
```

Mapping deployment:

```text
develop -> dev
staging -> staging
main    -> prod
```

Jalur PR yang diizinkan:

| Source | Target |
|--------|--------|
| `feature/*` | `develop` |
| `fix/*` | `develop` atau `staging` |
| `develop` | `staging` |
| `staging` | `main` |

Feature tidak boleh langsung ke `staging` atau `main`. Fix tidak boleh langsung ke `main`.

## Rebase vs Merge

### Rebase untuk feature/fix pribadi

Branch pribadi selalu dibuat dari `main` terbaru dan direbase sebelum PR:

```bash
git fetch origin
git switch main
git pull --ff-only origin main
git switch -c feature/nama-fitur

# sebelum membuka atau memperbarui PR
git fetch origin
git rebase origin/main
git push --force-with-lease origin feature/nama-fitur
```

Rebase hanya dilakukan pada branch pribadi yang tidak sedang dipakai bersama. Jangan rebase `main`, `develop`, atau `staging` secara langsung.

### Merge commit untuk promotion

Promotion PR wajib memakai **Create a merge commit**:

```text
develop -> staging
staging -> main
```

Merge commit mempertahankan ancestry branch sehingga jalur promotion dapat ditelusuri. Jangan gunakan Squash and merge untuk promotion karena commit asli akan diganti satu commit baru dan branch dapat terlihat divergen pada promotion berikutnya.

Repository mengizinkan merge commit sebagai metode utama. GitHub tetap mengizinkan squash sebagai fallback repository karena GitHub tidak menyediakan konfigurasi API untuk menonaktifkan merge commit sambil mempertahankan hanya satu metode; secara operasional, squash dilarang untuk promotion PR.

## Prosedur Promotion

1. Pastikan `main` terbaru.
2. Buat `feature/*` atau `fix/*` dari `main`.
3. Rebase branch pribadi ke `origin/main` sebelum PR.
4. Buat PR ke `develop` dan tunggu seluruh CI green.
5. Merge dengan **Create a merge commit**.
6. Buat PR `develop -> staging`; merge dengan **Create a merge commit**.
7. Buat PR `staging -> main`; merge dengan **Create a merge commit**.

Auto-CD berjalan setelah CI sukses:

```text
merge ke develop -> deploy dev
merge ke staging -> deploy staging
merge ke main    -> deploy prod
```

## Checklist Sebelum Merge

- [ ] Source dan target PR sesuai promotion policy.
- [ ] Feature/fix sudah rebase ke `origin/main`.
- [ ] CI wajib lulus: build/test, secret scan, SCA, SAST, DAST, IaC, image scan, dan Security Gate.
- [ ] Conflict status bersih.
- [ ] Promotion PR menggunakan **Create a merge commit**, bukan squash.
- [ ] Tidak ada secret, key, kubeconfig, atau credential baru yang masuk repository.

## Verifikasi Ancestry

Setelah promotion, periksa hubungan branch:

```bash
git fetch origin
git merge-base --is-ancestor origin/main origin/develop
git merge-base --is-ancestor origin/develop origin/staging
```

Exit code `0` berarti branch sumber merupakan ancestor target. Branch tidak harus memiliki hash yang identik; yang penting history promotion dapat ditelusuri dan tidak saling divergen tanpa merge commit.

## Jika Terjadi Konflik

Jangan force-push branch protected. Buat branch resolusi dari target branch, gabungkan source branch, resolve konflik dengan mempertahankan konfigurasi terbaru, lalu buka PR baru ke target:

```bash
git fetch origin
git switch -c fix/resolve-promotion origin/staging
git merge origin/develop
# resolve konflik, validasi, lalu:
git add .
git commit -m "fix: resolve promotion conflict"
git push -u origin fix/resolve-promotion
```

Pastikan workflow terbaru, Kustomize overlay `dev/staging/prod`, dan promotion policy tetap utuh sebelum merge.
