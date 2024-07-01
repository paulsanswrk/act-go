package bngx

import (
	"ACT_GO/db"
	"ACT_GO/utils"
	"time"
)

type i_bot interface {
	BuySellSide(do_open bool) string
	PosDirection() string
	addMessage(msg string, args ...interface{})
}

type i_bingx_order_api interface {
	place_futures_order(bot i_bot, orders ...swapOrderRequest) (res swapOrderResponse, response_string string, err error)
	cancel_all_orders(symbol string) (response_string string, err error)
	cancel_order(symbol string, ClientOrderID string) (response_string string, err error)
}

type bingx_order_api struct {
}

func (api *bingx_order_api) place_futures_order(bot i_bot, orders ...swapOrderRequest) (res swapOrderResponse, response_string string, err error) {
	for n, order := range orders {
		if n > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		params_map, _ := utils.Struct_To_Map(order)

		url := "https://open-api.bingx.com/openApi/swap/v2/trade/order?" + build_and_sign_url(params_map)

		//var res swapOrderResponse
		response_string, err = utils.HTTP_Request(url, "POST", map[string]string{"X-BX-APIKEY": apiKey}, &res)
		if err == nil {
			bot.addMessage("place_futures_order", order, response_string)
		} else {
			db.AddError(err, "place_futures_order", order)
			return
		}
	}

	return
	//println(url)
	//fmt.Printf("%v", res)
}

func (api *bingx_order_api) cancel_all_orders(symbol string) (response_string string, err error) {
	url := "https://open-api.bingx.com/openApi/swap/v2/trade/allOpenOrders?" + build_and_sign_url(map[string]interface{}{"symbol": symbol})
	response_string, err = utils.HTTP_Request(url, "DELETE", map[string]string{"X-BX-APIKEY": apiKey}, nil)
	if err != nil {
		db.AddError(err, "cancel_all_orders")
	} else {
		db.AddMessage("cancel_all_orders", url, response_string)
	}
	return
}

func (api *bingx_order_api) list_pending_orders(symbol string) (response_string string, err error) {
	url := "https://open-api.bingx.com/openApi/swap/v2/trade/openOrders?" + build_and_sign_url(map[string]interface{}{"symbol": symbol, "timestamp": time.Now().UnixNano() / 1e6})
	response_string, err = utils.HTTP_Request(url, "GET", map[string]string{"X-BX-APIKEY": apiKey}, nil)
	if err != nil {
		db.AddError(err, "list_pending_orders")
		return
	} else {
		db.AddMessage("list_pending_orders", url, response_string)
	}
	return

}

func (api *bingx_order_api) cancel_order(symbol string, ClientOrderID string) (response_string string, err error) {
	url := "https://open-api.bingx.com/openApi/swap/v2/trade/order?" + build_and_sign_url(map[string]interface{}{"symbol": symbol, "timestamp": time.Now().UnixNano() / 1e6, "clientOrderID": ClientOrderID})
	response_string, err = utils.HTTP_Request(url, "DELETE", map[string]string{"X-BX-APIKEY": apiKey}, nil)
	if err != nil {
		db.AddError(err, "cancel_order")
	} else {
		db.AddMessage("cancel_order", url, response_string)
	}
	return
}
