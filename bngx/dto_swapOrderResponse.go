package bngx

type swapOrderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Order struct {
			OrderID         int64   `json:"orderId"`
			Symbol          string  `json:"symbol"`
			PositionSide    string  `json:"positionSide"`
			Side            string  `json:"side"`
			Type            string  `json:"type"`
			Price           float64 `json:"price"`
			Quantity        float64 `json:"quantity"`
			StopPrice       float64 `json:"stopPrice"`
			WorkingType     string  `json:"workingType"`
			ClientOrderID   string  `json:"clientOrderID"`
			TimeInForce     string  `json:"timeInForce"`
			PriceRate       float64 `json:"priceRate"`
			StopLoss        string  `json:"stopLoss"`
			TakeProfit      string  `json:"takeProfit"`
			ReduceOnly      bool    `json:"reduceOnly"`
			ActivationPrice float64 `json:"activationPrice"`
			ClosePosition   string  `json:"closePosition"`
			StopGuaranteed  string  `json:"stopGuaranteed"`
		} `json:"order"`
	} `json:"data"`
}
