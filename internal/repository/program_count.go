package repository

import (
	"aurascan-backend/internal/config"
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	util2 "aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"ch-common-package/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// 计算最近一个月调用前十名的图表
func CalculateProgramCalledChartOneMonth() []*schema.ProgramCountMonthChart {
	var programCountMonthChart = make([]*schema.ProgramCountMonthChart, 0)
	ts := util.GetNullUTCPoint(time.Now().Unix()) - util2.DaySeconds*30
	genesis := util.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	dateList := util2.GetTimestampListByStart(ts, genesis)
	programs := getTop10ProgramCalledOneMonth(ts)

	for _, program := range programs {
		var programCalledTime = make([]*schema.ProgramCalledTime, 0)
		programCounts := GetProgramCalled(program, ts)

		for _, ts := range dateList {
			date := time.Unix(ts, 0).Format(util2.GoLangTimeFormat)
			value, ok := programCounts[ts]
			if ok {
				programCalledTime = append(programCalledTime, &schema.ProgramCalledTime{
					Date:  date,
					Value: value.TimesCalled,
				})
			} else {
				programCalledTime = append(programCalledTime, &schema.ProgramCalledTime{
					Date:  date,
					Value: 0,
				})
			}
		}

		programCountMonthChart = append(programCountMonthChart, &schema.ProgramCountMonthChart{
			Program: program,
			Times:   programCalledTime,
		})
	}
	return programCountMonthChart
}

// 获取近一个月总调用的前十名，对前十名进行每日统计
func getTop10ProgramCalledOneMonth(ts int64) []string {
	filter := bson.M{"tc": bson.M{"$gt": 2}, "pm": bson.M{"$nin": []string{"total", "credits.aleo"}}}
	if ts != 0 {
		filter["tp"] = bson.M{"$gte": ts}
	}
	match := bson.M{"$match": filter}
	group := bson.M{"$group": bson.M{"_id": "$pm", "totalCalled": bson.M{"$sum": "$tc"}}}
	sort := bson.M{"$sort": bson.M{"totalCalled": -1}}
	limit := bson.M{"$limit": 10}
	query := []bson.M{match, group, sort, limit}

	var data []struct {
		Program     string  `bson:"_id"`
		TotalCalled float64 `bson:"totalCalled"`
	}

	if err := mongodb.Aggregate(context.TODO(), (&model.ProgramCountInDb{}).TableName(), query, &data); err != nil {
		logger.Errorf("GetTop10ProgramCalledOneMonth Aggregate | %v", err)
		return nil
	}

	var res []string
	if len(data) > 0 {
		for _, v := range data {
			res = append(res, v.Program)
		}
	}
	return res
}

// 获取单个Program一个月的调用情况
func GetProgramCalled(program string, start int64) map[int64]*model.ProgramCountInDb {
	var programCountInDb []*model.ProgramCountInDb
	filter := bson.M{"pm": program, "tp": bson.M{"$gte": start}}
	if err := mongodb.Find(context.TODO(), (&model.ProgramCountInDb{}).TableName(), filter, nil, bson.D{{"tp", 1}}, 0, 0, &programCountInDb); err != nil {
		logger.Errorf("GetProgramCalled Find(program=%s, start=%s) | %v", program, time.Unix(start, 0).Format(util2.GoLangTimeFormat), err)
		return nil
	}

	var res = make(map[int64]*model.ProgramCountInDb)
	for _, v := range programCountInDb {
		res[v.Timestamp] = v
	}

	return res
}
