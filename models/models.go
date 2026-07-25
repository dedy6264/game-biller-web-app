package models

import "github.com/golang-jwt/jwt/v5"

// JwtCustomClaims represents custom claims inside JWT
type JwtCustomClaims struct {
	UserID     int64  `json:"user_id"`
	MerchantID int64  `json:"merchant_id"`
	RoleCode   string `json:"role_code"`
	jwt.RegisteredClaims
}

// 1. CENTRAL AUTHENTICATION & SECURITY
type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	PasswordHash string `json:"-"`
	Status       string `json:"status"` // active, suspended, unverified
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
}

type Role struct {
	ID        int64  `json:"id"`
	RoleCode  string `json:"role_code"` // super_admin, finance, cs, merchant_h2h, member_reseller, retail_guest
	RoleName  string `json:"role_name"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy string `json:"updated_by"`
}

type ModelHasRole struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	RoleID    int64  `json:"role_id"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	RoleName  string `json:"role_name,omitempty"`
	RoleCode  string `json:"role_code,omitempty"`
}

type OtpCode struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Identifier   string `json:"identifier"` // Email atau nomor HP target OTP
	OtpCode      string `json:"otp_code"`
	OtpType      string `json:"otp_type"` // REGISTER, LOGIN, RESET_PIN
	ExpiredAt    string `json:"expired_at"`
	IsUsed       bool   `json:"is_used"`
	AttemptCount int    `json:"attempt_count"`
	MaxAttempt   int    `json:"max_attempt"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
}

// 2. MERCHANT PROFILE & CREDENTIALS
type Merchant struct {
	ID            int64                  `json:"id"`
	UserID        int64                  `json:"user_id"`
	SegmentID     *int64                 `json:"segment_id"`
	MerchantName  string                 `json:"merchant_name"`
	MerchantType  string                 `json:"merchant_type"` // guest_retail, member_premium, h2h_api
	Status        string                 `json:"status"`
	CreatedAt     string                 `json:"created_at"`
	CreatedBy     string                 `json:"created_by"`
	UpdatedAt     string                 `json:"updated_at"`
	UpdatedBy     string                 `json:"updated_by"`
	UserName      string                 `json:"user_name"`
	UserEmail     string                 `json:"user_email"`
	SegmentName   string                 `json:"segment_name"`
	ClientKey     string                 `json:"client_key,omitempty"`
	SecretKey     string                 `json:"secret_key,omitempty"`
	WhitelistIPs  string                 `json:"whitelist_ips,omitempty"`
	ApiIsActive   bool                   `json:"api_is_active"`
	ApiCredential *MerchantApiCredential `json:"api_credential"`
}

type MerchantApiCredential struct {
	ID            int64  `json:"id"`
	MerchantID    int64  `json:"merchant_id"`
	ClientKey     string `json:"client_key"`
	SecretKey     string `json:"secret_key,omitempty"`
	SecretKeyHash string `json:"-"`
	WhitelistIPs  string `json:"whitelist_ips"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
}

// 3. BALANCE & MUTATION LOGS
type SavingAccount struct {
	ID             int64   `json:"id"`
	MerchantID     int64   `json:"merchant_id"`
	AccountNumber  string  `json:"account_number"`
	Balance        float64 `json:"balance"`
	AccountPinHash string  `json:"-"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	CreatedBy      string  `json:"created_by"`
	UpdatedAt      string  `json:"updated_at"`
	UpdatedBy      string  `json:"updated_by"`
}

type SavingTransaction struct {
	ID              int64   `json:"id"`
	SavingAccountID int64   `json:"saving_account_id"`
	TypeDC          string  `json:"type_dc"` // D = Debit, C = Credit
	Amount          float64 `json:"amount"`
	LastBalance     float64 `json:"last_balance"`
	ReferenceNumber string  `json:"reference_number"`
	TransactionCode string  `json:"transaction_code"` // DEPOSIT, GAME_TOPUP, REFUND, ADJUSTMENT
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
	CreatedBy       string  `json:"created_by"`
	CreatedByUser   *int64  `json:"created_by_user,omitempty"`
}

// 4. PRODUCT & PROVIDER MASTER
type Provider struct {
	ID           int64  `json:"id"`
	ProviderName string `json:"provider_name"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
}

type ProductType struct {
	ID              int64  `json:"id"`
	ProductTypeName string `json:"product_type_name"`
}

type ProductCategory struct {
	ID                  int64  `json:"id"`
	ProductCategoryName string `json:"product_category_name"`
}

type ProductReference struct {
	ID                   int64  `json:"id"`
	ProductReferenceCode string `json:"product_reference_code"`
	ProductReferenceName string `json:"product_reference_name"`
	CreatedAt            string `json:"created_at"`
	CreatedBy            string `json:"created_by"`
	UpdatedAt            string `json:"updated_at"`
	UpdatedBy            string `json:"updated_by"`
}

type ProductPrefix struct {
	ID                   int64  `json:"id"`
	ProductReferenceID   int64  `json:"product_reference_id"`
	PrefixNumber         string `json:"prefix_number"`
	CreatedAt            string `json:"created_at"`
	CreatedBy            string `json:"created_by"`
	UpdatedAt            string `json:"updated_at"`
	UpdatedBy            string `json:"updated_by"`
	ProductReferenceName string `json:"product_reference_name,omitempty"`
}

type Product struct {
	ID                   int64  `json:"id"`
	ProductReferenceID   *int64 `json:"product_reference_id,omitempty"`
	ProductTypeID        int64  `json:"product_type_id"`
	ProductCategoryID    int64  `json:"product_category_id"`
	ProductCode          string `json:"product_code"`
	ProductName          string `json:"product_name"`
	IsActive             bool   `json:"is_active"`
	CreatedAt            string `json:"created_at"`
	CreatedBy            string `json:"created_by"`
	UpdatedAt            string `json:"updated_at"`
	UpdatedBy            string `json:"updated_by"`
	ProductReferenceName string `json:"product_reference_name,omitempty"`
	ProductCategoryName  string `json:"product_category_name,omitempty"`
}

type ProductProvider struct {
	ID                         int64   `json:"id"`
	ProviderID                 int64   `json:"provider_id"`
	ProviderName               string  `json:"provider_name"`
	ProductProviderCode        string  `json:"product_provider_code"`
	ProductProviderName        string  `json:"product_provider_name"`
	ProductProviderPrice       float64 `json:"product_provider_price"`
	ProductProviderAdminFee    float64 `json:"product_provider_admin_fee"`
	ProductProviderMerchantFee float64 `json:"product_provider_merchant_fee"`
	ProviderIndex              int     `json:"provider_index"`
	IsAvailable                bool    `json:"is_available"`
	CreatedAt                  string  `json:"created_at"`
	CreatedBy                  string  `json:"created_by"`
	UpdatedAt                  string  `json:"updated_at"`
	UpdatedBy                  string  `json:"updated_by"`
}

type ProductSegment struct {
	ID                         int64   `json:"id"`
	SegmentID                  *int64  `json:"segment_id,omitempty"`
	ProductProviderID          *int64  `json:"product_provider_id,omitempty"`
	SegmentName                string  `json:"segment_name"` // Public_Retail, Gold_Reseller, H2H_Partner
	ProductID                  int64   `json:"product_id"`
	ProductName                string  `json:"product_name"`
	ProviderName               string  `json:"provider_name"`
	ProductPrice               float64 `json:"product_price"`
	AdminFee                   float64 `json:"admin_fee"`
	MerchantFee                float64 `json:"merchant_fee"`
	ProductProviderCode        string  `json:"product_provider_code,omitempty"`
	ProductProviderName        string  `json:"product_provider_name,omitempty"`
	ProductProviderPrice       float64 `json:"product_provider_price"`
	ProductProviderAdminFee    float64 `json:"product_provider_admin_fee"`
	ProductProviderMerchantFee float64 `json:"product_provider_merchant_fee"`
	CreatedAt                  string  `json:"created_at"`
	CreatedBy                  string  `json:"created_by"`
	UpdatedAt                  string  `json:"updated_at"`
	UpdatedBy                  string  `json:"updated_by"`
	ProductCode                string  `json:"product_code,omitempty"`
}

// 5. PAYMENT METHOD GATEWAY
type PaymentMethod struct {
	ID         int64  `json:"id"`
	MethodCode string `json:"method_code"` // VIRTUAL_ACCOUNT, E_WALLET, QRIS, OVER_THE_COUNTER, DEPOSIT
	MethodName string `json:"method_name"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	UpdatedAt  string `json:"updated_at"`
	UpdatedBy  string `json:"updated_by"`
}

type PaymentChannel struct {
	ID              int64   `json:"id"`
	PaymentMethodID int64   `json:"payment_method_id"`
	ChannelCode     string  `json:"channel_code"` // BCA_VA, MANDIRI_VA, OVO, DANA, QRIS_INTERNASIONAL, ALFAMART, BALANCE_INTERNAL
	ChannelName     string  `json:"channel_name"`
	FeeType         string  `json:"fee_type"` // FIXED, PERCENTAGE
	FeeValue        float64 `json:"fee_value"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	CreatedBy       string  `json:"created_by"`
	UpdatedAt       string  `json:"updated_at"`
	UpdatedBy       string  `json:"updated_by"`
	MethodName      string  `json:"method_name,omitempty"`
	MethodCode      string  `json:"method_code,omitempty"`
}

// 6. CORE TRANSACTION
type Transaction struct {
	ID                         int64   `json:"id"`
	MerchantID                 int64   `json:"merchant_id"`
	ProductID                  int64   `json:"product_id"`
	ProductSegmentID           int64   `json:"product_segment_id"`
	ProductProviderID          int64   `json:"product_provider_id"`
	ProviderID                 int64   `json:"provider_id"`
	ProductTypeID              int64   `json:"product_type_id"`
	ProductReferenceID         int64   `json:"product_reference_id"`
	PaymentChannelID           int64   `json:"payment_channel_id"`
	ProductCode                string  `json:"product_code"`
	SnapshotProductCode        string  `json:"snapshot_product_code"`
	SnapshotProductName        string  `json:"snapshot_product_name"`
	MerchantName               string  `json:"merchant_name"`
	ProductName                string  `json:"product_name"`
	ProductSegmentName         string  `json:"product_segment_name"`
	ProductProviderCode        string  `json:"product_provider_code"`
	ProductProviderName        string  `json:"product_provider_name"`
	ProviderName               string  `json:"provider_name"`
	ProductTypeName            string  `json:"product_type_name"`
	PaymentChannelName         string  `json:"payment_channel_name"`
	ProductProviderPrice       float64 `json:"product_provider_price"`
	ProductPrice               float64 `json:"product_price"`
	ProductAdminFee            float64 `json:"product_admin_fee"`
	ProductMerchantFee         float64 `json:"product_merchant_fee"`
	ProductProviderAdminFee    float64 `json:"product_provider_admin_fee"`
	ProductProviderMerchantFee float64 `json:"product_provider_merchant_fee"`
	PaymentAdminFee            float64 `json:"payment_admin_fee"`
	TotalAmount                float64 `json:"total_amount"`
	CustomerID                 string  `json:"customer_id"`
	OtherCustomerID            string  `json:"other_customer_id"`
	ReferenceNumberInternal    string  `json:"reference_number_internal"`
	ReferenceNumberMerchant    string  `json:"reference_number_merchant"`
	ReferenceNumberProvider    string  `json:"reference_number_provider"`
	SerialNumber               string  `json:"serial_number"`
	StatusCode                 string  `json:"status_code"`
	StatusMessage              string  `json:"status_message"`
	RetryCount                 int     `json:"retry_count"`
	CreatedAt                  string  `json:"created_at"`
	CreatedBy                  string  `json:"created_by"`
	UpdatedAt                  string  `json:"updated_at"`
	UpdatedBy                  string  `json:"updated_by"`
}

type TransactionPayloadLog struct {
	ID              int64  `json:"id"`
	TransactionID   int64  `json:"transaction_id"`
	RequestPayload  string `json:"request_payload"`
	ResponsePayload string `json:"response_payload"`
	CreatedAt       string `json:"created_at"`
	CreatedBy       string `json:"created_by"`
}

// WebApp Payloads
type RegisterRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Can be Email or Phone Number
	Password string `json:"password" validate:"required"`
}

type InquiryRequest struct {
	ProductCode             string `json:"product_code" validate:"required"`
	TargetUserID            string `json:"target_user_id" validate:"required"`
	PaymentChannelCode      string `json:"payment_channel_code" validate:"required"`
	ReferenceNumberMerchant string `json:"reference_number_merchant"`
	ZoneID                  string `json:"zone_id"`
	ServerID                string `json:"server_id"`
	// OtherCustomerID         string `json:"other_customer_id"`
}
type OtherCustomerID struct {
	ZoneID   string `json:"zone_id"`
	ServerID string `json:"server_id"`
}
type PaymentRequest struct {
	ReferenceNumberInternal string `json:"reference_number_internal" validate:"required"`
	PIN                     string `json:"pin"` // optional, required for DEPOSIT method
}

// JSON Generic Response Wrapper as per assignment.md
type ApiResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	StatusDesc    string `json:"status_desc"`
	UiMessage     string `json:"ui_message"`
	Result        any    `json:"result"`
}

type RequestUsers struct {
	Draw    int         `json:"draw"`
	Search  string      `json:"search"`
	Start   int         `json:"start"`
	Length  int         `json:"length"`
	Order   string      `json:"order"`
	Sort    string      `json:"sort"`
	Filters UserFilters `json:"filters"`
}

type UserFilters struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
}

type RequestRoles struct {
	Draw    int         `json:"draw"`
	Search  string      `json:"search"`
	Start   int         `json:"start"`
	Length  int         `json:"length"`
	Order   string      `json:"order"`
	Sort    string      `json:"sort"`
	Filters RoleFilters `json:"filters"`
}

type RoleFilters struct {
	ID       int64  `json:"id"`
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

type RequestModelHasRoles struct {
	Draw    int                 `json:"draw"`
	Search  string              `json:"search"`
	Start   int                 `json:"start"`
	Length  int                 `json:"length"`
	Order   string              `json:"order"`
	Sort    string              `json:"sort"`
	Filters ModelHasRoleFilters `json:"filters"`
}

type ModelHasRoleFilters struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type RequestProviders struct {
	Draw    int             `json:"draw"`
	Search  string          `json:"search"`
	Start   int             `json:"start"`
	Length  int             `json:"length"`
	Order   string          `json:"order"`
	Sort    string          `json:"sort"`
	Filters ProviderFilters `json:"filters"`
}

type ProviderFilters struct {
	ID           int64  `json:"id"`
	ProviderName string `json:"provider_name"`
	IsActive     bool   `json:"is_active"`
}

type RequestProductCategories struct {
	Draw    int                    `json:"draw"`
	Search  string                 `json:"search"`
	Start   int                    `json:"start"`
	Length  int                    `json:"length"`
	Order   string                 `json:"order"`
	Sort    string                 `json:"sort"`
	Filters ProductCategoryFilters `json:"filters"`
}

type ProductCategoryFilters struct {
	ID                  int64  `json:"id"`
	ProductCategoryName string `json:"product_category_name"`
}

type RequestProductReferences struct {
	Draw    int                     `json:"draw"`
	Search  string                  `json:"search"`
	Start   int                     `json:"start"`
	Length  int                     `json:"length"`
	Order   string                  `json:"order"`
	Sort    string                  `json:"sort"`
	Filters ProductReferenceFilters `json:"filters"`
}

type ProductReferenceFilters struct {
	ID                   int64  `json:"id"`
	ProductReferenceCode string `json:"product_reference_code"`
	ProductReferenceName string `json:"product_reference_name"`
}

type RequestProductPrefixes struct {
	Draw    int                  `json:"draw"`
	Search  string               `json:"search"`
	Start   int                  `json:"start"`
	Length  int                  `json:"length"`
	Order   string               `json:"order"`
	Sort    string               `json:"sort"`
	Filters ProductPrefixFilters `json:"filters"`
}

type ProductPrefixFilters struct {
	ID                 int64  `json:"id"`
	ProductReferenceID int64  `json:"product_reference_id"`
	PrefixNumber       string `json:"prefix_number"`
}

type RequestProducts struct {
	Draw    int            `json:"draw"`
	Search  string         `json:"search"`
	Start   int            `json:"start"`
	Length  int            `json:"length"`
	Order   string         `json:"order"`
	Sort    string         `json:"sort"`
	Filters ProductFilters `json:"filters"`
}

type ProductFilters struct {
	ID                 int64  `json:"id"`
	ProductReferenceID *int64 `json:"product_reference_id"`
	ProductTypeID      int64  `json:"product_type_id"`
	ProductCategoryID  int64  `json:"product_category_id"`
	ProductCode        string `json:"product_code"`
	ProductName        string `json:"product_name"`
}

type RequestProductProviders struct {
	Draw    int                    `json:"draw"`
	Search  string                 `json:"search"`
	Start   int                    `json:"start"`
	Length  int                    `json:"length"`
	Order   string                 `json:"order"`
	Sort    string                 `json:"sort"`
	Filters ProductProviderFilters `json:"filters"`
}

type ProductProviderFilters struct {
	ID                  int64  `json:"id"`
	ProviderID          int64  `json:"provider_id"`
	ProductProviderCode string `json:"product_provider_code"`
	ProductProviderName string `json:"product_provider_name"`
}

type RequestProductSegments struct {
	Draw    int                   `json:"draw"`
	Search  string                `json:"search"`
	Start   int                   `json:"start"`
	Length  int                   `json:"length"`
	Order   string                `json:"order"`
	Sort    string                `json:"sort"`
	Filters ProductSegmentFilters `json:"filters"`
}

type ProductSegmentFilters struct {
	ID                int64  `json:"id"`
	SegmentID         *int64 `json:"segment_id"`
	ProductProviderID *int64 `json:"product_provider_id"`
	ProductID         int64  `json:"product_id"`
	SegmentName       string `json:"segment_name"`
}

type RequestPaymentMethods struct {
	Draw    int                  `json:"draw"`
	Search  string               `json:"search"`
	Start   int                  `json:"start"`
	Length  int                  `json:"length"`
	Order   string               `json:"order"`
	Sort    string               `json:"sort"`
	Filters PaymentMethodFilters `json:"filters"`
}

type PaymentMethodFilters struct {
	ID         int64  `json:"id"`
	MethodCode string `json:"method_code"`
	MethodName string `json:"method_name"`
}

type RequestPaymentChannels struct {
	Draw    int                   `json:"draw"`
	Search  string                `json:"search"`
	Start   int                   `json:"start"`
	Length  int                   `json:"length"`
	Order   string                `json:"order"`
	Sort    string                `json:"sort"`
	Filters PaymentChannelFilters `json:"filters"`
}

type PaymentChannelFilters struct {
	ID              int64  `json:"id"`
	PaymentMethodID int64  `json:"payment_method_id"`
	ChannelCode     string `json:"channel_code"`
	ChannelName     string `json:"channel_name"`
}

type RequestTransactions struct {
	Draw    int                `json:"draw"`
	Search  string             `json:"search"`
	Start   int                `json:"start"`
	Length  int                `json:"length"`
	Order   string             `json:"order"`
	Sort    string             `json:"sort"`
	Filters TransactionFilters `json:"filters"`
}

type TransactionFilters struct {
	ID                      int64   `json:"id"`
	MerchantID              int64   `json:"merchant_id"`
	ProductID               *int64  `json:"product_id"`
	ProductSegmentID        *int64  `json:"product_segment_id"`
	ProductProviderID       *int64  `json:"product_provider_id"`
	ProductProviderCode     string  `json:"product_provider_code"`
	ProductProviderName     string  `json:"product_provider_name"`
	ProviderID              *int64  `json:"provider_id"`
	ProductTypeID           *int64  `json:"product_type_id"`
	ProductReferenceID      *int64  `json:"product_reference_id"`
	PaymentChannelID        *int64  `json:"payment_channel_id"`
	ProductCode             string  `json:"product_code"`
	CustomerID              string  `json:"customer_id"`
	OtherCustomerID         string  `json:"other_customer_id"`
	StatusCode              string  `json:"status_code"`
	ReferenceNumberInternal string  `json:"reference_number_internal"`
	ReferenceNumberMerchant *string `json:"reference_number_merchant"`
	StartDate               string  `json:"start_date"`
	EndDate                 string  `json:"end_date"`
}

type Segment struct {
	ID          int64  `json:"id"`
	SegmentName string `json:"segment_name"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
}

type RequestSegments struct {
	Draw    int            `json:"draw"`
	Search  string         `json:"search"`
	Start   int            `json:"start"`
	Length  int            `json:"length"`
	Order   string         `json:"order"`
	Sort    string         `json:"sort"`
	Filters SegmentFilters `json:"filters"`
}

type SegmentFilters struct {
	ID          int64  `json:"id"`
	SegmentName string `json:"segment_name"`
}

type ReqHeader struct {
	Header []HeaderItem `json:"header"`
}

type HeaderItem struct {
	Key      string `json:"key"`
	Val      string `json:"val"`
	IsUpCase bool   `json:"is_up_case"`
}
