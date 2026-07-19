# 1. CENTRAL AUTHENTICATION & SECURITY ===

Table users {
  id bigint [pk, increment]
  name varchar [note: "Nama Admin / Nama Owner H2H / Nama Pembeli Retail"]
  email varchar [unique]
  phone_number varchar [unique]
  password_hash varchar
  status varchar [note: "active, suspended, unverified"]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table roles {
  id bigint [pk, increment]
  role_code varchar [unique, note: "super_admin, finance, cs, merchant_h2h, member_reseller, retail_guest"]
  role_name varchar
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table model_has_roles {
  id bigint [pk, increment]
  user_id bigint
  role_id bigint
  created_at varchar
  created_by varchar
}

Table otp_codes {
  id bigint [pk, increment]
  user_id bigint
  identifier varchar [note: "Email atau nomor HP target OTP"]
  otp_code varchar
  otp_type varchar [note: "REGISTER, LOGIN, RESET_PIN"]
  expired_at varchar
  is_used boolean [default: false]
  attempt_count int [default: 0]
  max_attempt int [default: 3]
  created_at varchar
  created_by varchar
}

# 2. MERCHANT PROFILE & CREDENTIALS (UNTUK USER NON-STAFF KANTOR) ===

Table merchants {
  id bigint [pk, increment]
  user_id bigint [note: "Relasi ke tabel users (Owner dari akun merchant ini)"]
  merchant_name varchar [note: "Nama Brand / Nama Web Top-Up Mitra"]
  merchant_type varchar [note: "guest_retail, member_premium, h2h_api"]
  status varchar
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table merchant_api_credentials {
  id bigint [pk, increment]
  merchant_id bigint
  client_key varchar [unique]
  secret_key_hash varchar
  whitelist_ips text [note: "Diisi oleh user H2H di dashboard mereka untuk keamanan API"]
  is_active boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

# 3. BALANCE & MUTATION LOGS (DOUBLE-ENTRY READY) ===

Table saving_accounts {
  id bigint [pk, increment]
  merchant_id bigint [note: "Setiap member/merchant memiliki satu dompet deposit"]
  account_number varchar [unique]
  balance float64(15,2) [default: 0.00, note: "Cached value saldo saat ini"]
  account_pin_hash varchar
  status varchar
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table saving_transactions {
  id bigint [pk, increment]
  saving_account_id bigint
  type_dc varchar(1) [note: "D = Debit (Potong Saldo), C = Kredit (Tambah Saldo)"]
  amount float64(15,2)
  last_balance float64(15,2) [note: "Snapshot saldo sebelum mutasi"]
  reference_number varchar [unique, note: "ID Transaksi core terkait"]
  transaction_code varchar [note: "DEPOSIT, GAME_TOPUP, REFUND, ADJUSTMENT"]
  description text
  created_at varchar
  created_by varchar
  created_by_user bigint [note: "Bisa ID User Admin yang adjust, atau ID User H2H yang belanja"]
}

# 4. PRODUCT & PROVIDER MASTER ===

Table providers {
  id bigint [pk, increment]
  provider_name varchar
  is_active boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table product_types {
  id bigint [pk, increment]
  product_type_name varchar [note: "Prepaid, Postpaid"]
}

Table product_categories {
  id bigint [pk, increment]
  product_category_name varchar [note: "Game Top Up, Pulsa & Data, E-Wallet"]
}

Table product_references {
  id bigint [pk, increment]
  product_reference_code varchar [unique, note: "Contoh: ptsel, pind, gpubg, ewdn"]
  product_reference_name varchar [note: "Contoh: Pulsa Telkomsel, PUBG Mobile, E-Wallet DANA"]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table product_prefixes {
  id bigint [pk, increment]
  product_reference_id bigint [note: "Menghubungkan prefix ke kelompok produknya"]
  prefix_number varchar [unique, note: "Contoh: 0812, 0821, 0855"]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table products {
  id bigint [pk, increment]
  product_reference_id bigint [null]
  product_type_id bigint
  product_category_id bigint
  product_code varchar [unique, note: "SKU Internal, contoh: ML_86"]
  product_name varchar
  is_active boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table product_providers {
  id bigint [pk, increment]
  provider_id bigint
  product_id bigint
  provider_product_code varchar [note: "SKU dari sisi supplier/provider"]
  provider_price float64(15,2) [note: "Harga modal asli dari provider"]
  provider_admin_fee float64(15,2) [default: 0]
  provider_index int [default: 1, note: "Prioritas provider (1 = Utama, 2 = Backup)"]
  is_available boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table product_segments {
  id bigint [pk, increment]
  segment_name varchar [note: "Public_Retail, Gold_Reseller, H2H_Partner"]
  product_id bigint
  product_price float64(15,2) [note: "Harga jual dasar ke segment ini"]
  admin_fee float64(15,2) [default: 0]
  merchant_fee float64(15,2) [default: 0]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

# 5. PAYMENT METHOD GATEWAY ===

Table payment_methods {
  id bigint [pk, increment]
  method_code varchar [unique, note: "VIRTUAL_ACCOUNT, E_WALLET, QRIS, OVER_THE_COUNTER, DEPOSIT"]
  method_name varchar [note: "Contoh: Virtual Account, E-Wallet, QRIS, Gerai Ritel, Saldo Akun"]
  is_active boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table payment_channels {
  id bigint [pk, increment]
  payment_method_id bigint
  channel_code varchar [unique, note: "BCA_VA, MANDIRI_VA, OVO, DANA, QRIS_INTERNASIONAL, ALFAMART, BALANCE_INTERNAL"]
  channel_name varchar [note: "Contoh: BCA Virtual Account, OVO, QRIS Dana, Alfamart, Saldo Deposit"]
  fee_type varchar [note: "FIXED, PERCENTAGE"]
  fee_value float64(15,2) [default: 0.00, note: "Nilai fee PG (e.g. 4000 atau 1.5)"]
  is_active boolean [default: true]
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

# 6. CORE TRANSACTION ===

Table transactions {
  id bigint [pk, increment]
  merchant_id bigint
  product_id bigint
  product_segment_id bigint
  product_provider_id bigint
  payment_channel_id bigint [note: "Menyimpan metode pembayaran apa yang digunakan"]
  
  // Snapshots
  snapshot_product_code varchar
  snapshot_product_name varchar
  
  // Financial Details
  buy_price float64(15,2) [note: "Harga modal awal dari supplier"]
  sell_price float64(15,2) [note: "Harga produk dasar sebelum fee pembayaran"]
  admin_fee float64(15,2) [note: "Admin fee bawaan produk (jika ada)"]
  payment_fee float64(15,2) [note: "Fee charge tambahan dari payment gateway yang digunakan"]
  total_amount float64(15,2) [note: "Formula: sell_price + admin_fee + payment_fee"]
  
  // Target & Reff IDs
  target_user_id varchar [note: "ID Game / No HP Tujuan pembeli"]
  reference_number_internal varchar [unique, note: "Trx ID buatan sistem kita"]
  reference_number_merchant varchar [null, note: "Trx ID kiriman API merchant/web topup"]
  reference_number_provider varchar [note: "Trx ID dari pihak supplier"]
  serial_number varchar [note: "Token/SN hasil sukses dari provider game"]
  
  // Statuses
  status_code varchar(15) [note: "Menggunakan format komposit e.g. ERR-PVD1-302"]
  status_message varchar
  retry_count int [default: 0]
  
  created_at varchar
  created_by varchar
  updated_at varchar
  updated_by varchar
}

Table transaction_payload_logs {
  id bigint [pk, increment]
  transaction_id bigint
  request_payload text
  response_payload text
  created_at varchar
  created_by varchar
}

# RELATIONSHIPS (FOREIGN KEYS) ===

Ref: model_has_roles.user_id > users.id
Ref: model_has_roles.role_id > roles.id
Ref: otp_codes.user_id > users.id
Ref: merchants.user_id > users.id
Ref: merchant_api_credentials.merchant_id > merchants.id
Ref: saving_accounts.merchant_id - merchants.id
Ref: saving_transactions.saving_account_id > saving_accounts.id
Ref: saving_transactions.created_by_user > users.id
Ref: products.product_type_id > product_types.id
Ref: products.product_reference_id > product_references.id
Ref: product_prefixes.product_reference_id > product_references.id
Ref: products.product_category_id > product_categories.id
Ref: product_providers.provider_id > providers.id
Ref: product_providers.product_id > products.id
Ref: product_segments.product_id > products.id
Ref: payment_channels.payment_method_id > payment_methods.id
Ref: transactions.payment_channel_id > payment_channels.id
Ref: transactions.merchant_id > merchants.id
Ref: transactions.product_id > products.id
Ref: transactions.product_segment_id > product_segments.id
Ref: transactions.product_provider_id > product_providers.id
Ref: transaction_payload_logs.transaction_id - transactions.id