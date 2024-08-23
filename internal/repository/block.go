package repository

import (
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"aurascan-backend/util"
	"ch-common-package/cache"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	util2 "ch-common-package/util"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

// 删除某高度所有数据
func DeleteObsoleteDataByHeight(height int64) error {
	ctx := context.Background()
	filter := bson.M{"ht": height}

	_, err := mongodb.GetCollection((&model.AuthorityInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete authority by height(%d) err: %s", height, err.Error())
		return err
	}

	_, err = mongodb.GetCollection((&model.TransitionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete transition by height(%d) err: %s", height, err.Error())
		return err
	}

	_, err = mongodb.GetCollection((&model.TransactionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete transaction by height(%d) err: %s", height, err.Error())
		return err
	}

	_, err = mongodb.GetCollection((&model.ProgramInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete program by height(%d) err: %s", height, err.Error())
		return err
	}

	_, err = mongodb.GetCollection((&model.SolutionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete solution by height(%d) err: %s", height, err.Error())
		return err
	}

	_, err = mongodb.GetCollection((&model.BlockInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete block by height(%d) err: %s", height, err.Error())
		return err
	}

	return nil
}

// 删除指定高度范围所有数据（一次不能超过100个高度）
func DeleteObsoleteDataByHeightRange(start, end int64) {
	ctx := context.Background()
	filter := bson.M{"ht": bson.M{"$gte": start, "$lte": end}}

	_, err := mongodb.GetCollection((&model.AuthorityInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete authority by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	_, err = mongodb.GetCollection((&model.TransitionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete transition by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	_, err = mongodb.GetCollection((&model.TransactionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete transaction by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	_, err = mongodb.GetCollection((&model.ProgramInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete program by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	_, err = mongodb.GetCollection((&model.SolutionInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete solution by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	_, err = mongodb.GetCollection((&model.BlockInDb{}).TableName()).DeleteMany(ctx, filter)
	if err != nil {
		logger.Errorf("DeleteObsoleteData delete block by heightRange(%d~%d) err: %s", start, end, err.Error())
		return
	}

	return
}

func GetLatestHeightInDb() (int64, error) {
	var blocks []*model.BlockInDb
	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{},
		map[string]int{"ht": 1}, bson.D{{"ht", -1}}, 0, 1, &blocks); err != nil {
		return 0, err
	}
	if len(blocks) > 0 {
		return blocks[0].Height, nil
	}
	return 0, nil
}

func GetFirstHeightInDb() (int64, error) {
	var blocks []*model.BlockInDb
	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{},
		map[string]int{"ht": 1}, bson.D{{"ht", 1}}, 0, 1, &blocks); err != nil {
		return 0, err
	}
	if len(blocks) > 0 {
		return blocks[0].Height, nil
	}
	return 0, nil
}

func GetFirstHeightTimeInDb() (int64, error) {
	var blocks []*model.BlockInDb
	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"sn": bson.M{"$gt": 0}},
		map[string]int{"tp": 1}, bson.D{{"tp", 1}}, 0, 1, &blocks); err != nil {
		return 0, err
	}
	if len(blocks) > 0 {
		return blocks[0].Timestamp, nil
	}
	return 0, nil
}

func GetLatestBlockInDb() (*model.BlockInDb, error) {
	var blocks []*model.BlockInDb
	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{},
		nil, bson.D{{"ht", -1}}, 0, 1, &blocks); err != nil {
		return nil, err
	}
	if len(blocks) > 0 {
		return blocks[0], nil
	}
	return nil, mongo.ErrNoDocuments
}

func GetBlockByHeight(height int64) (*model.BlockInDb, error) {
	var block *model.BlockInDb
	found, err := mongodb.FindOne(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}
	return block, nil
}

func GetGenesisTime() int64 {
	var block *model.BlockInDb
	//官方0的时间戳有bug
	found, err := mongodb.FindOne(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"ht": 1}, &block)
	if err != nil {
		logger.Errorf("GetGenesisTime | %v", err)
		return 0
	}

	if !found {
		logger.Warnf("GetGenesisTime can't find height 1 in db")
		return 0
	}
	return block.Timestamp
}

// 获取创世区块当日0点
func GetGenesisNullPoint() int64 {
	var block *model.BlockInDb
	//官方0的时间戳有bug
	found, err := mongodb.FindOne(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"ht": 1}, &block)
	if err != nil {
		logger.Errorf("GetGenesisTime | %v", err)
		return 0
	}

	if !found {
		logger.Warnf("GetGenesisTime can't find height 1 in db")
		return 0
	}
	return util2.GetNullPoint(block.Timestamp)
}

func GetBlockByHash(hash string) (*model.BlockInDb, error) {
	var block *model.BlockInDb
	found, err := mongodb.FindOne(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"bh": hash}, &block)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}
	return block, nil
}

func JudgeBlockInDb(height int64) (bool, error) {
	var block *model.BlockInDb
	found, err := mongodb.FindOne(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false, err
	}

	if found {
		return true, nil
	}

	if model.JudgeProgramInDb(height) {
		if err = model.DeleteProgramByHeight(height); err != nil {
			logger.Errorf("JudgeBlockInDb | %v", err)
		}
	}
	if model.JudgeSolutionInDb(height) {
		if err = model.DeleteProgramByHeight(height); err != nil {
			logger.Errorf("JudgeBlockInDb | %v", err)
		}
	}
	if model.JudgeTransactionInDb(height) {
		if err = model.DeleteProgramByHeight(height); err != nil {
			logger.Errorf("JudgeBlockInDb | %v", err)
		}
	}
	if model.JudgeTransitionInDb(height) {
		if err = model.DeleteProgramByHeight(height); err != nil {
			logger.Errorf("JudgeBlockInDb | %v", err)
		}
	}
	if model.JudgeAuthorityInDb(height) {
		if err = model.DeleteProgramByHeight(height); err != nil {
			logger.Errorf("JudgeBlockInDb | %v", err)
		}
	}

	return false, nil
}

func GetSolutionNum(ts int64) int64 {
	var block = &model.BlockInDb{}
	var data []struct {
		TotalSolution int64 `bson:"TotalSolution"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gte": ts, "$lt": ts + 86400}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalSolution": bson.M{"$sum": "$sn"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetSolutionNum Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	return data[0].TotalSolution
}

func GetRewardByDay(ts int64) float64 {
	var block = &model.BlockInDb{}
	var data []struct {
		DayReward float64 `bson:"DayReward"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gte": ts, "$lt": ts + 86400}}}
	var group = bson.M{"$group": bson.M{"_id": "", "DayReward": bson.M{"$sum": "$br"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetRewardByDay Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	return data[0].DayReward / 1000000
}

func GetTotalReward() float64 {
	var block = &model.BlockInDb{}
	var data []struct {
		TotalReward float64 `bson:"TotalReward"`
	}
	var match = bson.M{"$match": bson.M{}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalReward": bson.M{"$sum": "$br"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetTotalReward Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	return data[0].TotalReward / 1000000
}

// 最近15m p/s <=> proof/second
func GetNetworkHashRate() float64 {
	var block = &model.BlockInDb{}
	var data []struct {
		TotalProofTarget float64 `bson:"TotalProofTarget"`
	}
	var match = bson.M{"$match": bson.M{"timestamp": bson.M{"$gt": time.Now().Unix() - 900}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalProofTarget": bson.M{"$sum": "$pt"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetNetworkHashRate Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	hashRate := data[0].TotalProofTarget / 900

	return hashRate
}

// 最近24h c/s <=> coinbase/second
func GetNetworkCS() float64 {
	var block = &model.BlockInDb{}
	var data []struct {
		TotalCoinbaseTarget float64 `bson:"TotalCoinbaseTarget"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gt": time.Now().Unix() - 86400}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalCoinbaseTarget": bson.M{"$sum": "$br"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetNetworkCS Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	hashRate := data[0].TotalCoinbaseTarget / 86400

	return hashRate
}

// 根据地址获取累计收益
func GetRewardByAddr(addr string) float64 {
	match := bson.M{"$match": bson.M{"as": addr}}
	group := bson.M{"$group": bson.M{"_id": "", "totalReward": bson.M{"$sum": "$rd"}}}
	query := []bson.M{match, group}

	var data []struct {
		Id          string  `bson:"_id"`
		TotalReward float64 `bson:"totalReward"`
	}
	if err := mongodb.Aggregate(context.TODO(), (&model.SolutionInDb{}).TableName(), query, &data); err != nil {
		logger.Errorf("GetRewardByAddr (%s) | %v", addr, err)
		return 0
	}
	if len(data) > 0 {
		return data[0].TotalReward
	}
	return 0
}

// 获取区块列表
func GetBlocksByPage(page, pageSize int) ([]*schema.BlockListResp, int64) {
	var blocks []*model.BlockInDb
	var blocksResp = make([]*schema.BlockListResp, 0)

	height, err := GetLatestHeightInDb()
	if err != nil {
		logger.Errorf("GetBlocksByPage GetLatestHeightInDb | %v", err)
		return blocksResp, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{}, nil, bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &blocks); err != nil {
		logger.Errorf("GetBlocksByPage Find | %v", err)
		return blocksResp, 0
	}

	if len(blocks) > 0 {
		for _, v := range blocks {
			var block = &schema.BlockListResp{
				Height:         v.Height,
				Epoch:          v.Height / 360,
				EpochIndex:     v.Height - v.Height/360*360,
				Round:          v.Round,
				Time:           time.Unix(v.Timestamp, 0).Format(util.GoLangTimeFormat),
				ProofTarget:    v.ProofTarget,
				CoinbaseTarget: v.CoinbaseTarget,
				BlockReward:    v.BlockReward,
				CoinbaseReward: v.PuzzleReward,
				Solutions:      v.SolutionNum,
				Transactions:   v.TransactionNum,
			}
			blocksResp = append(blocksResp, block)
		}
	}

	return blocksResp, height
}

func GetBlockDetailByHeight(height int64) (*schema.BlockDetailResp, error) {
	block, err := GetBlockByHeight(height)
	if err != nil {
		logger.Errorf("GetBlockDetailByHeight GetBlockByHeight(%d) | %v", height, err)
		return nil, err
	}
	var blockDetail = &schema.BlockDetailResp{
		Height:                block.Height,
		BlockHash:             block.BlockHash,
		PreviousHash:          block.PreviousHash,
		PreviousStateRoot:     block.PreviousStateRoot,
		TransactionsRoot:      block.TransactionsRoot,
		FinalizeRoot:          block.FinalizeRoot,
		RatificationsRoot:     block.RatificationsRoot,
		CumulativeWeight:      block.CumulativeWeight,
		CumulativeProofTarget: block.CumulativeProofTarget,
		AuthorityType:         block.AuthorityType,
		Round:                 block.Round,
		BlockReward:           block.BlockReward,
		CoinbaseReward:        block.PuzzleReward,
		ProofTarget:           block.ProofTarget,
		CoinbaseTarget:        block.CoinbaseTarget,
		Network:               block.Network,
		Time:                  time.Unix(block.Timestamp, 0).Format(util.GoLangTimeFormat),
	}
	//blockDetail.Authority = model.GetBlockAuthoritys(height)
	//blockDetail.Transactions = GetTransactionsByHeight(height)

	return blockDetail, nil
}

func Get24hProofTargetChart() []*schema.ProofTargetChart {
	var res = make([]*schema.ProofTargetChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "proof_target_24h_chart", &res)
	if err != nil {
		logger.Errorf("Get24hProofTargetChart GetValue(proof_target_24h_chart) | %v", err)
	} else if exist {
		return res
	}

	end := util2.GetPoint(time.Now().Unix())
	start := end - util.DaySeconds
	for i := start; i <= end; i += util.HourSeconds {
		proofTarget, height := GetNearProofTarget(i)
		res = append(res, &schema.ProofTargetChart{
			Timestamp:   i,
			Height:      height,
			ProofTarget: proofTarget,
		})
	}

	cache.Redis.SetValue(context.TODO(), "proof_target_24h_chart", res, time.Hour)
	return res
}

// 一天每个小时取一次，取24个高度的prooftarget
func GetNearProofTarget(timestamp int64) (float64, int64) {
	var blocksInDb = make([]*model.BlockInDb, 0)
	tableName := (&model.BlockInDb{}).TableName()
	err := mongodb.Find(context.TODO(), tableName, bson.M{"tp": bson.M{"$gte": timestamp}}, nil, bson.D{{"tp", 1}}, 0, 1, &blocksInDb)
	if err != nil {
		logger.Errorf("GetNearProofTarget Find(%d) | %v", timestamp, err)
		return 0, 0
	}
	if len(blocksInDb) > 0 {
		return blocksInDb[0].ProofTarget, blocksInDb[0].Height
	}
	return 0, 0
}

// 一周每六个小时取一次，取到现在的prooftarget
func Get7dProofTargetChart() []*schema.ProofTargetChart {
	var res = make([]*schema.ProofTargetChart, 0)

	exist, err := cache.Redis.GetValue(context.TODO(), "proof_target_7d_chart", &res)
	if err != nil {
		logger.Errorf("Get7dProofTargetChart GetValue(proof_target_7d_chart) | %v", err)
	} else if exist {
		return res
	}

	end := util2.GetPoint(time.Now().Unix())
	start := end - util.DaySeconds*7
	for i := start; i <= end; i += util.HourSeconds * 6 {
		proofTarget, height := GetNearProofTarget(i)
		res = append(res, &schema.ProofTargetChart{
			Timestamp:   i,
			Height:      height,
			ProofTarget: proofTarget,
		})
	}

	cache.Redis.SetValue(context.TODO(), "proof_target_7d_chart", res, time.Hour)

	return res
}

// 从第一天开始取，每日零点取一次，取到今日零点的prooftarget
func GetAllProofTargetChart() []*schema.ProofTargetChart {
	var res = make([]*schema.ProofTargetChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "proof_target_all_chart", &res)
	if err != nil {
		logger.Errorf("GetAllProofTargetChart GetValue(proof_target_all_chart) | %v", err)
	} else if exist {
		return res
	}

	end := util2.GetNullPoint(time.Now().Unix())
	startTime, err := GetFirstHeightTimeInDb()
	if err != nil {
		logger.Errorf("GetAllProofTargetChart GetFirstHeightTimeInDb | %v", err)
		return res
	}
	start := util2.GetNullPoint(startTime) + util.DaySeconds

	for i := start; i <= end; i += util.DaySeconds {
		proofTarget, height := GetNearProofTarget(i)
		res = append(res, &schema.ProofTargetChart{
			Timestamp:   i,
			Height:      height,
			ProofTarget: proofTarget,
		})
	}

	cache.Redis.SetValue(context.TODO(), "proof_target_all_chart", res, time.Hour*2)
	return res
}

func GetNetworkPowerChart() []*schema.PowerChart {
	var powerCharts = make([]*schema.PowerChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "network_power_chart", &powerCharts)
	if err != nil {
		logger.Errorf("GetNetworkPowerChart GetValue(network_power_chart) | %v", err)
	} else if exist {
		return powerCharts
	}

	end := util2.GetPoint(time.Now().Unix())
	start := end - util.DaySeconds
	for i := start + util.HourSeconds; i <= end; i += util.HourSeconds {
		power := GetNetworkPower(i-util.HourSeconds, i)
		powerCharts = append(powerCharts, &schema.PowerChart{
			Date:      time.Unix(i, 0).Format(util.GolangDayFormat),
			Timestamp: i,
			Power:     power,
		})
	}

	cache.Redis.SetValue(context.TODO(), "network_power_chart", powerCharts, time.Hour)
	return powerCharts
}

func GetNetwork7dPowerChart() []*schema.PowerChart {
	var powerCharts = make([]*schema.PowerChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "network_power_7d_chart", &powerCharts)
	if err != nil {
		logger.Errorf("GetNetwork7dPowerChart GetValue(network_power_7d_chart) | %v", err)
	} else if exist {
		return powerCharts
	}

	end := util2.GetPoint(time.Now().Unix())
	start := end - 7*util.DaySeconds
	for i := start; i <= end; i += util.HourSeconds * 6 {
		power := GetNeworkPowerSixHourAgo(i)
		powerCharts = append(powerCharts, &schema.PowerChart{
			Date:  time.Unix(i, 0).Format(util.GoLangTimeFormat),
			Power: power,
		})
	}

	cache.Redis.SetValue(context.TODO(), "network_power_7d_chart", powerCharts, time.Hour*2)
	return powerCharts
}

// TODO：待定时写入redis
func GetNetworkAllPowerChart() []*schema.PowerChart {
	var powerCharts = make([]*schema.PowerChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "network_power_all_chart", &powerCharts)
	if err != nil {
		logger.Errorf("GetNetwork7dPowerChart GetValue(network_power_all_chart) | %v", err)
	} else if exist {
		return powerCharts
	}

	end := util2.GetPoint(time.Now().Unix())
	startTime, err := GetFirstHeightTimeInDb()
	if err != nil {
		logger.Errorf("GetNetworkAllPowerChart GetFirstHeightInDb |%v", err)
		return powerCharts
	}
	start := util2.GetNullPoint(startTime)

	for i := start; i < end; i += util.DaySeconds {
		power := GetNetworkPower(i, i+util.DaySeconds)
		powerCharts = append(powerCharts, &schema.PowerChart{
			Date:  time.Unix(i, 0).Format(util.GoLangTimeFormat),
			Power: power,
		})
	}

	cache.Redis.SetValue(context.TODO(), "network_power_all_chart", powerCharts, time.Hour*4)
	return powerCharts
}

// 获取指定范围全网算力
func GetNetworkPower(start, end int64) float64 {
	var solution = &model.SolutionInDb{}
	var data []struct {
		TotalProofTarget float64 `bson:"TotalProofTarget"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gte": start, "$lt": end}}}
	//var group = bson.M{"$group": bson.M{"_id": "", "TotalProofTarget": bson.M{"$sum": bson.M{"$multiply": []string{"$sn", "$pt"}}}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalProofTarget": bson.M{"$sum": "$pt"}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), solution.TableName(), query, &data); err != nil {
		logger.Errorf("GetNetworkPower Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	timeRange := float64(end - start)
	hashRate := data[0].TotalProofTarget / timeRange

	res, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", hashRate), 64)
	return res
}

//获取上一个整点，往前推七天，每隔六小时取一次点

// 获取指定时间开始之前六小时的全网平均算力
func GetNeworkPowerSixHourAgo(end int64) float64 {
	timeRange := util.HourSeconds * 6

	var block = &model.BlockInDb{}
	var data []struct {
		TotalProofTarget float64 `bson:"TotalProofTarget"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gte": end - int64(timeRange), "$lt": end}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalProofTarget": bson.M{"$sum": bson.M{"$multiply": []string{"$sn", "$pt"}}}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetNetworkPower Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	hashRate := data[0].TotalProofTarget / float64(timeRange)

	res, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", hashRate), 64)
	return res
}

func GetNetworkPowerLastTime(start int64) float64 {
	blk, err := GetLatestBlockInDb()
	if err != nil {
		logger.Errorf("GetNetworkPowerLastTime GetLatestBlockInDb | %v", err)
		return 0
	}
	if start > blk.Timestamp {
		logger.Errorf("GetNetworkPowerLastTime start(%d) is out of range (%d)", start, blk.Timestamp)
		return 0
	}

	timeRange := blk.Timestamp - start

	var block = &model.BlockInDb{}
	var data []struct {
		TotalProofTarget float64 `bson:"TotalProofTarget"`
	}
	var match = bson.M{"$match": bson.M{"tp": bson.M{"$gte": start}}}
	var group = bson.M{"$group": bson.M{"_id": "", "TotalProofTarget": bson.M{"$sum": bson.M{"$multiply": []string{"$sn", "$pt"}}}}}
	var query = []bson.M{match, group}
	if err := mongodb.Aggregate(context.TODO(), block.TableName(), query, &data); err != nil {
		logger.Errorf("GetNetworkPowerLastTime Aggregate | %v", err)
		return 0
	}

	if len(data) == 0 {
		return 0
	}

	hashRate := data[0].TotalProofTarget / float64(timeRange)

	res, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", hashRate), 64)
	return res
}
