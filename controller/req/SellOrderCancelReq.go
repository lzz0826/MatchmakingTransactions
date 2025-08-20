package req

import "github.com/shopspring/decimal"

type SellOrderCancelReq struct {
	OrderId   string          `json:"orderId" binding:"required"  example:"BUY100011755067160543"` //撮合訂單號
	MemberId  int64           `json:"memberId" binding:"required"  example:"10001"`                //會員號
	Direction string          `json:"direction" binding:"required"  example:"BUY"`                 //訂單方向
	Price     decimal.Decimal `json:"price" binding:"required"  example:"30000"`                   //挂单价格
}
