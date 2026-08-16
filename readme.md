# Alur Registrasi — `POST /register`

## Request Body

| Field          | Type   | Wajib | Keterangan                              |
|----------------|--------|-------|-----------------------------------------|
| name           | string | ✅    | Nama lengkap pengguna                   |
| email          | string | ✅    | Email unik                              |
| phone_number   | string | ✅    | Nomor telepon unik                      |
| password       | string | ✅    | Minimal 8 karakter                      |
| referral_code  | string | ❌    | Kode referral agent (opsional)          |

---

## Alur Proses

### 1. Validasi Input
- Password wajib minimal 8 karakter
- Cek email — jika sudah terdaftar → error `422`
- Cek nomor telepon — jika sudah terdaftar → error `423`

### 2. Tentukan Agent & Segment (Binding Otomatis)

#### Jika `referral_code` diisi:
- Cari agent aktif berdasarkan kode referral tersebut
- Jika agent ditemukan dan berstatus `active`:
  - Gunakan `agent_id` agent tersebut
  - Gunakan `segment_id` dari segment pertama yang dimiliki agent tersebut
  - Jika agent tidak memiliki segment → fallback ke default segment (ID 2)
- Jika agent tidak ditemukan / tidak aktif → fallback ke default (agent_id=1, segment_id=2), catat warning di log

#### Jika `referral_code` kosong:
- Verifikasi keberadaan default segment (ID 2) di database
- Jika ditemukan: gunakan `agent_id` dan `segment_id` dari segment tersebut
- Jika tidak ditemukan: gunakan fallback terakhir (agent_id=1, segment_id=2)

### 3. Eksekusi Transaksi (Atomik)

Seluruh langkah berikut dieksekusi dalam satu database transaction:

1. **Buat User** — nama, email, phone, password (di-hash), status `active`
2. **Bind Role Merchant** — insert ke `model_has_roles` dengan `role_id=4`, `actor_id=0` (sementara)
3. **Buat Merchant** — dengan `agent_id` dan `segment_id` hasil penentuan di langkah 2
4. **Update `actor_id`** di `model_has_roles` → diisi `merchant_id` yang baru dibuat
5. **Buat Saving Account** — nomor akun `SA-XXXXXXXX`, saldo awal `1.000.000`, PIN default `123456`

> Jika salah satu langkah gagal, seluruh transaksi di-rollback.

---

## Response Sukses

```json
{
  "code": "201",
  "data": {
    "user_id": 10,
    "agent_id": 1
  }
}
```

---

## Catatan Teknis

- Default agent: `agent_id = 1`
- Default segment: `segment_id = 2` (Public Retail)
- Role merchant: `role_id = 4`
- PIN default saving account: `123456` (harus diganti setelah login pertama)
- Saldo awal `1.000.000` hanya untuk keperluan testing