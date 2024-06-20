package bngx

import (
	"testing"
)

const HOST = "open-api.bingx.com"

var API_KEY = apiKey
var API_SECRET = secretKey

func TestPlaceTestOrder2(t *testing.T) {
	new(bingx_order_api).place_futures_order(swapOrderRequest{
		Symbol:        "BTC-USDT",
		Side:          "BUY",
		PositionSide:  "LONG",
		Type:          "MARKET",
		Quantity:      "5",
		Price:         "60000",
		ClientOrderID: "ClientOrderID_1",
	})
	//place_futures_order("BTC-USDT", "BUY", "LONG", "MARKET", 5, 60000)
}

func TestCancelAllOrders(t *testing.T) {
	new(bingx_order_api).cancel_all_orders("WIF-USDT")
}
