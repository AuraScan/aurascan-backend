package schema

type PowerChart struct {
	Date      string  `json:"date"`
	Timestamp int64   `json:"timestamp"`
	Power     float64 `json:"power"`
}

type RewardChart struct {
	Date   string  `json:"date"`
	Reward float64 `json:"reward"`
}

type RewardTimestampChart struct {
	Timestamp int64   `json:"timestamp"`
	Reward    float64 `json:"reward"`
}

type PowerTimestampChart struct {
	Timestamp int64   `json:"timestamp"`
	Power     float64 `json:"power"`
}

type SolutionsTimestampChart struct {
	Timestamp int64 `json:"timestamp"`
	Count     int   `json:"count"`
}

type PowerAndRewardChart struct {
	PowerCharts  []*PowerChart  `json:"powers"`
	RewardCharts []*RewardChart `json:"rewards"`
	TimeList     []string       `json:"time_list"`
}

type AddrRewardChart struct {
	Address      string                  `json:"address"`
	RewardCharts []*RewardTimestampChart `json:"rewards"`
}

type AddrPowerChart struct {
	Address     string                 `json:"address"`
	PowerCharts []*PowerTimestampChart `json:"powers"`
}

type AddrSolutionsChart struct {
	Address         string                     `json:"address"`
	SolutionsCharts []*SolutionsTimestampChart `json:"solutions"`
}

type PriceChart struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

type FeeChart struct {
	Date string  `json:"date"`
	Fee  float64 `json:"fee"`
}
