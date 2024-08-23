package repository

import (
	"strconv"
	"strings"
)

// -1、未找到
// 1、通过高度查询区块 2、通过block hash查询区块 3、通过transaction id查询transaction
// 4、通过transition_id查询transition 5、通过地址查询prover/validator 6、通过包含字段查询program
func SearchByField(field string) (searchType int) {
	if num, err := strconv.ParseInt(field, 10, 64); err == nil {
		if exist, err := JudgeBlockInDb(num); err == nil && exist {
			return 1
		}
	} else if strings.HasPrefix(field, "ab1") {
		if tn, err := GetBlockByHash(field); tn != nil && err == nil {
			return 2
		}
	} else if strings.HasPrefix(field, "at1") {
		if tn, err := GetTransactionById(field); tn != nil && err == nil {
			return 3
		}
	} else if strings.HasPrefix(field, "au1") {
		if tn, err := GetTransitionById(field); tn != nil && err == nil {
			return 4
		}
	} else if strings.HasPrefix(field, "aleo1") {
		if tn, err := GetAddrInfoByAddress(field); tn != nil && err == nil {
			return 5
		}
	} else {
		if tn, err := GetProgramByRegexId(field); tn != nil && err == nil {
			return 6
		}
	}
	return -1
}
