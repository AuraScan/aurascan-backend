package router

import (
	"aurascan-backend/internal/handler/explorer"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(app *gin.Engine) {
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	g := app.Group("/api")

	v1 := g.Group("/v1")
	v1.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	//浏览器接口均使用v2
	v2 := g.Group("/v2")
	//获取全网概览
	v2.GET("/aleo/network/overview", explorer.GetNetwork)
	//获取首页全网算力图表
	v1.GET("/aleo/network/power/:range", explorer.GetNetworkPowerChart)
	//获取Proof Target变化曲线图接口
	v1.GET("/aleo/proof_target/chart/:range", explorer.GetProofTargetChart)
	//查询功能
	v2.GET("/aleo/search/:field", explorer.Search)
	//获取区块列表
	v2.POST("/aleo/block/list", explorer.GetBlockList)
	//获取区块详情
	v2.GET("/aleo/block/:height", explorer.GetBlockDetail)
	//获取区块Authority详情
	v2.GET("/aleo/block/authority/:height", explorer.GetBlockAuthority)
	//分页获取区块transaction列表
	v2.POST("/aleo/block/transaction", explorer.GetBlockTransaction)
	//分页获取区块solution列表
	v2.POST("/aleo/block/solution", explorer.GetBlockSolution)
	//获取transaction详情
	v2.GET("/aleo/transaction/:transaction_id", explorer.GetTransactionDetail)
	//获取交易列表
	v2.POST("/aleo/transaction/list", explorer.GetTransactionList)
	//获取Transition详情
	v2.GET("/aleo/transition/:transition_id", explorer.GetTransitionDetail)
	//获取Transition列表
	v2.POST("/aleo/transition/list", explorer.GetTransitionList)
	//获取validator列表
	v2.POST("/aleo/validator/list", explorer.GetValidatorList)
	//获取prover前三十
	v2.GET("/aleo/prover/list/:time_range", explorer.GetProverList)
	//获取地址的transfer列表
	v2.POST("/aleo/address/transfer/list", explorer.GetTransferByAddr)
	//获取Program详情
	v2.GET("/aleo/program/:program_id", explorer.GetProgramDetail)
	//Get program called count chart by id
	v2.GET("/aleo/program/chart/:program_id", explorer.GetProgramChartById)
	//获取Program列表
	v2.POST("/aleo/program/list", explorer.GetProgramList)
	//获取Program的Name列表
	v2.GET("/aleo/program/mapping/names/:program_id", explorer.GetMappingNameListByProgramId)
	//获取Program的Mapping的value
	v2.POST("/aleo/program/mapping/value", explorer.GetMappingValue)
	//获取Program的源码
	v2.GET("/aleo/program/source_code/:program_id", explorer.GetMappingSourceCode)
	//获取最近24小时Program调用前十
	v2.GET("/aleo/program/rank/day", explorer.Get24hTopProgram)
	//获取最近一个月Program调用前十
	v2.GET("/aleo/program/rank/month", explorer.GetOneMonthTopProgram)
	//获取最近一个月收益地址前十
	v2.GET("/aleo/reward/rank/month", explorer.GetTopAddrRewardChart)
	//获取最近一个月地址算力前十
	v2.GET("/aleo/power/rank/month", explorer.GetTopAddrPowerChart)
	//获取最近一个月地址solution总数前十
	v2.GET("/aleo/solutions/rank/month", explorer.GetTopAddrSolutionsChart)
	//获取指定Prover的Solution数量图表
	v2.GET("/aleo/prover/solutions/chart/:addr", explorer.GetAddrSolutionsChart)
	//获取指定Prover的收益图表
	v2.GET("/aleo/prover/reward/chart/:addr", explorer.GetAddrRewardChart)
	//获取指定Prover的算力图表
	v2.GET("/aleo/prover/power/chart/:addr", explorer.GetAddrPowerChart)
	//获取指定Prover的详情
	v2.GET("/aleo/addr/detail/:addr", explorer.GetAddrDetail)
	v2.GET("/aleo/puzzle_reward/chart", explorer.GetPuzzleRewardChart)
	//获取validator的stake图表
	v2.GET("/aleo/validator/stake/chart", explorer.GetStakeChart)
	//获取地址的solution列表
	v2.POST("/aleo/address/solutions", explorer.GetAddrSolution)
	//获取validator的overview
	v2.GET("/aleo/validator/overview", explorer.ValidatorOverview)
	//获取首页订阅
	v2.GET("/aleo/latest/sub", explorer.LatestHeightSub)
	//实现leo编译
	v2.POST("/aleo/leo/compile", explorer.CompileLeo)
}
