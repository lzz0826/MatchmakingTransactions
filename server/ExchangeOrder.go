package server

import (
	"TradeMatching/common/enum"
	"github.com/shopspring/decimal"
	"time"
)

// 撮合訂單構造
type ExchangeOrder struct {
	OrderId       string                      `json:"orderId"`       //撮合訂單號
	MemberId      int64                       `json:"memberId"`      //會員號
	Type          enum.ExchangeOrderType      `json:"type"`          //掛單類型
	Amount        decimal.Decimal             `json:"amount"`        //买入或卖出量，对于市价买入单表
	Symbol        enum.Symbol                 `json:"symbol"`        //交易符號
	TradedAmount  decimal.Decimal             `json:"tradedAmount"`  //成交量
	Turnover      decimal.Decimal             `json:"turnover"`      //成交額 對市價買賣有用
	CoinSymbol    enum.CoinSymbol             `json:"coinSymbol"`    //币单位
	BasesSymbol   enum.BaseSymbol             `json:"basesSymbol"`   //结算单位
	Status        enum.ExchangeOrderStatus    `json:"status"`        //订单状态
	Direction     enum.ExchangeOrderDirection `json:"direction"`     //订单方向
	Price         decimal.Decimal             `json:"price"`         //挂单价格
	Time          time.Time                   `json:"time"`          //挂单时间
	CompletedTime *time.Time                  `json:"completedTime"` //交易完成时间
	CanceledTime  *time.Time                  `json:"canceledTime"`  //取消时间
	UseDiscount   string                      `json:"useDiscount"`   //是否使用折扣 0 不使用 1使用
}

// 封裝 price、merge、orders 的結構
type QueueInfo struct {
	Price  decimal.Decimal  // 賣價
	Merge  *MergeOrder      // 該價位的合併掛單
	Orders []*ExchangeOrder // 該價位所有訂單
}
