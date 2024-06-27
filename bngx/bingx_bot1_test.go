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
	prev, next, _, _ := bot_long.closest_loi(decimal.NewFromFloat(1.25))
	assert.Equal(t, "1.2", prev.String())
	assert.Equal(t, "1.3", next.String())

	prev, next, _, _ = bot_short.closest_loi(decimal.NewFromFloat(1.25))
	assert.Equal(t, "1.3", prev.String())
	assert.Equal(t, "1.2", next.String())
}

func TestPos_size_rule(t *testing.T) {
	fmt.Printf("%v\n", bot_long.loi)
	for _, f := range []float64{0.5, 1, 1.5, 1.8, 1.9, 2, 2.1, 3, 4.5} {
		fmt.Printf("Long: %f -> %s\n", f, bot_long.pos_size_rule(decimal.NewFromFloat(f)).String())
		//fmt.Printf("Short: %f -> %s\n", f, bot_short.pos_size_rule(decimal.NewFromFloat(f)).String())
	}
}

func TestCalc_opening_closing_loi(t *testing.T) {
	fill_price, _ := decimal.NewFromString("3.45")
	bot_long.Pos_size = 2
	lo, hi, lo_index, hi_index := bot_long.closest_loi(fill_price)
	open_price, open_amt := bot_long.calc_opening_loi(lo_index)
	close_price, close_amt := bot_long.calc_closing_loi(hi_index)

	fmt.Println(map[string]interface{}{"fill_price": fill_price, "lo": lo, "hi": hi, "lo_index": lo_index, "hi_index": hi_index, "open_price": open_price, "open_amt": open_amt, "close_price": close_price, "close_amt": close_amt})
}

type mock_bingx_order_api struct {
	filled_orders_cnt int
	avg_entry         float64 //duplicating avg_entry/pos_size/balance for mock working
	pos_size          float64
	balance           float64
}

func (api *mock_bingx_order_api) cancel_all_orders(string) (response_string string, err error) {
	return
}
func (api *mock_bingx_order_api) cancel_order(string, string) (response_string string, err error) {
	return
}

var m sync.Mutex

func (api *mock_bingx_order_api) place_futures_order(testing_bot i_bot, orders ...swapOrderRequest) (res swapOrderResponse, response_string string, err error) {
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
	buy_or_sell_amt, _ := strconv.ParseFloat(order.Quantity, 64) //positive

	if order.Side == testing_bot.BuySellSide(true) { //opening order
		spent_amt_usdt := buy_or_sell_amt * price //positive //coins * usdt/coins = usdt
		api.avg_entry = (spent_amt_usdt + api.avg_entry*api.pos_size) / (api.pos_size + buy_or_sell_amt)
		api.balance -= spent_amt_usdt
		api.pos_size += buy_or_sell_amt
	} else { //closing
		received_amount_usdt := buy_or_sell_amt * price //positive
		api.balance += received_amount_usdt
		api.pos_size -= buy_or_sell_amt //pos_size is positive
	}

	response.Account.TradeInfo.Symbol = order.Symbol
	response.Account.TradeInfo.PosDirection = testing_bot.PosDirection()
	response.Account.TradeInfo.EntryPrice = api.avg_entry
	response.Account.TradeInfo.Position = fmt.Sprintf("%.2f", api.pos_size)
	response.Account.BalanceInfo.WalletBalance = api.balance

	testing_bot.addMessage(fmt.Sprintf("filling futures order (%s): %.2f @ %s", order.Side, buy_or_sell_amt, order.Price), order)

	time.Sleep(200 * time.Millisecond)
	//db.AddMessage("pub_account_update.RunAll")
	pub_account_update.RunAll(*response)

	time.Sleep(200 * time.Millisecond)
	pub_filled_order.RunAll(*response)
	m.Unlock()

	return
}

func TestBot1(t *testing.T) {
	db.TruncateLogs()
	db.DB.Exec("truncate table bots.bot_plan restart IDENTITY")
	bot_short.order_api = &mock_bingx_order_api{balance: bot_short.Balance}
	bot_short.start()
	<-stopC
}

func Test2Bots(t *testing.T) {
	db.TruncateLogs()
	db.DB.Exec("truncate table bots.bot_plan restart IDENTITY")
	bot_long.order_api = &mock_bingx_order_api{balance: bot_short.Balance}
	bot_short.order_api = &mock_bingx_order_api{balance: bot_short.Balance}
	go bot_long.start()
	go bot_short.start()
	<-stopC
}

func TestDbGet(t *testing.T) {
	var plan = BotPlan{BotID: "www"}
	//db.DB.Save(&plan)
	//res := db.DB.First(&plan, "bot_id=?", "www") //works
	res := db.DB.First(&plan) //works
	//res := db.DB.Model(BotPlan{BotID: "www"}).First(&plan) //doesn't work for non-existing key!

	//stmt := db.DB.Session(&gorm.Session{DryRun: true}).Where(&BotPlan{BotID: "qqq"}).First(&plan).Statement
	//stmt.SQL.String()
	//println(stmt.SQL.String())

	println(res.RowsAffected)
}

func TestLoi(t *testing.T) {
	fmt.Printf("%v\n", bot_long.loi)
	fmt.Printf("%v\n", bot_short.loi)
}

func TestDecimal(t *testing.T) {
	var (
		d1 = decimal.NewFromFloat(1)
		d2 = decimal.NewFromFloat(2)
	)

	println(d1.Compare(d2))
	println(d2.Compare(d1))
}
