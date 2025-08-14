package server

import (
	"github.com/shopspring/decimal"
)

// 「價格一樣」的訂單，可以把它們包在一個 MergeOrder 裡
type MergeOrder struct {
	orders []ExchangeOrder
}

// 最后位置添加一个
func (m *MergeOrder) Add(order ExchangeOrder) {
	m.orders = append(m.orders, order)
}

func (m *MergeOrder) GetFirst() *ExchangeOrder {
	return &m.orders[0]
}

func (m *MergeOrder) Size() int {
	return len(m.orders)
}

func (m *MergeOrder) GetPrice() decimal.Decimal {
	if len(m.orders) == 0 {
		return decimal.Zero
	}
	return m.orders[0].Price
}

func (m *MergeOrder) GetOrders() []ExchangeOrder {
	return m.orders
}

func (m *MergeOrder) GetTotalAmount() decimal.Decimal {
	total := decimal.Zero
	for _, v := range m.orders {
		total = total.Add(v.Amount)
	}
	return total
}

func (m *MergeOrder) GetOrdersMemoryAddress() *[]ExchangeOrder {
	return &m.orders
}
