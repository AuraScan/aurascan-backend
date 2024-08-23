package model

import (
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type SolutionInDb struct {
	Address string `json:"address" bson:"as"`
	//Nonce       float64 `json:"nonce" bson:"ne"` //19位
	//Commitment string `json:"commitment" bson:"ct"`
	Version     int     `json:"version" bson:"vn"`
	SolutionId  string  `json:"solution_id" bson:"si"`
	EpochHash   string  `json:"epoch_hash" bson:"eh"`
	Counter     float64 `json:"counter" bson:"cr"`
	Height      int64   `json:"height" bson:"ht"`
	Reward      float64 `json:"reward" bson:"rd"` //coinbase奖励
	Target      float64 `json:"target" bson:"tt"`
	ProofTarget float64 `json:"proof_target" bson:"pt"`
	Timestamp   int64   `json:"timestamp" bson:"tp"`
}

type SolutionRes struct {
	Address    string  `json:"address" bson:"as"`
	Commitment string  `json:"commitment" bson:"ct"`
	Version    int     `json:"version" bson:"vn"`
	SolutionId string  `json:"solution_id" bson:"si"`
	EpochHash  string  `json:"epoch_hash" bson:"eh"`
	Counter    float64 `json:"counter" bson:"cr"`
	Height     int64   `json:"height" bson:"ht"`
	BlockHash  string  `json:"block_hash" bson:"bh"`
	Reward     float64 `json:"reward" bson:"rd"` //coinbase奖励
	Target     float64 `json:"target" bson:"tt"`
	Timestamp  int64   `json:"timestamp" bson:"timestamp"` //
	Epoch      int64   `json:"epoch" bson:"epoch"`
}

func (*SolutionInDb) TableName() string {
	return "solution"
}

func (t *SolutionInDb) Save() {
	if err := mongodb.InsertOne(context.TODO(), t.TableName(), t); err != nil {
		logger.Errorf("SolutionInDb.Save InsertOne | %v", err)
	}
}

func SaveSolutionList(ts []*SolutionInDb) error {
	if len(ts) == 0 {
		return nil
	}

	var solutions []interface{}
	for _, m := range ts {
		solutions = append(solutions, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&SolutionInDb{}).TableName(), solutions); err != nil {
		return fmt.Errorf("SaveSolutionList InsertMany | %v", err)
	}
	return nil
}

func GetSolutions(addr string, page, pageSize int, height int64) ([]*SolutionRes, int, error) {
	var solutions []*SolutionInDb
	query := bson.M{"as": addr}
	if height != 0 {
		query["ht"] = height
	}

	count, err := mongodb.Count(context.TODO(), (&SolutionInDb{}).TableName(), query)
	if err != nil {
		return nil, 0, err
	}

	err = mongodb.Find(context.TODO(), (&SolutionInDb{}).TableName(), query, nil, bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &solutions)
	if err != nil {
		return nil, 0, err
	}
	var solutionRess = make([]*SolutionRes, 0)
	for _, v := range solutions {
		var solutionRes = &SolutionRes{
			Address:    v.Address,
			Commitment: "",
			Version:    v.Version,
			SolutionId: v.SolutionId,
			EpochHash:  v.EpochHash,
			Counter:    v.Counter,
			Height:     v.Height,
			BlockHash:  "",
			Reward:     v.Reward,
			Target:     v.Target,
			Timestamp:  v.Timestamp,
			Epoch:      v.Height / 360,
		}
		solutionRess = append(solutionRess, solutionRes)
	}
	return solutionRess, int(count), nil
}

func JudgeSolutionInDb(height int64) bool {
	var block *SolutionInDb
	found, err := mongodb.FindOne(context.TODO(), (&SolutionInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false
	}

	if found {
		return true
	}
	return false
}

func DeleteSolutionByHeight(height int64) error {
	_, err := mongodb.DeleteMany(context.TODO(), (&SolutionInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		return fmt.Errorf("DeleteSolutionByHeight DeleteMany(height=%d) | %v", height, err)
	}
	return nil
}
