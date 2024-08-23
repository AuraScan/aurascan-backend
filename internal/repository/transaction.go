package repository

import (
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func GetTransactionDetail(ti string) (*schema.TransactionDetailResp, error) {
	var tdr *schema.TransactionDetailResp
	txn, err := GetTransactionById(ti)
	if err != nil {
		logger.Errorf("GetTransactionDetail(%s) GetTransactionById | %v", ti, err)
		return nil, err
	}
	tdr = &schema.TransactionDetailResp{
		Id:       txn.TransactionId,
		Height:   txn.Height,
		Type:     txn.Type,
		Status:   txn.OuterStatus,
		Time:     time.Unix(txn.Timestamp, 0).Format(util.GoLangTimeFormat),
		Fee:      txn.Fee,
		Finalize: txn.Finalize,
	}
	var transitionsResp = make([]*schema.TransitionListInTransactionResp, 0)
	transitions := GetTransitionByTransactionId(txn.TransactionId)
	if len(transitions) > 0 {
		for _, v := range transitions {
			transitionsResp = append(transitionsResp, &schema.TransitionListInTransactionResp{
				ID:       v.Id,
				Program:  v.Program,
				Function: v.Function,
				State:    v.State,
			})
		}
	}
	tdr.Transitions = transitionsResp
	return tdr, nil
}

func GetTransactionById(ti string) (*model.TransactionInDb, error) {
	var td *model.TransactionInDb
	_, err := mongodb.FindOne(context.TODO(), (&model.TransactionInDb{}).TableName(), bson.M{"ti": ti}, &td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GetTransactionsByHeight(height int64, page, pageSize int) ([]*schema.TransactionListInBlockResp, int64) {
	var transactions []*model.TransactionInDb
	var transactionsResp = make([]*schema.TransactionListInBlockResp, 0)

	total, err := mongodb.Count(context.TODO(), (&model.SolutionInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		logger.Errorf("GetSolutionListByHeight Count | %v", err)
		return transactionsResp, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.TransactionInDb{}).TableName(), bson.M{"ht": height}, nil, bson.D{{"ti", 1}}, util.Offset(pageSize, page), int64(pageSize), &transactions); err != nil {
		logger.Errorf("GetTransactionsByHeight Find(%d) | %v", height, err)
		return transactionsResp, 0
	}

	if len(transactions) > 0 {
		for _, v := range transactions {
			var transaction = &schema.TransactionListInBlockResp{
				Id:     v.TransactionId,
				Type:   v.Type,
				Status: v.OuterStatus,
				Fee:    v.Fee,
			}
			transactionsResp = append(transactionsResp, transaction)
		}
	}

	return transactionsResp, total
}

// 获取交易列表
func GetTransactionsByPage(page, pageSize int) ([]*schema.TransactionListResp, int64) {
	var transactions []*model.TransactionInDb
	var transactionsResp = make([]*schema.TransactionListResp, 0)

	count, err := mongodb.GetCollection((&model.TransactionInDb{}).TableName()).EstimatedDocumentCount(context.TODO())
	if err != nil {
		logger.Errorf("GetTransactionsByPage Count | %v", err)
		return transactionsResp, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.TransactionInDb{}).TableName(), bson.M{}, nil, bson.D{{"ht", -1}, {"ti", 1}}, util.Offset(pageSize, page), int64(pageSize), &transactions); err != nil {
		logger.Errorf("GetTransactionsByPage Find | %v", err)
		return transactionsResp, 0
	}

	if len(transactions) > 0 {
		for _, v := range transactions {
			var transaction = &schema.TransactionListResp{
				Id:        v.TransactionId,
				Height:    v.Height,
				Time:      time.Unix(v.Timestamp, 0).Format(util.GoLangTimeFormat),
				Timestamp: v.Timestamp,
				Type:      v.Type,
				Status:    v.OuterStatus,
				Fee:       v.Fee,
			}
			transactionsResp = append(transactionsResp, transaction)
		}
	}

	return transactionsResp, count
}
