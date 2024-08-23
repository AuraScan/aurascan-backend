package model

import (
	"ch-common-package/mongodb"
	"context"
	"fmt"
)

// 储存十五分钟的算力信息
type PoolInfoInDb struct {
	Date            string  `bson:"date"`             //时刻
	SelfSpeed       float64 `bson:"self_speed"`       //自身算力
	CommissionSpeed float64 `bson:"commission_speed"` //抽成的算力
	TotalSpeed      float64 `bson:"total_speed"`      //总算力
	IsPoint         int     `bson:"is_point"`         //是否整点
	Multiple        float64 `bson:"multiple"`         //调整的比率
	Timestamp       int64   `bson:"timestamp"`
	CreateAt        int64   `bson:"create_at"`
	UpdateAt        int64   `bson:"update_at"`
}

func (*PoolInfoInDb) TableName() string {
	return "pool_info"
}

func (*PoolInfoInDb) TableNameHourly() string {
	return "pool_info_hour"
}

func InsertPoolInfoHourly(p *PoolInfoInDb) error {
	if err := mongodb.InsertOne(context.TODO(), p.TableNameHourly(), p); err != nil {
		return err
	}
	return nil
}

func InsertPoolInfo(p *PoolInfoInDb) error {
	if err := mongodb.InsertOne(context.TODO(), p.TableName(), p); err != nil {
		return err
	}
	return nil
}

func InsertManyPoolInfo(p []*PoolInfoInDb) error {
	if len(p) == 0 {
		return nil
	}

	var poolInfos []interface{}
	for _, m := range p {
		poolInfos = append(poolInfos, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&PoolInfoInDb{}).TableName(), poolInfos); err != nil {
		return fmt.Errorf("InsertManyPoolInfo InsertMany | %v", err)
	}
	return nil
}
