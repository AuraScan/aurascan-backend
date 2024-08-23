package schema

type ProgramListResp struct {
	ProgramId     string `json:"program_id"`
	TransactionId string `json:"transaction_id"`
	Height        int64  `json:"height"`
	Time          string `json:"time"`
	Timestamp     int64  `json:"timestamp"`
	TimesCalled   int    `json:"times_called"`
}

type ProgramDetailResp struct {
	ProgramId         string  `json:"program_id"`
	Owner             string  `json:"owner"`
	DeployHeight      int64   `json:"deploy_height"`
	DeployTime        string  `json:"deploy_time"`
	DeployTimestamp   int64   `json:"deploy_timestamp"`
	DeployTransaction string  `json:"deploy_transaction"`
	DeployFee         float64 `json:"deploy_fee"`
	TimesCalled       int     `json:"times_called"`
}

type ProgramCallingCount struct {
	Timestamp int64 `json:"timestamp"`
	Value     int   `json:"value"`
}

type MappingInfo struct {
	ProgramId   string `json:"program_id"`
	MappingName string `json:"mapping_name"`
	MappingKey  string `json:"mapping_key"`
}

type ProgramSource struct {
	SourceCode      string            `json:"source_code"`
	ProgramFunction []ProgramFunction `json:"functions"`
}

type ProgramFunction struct {
	Name   string   `json:"name"`
	Inputs []string `json:"inputs"`
}
