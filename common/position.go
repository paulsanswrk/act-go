package common

import "fmt"

type Position struct {
	Symbol string
	Amount float64
}

func (p *Position) String() string {
	return fmt.Sprintf("%s %f", p.Symbol, p.Amount)
}

func (p *Position) GetEntryPrice(orders []IOrder) (entry float64) {
	sum_coin_qty := 0.0
	sum_usdt_amt := 0.0

	for _, order := range orders {
		if order.IsOpen() {
			sum_coin_qty += order.GetOrderCoinQty()
			sum_usdt_amt += order.GetOrderUSDT_Total()
			entry = sum_usdt_amt / sum_coin_qty
		} else {
			sum_coin_qty -= order.GetOrderCoinQty()
			sum_usdt_amt -= order.GetOrderCoinQty() * entry
		}
	}

	return
}
