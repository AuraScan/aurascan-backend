package model

import (
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"time"
)

type TransitionInDb struct {
	Id            string   `bson:"id"`
	TransactionId string   `bson:"ti"`
	Height        int64    `bson:"ht"`
	Timestamp     int64    `bson:"tp"`
	State         string   `bson:"se"`
	Program       string   `bson:"pm"`
	Function      string   `bson:"fn"`
	Inputs        []Input  `bson:"is"`
	Outputs       []Output `bson:"os"`
	Tpk           string   `bson:"tk"`
	Tcm           string   `bson:"tm"`
}

func (*TransitionInDb) TableName() string {
	return "transition"
}

func (t *TransitionInDb) Save() {
	if err := mongodb.InsertOne(context.TODO(), t.TableName(), t); err != nil {
		logger.Errorf("TransitionInDb.Save InsertOne | %v", err)
	}
}

func SaveTransitionList(ts []*TransitionInDb) (err error) {
	if len(ts) == 0 {
		return
	}

	var transactions []interface{}
	for _, m := range ts {
		transactions = append(transactions, m)
	}

	if err = mongodb.InsertMany(context.TODO(), (&TransitionInDb{}).TableName(), transactions); err != nil {
		return fmt.Errorf("SaveTransitionList InsertMany | %v", err)
	}
	return nil
}

func JudgeTransitionInDb(height int64) bool {
	var block *TransitionInDb
	found, err := mongodb.FindOne(context.TODO(), (&TransitionInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false
	}

	if found {
		return true
	}
	return false
}

// 通过时间范围查询program被调用的次数
func GetTimesCalledByTimeRange(start, end int64, isTotal bool, exclude []string) ([]*ProgramCountInDb, []interface{}, int) {
	var programCounts []interface{}
	var programCountsStruct = make([]*ProgramCountInDb, 0)
	date := time.Unix(start, 0).Format(util.GolangDateFormat)
	tn := time.Now()

	filter := bson.M{"tp": bson.M{"$gte": start, "$lt": end}}
	if len(exclude) > 0 {
		filter["pm"] = bson.M{"$nin": exclude}
	}
	match := bson.M{"$match": filter}
	group := bson.M{"$group": bson.M{"_id": "$pm", "TimesCalled": bson.M{"$sum": 1}}}
	if isTotal {
		group = bson.M{"$group": bson.M{"_id": "", "TimesCalled": bson.M{"$sum": 1}}}
	}
	query := []bson.M{match, group}

	var data []struct {
		Program     string  `bson:"_id"`
		TimesCalled float64 `bson:"TimesCalled"`
	}
	if err := mongodb.Aggregate(context.TODO(), (&TransitionInDb{}).TableName(), query, &data); err != nil {
		logger.Errorf("GetTimesCalledByTimeRange time(%s~%s) | %v", date, time.Unix(end, 0).Format(util.GolangDateFormat), err)
		return nil, nil, 0
	}
	if len(data) > 0 {
		for _, v := range data {
			var programCount = &ProgramCountInDb{
				Program:     v.Program,
				Date:        date,
				Timestamp:   start,
				TimesCalled: v.TimesCalled,
				UpdateAt:    tn.Unix(),
			}
			if v.Program == "" {
				programCount.Program = "total"
			}
			programCounts = append(programCounts, programCount)
			programCountsStruct = append(programCountsStruct, programCount)
		}
	}
	sort.Slice(programCountsStruct, func(i, j int) bool {
		return programCountsStruct[i].TimesCalled > programCountsStruct[j].TimesCalled
	})

	return programCountsStruct, programCounts, len(programCountsStruct)
}

func DeleteTransitionInDb(height int64) error {
	_, err := mongodb.DeleteMany(context.TODO(), (&TransitionInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		return fmt.Errorf("DeleteTransitionInDb DeleteMany(height=%d) | %v", height, err)
	}
	return nil
}

type TransitionSketch struct {
	Id       string `json:"id" bson:"id"`
	Program  string `json:"program" bson:"program"`
	Function string `json:"function" bson:"function"`
	State    string `json:"state" bson:"state"`
}

type Transition struct {
	Id       string   `json:"id" bson:"id"`
	Program  string   `json:"program" bson:"program"`
	Function string   `json:"function" bson:"function"`
	Inputs   []Input  `json:"inputs" bson:"inputs"`
	Outputs  []Output `json:"outputs" bson:"outputs"`
	//Proof    string   `json:"proof" bson:"proof"`
	Tpk string `json:"tpk" bson:"tpk"`
	Tcm string `json:"tcm" bson:"tcm"`
}

type Input struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	//Tag   string `json:"tag"`
	Value string `json:"value"`
}

type Output struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	//Checksum string `json:"checksum"`
	Value string `json:"value"`
}
