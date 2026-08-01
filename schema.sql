-- DROP tables if exist
DROP TABLE IF EXISTS transaction_payload_logs CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS payment_channels CASCADE;
DROP TABLE IF EXISTS payment_methods CASCADE;
DROP TABLE IF EXISTS product_segments CASCADE;
DROP TABLE IF EXISTS segments CASCADE;
DROP TABLE IF EXISTS product_providers CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS product_prefixes CASCADE;
DROP TABLE IF EXISTS product_references CASCADE;
DROP TABLE IF EXISTS product_categories CASCADE;
DROP TABLE IF EXISTS product_types CASCADE;
DROP TABLE IF EXISTS providers CASCADE;
DROP TABLE IF EXISTS saving_transactions CASCADE;
DROP TABLE IF EXISTS saving_accounts CASCADE;
DROP TABLE IF EXISTS merchant_api_credentials CASCADE;
DROP TABLE IF EXISTS merchants CASCADE;
DROP TABLE IF EXISTS otp_codes CASCADE;
DROP TABLE IF EXISTS model_has_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 1. CENTRAL AUTHENTICATION & SECURITY
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  phone_number VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  status VARCHAR(50), -- active, suspended, unverified
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  role_code VARCHAR(255) UNIQUE, -- super_admin, finance, cs, merchant_h2h, member_reseller, retail_guest
  role_name VARCHAR(255),
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE model_has_roles (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE UNIQUE,
  role_id BIGINT REFERENCES roles(id) ON DELETE CASCADE,
  created_at VARCHAR(255),
  created_by VARCHAR(255)
);

CREATE TABLE otp_codes (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
  identifier VARCHAR(255), -- Email atau nomor HP target OTP
  otp_code VARCHAR(255),
  otp_type VARCHAR(255), -- REGISTER, LOGIN, RESET_PIN
  expired_at VARCHAR(255),
  is_used BOOLEAN DEFAULT FALSE,
  attempt_count INT DEFAULT 0,
  max_attempt INT DEFAULT 3,
  created_at VARCHAR(255),
  created_by VARCHAR(255)
);

CREATE TABLE segments (
  id BIGSERIAL PRIMARY KEY,
  segment_name VARCHAR(255) UNIQUE NOT NULL,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);
-- 2. MERCHANT PROFILE & CREDENTIALS
CREATE TABLE merchants (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE UNIQUE,
  segment_id BIGINT REFERENCES segments(id) ON DELETE SET NULL,
  merchant_name VARCHAR(255),
  merchant_type VARCHAR(255), -- guest_retail, member_premium, h2h_api
  status VARCHAR(50),
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE merchant_api_credentials (
  id BIGSERIAL PRIMARY KEY,
  merchant_id BIGINT REFERENCES merchants(id) ON DELETE CASCADE,
  client_key VARCHAR(255) UNIQUE,
  secret_key_hash VARCHAR(255),
  whitelist_ips TEXT,
  is_active BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

-- 3. BALANCE & MUTATION LOGS
CREATE TABLE saving_accounts (
  id BIGSERIAL PRIMARY KEY,
  merchant_id BIGINT UNIQUE REFERENCES merchants(id) ON DELETE CASCADE,
  account_number VARCHAR(255) UNIQUE,
  balance NUMERIC(15,2) DEFAULT 0.00,
  account_pin_hash VARCHAR(255),
  status VARCHAR(50),
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE saving_transactions (
  id BIGSERIAL PRIMARY KEY,
  saving_account_id BIGINT REFERENCES saving_accounts(id) ON DELETE CASCADE,
  type_dc VARCHAR(1), -- D = Debit, C = Credit
  amount NUMERIC(15,2),
  last_balance NUMERIC(15,2),
  reference_number VARCHAR(255) UNIQUE,
  transaction_code VARCHAR(255), -- DEPOSIT, GAME_TOPUP, REFUND, ADJUSTMENT
  description TEXT,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  created_by_user BIGINT REFERENCES users(id) ON DELETE SET NULL
);

-- 4. PRODUCT & PROVIDER MASTER
CREATE TABLE providers (
  id BIGSERIAL PRIMARY KEY,
  provider_name VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE product_types (
  id BIGSERIAL PRIMARY KEY,
  product_type_name VARCHAR(255) -- Prepaid, Postpaid
);

CREATE TABLE product_categories (
  id BIGSERIAL PRIMARY KEY,
  product_category_name VARCHAR(255) -- Game Top Up, Pulsa & Data, E-Wallet
);

CREATE TABLE product_references (
  id BIGSERIAL PRIMARY KEY,
  product_reference_code VARCHAR(255) UNIQUE,
  product_reference_name VARCHAR(255),
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE product_prefixes (
  id BIGSERIAL PRIMARY KEY,
  product_reference_id BIGINT REFERENCES product_references(id) ON DELETE CASCADE,
  prefix_number VARCHAR(255) UNIQUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  product_reference_id BIGINT REFERENCES product_references(id) ON DELETE SET NULL,
  product_type_id BIGINT REFERENCES product_types(id) ON DELETE CASCADE,
  product_category_id BIGINT REFERENCES product_categories(id) ON DELETE CASCADE,
  product_code VARCHAR(255) UNIQUE,
  product_name VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE product_providers (
  id BIGSERIAL PRIMARY KEY,
  provider_id BIGINT REFERENCES providers(id) ON DELETE CASCADE,
  provider_name VARCHAR(255),
  product_provider_code VARCHAR(255),
  product_provider_name VARCHAR(255),
  product_provider_price NUMERIC(15,2),
  product_provider_admin_fee NUMERIC(15,2) DEFAULT 0.00,
  product_provider_merchant_fee NUMERIC(15,2) DEFAULT 0.00,
  provider_index INT DEFAULT 1,
  is_available BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);


CREATE TABLE product_segments (
  id BIGSERIAL PRIMARY KEY,
  segment_id BIGINT REFERENCES segments(id) ON DELETE SET NULL,
  product_provider_id BIGINT REFERENCES product_providers(id) ON DELETE SET NULL,
  segment_name VARCHAR(255), -- Public_Retail, Gold_Reseller, H2H_Partner
  product_name VARCHAR(255),
  provider_name VARCHAR(255),
  product_provider_code VARCHAR(255),
  product_provider_name VARCHAR(255),
  product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,
  product_price NUMERIC(15,2),
  admin_fee NUMERIC(15,2) DEFAULT 0.00,
  merchant_fee NUMERIC(15,2) DEFAULT 0.00,
  product_provider_price NUMERIC(15,2) DEFAULT 0.00,
  product_provider_admin_fee NUMERIC(15,2) DEFAULT 0.00,
  product_provider_merchant_fee NUMERIC(15,2) DEFAULT 0.00,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

-- 5. PAYMENT METHOD GATEWAY
CREATE TABLE payment_methods (
  id BIGSERIAL PRIMARY KEY,
  method_code VARCHAR(255) UNIQUE, -- VIRTUAL_ACCOUNT, E_WALLET, QRIS, OVER_THE_COUNTER, DEPOSIT
  method_name VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE payment_channels (
  id BIGSERIAL PRIMARY KEY,
  payment_method_id BIGINT REFERENCES payment_methods(id) ON DELETE CASCADE,
  channel_code VARCHAR(255) UNIQUE, -- BCA_VA, MANDIRI_VA, OVO, DANA, QRIS_INTERNASIONAL, ALFAMART, BALANCE_INTERNAL
  channel_name VARCHAR(255),
  fee_type VARCHAR(50), -- FIXED, PERCENTAGE
  fee_value NUMERIC(15,2) DEFAULT 0.00,
  is_active BOOLEAN DEFAULT TRUE,
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

-- 6. CORE TRANSACTION
CREATE TABLE transactions (
  id BIGSERIAL PRIMARY KEY,
  --merchant
  merchant_id BIGINT REFERENCES merchants(id) ON DELETE CASCADE,
  merchant_name VARCHAR(255),
  --product type
  product_type_id BIGINT REFERENCES product_types(id) ON DELETE SET NULL,
  --product
  product_id BIGINT REFERENCES products(id) ON DELETE SET NULL,
  product_code VARCHAR(255),
  product_name VARCHAR(255),
  product_price NUMERIC(15,2),
  product_admin_fee NUMERIC(15,2),
  product_merchant_fee NUMERIC(15,2),
  --product provider
  product_provider_id BIGINT REFERENCES product_providers(id) ON DELETE SET NULL,
  product_provider_code VARCHAR(255),
  product_provider_name VARCHAR(255),
  product_provider_price NUMERIC(15,2),
  product_provider_admin_fee NUMERIC(15,2),
  product_provider_merchant_fee NUMERIC(15,2),
  --provider
  provider_id BIGINT REFERENCES providers(id) ON DELETE SET NULL,
  provider_name VARCHAR(255),
  --segment
  --product segment
  product_segment_id BIGINT REFERENCES product_segments(id) ON DELETE SET NULL,
  product_segment_name VARCHAR(255),
  product_reference_id BIGINT REFERENCES product_references(id) ON DELETE SET NULL,
  payment_channel_id BIGINT REFERENCES payment_channels(id) ON DELETE SET NULL,
  
  --snapshot_product_code VARCHAR(255),
  --snapshot_product_name VARCHAR(255),
  product_type_name VARCHAR(255),
  payment_channel_name VARCHAR(255),
  
  
  
  payment_admin_fee NUMERIC(15,2),
  total_amount NUMERIC(15,2),
  
  customer_id VARCHAR(255),
  other_customer_id TEXT,
  customer_phone VARCHAR(255),
  reference_number_internal VARCHAR(255) UNIQUE,
  reference_number_merchant VARCHAR(255),
  reference_number_provider VARCHAR(255),
  serial_number VARCHAR(255),
  
  status_code VARCHAR(50),
  status_message VARCHAR(255),
  retry_count INT DEFAULT 0,
  
  created_at VARCHAR(255),
  created_by VARCHAR(255),
  updated_at VARCHAR(255),
  updated_by VARCHAR(255)
);

CREATE TABLE transaction_payload_logs (
  id BIGSERIAL PRIMARY KEY,
  transaction_id BIGINT REFERENCES transactions(id) ON DELETE CASCADE,
  request_payload TEXT,
  response_payload TEXT,
  created_at VARCHAR(255),
  created_by VARCHAR(255)
);

-- ============================================================
-- INDEXES — fields commonly used for searching / filtering
-- ============================================================

-- users
CREATE INDEX IF NOT EXISTS idx_users_email            ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone_number     ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_status           ON users(status);

-- roles
CREATE INDEX IF NOT EXISTS idx_roles_role_code        ON roles(role_code);

-- model_has_roles
CREATE INDEX IF NOT EXISTS idx_model_has_roles_user_id ON model_has_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_model_has_roles_role_id ON model_has_roles(role_id);

-- otp_codes
CREATE INDEX IF NOT EXISTS idx_otp_codes_user_id      ON otp_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_otp_codes_identifier   ON otp_codes(identifier);
CREATE INDEX IF NOT EXISTS idx_otp_codes_otp_type     ON otp_codes(otp_type);
CREATE INDEX IF NOT EXISTS idx_otp_codes_is_used      ON otp_codes(is_used);

-- merchants
CREATE INDEX IF NOT EXISTS idx_merchants_user_id       ON merchants(user_id);
CREATE INDEX IF NOT EXISTS idx_merchants_segment_id    ON merchants(segment_id);
CREATE INDEX IF NOT EXISTS idx_merchants_merchant_type ON merchants(merchant_type);
CREATE INDEX IF NOT EXISTS idx_merchants_status        ON merchants(status);

-- merchant_api_credentials
CREATE INDEX IF NOT EXISTS idx_mac_merchant_id         ON merchant_api_credentials(merchant_id);
CREATE INDEX IF NOT EXISTS idx_mac_client_key          ON merchant_api_credentials(client_key);
CREATE INDEX IF NOT EXISTS idx_mac_is_active           ON merchant_api_credentials(is_active);

-- saving_accounts
CREATE INDEX IF NOT EXISTS idx_saving_accounts_merchant_id     ON saving_accounts(merchant_id);
CREATE INDEX IF NOT EXISTS idx_saving_accounts_account_number  ON saving_accounts(account_number);
CREATE INDEX IF NOT EXISTS idx_saving_accounts_status          ON saving_accounts(status);

-- saving_transactions
CREATE INDEX IF NOT EXISTS idx_saving_transactions_saving_account_id ON saving_transactions(saving_account_id);
CREATE INDEX IF NOT EXISTS idx_saving_transactions_reference_number  ON saving_transactions(reference_number);
CREATE INDEX IF NOT EXISTS idx_saving_transactions_transaction_code  ON saving_transactions(transaction_code);
CREATE INDEX IF NOT EXISTS idx_saving_transactions_type_dc           ON saving_transactions(type_dc);

-- providers
CREATE INDEX IF NOT EXISTS idx_providers_is_active     ON providers(is_active);

-- product_references
CREATE INDEX IF NOT EXISTS idx_product_references_code ON product_references(product_reference_code);

-- product_prefixes
CREATE INDEX IF NOT EXISTS idx_product_prefixes_product_reference_id ON product_prefixes(product_reference_id);
CREATE INDEX IF NOT EXISTS idx_product_prefixes_prefix_number        ON product_prefixes(prefix_number);

-- products
CREATE INDEX IF NOT EXISTS idx_products_product_reference_id ON products(product_reference_id);
CREATE INDEX IF NOT EXISTS idx_products_product_type_id      ON products(product_type_id);
CREATE INDEX IF NOT EXISTS idx_products_product_category_id  ON products(product_category_id);
CREATE INDEX IF NOT EXISTS idx_products_product_code         ON products(product_code);
CREATE INDEX IF NOT EXISTS idx_products_is_active            ON products(is_active);

-- product_providers
CREATE INDEX IF NOT EXISTS idx_product_providers_provider_id          ON product_providers(provider_id);
CREATE INDEX IF NOT EXISTS idx_product_providers_is_available         ON product_providers(is_available);
CREATE INDEX IF NOT EXISTS idx_product_providers_product_provider_code ON product_providers(product_provider_code);

-- segments
CREATE INDEX IF NOT EXISTS idx_segments_segment_name  ON segments(segment_name);

-- product_segments
CREATE INDEX IF NOT EXISTS idx_product_segments_segment_id         ON product_segments(segment_id);
CREATE INDEX IF NOT EXISTS idx_product_segments_product_id         ON product_segments(product_id);
CREATE INDEX IF NOT EXISTS idx_product_segments_product_provider_id ON product_segments(product_provider_id);
CREATE INDEX IF NOT EXISTS idx_product_segments_segment_name       ON product_segments(segment_name);
-- Composite: digunakan pada query GetProductSegmentByProductAndSegment
CREATE INDEX IF NOT EXISTS idx_product_segments_product_segment_lookup ON product_segments(product_id, segment_name);

-- payment_methods
CREATE INDEX IF NOT EXISTS idx_payment_methods_method_code ON payment_methods(method_code);
CREATE INDEX IF NOT EXISTS idx_payment_methods_is_active   ON payment_methods(is_active);

-- payment_channels
CREATE INDEX IF NOT EXISTS idx_payment_channels_payment_method_id ON payment_channels(payment_method_id);
CREATE INDEX IF NOT EXISTS idx_payment_channels_channel_code      ON payment_channels(channel_code);
CREATE INDEX IF NOT EXISTS idx_payment_channels_is_active         ON payment_channels(is_active);

-- transactions
CREATE INDEX IF NOT EXISTS idx_transactions_merchant_id               ON transactions(merchant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_id                ON transactions(product_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_segment_id        ON transactions(product_segment_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_provider_id       ON transactions(product_provider_id);
CREATE INDEX IF NOT EXISTS idx_transactions_provider_id               ON transactions(provider_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_type_id           ON transactions(product_type_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_reference_id      ON transactions(product_reference_id);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_channel_id        ON transactions(payment_channel_id);
CREATE INDEX IF NOT EXISTS idx_transactions_product_code              ON transactions(product_code);
CREATE INDEX IF NOT EXISTS idx_transactions_product_provider_code     ON transactions(product_provider_code);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_number_internal ON transactions(reference_number_internal);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_number_merchant ON transactions(reference_number_merchant);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_number_provider ON transactions(reference_number_provider);
CREATE INDEX IF NOT EXISTS idx_transactions_status_code               ON transactions(status_code);
CREATE INDEX IF NOT EXISTS idx_transactions_customer_id               ON transactions(customer_id);
CREATE INDEX IF NOT EXISTS idx_transactions_other_customer_id         ON transactions(other_customer_id);
-- CREATE INDEX IF NOT EXISTS idx_transactions_snapshot_product_code     ON transactions(snapshot_product_code);
-- Composite: merchant + status sering difilter bersamaan di dashboard
CREATE INDEX IF NOT EXISTS idx_transactions_merchant_status           ON transactions(merchant_id, status_code);

-- transaction_payload_logs
CREATE INDEX IF NOT EXISTS idx_transaction_payload_logs_transaction_id ON transaction_payload_logs(transaction_id);

--- Seed Initial Master Data

INSERT INTO roles (role_code, role_name, created_at, updated_at) VALUES
('super_admin', 'Super Administrator Internal', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('finance', 'Finance & Billing Internal', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('merchant_h2h', 'Mitra Host-to-Host API', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('member_reseller', 'Reseller VIP Dashboard', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('retail_guest', 'Pembeli Lepas Web Retail', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

INSERT INTO product_types (product_type_name) VALUES
('Prepaid'),
('Postpaid');

INSERT INTO product_categories (product_category_name) VALUES
('Game Top Up'),
('Pulsa & Data'),
('E-Wallet'),
('PLN');

INSERT INTO product_references (product_reference_code, product_reference_name, created_at, updated_at) VALUES
('ref_tsel', 'TELKOMSEL', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_isat', 'INDOSAT', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_three', 'THREE', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_axis', 'AXIS', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_smart', 'SMARTFREN', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_xl', 'XL', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_byu', 'BY.U', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_mlbb', 'MOBILE LEGEND', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_genshin', 'GENSHIN IMPACT', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_ff', 'FREE FIRE', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_pubg', 'PUBG ID', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_ragnarok', 'RAGNAROK', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_pb', 'POINT BLANK', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_speed_drifters', 'Speed Drifters', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_aof', 'Arena of Valor', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_valoran', 'Valorant', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_steam_wallet', 'Steam Sea', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('ref_garena', 'Garena', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(1, '08120', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08121', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08122', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08123', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08124', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08125', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08126', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08127', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08128', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08129', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08130', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08131', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08132', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08133', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08134', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08135', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08136', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08137', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08138', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08139', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08520', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08521', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08522', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08523', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08524', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08525', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08526', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08527', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08528', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08529', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08530', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08531', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08532', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08533', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08534', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08535', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08536', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08537', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08538', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08539', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08210', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08211', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08212', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08213', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08214', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08215', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08216', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08217', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08218', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08219', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08230', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08231', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08232', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08233', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08234', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08235', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08236', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08237', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08238', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08239', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08220', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08221', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08222', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08223', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08224', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(1, '08225', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08226', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08227', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08228', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (1, '08229', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- INDOSAT prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(2, '08140', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08141', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08142', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08143', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08144', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08145', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08146', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08147', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08148', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08149', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08150', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08151', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08152', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08153', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08154', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08155', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08156', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08157', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08158', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08159', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08160', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08161', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08162', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08163', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08164', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08165', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08166', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08167', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08168', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08169', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08550', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08551', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08552', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08553', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08554', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08555', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08556', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08557', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08558', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08559', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08560', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08561', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08562', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08563', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08564', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08565', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08566', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08567', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08568', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08569', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08570', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08571', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08572', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08573', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08574', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08575', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08576', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08577', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08578', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08579', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08580', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08581', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08582', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08583', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08584', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, '08585', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08586', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08587', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08588', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (2, '08589', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- THREE prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(3, '08960', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08961', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08962', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08963', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08964', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08965', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08966', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08967', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08968', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08969', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08970', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08971', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08972', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08973', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08974', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08975', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08976', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08977', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08978', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08979', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08980', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08981', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08982', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08983', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08984', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08985', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08986', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08987', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08988', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08989', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08990', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08991', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08992', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08993', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08994', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08995', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08996', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08997', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08998', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08999', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08950', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08951', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08952', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08953', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08954', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(3, '08955', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08956', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08957', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08958', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (3, '08959', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- AXIS prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(4, '08380', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08381', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08382', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08383', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08384', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08385', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08386', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08387', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08388', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08389', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08370', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08371', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08372', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08373', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08374', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08375', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08376', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08377', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08378', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08379', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08310', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08311', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08312', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08313', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08314', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08315', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08316', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08317', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08318', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08319', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08320', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08321', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08322', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08323', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08324', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(4, '08325', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08326', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08327', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08328', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (4, '08329', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- SMARTFREN prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(5, '08810', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08811', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08812', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08813', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08814', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08815', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08816', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08817', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08818', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08819', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08820', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08821', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08822', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08823', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08824', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08825', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08826', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08827', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08828', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08829', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08830', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08831', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08832', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08833', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08834', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08835', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08836', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08837', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08838', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08839', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08840', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08841', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08842', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08843', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08844', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08845', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08846', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08847', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08848', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08849', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08850', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08851', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08852', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08853', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08854', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08855', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08856', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08857', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08858', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08859', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08860', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08861', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08862', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08863', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08864', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08865', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08866', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08867', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08868', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08869', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08870', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08871', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08872', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08873', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08874', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08875', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08876', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08877', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08878', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08879', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08880', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08881', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08882', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08883', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08884', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(5, '08885', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08886', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08887', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08888', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (5, '08889', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- XL prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(6, '08170', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08171', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08172', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08173', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08174', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08175', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08176', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08177', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08178', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08179', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08180', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08181', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08182', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08183', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08184', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08185', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08186', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08187', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08188', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08189', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08190', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08191', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08192', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08193', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08194', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08195', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08196', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08197', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08198', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08199', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08590', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08591', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08592', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08593', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08594', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08595', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08596', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08597', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08598', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08599', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08780', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08781', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08782', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08783', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08784', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08785', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08786', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08787', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08788', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08789', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08770', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08771', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08772', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08773', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08774', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(6, '08775', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08776', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08777', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08778', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (6, '08779', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- BY.U prefix
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(7, '08510', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08511', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08512', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08513', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08514', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(7, '08515', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08516', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08517', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08518', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), (7, '08519', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- Master Data Payment Methods
INSERT INTO payment_methods (method_code, method_name, created_at, updated_at) VALUES
('DEPOSIT', 'Saldo Deposit internal', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('VIRTUAL_ACCOUNT', 'Virtual Account Transfer', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('QRIS', 'QR Code QRIS', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
('E_WALLET', 'Dompet Digital / E-Wallet', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'));

-- Master Data Payment Channels
INSERT INTO payment_channels (payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, updated_at) VALUES
-- Jalur potong saldo deposit (Utamanya untuk mitra H2H / Member Premium)
(1, 'BALANCE_INTERNAL', 'Saldo Deposit Akun', 'FIXED', 0.00, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
-- Jalur Ritel Web Top Up via Virtual Account (Fee flat per transaksi)
(2, 'BCA_VA', 'BCA Virtual Account', 'FIXED', 4000.00, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
(2, 'MANDIRI_VA', 'Mandiri Virtual Account', 'FIXED', 3500.00, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')),
-- Jalur Ritel Web Top Up via QRIS (Fee persentase MDR dari nilai transaksi)
(3, 'QRIS_GATEWAY', 'QRIS Dana/LinkAja (All Shopee/Gopay)', 'PERCENTAGE', 0.70, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')), -- 0.7% MDR
-- Jalur Ritel via E-Wallet Direct Link
(4, 'OVO_DIRECT', 'OVO Instant Payment', 'PERCENTAGE', 1.50, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')); -- 1.5% Fee

-- Insert Test Provider
INSERT INTO providers (provider_name, is_active, created_at, created_by, updated_at, updated_by) VALUES
('IAK Provider', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Test Product ML_86
INSERT INTO products (product_reference_id, product_type_id, product_category_id, product_code, product_name, is_active, created_at, created_by, updated_at, updated_by) VALUES
 (8, 1, 1, 'ML_86', 'Mobile Legends 86 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam600k', 'Steam Wallet Code 600K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam400k', 'Steam Wallet Code 400K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam250k', 'Steam Wallet Code 250K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam120k', 'Steam Wallet Code 120K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam90k', 'Steam Wallet Code 90K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam60k', 'Steam Wallet Code 60K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
 (17, 1, 1, 'steam45k', 'Steam Wallet Code 45K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (17, 1, 1, 'steam12k', 'Steam Wallet Code 12K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (17, 1, 1, 'steam8k', 'Steam Wallet Code 8K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (17, 1, 1, 'steam6k', 'Steam Wallet Code 6K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant11KP', 'Valorant 11.000 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant5350P', 'Valorant 5.350 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant4125P', 'Valorant 4.125 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant3650P', 'Valorant 3.650 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant2050P', 'Valorant 2.050 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valorant1KP', 'Valorant 1.000 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (16, 1, 1, 'valoran475P', 'Valorant 475 P', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (12, 1, 1, 'ragnarok18bcc', 'Ragnarok 18 BCC', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (12, 1, 1, 'ragnarog12bcc', 'Ragnarok 12 BCC', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (12, 1, 1, 'ragnarok6bcc', 'Ragnarok 6 BCC', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (10, 1, 1, 'FF_12', 'Free Fire 12 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (10, 1, 1, 'FF_5', 'Free Fire 5 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (10, 1, 1, 'FF_10', 'Free Fire 10 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (8, 1, 1, 'ML_10', 'Mobile Legends 10 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (8, 1, 1, 'ML_5', 'Mobile Legends 5 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (8, 1, 1, 'ML_3', 'Mobile Legends 3 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 2, 'PI5K', 'Pulsa Isat 5K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 2, 'PI10K', 'Pulsa Isat 10K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 2, 'PI50K', 'Pulsa Isat 50K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 2, 'PI100K', 'Pulsa Isat 100K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 1, 'PT100K', 'Pulsa Telkomsel 100K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 1, 'PT50K', 'Pulsa Telkomsel 50K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 1, 'PT10K', 'Pulsa Telkomsel 10K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys'),
 (2, 1, 1, 'PT5K', 'Pulsa Telkomsel 5K', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'sys');

INSERT INTO product_segments (segment_name, product_id, product_price, admin_fee, merchant_fee, created_at, created_by, updated_at, updated_by) VALUES
('Open_Biller', 1, 20000.00, 1000.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
('Public_Retail', 1, 20000.00, 1000.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
('Gold_Reseller', 1, 19000.00, 500.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
('H2H_Partner', 1, 18500.00, 200.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Product Provider
INSERT INTO product_providers (provider_id, product_provider_code, product_provider_price, product_provider_admin_fee, product_provider_merchant_fee, provider_index, is_available, created_at, created_by, updated_at, updated_by) VALUES
(1, 'ML_86_IAK', 18000.00, 0.00, 0.00, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(1, 'hsteam400000', 400000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam250000', 250000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam120000', 120000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam90000', 90000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam60000', 60000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam45000', 45000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam12000', 12000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam8000', 9500, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'hsteam6000', 7000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'SYSTEM_SEED'),
(1, 'valorant11000', 1030000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant5350', 525000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant4125', 418000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant3650', 368800, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant2050', 212400, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant1000', 106200, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'valorant475', 53400, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'ragnarokbcc18', 45000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'ragnarokbcc12', 30000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'ragnarokbcc6', 15000, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'freefire12', 1940, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'freefire10', 1800, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'freefire5', 970, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'hmobilelegend10', 3100, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'hmobilelegend5', 1620, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin'),
(1, 'hmobilelegend3', 1100, 0, 0, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'admin');
