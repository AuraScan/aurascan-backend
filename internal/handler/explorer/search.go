package explorer

import (
	"aurascan-backend/internal/repository"
	"ch-common-package/ginx"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	field := ginx.Param(c, "field")
	searchType := repository.SearchByField(field)
	ginx.ResSuccess(c, searchType)
}
