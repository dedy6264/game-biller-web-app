# Webpage
## Kelompok Kode: Transaksi Sukses / Pending (SUC / PEN)
['SUC-INT-000','SUCCESS','Transaction successfully processed and delivery confirmed.','Top-up berhasil! Produk telah masuk ke akun game Anda. Terima kasih.')]
['PEN-SYS-001','PENDING_UPSTREAM','Transaction accepted by core system and currently queued in the upstream gateway.','Pembayaran diterima. Pesanan Anda sedang diproses oleh sistem game. Mohon tunggu.')]
['INQ-INT-002','PENDING_MANUAL_REVIEW','Transaction suspended internally for security verification or compliance approval.','Transaksi Anda sedang diverifikasi demi keamanan. Mohon tunggu maksimal 5 menit.')]
['INQ-INT-003','PENDING_RETRY','Upstream network glitch detected. Core system is auto-switching or retrying the request.','Jaringan server game sedang padat. Sistem kami sedang mencoba mengirimkan ulang pesanan Anda.')]

## Kelompok Kode: Kesalahan Input Pengguna (ERR - User Input)
['ERR-VAL-100','INVALID_PLAYER_ID','External game validator service rejected the provided Player ID, Zone ID, or Server ID.','ID Game tidak terdaftar atau salah server. Silakan periksa kembali data Anda dan coba lagi.')]
['ERR-INT-101','PRODUCT_CODE_NOT_FOUND','The requested target product code or SKU is missing or inactive in our catalog.','Produk yang Anda pilih saat ini tidak tersedia. Silakan pilih nominal atau paket lainnya.')]
['ERR-INT-102','SUSPECTED_DUPLICATE','Identical transaction request detected for the same target and SKU within the restriction time.','Pembelian ganda terdeteksi! Harap tunggu 5 menit jika ingin membeli item yang sama ke akun yang sama.')]
['ERR-INT-103','INSUFFICIENT_MERCHANT_BALANCE','Internal ledger check indicates the merchant deposit balance is insufficient for this checkout.','Saldo Anda tidak mencukupi untuk melakukan pembelian ini. Silakan top-up saldo terlebih dahulu.')]
['ERR-VAL-104','INVALID_TARGET_FORMAT','The destination identity string failed regex or basic character formatting check.','Nomor tujuan atau format ID yang Anda masukkan salah. Mohon periksa kembali aturan penulisan ID.')]
['ERR-INT-105','AMOUNT_MISMATCH','The transaction amount payload sent does not match the actual designated product price segment.','Terjadi ketidaksesuaian harga produk. Silakan segarkan halaman aplikasi Anda dan coba lagi.')]

## Kelompok Kode: Masalah Operasional & Finansial Mitra (ERR - Merchant/B2B Side)
['ERR-INT-200','API_AUTHENTICATION_FAILED','Invalid client key, invalid signature hash, or IP address not registered in the whitelist.','Sesi masuk atau koneksi tidak sah. Silakan hubungi bagian administrasi akun/dukungan teknis.')]
['ERR-INT-201','MERCHANT_SUSPENDED','Merchant corporate profile has been deactivated or restricted due to compliance or risk control.','Akun Anda dinonaktifkan sementara oleh sistem. Silakan hubungi Customer Service untuk verifikasi data.')]
['ERR-INT-202','DAILY_LIMIT_EXCEEDED','The total value or volume of transactions has exceeded the agreed maximum daily limit.','Batas kuota transaksi harian Anda telah tercapai. Anda dapat melakukan transaksi kembali esok hari.')]
['ERR-INT-203','INVALID_API_METHOD','The API endpoint requested does not support the HTTP verb or request method applied.','Permintaan sistem tidak didukung. Pastikan aplikasi Anda terintegrasi dengan benar.')]
['ERR-INT-204','MAINTENANCE_SCHEDULED','The client-facing dashboard and payment processing gateway are offline for scheduled system upgrade.','Layanan kami sedang ditingkatkan untuk kenyamanan Anda. Kami akan segera kembali dalam beberapa saat.')]

## Kelompok Kode: Gangguan Sisi Provider / Supplier (ERR - Upstream Side)
['ERR-PVD1-300','UPSTREAM_MAINTENANCE','The upstream vendor gateway for this specific product is currently undergoing maintenance.','Server pusat untuk game ini sedang mengalami pemeliharaan sistem. Silakan coba beberapa saat lagi.')]
['ERR-PVD2-301','PRODUCT_OUT_OF_STOCK','The upstream distributor core reported that the stock, voucher pool, or product quota is exhausted.','Stok untuk item game ini sedang habis di pusat. Tim kami sedang melakukan pengisian ulang stok.')]
['ERR-PVD1-302','UPSTREAM_TIMEOUT','Upstream server failed to respond within the designated execution window. Money safely refunded.','Jaringan ke server game sedang padat. Transaksi dibatalkan secara aman dan saldo Anda tetap utuh.')]
['ERR-PVD3-303','UPSTREAM_DECLINED','The execution request was rejected by the distribution core gateway due to structural rule conflicts.','Pembelian ditolak oleh sistem pusat game. Silakan pilih nominal atau paket game lainnya.')]
['ERR-PVD2-304','UPSTREAM_PRICE_CHANGED','The purchase cost from the upstream gateway has changed and exceeded the configured system margin.','Terjadi pembaruan harga dari sistem pusat game. Silakan ulangi transaksi Anda untuk memperbarui harga.')]
['ERR-PVD4-305','UPSTREAM_UNKNOWN_ERROR','Upstream provider responded with an unmapped error structure or critical raw payload exception.','Terjadi kendala pada jaringan distribusi game. Saldo tidak terpotong, silakan coba lagi nanti.')]

# Dashboard
## Kelompok Modul: Autentikasi & Hak Akses (AUTH)
['SUC-AUTH-200','AUTHENTICATION_SUCCESS','User credentials verified. Session token generated successfully.','Login berhasil! Mengalihkan Anda ke halaman dashboard...')]
['SUC-AUTH-201','REGISTRATION_SUCCESS','New user entity successfully created inside the core database.','Pendaftaran berhasil! Selamat bergabung di platform kami.')]
['ERR-AUTH-401','INVALID_CREDENTIALS','The password hash or username/email combination did not match database records.','Email atau password yang Anda masukkan salah. Silakan coba lagi.')]
['ERR-AUTH-403','ACCESS_DENIED','The authenticated user does not have the required role assigned in model_has_roles.','Akun Anda tidak memiliki izin untuk mengakses halaman ini.')]
['ERR-AUTH-419','TOKEN_EXPIRED','The bearer token or session cookie has expired or been blacklisted.','Sesi Anda telah berakhir demi keamanan. Silakan masuk (login) kembali.')]

## Kelompok Modul: Validasi & Pembaruan Data (USER / VAL)
['SUC-USER-200','PROFILE_UPDATED','Profile fields successfully modified and committed to the user entity.','Perubahan profil Anda telah berhasil disimpan.')]
['VAL-USER-422','EMAIL_ALREADY_EXISTS','Unique constraint failed. The provided email address is already taken by another account.','Email sudah terdaftar. Silakan gunakan email lain atau gunakan fitur lupa password.')]
['VAL-USER-423','PHONE_ALREADY_EXISTS','Unique constraint failed. The phone number is already registered in the system.','Nomor HP sudah digunakan. Mohon periksa kembali nomor Anda.')]
['VAL-USER-424','WEAK_PASSWORD','Validation failed. The password policy string requirements are not met.','Password terlalu lemah. Gunakan minimal 8 karakter dengan kombinasi angka dan huruf.')]
['ERR-USER-404','USER_NOT_FOUND','Query returned no records for the specified User ID or identifier.','Data pengguna tidak ditemukan di dalam sistem kami.')]

## Kelompok Modul: Manajemen API & Konfigurasi Partner H2H (MERCH)
['SUC-MERCH-200','API_CREDENTIALS_REGENERATED','New client_key and secret_key successfully rolled over for the merchant entity.','API Key baru berhasil dibuat. Pastikan Anda segera memperbarui sistem H2H Anda.')]
['SUC-MERCH-201','IP_WHITELIST_UPDATED','Whitelist IP text area parsed and updated inside merchant_api_credentials.','Daftar IP Whitelist Anda telah berhasil diperbarui.')]
['ERR-MERCH-400','INVALID_IP_FORMAT','Failed to parse the provided whitelist string. IP address format is invalid.','Format IP Address yang Anda masukkan salah. Mohon periksa kembali.')]
['ERR-MERCH-403','WRONG_TRANSACTION_PIN','The secure transaction PIN validation failed against account_pin_hash.','PIN Transaksi yang Anda masukkan salah. Sisa kesempatan 2 kali lagi.')]

## Kelompok Modul: OTP & Keamanan Utilitas (UTIL)
['SUC-UTIL-200','OTP_SENT','Verification code generated and successfully pushed to the notification provider.','Kode verifikasi (OTP) telah dikirimkan ke nomor HP / email Anda.')]
['SUC-UTIL-202','OTP_VERIFIED','OTP token matches the active otp_codes row and is successfully marked as used.','Verifikasi berhasil. Silakan melanjutkan ke proses berikutnya.')]
['ERR-UTIL-408','OTP_EXPIRED','The system matched the code but the active window passed the expired_at timestamp.','Kode OTP sudah kedaluwarsa. Silakan klik tombol "Kirim Ulang".')]
['ERR-UTIL-410','OTP_INCORRECT','The input OTP string does not match the active sequence stored.','Kode OTP yang Anda masukkan salah. Mohon periksa kembali pesan masuk Anda.')]
['ERR-UTIL-429','OTP_MAX_ATTEMPT_REACHED','The counter on attempt_count has exceeded max_attempt. The token is invalidated.','Anda telah salah memasukkan OTP sebanyak 3 kali. Silakan minta kode OTP baru.')]