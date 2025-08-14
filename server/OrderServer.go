package server

import (
	"TradeMatching/common/enum"
	"TradeMatching/common/errors"
	"TradeMatching/common/glog"
	"TradeMatching/common/myContext"
	"TradeMatching/common/tool"
	"TradeMatching/dao"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

// GetOrderId 生成訂單 ID：前綴 + memberId + 掛單時間
func GetOrderId(direction enum.ExchangeOrderDirection, memberId string, milliTimestamp string) string {
	return fmt.Sprintf("%s%s%s", direction, memberId, milliTimestamp)
}

// CreateOrder 創建訂單。
// 參數:
// - e: 當前下單的訂單資訊。
// 回傳:
// - *errors.Errx: 如果發生錯誤則回傳錯誤，成功則為 nil。
func CreateOrder(ctx *myContext.MyContext, e ExchangeOrder) *errors.Errx {
	glog.Infof("創建訂單 CreateOrder Trace: %v ExchangeOrder: %v", ctx.Trace, e)
	order := &dao.ExchangeOrder{
		OrderId:       e.OrderId,
		MemberId:      e.MemberId,
		Type:          string(e.Type), // enum to string
		Amount:        e.Amount,
		Symbol:        string(e.Symbol),
		TradedAmount:  e.TradedAmount,
		Turnover:      e.Turnover,
		CoinSymbol:    string(e.CoinSymbol),
		BasesSymbol:   string(e.BasesSymbol),
		Status:        string(e.Status),
		Direction:     string(e.Direction),
		Price:         e.Price,
		Time:          e.Time,
		CompletedTime: e.CompletedTime,
		CanceledTime:  e.CanceledTime,
		UseDiscount:   e.UseDiscount,
	}
	err := dao.InsertExchangeOrder(ctx, order)
	if err != nil {
		return errors.NewBizErrx(tool.CreateOrderError.Code, tool.CreateOrderError.Msg)
	}
	return nil
}

// RecordTradeDetail 紀錄此次 "撮合" 交易的詳細記錄，包含下單方與對手方資訊。
// 參數:
// - ctx: 請求上下文，包含請求相關資料。
// - current: 當前下單的訂單資訊。
// - opponent: 對手方的訂單資訊。
// - price: 撮合後的成交價格。
// - matchAmount: 本次撮合的成交數量。
// 回傳:
// - *errors.Errx: 如果發生錯誤則回傳錯誤，成功則為 nil。
func RecordTradeDetail(ctx *myContext.MyContext, current ExchangeOrder, opponent ExchangeOrder, price, matchAmount decimal.Decimal) *errors.Errx {
	detail := dao.TradeDetail{
		Price:       price,
		DealAmount:  matchAmount,
		Symbol:      string(current.Symbol),
		Remark:      "交易備注",
		TradeTime:   time.Now(),
		CreatedTime: time.Now(),
	}
	if current.Direction == enum.BUY {
		detail.BuyOrderId = current.OrderId
		detail.SellOrderId = opponent.OrderId
	} else {
		detail.BuyOrderId = opponent.OrderId
		detail.SellOrderId = current.OrderId
	}
	glog.Info("紀錄交易紀錄...", current)
	err := dao.InsertTradeDetail(ctx, &detail)
	if err != nil {
		glog.Errorf("recordTradeDetail err:%v", err)
		return errors.NewBizErrx(tool.RecordTradeDetailError.Code, tool.RecordTradeDetailError.Msg)
	}
	return nil
}

// UpdateTradedOrder 更新撮合交易訂單 更新訂單 更新狀態
// 參數:
// - e: 當前下單的訂單資訊。
// - status: 订单状态。
// 回傳:
// - *errors.Errx: 如果發生錯誤則回傳錯誤，成功則為 nil。
func UpdateTradedOrder(ctx *myContext.MyContext, e ExchangeOrder, status enum.ExchangeOrderStatus) *errors.Errx {
	glog.Infof("更新撮合交易訂單 updateTradedOrder: Trace: %v ExchangeOrder: %v", ctx.Trace, e)
	updates := map[string]interface{}{
		"amount":        e.Amount,       //买入或卖出量，对于市价买入单表
		"traded_amount": e.TradedAmount, //成交量
		"turnover":      e.Turnover,     //成交額 對市價買賣有用
		"status":        status,
	}
	now := time.Now()
	switch status {
	case enum.COMPLETED, enum.PARTIAL_COMPLETED:
		updates["completed_time"] = now
	case enum.CANCELED, enum.PARTIAL_CANCELED:
		updates["canceled_time"] = now
	}
	rowsAffected, err := dao.UpdateExchangeOrderByOderId(e.OrderId, updates)
	if err != nil {
		glog.Errorf("更新撮合交易訂單 updateTradedOrder: Trace: %v ExchangeOrder: %v", ctx.Trace, e)
		return errors.NewBizErrx(tool.ExchangeOrderUpdateError.Code, tool.ExchangeOrderUpdateError.Msg)
	}

	if rowsAffected != 1 {
		glog.Errorf("更新撮合交易訂單 沒有找到對應訂單 updateTradedOrder: Trace: %v ExchangeOrder: %v", ctx.Trace, e)
		return errors.NewBizErrx(tool.ExchangeOrderUpdateError.Code, tool.ExchangeOrderUpdateError.Msg)
	}
	return nil
}

// RecordTradeDetailDelect 紀錄此次 "訂單" 交易取消明細
// 參數:
// - ctx: 請求上下文，包含請求相關資料。
// - current: 當前下單的訂單資訊。
// 回傳:
// - *errors.Errx: 如果發生錯誤則回傳錯誤，成功則為 nil。
func RecordTradeDetailDelect(ctx *myContext.MyContext, current ExchangeOrder, et enum.EventType) *errors.Errx {
	oderDetail := dao.OrderDetail{
		OrderId:        current.OrderId,
		MemberId:       current.MemberId,
		EventType:      et,
		Symbol:         current.Symbol,
		Direction:      current.Direction,
		Type:           current.Type,
		Price:          current.Price,
		Amount:         current.Amount,
		TradedAmount:   current.TradedAmount,
		UntradedAmount: current.Amount.Sub(current.TradedAmount),
		Turnover:       current.Turnover,
		Reason:         "",
		IPAddress:      "",
		DeviceInfo:     "",
		ApiKeyId:       "",
		CreatedTime:    time.Now(),
	}
	glog.Infof("紀錄交易取消明細 recordTradeDetailDelect Trace: %v oderDetail: %v", ctx.Trace, oderDetail)
	err := dao.InsertOrderDetail(ctx, &oderDetail)
	if err != nil {
		glog.Errorf("紀錄交易取消明細失敗: recordTradeDetailDelect Trace: %v error: %v", ctx.Trace, err)
		return errors.NewBizErrx(tool.RecordOrderDetailError.Code, tool.RecordOrderDetailError.Msg)
	}
	return nil
}
