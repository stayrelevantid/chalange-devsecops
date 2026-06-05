# 📖 Glosarium Istilah Teknis & Asing (DevSecOps)

Dokumen ini dibuat untuk membantu pemula memahami berbagai istilah teknis, singkatan, dan bahasa asing yang sering muncul dalam pembelajaran DevSecOps dan silabus ini.

---

## 1. Konsep Dasar Pengembangan Perangkat Lunak
* **Repository (Repo):** Tempat atau "folder online" untuk menyimpan seluruh baris kode aplikasi beserta riwayat perubahannya (biasanya menggunakan Git/GitHub).
* **Source Code:** Baris kode mentah yang ditulis oleh programmer sebelum menjadi aplikasi.
* **Endpoint:** URL atau alamat spesifik dari sebuah API. (Contoh: `http://websaya.com/api/balance` adalah sebuah endpoint untuk mengecek saldo).
* **Unit Test:** Pengujian otomatis pada bagian terkecil kode (misal: satu fungsi) untuk memastikan bagian tersebut berjalan dengan benar.
* **Pull Request (PR) / Merge Request (MR):** Permintaan dari seorang developer untuk menggabungkan kode baru yang baru ia tulis ke dalam kode utama (main branch).

## 2. CI/CD & Pipeline
* **CI/CD (Continuous Integration / Continuous Deployment):** Praktik mengotomatiskan proses penggabungan kode, pengujian, hingga rilis aplikasi agar lebih cepat dan aman.
* **Pipeline:** Jalur perakitan otomatis. Serangkaian perintah yang berjalan berurutan setiap kali ada kode baru. (Misal: Pipeline melakukan build -> test -> scan keamanan -> deploy).
* **Job / Matrix Job:** Satu tugas spesifik di dalam pipeline. Matrix job berarti menjalankan banyak tugas sekaligus secara paralel (bersamaan).
* **Environment Variables (Env Vars):** Variabel atau nilai konfigurasi yang disimpan di sistem operasi komputer (atau server), bukan ditulis langsung di dalam kode.
* **Quality Gate (Gerbang Kualitas):** Titik pengecekan dalam pipeline. Jika kode tidak lulus syarat keamanan tertentu (misal: ditemukan bug bahaya), pipeline akan dihentikan paksa (Build Failed/Merah).

## 3. Istilah Keamanan (Security)
* **Hardcode (Hardcoded):** Kebiasaan buruk menuliskan data sensitif (seperti password atau token) secara langsung berupa teks biasa di dalam source code.
* **Secret / Credential:** Data rahasia yang digunakan untuk otentikasi, seperti password, token API, kunci enkripsi, dan kunci akses AWS.
* **Vulnerability:** Kerentanan atau "celah keamanan" pada aplikasi atau sistem yang bisa dimanfaatkan oleh peretas (hacker).
* **CVE (Common Vulnerabilities and Exposures):** Nomor seri resmi untuk sebuah celah keamanan publik. (Contoh: CVE-2021-44228 adalah nama resmi untuk celah Log4j).
* **Remediation (Remediasi):** Tindakan perbaikan, penambalan (patch), atau solusi untuk menutup sebuah kerentanan.
* **Red Teaming:** Simulasi serangan siber yang dilakukan oleh tim ahli (bertindak sebagai hacker) untuk menguji seberapa kuat pertahanan sebuah sistem.
* **Threat Modeling:** Proses menganalisis sistem sejak fase desain untuk mencari tahu dari mana saja sistem ini bisa diserang dan merancang cara mencegahnya (contoh: metode STRIDE).
* **Least Privilege:** Prinsip keamanan di mana seseorang atau sistem hanya diberikan hak akses seminimal mungkin, sekadar cukup untuk melakukan pekerjaannya.

## 4. Jenis-Jenis Pemindaian Keamanan (Scanning)
* **Secret Scanning:** Proses memindai source code untuk mencari password atau kunci rahasia yang tidak sengaja tertulis (bocor). (Tool: Gitleaks).
* **SCA (Software Composition Analysis):** Pemindaian untuk mencari kerentanan pada "perpustakaan pihak ketiga" (library/dependency) yang digunakan oleh aplikasi. (Tool: Trivy).
* **SAST (Static Application Security Testing):** Pemindaian pada source code (sebelum aplikasi dijalankan) untuk mencari gaya penulisan kode yang salah atau berbahaya, seperti potensi SQL Injection. (Tool: Semgrep).
* **DAST (Dynamic Application Security Testing):** Pemindaian dengan cara "menyerang" aplikasi yang sedang menyala/berjalan (runtime) dari luar untuk menemukan celah yang bisa dieksploitasi. (Tool: OWASP ZAP).
* **CSPM (Cloud Security Posture Management):** Alat untuk memeriksa konfigurasi keamanan akun Cloud (seperti AWS atau GCP) apakah sudah sesuai standar atau belum (misal: mengecek apakah S3 Bucket terbuka untuk publik). (Tool: Prowler).

## 5. Infrastruktur, Container & Cloud
* **Container / Docker:** Teknologi untuk membungkus aplikasi beserta seluruh kebutuhannya ke dalam satu kotak (container) agar bisa dijalankan di mana saja dengan hasil yang sama.
* **Image (Docker Image):** Cetakan atau blueprint dari sebuah container. Container adalah Image yang sedang dijalankan (running).
* **Multi-stage Build:** Teknik membuat Docker Image di mana proses kompilasi (build) dipisah dari proses menjalankan aplikasi, agar ukuran image akhir lebih kecil dan aman.
* **IaC (Infrastructure as Code):** Praktik mengelola infrastruktur server dan jaringan menggunakan baris kode (bukan klik-klik manual di web). (Tool: Terraform).
* **Kubernetes (K8s):** Sistem manajemen otomatis (orkestrasi) untuk mengatur ribuan container Docker.
* **Pod:** Unit terkecil di Kubernetes. Satu pod biasanya berisi satu atau beberapa container yang berjalan bersama.
* **Cluster:** Kumpulan beberapa server/komputer (node) yang disatukan dan dikelola oleh Kubernetes.

## 6. Runtime Security & Policies
* **OPA (Open Policy Agent) Gatekeeper:** Penjaga pintu di Kubernetes. Ia bertugas memastikan tidak ada hal yang di-deploy ke cluster jika melanggar aturan keamanan (misal: menolak aplikasi yang berjalan sebagai root).
* **RBAC (Role-Based Access Control):** Sistem pengaturan hak akses berdasarkan peran pengguna (siapa boleh melakukan apa).
* **Network Policy:** Layaknya firewall internal di dalam Kubernetes yang mengatur aplikasi A boleh berkomunikasi dengan aplikasi B atau tidak.
* **Lateral Movement:** Pergerakan menyamping. Taktik peretas yang setelah berhasil menembus satu server lemah, menggunakannya sebagai batu loncatan untuk menyerang server lain di dalam jaringan yang sama.
* **Escape / Privilege Escalation:** Taktik peretas untuk menaikkan hak aksesnya (dari pengguna biasa menjadi admin/root) atau kabur dari dalam container ke sistem utama (host).

## 7. Singkatan Tambahan
* **JWT (JSON Web Token):** Standar untuk membuat token digital sebagai bukti bahwa seseorang sudah login dan berhak mengakses aplikasi.
* **WAF (Web Application Firewall):** Tembok pelindung khusus untuk aplikasi web yang bisa memblokir lalu lintas jahat, seperti upaya SQL Injection atau DDoS.
* **CIS Benchmarks (Center for Internet Security):** Panduan standar global berisi daftar "cara mengonfigurasi sistem dengan aman". 
* **LLM (Large Language Model):** AI pemroses teks seperti ChatGPT, Claude, atau Gemini.
* **Webhook:** Cara sebuah aplikasi untuk mengirimkan data secara instan (real-time) ke aplikasi lain segera setelah ada kejadian tertentu (event).
