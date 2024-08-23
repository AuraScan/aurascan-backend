package explorer

import (
	"aurascan-backend/internal/repository"
	"aurascan-backend/model/schema"
	"ch-common-package/cache"
	"ch-common-package/ginx"
	"ch-common-package/logger"
	"context"
	"github.com/gin-gonic/gin"
)

func GetNetwork(c *gin.Context) {
	networkInfo := repository.GetNetworkOverview()
	ginx.ResSuccess(c, networkInfo)
}

// 24h、7d、All
func GetNetworkPowerChart(c *gin.Context) {
	query := ginx.Param(c, "range")
	var res = make([]*schema.PowerChart, 0)
	switch query {
	case "day":
		res = repository.GetNetworkPowerChart()
	case "week":
		res = repository.GetNetwork7dPowerChart()
	case "all":
		res = repository.GetNetworkAllPowerChart()
	}
	ginx.ResSuccess(c, res)
}

// 24h、7d、All
func GetProofTargetChart(c *gin.Context) {
	query := ginx.Param(c, "range")
	var res = make([]*schema.ProofTargetChart, 0)
	switch query {
	case "day":
		res = repository.Get24hProofTargetChart()
	case "week":
		res = repository.Get7dProofTargetChart()
	case "all":
		res = repository.GetAllProofTargetChart()
	}
	ginx.ResSuccess(c, res)
}

func ValidatorOverview(c *gin.Context) {
	var res = make(map[string]string)
	var err error
	res, err = cache.Redis.HGetAll(context.TODO(), "validator_overview").Result()
	if err != nil {
		logger.Errorf("ValidatorOverview HGetAll | %v", err)
	}
	ginx.ResSuccess(c, res)
}
