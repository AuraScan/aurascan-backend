package model

import (
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type TransferInDb struct {
	TransitionId    string  `json:"transition_id" bson:"ti"`
	TransferType    int     `json:"transfer_type" bson:"tt"` //1表示transfer_public、2表示transfer_private、3表示transfer_private_to_public、4表示transfer_public_to_private
	TransferTypeStr string  `json:"transfer_type_str" bson:"tts"`
	From            string  `json:"from" bson:"fm"`
	FromPrivate     string  `bson:"fp"`
	To              string  `bson:"to"`
	ToPrivate       string  `bson:"te"`
	Credits         float64 `bson:"cs"`
	AmountPrivate   string  `bson:"ap"`
	Height          int64   `bson:"ht"`
	Timestamp       int64   `bson:"tp"`
	Status          string  `bson:"ss"`
}

type TransferRes struct {
	TransitionId string `json:"transition_id"`
	TransferType string `json:"transfer_type"`
	From         string `json:"from"`
	To           string `json:"to"`
	Credits      string `json:"credits"`
	Height       int64  `json:"height"`
	Timestamp    int64  `json:"timestamp"`
	Status       string `json:"status"`
}

func (*TransferInDb) TableName() string {
	return "transfer"
}

func GetTransferByAddr(addr string, page, pageSize int) ([]*TransferRes, int64) {
	var transfersRes = make([]*TransferRes, 0)

	tableName := (&TransferInDb{}).TableName()
	filter := bson.M{"$or": []bson.M{{"fm": addr}, {"to": addr}}}

	count, err := mongodb.Count(context.TODO(), tableName, filter)
	if err != nil {
		logger.Errorf("GetTransferByAddr Count (%s) | %v", addr, err)
		return nil, 0
	}

	var transferInDbs []*TransferInDb
	if err := mongodb.Find(context.TODO(), tableName, filter, nil,
		bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &transferInDbs); err != nil {
		return nil, 0
	}
	for _, v := range transferInDbs {
		transferRes := &TransferRes{
			TransitionId: v.TransitionId,
			TransferType: v.TransferTypeStr,
			From:         v.From,
			To:           v.To,
			Credits:      fmt.Sprintf("%.6f", v.Credits/1000000),
			Height:       v.Height,
			Timestamp:    v.Timestamp,
			Status:       v.Status,
		}
		transfersRes = append(transfersRes, transferRes)
	}

	return transfersRes, count
}
