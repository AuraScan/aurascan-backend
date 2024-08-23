package schema

type StakingDailyResp struct {
	Date          string  `json:"date"`
	TotalBond     float64 `json:"total_bond"`
	ValidatorBond float64 `json:"validator_bond"`
	DelegatorBond float64 `json:"delegator_bond"`
}
