package bngx

import (
	"ACT_GO/db"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var stopC = make(chan bool)

func TestClosestLoi(t *testing.T) {
	lo, hi, _, _ := bot1.closest_loi(decimal.NewFromFloat(1.2))
	assert.Equal(t, "1", lo.String())
	assert.Equal(t, "1.25", hi.String())
}

func TestPos_size_rule(t *testing.T) {
	for _, f := range []float64{0.25, 2.15, 3, 5} {
		fmt.Printf("%f -> %s\n", f, bot1.pos_size_rule(decimal.NewFromFloat(f)).String())
	}
}

func TestCalc_opening_closing_loi(t *testing.T) {
	fill_price, _ := decimal.NewFromString("3.45")
	bot1.Pos_size = 2
	lo, hi, lo_index, hi_index := bot1.closest_loi(fill_price)
	open_price, open_amt := bot1.calc_opening_loi(lo_index)
	close_price, close_amt := bot1.calc_closing_loi(hi_index)

	fmt.Println(map[string]interface{}{"fill_price": fill_price, "lo": lo, "hi": hi, "lo_index": lo_index, "hi_index": hi_index, "open_price": open_price, "open_amt": open_amt, "close_price": close_price, "close_amt": close_amt})
}

type mock_bingx_order_api struct {
	filled_orders_cnt int
	avg_entry         float64 //duplicating avg_entry/pos_size/balance for mock working
	pos_size          float64
	balance           float64
}

func (api *mock_bingx_order_api) cancel_all_orders(symbol string) {
}

var m sync.Mutex

func (api *mock_bingx_order_api) place_futures_order(orders ...swapOrderRequest) {
	m.Lock()
	if api.filled_orders_cnt >= 10 {
		close(stopC)
		return
	}
	api.filled_orders_cnt++

	//"fill" a random order
	order := orders[rand.Intn(len(orders))]

	response := new(bingx_aop_response)
	response.Order.ClientOrderID = order.ClientOrderID
	response.Order.AvgPrice = order.Price

	//calc

	price, _ := strconv.ParseFloat(order.Price, 64)
	buy_or_sell_amt, _ := strconv.ParseFloat(order.Quantity, 64)
	if order.Side == "SELL" {
		buy_or_sell_amt = -buy_or_sell_amt
	}
	spent_amt := buy_or_sell_amt * price

	if order.Side == "BUY" {
		api.avg_entry = (spent_amt + api.avg_entry*api.pos_size) / (api.pos_size + buy_or_sell_amt)
	}

	api.pos_size += buy_or_sell_amt
	api.balance -= spent_amt

	response.Account.TradeInfo.Symbol = bot1.symbol
	response.Account.TradeInfo.EntryPrice = api.avg_entry
	response.Account.TradeInfo.Position = fmt.Sprintf("%.2f", api.pos_size)
	response.Account.BalanceInfo.WalletBalance = api.balance

	db.AddMessage(fmt.Sprintf("filling futures order (%s): %.2f @ %s", order.Side, buy_or_sell_amt, order.Price), order)

	time.Sleep(200 * time.Millisecond)
	//db.AddMessage("pub_account_update.RunAll")
	pub_account_update.RunAll(*response)

	time.Sleep(200 * time.Millisecond)
	pub_filled_order.RunAll(*response)
	m.Unlock()

}

func TestBot1(t *testing.T) {
	db.TruncateLogs()
	bot1.order_api = &mock_bingx_order_api{balance: 43}
	Start_Bot1()
	<-stopC
}
