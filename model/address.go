package model

import (
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type AddressInDb struct {
	Addr string `json:"addr" bson:"addr"`
	//Rank                int                   `json:"rank" bson:"rk"`
	AddrType      int     `json:"addr_type" bson:"addr_type"` //0表示未知，1表示验证者，2表示普通地址（不需要判定是delegator还是prover）
	PublicCredits float64 `json:"public_credits" bson:"public_credits"`

	//作为验证者/委托者的基本信息
	DelegatedAmount float64        `json:"delegated_amount" bson:"delegated_amount"`
	WithdrawalAddr  string         `json:"withdrawal_addr" bson:"withdrawal_addr"`
	CommitteeState  CommitteeState `json:"committee_state" bson:"committee_state"`
	BondState       BondedState    `json:"bond_state" bson:"bond_state"`       //绑定的验证者, 可能是自身，如果不是自身则
	UnBondState     UnBondState    `json:"unbond_state" bson:"unbond_state"`   //解绑状态
	InitialStake    float64        `json:"initial_stake" bson:"initial_stake"` //初始质押
	StakeRatio      float64        `json:"stake_ratio" bson:"stake_ratio"`     //质押占比
	Vote            float64        `json:"vote" bson:"vote"`                   //得票率

	//作为prover的出快信息
	LatestBlock         int64   `json:"latest_block" bson:"latest_block"`             //最近出块
	TotalReward         float64 `json:"total_reward" bson:"total_reward"`             //出块收益
	TotalSolutionsFound int64   `json:"solutions_found" bson:"total_solutions_found"` //出块数
	Power1h             float64 `json:"power_1h" bson:"power_1h"`                     //1小时平均算力
	Power24h            float64 `json:"power_24h" bson:"power_24h"`                   //24小时平均算力
	Power7d             float64 `json:"power_7d" bson:"power_7d"`                     //7天平均算力
	Reward1h            float64 `json:"reward_1h" bson:"reward_1h"`                   //1小时收益
	Reward24h           float64 `json:"reward_24h" bson:"reward_24h"`                 //24小时收益
	Reward7d            float64 `json:"reward_7d" bson:"reward_7d"`                   //7天收益

	//创建高度
	CreateHeight int64 `json:"create_Height" bson:"create_height"`
	//更新高度
	UpdateHeight int64 `json:"update_Height" bson:"update_height"`
}

type CommitteeState struct {
	IsOpen     string  `json:"is_open" bson:"is_open"`
	Commission float64 `json:"commission" bson:"commission"` //1~100, 验证者保留的奖励的百分比
}

type BondedState struct {
	Validator  string  `json:"validator" bson:"validator"`
	BondAmount float64 `json:"bond_amount" bson:"bond_amount"`
}

type UnBondState struct {
	UnBondingAmount float64 `json:"unbonding_amount" bson:"unbonding_amount"`
	UnBondHeight    int64   `json:"unbond_height" bson:"unbond_height"`
}

type ValidatorDetailRes struct {
	Address             string  `json:"address"`
	TotalPower          float64 `json:"total_power"`
	TotalSolutionsFound int     `json:"total_solutions_found"`
}

type ValidatorListRes struct {
	Address          string  `json:"address"`
	PublicCredits    float64 `json:"public_credits"`
	CommitteeCredits float64 `json:"committee_credits"`
	BondCredits      float64 `json:"bond_credits"`
	Ratio            float64 `json:"ratio"`
	IsOpen           int     `json:"is_open"`
	CreateHeight     int64   `json:"create_height"`
}

func (*AddressInDb) TableName() string {
	return "address"
}

//func GetBlocksByPage(page, pageSize int) ([]*schema.BlockListResp, int64) {
//	var blocks []*model.BlockInDb
//	var blocksResp = make([]*schema.BlockListResp, 0)
//
//	height, err := GetLatestHeightInDb()
//	if err != nil {
//		logger.Errorf("GetBlocksByPage GetLatestHeightInDb | %v", err)
//		return blocksResp, 0
//	}
//
//	if err := mongodb.Find(context.TODO(), (&model.BlockInDb{}).TableName(), bson.M{}, nil, bson.D{{"ht", -1}}, util.Offset(pageSize, page), int64(pageSize), &blocks); err != nil {
//		logger.Errorf("GetBlocksByPage Find | %v", err)
//		return blocksResp, 0
//	}

func GetValidatorList(page, pageSize int) ([]*ValidatorListRes, int64) {
	var addrs []*AddressInDb
	var validatorList = make([]*ValidatorListRes, 0)

	tableName := (&AddressInDb{}).TableName()

	count, err := mongodb.Count(context.TODO(), tableName, bson.M{"addr_type": 1})
	if err != nil {
		logger.Errorf("GetValidatorList Count(%d %d) | %v", page, pageSize, err)
		return validatorList, 0
	}

	if err := mongodb.Find(context.TODO(), tableName, bson.M{"addr_type": 1}, nil, bson.D{{"bond_credits", -1}}, util.Offset(pageSize, page), int64(pageSize), &addrs); err != nil {
		logger.Errorf("GetValidatorList Find(%d %d) | %v", page, pageSize, err)
		return validatorList, 0
	}

	for _, v := range addrs {
		isOpen := 0
		if v.CommitteeState.IsOpen == "true" {
			isOpen = 1
		}
		validator := &ValidatorListRes{
			Address:          v.Addr,
			PublicCredits:    v.PublicCredits,
			CommitteeCredits: v.DelegatedAmount,
			BondCredits:      v.BondState.BondAmount,
			Ratio:            v.StakeRatio,
			IsOpen:           isOpen,
			CreateHeight:     v.CreateHeight,
		}
		validatorList = append(validatorList, validator)
	}

	return validatorList, count
}

//
