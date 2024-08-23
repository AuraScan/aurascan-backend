package explorer

import (
	"aurascan-backend/model"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
)

func GetTransferByAddr(c *gin.Context) {
	var listInfo schema.TransferListReq
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetTransferByAddr BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	transfers, count := model.GetTransferByAddr(listInfo.Address, page, size)
	var data = struct {
		Validators []*model.TransferRes `json:"transfers"`
		Count      int64                `json:"count"`
	}{transfers, count}

	ginx.ResSuccess(c, data)
}
