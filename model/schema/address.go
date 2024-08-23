package schema

type ProverListResp struct {
	Rank      int     `json:"rank"`
	Address   string  `json:"address"`
	LastBlock int64   `json:"last_block"`
	Power     float64 `json:"power"`
	Reward    float64 `json:"reward"`
}

type AddrBondPart struct {
	Address   string          `bson:"addr"`
	BondState BondedStateResp `bson:"bond_state"`
}

type BondedStateResp struct {
	Validator  string  `json:"validator" bson:"validator"`
	BondAmount float64 `json:"bond_amount" bson:"bond_amount"`
}

type AddrDetailResp struct {
	Addr                string  `json:"addr"`
	AddrType            string  `json:"addr_type"`
	PublicCredits       float64 `json:"public_credits"`
	TotalPuzzleReward   float64 `json:"total_puzzle_reward"`
	TotalSolutionsFound int64   `json:"total_solutions_found"`
	Power1h             float64 `json:"power_1h"`
	Power24h            float64 `json:"power_24h"`
	Power7d             float64 `json:"power_7d"`

	BondValidatorState *ValidatorPartResp `json:"validator_state"`

	//validator信息
	CommitteeCreditsStake float64        `json:"committee_credits"`
	BondCreditsStake      float64        `json:"bond_credits"`
	DelegatorStake        float64        `json:"delegator_stake"`
	TotalCommitteeEarned  float64        `json:"total_committee_earned"`
	TotalValidatorEarned  float64        `json:"total_validator_earned"`
	TotalDelegatorEarned  float64        `json:"total_delegator_earned"`
	IsOpen                int            `json:"is_open"`
	Ratio                 float64        `json:"ratio"`
	TotalDelegators       int            `json:"delegators"`
	BondDelegatorList     []AddrBondPart `json:"delegator_list"`
}

type ValidatorPartResp struct {
	Validator string  `json:"validator"`
	Staked    float64 `json:"staked"`
	Earned    float64 `json:"earned"`
}
