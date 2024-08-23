package repository

import (
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/mongodb"
	"ch-common-package/ssdb"
	util2 "ch-common-package/util"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"regexp"
	"strings"
	"time"
)

// 获取交易列表
func GetProgramsByPage(page, pageSize int) ([]*schema.ProgramListResp, int64) {
	var programs []*model.ProgramInDb
	var programsResp = make([]*schema.ProgramListResp, 0)

	count, err := mongodb.Count(context.TODO(), (&model.ProgramInDb{}).TableName(), bson.M{})
	if err != nil {
		logger.Errorf("GetProgramsByPage Count | %v", err)
		return programsResp, 0
	}

	if err := mongodb.Find(context.TODO(), (&model.ProgramInDb{}).TableName(), bson.M{}, nil, bson.D{{"tc", -1}}, util.Offset(pageSize, page), int64(pageSize), &programs); err != nil {
		logger.Errorf("GetProgramsByPage Find | %v", err)
		return programsResp, 0
	}

	if len(programs) > 0 {
		for _, v := range programs {
			var program = &schema.ProgramListResp{
				ProgramId:     v.ProgramID,
				TransactionId: v.TransactionID,
				Height:        v.Height,
				Time:          time.Unix(v.DeployTime, 0).Format(util.GoLangTimeFormat),
				Timestamp:     v.DeployTime,
				TimesCalled:   v.TimesCalled,
			}
			programsResp = append(programsResp, program)
		}
	}

	return programsResp, count
}

func GetProgramByRegexId(id string) (*model.ProgramInDb, error) {
	var td *model.ProgramInDb
	_, err := mongodb.FindOne(context.TODO(), (&model.ProgramInDb{}).TableName(), bson.M{"pd": bson.M{"$regex": id}}, &td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GetProgramDetail(programId string) (*schema.ProgramDetailResp, error) {
	var programInDb *model.ProgramInDb
	filter := bson.M{"pd": programId}
	exist, err := mongodb.FindOne(context.TODO(), (&model.ProgramInDb{}).TableName(), filter, &programInDb)
	if err != nil {
		logger.Errorf("GetProgramDetail FindOne(%s) | %v", programId, err)
		return nil, err
	}
	if exist {
		var programDetail = &schema.ProgramDetailResp{
			ProgramId:         programInDb.ProgramID,
			Owner:             programInDb.Owner,
			DeployHeight:      programInDb.Height,
			DeployTime:        time.Unix(programInDb.DeployTime, 0).Format(util.GoLangTimeFormat),
			DeployTimestamp:   programInDb.DeployTime,
			DeployTransaction: programInDb.TransactionID,
			TimesCalled:       programInDb.TimesCalled,
		}
		txn, err := GetTransactionById(programInDb.TransactionID)
		if err != nil {
			logger.Errorf("GetProgramDetail(%s) GetTransactionById(%s) | %v", programId, programInDb.TransactionID, err)
		}
		if txn != nil {
			programDetail.DeployFee = txn.Fee
		}
		return programDetail, nil
	}
	return &schema.ProgramDetailResp{}, nil
}

// 获取program每日被调用次数
func GetSingleProgramCallingChart(programId string) ([]*schema.ProgramCallingCount, error) {
	start := util.GetOneMonthAgoTime()

	var programsInDb []*model.ProgramCountInDb
	if err := mongodb.Find(context.TODO(), (&model.ProgramCountInDb{}).TableName(),
		bson.M{"pm": programId, "tp": bson.M{"$gte": start}}, nil, bson.D{{"tp", 1}}, 0, 0, &programsInDb); err != nil {
		return nil, fmt.Errorf("find program(%s) | %v", programId, err)
	}

	var programsCallingCount = make([]*schema.ProgramCallingCount, 0)

	var programList = make(map[int64]float64)

	if len(programsInDb) > 0 {
		for _, v := range programsInDb {
			programList[v.Timestamp] = v.TimesCalled
		}
	}

	timeList := util2.Generate30DTime(start)
	for _, v := range timeList {
		count, ok := programList[v]
		if !ok {
			programsCallingCount = append(programsCallingCount, &schema.ProgramCallingCount{
				Timestamp: v,
				Value:     0,
			})
		} else {
			programsCallingCount = append(programsCallingCount, &schema.ProgramCallingCount{
				Timestamp: v,
				Value:     int(count),
			})
		}
	}
	return programsCallingCount, nil

}

// 获取program源码
func GetProgramSourceCode(programId string) string {
	res, err := ssdb.Client.Get(programId)
	if err != nil {
		logger.Errorf("GetProgramSourceCode(%s) | %v", programId, err)
		return ""
	}
	return res
}

func GetFunctionFromSourceCode(sourceCode string) []schema.ProgramFunction {
	var res = make([]schema.ProgramFunction, 0)

	// 正则表达式匹配以 function 开头的函数名称
	funcRe := regexp.MustCompile(`(?m)^function ([a-zA-Z0-9_]+):`)
	// 正则表达式匹配 input 行并提取类型
	inputRe := regexp.MustCompile(`(?m)^\s*input [a-zA-Z0-9_]+ as ([a-zA-Z0-9_.]+)/?([a-zA-Z0-9_.]*)\.([a-zA-Z0-9_.]+);`)
	// 正则表达式匹配所有关键字
	keywordRe := regexp.MustCompile(`(?m)^(program|import|function|closure|struct|record|mapping|finalize) `)

	// 查找所有以 function 开头的函数名称
	funcMatches := funcRe.FindAllStringSubmatchIndex(sourceCode, -1)

	// 解析每个函数块
	for _, match := range funcMatches {
		start := match[0]
		end := len(sourceCode)

		// 查找下一个关键字的位置
		keywordMatches := keywordRe.FindAllStringSubmatchIndex(sourceCode[start+1:], -1)
		for _, km := range keywordMatches {
			keywordStart := start + 1 + km[0]
			if keywordStart > start {
				end = keywordStart
				break
			}
		}

		funcBlock := sourceCode[start:end]

		// 获取函数名称
		funcNameMatch := funcRe.FindStringSubmatch(funcBlock)
		if len(funcNameMatch) > 1 {
			var programFunction = schema.ProgramFunction{}
			programFunction.Name = funcNameMatch[1]

			// 查找输入行
			inputMatches := inputRe.FindAllStringSubmatch(funcBlock, -1)
			var inputTypes []string
			for _, inputMatch := range inputMatches {
				if len(inputMatch) > 3 {
					fullType := inputMatch[1]
					if inputMatch[2] != "" {
						fullType += "/" + inputMatch[2]
					}
					fullType += "." + inputMatch[3]
					inputTypes = append(inputTypes, fullType)
				}
			}

			// 输出 input 类型
			if len(inputTypes) > 0 {
				programFunction.Inputs = inputTypes
			}
			res = append(res, programFunction)
		}
	}
	return res
}

func RemoveSubstring(s, sep string) string {
	if idx := strings.Index(s, sep); idx != -1 {
		return s[:idx]
	}
	return s
}
