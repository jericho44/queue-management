# Skenario Penggunaan Nyata (End-to-End Demo Scenario)
## Multi-Tenant Queue Management System (QMS) SaaS

Dokumen ini berisi panduan skenario pengujian dan penggunaan nyata (Real-World Usage Scenario) dari sistem **Multi-Tenant Queue Management System (QMS)** mulai dari pendaftaran tenant, konfigurasi cabang & kiosk thermal printer, pendaftaran antrean pelanggan, pemanggilan loket real-time, hingga tagihan SaaS bulanan via Midtrans Payment Gateway.

---

## 👥 Pemeran / Aktor Skenario

| Aktor | Role / Entitas | Deskripsi & Tugas |
| :--- | :--- | :--- |
| **Superadmin SaaS** | `SUPER_ADMIN` | Pengelola utama platform SaaS QMS, memantau revenue & seluruh tenant |
| **Dr. Hendra** | `OWNER` | Pemilik jaringan klinik **"Klinik Sehat Sejahtera"** |
| **Rina** | `MANAGER` | Manajer Operasional cabang **"Sudirman Main Branch"** |
| **Siti** | `RECEPTIONIST` | Staf pendaftaran di kiosk / terminal pendaftaran depan |
| **Budi & Ani** | `CUSTOMER` | Pasien / Pelanggan yang mengambil nomor antrean |
| **Bambang** | `STAFF` | Staf Petugas di **Loket 01 (Poli Umum)** |

---

## 🎬 Alur Skenario Operasional (Step-by-Step)

### 📍 Langkah 1: Registrasi Organisasi Baru (Tenant Onboarding)
1. **Dr. Hendra** mengakses halaman pendaftaran tenant:
   - Jika running via **Vite Dev Server**: `http://localhost:5173/register-org`
   - Jika running via **Docker Compose**: `http://localhost/register-org`
2. Dr. Hendra mengisi formulir:
   - Nama Organisasi: `Klinik Sehat Sejahtera`
   - Kode Organisasi: `KSS`
   - Email Owner: `owner@kss.com`
   - Password: `password123`
3. **Hasil**: Sistem membuat tenant `organizations`, akun `OWNER`, dan langganan paket `POSTPAID_STANDARD` secara otomatis.


---

### 📍 Langkah 2: Konfigurasi Cabang, Layanan, & Loket
1. Dr. Hendra atau Rina (Manager) masuk ke Dashboard Admin `http://localhost/login` dan membuka menu **Branches & Counters** (`/branches`).
2. **Membuat Cabang**:
   - Nama Cabang: `Sudirman Main Branch`
   - Kode Cabang: `SDR`
3. **Membuat Jenis Layanan**:
   - Layanan 1: `Poli Umum` (Prefix `A`, Est. Servis ~8 menit)
   - Layanan 2: `Poli Gigi` (Prefix `B`, Est. Servis ~15 menit)
   - Layanan 3: `Apotek / Kasir` (Prefix `C`, Est. Servis ~3 menit)
4. **Membuat Loket Pelayanan**:
   - `Loket 01` (Ditugaskan melayani Poli Umum & Apotek)
   - `Loket 02` (Ditugaskan melayani Poli Gigi)

---

### 📍 Langkah 3: Konfigurasi Kiosk Terminal & Thermal Printer
1. Pada menu **Branches & Counters**, Rina menekan tombol **"Kiosk & Printer Settings"** untuk cabang Sudirman.
2. Rina mengatur konfigurasi:
   - **Aktifkan Kiosk**: `TRUE`
   - **Mode Cetak**: `DUAL` (Memberikan opsi Cetak Struk Fisik ATAU QR Paperless)
   - **Ukuran Kertas**: `58mm` (Thermal Paper Roll)
   - **Header Struk**: `Selamat Datang di Klinik Sehat Sejahtera - Cabang Sudirman`
   - **Footer Struk**: `Semoga lekas sembuh! Simpan struk ini hingga dipanggil.`
3. **Hasil**: URL Terminal Kiosk Publik `http://localhost/kiosk/sdr` (atau `http://localhost:5173/kiosk/sdr`) siap digunakan pada iPad / Touchscreen Kiosk di pintu masuk.

---

### 📍 Langkah 4: Pengambilan Tiket Antrean oleh Pelanggan (Touchscreen Kiosk)
1. Pasien **Budi** tiba di klinik dan mendatangi layar Touchscreen Kiosk (`/kiosk/sdr` atau `/kiosk/sudirman-main-branch`).

2. Budi menekan tombol besar **"Poli Umum (Prefix A)"**.
3. Sistem secara *concurrency-safe* (PostgreSQL `SELECT ... FOR UPDATE`) memproses antrean dan menampilkan modal pilihan **Dual Mode**:
   - Option A: **📄 Cetak Struk Fisik** (Mengeluarkan kertas thermal 58mm via WebUSB / `window.print()`).
   - Option B: **📱 Tiket Digital Paperless** (Menampilkan QR Code raksasa di layar Kiosk).
4. Budi memilih **Tiket Digital Paperless** dan memindai QR Code di layar menggunakan Kamera Smartphone miliknya.
5. **Tiket Terbit**: Nomor Antrean **`A101`** berhasil dibuat.
6. Layar Kiosk otomatis melakukan reset dalam 12 detik untuk menyambut pelanggan berikutnya.

---

### 📍 Langkah 5: Pelacakan Antrean Mandiri di Smartphone Pelanggan
1. Setelah memindai QR Code, smartphone Budi mengarah ke halaman tracking publik (`http://localhost/ticket/a3b1c2d4-...`).
2. Di layar HP Budi tampil informasi real-time:
   - **Nomor Antrean**: `A101`
   - **Status**: `WAITING` (Menunggu)
   - **Sisa Antrean Di Depan**: `0 Orang`
   - **Estimasi Waktu Tunggu**: `~8 Menit`
3. Budi dapat duduk santai di kantin tanpa takut terlewat panggilan.

---

### 📍 Langkah 6: Tampilan Public Display Layar Utama Klinik
1. Di ruang tunggu utama, TV Smart Display membuka alamat `http://localhost/display`.
2. Layar fullscreen menampilkan:
   - Header Logo & Nama Klinik `Klinik Sehat Sejahtera`.
   - Papan Informasi Loket Aktif (Loket 01 & Loket 02).
   - Daftar Antrean Sedang Dilayani & Antrean Menunggu berikutnya.
   - Efek Audio Chime Synthesizer setiap kali ada panggilan baru.

---

### 📍 Langkah 7: Pemanggilan Antrean oleh Staf Loket (Counter Station)
1. Staf **Bambang** membuka komputer `Loket 01` dan masuk ke halaman **Counter Station** (`/counter`).
2. Bambang menekan tombol **`OPEN COUNTER`** untuk membuka Loket 01.
3. Bambang menekan tombol **`CALL NEXT TICKET`**:
   - Database PostgreSQL mengeksekusi `FOR UPDATE SKIP LOCKED` pada tiket dengan status `WAITING`.
   - Tiket **`A101`** terpanggil secara instan.
4. **Efek Real-Time Sync (WebSocket + Redis)**:
   - **Public Display TV**: Bunyi *chime bell* berbunyi dan nomor `A101 ➔ Loket 01` berkedip animasi emas.
   - **Smartphone Budi**: Status berubah menjadi **`CALLED` (Harap Menuju Loket 01)** dengan notifikasi bergetar.
5. Budi berjalan menuju Loket 01. Bambang menekan **`START SERVING`** ➔ status menjadi `SERVING`.
6. Setelah pemeriksaan selesai (5 menit kemudian), Bambang menekan **`COMPLETE TICKET`** ➔ status menjadi `COMPLETED`.

---

### 📍 Langkah 8: Pengujian Concurrency / Anti-Race Condition
1. **Skenario Uji Beban**: 10 Petugas Loket menekan tombol `CALL NEXT` secara bersamaan pada milidetik yang sama di saat ada 50 antrean menunggu.
2. **Hasil Eksekusi System**:
   - Berkat penguncian `FOR UPDATE SKIP LOCKED`, **0 tiket ganda / 0 tiket yang terpanggil oleh 2 loket berbeda**.
   - Setiap loket mendapatkan 1 tiket yang unik dan presisi.

---

### 📍 Langkah 9: Rekapitulasi Tagihan Bulanan SaaS & Pembayaran via Midtrans
1. Di akhir bulan (tanggal 1), sistem secara otomatis menghitung total tiket yang diterbitkan oleh `Klinik Sehat Sejahtera` (contoh: 2.400 tiket).
2. Sistem merekapitulasi tagihan pascabayar (*Postpaid Metered Billing*):
   - 2.400 Tiket x Rp 500 = **Rp 1.200.000**.
3. **Invoice Terbit**: Invoice `INV-2026-08-0012` dengan status `UNPAID` diterbitkan.
4. Dr. Hendra (Owner) membuka menu **Billing & Invoices** (`/billing`):
   - Di layar terlihat grafik akumulasi tiket dan tagihan Rp 1.200.000.
   - Dr. Hendra menekan tombol **"Bayar Sekarang"**.
5. **Payment Gateway Midtrans Snap**:
   - Modal Pop-up Midtrans muncul dengan opsi pembayaran **QRIS (Gopay/ShopeePay)** atau **Bank Virtual Account (BCA/Mandiri/BRI)**.
6. Dr. Hendra menyelesaikan pembayaran via QRIS.
7. **Webhook Midtrans (`POST /api/v1/billing/webhooks`)**:
   - Midtrans mengirimkan callback settlement ➔ backend memverifikasi signature SHA512.
   - Status invoice otomatis berubah menjadi **`PAID` (Lunas)** secara instan.

---

## 🎯 Summary Hasil Skenario

- ✅ **Multi-Tenant Data Isolation**: Data antrean Klinik Sehat Sejahtera terisolasi penuh dari organisasi lain.
- ✅ **Concurrency-Safe Guaranteed**: Bebas dari *race condition* nomor antrean ganda atau pemanggilan bentrok.
- ✅ **Realtime Sync Multi-Device**: Sinkronisasi dalam orde milidetik antara Kiosk, Smartphone Pelanggan, Public Display TV, dan Counter Station.
- ✅ **Dual-Mode Kiosk Thermal Printer**: Mendukung struk fisik 58mm/80mm maupun Paperless QR Code.
- ✅ **Postpaid Metered Billing**: Pembayaran berbasis penggunaan tiket otomatis terintegrasi dengan Payment Gateway Midtrans.
