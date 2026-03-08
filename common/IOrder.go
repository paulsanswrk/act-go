package common

type IOrder interface {
	GetOrderID() string
	//GetOrderType() string
	GetOrderSide() string
	GetOrderPrice() float64
	GetOrderUSDT_Total() float64
	GetOrderCoinQty() float64
	GetOrderStatus() string
	GetOrderTime() int64

	IsOpen() bool
	IsClose() bool
}
