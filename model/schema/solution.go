package schema

type SolutionInBlockResp struct {
	Address    string  `json:"address"`
	Commitment string  `json:"commitment"`
	Target     float64 `json:"target"`
	Reward     float64 `json:"reward"`
}

type SolutionInAddrResp struct {
	BlockHeight int64   `json:"block_height"`
	SolutionId  string  `json:"solution_id"`
	Time        string  `json:"time"`
	Target      float64 `json:"target"`
	Reward      float64 `json:"reward"`
}

type SolutionInAddrReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Address  string `json:"address"`
}
