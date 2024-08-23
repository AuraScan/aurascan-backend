package schema

type TransactionListResp struct {
	Id        string  `json:"id"`
	Height    int64   `json:"height"`
	Time      string  `json:"time"`
	Timestamp int64   `json:"timestamp"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Fee       float64 `json:"fee"`
}

type TransactionListInBlockResp struct {
	Id     string  `json:"id"`
	Type   string  `json:"type"`
	Status string  `json:"status"`
	Fee    float64 `json:"fee"`
}

type TransactionDetailResp struct {
	Id     string `json:"id"`
	Height int64  `json:"height"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Time   string `json:"time"`
	//GlobalStateRoot string                            `json:"global_state_root"`
	Fee         float64                            `json:"fee"`
	Finalize    interface{}                        `json:"finalize"`
	Transitions []*TransitionListInTransactionResp `json:"transitions"`
}
