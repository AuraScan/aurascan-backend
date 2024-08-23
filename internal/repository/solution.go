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

func GetSolutions() int64 {
	var solution = &model.SolutionInDb{}
	total, err := mongodb.Count(context.TODO(), solution.TableName(), bson.M{})
	if err != nil {
		logger.Errorf("GetSolutions Count | %v", err)
		return 0
	}
	return total
}

func GetSolutionListByAddr(addr string, page, pageSize int) ([]*schema.SolutionInAddrResp, int64) {
	var solutions []*model.SolutionInDb
	var solutionsResp = make([]*schema.SolutionInAddrResp, 0)

	total, err := mongodb.Count(context.TODO(), (&model.SolutionInDb{}).TableName(), bson.M{"as": addr})
	if err != nil {
		logger.Errorf("GetSolutionListByAddr Count | %v", err)
		return solutionsResp, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.SolutionInDb{}).TableName(), bson.M{"as": addr},
		nil, bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &solutions); err != nil {
		logger.Errorf("GetSolutionListByAddr Find(%d) | %v", addr, err)
		return solutionsResp, 0
	}
	if len(solutions) > 0 {
		for _, v := range solutions {
			var solution = &schema.SolutionInAddrResp{
				BlockHeight: v.Height,
				SolutionId:  v.SolutionId,
				Time:        time.Unix(v.Timestamp, 0).Format(util.GoLangTimeFormat),
				Target:      v.Target,
				Reward:      v.Reward,
			}
			solutionsResp = append(solutionsResp, solution)
		}
	}

	return solutionsResp, total
}

func GetSolutionsByAddr(addr string) int64 {
	var solution = &model.SolutionInDb{}
	total, err := mongodb.Count(context.TODO(), solution.TableName(), bson.M{"as": addr})
	if err != nil {
		logger.Errorf("GetSolutions Count | %v", err)
		return 0
	}
	return total
}

type CommitmentDetail struct {
	Commitment string
	Height     int64
	Epoch      int64
	PushTime   string
	Target     float64
	Rewards    float64
}

func GetCommitmentListByAddr(addr string, page, pageSize int) ([]*CommitmentDetail, int64, error) {
	var solutions []*model.SolutionInDb
	query := bson.M{"as": addr}
	err := mongodb.Find(context.TODO(), (&model.SolutionInDb{}).TableName(), query, nil, bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &solutions)
	if err != nil {
		return nil, 0, err
	}
	var commits []*CommitmentDetail
	for _, v := range solutions {
		commits = append(commits, &CommitmentDetail{
			Commitment: "",
			Height:     v.Height,
			Epoch:      v.Height / 360,
			PushTime:   time.Unix(time.Now().Unix(), 0).Format(util.GoLangTimeFormat),
			Target:     0,
			Rewards:    0,
		})
	}

	count, err := mongodb.Count(context.TODO(), (&model.SolutionInDb{}).TableName(), query)
	if err != nil {
		return nil, 0, err
	}

	return commits, count, nil
}

func GetSolutionCountByAddress(addr string) int64 {
	filter := bson.M{"as": addr}
	count, err := mongodb.Count(context.TODO(), (&model.SolutionInDb{}).TableName(), filter)
	if err != nil {
		logger.Errorf("GetSolutionCountByAddress Count(%s) | %v", addr, err)
		return 0
	}
	return count
}

// 传入起始高度，计算指定节点指定日期累计收益
func CalculateAddrDayReward(start int64, addr string) float64 {
	end := start + util.DaySeconds
	match := bson.M{"$match": bson.M{"as": addr, "tp": bson.M{"$gte": start, "$lt": end}}}
	group := bson.M{"$group": bson.M{"_id": "", "dayReward": bson.M{"$sum": "$rd"}}}
	query := []bson.M{match, group}
	var data []struct {
		Id        string  `bson:"_id"`
		DayReward float64 `bson:"dayReward"`
	}
	if err := mongodb.Aggregate(context.TODO(), (&model.SolutionInDb{}).TableName(), query, &data); err != nil {
		logger.Errorf("CalculateAddrDayReward date(%s) addr(%s) | %v", time.Unix(start, 0).Format(util.GolangDateFormat), addr, err)
		return 0
	}
	if len(data) > 0 {
		return data[0].DayReward
	}
	return 0
}

func GetSolutionListByHeight(height int64, page, pageSize int) ([]*schema.SolutionInBlockResp, int64, float64) {
	var solutions []*model.SolutionInDb
	var solutionsResp = make([]*schema.SolutionInBlockResp, 0)
	var targetTotal float64

	total, err := mongodb.Count(context.TODO(), (&model.SolutionInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		logger.Errorf("GetSolutionListByHeight Count | %v", err)
		return solutionsResp, 0, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.SolutionInDb{}).TableName(), bson.M{"ht": height},
		nil, bson.D{{"rd", -1}}, util.Offset(pageSize, page), int64(pageSize), &solutions); err != nil {
		logger.Errorf("GetSolutionListByHeight Find(%d) | %v", height, err)
		return solutionsResp, 0, 0
	}
	if len(solutions) > 0 {
		for _, v := range solutions {
			var solution = &schema.SolutionInBlockResp{
				Address:    v.Address,
				Commitment: "",
				Target:     v.Target,
				Reward:     v.Reward,
			}
			targetTotal += v.Target
			solutionsResp = append(solutionsResp, solution)
		}
	}

	return solutionsResp, total, targetTotal
}
