package schema

type ProgramCountMonthChart struct {
	Program string               `json:"program"`
	Times   []*ProgramCalledTime `json:"times"`
}

type ProgramCalledTime struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}
