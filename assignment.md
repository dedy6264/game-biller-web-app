# Buat proyek go di lokasi existing ini dengan menyamakan framework sesuai yang ada di ../makarios
# Buat repository untuk crud berdasarkan file db.md lengkap dengan additional filter
# data master sesuai pada dbseed.md
# Bedakan route api untuk dashboard, webapp dan utils
# Kebutuhan endpoint Api atau service
    1. webapp - register
    2. webapp - login
        login akan menghasilkan jwt token
    3. webapp - get reference product
        dapat di hit tanpa login atau membawa kredensial
    4. webapp - get product segment by reference code
        dapat di hit tanpa login atau membawa kredensial
    5. webapp - get popular product by summaries reference code
        dapat di hit tanpa login atau membawa kredensial
    6. webapp - get payment method
        dapat di hit tanpa login atau membawa kredensial
    7. webapp - inquiry
        harus membawa kredensial token 
    8. webapp - payment
        harus membawa kredensial token 
    9. webapp - transaction history
        dapat di hit harus dengan membawa kredensial atau harus login dahulu
    10. dashboard - login(user admin)
    11. dashboard - crud semua fitur untuk admin
# untuk list response code mengacu pada responsecodelist.md
# untuk seluruh format json menggunakan format
{
  "status_code": "",
  "status_message": "",
  "status_desc": "",
  "ui_message": "",
  "result":{}
}


lanjut kepada api inquiry, 
memiliki alur pertama setelah validasi kredensial pengguna, lanjut validasi customer id dan validasi product melalui get product segment sesuai product code yang dikirim, jika product tidak tersedia maka inquiry gagal dan jika tersedia maka dilanjutkan generate noreff transaksi lalu dilanjutkan penyimpanan data transaksi ke database menggunakan db transaction

api payment, 
memiliki alur pengecekan transaksi berdasarkan 