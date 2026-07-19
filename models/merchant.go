package models

type RequestMerchants struct {
	Draw    int             `json:"draw"`
	Search  string          `json:"search"`
	Start   int             `json:"start"`
	Length  int             `json:"length"`
	Order   string          `json:"order"`
	Sort    string          `json:"sort"`
	Filters MerchantFilters `json:"filters"`
}

type MerchantFilters struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	SegmentID    *int64 `json:"segment_id"`
	MerchantName string `json:"merchant_name"`
	MerchantType string `json:"merchant_type"`
	Status       string `json:"status"`
}
