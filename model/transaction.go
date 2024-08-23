package model

import (
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionInDb struct {
	BlockHash string `bson:"bh"`
	Height    int64  `bson:"ht"`
	Timestamp int64  `bson:"tp"`

	TransactionId string `bson:"ti"`
	Type          string `bson:"te"`

	OuterStatus string `bson:"os"` // 外层Status
	OuterType   string `bson:"ot"` // 外层Type
	OuterIndex  int    `bson:"oi"` // 外层Index

	DeploymentId string  `bson:"di"`
	Fee          float64 `bson:"fee"`
	BaseFee      float64 `bson:"bf"`
	PriorityFee  float64 `bson:"pf"`

	Finalize []Finalize `bson:"fe"`
}

func (*TransactionInDb) TableName() string {
	return "transaction"
}

func (t *TransactionInDb) Save() {
	if err := mongodb.InsertOne(context.TODO(), t.TableName(), t); err != nil {
		logger.Errorf("TransactionInDb.Save InsertOne | %v", err)
	}
}

func SaveTransactionList(ts []*TransactionInDb) error {
	if len(ts) == 0 {
		return nil
	}

	var transactions []interface{}
	for _, m := range ts {
		transactions = append(transactions, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&TransactionInDb{}).TableName(), transactions); err != nil {
		return fmt.Errorf("SaveTransactionList InsertMany | %v", err)
	}
	return nil
}

func JudgeTransactionInDb(height int64) bool {
	var block *TransactionInDb
	found, err := mongodb.FindOne(context.TODO(), (&TransactionInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false
	}

	if found {
		return true
	}
	return false
}

func DeleteTransactionInDb(height int64) error {
	_, err := mongodb.DeleteMany(context.TODO(), (&TransactionInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		return fmt.Errorf("DeleteTransactionInDb DeleteMany(height=%d) | %v", height, err)
	}
	return nil
}

type Transaction struct {
	Type       string     `json:"type"`
	Id         string     `json:"id"`
	Owner      Owner      `json:"owner"`
	Deployment Deployment `json:"deployment"`
	Execution  Execution  `json:"execution"`
	Fee        Fee        `json:"fee"`
}

type Fee struct {
	Transition      Transition `json:"transition" bson:"transition"`
	GlobalStateRoot string     `json:"global_state_root" bson:"global_state_root"`
	//Inclusion       string     `json:"inclusion" bson:"inclusion"`
	//Proof string `json:"proof"` //过长，考虑是否入库或者存入ssdb中
}

type Owner struct {
	Address   string `json:"address" bson:"address"`
	Signature string `json:"signature" bson:"signature"`
}

type Deployment struct {
	Edition       int64           `json:"edition" bson:"edition"`
	Program       string          `json:"program" bson:"program"`
	VerifyingKeys [][]interface{} `json:"verifying_keys" bson:"verifying_keys"`
}

type Execution struct {
	Transitions     []Transition `json:"transitions" bson:"transitions"`
	GlobalStateRoot string       `json:"global_state_root" bson:"global_state_root"`
	//Proof string `json:"proof"` //过长，考虑是否入库或者存入ssdb中
}
