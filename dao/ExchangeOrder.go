package dao

import (
	"TradeMatching/common/myContext"
	"TradeMatching/common/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
	"time"
)

// ExchangeOrder 撮合訂單構造
type ExchangeOrder struct {
	ID            int             `gorm:"column:id;primaryKey;autoIncrement;comment:'id'" json:"id"`
	OrderId       string          `gorm:"column:order_id;comment:'撮合訂單號'" json:"orderId"`                   //撮合訂單號
	MemberId      int64           `gorm:"column:member_id;comment:'會員號'" json:"memberId"`                   //會員號
	Type          string          `gorm:"column:type;comment:'掛單類型'" json:"type"`                           //掛單類型
	Amount        decimal.Decimal `gorm:"column:amount;comment:'买入或卖出量，对于市价买入单表'" json:"amount"`            //买入或卖出量，对于市价买入单表
	Symbol        string          `gorm:"column:symbol;comment:'交易符號'" json:"symbol"`                       //交易符號
	TradedAmount  decimal.Decimal `gorm:"column:traded_amount;comment:'成交量'" json:"tradedAmount"`           //成交量
	Turnover      decimal.Decimal `gorm:"column:turnover;comment:'成交額 對市價買賣有用'" json:"turnover"`            //成交額 對市價買賣有用
	CoinSymbol    string          `gorm:"column:coin_symbol;comment:'币单位'" json:"coinSymbol"`               //币单位
	BasesSymbol   string          `gorm:"column:bases_symbol;comment:'结算单位'" json:"basesSymbol"`            //结算单位
	Status        string          `gorm:"column:status;comment:'订单状态'" json:"status"`                       //订单状态
	Direction     string          `gorm:"column:direction;comment:'订单方向'" json:"direction"`                 //订单方向
	Price         decimal.Decimal `gorm:"column:price;comment:'挂单价格'" json:"price"`                         //挂单价格
	Time          time.Time       `gorm:"column:time;comment:'挂单时间'" json:"time"`                           //挂单时间
	CompletedTime *time.Time      `gorm:"column:completed_time;comment:'交易完成时间'" json:"completedTime"`      //交易完成时间
	CanceledTime  *time.Time      `gorm:"column:canceled_time;comment:'取消时间'" json:"canceledTime"`          //取消时间
	UseDiscount   string          `gorm:"column:use_discount;comment:'是否使用折扣 0 不使用 1使用'" son:"useDiscount"` //是否使用折扣 0 不使用 1使用
}

func InsertExchangeOrder(ctx *myContext.MyContext, e *ExchangeOrder) (err error) {
	err = mysql.GormDb.Table("exchange_order").Omit("id").Create(&e).Error
	if err != nil {
		log.Println(err.Error())
		return
	}
	return nil
}

// UpdateExchangeOrderByOderId 更新 exchange_order 表中指定 order_id 的欄位 只有交易中才能被更新
func UpdateExchangeOrderByOderId(orderId string, updatesReq interface{}) (int64, error) {
	result := mysql.GormDb.Table("exchange_order").Where("order_id = ? AND status = 'TRADING' ", orderId).Updates(updatesReq)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// InsertExchangeOrderTransaction 事務
func InsertExchangeOrderTransaction(tx *gorm.DB, e *ExchangeOrder) (err error) {
	return tx.Table("exchange_order").Omit("id").Create(&e).Error
}

// UpdateExchangeOrderByOderIdTransaction 事務
func UpdateExchangeOrderByOderIdTransaction(tx *gorm.DB, orderId string, updatesReq interface{}) (int64, error) {
	result := tx.Table("exchange_order").Where("order_id = ? AND status = 'TRADING' ", orderId).Updates(updatesReq)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
