package dao

import (
	"TradeMatching/common/myContext"
	"TradeMatching/common/mysql"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

// TradeDetail 撮合交易明細構造
type TradeDetail struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement;comment:'id'" json:"id"`
	BuyOrderId  string          `gorm:"column:buy_order_id;comment:'買方訂單 ID'" json:"buyOrderId"`   //對應買方 ExchangeOrder 表order_id
	SellOrderId string          `gorm:"column:sell_order_id;comment:'賣方訂單 ID'" json:"sellOrderId"` //對應賣方 ExchangeOrder 表order_id
	Price       decimal.Decimal `gorm:"column:price;comment:'成交價格'" json:"price"`
	DealAmount  decimal.Decimal `gorm:"column:amount;comment:'成交數量'" json:"amount"`
	Symbol      string          `gorm:"column:symbol;comment:'交易符號'" json:"symbol"`
	Remark      string          `gorm:"column:remark;comment:'備注'" json:"remark"`
	TradeTime   time.Time       `gorm:"column:trade_time;comment:'审核时间'" json:"trade_time"`
	CreatedTime time.Time       `gorm:"column:created_time;comment:'创建时间'" json:"createdTime"`
}

func InsertTradeDetail(ctx *myContext.MyContext, tradeDetail *TradeDetail) (err error) {
	err = mysql.GormDb.Table("trade_detail").Omit("id").Create(&tradeDetail).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
	return nil
}
