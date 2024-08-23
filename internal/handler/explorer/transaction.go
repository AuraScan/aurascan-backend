package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
)

func GetTransactionList(c *gin.Context) {
	var listInfo schema.PageListReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetTransactionList BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	transactions, count := repository.GetTransactionsByPage(page, size)
	var data = struct {
		Blocks []*schema.TransactionListResp `json:"transactions"`
		Count  int64                         `json:"count"`
	}{transactions, count}

	ginx.ResSuccess(c, data)
}

func GetTransactionDetail(c *gin.Context) {
	ti := ginx.Param(c, "transaction_id")
	tdr, err := repository.GetTransactionDetail(ti)
	if err != nil {
		ginx.ResFailed(c, "record not found")
		return
	}

	ginx.ResSuccess(c, tdr)
}
