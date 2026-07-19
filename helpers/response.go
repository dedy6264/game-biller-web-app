package helpers

import (
	"fmt"
	"gamebiller/models"

	"github.com/labstack/echo/v4"
)

type ResponseMetadata struct {
	StatusCode    string
	StatusMessage string
	StatusDesc    string
	UiMessage     string
}

var responseCodes = map[string]ResponseMetadata{
	// Webpage - SUC / PEN / INQ
	"SUC-INT-000": {
		StatusCode:    "SUC-INT-000",
		StatusMessage: "PAYMENT_SUCCESS",
		StatusDesc:    "Payment confirmed and transaction successfully processed. Delivery is in progress.",
		UiMessage:     "Pembayaran berhasil! Produk sedang dikirimkan ke akun game Anda. Terima kasih.",
	},
	"INQ-SYS-001": {
		StatusCode:    "INQ-SYS-001",
		StatusMessage: "INQUIRY_SUCCESS",
		StatusDesc:    "Inquiry accepted. Product and pricing verified. Awaiting payment confirmation from merchant.",
		UiMessage:     "Informasi produk berhasil dikonfirmasi. Silakan lanjutkan ke proses pembayaran.",
	},
	"PEN-SYS-001": {
		StatusCode:    "PEN-SYS-001",
		StatusMessage: "PENDING_UPSTREAM",
		StatusDesc:    "Transaction accepted by core system and currently queued in the upstream gateway.",
		UiMessage:     "Pembayaran diterima. Pesanan Anda sedang diproses oleh sistem game. Mohon tunggu.",
	},
	"PEN-INT-002": {
		StatusCode:    "PEN-INT-002",
		StatusMessage: "PENDING_MANUAL_REVIEW",
		StatusDesc:    "Transaction suspended internally for security verification or compliance approval.",
		UiMessage:     "Transaksi Anda sedang diverifikasi demi keamanan. Mohon tunggu maksimal 5 menit.",
	},
	"PEN-SYS-003": {
		StatusCode:    "PEN-SYS-003",
		StatusMessage: "PENDING_RETRY",
		StatusDesc:    "Upstream network glitch detected. Core system is auto-switching or retrying the request.",
		UiMessage:     "Jaringan server game sedang padat. Sistem kami sedang mencoba mengirimkan ulang pesanan Anda.",
	},

	// Webpage - ERR User Input
	"ERR-VAL-100": {
		StatusCode:    "ERR-VAL-100",
		StatusMessage: "INVALID_PLAYER_ID",
		StatusDesc:    "External game validator service rejected the provided Player ID, Zone ID, or Server ID.",
		UiMessage:     "ID Game tidak terdaftar atau salah server. Silakan periksa kembali data Anda dan coba lagi.",
	},
	"ERR-INT-101": {
		StatusCode:    "ERR-INT-101",
		StatusMessage: "PRODUCT_CODE_NOT_FOUND",
		StatusDesc:    "The requested target product code or SKU is missing or inactive in our catalog.",
		UiMessage:     "Produk yang Anda pilih saat ini tidak tersedia. Silakan pilih nominal atau paket lainnya.",
	},
	"ERR-INT-102": {
		StatusCode:    "ERR-INT-102",
		StatusMessage: "INVALID_TRANSACTION_STATUS",
		StatusDesc:    "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.",
		UiMessage:     "Transaksi tidak dapat diproses. Status transaksi tidak valid atau sudah diselesaikan sebelumnya.",
	},
	"ERR-INT-103": {
		StatusCode:    "ERR-INT-103",
		StatusMessage: "INSUFFICIENT_MERCHANT_BALANCE",
		StatusDesc:    "Internal ledger check indicates the merchant deposit balance is insufficient for this checkout.",
		UiMessage:     "Saldo Anda tidak mencukupi untuk melakukan pembelian ini. Silakan top-up saldo terlebih dahulu.",
	},
	"ERR-VAL-104": {
		StatusCode:    "ERR-VAL-104",
		StatusMessage: "INVALID_TARGET_FORMAT",
		StatusDesc:    "The destination identity string failed regex or basic character formatting check.",
		UiMessage:     "Nomor tujuan atau format ID yang Anda masukkan salah. Mohon periksa kembali aturan penulisan ID.",
	},
	"ERR-INT-105": {
		StatusCode:    "ERR-INT-105",
		StatusMessage: "PRODUCT_SEGMENT_NOT_FOUND",
		StatusDesc:    "No active product segment pricing found for this product and the merchant's account type.",
		UiMessage:     "Produk ini belum tersedia untuk tipe akun Anda. Silakan hubungi Customer Service untuk informasi lebih lanjut.",
	},

	// Webpage - ERR Merchant/B2B Side
	"ERR-INT-200": {
		StatusCode:    "ERR-INT-200",
		StatusMessage: "API_AUTHENTICATION_FAILED",
		StatusDesc:    "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.",
		UiMessage:     "Sesi masuk atau koneksi tidak sah. Silakan hubungi bagian administrasi akun/dukungan teknis.",
	},
	"ERR-INT-201": {
		StatusCode:    "ERR-INT-201",
		StatusMessage: "MERCHANT_SUSPENDED",
		StatusDesc:    "Merchant corporate profile has been deactivated or restricted due to compliance or risk control.",
		UiMessage:     "Akun Anda dinonaktifkan sementara oleh sistem. Silakan hubungi Customer Service untuk verifikasi data.",
	},
	"ERR-INT-202": {
		StatusCode:    "ERR-INT-202",
		StatusMessage: "DAILY_LIMIT_EXCEEDED",
		StatusDesc:    "The total value or volume of transactions has exceeded the agreed maximum daily limit.",
		UiMessage:     "Batas kuota transaksi harian Anda telah tercapai. Anda dapat melakukan transaksi kembali esok hari.",
	},
	"ERR-INT-203": {
		StatusCode:    "ERR-INT-203",
		StatusMessage: "INVALID_API_METHOD",
		StatusDesc:    "The API endpoint requested does not support the HTTP verb or request method applied.",
		UiMessage:     "Permintaan sistem tidak didukung. Pastikan aplikasi Anda terintegrasi dengan benar.",
	},
	"ERR-INT-204": {
		StatusCode:    "ERR-INT-204",
		StatusMessage: "MAINTENANCE_SCHEDULED",
		StatusDesc:    "The client-facing dashboard and payment processing gateway are offline for scheduled system upgrade.",
		UiMessage:     "Layanan kami sedang ditingkatkan untuk kenyamanan Anda. Kami akan segera kembali dalam beberapa saat.",
	},

	// Webpage - ERR Upstream Side
	"ERR-PVD1-300": {
		StatusCode:    "ERR-PVD1-300",
		StatusMessage: "UPSTREAM_MAINTENANCE",
		StatusDesc:    "The upstream vendor gateway for this specific product is currently undergoing maintenance.",
		UiMessage:     "Server pusat untuk game ini sedang mengalami pemeliharaan sistem. Silakan coba beberapa saat lagi.",
	},
	"ERR-PVD2-301": {
		StatusCode:    "ERR-PVD2-301",
		StatusMessage: "PRODUCT_OUT_OF_STOCK",
		StatusDesc:    "The upstream distributor core reported that the stock, voucher pool, or product quota is exhausted.",
		UiMessage:     "Stok untuk item game ini sedang habis di pusat. Tim kami sedang melakukan pengisian ulang stok.",
	},
	"ERR-PVD1-302": {
		StatusCode:    "ERR-PVD1-302",
		StatusMessage: "UPSTREAM_TIMEOUT",
		StatusDesc:    "Upstream server failed to respond within the designated execution window. Money safely refunded.",
		UiMessage:     "Jaringan ke server game sedang padat. Transaksi dibatalkan secara aman dan saldo Anda tetap utuh.",
	},
	"ERR-PVD3-303": {
		StatusCode:    "ERR-PVD3-303",
		StatusMessage: "UPSTREAM_DECLINED",
		StatusDesc:    "The execution request was rejected by the distribution core gateway due to structural rule conflicts.",
		UiMessage:     "Pembelian ditolak oleh sistem pusat game. Silakan pilih nominal atau paket game lainnya.",
	},
	"ERR-PVD2-304": {
		StatusCode:    "ERR-PVD2-304",
		StatusMessage: "UPSTREAM_PRICE_CHANGED",
		StatusDesc:    "The purchase cost from the upstream gateway has changed and exceeded the configured system margin.",
		UiMessage:     "Terjadi pembaruan harga dari sistem pusat game. Silakan ulangi transaksi Anda untuk memperbarui harga.",
	},
	"ERR-PVD4-305": {
		StatusCode:    "ERR-PVD4-305",
		StatusMessage: "UPSTREAM_UNKNOWN_ERROR",
		StatusDesc:    "Upstream provider responded with an unmapped error structure or critical raw payload exception.",
		UiMessage:     "Terjadi kendala pada jaringan distribusi game. Saldo tidak terpotong, silakan coba lagi nanti.",
	},

	// Dashboard - AUTH
	"SUC-AUTH-200": {
		StatusCode:    "SUC-AUTH-200",
		StatusMessage: "AUTHENTICATION_SUCCESS",
		StatusDesc:    "User credentials verified. Session token generated successfully.",
		UiMessage:     "Login berhasil! Mengalihkan Anda ke halaman dashboard...",
	},
	"SUC-AUTH-201": {
		StatusCode:    "SUC-AUTH-201",
		StatusMessage: "REGISTRATION_SUCCESS",
		StatusDesc:    "New user entity successfully created inside the core database.",
		UiMessage:     "Pendaftaran berhasil! Selamat bergabung di platform kami.",
	},
	"ERR-AUTH-401": {
		StatusCode:    "ERR-AUTH-401",
		StatusMessage: "INVALID_CREDENTIALS",
		StatusDesc:    "The password hash or username/email combination did not match database records.",
		UiMessage:     "Email atau password yang Anda masukkan salah. Silakan coba lagi.",
	},
	"ERR-AUTH-403": {
		StatusCode:    "ERR-AUTH-403",
		StatusMessage: "ACCESS_DENIED",
		StatusDesc:    "The authenticated user does not have the required role assigned in model_has_roles.",
		UiMessage:     "Akun Anda tidak memiliki izin untuk mengakses halaman ini.",
	},
	"ERR-AUTH-419": {
		StatusCode:    "ERR-AUTH-419",
		StatusMessage: "TOKEN_EXPIRED",
		StatusDesc:    "The bearer token or session cookie has expired or been blacklisted.",
		UiMessage:     "Sesi Anda telah berakhir demi keamanan. Silakan masuk (login) kembali.",
	},

	// Dashboard - USER / VAL
	"SUC-USER-200": {
		StatusCode:    "SUC-USER-200",
		StatusMessage: "PROFILE_UPDATED",
		StatusDesc:    "Profile fields successfully modified and committed to the user entity.",
		UiMessage:     "Perubahan profil Anda telah berhasil disimpan.",
	},
	"VAL-USER-422": {
		StatusCode:    "VAL-USER-422",
		StatusMessage: "EMAIL_ALREADY_EXISTS",
		StatusDesc:    "Unique constraint failed. The provided email address is already taken by another account.",
		UiMessage:     "Email sudah terdaftar. Silakan gunakan email lain atau gunakan fitur lupa password.",
	},
	"VAL-USER-423": {
		StatusCode:    "VAL-USER-423",
		StatusMessage: "PHONE_ALREADY_EXISTS",
		StatusDesc:    "Unique constraint failed. The phone number is already registered in the system.",
		UiMessage:     "Nomor HP sudah digunakan. Mohon periksa kembali nomor Anda.",
	},
	"VAL-USER-424": {
		StatusCode:    "VAL-USER-424",
		StatusMessage: "WEAK_PASSWORD",
		StatusDesc:    "Validation failed. The password policy string requirements are not met.",
		UiMessage:     "Password terlalu lemah. Gunakan minimal 8 karakter dengan kombinasi angka dan huruf.",
	},
	"ERR-USER-404": {
		StatusCode:    "ERR-USER-404",
		StatusMessage: "USER_NOT_FOUND",
		StatusDesc:    "Query returned no records for the specified User ID or identifier.",
		UiMessage:     "Data pengguna tidak ditemukan di dalam sistem kami.",
	},

	// Dashboard - MERCH
	"SUC-MERCH-200": {
		StatusCode:    "SUC-MERCH-200",
		StatusMessage: "API_CREDENTIALS_REGENERATED",
		StatusDesc:    "New client_key and secret_key successfully rolled over for the merchant entity.",
		UiMessage:     "API Key baru berhasil dibuat. Pastikan Anda segera memperbarui sistem H2H Anda.",
	},
	"SUC-MERCH-201": {
		StatusCode:    "SUC-MERCH-201",
		StatusMessage: "IP_WHITELIST_UPDATED",
		StatusDesc:    "Whitelist IP text area parsed and updated inside merchant_api_credentials.",
		UiMessage:     "Daftar IP Whitelist Anda telah berhasil diperbarui.",
	},
	"ERR-MERCH-400": {
		StatusCode:    "ERR-MERCH-400",
		StatusMessage: "INVALID_IP_FORMAT",
		StatusDesc:    "Failed to parse the provided whitelist string. IP address format is invalid.",
		UiMessage:     "Format IP Address yang Anda masukkan salah. Mohon periksa kembali.",
	},
	"ERR-MERCH-403": {
		StatusCode:    "ERR-MERCH-403",
		StatusMessage: "WRONG_TRANSACTION_PIN",
		StatusDesc:    "The secure transaction PIN validation failed against account_pin_hash.",
		UiMessage:     "PIN Transaksi yang Anda masukkan salah. Sisa kesempatan 2 kali lagi.",
	},

	// Dashboard - UTIL
	"SUC-UTIL-200": {
		StatusCode:    "SUC-UTIL-200",
		StatusMessage: "OTP_SENT",
		StatusDesc:    "Verification code generated and successfully pushed to the notification provider.",
		UiMessage:     "Kode verifikasi (OTP) telah dikirimkan ke nomor HP / email Anda.",
	},
	"SUC-UTIL-202": {
		StatusCode:    "SUC-UTIL-202",
		StatusMessage: "OTP_VERIFIED",
		StatusDesc:    "OTP token matches the active otp_codes row and is successfully marked as used.",
		UiMessage:     "Verifikasi berhasil. Silakan melanjutkan ke proses berikutnya.",
	},
	"ERR-UTIL-408": {
		StatusCode:    "ERR-UTIL-408",
		StatusMessage: "OTP_EXPIRED",
		StatusDesc:    "The system matched the code but the active window passed the expired_at timestamp.",
		UiMessage:     "Kode OTP sudah kedaluwarsa. Silakan klik tombol \"Kirim Ulang\".",
	},
	"ERR-UTIL-410": {
		StatusCode:    "ERR-UTIL-410",
		StatusMessage: "OTP_INCORRECT",
		StatusDesc:    "The input OTP string does not match the active sequence stored.",
		UiMessage:     "Kode OTP yang Anda masukkan salah. Mohon periksa kembali pesan masuk Anda.",
	},
	"ERR-UTIL-429": {
		StatusCode:    "ERR-UTIL-429",
		StatusMessage: "OTP_MAX_ATTEMPT_REACHED",
		StatusDesc:    "The counter on attempt_count has exceeded max_attempt. The token is invalidated.",
		UiMessage:     "Anda telah salah memasukkan OTP sebanyak 3 kali. Silakan minta kode OTP baru.",
	},
}

// BuildResponse generates a standardized models.ApiResponse
func BuildResponse(code string, result any) models.ApiResponse {
	meta, exists := responseCodes[code]
	if !exists {
		return models.ApiResponse{
			StatusCode:    "ERR-SYS-500",
			StatusMessage: "SYSTEM_ERROR",
			StatusDesc:    "Internal system error occurred.",
			UiMessage:     "Terjadi gangguan sistem. Silakan coba beberapa saat lagi.",
			Result:        result,
		}
	}

	if result == nil {
		result = map[string]any{}
	}

	return models.ApiResponse{
		StatusCode:    meta.StatusCode,
		StatusMessage: meta.StatusMessage,
		StatusDesc:    meta.StatusDesc,
		UiMessage:     meta.UiMessage,
		Result:        result,
	}
}
func ProcessLogger(c echo.Context, svc string, message any, desc string) {
	// Log request details
	fmt.Println("Error :: ", svc, message, desc)
	// log.Printf("Request: %s %s", c.Request().Method, c.Request().URL)
}
