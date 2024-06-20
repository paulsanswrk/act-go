package bngx

import (
	"ACT_GO/db"
	"ACT_GO/utils"
	"time"
)

type i_bingx_order_api interface {
	place_futures_order(orders ...swapOrderRequest)
	cancel_all_orders(symbol string)
}

type bingx_order_api struct {
}

func (api *bingx_order_api) place_futures_order(orders ...swapOrderRequest) {
	for n, order := range orders {
		if n > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		params_map, _ := utils.Struct_To_Map(order)

		url := "https://open-api.bingx.com/openApi/swap/v2/trade/order?" + build_and_sign_url(params_map)

		var res swapOrderResponse
		res_body, err := utils.HTTP_Request(url, "POST", map[string]string{"X-BX-APIKEY": apiKey}, &res)
		if err == nil {
			db.AddMessage("place_futures_order", order, res_body)
		} else {
			db.AddError(err, "place_futures_order", order)
			return
		}
	}

	//println(url)
	//fmt.Printf("%v", res)
}

func (api *bingx_order_api) cancel_all_orders(symbol string) {
	url := "https://open-api.bingx.com/openApi/swap/v2/trade/allOpenOrders?" + build_and_sign_url(map[string]interface{}{"symbol": symbol})
	response_string, err := utils.HTTP_Request(url, "DELETE", map[string]string{"X-BX-APIKEY": apiKey}, nil)
	if err != nil {
		db.AddError(err, "place_futures_order")
		return
	} else {
		db.AddMessage("cancel_all_orders", url, response_string)
	}

}
