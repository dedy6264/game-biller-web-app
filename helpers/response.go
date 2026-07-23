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

// Response Code Constants
const (
	// Webpage - SUC / PEN / INQ
	CodeSuccess     = "SUC-INT-000" //Success
	CodeInqSuccess  = "INQ-INT-001" //Inq Success
	CodeIntrPending = "PEN-SYS-001" //pending

	// Webpage - ERR User Input
	CodeInvalidIdGame          = "ERR-VAL-100"
	CodeInvalidProductNotFound = "ERR-INT-101"
	CodeInvalidTransaction     = "ERR-INT-102"
	CodeErrIntBalance          = "ERR-INT-103"
	CodeInvalidCustId          = "ERR-VAL-104"
	CodeInvalidProductSegmnt   = "ERR-INT-105"

	// Webpage - ERR Merchant/B2B Side
	CodeErrInt200 = "ERR-INT-200"
	CodeErrInt201 = "ERR-INT-201"
	CodeErrInt202 = "ERR-INT-202"
	CodeErrInt203 = "ERR-INT-203"
	CodeErrInt204 = "ERR-INT-204"

	// Webpage - ERR Upstream Side
	CodeErrPvd1300 = "ERR-PVD1-300"
	CodeErrPvd2301 = "ERR-PVD2-301"
	CodeErrPvd1302 = "ERR-PVD1-302"
	CodeErrPvd3303 = "ERR-PVD3-303"
	CodeErrPvd2304 = "ERR-PVD2-304"
	CodeErrPvd4305 = "ERR-PVD4-305"

	// Dashboard - AUTH
	CodeSuccessAuth = "SUC-AUTH-200"
	CodeSucAuth201  = "SUC-AUTH-201"
	CodeErrAuth401  = "ERR-AUTH-401"
	CodeErrAuth403  = "ERR-AUTH-403"
	CodeErrAuth419  = "ERR-AUTH-419"

	// Dashboard - USER / VAL
	CodeSucUser200 = "SUC-USER-200"
	CodeValUser422 = "VAL-USER-422"
	CodeValUser423 = "VAL-USER-423"
	CodeValUser424 = "VAL-USER-424"
	CodeErrUser404 = "ERR-USER-404"

	// Dashboard - MERCH
	CodeSucMerch200 = "SUC-MERCH-200"
	CodeSucMerch201 = "SUC-MERCH-201"
	CodeErrMerch400 = "ERR-MERCH-400"
	CodeErrMerch403 = "ERR-MERCH-403"

	// Dashboard - UTIL
	CodeSucUtil200 = "SUC-UTIL-200"
	CodeSucUtil202 = "SUC-UTIL-202"
	CodeErrUtil408 = "ERR-UTIL-408"
	CodeErrUtil410 = "ERR-UTIL-410"
	CodeErrUtil429 = "ERR-UTIL-429"

	// System Error
	CodeErrSys404 = "ERR-SYS-404"
	CodeErrSys500 = "ERR-SYS-500"
)

var responseCodes = map[string]ResponseMetadata{
	// Webpage - SUC / PEN / INQ
	CodeSuccess: {
		StatusCode:    CodeSuccess,
		StatusMessage: "PAYMENT_SUCCESS",
		StatusDesc:    "Payment confirmed and transaction successfully processed. Delivery is in progress.",
		UiMessage:     "Pembayaran berhasil! Produk sedang dikirimkan ke akun game Anda. Terima kasih.",
	},
	CodeInqSuccess: {
		StatusCode:    CodeInqSuccess,
		StatusMessage: "INQUIRY_SUCCESS",
		StatusDesc:    "Inquiry accepted. Product and pricing verified. Awaiting payment confirmation from merchant.",
		UiMessage:     "Informasi produk berhasil dikonfirmasi. Silakan lanjutkan ke proses pembayaran.",
	},
	CodeIntrPending: {
		StatusCode:    CodeIntrPending,
		StatusMessage: "PENDING_UPSTREAM",
		StatusDesc:    "Transaction accepted by core system and currently queued in the upstream gateway.",
		UiMessage:     "Pembayaran diterima. Pesanan Anda sedang diproses oleh sistem game. Mohon tunggu.",
	},

	// Webpage - ERR User Input
	CodeInvalidIdGame: {
		StatusCode:    CodeInvalidIdGame,
		StatusMessage: "INVALID_PLAYER_ID",
		StatusDesc:    "External game validator service rejected the provided Player ID, Zone ID, or Server ID.",
		UiMessage:     "ID Game tidak terdaftar atau salah server. Silakan periksa kembali data Anda dan coba lagi.",
	},
	CodeInvalidProductNotFound: {
		StatusCode:    CodeInvalidProductNotFound,
		StatusMessage: "PRODUCT_CODE_NOT_FOUND",
		StatusDesc:    "The requested target product code or SKU is missing or inactive in our catalog.",
		UiMessage:     "Produk yang Anda pilih saat ini tidak tersedia. Silakan pilih nominal atau paket lainnya.",
	},
	CodeInvalidTransaction: {
		StatusCode:    CodeInvalidTransaction,
		StatusMessage: "INVALID_TRANSACTION_STATUS",
		StatusDesc:    "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.",
		UiMessage:     "Transaksi tidak dapat diproses. Status transaksi tidak valid atau sudah diselesaikan sebelumnya.",
	},
	CodeErrIntBalance: {
		StatusCode:    CodeErrIntBalance,
		StatusMessage: "INSUFFICIENT_MERCHANT_BALANCE",
		StatusDesc:    "Internal ledger check indicates the merchant deposit balance is insufficient for this checkout.",
		UiMessage:     "Keliatannya saldomu kurang. Tenang, cukup top up ulang lalu mari lanjutkan transaksi",
	},
	CodeInvalidCustId: {
		StatusCode:    CodeInvalidCustId,
		StatusMessage: "INVALID_TARGET_FORMAT",
		StatusDesc:    "The destination identity string failed regex or basic character formatting check.",
		UiMessage:     "Nomor tujuan atau format ID yang Anda masukkan salah. Mohon periksa kembali aturan penulisan ID.",
	},
	CodeInvalidProductSegmnt: {
		StatusCode:    CodeInvalidProductSegmnt,
		StatusMessage: "PRODUCT_SEGMENT_NOT_FOUND",
		StatusDesc:    "No active product segment pricing found for this product and the merchant's account type.",
		UiMessage:     "Produk ini belum tersedia untuk tipe akun Anda. Silakan hubungi Customer Service untuk informasi lebih lanjut.",
	},

	// Webpage - ERR Merchant/B2B Side
	CodeErrInt200: {
		StatusCode:    CodeErrInt200,
		StatusMessage: "API_AUTHENTICATION_FAILED",
		StatusDesc:    "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.",
		UiMessage:     "Sesi masuk atau koneksi tidak sah. Silakan hubungi bagian administrasi akun/dukungan teknis.",
	},
	CodeErrInt201: {
		StatusCode:    CodeErrInt201,
		StatusMessage: "MERCHANT_SUSPENDED",
		StatusDesc:    "Merchant corporate profile has been deactivated or restricted due to compliance or risk control.",
		UiMessage:     "Akun Anda dinonaktifkan sementara oleh sistem. Silakan hubungi Customer Service untuk verifikasi data.",
	},
	CodeErrInt202: {
		StatusCode:    CodeErrInt202,
		StatusMessage: "DAILY_LIMIT_EXCEEDED",
		StatusDesc:    "The total value or volume of transactions has exceeded the agreed maximum daily limit.",
		UiMessage:     "Batas kuota transaksi harian Anda telah tercapai. Anda dapat melakukan transaksi kembali esok hari.",
	},
	CodeErrInt203: {
		StatusCode:    CodeErrInt203,
		StatusMessage: "INVALID_API_METHOD",
		StatusDesc:    "The API endpoint requested does not support the HTTP verb or request method applied.",
		UiMessage:     "Permintaan sistem tidak didukung. Pastikan aplikasi Anda terintegrasi dengan benar.",
	},
	CodeErrInt204: {
		StatusCode:    CodeErrInt204,
		StatusMessage: "MAINTENANCE_SCHEDULED",
		StatusDesc:    "The client-facing dashboard and payment processing gateway are offline for scheduled system upgrade.",
		UiMessage:     "Layanan kami sedang ditingkatkan untuk kenyamanan Anda. Kami akan segera kembali dalam beberapa saat.",
	},

	// Webpage - ERR Upstream Side
	CodeErrPvd1300: {
		StatusCode:    CodeErrPvd1300,
		StatusMessage: "UPSTREAM_MAINTENANCE",
		StatusDesc:    "The upstream vendor gateway for this specific product is currently undergoing maintenance.",
		UiMessage:     "Server pusat untuk game ini sedang mengalami pemeliharaan sistem. Silakan coba beberapa saat lagi.",
	},
	CodeErrPvd2301: {
		StatusCode:    CodeErrPvd2301,
		StatusMessage: "PRODUCT_OUT_OF_STOCK",
		StatusDesc:    "The upstream distributor core reported that the stock, voucher pool, or product quota is exhausted.",
		UiMessage:     "Stok untuk item game ini sedang habis di pusat. Tim kami sedang melakukan pengisian ulang stok.",
	},
	CodeErrPvd1302: {
		StatusCode:    CodeErrPvd1302,
		StatusMessage: "UPSTREAM_TIMEOUT",
		StatusDesc:    "Upstream server failed to respond within the designated execution window. Money safely refunded.",
		UiMessage:     "Jaringan ke server game sedang padat. Transaksi dibatalkan secara aman dan saldo Anda tetap utuh.",
	},
	CodeErrPvd3303: {
		StatusCode:    CodeErrPvd3303,
		StatusMessage: "UPSTREAM_DECLINED",
		StatusDesc:    "The execution request was rejected by the distribution core gateway due to structural rule conflicts.",
		UiMessage:     "Pembelian ditolak oleh sistem pusat game. Silakan pilih nominal atau paket game lainnya.",
	},
	CodeErrPvd2304: {
		StatusCode:    CodeErrPvd2304,
		StatusMessage: "UPSTREAM_PRICE_CHANGED",
		StatusDesc:    "The purchase cost from the upstream gateway has changed and exceeded the configured system margin.",
		UiMessage:     "Terjadi pembaruan harga dari sistem pusat game. Silakan ulangi transaksi Anda untuk memperbarui harga.",
	},
	CodeErrPvd4305: {
		StatusCode:    CodeErrPvd4305,
		StatusMessage: "UPSTREAM_UNKNOWN_ERROR",
		StatusDesc:    "Upstream provider responded with an unmapped error structure or critical raw payload exception.",
		UiMessage:     "Terjadi kendala pada jaringan distribusi game. Saldo tidak terpotong, silakan coba lagi nanti.",
	},

	// Dashboard - AUTH
	CodeSuccessAuth: {
		StatusCode:    CodeSuccessAuth,
		StatusMessage: "AUTHENTICATION_SUCCESS",
		StatusDesc:    "User credentials verified. Session token generated successfully.",
		UiMessage:     "Login berhasil! Mengalihkan Anda ke halaman dashboard...",
	},
	CodeSucAuth201: {
		StatusCode:    CodeSucAuth201,
		StatusMessage: "REGISTRATION_SUCCESS",
		StatusDesc:    "New user entity successfully created inside the core database.",
		UiMessage:     "Pendaftaran berhasil! Selamat bergabung di platform kami.",
	},
	CodeErrAuth401: {
		StatusCode:    CodeErrAuth401,
		StatusMessage: "INVALID_CREDENTIALS",
		StatusDesc:    "The password hash or username/email combination did not match database records.",
		UiMessage:     "Email atau password yang Anda masukkan salah. Silakan coba lagi.",
	},
	CodeErrAuth403: {
		StatusCode:    CodeErrAuth403,
		StatusMessage: "ACCESS_DENIED",
		StatusDesc:    "The authenticated user does not have the required role assigned in model_has_roles.",
		UiMessage:     "Akun Anda tidak memiliki izin untuk mengakses halaman ini.",
	},
	CodeErrAuth419: {
		StatusCode:    CodeErrAuth419,
		StatusMessage: "TOKEN_EXPIRED",
		StatusDesc:    "The bearer token or session cookie has expired or been blacklisted.",
		UiMessage:     "Sesi Anda telah berakhir demi keamanan. Silakan masuk (login) kembali.",
	},

	// Dashboard - USER / VAL
	CodeSucUser200: {
		StatusCode:    CodeSucUser200,
		StatusMessage: "PROFILE_UPDATED",
		StatusDesc:    "Profile fields successfully modified and committed to the user entity.",
		UiMessage:     "Perubahan profil Anda telah berhasil disimpan.",
	},
	CodeValUser422: {
		StatusCode:    CodeValUser422,
		StatusMessage: "EMAIL_ALREADY_EXISTS",
		StatusDesc:    "Unique constraint failed. The provided email address is already taken by another account.",
		UiMessage:     "Email sudah terdaftar. Silakan gunakan email lain atau gunakan fitur lupa password.",
	},
	CodeValUser423: {
		StatusCode:    CodeValUser423,
		StatusMessage: "PHONE_ALREADY_EXISTS",
		StatusDesc:    "Unique constraint failed. The phone number is already registered in the system.",
		UiMessage:     "Nomor HP sudah digunakan. Mohon periksa kembali nomor Anda.",
	},
	CodeValUser424: {
		StatusCode:    CodeValUser424,
		StatusMessage: "WEAK_PASSWORD",
		StatusDesc:    "Validation failed. The password policy string requirements are not met.",
		UiMessage:     "Password terlalu lemah. Gunakan minimal 8 karakter dengan kombinasi angka dan huruf.",
	},
	CodeErrUser404: {
		StatusCode:    CodeErrUser404,
		StatusMessage: "USER_NOT_FOUND",
		StatusDesc:    "Query returned no records for the specified User ID or identifier.",
		UiMessage:     "Data pengguna tidak ditemukan di dalam sistem kami.",
	},

	// Dashboard - MERCH
	CodeSucMerch200: {
		StatusCode:    CodeSucMerch200,
		StatusMessage: "API_CREDENTIALS_REGENERATED",
		StatusDesc:    "New client_key and secret_key successfully rolled over for the merchant entity.",
		UiMessage:     "API Key baru berhasil dibuat. Pastikan Anda segera memperbarui sistem H2H Anda.",
	},
	CodeSucMerch201: {
		StatusCode:    CodeSucMerch201,
		StatusMessage: "IP_WHITELIST_UPDATED",
		StatusDesc:    "Whitelist IP text area parsed and updated inside merchant_api_credentials.",
		UiMessage:     "Daftar IP Whitelist Anda telah berhasil diperbarui.",
	},
	CodeErrMerch400: {
		StatusCode:    CodeErrMerch400,
		StatusMessage: "INVALID_IP_FORMAT",
		StatusDesc:    "Failed to parse the provided whitelist string. IP address format is invalid.",
		UiMessage:     "Format IP Address yang Anda masukkan salah. Mohon periksa kembali.",
	},
	CodeErrMerch403: {
		StatusCode:    CodeErrMerch403,
		StatusMessage: "WRONG_TRANSACTION_PIN",
		StatusDesc:    "The secure transaction PIN validation failed against account_pin_hash.",
		UiMessage:     "PIN Transaksi yang Anda masukkan salah. Sisa kesempatan 2 kali lagi.",
	},

	// Dashboard - UTIL
	CodeSucUtil200: {
		StatusCode:    CodeSucUtil200,
		StatusMessage: "OTP_SENT",
		StatusDesc:    "Verification code generated and successfully pushed to the notification provider.",
		UiMessage:     "Kode verifikasi (OTP) telah dikirimkan ke nomor HP / email Anda.",
	},
	CodeSucUtil202: {
		StatusCode:    CodeSucUtil202,
		StatusMessage: "OTP_VERIFIED",
		StatusDesc:    "OTP token matches the active otp_codes row and is successfully marked as used.",
		UiMessage:     "Verifikasi berhasil. Silakan melanjutkan ke proses berikutnya.",
	},
	CodeErrUtil408: {
		StatusCode:    CodeErrUtil408,
		StatusMessage: "OTP_EXPIRED",
		StatusDesc:    "The system matched the code but the active window passed the expired_at timestamp.",
		UiMessage:     "Kode OTP sudah kedaluwarsa. Silakan klik tombol \"Kirim Ulang\".",
	},
	CodeErrUtil410: {
		StatusCode:    CodeErrUtil410,
		StatusMessage: "OTP_INCORRECT",
		StatusDesc:    "The input OTP string does not match the active sequence stored.",
		UiMessage:     "Kode OTP yang Anda masukkan salah. Mohon periksa kembali pesan masuk Anda.",
	},
	CodeErrUtil429: {
		StatusCode:    CodeErrUtil429,
		StatusMessage: "OTP_MAX_ATTEMPT_REACHED",
		StatusDesc:    "The counter on attempt_count has exceeded max_attempt. The token is invalidated.",
		UiMessage:     "Anda telah salah memasukkan OTP sebanyak 3 kali. Silakan minta kode OTP baru.",
	},

	// System Error
	CodeErrSys404: {
		StatusCode:    CodeErrSys404,
		StatusMessage: "NOT_FOUND",
		StatusDesc:    "Resource or reference not found.",
		UiMessage:     "Data yang Anda cari tidak ditemukan.",
	},
	CodeErrSys500: {
		StatusCode:    CodeErrSys500,
		StatusMessage: "SYSTEM_ERROR",
		StatusDesc:    "Internal system error occurred.",
		UiMessage:     "Terjadi gangguan sistem. Silakan coba beberapa saat lagi.",
	},
}

// BuildResponse generates a standardized models.ApiResponse
func BuildResponse(code string, result any) models.ApiResponse {
	meta, exists := responseCodes[code]
	if !exists {
		return models.ApiResponse{
			StatusCode:    CodeErrSys500,
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
}
