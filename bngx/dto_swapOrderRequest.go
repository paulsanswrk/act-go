package bngx

type swapOrderRequest struct {
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Type          string `json:"type"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	ClientOrderID string `json:"clientOrderID"` //up to 40 characters

	TakeProfit string `json:"takeProfit"`
}
