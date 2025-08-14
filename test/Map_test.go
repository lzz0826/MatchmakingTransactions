package model

import (
	"TradeMatching/server"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMatching(t *testing.T) {

	// 建立一筆賣單（價格 100，數量 5）
	sellOrder := server.ExchangeOrder{
		OrderId:      "S1",
		MemberId:     1,
		Type:         server.LIMIT_PRICE,
		Amount:       decimal.NewFromInt(5),
		Symbol:       "BTC/USDT",
		CoinSymbol:   "BTC",
		BasesSymbol:  "USDT",
		Status:       server.TRADING,
		Direction:    server.SELL,
		Price:        decimal.NewFromInt(100),
		Time:         time.Now(),
		UseDiscount:  "0",
		TradedAmount: decimal.Zero,
	}

	// 建立一筆買單（價格 100，數量 3）
	buyOrder := server.ExchangeOrder{
		OrderId:      "B1",
		MemberId:     2,
		Type:         server.LIMIT_PRICE,
		Amount:       decimal.NewFromInt(3),
		Symbol:       "BTC/USDT",
		CoinSymbol:   "BTC",
		BasesSymbol:  "USDT",
		Status:       server.TRADING,
		Direction:    server.BUY,
		Price:        decimal.NewFromInt(100),
		Time:         time.Now(),
		UseDiscount:  "0",
		TradedAmount: decimal.Zero,
	}

	// 先掛賣單
	fmt.Println("----------- 掛入賣單 S1 -----------")
	server.TradeMatchingService.AddLimitPriceOrder(sellOrder)

	// 檢查掛單後的賣單佇列
	fmt.Println("掛單後賣單隊列:")
	checkQueueMap()
	fmt.Println("掛單後買單隊列:")
	checkBuyQueueMap()

	// 掛入買單（應立即撮合）
	fmt.Println("----------- 掛入買單 B1，價格相符應立即撮合 -----------")
	server.TradeMatchingService.AddLimitPriceOrder(buyOrder)

	// 撮合完後再檢查賣單佇列
	fmt.Println("撮合後賣單隊列:")
	checkQueueMap()
	fmt.Println("掛單後買單隊列:")
	checkBuyQueueMap()
}

// 檢查賣單QueueMap
func checkQueueMap() {
	iterator := server.TradeMatchingService.GetSellQueueMap().Iterator()
	for iterator.Next() {
		price := iterator.Key().(decimal.Decimal)
		merge := iterator.Value().(*server.MergeOrder)
		fmt.Printf("  賣價: %s，訂單數: %d\n", price, merge.Size())
		for _, order := range merge.GetOrders() {
			fmt.Printf("    → %+v\n", order)
		}
	}
}

// 檢查買單QueueMap
func checkBuyQueueMap() {
	iterator := server.TradeMatchingService.GetBuyQueueMap().Iterator()
	for iterator.Next() {
		price := iterator.Key().(decimal.Decimal)
		merge := iterator.Value().(*server.MergeOrder)
		fmt.Printf("  買價: %s，訂單數: %d\n", price, merge.Size())
		for _, order := range merge.GetOrders() {
			fmt.Printf("    → %+v\n", order)
		}
	}
}

// 測試：向 buyQueue 插入多個相同價格的 ExchangeOrder
//func TestPutMultipleOrdersToBuyQueue(t *testing.T) {
//	price := decimal.NewFromFloat(100.5)
//
//	// 建立第一個 MergeOrder，加入一筆訂單
//	order1 := &model.MergeOrder{}
//	order1.Add(model.ExchangeOrder{
//		OrderId:   "order1",
//		Amount:    decimal.NewFromFloat(1.5),
//		Price:     price,
//		Direction: model.BUY,
//		Time:      time.Now(),
//	})
//	if err := model.PutBuyQueue(price, order1); err != nil {
//		t.Fatalf("PutBuyQueue error: %v", err)
//	}
//
//	// 建立第二筆訂單，加入到同一價格
//	order2 := model.ExchangeOrder{
//		OrderId:   "order2",
//		Amount:    decimal.NewFromFloat(2.3),
//		Price:     price,
//		Direction: model.BUY,
//		Time:      time.Now(),
//	}
//
//	// 取出現有 MergeOrder 並追加
//	if existing, found := model.GetBuyQueue(price); found {
//		existing.Add(order2)
//	} else {
//		t.Fatal("預期找到價格 100.5 的 MergeOrder，但未找到")
//	}
//
//	// 建立第三筆訂單，加入到同一價格
//	order3 := model.ExchangeOrder{
//		OrderId:   "order3",
//		Amount:    decimal.NewFromFloat(0.7),
//		Price:     price,
//		Direction: model.BUY,
//		Time:      time.Now(),
//	}
//
//	if existing, found := model.GetBuyQueue(price); found {
//		existing.Add(order3)
//	} else {
//		t.Fatal("預期找到價格 100.5 的 MergeOrder，但未找到")
//	}
//
//	// 驗證
//	result, found := model.GetBuyQueue(price)
//	if !found {
//		t.Fatal("找不到預期的價格 key")
//	}
//
//	if result.Size() != 3 {
//		t.Errorf("預期累加 3 筆訂單，實際只有 %d 筆", result.Size())
//	}
//
//	expectedTotal := decimal.NewFromFloat(1.5 + 2.3 + 0.7)
//	if !result.GetTotalAmount().Equal(expectedTotal) {
//		t.Errorf("總數量錯誤，預期 %s，實際 %s", expectedTotal, result.GetTotalAmount())
//	}
//
//	// 輸出結果
//	fmt.Println("價格:", price)
//	for _, ord := range result.GetOrders() {
//		fmt.Printf("  訂單ID: %s, 數量: %s\n", ord.OrderId, ord.Amount.String())
//	}
//	fmt.Println("總數量:", result.GetTotalAmount())
//}

func TestMergeOrderBasic(t *testing.T) {
	// 新建 MergeOrder
	m := &server.MergeOrder{}

	// 新增三筆 ExchangeOrder
	m.Add(server.ExchangeOrder{
		OrderId:   "order1",
		Amount:    decimal.NewFromFloat(1.5),
		Price:     decimal.NewFromFloat(100.5),
		Direction: server.BUY,
		Time:      time.Now(),
	})
	m.Add(server.ExchangeOrder{
		OrderId:   "order2",
		Amount:    decimal.NewFromFloat(2.3),
		Price:     decimal.NewFromFloat(100.5),
		Direction: server.BUY,
		Time:      time.Now(),
	})
	m.Add(server.ExchangeOrder{
		OrderId:   "order3",
		Amount:    decimal.NewFromFloat(0.7),
		Price:     decimal.NewFromFloat(100.5),
		Direction: server.BUY,
		Time:      time.Now(),
	})

	// 測試 Size
	if size := m.Size(); size != 3 {
		t.Errorf("預期Size為3，實際為 %d", size)
	}

	// 測試 GetFirst
	first := m.GetFirst()
	if first.OrderId != "order1" {
		t.Errorf("預期第一筆OrderId為order1，實際為 %s", first.OrderId)
	}

	// 測試 GetPrice
	price := m.GetPrice()
	expectedPrice := decimal.NewFromFloat(100.5)
	if !price.Equal(expectedPrice) {
		t.Errorf("預期價格為 %s，實際為 %s", expectedPrice.String(), price.String())
	}

	// 測試 GetTotalAmount
	total := m.GetTotalAmount()
	expectedTotal := decimal.NewFromFloat(1.5 + 2.3 + 0.7) // 4.5
	if !total.Equal(expectedTotal) {
		t.Errorf("預期總量為 %s，實際為 %s", expectedTotal.String(), total.String())
	}

	// 輸出確認
	fmt.Printf("訂單數量: %d\n", m.Size())
	fmt.Printf("第一筆訂單ID: %s\n", first.OrderId)
	fmt.Printf("價格: %s\n", price.String())
	fmt.Printf("總數量: %s\n", total.String())
}

//func TestBuyQueueViaAPI(t *testing.T) {
//	prices := []string{"100.5", "110.1", "101.0", "444.9", "122.0"}
//
//	for i, p := range prices {
//		order := &model.MergeOrder{}
//		order.Add(model.ExchangeOrder{
//			OrderId: strconv.Itoa(i),
//			Price:   decimal.New(11, 2),
//		})
//		price, err := decimal.NewFromString(p)
//		if err != nil {
//			t.Fatalf("decimal parse error: %v", err)
//		}
//		if err := model.PutBuyQueue(price, order); err != nil {
//			t.Fatalf("PutBuyQueue error: %v", err)
//		}
//	}
//	it := model.GetBuyQueueSnapshot()
//	for k, v := range it {
//		fmt.Println("價格:", k)
//		fmt.Println("ExchangeOrder:", v.GetFirst())
//	}
//}
