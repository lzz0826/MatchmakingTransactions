package req

import "github.com/shopspring/decimal"

type BuyOrderCancelReq struct {
	OrderId   string          `json:"orderId" binding:"required" example:"SELL100011755067521826"` //撮合訂單號
	MemberId  int64           `json:"memberId" binding:"required" example:"10001"`                 //會員號
	Direction string          `json:"direction" binding:"required" example:"SELL"`                 //訂單方向
	Price     decimal.Decimal `json:"price" binding:"required" example:"30000"`                    //挂单价格
}
