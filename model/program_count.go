package model

import (
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type ProgramCountInDb struct {
	Program     string  `bson:"pm"`
	Date        string  `bson:"de"` //"2024-01-01"
	Timestamp   int64   `bson:"tp"`
	TimesCalled float64 `bson:"tc"`
	UpdateAt    int64   `bson:"ua"`
}

func (*ProgramCountInDb) TableName() string {
	return "program_count"
}

type ProgramCountResp struct {
	Program string `json:"program"`
	Count   int64  `json:"count"`
}

func FindProgramCountByDate(date string) bool {
	var pc *ProgramCountInDb
	found, err := mongodb.FindOne(context.TODO(), (&ProgramCountInDb{}).TableName(), bson.M{"de": date}, &pc)
	if err != nil {
		logger.Errorf("FindProgramCountByDate FindOne(%s) | %v", date, err)
		return false
	}
	return found
}
