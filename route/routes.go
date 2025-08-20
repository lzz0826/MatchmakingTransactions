package routes

import (
	"TradeMatching/common/myContext"
	"TradeMatching/controller/index"
	"TradeMatching/controller/trade"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	unprotected := router.Group("/tradeMatching")
	{

		unprotected.GET("/test", myContext.Wrap(index.GetById))

		unprotected.POST("/buyOrder", myContext.Wrap(trade.BuyOrder))               //買單
		unprotected.POST("/sellOrder", myContext.Wrap(trade.SellOrder))             //賣單
		unprotected.POST("/buyOrderCancel", myContext.Wrap(trade.BuyOrderCancel))   //買單取消
		unprotected.POST("/sellOrderCancel", myContext.Wrap(trade.SellOrderCancel)) //賣單取消
		unprotected.GET("/checkQueueMap", myContext.Wrap(trade.CheckQueueMap))      //查詢 買 賣 單MAP

	}

}
