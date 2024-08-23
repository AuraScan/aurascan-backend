package repository

import (
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"time"
)

func GetTransitionDetail(transitionId string) (*schema.TransitionDetailResp, error) {
	var transitionInDb *model.TransitionInDb
	if _, err := mongodb.FindOne(context.TODO(), (&model.TransitionInDb{}).TableName(), bson.M{"id": transitionId}, &transitionInDb); err != nil {
		logger.Errorf("GetTransitionDetail FindOne(%s) | %v", transitionId, err)
		return nil, err
	}
	if transitionInDb != nil {
		var transitionDetailResp = &schema.TransitionDetailResp{
			TransitionId:  transitionInDb.Id,
			TransactionId: transitionInDb.TransactionId,
			State:         transitionInDb.State,
			Program:       transitionInDb.Program,
			Function:      transitionInDb.Function,
			Tpk:           transitionInDb.Tpk,
			Tcm:           transitionInDb.Tcm,
			Input:         transitionInDb.Inputs,
			Output:        transitionInDb.Outputs,
		}
		return transitionDetailResp, nil
	}

	return nil, nil
}

func GetTransitionById(ti string) (*model.TransitionInDb, error) {
	var td *model.TransitionInDb
	_, err := mongodb.FindOne(context.TODO(), (&model.TransitionInDb{}).TableName(), bson.M{"id": ti}, &td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GetTransitionByTransactionId(ti string) []*model.TransitionInDb {
	var transitionInDb []*model.TransitionInDb
	if err := mongodb.Find(context.TODO(), (&model.TransitionInDb{}).TableName(), bson.M{"ti": ti}, nil, bson.D{{"fn", -1}}, 0, 0, &transitionInDb); err != nil {
		logger.Errorf("GetTransitionByTransactionId FindOne(%s) | %v", ti, err)
		return nil
	}
	return transitionInDb
}

// 获取交易列表
func GetTransitionsByPage(page, pageSize int, programId string) ([]*schema.TransitionListResp, int64) {
	var transitions []*model.TransitionInDb
	var transitionsResp = make([]*schema.TransitionListResp, 0)

	filter := bson.M{}
	var count int64
	var err error

	if programId != "" {
		filter = bson.M{"pm": programId}
		count, err = mongodb.GetCollection((&model.TransitionInDb{}).TableName()).CountDocuments(context.TODO(), filter)
		if err != nil {
			logger.Errorf("GetTransitionsByPage CountDocuments | %v", err)
			return transitionsResp, 0
		}
	} else {
		count, err = mongodb.GetCollection((&model.TransitionInDb{}).TableName()).EstimatedDocumentCount(context.TODO())
		if err != nil {
			logger.Errorf("GetTransitionsByPage EstimatedDocumentCount | %v", err)
			return transitionsResp, 0
		}
	}

	if err := mongodb.Find(context.TODO(), (&model.TransitionInDb{}).TableName(), filter, nil, bson.D{{"ht", -1}, {"id", 1}}, util.Offset(pageSize, page), int64(pageSize), &transitions); err != nil {
		logger.Errorf("GetTransitionsByPage Find | %v", err)
		return transitionsResp, 0
	}

	if len(transitions) > 0 {
		for _, v := range transitions {
			var transaction = &schema.TransitionListResp{
				ID:        v.Id,
				Height:    v.Height,
				Timestamp: v.Timestamp,
				Time:      time.Unix(v.Timestamp, 0).Format(util.GoLangTimeFormat),
				Program:   v.Program,
				Function:  v.Function,
			}
			transitionsResp = append(transitionsResp, transaction)
		}
	}

	return transitionsResp, count
}

// 获取最近24H的前十Program调用
func GetTopProgramCalledIn24h() []*schema.ProgramCalledChartResp {
	end := time.Now().Unix()
	start := time.Now().Unix() - util.DaySeconds
	programCount, _, _ := model.GetTimesCalledByTimeRange(start, end, false, []string{"credits.aleo"})
	sort.Slice(programCount, func(i, j int) bool {
		return programCount[i].TimesCalled > programCount[j].TimesCalled
	})
	if len(programCount) > 10 {
		programCount = programCount[:10]
	}

	programCountAll, _, _ := model.GetTimesCalledByTimeRange(start, end, true, []string{"credits.aleo"})
	//
	var programCharts = make([]*schema.ProgramCalledChartResp, 0)
	var topCalled float64
	if len(programCount) > 0 && len(programCountAll) > 0 {
		for _, v := range programCount {
			var programChart = &schema.ProgramCalledChartResp{
				Program:   v.Program,
				CallTimes: int(v.TimesCalled),
			}
			topCalled += v.TimesCalled
			programCharts = append(programCharts, programChart)
		}
	}

	var otherCalled float64
	if len(programCountAll) > 0 {
		otherCalled = programCountAll[0].TimesCalled - topCalled
		if otherCalled > 0 {
			programCharts = append(programCharts, &schema.ProgramCalledChartResp{
				Program:   "others",
				CallTimes: int(otherCalled),
			})
		}
	}
	return programCharts
}
