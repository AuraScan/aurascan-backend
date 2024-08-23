package model

type AddressSumDailyInDb struct {
	Address   string  `bson:"address"`
	Date      string  `bson:"date"`
	Power     float64 `bson:"power"`
	Reward    float64 `bson:"reward"`
	Solutions int     `bson:"solutions"`
	Timestamp int64   `bson:"timestamp"`
	CreatedAt int64   `bson:"created_at"`
}

func (*AddressSumDailyInDb) TableName() string {
	return "address_sum_daily"
}
