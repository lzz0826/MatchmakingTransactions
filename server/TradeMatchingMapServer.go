package server

import (
	"TradeMatching/common/enum"
	"TradeMatching/common/errors"
	"TradeMatching/common/glog"
	"TradeMatching/common/myContext"
	"TradeMatching/common/mysql"
	"TradeMatching/common/tool"
	"context"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
	"time"
)

func init() {
	buyLimitPriceQueue = BuyLimitPriceQueueMap{
		buyLimitPriceQueue: treemap.NewWith(buyPriceComparator),
	}
	sellLimitPriceQueue = SellLimitPriceQueueMap{
		sellLimitPriceQueue: treemap.NewWith(sellPriceComparator),
	}
	TradeMatchingService = tradeMatchingService{}
}

type tradeMatchingService struct {
	lock sync.RWMutex
}

var (
	buyLimitPriceQueue  BuyLimitPriceQueueMap  //限價買單
	sellLimitPriceQueue SellLimitPriceQueueMap //限價賣單
	//給外部調用
	TradeMatchingService tradeMatchingService
)

// 買入限價單，價格由高到低排列
type BuyLimitPriceQueueMap struct {
	buyLimitPriceQueue *treemap.Map
}

// 賣出限價訂單鍊錶，價格由低到高排列
type SellLimitPriceQueueMap struct {
	sellLimitPriceQueue *treemap.Map
}

// sellPriceComparator 賣單價格比較器 (價格從低到高排列)
func sellPriceComparator(a, b interface{}) int {
	da := a.(decimal.Decimal)
	db := b.(decimal.Decimal)
	return da.Cmp(db)
}

// buyPriceComparator 買單價格比較器（價格由高到低）
func buyPriceComparator(a, b interface{}) int {
	da := a.(decimal.Decimal)
	db := b.(decimal.Decimal)
	return db.Cmp(da) // 反轉
}

// addToBuyQueue 添加買單列
func addToBuyQueue(e ExchangeOrder) {
	if value, found := buyLimitPriceQueue.buyLimitPriceQueue.Get(e.Price); found {
		order := value.(*MergeOrder)
		order.Add(e)
	} else {
		mer := &MergeOrder{}
		mer.Add(e)
		buyLimitPriceQueue.buyLimitPriceQueue.Put(e.Price, mer)
	}
}

// addToSellQueue 添加賣單列
func addToSellQueue(e ExchangeOrder) {
	if value, found := sellLimitPriceQueue.sellLimitPriceQueue.Get(e.Price); found {
		order := value.(*MergeOrder)
		order.Add(e)
	} else {
		mer := &MergeOrder{}
		mer.Add(e)
		sellLimitPriceQueue.sellLimitPriceQueue.Put(e.Price, mer)
	}
}

func (t *tradeMatchingService) GetSellQueueMap() *treemap.Map {
	return sellLimitPriceQueue.sellLimitPriceQueue
}

func (t *tradeMatchingService) GetBuyQueueMap() *treemap.Map {
	return buyLimitPriceQueue.buyLimitPriceQueue
}

// 檢查買單QueueMap
func (t *tradeMatchingService) CheckBuyQueueMap() []QueueInfo {
	var result []QueueInfo
	iterator := t.GetBuyQueueMap().Iterator()
	for iterator.Next() {
		price := iterator.Key().(decimal.Decimal)
		merge := iterator.Value().(*MergeOrder)

		var orders []*ExchangeOrder
		for _, order := range merge.GetOrders() {
			orders = append(orders, &order)
		}
		result = append(result, QueueInfo{
			Price:  price,
			Merge:  merge,
			Orders: orders,
		})
	}
	return result
}

// 檢查賣單QueueMap
func (t *tradeMatchingService) CheckSellQueueMap() []QueueInfo {
	var result []QueueInfo
	iterator := t.GetSellQueueMap().Iterator()
	for iterator.Next() {
		price := iterator.Key().(decimal.Decimal)
		merge := iterator.Value().(*MergeOrder)

		var orders []*ExchangeOrder
		for _, order := range merge.GetOrders() {
			orders = append(orders, &order)
		}
		result = append(result, QueueInfo{
			Price:  price,
			Merge:  merge,
			Orders: orders,
		})
	}
	return result
}

// AddMarketPriceOrder 市價單撮合主要業務邏輯
// 參數：
// - ctx: 請求上下文（用於記錄日誌、交易資料等）
// - e:   待撮合的市價訂單（包含方向、數量等）
// 撮合流程：
// 1. 根據訂單方向（BUY / SELL）選擇對手方的訂單簿
// 2. 從最佳價位開始遍歷對手方訂單，直到：該市價單數量用完，或對手方訂單簿沒有剩餘可成交的單
// 3. 在每個價位呼叫 matchOrders 進行成交撮合，若該價位的對手單全部成交，從訂單簿移除該價位
// 4. 撮合結束後：若數量還有剩，狀態標記為 PARTIAL_COMPLETED（部分成交），若全部成交，狀態標記為 COMPLETED（完成成交）
// 5. 市價單不會掛單簿
func (t *tradeMatchingService) AddMarketPriceOrder(ctx *myContext.MyContext, e ExchangeOrder) *errors.Errx {
	t.lock.Lock()
	defer t.lock.Unlock()
	if e.Direction == enum.BUY {
		iterator := sellLimitPriceQueue.sellLimitPriceQueue.Iterator()
		for iterator.Next() && e.Amount.GreaterThan(decimal.Zero) {
			price := iterator.Key().(decimal.Decimal)
			//市價單則是指不指定價格，以當時市場上的最佳價格立即成交。
			mergeOrder := iterator.Value().(*MergeOrder)
			err := matchOrders(ctx, mysql.GormDb, &e, &mergeOrder.orders)
			if err != nil {
				return err
			}
			// 若整個價格點的賣單都吃光了，從 map 移除該價位
			if len(mergeOrder.orders) == 0 {
				sellLimitPriceQueue.sellLimitPriceQueue.Remove(price)
			}
		}
		// 撮合完 立即成交，未成交的部分應直接取消
		if e.Amount.GreaterThan(decimal.Zero) {
			//部分成交
			err := UpdateTradedOrder(ctx, mysql.GormDb, e, enum.PARTIAL_COMPLETED)
			if err != nil {
				return err
			}
		} else {
			//當前訂完整搓合完更新訂單
			err := UpdateTradedOrder(ctx, mysql.GormDb, e, enum.COMPLETED)
			if err != nil {
				return err
			}
		}
	} else if e.Direction == enum.SELL {
		iterator := buyLimitPriceQueue.buyLimitPriceQueue.Iterator()
		for iterator.Next() && e.Amount.GreaterThan(decimal.Zero) {
			price := iterator.Key().(decimal.Decimal)
			//市價單則是指不指定價格，以當時市場上的最佳價格立即成交。
			mergeOrder := iterator.Value().(*MergeOrder)
			err := matchOrders(ctx, mysql.GormDb, &e, &mergeOrder.orders)
			if err != nil {
				return err
			}
			// 若整個價格點的賣單都吃光了，從 map 移除該價位
			if len(mergeOrder.orders) == 0 {
				buyLimitPriceQueue.buyLimitPriceQueue.Remove(price)
			}
		}
		// 撮合完 立即成交，未成交的部分應直接取消
		if e.Amount.GreaterThan(decimal.Zero) {
			//部分成交
			err := UpdateTradedOrder(ctx, mysql.GormDb, e, enum.PARTIAL_COMPLETED)
			if err != nil {
				return err
			}
		} else {
			//當前訂完整搓合完更新訂單
			err := UpdateTradedOrder(ctx, mysql.GormDb, e, enum.COMPLETED)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddLimitPriceOrder 限價撮合主要業務邏輯
// 根據訂單方向（買/賣）與價格，嘗試與對方限價訂單撮合。
// 若有剩餘數量，則將訂單加入相應掛單隊列。
// 拍撮合完成的訂單會更新狀態為完成。
//
// 參數：
// - ctx: 請求上下文
// - e: 待撮合的限價訂單，包含訂單信息與剩餘數量
//
// 流程：
// 1. 若為買單，遍歷賣方限價訂單價格從低到高，找出價格 <= 買單價格的對手訂單進行撮合。
// 2. 若為賣單，遍歷買方限價訂單價格從高到低，找出價格 >= 賣單價格的對手訂單進行撮合。
// 3. 使用 matchOrders 進行撮合，更新對手隊列。
// 4. 若該價位所有訂單已撮合完，從隊列中移除該價格。
// 5. 若撮合後仍有剩餘數量，將訂單加入對應隊列；否則標記訂單為完成。
// 買進單需要匹配的價格不大於委託價，否則退出
// 賣出單需要匹配的價格不小於委託價，否則退出
func (t *tradeMatchingService) AddLimitPriceOrder(ctx *myContext.MyContext, e ExchangeOrder) *errors.Errx {
	t.lock.Lock()
	defer t.lock.Unlock()
	if e.Direction == enum.BUY {
		//帶事務 參考
		err := AddLimitPriceOrderBuyTransaction(ctx, e)
		if err != nil {
			return err
		}
	} else if e.Direction == enum.SELL {
		iterator := buyLimitPriceQueue.buyLimitPriceQueue.Iterator()
		for iterator.Next() && e.Amount.GreaterThan(decimal.Zero) {
			price := iterator.Key().(decimal.Decimal)
			// 只匹配買價 <= 賣價
			if e.Price.LessThanOrEqual(price) {
				mergeOrder := iterator.Value().(*MergeOrder)
				matchOrders(ctx, mysql.GormDb, &e, &mergeOrder.orders)
				// 若整個價格點的賣單都吃光了，從 map 移除該價位
				if len(mergeOrder.orders) == 0 {
					buyLimitPriceQueue.buyLimitPriceQueue.Remove(price)
				}
			} else {
				break // 價格不符合，停止遍歷
			}
		}
		// 撮合完還有剩，掛入賣單隊列
		if e.Amount.GreaterThan(decimal.Zero) {
			addToSellQueue(e)
		} else {
			//當前訂完整搓合完更新訂單
			UpdateTradedOrder(ctx, mysql.GormDb, e, enum.COMPLETED)
		}
	}
	return nil
}

// 撮合函數 當前訂單 對手訂單表
// matchOrders 撮合當前訂單與對手訂單列表。
// 會根據撮合規則，不斷撮合，直到當前訂單成交完畢或對手訂單耗盡。
// 每次撮合會更新訂單的成交數量、成交金額，並記錄交易明細與更新訂單狀態。
//
// 參數：
// - ctx: 請求上下文，包含請求相關資料與狀態。
// - tx:  事務上下文
// - current: 指向當前待撮合訂單的指標。
// - opponents: 指向對手訂單切片指標，按照先進先出順序撮合。
//
// 執行流程：
// 1. 循環撮合對手訂單直到對手訂單用盡或當前訂單數量為零。
// 2. 取出最早的對手訂單，計算本次撮合數量（雙方剩餘量的最小值）。
// 3. 根據撮合數量計算成交金額，更新雙方成交數量與成交額。
// 4. 紀錄本次撮合明細，更新訂單狀態為交易中。
// 5. 若對手訂單成交完畢，更新其狀態為完成並從對手列表中移除；否則停止撮合等待下一輪。
func matchOrders(ctx *myContext.MyContext, tx *gorm.DB, current *ExchangeOrder, opponents *[]ExchangeOrder) *errors.Errx {
	// 只要還有對手單、並且當前訂單還沒撮合完，就持續撮合
	for len(*opponents) > 0 && current.Amount.GreaterThan(decimal.Zero) {

		// 取得最早一筆對手單（FIFO）
		opponent := &(*opponents)[0]

		// 撮合數量為兩者之間量的最小值
		matchAmount := decimal.Min(opponent.Amount, current.Amount)

		//計算此次撮合的成交金額
		turnover := matchAmount.Mul(opponent.Price)

		// 從對手單扣除撮合數量
		opponent.Amount = opponent.Amount.Sub(matchAmount)
		// 累加對手單的成交量
		opponent.TradedAmount = opponent.TradedAmount.Add(matchAmount)
		//對手成交額
		opponent.Turnover = opponent.Turnover.Add(turnover)

		// 從當前訂單扣除撮合數量
		current.Amount = current.Amount.Sub(matchAmount)
		// 累加當前訂單的成交量
		current.TradedAmount = current.TradedAmount.Add(matchAmount)
		// 當前成交額
		current.Turnover = current.Turnover.Add(turnover)

		glog.Infof("Trace %v BuyOrder matchOrders: 當前方orderID: %v 對手方orderID: %v", ctx.Trace, current.OrderId, current.OrderId)

		// 紀錄此次撮合交易紀錄 下單方 對手方
		repError := RecordTradeDetail(ctx, tx, *current, *opponent, turnover, matchAmount)
		if repError != nil {
			return repError
		}
		// 更新撮合交易訂單 當前 交易中
		repError = UpdateTradedOrder(ctx, tx, *current, enum.TRADING)
		if repError != nil {
			return repError
		}
		// 更新撮合交易訂單 對手 交易中
		repError = UpdateTradedOrder(ctx, tx, *opponent, enum.TRADING)
		if repError != nil {
			return repError
		}
		// 如果對手單已全部成交，從佇列中移除
		if opponent.Amount.IsZero() {
			// 更新撮合交易訂單 對手 交易完成
			repError = UpdateTradedOrder(ctx, tx, *opponent, enum.COMPLETED)
			if repError != nil {
				return repError
			}
			*opponents = (*opponents)[1:] // 移除第一筆
		} else {
			// 對手單還有剩，就不再往下撮合（等下一次進來再繼續）
			break
		}
	}
	return nil
}

// OrderDelete 取消訂單 (只能取消尚未撮合完成的訂單)
// 參數：
// - ctx: 請求上下文
// - e: 訂單方向 (BUY 或 SELL)
// - orderId: 要取消的訂單號
// - orderPrice: 訂單價格，用來快速定位價格隊列
//
// 回傳：
// - *errors.Errx 可能包含錯誤訊息
func (t *tradeMatchingService) OrderDelete(ctx *myContext.MyContext, e enum.ExchangeOrderDirection, orderId string, orderPrice decimal.Decimal) *errors.Errx {
	t.lock.Lock()
	defer t.lock.Unlock()
	// 只允許取消尚未撮合完成的部分
	if e == enum.BUY {
		iterator := buyLimitPriceQueue.buyLimitPriceQueue.Iterator()
		isBreak := false
		found := false
		for iterator.Next() && !isBreak {
			price := iterator.Key().(decimal.Decimal)
			// 找到該筆下單價位
			if orderPrice.Equal(price) {
				mergeOrder := iterator.Value().(*MergeOrder)
				for i, v := range mergeOrder.orders {
					if v.OrderId == orderId {
						// 紀錄取消交易明細
						RecordTradeDetailDelect(ctx, mysql.GormDb, v, enum.CANCEL)
						// 更新資料庫，將訂單設為取消狀態
						UpdateTradedOrder(ctx, mysql.GormDb, v, enum.CANCELED)
						// 從該價格隊列中移除該筆訂單
						mergeOrder.orders = append(mergeOrder.orders[:i], mergeOrder.orders[i+1:]...)
						// 若該價位訂單已清空，則從整個撮合佇列中移除該價位
						if len(mergeOrder.orders) == 0 {
							sellLimitPriceQueue.sellLimitPriceQueue.Remove(price)
						}
						isBreak = true
						found = true
						break
					}
				}
			}
		}
		if !found {
			//找不到符合的買單
			return errors.NewBizErrx(tool.BuyOrdersNotFound.Code, tool.BuyOrdersNotFound.Msg)
		}
	} else if e == enum.SELL {
		iterator := sellLimitPriceQueue.sellLimitPriceQueue.Iterator()
		isBreak := false
		found := false
		for iterator.Next() && !isBreak {
			price := iterator.Key().(decimal.Decimal)
			// 找到該筆下單價位
			if orderPrice.Equal(price) {
				mergeOrder := iterator.Value().(*MergeOrder)
				for i, v := range mergeOrder.orders {
					if v.OrderId == orderId {
						// 紀錄取消交易明細
						RecordTradeDetailDelect(ctx, mysql.GormDb, v, enum.CANCEL)
						// 更新資料庫，將訂單設為取消狀態
						UpdateTradedOrder(ctx, mysql.GormDb, v, enum.CANCELED)
						// 從該價格隊列中移除該筆訂單
						mergeOrder.orders = append(mergeOrder.orders[:i], mergeOrder.orders[i+1:]...)
						// 若該價位訂單已清空，則從整個撮合佇列中移除該價位
						if len(mergeOrder.orders) == 0 {
							buyLimitPriceQueue.buyLimitPriceQueue.Remove(price)
						}
						isBreak = true
						found = true
						break
					}
				}
			}
		}
		if !found {
			//找不到符合的賣單
			return errors.NewBizErrx(tool.SellOrdersNotFound.Code, tool.SellOrdersNotFound.Msg)
		}
	}
	return nil
}

// AddLimitPriceOrderBuyTransaction 事務 TODO **事務確實回滾但會與 買賣簿MAP不同步
func AddLimitPriceOrderBuyTransaction(muCtx *myContext.MyContext, e ExchangeOrder) *errors.Errx {
	// 加個 context timeout，避免交易卡死
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var repError *errors.Errx
	err := mysql.GormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		iterator := sellLimitPriceQueue.sellLimitPriceQueue.Iterator()
		for iterator.Next() && e.Amount.GreaterThan(decimal.Zero) {
			price := iterator.Key().(decimal.Decimal)
			// 只匹配賣價 <= 買價
			if e.Price.GreaterThanOrEqual(price) {
				mergeOrder := iterator.Value().(*MergeOrder)
				e1 := matchOrders(muCtx, tx, &e, &mergeOrder.orders)
				if e1 != nil {
					repError = e1
					return e1
				}
				// 若整個價格點的賣單都吃光了，從 map 移除該價位
				if len(mergeOrder.orders) == 0 {
					sellLimitPriceQueue.sellLimitPriceQueue.Remove(price)
				}
			} else {
				break // 價格不符合，停止遍歷
			}
		}
		// 撮合完還有剩，掛入買單隊列
		if e.Amount.GreaterThan(decimal.Zero) {
			addToBuyQueue(e)
		} else {
			//當前訂完整搓合完更新訂單
			e2 := UpdateTradedOrder(muCtx, tx, e, enum.COMPLETED)
			if e2 != nil {
				repError = e2
				return e2
			}
		}
		// 任一錯誤直接 return -> 自動 Rollback；正常 return nil -> 自動 Commit
		return nil
	})
	if err != nil {
		// 統一處理失敗
		return repError
	}
	return nil
}
