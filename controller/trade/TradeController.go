package trade

import (
	"TradeMatching/common/enum"
	"TradeMatching/common/glog"
	"TradeMatching/common/myContext"
	"TradeMatching/common/tool"
	"TradeMatching/common/utils"
	"TradeMatching/controller/req"
	"TradeMatching/server"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

// SellOrder掛賣單
func SellOrder(ctx *myContext.MyContext) {
	orderReq := req.SellOrderReq{}
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	glog.Infof("Trace %v SellOrder SellOrderReq:%v", ctx.Trace, orderReq)
	//OrderId = 方向前綴+memberId+挂单时间(時間搓)
	memberId := orderReq.MemberId
	orderType := orderReq.Type
	amount := orderReq.Amount
	price := orderReq.Price
	orderDirection := orderReq.Direction
	sellTime := time.Now()

	//訂單方向錯誤
	if orderDirection != enum.SELLSTR {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderDirectionError, nil))
		return
	}

	// 檢查掛單類型 市價:MARKET_PRICE 限價:LIMIT_PRICE
	if orderType != enum.MARKET_PRICE && orderType != enum.LIMIT_PRICE {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderTypeError, nil))
		return
	}

	// 限價單必須帶挂單價格且必須大於 0
	if orderType == enum.LIMIT_PRICE {
		if price.LessThanOrEqual(decimal.Zero) {
			ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.LimitNeedPriceError, nil))
			return
		}
	}

	orderId := server.GetOrderId(enum.SELL, strconv.FormatInt(memberId, 10), utils.GetMilliTimestamp(sellTime))

	order := server.ExchangeOrder{
		OrderId:      orderId,
		MemberId:     memberId,
		Type:         orderType,
		Amount:       amount,
		Symbol:       enum.SymbolBTCUSDT,
		CoinSymbol:   enum.CoinBTC,
		BasesSymbol:  enum.BaseUSDT,
		Status:       enum.TRADING,
		Direction:    enum.SELL,
		Price:        price,
		Time:         sellTime,
		UseDiscount:  "0",
		TradedAmount: decimal.Zero,
	}
	
	//創建訂單 紀錄訂明細(事務)
	err := server.CreateOrderANDRecordTradeDetailDelect(ctx, order)
	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}

	//判斷 市價單 限價單
	if orderType == enum.MARKET_PRICE {
		err = server.TradeMatchingService.AddMarketPriceOrder(ctx, order)

	} else if orderType == enum.LIMIT_PRICE {
		err = server.TradeMatchingService.AddLimitPriceOrder(ctx, order)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}
	ctx.JSON(http.StatusOK, tool.RespOkStatus(&tool.SellOrderSuccess, tool.RespOk(order)))
}

// BuyOrder掛買單
func BuyOrder(ctx *myContext.MyContext) {
	orderReq := req.BuyOrderReq{}
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	glog.Infof("Trace %v BuyOrder SellOrderReq:%v", ctx.Trace, orderReq)

	//OrderId = 方向前綴+memberId+挂单时间(時間搓)
	memberId := orderReq.MemberId
	orderType := orderReq.Type
	amount := orderReq.Amount
	price := orderReq.Price
	orderDirection := orderReq.Direction
	sellTime := time.Now()

	//訂單方向錯誤
	if orderDirection != enum.BUYSTR {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderDirectionError, nil))
		return
	}

	// 檢查掛單類型 市價:MARKET_PRICE 限價:LIMIT_PRICE
	if orderType != enum.MARKET_PRICE && orderType != enum.LIMIT_PRICE {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderTypeError, nil))
		return
	}

	// 限價單必須帶挂單價格且必須大於 0
	if orderType == enum.LIMIT_PRICE {
		if price.LessThanOrEqual(decimal.Zero) {
			ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.LimitNeedPriceError, nil))
			return
		}
	}

	//OrderId = 前綴+memberId+挂单时间(時間搓)
	orderId := server.GetOrderId(enum.BUY, strconv.FormatInt(memberId, 10), utils.GetMilliTimestamp(sellTime))

	order := server.ExchangeOrder{
		OrderId:      orderId,
		MemberId:     memberId,
		Type:         orderType,
		Amount:       amount,
		Symbol:       enum.SymbolBTCUSDT,
		CoinSymbol:   enum.CoinBTC,
		BasesSymbol:  enum.BaseUSDT,
		Status:       enum.TRADING,
		Direction:    enum.BUY,
		Price:        price,
		Time:         sellTime,
		UseDiscount:  "0",
		TradedAmount: decimal.Zero,
	}

	//創建訂單 紀錄訂明細(事務)
	err := server.CreateOrderANDRecordTradeDetailDelect(ctx, order)
	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}

	//判斷 市價單 限價單
	if orderType == enum.MARKET_PRICE {
		err = server.TradeMatchingService.AddMarketPriceOrder(ctx, order)
	} else if orderType == enum.LIMIT_PRICE {
		err = server.TradeMatchingService.AddLimitPriceOrder(ctx, order)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}
	ctx.JSON(http.StatusOK, tool.RespOkStatus(&tool.SellOrderSuccess, tool.RespOk(order)))
}

// BuyOrderCancel 買單取消->進MAP移除剩下訂單->檢查DB之前部分成功交易訂單
func BuyOrderCancel(ctx *myContext.MyContext) {
	sellOrderCancelReq := req.SellOrderCancelReq{}
	if err := ctx.ShouldBindJSON(&sellOrderCancelReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	glog.Infof("Trace %v SellOrder BuyOrderCancel:%v", ctx.Trace, sellOrderCancelReq)

	//訂單方向錯誤
	if sellOrderCancelReq.Direction != enum.BUYSTR {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderDirectionError, nil))
		return
	}

	orderId := sellOrderCancelReq.OrderId
	//memberId := sellOrderCancelReq.MemberId
	price := sellOrderCancelReq.Price
	glog.Infof("Trace %v SellOrderCancel SellOrderCancelReq:%v", ctx.Trace, sellOrderCancelReq)
	err := server.TradeMatchingService.OrderDelete(ctx, enum.BUY, orderId, price)
	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}
	ctx.JSON(http.StatusOK, tool.RespOk(nil))
}

// SellOrderCancel 賣單取消->進MAP移除剩下訂單->檢查DB之前部分成功交易訂單
func SellOrderCancel(ctx *myContext.MyContext) {
	sellOrderCancelReq := req.SellOrderCancelReq{}
	if err := ctx.ShouldBindJSON(&sellOrderCancelReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	glog.Infof("Trace %v SellOrder SellOrderCancel:%v", ctx.Trace, sellOrderCancelReq)

	//訂單方向錯誤
	if sellOrderCancelReq.Direction != enum.SELLSTR {
		ctx.JSON(http.StatusOK, tool.RespFailStatus(&tool.ExchangeOrderDirectionError, nil))
		return
	}

	orderId := sellOrderCancelReq.OrderId
	//memberId := sellOrderCancelReq.MemberId
	price := sellOrderCancelReq.Price
	glog.Infof("Trace %v SellOrderCancel SellOrderCancelReq:%v", ctx.Trace, sellOrderCancelReq)
	err := server.TradeMatchingService.OrderDelete(ctx, enum.SELL, orderId, price)
	if err != nil {
		ctx.JSON(http.StatusOK, tool.RespFail(err.Code(), err.Msg(), nil))
		return
	}
	ctx.JSON(http.StatusOK, tool.RespOk(nil))
}

// CheckQueueMap 檢查 買 賣 簿
func CheckQueueMap(ctx *myContext.MyContext) {
	var queueInfo []server.QueueInfo
	direction := ctx.Query("direction") //訂單方向
	if direction == "" {
		ctx.JSON(http.StatusOK, tool.RespFail(tool.ExchangeOrderDirectionError.Code, tool.ExchangeOrderDirectionError.Msg, nil))
		return
	}
	if direction == enum.BUYSTR {
		queueInfo = server.TradeMatchingService.CheckBuyQueueMap()

	} else if direction == enum.SELLSTR {
		queueInfo = server.TradeMatchingService.CheckSellQueueMap()
	}
	ctx.JSON(http.StatusOK, tool.RespOk(queueInfo))
}
