package model

type StakingDailyInDb struct {
	Date          string  `bson:"date"`
	TotalBond     float64 `bson:"total_bond"`
	ValidatorBond float64 `bson:"validator_bond"`
	DelegatorBond float64 `bson:"delegator_bond"`
	Timestamp     int64   `bson:"timestamp"`
	CreatedAt     int64   `bson:"created_at"`
}

func (*StakingDailyInDb) TableName() string {
	return "bond_sum_daily"
}
