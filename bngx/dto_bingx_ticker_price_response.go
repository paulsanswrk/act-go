package bngx

type bingx_ticker_price_response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
		Time   int64  `json:"time"`
	} `json:"data"`
}
