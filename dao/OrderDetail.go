package dao

import (
	"TradeMatching/common/enum"
	"TradeMatching/common/myContext"
	"TradeMatching/common/mysql"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

// OrderDetail 訂單交易明細
type OrderDetail struct {
	ID        int64                       `gorm:"primaryKey;autoIncrement;comment:'ID'" json:"id"`
	OrderId   string                      `gorm:"column:order_id;index:idx_order_id;comment:'訂單號'" json:"orderId"` //對應 ExchangeOrder 表order_id
	MemberId  int64                       `gorm:"column:member_id;index:idx_member_id;comment:'會員ID'" json:"memberId"`
	EventType enum.EventType              `gorm:"column:event_type;comment:'事件類型(CREATE, TRADE, CANCEL, AMEND, EXPIRE)'" json:"eventType"`
	Symbol    enum.Symbol                 `gorm:"column:symbol;comment:'交易符號'" json:"symbol"`
	Direction enum.ExchangeOrderDirection `gorm:"column:direction;comment:'訂單方向 BUY/SELL'" json:"direction"`
	Type      enum.ExchangeOrderType      `gorm:"column:order_type;comment:'訂單類型 LIMIT/MARKET'" json:"type"`

	// 僅部分事件會用到的欄位
	Price          decimal.Decimal `gorm:"column:price;comment:'價格(下單價或成交價)'" json:"price"`
	Amount         decimal.Decimal `gorm:"column:amount;comment:'掛單數量'" json:"amount"`
	TradedAmount   decimal.Decimal `gorm:"column:traded_amount;comment:'已成交數量'" json:"tradedAmount"`
	UntradedAmount decimal.Decimal `gorm:"column:untraded_amount;comment:'未成交數量'" json:"untradedAmount"`
	Turnover       decimal.Decimal `gorm:"column:turnover;comment:'成交金額(成交價×成交量)'" json:"turnover"`

	Reason     string `gorm:"column:reason;comment:'事件原因或備註'" json:"reason"`
	IPAddress  string `gorm:"column:ip_address;comment:'操作IP'" json:"ipAddress"`
	DeviceInfo string `gorm:"column:device_info;comment:'設備信息'" json:"deviceInfo"`
	ApiKeyId   string `gorm:"column:api_key_id;comment:'API Key ID'" json:"apiKeyId"`

	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime;comment:'事件時間'" json:"createdTime"`
}

func InsertOrderDetail(ctx *myContext.MyContext, orderDetail *OrderDetail) (err error) {
	err = mysql.GormDb.Table("order_detail").Omit("id").Create(&orderDetail).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
	return nil
}
