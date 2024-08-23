package explorer

import (
	"aurascan-backend/chain"
	"aurascan-backend/internal/repository"
	"aurascan-backend/model/schema"
	"ch-common-package/cache"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func GetProgramList(c *gin.Context) {
	var listInfo schema.PageListReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetProgramList BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	programs, count := repository.GetProgramsByPage(page, size)
	var data = struct {
		Blocks []*schema.ProgramListResp `json:"programs"`
		Count  int64                     `json:"count"`
	}{programs, count}

	ginx.ResSuccess(c, data)
}

func GetProgramDetail(c *gin.Context) {
	programId := ginx.Param(c, "program_id")
	programDetail, err := repository.GetProgramDetail(programId)
	if err != nil {
		ginx.ResFailed(c, "program not find")
		return
	}

	ginx.ResSuccess(c, programDetail)
}

func GetProgramChartById(c *gin.Context) {
	programId := ginx.Param(c, "program_id")
	programChart, err := repository.GetSingleProgramCallingChart(programId)
	if err != nil {
		ginx.ResFailed(c, "unknown error")
		return
	}
	ginx.ResSuccess(c, programChart)
}

func Get24hTopProgram(c *gin.Context) {
	var n []*schema.ProgramCalledChartResp

	exist, err := cache.Redis.GetValue(context.TODO(), "program_called_24h_top_10", &n)
	if err != nil {
		logger.Errorf("Get24hTopProgram Redis GetValue | %v", err)
	}
	if exist {
		ginx.ResSuccess(c, n)
		return
	}

	programChart := repository.GetTopProgramCalledIn24h()
	cache.Redis.SetValue(context.TODO(), "program_called_24h_top_10", programChart, 15*time.Minute)

	ginx.ResSuccess(c, programChart)
}

func GetOneMonthTopProgram(c *gin.Context) {
	var n []*schema.ProgramCountMonthChart

	exist, err := cache.Redis.GetValue(context.TODO(), "program_called_1month_top_10", &n)
	if err != nil {
		logger.Errorf("GetOneMonthTopProgram Redis GetValue | %v", err)
	}
	if exist {
		ginx.ResSuccess(c, n)
		return
	}

	programChart := repository.CalculateProgramCalledChartOneMonth()
	cache.Redis.SetValue(context.TODO(), "program_called_1month_top_10", programChart, 1*time.Hour)

	ginx.ResSuccess(c, programChart)
}

func GetMappingNameListByProgramId(c *gin.Context) {
	programId := ginx.Param(c, "program_id")
	names := chain.GetProgramNameById(programId)
	ginx.ResSuccess(c, names)
}

func GetMappingValue(c *gin.Context) {
	var listInfo schema.MappingInfo
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetMappingValue BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}
	if listInfo.MappingKey == "" || listInfo.MappingName == "" || listInfo.ProgramId == "" {
		logger.Error("GetMappingValue missing info | %v", listInfo)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	value := chain.GetProgramMapValue(listInfo.ProgramId, listInfo.MappingName, listInfo.MappingKey)

	ginx.ResSuccess(c, value)
}

func GetMappingSourceCode(c *gin.Context) {
	programId := ginx.Param(c, "program_id")
	sourceCode := repository.GetProgramSourceCode(programId)
	functions := repository.GetFunctionFromSourceCode(sourceCode)
	var res = schema.ProgramSource{
		SourceCode:      sourceCode,
		ProgramFunction: functions,
	}
	ginx.ResSuccess(c, res)
}
