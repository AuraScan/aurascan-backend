package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
)

func GetTopAddrRewardChart(c *gin.Context) {
	addrRewardChart := repository.GetTop10AddrDailyRewardChartOneMonth()
	ginx.ResSuccess(c, addrRewardChart)
}

func GetTopAddrSolutionsChart(c *gin.Context) {
	addrSolutionsChart := repository.GetTop10AddrDailySolutionsChartOneMonth()
	ginx.ResSuccess(c, addrSolutionsChart)
}

func GetAddrSolution(c *gin.Context) {
	var listInfo schema.SolutionInAddrReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetAddrSolution BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}
	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)
	solutions, total := repository.GetSolutionListByAddr(listInfo.Address, page, size)
	var data = struct {
		Solutions []*schema.SolutionInAddrResp `json:"solutions"`
		Count     int64                        `json:"count"`
	}{solutions, total}

	ginx.ResSuccess(c, data)
}
