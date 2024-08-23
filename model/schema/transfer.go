package schema

type TransferListReq struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Address  string `json:"address"`
}
