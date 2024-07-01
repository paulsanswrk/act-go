package bngx

type bingx_aop_response struct {
	EventType string   `json:"e"` //listenKeyExpired, ORDER_TRADE_UPDATE, ACCOUNT_UPDATE, ACCOUNT_CONFIG_UPDATE
	EventTime uint64   `json:"E"`
	Timestamp string   `json:"T"`
	Order     struct { // EventType == ORDER_TRADE_UPDATE https://bingx-api.github.io/docs/#/en-us/swapV2/socket/account.html#Order%20update%20push
		AvgPrice        string `json:"ap"`
		ClientOrderID   string `json:"c"`
		HandlingFee     string `json:"n"`
		OrderID         int64  `json:"i"`
		OrderQuantity   string `json:"q"`
		Status          string `json:"X"` //NEW, PARTIALLY_FILLED, FILLED, CANCELED, EXPIRED
		TransactionTime int64  `json:"T"`
	} `json:"o"`
	Account struct { //EventType == ACCOUNT_UPDATE https://bingx-api.github.io/docs/#/en-us/swapV2/socket/account.html#Account%20balance%20and%20position%20update%20push
		EventLaunchReason string `json:"m"` //DEPOSIT, WITHDRAW, ORDER, FUNDING_FEE
		BalanceInfo       []struct {
			AssetName     string `json:"a"`
			WalletBalance string `json:"wb"`
			//WalletBalanceExclIsolatedMargin       string `json:"cw"`
			WalletBalanceChangeAmount string `json:"bc"`
		} `json:"B"`
		TradeInfo []struct { //empty for FUNDING_FEE event
			Symbol                  string `json:"s"`
			Position                string `json:"pa"`
			EntryPrice              string `json:"ep"`
			UnrealizedProfitAndLoss string `json:"up"`
			MarginType              string `json:"mt"`
			IsolatedMargin          string `json:"iw"`
			PosDirection            string `json:"ps"`
		} `json:"P"`
	} `json:"a"`
}
