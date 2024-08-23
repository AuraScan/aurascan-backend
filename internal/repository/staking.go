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

func GetStakeChartOneMonth() []*schema.StakingDailyResp {
	ts := util.GetNullPoint(time.Now().Unix())
	startTs := ts - util2.DaySeconds*30
	var stakingDailyInDb []*model.StakingDailyInDb
	if err := mongodb.Find(context.TODO(), (&model.StakingDailyInDb{}).TableName(), bson.M{"timestamp": bson.M{"$gte": startTs}},
		nil, bson.D{{"timestamp", -1}}, 0, 0, &stakingDailyInDb); err != nil {
		logger.Errorf("GetStakeChartOneMonth Find | %v", err)
	}
	var stakingDailyMap = make(map[string]*schema.StakingDailyResp, 0)

	for _, v := range stakingDailyInDb {
		stakingDailyMap[v.Date] = &schema.StakingDailyResp{
			Date:          v.Date,
			TotalBond:     v.TotalBond,
			ValidatorBond: v.ValidatorBond,
			DelegatorBond: v.DelegatorBond,
		}
	}

	var stakingDailyResps = make([]*schema.StakingDailyResp, 0)
	genesis := util.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	timeList := util2.GetDateListByStart(startTs, genesis)
	for _, v := range timeList {
		value, ok := stakingDailyMap[v]
		if !ok {
			stakingDailyResps = append(stakingDailyResps, &schema.StakingDailyResp{
				Date:          v,
				TotalBond:     0,
				ValidatorBond: 0,
				DelegatorBond: 0,
			})
		} else {
			stakingDailyResps = append(stakingDailyResps, value)
		}
	}
	return stakingDailyResps
}
