package models

type (
	RequestInquiry struct {
		RefID       string      `json:"ref_id"` // Reference ID dari Main Service
		CustomerID  string      `json:"customer_id"`
		DataProduct DataProduct `json:"data_product"`
	}
	RequestPayment struct {
		RefID           string      `json:"ref_id"`          // Reference ID dari Main Service
		ProviderRefID   string      `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
		DataProduct     DataProduct `json:"data_product"`
		CustomerID      string      `json:"customer_id"`
		OtherCustomerID string      `json:"other_customer_id"`
		BillDesc        string      `json:"bill_desc"` //kumpulan informasi terkait product format baku
	}
	DataProduct struct {
		ProviderID         int64  `json:"provider_id"`
		ProductCategoryID  int64  `json:"product_category_id"`
		ProductTypeID      int64  `json:"product_type_id"`
		ProductCode        string `json:"product_code"`
		ProductReferenceID int64  `json:"product_reference_id"`
	}
)
type (
	PaymentResult struct {
		StatusCode      string           `json:"status_code"`
		ProviderRefID   string           `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
		RefID           string           `json:"ref_id"`          // Reference ID dari Main Service
		DataTransaction DataTransaction  `json:"data_transaction"`
		ProviderDetail  ProviderFeedback `json:"provider_detail"` //deskripsi status dari provider
		BillDesc        string           `json:"bill_desc"`       //kumpulan informasi terkait product format baku
		ProcessedAt     string           `json:"processed_at"`
	}
	DataTransaction struct {
		CustomerID   string  `json:"customer_id"`
		SerialNumber string  `json:"serial_number"` // SN / Token / Kode Voucher (Penting untuk produk game/pulsa)
		Price        float64 `json:"price"`         //nominal tagihan
		AdminFee     float64 `json:"AdminFee"`      //nominal admin dari provider
		MerchantFee  float64 `json:"MerchantFee"`   //nominal fee dari provider ke kita
		GrandTotal   float64 `json:"grand_total"`   //total pembayaran yang perlu dibayar konsumen
		LastBalance  float64 `json:"last_balance"`  // Sisa saldo deposit di provider (jika ada)
	}
	// ProviderFeedback menampung data mentah dari provider tanpa merusak format standar
	ProviderFeedback struct {
		Code        string `json:"code"`        // RC Asli Provider (misal: "00", "141", "39")
		Message     string `json:"message"`     // Message Murni Provider
		Description string `json:"description"` // Description Murni Provider
	}
)
type (
	InquiryResult struct {
		StatusCode      string           `json:"status_code"`
		IsInquiry       string           `json:"is_inquiry"`
		RefID           string           `json:"ref_id"`          // Reference ID dari Main Service
		ProviderRefID   string           `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
		DataTransaction DataTransaction  `json:"data_transaction"`
		ProviderDetail  ProviderFeedback `json:"provider_detail"` //deskripsi status dari provider
		BillDesc        string           `json:"bill_desc"`       //kumpulan informasi terkait product format baku
		ProcessedAt     string           `json:"processed_at"`
	}
)

type IAKCallbackData struct {
	RefID       string `json:"ref_id"`
	Status      string `json:"status"`
	ProductCode string `json:"product_code"`
	CustomerID  string `json:"customer_id"`
	Price       string `json:"price"`
	Message     string `json:"message"`
	Balance     string `json:"balance"`
	TrID        string `json:"tr_id"`
	RC          string `json:"rc"`
	SN          string `json:"sn"`
	PIN         string `json:"pin"`
	Sign        string `json:"sign"`
}

type IAKCallbackPayload struct {
	Data IAKCallbackData `json:"data"`
}
