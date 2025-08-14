package req

import (
	"TradeMatching/common/enum"
	"github.com/shopspring/decimal"
)

type BuyOrderReq struct {
	MemberId   int64                  `json:"memberId" binding:"required"`   // 會員號
	Type       enum.ExchangeOrderType `json:"type" binding:"required"`       //掛單類型 MARKET_PRICE LIMIT_PRICE
	Direction  string                 `json:"direction" binding:"required"`  //訂單方向
	Amount     decimal.Decimal        `json:"amount" binding:"required"`     //买入或卖出量，对于市价买入单表
	Symbol     string                 `json:"symbol" binding:"required"`     //交易符號
	CoinSymbol string                 `json:"coinSymbol" binding:"required"` //币单位
	Price      decimal.Decimal        `json:"price"`                         //挂单价格
}
