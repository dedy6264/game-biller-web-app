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
(1, 'Public_Retail', 1, 20000.00, 1000.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(2, 'Gold_Reseller', 1, 19000.00, 500.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system'),
(3, 'H2H_Partner', 1, 18500.00, 200.00, 0.00, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');

-- Insert Product Provider
INSERT INTO product_providers (id, provider_id, product_id, provider_product_code, provider_price, provider_admin_fee, provider_index, is_available, created_at, created_by, updated_at, updated_by) OVERRIDING SYSTEM VALUE VALUES
(1, 1, 1, 'ML_86_IAK', 18000.00, 0.00, 1, true, TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system', TO_CHAR(NOW(), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'system');
