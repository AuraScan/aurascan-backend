package repository

import (
	"aurascan-backend/chain"
	"aurascan-backend/internal/config"
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"aurascan-backend/util"
	"ch-common-package/cache"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	util2 "ch-common-package/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

func GetNetworkOverview() *schema.Network {
	var networkOverview = &schema.Network{}

	var blockInRedis = &schema.BlockInRedis{}
	if exist, err := cache.Redis.GetValue(context.TODO(), "latest_block_info", &blockInRedis); err != nil {
		logger.Errorf("GetNetworkOverview GetValue(latest_block_info) | %v", err)
	} else if exist {
		networkOverview.BlockHeight = blockInRedis.BlockHeight
		networkOverview.LatestBlockTime = blockInRedis.LatestBlockTime
		networkOverview.CoinbaseTarget = blockInRedis.CoinbaseTarget
		networkOverview.ProofTarget = blockInRedis.ProofTarget
		networkOverview.Epoch = blockInRedis.BlockHeight / 360
		networkOverview.EstimatedNetworkSpeed = blockInRedis.EstimatedNetworkSpeed
	} else {
		block, err := chain.GetLatestBlock()
		if err != nil {
			logger.Errorf("GetNetworkOverview GetLatestBlock | %v", err)
		} else {
			metaData := block.Header.MetaData
			networkOverview.BlockHeight = metaData.Height
			networkOverview.LatestBlockTime = metaData.Timestamp
			networkOverview.CoinbaseTarget = metaData.CoinbaseTarget
			networkOverview.ProofTarget = metaData.ProofTarget
			networkOverview.Epoch = metaData.Height / 360
			speed := GetEstimatedNetworkSpeed15m()
			networkOverview.EstimatedNetworkSpeed = speed
			blockInRedis = &schema.BlockInRedis{
				BlockHeight:           metaData.Height,
				LatestBlockTime:       metaData.Timestamp,
				CoinbaseTarget:        metaData.CoinbaseTarget,
				ProofTarget:           metaData.ProofTarget,
				EstimatedNetworkSpeed: speed,
			}
			cache.Redis.SetValue(context.TODO(), "latest_block_info", blockInRedis, time.Second*30)
		}
	}

	res, err := cache.Redis.HGetAll(context.TODO(), "network_overview").Result()
	if err != nil {
		logger.Errorf("GetNetworkOverview HGetAll | %v", err)
	} else {
		for k, v := range res {
			if k == "total_network_Staking" {
				staking, err := strconv.ParseFloat(v, 64)
				if err != nil {
					logger.Errorf("GetNetworkOverview ParseFloat(%s) | %v", k, err)
				}
				networkOverview.NetworkStaking = staking
			}

			if k == "program_count" {
				programCount, err := strconv.Atoi(v)
				if err != nil {
					logger.Errorf("GetNetworkOverview Atoi(%s) | %v", k, err)
				}
				networkOverview.ProgramCount = programCount
			}

			if k == "miner_count" {
				minerCount, err := strconv.Atoi(v)
				if err != nil {
					logger.Errorf("GetNetworkOverview Atoi(%s) | %v", k, err)
				}
				networkOverview.NetworkMiners = minerCount
			}

			if k == "effective_proof_24h" {
				effectiveProof, err := strconv.ParseFloat(v, 64)
				if err != nil {
					logger.Errorf("GetNetworkOverview ParseFloat(%s) | %v", k, err)
				}
				networkOverview.NetworkEffectiveProof24h = effectiveProof
			}

			if k == "total_puzzle_reward" {
				puzzleReward, err := strconv.ParseFloat(v, 64)
				if err != nil {
					logger.Errorf("GetNetworkOverview ParseFloat(%s) | %v", k, err)
				}
				networkOverview.NetworkPuzzleReward = puzzleReward
			}

			if k == "committee_members" {
				committees, err := strconv.Atoi(v)
				if err != nil {
					logger.Errorf("GetNetworkOverview Atoi(%s) | %v", k, err)
				}
				networkOverview.NetworkValidators = committees
			}

			if k == "delegator_numbers" {
				delegators, err := strconv.Atoi(v)
				if err != nil {
					logger.Errorf("GetNetworkOverview Atoi(%s) | %v", k, err)
				}
				networkOverview.NetworkDelegators = delegators
			}
		}
	}
	return networkOverview
}

func GetEstimatedNetworkSpeed15m() float64 {
	start := time.Now().Unix() - util.QuarterSeconds
	table := (&model.SolutionInDb{}).TableName()
	match := bson.M{"$match": bson.M{"tp": bson.M{"$gte": start}}}
	group := bson.M{"$group": bson.M{"_id": "", "total_proof": bson.M{"$sum": "$pt"}}}
	query := []bson.M{match, group}

	var data []struct {
		TotalProof float64 `bson:"total_proof"`
	}

	if err := mongodb.Aggregate(context.TODO(), table, query, &data); err != nil {
		logger.Errorf("GetEffectiveProof24h Aggregate | %v", err)
		return 0
	}
	if len(data) > 0 {
		return data[0].TotalProof / util.QuarterSeconds
	}
	return 0
}

func GetAddrInfoByAddress(addr string) (*model.AddressInDb, error) {
	var td *model.AddressInDb
	_, err := mongodb.FindOne(context.TODO(), (&model.AddressInDb{}).TableName(), bson.M{"addr": addr}, &td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GetAddrDetailByAddress(addr string) *schema.AddrDetailResp {
	var td = &model.AddressInDb{}
	var addrRes = &schema.AddrDetailResp{}

	_, err := mongodb.FindOne(context.TODO(), (&model.AddressInDb{}).TableName(), bson.M{"addr": addr}, &td)
	if err != nil {
		logger.Errorf("GetAddrDetailByAddress FindOne(%s) | %v", addr, err)
		return addrRes
	}

	addrRes = &schema.AddrDetailResp{
		Addr:          addr,
		AddrType:      "unknown",
		PublicCredits: td.PublicCredits,
	}

	var bondPart []schema.AddrBondPart
	err = mongodb.Find(context.TODO(), (&model.AddressInDb{}).TableName(), bson.M{"bond_state.validator": addr}, map[string]int{"addr": 1,
		"bond_state": 1}, bson.D{{"bond_state.bond_amount", -1}}, 0, 0, &bondPart)
	if err != nil {
		logger.Errorf("GetAddrDetailByAddress Find | %v", err)
	} else {
		addrRes.BondDelegatorList = bondPart
	}

	if td.Addr != td.BondState.Validator {
		addrRes.BondValidatorState = &schema.ValidatorPartResp{
			Validator: td.BondState.Validator,
			Staked:    td.BondState.BondAmount,
			Earned:    0,
		}
	}

	switch td.AddrType {
	case 1:
		addrRes.AddrType = "validator"
		addrRes.CommitteeCreditsStake = td.DelegatedAmount
		addrRes.BondCreditsStake = td.BondState.BondAmount
		addrRes.DelegatorStake = td.DelegatedAmount - td.BondState.BondAmount
		// TODO 取缓存
		addrRes.Ratio = td.StakeRatio
		addrRes.TotalDelegators = 0
		// TODO 质押收益待计算
		addrRes.TotalCommitteeEarned = 0
		addrRes.TotalValidatorEarned = 0
		addrRes.TotalDelegatorEarned = 0
	case 2:
		addrRes.AddrType = "prover"
		addrRes.TotalPuzzleReward = td.TotalReward
		addrRes.TotalSolutionsFound = td.TotalSolutionsFound
		addrRes.Power1h = td.Power1h
		addrRes.Power24h = td.Power24h
		addrRes.Power7d = td.Power7d
	}

	return addrRes
}

func GetAddrRankByTimeRange(timeRange string) []*schema.ProverListResp {
	var proverList = make([]*schema.ProverListResp, 0)

	//TODO: 默认按算力排行，后面需要添加按收益排行
	var timeField string
	switch timeRange {
	case util.HourStr:
		timeField = "power_1h"
	case util.DayStr:
		timeField = "power_24h"
	case util.WeekStr:
		timeField = "power_7d"
	}

	var addrs []*model.AddressInDb
	if err := mongodb.Find(context.TODO(), (&model.AddressInDb{}).TableName(), bson.M{timeField: bson.M{"$gt": 0}},
		nil, bson.D{{timeField, -1}}, 0, 30, &addrs); err != nil {
		logger.Errorf("GetAddrRankByTimeRange Find | %v", err)
		return proverList
	}

	for index, v := range addrs {
		power, reward := 0.0, 0.0
		switch timeRange {
		case util.HourStr:
			power = v.Power1h
			reward = v.Reward1h
		case util.DayStr:
			power = v.Power24h
			reward = v.Reward24h
		case util.WeekStr:
			power = v.Power7d
			reward = v.Reward7d
		}

		proverList = append(proverList, &schema.ProverListResp{
			Rank:      index + 1,
			Address:   v.Addr,
			LastBlock: v.LatestBlock,
			Power:     power,
			Reward:    reward,
		})
	}

	return proverList
}

// 获取前十地址列表
func getTop10DailyPowerAddrsLastMonth(sortBy string) []string {
	start, end := util.GetOneMonthUTCTimeRange()
	match := bson.M{"$match": bson.M{"timestamp": bson.M{"$gte": start, "$lte": end}, "address": bson.M{"$ne": "all"}}}
	group := bson.M{"$group": bson.M{"_id": "$address", "total": bson.M{"$sum": "$" + sortBy}}}
	sort := bson.M{"$sort": bson.M{"total": -1}}
	limit := bson.M{"$limit": 10}
	query := []bson.M{match, group, sort, limit}

	var data []struct {
		Address string  `bson:"_id"`
		Total   float64 `bson:"total"`
	}
	if err := mongodb.Aggregate(context.TODO(), (&model.AddressSumDailyInDb{}).TableName(), query, &data); err != nil {
		logger.Errorf("getTop10DailyPowerAddrsLastMonth Aggregate | %v", err)
		return nil
	}
	if len(data) > 0 {
		var res = make([]string, 0)
		for _, v := range data {
			res = append(res, v.Address)
		}
		return res
	}
	return nil
}

// 获取指定地址的Power图表
func GetAddrPowerOneMonth(addr string) []*schema.PowerTimestampChart {
	start, end := util.GetOneMonthUTCTimeRange()
	filter := bson.M{"timestamp": bson.M{"$gte": start, "$lte": end}, "address": addr}
	var asdi []*model.AddressSumDailyInDb

	var ptc = make([]*schema.PowerTimestampChart, 0)

	if err := mongodb.Find(context.TODO(), (&model.AddressSumDailyInDb{}).TableName(), filter, nil,
		bson.D{{"timestamp", 1}}, 0, 0, &asdi); err != nil {
		logger.Errorf("GetAddrPowerOneMonth Find | %v", err)
		return ptc
	}

	var asdiMap = make(map[int64]*model.AddressSumDailyInDb)
	for _, v := range asdi {
		asdiMap[v.Timestamp] = v
	}

	genesis := util2.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	timeList := util.GetTimestampListByTimeRange(start, end, genesis)
	for _, timestamp := range timeList {
		asdi, ok := asdiMap[timestamp]
		if ok {
			ptc = append(ptc, &schema.PowerTimestampChart{
				Timestamp: asdi.Timestamp,
				Power:     asdi.Power,
			})
		} else {
			ptc = append(ptc, &schema.PowerTimestampChart{
				Timestamp: timestamp,
				Power:     0,
			})
		}
	}
	return ptc
}

// 获取前十地址Power图表
func GetTop10AddrDailyPowerChartOneMonth() []*schema.AddrPowerChart {
	var addrPowerChart = make([]*schema.AddrPowerChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "top_address_power_one_month", &addrPowerChart)
	if err != nil {
		logger.Errorf("GetTop10AddrDailyPowerChartOneMonth GetValue(top_address_power_one_month) | %v", err)
	} else if exist {
		return addrPowerChart
	}

	addrs := getTop10DailyPowerAddrsLastMonth("power")
	if len(addrs) > 0 {
		for _, v := range addrs {
			powerChart := GetAddrPowerOneMonth(v)
			addrPowerChart = append(addrPowerChart, &schema.AddrPowerChart{
				Address:     v,
				PowerCharts: powerChart,
			})
		}
	}
	cache.Redis.SetValue(context.TODO(), "top_address_power_one_month", addrPowerChart, time.Hour*6)
	return addrPowerChart
}

// 获取指定地址的Reward列表
func GetAddrRewardOneMonth(addr string) []*schema.RewardTimestampChart {
	start, end := util.GetOneMonthUTCTimeRange()
	filter := bson.M{"timestamp": bson.M{"$gte": start, "$lte": end}, "address": addr}
	var asdi []*model.AddressSumDailyInDb

	var ptc = make([]*schema.RewardTimestampChart, 0)

	if err := mongodb.Find(context.TODO(), (&model.AddressSumDailyInDb{}).TableName(), filter, nil,
		bson.D{{"timestamp", 1}}, 0, 0, &asdi); err != nil {
		logger.Errorf("GetAddrRewardOneMonth Find | %v", err)
		return ptc
	}

	var asdiMap = make(map[int64]*model.AddressSumDailyInDb)
	for _, v := range asdi {
		asdiMap[v.Timestamp] = v
	}

	genesis := util2.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	timeList := util.GetTimestampListByTimeRange(start, end, genesis)
	for _, timestamp := range timeList {
		asdi, ok := asdiMap[timestamp]
		if ok {
			ptc = append(ptc, &schema.RewardTimestampChart{
				Timestamp: asdi.Timestamp,
				Reward:    asdi.Reward,
			})
		} else {
			ptc = append(ptc, &schema.RewardTimestampChart{
				Timestamp: timestamp,
				Reward:    0,
			})
		}
	}
	return ptc
}

// 获取前十地址Reward图表
func GetTop10AddrDailyRewardChartOneMonth() []*schema.AddrRewardChart {
	var addrRewardChart = make([]*schema.AddrRewardChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "top_address_reward_one_month", &addrRewardChart)
	if err != nil {
		logger.Errorf("GetTop10AddrDailyRewardChartOneMonth GetValue(top_address_reward_one_month) | %v", err)
	} else if exist {
		return addrRewardChart
	}

	addrs := getTop10DailyPowerAddrsLastMonth("reward")
	if len(addrs) > 0 {
		for _, v := range addrs {
			rewardChart := GetAddrRewardOneMonth(v)
			addrRewardChart = append(addrRewardChart, &schema.AddrRewardChart{
				Address:      v,
				RewardCharts: rewardChart,
			})
		}
	}
	cache.Redis.SetValue(context.TODO(), "top_address_reward_one_month", addrRewardChart, time.Hour*6)
	return addrRewardChart
}

// 获取指定地址的Solution数量列表
func GetAddrSolutionOneMonth(addr string) []*schema.SolutionsTimestampChart {
	start, end := util.GetOneMonthUTCTimeRange()
	filter := bson.M{"timestamp": bson.M{"$gte": start, "$lte": end}, "address": addr}
	var asdi []*model.AddressSumDailyInDb

	var ptc = make([]*schema.SolutionsTimestampChart, 0)

	if err := mongodb.Find(context.TODO(), (&model.AddressSumDailyInDb{}).TableName(), filter, nil,
		bson.D{{"timestamp", 1}}, 0, 0, &asdi); err != nil {
		logger.Errorf("GetAddrSolutionOneMonth Find | %v", err)
		return ptc
	}

	var asdiMap = make(map[int64]*model.AddressSumDailyInDb)
	for _, v := range asdi {
		asdiMap[v.Timestamp] = v
	}

	genesis := util2.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	timeList := util.GetTimestampListByTimeRange(start, end, genesis)
	for _, timestamp := range timeList {
		asdi, ok := asdiMap[timestamp]
		if ok {
			ptc = append(ptc, &schema.SolutionsTimestampChart{
				Timestamp: asdi.Timestamp,
				Count:     asdi.Solutions,
			})
		} else {
			ptc = append(ptc, &schema.SolutionsTimestampChart{
				Timestamp: timestamp,
				Count:     0,
			})
		}
	}
	return ptc
}

// 获取前十地址Solution数量图表
func GetTop10AddrDailySolutionsChartOneMonth() []*schema.AddrSolutionsChart {
	var addrSolutionsChart = make([]*schema.AddrSolutionsChart, 0)
	exist, err := cache.Redis.GetValue(context.TODO(), "top_address_solution_one_month", &addrSolutionsChart)
	if err != nil {
		logger.Errorf("GetTop10AddrDailyPowerChartOneMonth GetValue(top_address_power_one_month) | %v", err)
	} else if exist {
		return addrSolutionsChart
	}

	addrs := getTop10DailyPowerAddrsLastMonth("solutions")
	if len(addrs) > 0 {
		for _, v := range addrs {
			solutionChart := GetAddrSolutionOneMonth(v)
			addrSolutionsChart = append(addrSolutionsChart, &schema.AddrSolutionsChart{
				Address:         v,
				SolutionsCharts: solutionChart,
			})
		}
	}
	cache.Redis.SetValue(context.TODO(), "top_address_solution_one_month", addrSolutionsChart, time.Hour*6)
	return addrSolutionsChart
}

func GetRewardChartOneMonth() []*schema.RewardTimestampChart {
	start := util.GetOneMonthAgoUTCTimeStart()
	filter := bson.M{"timestamp": bson.M{"$gte": start}, "address": "all"}
	var asdi []*model.AddressSumDailyInDb

	var ptc = make([]*schema.RewardTimestampChart, 0)

	exist, err := cache.Redis.GetValue(context.TODO(), "puzzle_reward_month", &ptc)
	if err != nil {
		logger.Errorf("GetRewardChartOneMonth GetValue(puzzle_reward_month) | %v", err)
	} else if exist {
		return ptc
	}

	if err := mongodb.Find(context.TODO(), (&model.AddressSumDailyInDb{}).TableName(), filter, nil,
		bson.D{{"timestamp", 1}}, 0, 0, &asdi); err != nil {
		logger.Errorf("GetRewardChartOneMonth Find | %v", err)
		return ptc
	}

	var asdiMap = make(map[int64]*model.AddressSumDailyInDb)
	for _, v := range asdi {
		asdiMap[v.Timestamp] = v
	}

	genesis := util2.GetNullPoint(config.Global.Chain.GenesisTimestamp)
	timeList := util.GetTimestampListByStart(start, genesis)
	for _, timestamp := range timeList {
		asdi, ok := asdiMap[timestamp]
		if ok {
			ptc = append(ptc, &schema.RewardTimestampChart{
				Timestamp: asdi.Timestamp,
				Reward:    asdi.Reward,
			})
		} else {
			ptc = append(ptc, &schema.RewardTimestampChart{
				Timestamp: timestamp,
				Reward:    0,
			})
		}
	}

	cache.Redis.SetValue(context.TODO(), "puzzle_reward_month", ptc, time.Hour*2)
	return ptc
}
