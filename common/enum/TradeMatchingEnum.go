package enum

// 掛單類型
type ExchangeOrderType string

const (
	//市價單沒有價格，所以一定吃對手單的價格
	//限價單遇限價單，也用對手單（先掛單）的價格
	//
	//logger.info(">>>>>市价单>>>交易与限价单交易");
	//与限价单交易
	MARKET_PRICE     ExchangeOrderType = "MARKET_PRICE" //市價 市價單則是指不指定價格，以當時市場上的最佳價格立即成交。
	MARKET_PRICE_STR string            = "MARKET_PRICE"
	LIMIT_PRICE      ExchangeOrderType = "LIMIT_PRICE" //限價 限價單是指投資人設定一個特定的價格，只有在市場價格達到或優於該價格時才會成交。
	LIMIT_PRICE_STR  string            = "LIMIT_PRICE"
	//先与限价单交易 后与市价单交易
	//logger.info(">>>>>限价单>>>交易与限价单交易");
	//logger.info(">>>>限价单未交易完>>>>与市价单交易>>>>");

)

// 订单状态
type ExchangeOrderStatus string

const (
	TRADING           ExchangeOrderStatus = "TRADING"           //交易中
	PARTIAL_COMPLETED ExchangeOrderStatus = "PARTIAL_COMPLETED" //部分成交
	COMPLETED         ExchangeOrderStatus = "COMPLETED"         //完全成交
	PARTIAL_CANCELED  ExchangeOrderStatus = "PARTIAL_CANCELED"  //部分取消
	CANCELED          ExchangeOrderStatus = "CANCELED"          //完全取消
	REJECTED          ExchangeOrderStatus = "REJECTED"          // 訂單被拒
	OVERTIMED         ExchangeOrderStatus = "OVERTIMED"         //超時
)

// 判斷是否為終態
func IsFinalStatus(status ExchangeOrderStatus) bool {
	switch status {
	case PARTIAL_COMPLETED, COMPLETED, PARTIAL_CANCELED, CANCELED, REJECTED, OVERTIMED:
		return true
	default:
		return false
	}
}

// 驗證是否允許狀態轉移
func CanTransition(from, to ExchangeOrderStatus) bool {
	if IsFinalStatus(from) {
		// 如果是終態，就不能再轉換
		return false
	}
	// 只有從 TRADING 可以轉到其他終態
	if from == TRADING && IsFinalStatus(to) {
		return true
	}

	return false
}

// 訂單方向
type ExchangeOrderDirection string

const (
	BUY     ExchangeOrderDirection = "BUY"
	BUYSTR  string                 = "BUY"
	SELL    ExchangeOrderDirection = "SELL"
	SELLSTR string                 = "SELL"
)

// Symbol 市場對（例如：BTC/USDT）
type Symbol string

const (
	SymbolBTCUSDT Symbol = "BTC/USDT"
	SymbolETHUSDT Symbol = "ETH/USDT"
	SymbolBNBUSDT Symbol = "BNB/USDT"
	SymbolADAUSDT Symbol = "ADA/USDT"
	SymbolSOLUSDT Symbol = "SOL/USDT"
)

// CoinSymbol 幣的標記（交易幣）
type CoinSymbol string

const (
	CoinBTC CoinSymbol = "BTC"
	CoinETH CoinSymbol = "ETH"
	CoinBNB CoinSymbol = "BNB"
	CoinADA CoinSymbol = "ADA"
	CoinSOL CoinSymbol = "SOL"
)

// BaseSymbol 結算單位
type BaseSymbol string

const (
	BaseUSDT BaseSymbol = "USDT"
	BaseUSD  BaseSymbol = "USD"
	BaseBUSD BaseSymbol = "BUSD"
)

// 事件類型

type EventType string

const (
	CREATE EventType = "CREATE" //建立
	TRADE  EventType = "TRADE"  //交易
	CANCEL EventType = "CANCEL" //取消
	AMEND  EventType = "AMEND"  //修改
	EXPIRE EventType = "EXPIRE" // 到期
)
