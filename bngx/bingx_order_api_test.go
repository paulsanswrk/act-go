package bngx

import (
	"testing"
)

const HOST = "open-api.bingx.com"

var API_KEY = apiKey
var API_SECRET = secretKey

func TestPlaceTestOrder2(t *testing.T) {
	_, response_string, _ := new(bingx_order_api).place_futures_order(&bot_long, swapOrderRequest{
		Symbol:        "BTC-USDT",
		Side:          "BUY",
		PositionSide:  "LONG",
		Type:          "LIMIT",
		Quantity:      "0.01",
		Price:         "15000",
		ClientOrderID: "ClientOrderID_3",
	})
	//place_futures_order("BTC-USDT", "BUY", "LONG", "MARKET", 5, 60000)

	println(response_string)
}

func TestCancelAllOrders(t *testing.T) {
	new(bingx_order_api).cancel_all_orders("WIF-USDT")
}

func TestList_pending_orders(t *testing.T) {
	resp, _ := new(bingx_order_api).list_pending_orders("WIF-USDT")

	println(resp)
}

func Test_cancel_order(t *testing.T) {
	resp, _ := new(bingx_order_api).cancel_order("BTC-USDT", "ClientOrderID_3")

	println(resp)
}
