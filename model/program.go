package model

import (
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ProgramInDb struct {
	ProgramID      string `bson:"pd"`
	Height         int64  `bson:"ht"`
	Owner          string `bson:"or"` //合约部署地址
	OwnerSignature string `bson:"os"` //合约部署地址的签名
	TransactionID  string `bson:"ti"` //部署交易ID
	TimesCalled    int    `bson:"tc"` //部署时初始化为0次，调用时增加调用次数
	DeployTime     int64  `bson:"dt"`
	UpdateAt       int64  `bson:"ua"`
}

func (*ProgramInDb) TableName() string {
	return "program"
}

func SaveProgramList(ts map[string]*ProgramInDb) error {
	if len(ts) == 0 {
		return nil
	}

	var programs []interface{}
	for _, m := range ts {
		programs = append(programs, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&ProgramInDb{}).TableName(), programs); err != nil {
		return fmt.Errorf("SaveProgramList InsertMany | %v", err)
	}
	return nil
}

func AddTimesByProgramID(programId string) {
	filter := bson.M{"pd": programId}
	updates := bson.M{"$inc": bson.M{"tc": 1}, "$set": bson.M{"ua": time.Now().Unix()}}
	opts := options.Update().SetUpsert(false)

	if _, err := mongodb.GetCollection((&ProgramInDb{}).TableName()).UpdateOne(context.TODO(), filter, updates, opts); err != nil {
		logger.Errorf("AddTimesByProgramID UpdateOne (programID=%s) | %v", programId, err)
	}
}

func JudgeProgramInDb(height int64) bool {
	var block *ProgramInDb
	found, err := mongodb.FindOne(context.TODO(), (&ProgramInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false
	}

	if found {
		return true
	}
	return false
}

func DeleteProgramByHeight(height int64) error {
	_, err := mongodb.DeleteMany(context.TODO(), (&ProgramInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		return fmt.Errorf("DeleteProgramByHeight DeleteMany(height=%d) | %v", height, err)
	}
	return nil
}

//func UpdateTimes(programId string, times int) {
//	mongodb.UpdateOne(context.TODO(),(&ProgramInDb{}).TableName(),)
//}
