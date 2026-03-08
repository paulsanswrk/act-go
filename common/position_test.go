package common

import (
	"log"
	"testing"
)

type TestOrder struct {
	OrderID      string
	OrderSide    string
	OrderPrice   float64
	OrderUSDT    float64
	OrderCoinQty float64
	OrderStatus  string
	OrderTime    int64
}

// GetOrderID returns the ID of the order
func (o TestOrder) GetOrderID() string {
	return o.OrderID
}

// GetOrderSide returns the side of the order
func (o TestOrder) GetOrderSide() string {
	return o.OrderSide
}

// GetOrderPrice returns the price of the order
func (o TestOrder) GetOrderPrice() float64 {
	return o.OrderPrice
}

// GetOrderUSDT_Total returns the total USDT value of the order
func (o TestOrder) GetOrderUSDT_Total() float64 {
	return o.OrderUSDT
}

// GetOrderCoinQty returns the coin quantity of the order
func (o TestOrder) GetOrderCoinQty() float64 {
	return o.OrderCoinQty
}

// GetOrderStatus returns the status of the order
func (o TestOrder) GetOrderStatus() string {
	return o.OrderStatus
}

// GetOrderTime returns the time of the order
func (o TestOrder) GetOrderTime() int64 {
	return o.OrderTime
}

// IsOpen checks if the order is open
func (o TestOrder) IsOpen() bool {
	return o.OrderSide == "buy"
}

// IsClose checks if the order is closed
func (o TestOrder) IsClose() bool {
	return o.OrderSide == "sell"
}

func TestGetEntryPrice(t *testing.T) {

	orders := []IOrder{
		TestOrder{OrderSide: "buy", OrderPrice: 1, OrderUSDT: 100, OrderCoinQty: 100},
		TestOrder{OrderSide: "buy", OrderPrice: 2, OrderUSDT: 100, OrderCoinQty: 50},
		TestOrder{OrderSide: "sell", OrderPrice: 3, OrderUSDT: 150, OrderCoinQty: 50},
		//TestOrder{OrderSide: "buy", OrderPrice: 3, OrderUSDT: 150, OrderCoinQty: 50},
	}

	position := Position{}

	entry := position.GetEntryPrice(orders)

	log.Printf("entry price: %f", entry)
}
