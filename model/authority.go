package model

import (
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// Subdag Details
type AuthorityInDb struct {
	Type                   string      `bson:"te"`
	Round                  int64       `bson:"rd"`  // 区块高度和Round做唯一索引
	Idx                    int         `bson:"idx"` // 同一Round下的的索引
	Height                 int64       `bson:"ht"`  // 绑定指定区块，需索引
	CertificateId          string      `bson:"cd" json:"certificate_id"`
	BatchId                string      `bson:"bd"` //out
	Author                 string      `bson:"ar"`
	Timestamp              int64       `bson:"tp"`
	TransmissionIds        []string    `bson:"ts"`
	PreviousCertificateIds []string    `bson:"ps"`
	Signature              string      `bson:"se"`
	Signatures             interface{} `bson:"ss"`
}

func (*AuthorityInDb) TableName() string {
	return "authority"
}

func (t *AuthorityInDb) Save() error {
	if err := mongodb.InsertOne(context.TODO(), t.TableName(), t); err != nil {
		return fmt.Errorf("AuthorityInDb.Save InsertOne | %v", err)
	}
	return nil
}

// 通过高度获取
// TODO:高度为0的数据待处理
func GetBlockAuthoritys(height int64) []*AuthorityInDb {
	var authoritys = make([]*AuthorityInDb, 0)
	if err := mongodb.Find(context.TODO(), (&AuthorityInDb{}).TableName(), bson.M{"ht": height},
		nil, bson.D{{"rd", -1}, {"idx", 1}}, 0, 0, &authoritys); err != nil {
		logger.Errorf("GetBlockAuthoritys Find(%d) | %v", height, err)
	}
	return authoritys
}

func SaveAuthorityList(ts []*AuthorityInDb) error {
	if len(ts) == 0 {
		return nil
	}

	var authoritys []interface{}
	for _, m := range ts {
		authoritys = append(authoritys, m)
	}

	if err := mongodb.InsertMany(context.TODO(), (&AuthorityInDb{}).TableName(), authoritys); err != nil {
		return fmt.Errorf("SaveAuthorityList InsertMany | %v", err)
	}
	return nil
}

func JudgeAuthorityInDb(height int64) bool {
	var block *AuthorityInDb
	found, err := mongodb.FindOne(context.TODO(), (&AuthorityInDb{}).TableName(), bson.M{"ht": height}, &block)
	if err != nil {
		return false
	}

	if found {
		return true
	}
	return false
}

func DeleteAuthorityByHeight(height int64) error {
	_, err := mongodb.DeleteMany(context.TODO(), (&AuthorityInDb{}).TableName(), bson.M{"ht": height})
	if err != nil {
		return fmt.Errorf("DeleteAuthorityByHeight DeleteMany(height=%d) | %v", height, err)
	}
	return nil
}
