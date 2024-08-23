package schema

type Network struct {
	BlockHeight     int64   `json:"block_height"`      //最新出块高度
	LatestBlockTime int64   `json:"latest_block_time"` //最近出块时间
	ProofTarget     float64 `json:"proof_target"`      //Proof目标
	CoinbaseTarget  float64 `json:"coinbase_target"`   //Coinbase目标
	Epoch           int64   `json:"epoch"`             //当前Epoch

	NetworkStaking           float64 `json:"network_staking"`             //全网Staking总数
	EstimatedNetworkSpeed    float64 `json:"estimated_network_speed"`     //15m全网算力
	NetworkValidators        int     `json:"network_validators"`          // 全网Validator数量
	NetworkDelegators        int     `json:"network_delegators"`          // 全网Delegator数量
	NetworkMiners            int     `json:"network_miners"`              //全网矿工数
	ProgramCount             int     `json:"program_count"`               //Program总数
	NetworkPuzzleReward      float64 `json:"network_reward"`              //全网总出块收益
	NetworkEffectiveProof24h float64 `json:"network_effective_proof_24h"` //全网24小时有效证明提交数
}

type ProofTargetChart struct {
	Timestamp   int64   `json:"timestamp"`
	Height      int64   `json:"height"`
	ProofTarget float64 `json:"proof_target"`
}
