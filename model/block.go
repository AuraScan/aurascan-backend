package model

import (
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
)

type BlockInDb struct {
	BlockHash    string `bson:"bh"`
	PreviousHash string `bson:"ph"`

	Epoch      int64 `bson:"eh"`
	EpochIndex int64 `bson:"ei"`

	//CoinbaseReward float64 `bson:"coinbase_reward"`
	BlockReward  float64 `bson:"br"`
	PuzzleReward float64 `bson:"pr"`

	//Header
	PreviousStateRoot string `bson:"psr"`
	TransactionsRoot  string `bson:"tr"`
	FinalizeRoot      string `bson:"fr"`
	RatificationsRoot string `bson:"rr"`
	//CoinbaseAccumulatorPoint string `bson:"coinbase_accumulator_point"`
	SolutionsRoot string `bson:"sr"`
	SubdagRoot    string `bson:"sgr"`

	//MetaData
	Network int   `bson:"nk"`
	Round   int64 `bson:"rd"`
	Height  int64 `bson:"ht"`
	//TotalSupplyInMicrocredits float64 `bson:"total_supply_in_microcredits"`
	CumulativeWeight      float64 `bson:"cw"`
	CumulativeProofTarget float64 `bson:"cpt"`
	CoinbaseTarget        float64 `bson:"ct"`
	//TargetReached         float64 `bson:"target_reached"` //保留六位小数
	ProofTarget           float64 `bson:"pt"`
	LastCoinbaseTarget    int64   `bson:"lct"`
	LastCoinbaseTimestamp int64   `bson:"ltp"`
	Timestamp             int64   `bson:"tp"`

	//Authority     Authority         `json:"authority"`
	//Ratifications []Ratification    `json:"ratifications"`
	//Solutions     Solutions         `json:"solutions"`
	//Transactions  []TransactionSpec `json:"transactions"`
	AuthorityType string `bson:"at"`
	//Subdag SubdagOut `json:"subdag"`

	AbortedTransactionIds []string `json:"ats"`
	//X        string `bson:"x"`
	//Y        string `bson:"y"`
	//Infinity bool   `bson:"infinity"`

	//Signature string `bson:"signature"`

	//analysis
	SolutionNum    int     `bson:"sn"`
	TransactionNum int     `bson:"tn"`
	TotalFee       float64 `bson:"tf"`
	BaseFee        float64 `bson:"bf"`
	PriorityFee    float64 `bson:"pf"`
}

func (*BlockInDb) TableName() string {
	return "block"
}

func (b *BlockInDb) Save() error {
	if err := mongodb.InsertOne(context.TODO(), b.TableName(), b); err != nil {
		logger.Errorf("BlockInDb.Save InsertOne | %v", err)
		return err
	}
	return nil
}

func SaveBlockList(bs []*BlockInDb) error {
	if len(bs) == 0 {
		return nil
	}

	var blocks []interface{}
	for _, m := range bs {
		blocks = append(blocks, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&BlockInDb{}).TableName(), blocks); err != nil {
		return fmt.Errorf("BlockInDb.SaveBlockList InsertMany | %v", err)
	}
	return nil
}

type Block struct {
	BlockHash             string            `json:"block_hash"`
	PreviousHash          string            `json:"previous_hash"`
	Header                Header            `json:"header"`
	Authority             Authority         `json:"authority"`
	Ratifications         []Ratification    `json:"ratifications"`
	Solutions             Solutions         `json:"solutions"`
	AbortedSolutionIds    []string          `json:"aborted_solution_ids"`
	Transactions          []TransactionSpec `json:"transactions"`
	AbortedTransactionIds []string          `json:"aborted_transaction_ids"`
}

type Solutions struct {
	Version   int         `json:"version"`
	Solutions SolutionsIn `json:"solutions"`
}

type SolutionsIn struct {
	Solutions []Solution `json:"solutions"`
}

type Solution struct {
	PartialSolution PartialSolution `json:"partial_solution"`
	Target          float64         `json:"target"`
}

type PartialSolution struct {
	SolutionId string  `json:"solution_id"`
	EpochHash  string  `json:"epoch_hash"`
	Address    string  `json:"address"`
	Counter    float64 `json:"counter"`
}

type Authority struct {
	Type   string    `json:"type"`
	Subdag SubdagOut `json:"subdag"`
}

type SubdagOut struct {
	SubdagIn               map[string][]SubdagDetail `json:"subdag"`
	ElectionCertificateIds []string                  `json:"election_certificate_ids"`
}

type SubdagDetail struct {
	CertificateId string      `json:"certificate_id"`
	BatchHeader   BatchHeader `json:"batch_header"`
	Signatures    interface{} `json:"signatures"`
}

type BatchHeader struct {
	Version                    int      `json:"version"`
	BatchId                    string   `json:"batch_id"` //out
	Author                     string   `json:"author"`
	Round                      int64    `json:"round"` //out
	Timestamp                  int64    `json:"timestamp"`
	CommitteeId                string   `json:"committee_id"`
	TransmissionIds            []string `json:"transmission_ids"`
	PreviousCertificateIds     []string `json:"previous_certificate_ids"`
	LastElectionCertificateIds []string `json:"last_election_certificate_ids"`
	Signature                  string   `json:"signature"`
}

type Header struct {
	PreviousStateRoot string   `json:"previous_state_root"`
	TransactionsRoot  string   `json:"transactions_root"`
	FinalizeRoot      string   `json:"finalize_root"`
	RatificationsRoot string   `json:"ratifications_root"`
	SolutionsRoot     string   `json:"solutions_root"`
	SubdagRoot        string   `json:"subdag_root"`
	MetaData          MetaData `json:"metadata"`
}

type MetaData struct {
	Network int   `json:"network"`
	Round   int64 `json:"round"`
	Height  int64 `json:"height"`
	//TotalSupplyInMicrocredits float64 `json:"total_supply_in_microcredits"`
	CumulativeWeight      float64 `json:"cumulative_weight"`
	CumulativeProofTarget float64 `json:"cumulative_proof_target"`
	CoinbaseTarget        float64 `json:"coinbase_target"`
	ProofTarget           float64 `json:"proof_target"`
	LastCoinbaseTarget    int64   `json:"last_coinbase_target"`
	LastCoinbaseTimestamp int64   `json:"last_coinbase_timestamp"`
	Timestamp             int64   `json:"timestamp"`
}

type TransactionSpec struct {
	Status      string      `json:"status"`
	Type        string      `json:"type"`
	Index       int         `json:"index"`
	Transaction Transaction `json:"transaction"`
	Finalize    []Finalize  `json:"finalize"`
}

type Finalize struct {
	Type      string `bson:"type" json:"type"`
	MappingId string `bson:"mapping_id" json:"mapping_id"`
	Index     int    `bson:"index" json:"index"`
	KeyId     string `bson:"key_id" json:"key_id"`
	ValueId   string `bson:"value_id" json:"value_id"`
}

type Ratification struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

type Proofw struct {
	X        string `json:"x"`
	Y        string `json:"y"`
	Infinity bool   `json:"infinity"`
}
