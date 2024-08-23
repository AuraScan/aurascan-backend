package schema

type TransitionListResp struct {
	ID        string `json:"id"`
	Height    int64  `json:"height"`
	Timestamp int64  `json:"timestamp"`
	Time      string `json:"time"`
	Program   string `json:"program"`
	Function  string `json:"function"`
}

type TransitionPageResp struct {
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	ProgramId string `json:"program_id,omitempty"`
}

type TransitionListInTransactionResp struct {
	ID       string `json:"id"`
	Program  string `json:"program"`
	Function string `json:"function"`
	State    string `json:"state"`
}

type TransitionDetailResp struct {
	TransitionId  string      `json:"transition_id"`
	TransactionId string      `json:"transaction_id"`
	State         string      `json:"state"`
	Program       string      `json:"program"`
	Function      string      `json:"function"`
	Tpk           string      `json:"tpk"`
	Tcm           string      `json:"tcm"`
	Input         interface{} `json:"input"`
	Output        interface{} `json:"output"`
}

type ProgramCalledChartResp struct {
	Program   string `json:"program"`
	CallTimes int    `json:"call_times"`
}
