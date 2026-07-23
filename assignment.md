# Inquiry
1. setelah simpan data inquiry ke database, validasi provider, jika provider_id = 1 (IAK), maka lanjutkan request ke 
host: localhost:10003/api/iak/inquiry
method : post
payload : menggunakan models.RequestInquiry
response : menggunakan models.InquiryResult

2. handle response ketika sukses update data transaksi dengan data produk dan kelengkapannya di database dengan status inquiry sukses. dan update dengan status yang lain sesuai hasil dari worker

# payment
1. setelah validasi noreff dengan mengecek didatabase bahwa transaksi valid(inquiry sukses), lanjutkan validasi provider, jika provider_id = 1 (IAK), maka lanjutkan request ke 
host: localhost:10003/api/iak/payment
method : post
payload : menggunakan models.RequestPayment
response : menggunakan models.PaymentResult

2. handle response ketika sukses update data transaksi dengan data produk dan kelengkapannya di database dengan status payment sukses.dan update dengan status yang lain sesuai hasil dari worker
3. jika saat validasi noreff invalid (selain inquiry sukses) maka return transaksi dengan status iinvalid transaksi
4. jika saat validasi noreff status transaksi payment sukses atau gagal, maka kembalikan informasi data sesuai pada di database
