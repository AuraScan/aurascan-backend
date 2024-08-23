package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetBlockList(c *gin.Context) {
	var listInfo schema.PageListReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetBlocksByPage BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	blocks, count := repository.GetBlocksByPage(page, size)
	var data = struct {
		Blocks []*schema.BlockListResp `json:"blocks"`
		Count  int64                   `json:"count"`
	}{blocks, count}

	ginx.ResSuccess(c, data)
}

func GetBlockDetail(c *gin.Context) {
	query := ginx.Param(c, "height")
	height, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		logger.Errorf("GetBlockDetail ParseInt(%v) to int | %v", query, err)
		ginx.ResFailed(c, "invalid parameter")
		return
	}
	blockDetail, err := repository.GetBlockDetailByHeight(height)
	if err != nil {
		logger.Errorf("GetBlockDetail | %v", err)
		ginx.ResFailed(c, "internal error")
		return
	}
	ginx.ResSuccess(c, blockDetail)
}

func GetBlockAuthority(c *gin.Context) {
	query := ginx.Param(c, "height")
	height, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		logger.Errorf("GetBlockAuthority ParseInt(%v) to int | %v", query, err)
		ginx.ResFailed(c, "invalid parameter")
		return
	}
	authoritys := model.GetBlockAuthoritys(height)
	ginx.ResSuccess(c, authoritys)
}

func GetBlockSolution(c *gin.Context) {
	var listInfo schema.BlockSpecReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetBlockSolution BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}
	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)
	solutions, total, totalTarget := repository.GetSolutionListByHeight(listInfo.Height, page, size)
	var data = struct {
		Solutions   []*schema.SolutionInBlockResp `json:"solutions"`
		TotalTarget float64                       `json:"total_target"`
		Count       int64                         `json:"count"`
	}{solutions, totalTarget, total}

	ginx.ResSuccess(c, data)
}

func GetBlockTransaction(c *gin.Context) {
	var listInfo schema.BlockSpecReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetBlockSolution BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}
	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)
	transactions, count := repository.GetTransactionsByHeight(listInfo.Height, page, size)
	var data = struct {
		Transactions []*schema.TransactionListInBlockResp `json:"transactions"`
		Count        int64                                `json:"count"`
	}{transactions, count}

	ginx.ResSuccess(c, data)
}
