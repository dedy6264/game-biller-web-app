

# Mengisi Master Roles
INSERT INTO roles (id, role_code, role_name, created_at, updated_at) VALUES
(1, 'super_admin', 'Super Administrator Internal', NOW(), NOW()),
(2, 'finance', 'Finance & Billing Internal', NOW(), NOW()),
(3, 'merchant_h2h', 'Mitra Host-to-Host API', NOW(), NOW()),
(4, 'member_reseller', 'Reseller VIP Dashboard', NOW(), NOW()),
(5, 'retail_guest', 'Pembeli Lepas Web Retail', NOW(), NOW());


# Product Types (Pembeda Logika Finansial Makro)
INSERT INTO product_types (id, product_type_name) VALUES
(1, 'Prepaid'),
(2, 'Postpaid');

# Product Categories (Pengelompokan Komoditas Ritel)
INSERT INTO product_categories (id, product_category_name) VALUES
(1, 'Game Top Up'),
(2, 'Pulsa & Data'),
(3, 'E-Wallet'),
(4, 'Tagihan PLN & PDAM'),


# Product References (Grouping untuk Jalur Switching/Routing Gateway)
INSERT INTO product_references (id, product_reference_code, product_reference_name, created_at, updated_at) VALUES
(1, 'ref_tsel', 'TELKOMSEL', NOW(), NOW()),
(2, 'ref_isat', 'INDOSAT', NOW(), NOW()),
(3, 'ref_three', 'THREE', NOW(), NOW()),
(4, 'ref_axis', 'AXIS', NOW(), NOW()),
(5, 'ref_smart', 'SMARTFREN', NOW(), NOW()),
(6, 'ref_xl', 'XL', NOW(), NOW()),
(7, 'ref_byu', 'BY.U', NOW(), NOW()),
(8, 'ref_mlbb', 'MOBILE LEGEND', NOW(), NOW()),
(9, 'ref_genshin', 'GENSHIN IMPACT', NOW(), NOW()),
(10, 'ref_ff', 'FREE FIRE', NOW(), NOW());

# Product Prefixes (Otomatisasi Deteksi Operator Berdasarkan Input User)
## PREFIX TELKOMSEL (Reference ID: 1)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(1, '08120', NOW(), NOW()), (1, '08121', NOW(), NOW()), (1, '08122', NOW(), NOW()), (1, '08123', NOW(), NOW()), (1, '08124', NOW(), NOW()),
(1, '08125', NOW(), NOW()), (1, '08126', NOW(), NOW()), (1, '08127', NOW(), NOW()), (1, '08128', NOW(), NOW()), (1, '08129', NOW(), NOW()),
(1, '08130', NOW(), NOW()), (1, '08131', NOW(), NOW()), (1, '08132', NOW(), NOW()), (1, '08133', NOW(), NOW()), (1, '08134', NOW(), NOW()),
(1, '08135', NOW(), NOW()), (1, '08136', NOW(), NOW()), (1, '08137', NOW(), NOW()), (1, '08138', NOW(), NOW()), (1, '08139', NOW(), NOW()),
(1, '08520', NOW(), NOW()), (1, '08521', NOW(), NOW()), (1, '08522', NOW(), NOW()), (1, '08523', NOW(), NOW()), (1, '08524', NOW(), NOW()),
(1, '08525', NOW(), NOW()), (1, '08526', NOW(), NOW()), (1, '08527', NOW(), NOW()), (1, '08528', NOW(), NOW()), (1, '08529', NOW(), NOW()),
(1, '08530', NOW(), NOW()), (1, '08531', NOW(), NOW()), (1, '08532', NOW(), NOW()), (1, '08533', NOW(), NOW()), (1, '08534', NOW(), NOW()),
(1, '08535', NOW(), NOW()), (1, '08536', NOW(), NOW()), (1, '08537', NOW(), NOW()), (1, '08538', NOW(), NOW()), (1, '08539', NOW(), NOW()),
(1, '08210', NOW(), NOW()), (1, '08211', NOW(), NOW()), (1, '08212', NOW(), NOW()), (1, '08213', NOW(), NOW()), (1, '08214', NOW(), NOW()),
(1, '08215', NOW(), NOW()), (1, '08216', NOW(), NOW()), (1, '08217', NOW(), NOW()), (1, '08218', NOW(), NOW()), (1, '08219', NOW(), NOW()),
(1, '08230', NOW(), NOW()), (1, '08231', NOW(), NOW()), (1, '08232', NOW(), NOW()), (1, '08233', NOW(), NOW()), (1, '08234', NOW(), NOW()),
(1, '08235', NOW(), NOW()), (1, '08236', NOW(), NOW()), (1, '08237', NOW(), NOW()), (1, '08238', NOW(), NOW()), (1, '08239', NOW(), NOW()),
(1, '08220', NOW(), NOW()), (1, '08221', NOW(), NOW()), (1, '08222', NOW(), NOW()), (1, '08223', NOW(), NOW()), (1, '08224', NOW(), NOW()),
(1, '08225', NOW(), NOW()), (1, '08226', NOW(), NOW()), (1, '08227', NOW(), NOW()), (1, '08228', NOW(), NOW()), (1, '08229', NOW(), NOW());

## PREFIX INDOSAT (Reference ID: 2)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(2, '08140', NOW(), NOW()), (2, '08141', NOW(), NOW()), (2, '08142', NOW(), NOW()), (2, '08143', NOW(), NOW()), (2, '08144', NOW(), NOW()),
(2, '08145', NOW(), NOW()), (2, '08146', NOW(), NOW()), (2, '08147', NOW(), NOW()), (2, '08148', NOW(), NOW()), (2, '08149', NOW(), NOW()),
(2, '08150', NOW(), NOW()), (2, '08151', NOW(), NOW()), (2, '08152', NOW(), NOW()), (2, '08153', NOW(), NOW()), (2, '08154', NOW(), NOW()),
(2, '08155', NOW(), NOW()), (2, '08156', NOW(), NOW()), (2, '08157', NOW(), NOW()), (2, '08158', NOW(), NOW()), (2, '08159', NOW(), NOW()),
(2, '08160', NOW(), NOW()), (2, '08161', NOW(), NOW()), (2, '08162', NOW(), NOW()), (2, '08163', NOW(), NOW()), (2, '08164', NOW(), NOW()),
(2, '08165', NOW(), NOW()), (2, '08166', NOW(), NOW()), (2, '08166', NOW(), NOW()), (2, '08167', NOW(), NOW()), (2, '08168', NOW(), NOW()),
(2, '08169', NOW(), NOW()), (2, '08550', NOW(), NOW()), (2, '08551', NOW(), NOW()), (2, '08552', NOW(), NOW()), (2, '08553', NOW(), NOW()),
(2, '08554', NOW(), NOW()), (2, '08555', NOW(), NOW()), (2, '08556', NOW(), NOW()), (2, '08557', NOW(), NOW()), (2, '08558', NOW(), NOW()),
(2, '08559', NOW(), NOW()), (2, '08560', NOW(), NOW()), (2, '08561', NOW(), NOW()), (2, '08562', NOW(), NOW()), (2, '08563', NOW(), NOW()),
(2, '08564', NOW(), NOW()), (2, '08565', NOW(), NOW()), (2, '08566', NOW(), NOW()), (2, '08567', NOW(), NOW()), (2, '08568', NOW(), NOW()),
(2, '08569', NOW(), NOW()), (2, '08570', NOW(), NOW()), (2, '08571', NOW(), NOW()), (2, '08572', NOW(), NOW()), (2, '08573', NOW(), NOW()),
(2, '08574', NOW(), NOW()), (2, '08575', NOW(), NOW()), (2, '08576', NOW(), NOW()), (2, '08577', NOW(), NOW()), (2, '08578', NOW(), NOW()),
(2, '08579', NOW(), NOW()), (2, '08580', NOW(), NOW()), (2, '08581', NOW(), NOW()), (2, '08582', NOW(), NOW()), (2, '08583', NOW(), NOW()),
(2, '08584', NOW(), NOW()), (2, '08585', NOW(), NOW()), (2, '08586', NOW(), NOW()), (2, '08587', NOW(), NOW()), (2, '08588', NOW(), NOW()),
(2, '08589', NOW(), NOW());

## PREFIX THREE (Reference ID: 3)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(3, '08960', NOW(), NOW()), (3, '08961', NOW(), NOW()), (3, '08962', NOW(), NOW()), (3, '08963', NOW(), NOW()), (3, '08964', NOW(), NOW()),
(3, '08965', NOW(), NOW()), (3, '08966', NOW(), NOW()), (3, '08967', NOW(), NOW()), (3, '08968', NOW(), NOW()), (3, '08969', NOW(), NOW()),
(3, '08970', NOW(), NOW()), (3, '08971', NOW(), NOW()), (3, '08972', NOW(), NOW()), (3, '08973', NOW(), NOW()), (3, '08974', NOW(), NOW()),
(3, '08975', NOW(), NOW()), (3, '08976', NOW(), NOW()), (3, '08977', NOW(), NOW()), (3, '08978', NOW(), NOW()), (3, '08979', NOW(), NOW()),
(3, '08980', NOW(), NOW()), (3, '08981', NOW(), NOW()), (3, '08982', NOW(), NOW()), (3, '08983', NOW(), NOW()), (3, '08984', NOW(), NOW()),
(3, '08985', NOW(), NOW()), (3, '08986', NOW(), NOW()), (3, '08987', NOW(), NOW()), (3, '08988', NOW(), NOW()), (3, '08989', NOW(), NOW()),
(3, '08990', NOW(), NOW()), (3, '08991', NOW(), NOW()), (3, '08992', NOW(), NOW()), (3, '08993', NOW(), NOW()), (3, '08994', NOW(), NOW()),
(3, '08995', NOW(), NOW()), (3, '08996', NOW(), NOW()), (3, '08997', NOW(), NOW()), (3, '08998', NOW(), NOW()), (3, '08999', NOW(), NOW()),
(3, '08950', NOW(), NOW()), (3, '08951', NOW(), NOW()), (3, '08952', NOW(), NOW()), (3, '08953', NOW(), NOW()), (3, '08954', NOW(), NOW()),
(3, '08955', NOW(), NOW()), (3, '08956', NOW(), NOW()), (3, '08957', NOW(), NOW()), (3, '08958', NOW(), NOW()), (3, '08959', NOW(), NOW());

## PREFIX AXIS (Reference ID: 4)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(4, '08380', NOW(), NOW()), (4, '08381', NOW(), NOW()), (4, '08382', NOW(), NOW()), (4, '08383', NOW(), NOW()), (4, '08384', NOW(), NOW()),
(4, '08385', NOW(), NOW()), (4, '08386', NOW(), NOW()), (4, '08387', NOW(), NOW()), (4, '08388', NOW(), NOW()), (4, '08389', NOW(), NOW()),
(4, '08370', NOW(), NOW()), (4, '08371', NOW(), NOW()), (4, '08372', NOW(), NOW()), (4, '08373', NOW(), NOW()), (4, '08374', NOW(), NOW()),
(4, '08375', NOW(), NOW()), (4, '08376', NOW(), NOW()), (4, '08377', NOW(), NOW()), (4, '08378', NOW(), NOW()), (4, '08379', NOW(), NOW()),
(4, '08310', NOW(), NOW()), (4, '08311', NOW(), NOW()), (4, '08312', NOW(), NOW()), (4, '08313', NOW(), NOW()), (4, '08314', NOW(), NOW()),
(4, '08315', NOW(), NOW()), (4, '08316', NOW(), NOW()), (4, '08317', NOW(), NOW()), (4, '08318', NOW(), NOW()), (4, '08319', NOW(), NOW()),
(4, '08320', NOW(), NOW()), (4, '08321', NOW(), NOW()), (4, '08322', NOW(), NOW()), (4, '08323', NOW(), NOW()), (4, '08324', NOW(), NOW()),
(4, '08325', NOW(), NOW()), (4, '08326', NOW(), NOW()), (4, '08327', NOW(), NOW()), (4, '08328', NOW(), NOW()), (4, '08329', NOW(), NOW());

## PREFIX SMARTFREN (Reference ID: 5)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(5, '08810', NOW(), NOW()), (5, '08811', NOW(), NOW()), (5, '08812', NOW(), NOW()), (5, '08813', NOW(), NOW()), (5, '08814', NOW(), NOW()),
(5, '08815', NOW(), NOW()), (5, '08816', NOW(), NOW()), (5, '08817', NOW(), NOW()), (5, '08818', NOW(), NOW()), (5, '08819', NOW(), NOW()),
(5, '08820', NOW(), NOW()), (5, '08821', NOW(), NOW()), (5, '08822', NOW(), NOW()), (5, '08823', NOW(), NOW()), (5, '08824', NOW(), NOW()),
(5, '08825', NOW(), NOW()), (5, '08826', NOW(), NOW()), (5, '08827', NOW(), NOW()), (5, '08828', NOW(), NOW()), (5, '08829', NOW(), NOW()),
(5, '08830', NOW(), NOW()), (5, '08831', NOW(), NOW()), (5, '08832', NOW(), NOW()), (5, '08833', NOW(), NOW()), (5, '08834', NOW(), NOW()),
(5, '08835', NOW(), NOW()), (5, '08836', NOW(), NOW()), (5, '08837', NOW(), NOW()), (5, '08838', NOW(), NOW()), (5, '08839', NOW(), NOW()),
(5, '08840', NOW(), NOW()), (5, '08841', NOW(), NOW()), (5, '08842', NOW(), NOW()), (5, '08843', NOW(), NOW()), (5, '08844', NOW(), NOW()),
(5, '08845', NOW(), NOW()), (5, '08846', NOW(), NOW()), (5, '08847', NOW(), NOW()), (5, '08848', NOW(), NOW()), (5, '08849', NOW(), NOW()),
(5, '08850', NOW(), NOW()), (5, '08851', NOW(), NOW()), (5, '08852', NOW(), NOW()), (5, '08853', NOW(), NOW()), (5, '08854', NOW(), NOW()),
(5, '08855', NOW(), NOW()), (5, '08856', NOW(), NOW()), (5, '08857', NOW(), NOW()), (5, '08858', NOW(), NOW()), (5, '08859', NOW(), NOW()),
(5, '08860', NOW(), NOW()), (5, '08861', NOW(), NOW()), (5, '08862', NOW(), NOW()), (5, '08863', NOW(), NOW()), (5, '08864', NOW(), NOW()),
(5, '08865', NOW(), NOW()), (5, '08866', NOW(), NOW()), (5, '08867', NOW(), NOW()), (5, '08868', NOW(), NOW()), (5, '08869', NOW(), NOW()),
(5, '08870', NOW(), NOW()), (5, '08871', NOW(), NOW()), (5, '08872', NOW(), NOW()), (5, '08873', NOW(), NOW()), (5, '08874', NOW(), NOW()),
(5, '08875', NOW(), NOW()), (5, '08876', NOW(), NOW()), (5, '08877', NOW(), NOW()), (5, '08878', NOW(), NOW()), (5, '08879', NOW(), NOW()),
(5, '08880', NOW(), NOW()), (5, '08881', NOW(), NOW()), (5, '08882', NOW(), NOW()), (5, '08883', NOW(), NOW()), (5, '08884', NOW(), NOW()),
(5, '08885', NOW(), NOW()), (5, '08886', NOW(), NOW()), (5, '08887', NOW(), NOW()), (5, '08888', NOW(), NOW()), (5, '08889', NOW(), NOW());

## PREFIX XL (Reference ID: 6)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(6, '08170', NOW(), NOW()), (6, '08171', NOW(), NOW()), (6, '08172', NOW(), NOW()), (6, '08173', NOW(), NOW()), (6, '08174', NOW(), NOW()),
(6, '08175', NOW(), NOW()), (6, '08176', NOW(), NOW()), (6, '08177', NOW(), NOW()), (6, '08178', NOW(), NOW()), (6, '08179', NOW(), NOW()),
(6, '08180', NOW(), NOW()), (6, '08181', NOW(), NOW()), (6, '08182', NOW(), NOW()), (6, '08183', NOW(), NOW()), (6, '08184', NOW(), NOW()),
(6, '08185', NOW(), NOW()), (6, '08186', NOW(), NOW()), (6, '08187', NOW(), NOW()), (6, '08188', NOW(), NOW()), (6, '08189', NOW(), NOW()),
(6, '08190', NOW(), NOW()), (6, '08191', NOW(), NOW()), (6, '08192', NOW(), NOW()), (6, '08193', NOW(), NOW()), (6, '08194', NOW(), NOW()),
(6, '08195', NOW(), NOW()), (6, '08196', NOW(), NOW()), (6, '08197', NOW(), NOW()), (6, '08198', NOW(), NOW()), (6, '08199', NOW(), NOW()),
(6, '08590', NOW(), NOW()), (6, '08591', NOW(), NOW()), (6, '08592', NOW(), NOW()), (6, '08593', NOW(), NOW()), (6, '08594', NOW(), NOW()),
(6, '08595', NOW(), NOW()), (6, '08596', NOW(), NOW()), (6, '08597', NOW(), NOW()), (6, '08598', NOW(), NOW()), (6, '08599', NOW(), NOW()),
(6, '08780', NOW(), NOW()), (6, '08781', NOW(), NOW()), (6, '08782', NOW(), NOW()), (6, '08783', NOW(), NOW()), (6, '08784', NOW(), NOW()),
(6, '08785', NOW(), NOW()), (6, '08786', NOW(), NOW()), (6, '08787', NOW(), NOW()), (6, '08788', NOW(), NOW()), (6, '08789', NOW(), NOW()),
(6, '08770', NOW(), NOW()), (6, '08771', NOW(), NOW()), (6, '08772', NOW(), NOW()), (6, '08773', NOW(), NOW()), (6, '08774', NOW(), NOW()),
(6, '08775', NOW(), NOW()), (6, '08776', NOW(), NOW()), (6, '08777', NOW(), NOW()), (6, '08778', NOW(), NOW()), (6, '08779', NOW(), NOW());

## PREFIX BY.U (Reference ID: 7)
INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, updated_at) VALUES
(7, '08510', NOW(), NOW()), (7, '08511', NOW(), NOW()), (7, '08512', NOW(), NOW()), (7, '08513', NOW(), NOW()), (7, '08514', NOW(), NOW()),
(7, '08515', NOW(), NOW()), (7, '08516', NOW(), NOW()), (7, '08517', NOW(), NOW()), (7, '08518', NOW(), NOW()), (7, '08519', NOW(), NOW());

# Master Data Payment Methods
INSERT INTO payment_methods (id, method_code, method_name, created_at, updated_at) VALUES
(1, 'DEPOSIT', 'Saldo Deposit internal', NOW(), NOW()),
(2, 'VIRTUAL_ACCOUNT', 'Virtual Account Transfer', NOW(), NOW()),
(3, 'QRIS', 'QR Code QRIS', NOW(), NOW()),
(4, 'E_WALLET', 'Dompet Digital / E-Wallet', NOW(), NOW());

# Master Data Payment Channels
INSERT INTO payment_channels (id, payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, updated_at) VALUES
-- Jalur potong saldo deposit (Utamanya untuk mitra H2H / Member Premium)
(1, 1, 'BALANCE_INTERNAL', 'Saldo Deposit Akun', 'FIXED', 0.00, true, NOW(), NOW()),

-- Jalur Ritel Web Top Up via Virtual Account (Fee flat per transaksi)
(2, 2, 'BCA_VA', 'BCA Virtual Account', 'FIXED', 4000.00, true, NOW(), NOW()),
(3, 2, 'MANDIRI_VA', 'Mandiri Virtual Account', 'FIXED', 3500.00, true, NOW(), NOW()),

-- Jalur Ritel Web Top Up via QRIS (Fee persentase MDR dari nilai transaksi)
(4, 3, 'QRIS_GATEWAY', 'QRIS Dana/LinkAja (All Shopee/Gopay)', 'PERCENTAGE', 0.70, true, NOW(), NOW()), -- 0.7% MDR

-- Jalur Ritel via E-Wallet Direct Link
(5, 4, 'OVO_DIRECT', 'OVO Instant Payment', 'PERCENTAGE', 1.50, true, NOW(), NOW()); -- 1.5% Fee