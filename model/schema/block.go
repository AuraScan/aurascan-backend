package schema

type PageListReq struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type BlockSpecReq struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Height   int64 `json:"height"`
}

type BlockListResp struct {
	Height         int64   `json:"height"`
	Epoch          int64   `json:"epoch"`       // height / 360
	EpochIndex     int64   `json:"epoch_index"` // height - epoch*360
	Round          int64   `json:"round"`
	Time           string  `json:"time"` // 2024-01-01 13:14:15
	ProofTarget    float64 `json:"proof_target"`
	CoinbaseTarget float64 `json:"coinbase_target"`
	BlockReward    float64 `json:"block_reward"`
	CoinbaseReward float64 `json:"coinbase_reward"`
	Solutions      int     `json:"solutions"`
	Transactions   int     `json:"transactions"`
}

type BlockDetailResp struct {
	Height                int64   `json:"height"`
	BlockHash             string  `json:"block_hash"`
	PreviousHash          string  `json:"previous_hash"`
	PreviousStateRoot     string  `json:"previous_state_root"`
	TransactionsRoot      string  `json:"transactions_root"`
	FinalizeRoot          string  `json:"finalize_root"`
	RatificationsRoot     string  `json:"ratifications_root"`
	CumulativeWeight      float64 `json:"cumulative_weight"`
	CumulativeProofTarget float64 `json:"cumulative_proof_target"`
	AuthorityType         string  `json:"authority_type"`

	Round          int64   `json:"round"`
	BlockReward    float64 `json:"block_reward"`
	CoinbaseReward float64 `json:"coinbase_reward"`
	ProofTarget    float64 `json:"proof_target"`
	CoinbaseTarget float64 `json:"coinbase_target"`
	Network        int     `json:"network"`
	Time           string  `json:"time"`
}

type BlockInRedis struct {
	BlockHeight           int64   `json:"block_height"`
	LatestBlockTime       int64   `json:"latest_block_time"`
	CoinbaseTarget        float64 `json:"coinbase_target"`
	ProofTarget           float64 `json:"proof_target"`
	EstimatedNetworkSpeed float64 `json:"estimated_network_speed"`
}
