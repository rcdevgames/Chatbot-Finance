# Fitur Hapus Data

## Overview

Fitur hapus data memungkinkan pengguna untuk menghapus data keuangan mereka dengan berbagai opsi periode waktu. Fitur ini dilengkapi dengan sistem konfirmasi untuk mencegah penghapusan tidak sengaja.

## Commands

### Menu Hapus Data
```
/hapus
```
Menampilkan menu hapus data dengan semua opsi yang tersedia.

### Hapus Data Per Bulan
```
/hapus bulan 2024 1
```
Menghapus semua data untuk bulan Januari 2024.

Format: `/hapus bulan [tahun] [bulan]`
Contoh:
- `/hapus bulan 2024 1` - Hapus data Januari 2024
- `/hapus bulan 2024 12` - Hapus data Desember 2024

### Hapus Data Per Tahun
```
/hapus tahun 2024
```
Menghapus semua data untuk tahun 2024.

Format: `/hapus tahun [tahun]`
Contoh:
- `/hapus tahun 2024` - Hapus semua data 2024
- `/hapus tahun 2023` - Hapus semua data 2023

### Hapus Semua Data
```
/hapus semua
```
Menghapus SEMUA data pengguna (transaksi dan chat history).

## Proses Konfirmasi

1. Setelah menjalankan perintah hapus, sistem akan menampilkan pesan konfirmasi
2. Pesan konfirmasi berisi:
   - Data yang akan dihapus
   - Peringatan bahwa tindakan tidak dapat dibatalkan
   - Kode konfirmasi unik
   - Instruksi untuk konfirmasi

3. Untuk melanjutkan penghapusan, gunakan:
   ```
   /konfirmasi [kode]
   ```
   Contoh: `/konfirmasi 3851`

4. Jika tidak ingin melanjutkan, abaikan pesan konfirmasi

## Keamanan

- Kode konfirmasi unik untuk setiap permintaan hapus
- Kode konfirmasi berlaku selama 5 menit
- Pesan peringatan jelas tentang konsekuensi penghapusan
- Validasi input untuk mencegah format salah

## Data yang Dihapus

### Periode Bulan/Tahun
- Semua transaksi dalam periode tersebut
- Semua chat history dalam periode tersebut

### Semua Data
- Semua transaksi pengguna
- Semua chat history pengguna

## Implementasi Technical

### Repository Functions
- `DeleteTransactionsByUserIDAndMonth`
- `DeleteTransactionsByUserIDAndYear`
- `DeleteAllTransactionsByUserID`
- `DeleteChatHistoryByUserIDAndMonth`
- `DeleteChatHistoryByUserIDAndYear`

### Service Functions
- `DeleteDataByPeriod(userID, period, year, month)`

### Handler Functions
- `handleDeleteDataMenu`
- `handleDeleteCommand`
- `requestDeleteConfirmation`
- `handleDeleteConfirmation`
- `validateConfirmationCode`

## Error Handling

- Validasi format input
- Validasi range tahun (2000-2100)
- Validasi range bulan (1-12)
- Validasi kode konfirmasi
- Pesan error yang jelas dan informatif