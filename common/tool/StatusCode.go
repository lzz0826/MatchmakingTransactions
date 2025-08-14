package tool

type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var statusList = []Status{
	Success,
	SystemError,
	SetRedisKeyError,
	LoginError,
	UserError,
	PasswordError,
	TokenError,
	CreateTokenError,
	MissingParameters,
	ExchangeOrderDirectionError,
	BuyOrdersNotFound,
	SellOrdersNotFound,
	ExchangeOrderUpdateError,
	RecordTradeDetailError,
	RecordOrderDetailError,
	SellOrderSuccess,
	BuyOrderSuccess,
	CreateOrderError,
	ExchangeOrderTypeError,
	LimitNeedPriceError,
}

// 添加狀態 上面statusList也需要添加 才能收尋到
var (
	//系統 0
	Success           = Status{Code: 0, Msg: "成功"}
	SystemError       = Status{Code: -1, Msg: "失敗"}
	MissingParameters = Status{Code: 1, Msg: "缺少必要參數"}

	//Redis 100
	SetRedisKeyError = Status{Code: 100, Msg: "Redis Set錯誤"}

	//登入 1000
	LoginError       = Status{Code: 1000, Msg: "登入失敗"}
	UserError        = Status{Code: 1001, Msg: "帳號錯誤"}
	PasswordError    = Status{Code: 1002, Msg: "密碼錯誤"}
	TokenError       = Status{Code: 1003, Msg: "Token無效"}
	CreateTokenError = Status{Code: 1004, Msg: "生產Token失敗"}

	//撮合交易
	SellOrderSuccess = Status{Code: 2001, Msg: "賣單下單成功"}
	BuyOrderSuccess  = Status{Code: 2002, Msg: "買單下單成功"}

	ExchangeOrderDirectionError = Status{Code: 3000, Msg: "訂單方向錯誤"}
	BuyOrdersNotFound           = Status{Code: 3001, Msg: "找不到符合的買單"}
	SellOrdersNotFound          = Status{Code: 3002, Msg: "找不到符合的賣單"}
	ExchangeOrderUpdateError    = Status{Code: 3003, Msg: "更新訂單失敗"}
	RecordTradeDetailError      = Status{Code: 3004, Msg: "紀錄撮合交易明細失敗"}
	RecordOrderDetailError      = Status{Code: 3005, Msg: "紀錄訂單交易明細失敗"}
	CreateOrderError            = Status{Code: 3006, Msg: "創建訂單失敗"}
	ExchangeOrderTypeError      = Status{Code: 3007, Msg: "掛單類型錯誤"}
	LimitNeedPriceError         = Status{Code: 3008, Msg: "限價單必須要掛單價"}
)

func GetStatusByCode(code int) Status {
	for _, status := range statusList {
		if status.Code == code {
			return status
		}
	}
	return Status{Code: -1, Msg: "未知狀態"}
}

func GetStatusByMsg(msg string) Status {
	for _, status := range statusList {
		if status.Msg == msg {
			return status
		}
	}
	return Status{Code: -1, Msg: "未知狀態"}
}

func GetStatusCodeFromError(err error) int {
	status := GetStatusByMsg(err.Error())
	return status.Code
}

func GetStatusMsgFromError(err error) string {
	status := GetStatusByMsg(err.Error())
	return status.Msg
}
