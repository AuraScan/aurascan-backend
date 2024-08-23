package model

// 创世区块
type BlockGenesis struct {
	BlockHash             string                `json:"block_hash"`
	PreviousHash          string                `json:"previous_hash"`
	Header                Header                `json:"header"`
	Authority             AuthorityGenesis      `json:"authority"`
	Ratifications         []RatificationGenesis `json:"ratifications"`
	Transactions          []TransactionSpec     `json:"transactions"`
	AbortedTransactionIds []string              `json:"aborted_transaction_ids"`
}

type AuthorityGenesis struct {
	Type      string `json:"type"`
	Signature string `json:"signature"`
}

type RatificationGenesis struct {
	Type           string             `json:"type"` //类型为genesis
	Committee      Committee          `json:"committee"`
	PublicBalances map[string]float64 `json:"public_balances"`
}

type Committee struct {
	StartingRound int64                    `json:"starting_round"`
	Members       map[string][]interface{} `json:"members"`
	TotalStake    float64                  `json:"total_stake"`
}
