# tambahkan field pada tabel transactions field berikut:
    merchant_name
    product_name
    product_segment_name
    product_provider_name
    provider_name
    product_type_name
    payment_channel_name
    product_merchant_fee
    product_provider_admin_fee
    product_provider_merchant_fee
# ubah nama field berikut menjadi yang telah di deskripsikan
    buy_price -> product_provider_price
    sell_price -> product_price
    admin_fee -> product_admin_fee
    payment_fee -> payment_admin_fee
    target_user_id -> customer_id
# sesuaikan juga pada tabel lainnya yang berkaitan dengan tabel transactions
