package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
)

func GetValidatorList(c *gin.Context) {
	var listInfo schema.PageListReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetValidatorList BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	transitions, count := model.GetValidatorList(page, size)
	var data = struct {
		Validators []*model.ValidatorListRes `json:"validators"`
		Count      int64                     `json:"count"`
	}{transitions, count}

	ginx.ResSuccess(c, data)
}

func GetProverList(c *gin.Context) {
	timeRange := ginx.Param(c, "time_range")
	addrs := repository.GetAddrRankByTimeRange(timeRange)
	ginx.ResSuccess(c, addrs)
}

func GetStakeChart(c *gin.Context) {
	data := repository.GetStakeChartOneMonth()
	ginx.ResSuccess(c, data)
}

func GetTopAddrPowerChart(c *gin.Context) {
	addrPowerChart := repository.GetTop10AddrDailyPowerChartOneMonth()
	ginx.ResSuccess(c, addrPowerChart)
}

func GetAddrSolutionsChart(c *gin.Context) {
	addr := ginx.Param(c, "addr")
	addrSolutionsChart := repository.GetAddrSolutionOneMonth(addr)
	ginx.ResSuccess(c, addrSolutionsChart)
}

func GetAddrRewardChart(c *gin.Context) {
	addr := ginx.Param(c, "addr")
	addrRewardChart := repository.GetAddrRewardOneMonth(addr)
	ginx.ResSuccess(c, addrRewardChart)
}

func GetAddrPowerChart(c *gin.Context) {
	addr := ginx.Param(c, "addr")
	addrPowerChart := repository.GetAddrPowerOneMonth(addr)
	ginx.ResSuccess(c, addrPowerChart)
}

func GetAddrDetail(c *gin.Context) {
	addr := ginx.Param(c, "addr")
	proverDetail := repository.GetAddrDetailByAddress(addr)
	ginx.ResSuccess(c, proverDetail)
}

func GetPuzzleRewardChart(c *gin.Context) {
	chart := repository.GetRewardChartOneMonth()
	ginx.ResSuccess(c, chart)
}
