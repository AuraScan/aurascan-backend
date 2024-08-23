package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model/schema"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"github.com/gin-gonic/gin"
)

func GetTransitionList(c *gin.Context) {
	var listInfo schema.TransitionPageResp
	if err := c.BindJSON(&listInfo); err != nil {
		logger.Errorf("GetTransitionList BindJSON | %v", err)
		ginx.ResFailed(c, "invalid parameter!")
		return
	}

	page, size := ginx.GinPostPagination(listInfo.Page, listInfo.PageSize)

	transitions, count := repository.GetTransitionsByPage(page, size, listInfo.ProgramId)
	var data = struct {
		Blocks []*schema.TransitionListResp `json:"transitions"`
		Count  int64                        `json:"count"`
	}{transitions, count}

	ginx.ResSuccess(c, data)
}

func GetTransitionDetail(c *gin.Context) {
	ti := ginx.Param(c, "transition_id")
	tdr, err := repository.GetTransitionDetail(ti)
	if err != nil {
		ginx.ResFailed(c, "record not found")
		return
	}

	ginx.ResSuccess(c, tdr)
}
