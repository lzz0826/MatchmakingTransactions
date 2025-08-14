package req

import "github.com/shopspring/decimal"

type SellOrderCancelReq struct {
	OrderId   string          `json:"orderId" binding:"required"`   //撮合訂單號
	MemberId  int64           `json:"memberId" binding:"required"`  //會員號
	Direction string          `json:"direction" binding:"required"` //訂單方向
	Price     decimal.Decimal `json:"price" binding:"required"`     //挂单价格
}
