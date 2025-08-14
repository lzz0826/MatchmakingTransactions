package rep

import (
	"TradeMatching/server"
	"github.com/shopspring/decimal"
)

type CheckQueueMapRep struct {
	Price  decimal.Decimal         // 賣價
	Merge  *server.MergeOrder      // 該價位的合併掛單
	Orders []*server.ExchangeOrder // 該價位所有訂單
}
