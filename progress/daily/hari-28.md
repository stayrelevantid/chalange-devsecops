# Hari 28 — Compliance as Code (Chef InSpec)

**📅 Tanggal:** 2026-07-05  
**⏱️ Durasi Belajar:** ~1 jam  
**🏷️ Fase:** Fase 2 — IaC & Container Security  
**📊 Status:** ✅ Selesai  

---

## 🎯 Tujuan Hari Ini

- [x] Install Chef InSpec
- [x] Buat InSpec profile untuk SecureBank API
- [x] Tulis 3 compliance controls
- [x] Run InSpec — semua controls pass

---

## ✅ Yang Berhasil Dikerjakan

- Install InSpec v5.22.3 via `brew install --cask chef/chef/inspec`
- Build + run Go binary di local machine (port 8080)
- Buat InSpec profile di `security/inspec-profiles/securebank/`
- 3 controls: SSH port closed, API port listening, no root process
- Run InSpec: **3 successful controls, 0 failures, 0 skipped** ✅

---

## 📝 Catatan Teknis

### InSpec Profile Structure
```
security/inspec-profiles/securebank/
├── inspec.yml          # Profile metadata
├── inspec.lock        # Dependency lock (auto-generated)
└── controls/
    ├── network.rb     # ssh-port-closed + api-port-listening
    └── process.rb    # no-root-process
```

### Controls

| Control | Impact | Matcher | Result |
|---------|--------|---------|--------|
| ssh-port-closed | 1.0 | `port(22) should_not be_listening` | PASS ✅ |
| api-port-listening | 0.7 | `port(8080) should be_listening` | PASS ✅ |
| no-root-process | 1.0 | `processes('securebank') users should_not include 'root'` | PASS ✅ |

### Run Command
```bash
inspec exec security/inspec-profiles/securebank/ --chef-license accept
```

---

## 🚧 Hambatan & Solusi

| Hambatan | Solusi / Workaround |
|----------|---------------------|
| `gem install inspec` tidak menyertakan CLI executable | InSpec 7.x gem tidak ship CLI — hanya library. Pakai `brew install --cask chef/chef/inspec` (installer pkg, butuh sudo) |
| `brew install --cask chef/chef/inspec` butuh sudo password (gagal di background) | User install manual di terminal dengan sudo |
| Chef InSpec license acceptance required |Tambah `--chef-license accept` flag di command pertama kali |

---

## 📤 Output Hari Ini

- [x] InSpec profile: 3 files (inspec.yml, network.rb, process.rb)
- [x] InSpec result: 3/3 controls pass
- [x] Commit: `a2a4fe7`

---

## 💡 Pelajaran Baru

- **Compliance as Code = testable security assertions.** Daripada checklist manual "apakah SSH port tertutup?", InSpec bikin assertions yang bisa di-run berulang dan di-integrate ke CI. Kalau control fail, langsung ketahuan.

- **InSpec matchers adalah RSpec-based.** `port()`, `processes()`, `file()`, `command()` — semua matcher InSpec pakai RSpec syntax (`it { should be_listening }`). Familiar untuk yang pernah pakai Ruby testing.

- **Impact level (0.0-1.0) menentukan severity.** 1.0 = critical (SSH open = critical), 0.7 = high (API harus listening = important tapi bukan security fatal). InSpec report mengikuti impact untuk prioritas.

- **`gem install inspec` vs `brew install --cask`.** InSpec 7.x gem hanya ship library untuk Ruby, bukan CLI executable. Untuk CLI, butuh Chef Workstation installer (`.pkg` yang butuh sudo). Homebrew cask handle ini.

- **Chef license acceptance.** InSpec sekarang butuh `--chef-license accept` di first run. Tanpa ini, InSpec refuse to execute. License acceptance disimpan di `~/.chef/accepted-licenses/`.

---

## 🔗 Referensi

- [Chef InSpec documentation](https://docs.chef.io/inspec/)
- [InSpec matchers reference](https://docs.chef.io/inspec/matchers/)
- [InSpec profile structure](https://docs.chef.io/inspec/profiles/)

---

## 📈 Mood & Energi

| Aspek | Skor (1–5) | Catatan |
|-------|-----------|---------|
| Semangat belajar | 4 | Compliance as Code pertama |
| Pemahaman materi | 4 | InSpec concept, matchers, impact level |
| Progres sesuai target | 5 | 3/3 pass, satu commit |

---

## ➡️ Rencana Besok

- [ ] Hari 29: Pipeline Consolidation — gabungkan semua scan ke satu pipeline

---

*[← Hari 27](hari-27.md) | [Hari 29 →](hari-29.md)*