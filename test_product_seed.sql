-- Insert Test Provider
INSERT INTO providers (id, provider_name, is_active, created_at, created_by, updated_at, updated_by) OVERRIDING SYSTEM VALUE VALUES
(1, 'IAK Provider', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Test Product ML_86
INSERT INTO products (id, product_reference_id, product_type_id, product_category_id, product_code, product_name, is_active, created_at, created_by, updated_at, updated_by) OVERRIDING SYSTEM VALUE VALUES
(1, 8, 1, 1, 'ML_86', 'Mobile Legends 86 Diamonds', true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Product Segment Pricing
-- 1. guest_retail -> Public_Retail
-- 2. member_premium -> Gold_Reseller
-- 3. h2h_api -> H2H_Partner
INSERT INTO product_segments (id, segment_name, product_id, product_price, admin_fee, merchant_fee, created_at, created_by, updated_at, updated_by) OVERRIDING SYSTEM VALUE VALUES
(1, 'Open_Biller', 1, 20000.00, 1000.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(2, 'Public_Retail', 1, 20000.00, 1000.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(3, 'Gold_Reseller', 1, 19000.00, 500.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(4, 'H2H_Partner', 1, 18500.00, 200.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Product Provider
INSERT INTO product_providers (provider_id, provider_product_code, provider_price, provider_admin_fee, provider_merchant_fee, provider_index, is_available, created_at, created_by, updated_at, updated_by) VALUES
(1, 'ML_86_IAK', 18000.00, 0.00, 0.00, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),

 (28, 1, "hsteam400000", 400000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (27, 1, "hsteam250000", 250000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (26, 1, "hsteam120000", 120000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (25, 1, "hsteam90000", 90000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (24, 1, "hsteam60000", 60000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (23, 1, "hsteam45000", 45000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (22, 1, "hsteam12000", 12000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (21, 1, "hsteam8000", 9500, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (20, 1, "hsteam6000", 7000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "SYSTEM_SEED"),
 (17, 1, "valorant11000", 1030000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (16, 1, "valorant5350", 525000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (15, 1, "valorant4125", 418000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (14, 1, "valorant3650", 368800, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (13, 1, "valorant2050", 212400, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (12, 1, "valorant1000", 106200, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (11, 1, "valorant475", 53400, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (10, 1, "ragnarokbcc18", 45000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (9, 1, "ragnarokbcc12", 30000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (8, 1, "ragnarokbcc6", 15000, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (7, 1, "freefire12", 1940, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (6, 1, "freefire10", 1800, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (5, 1, "freefire5", 970, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (4, 1, "hmobilelegend10", 3100, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (3, 1, "hmobilelegend5", 1620, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin"),
 (2, 1, "hmobilelegend3", 1100, 0, 0, 1, true,TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin",TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), "admin")
ON CONFLICT DO NOTHING;
