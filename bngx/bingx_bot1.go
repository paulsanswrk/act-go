package bngx

import (
	"ACT_GO/db"
	"ACT_GO/utils"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"log"
	"math"
	"slices"
	"strconv"
)

type bingx_bot1 struct {
	ID              string
	order_api       i_bingx_order_api
	bot_started     bool
	symbol          string
	is_long         bool
	loi             []decimal.Decimal
	Balance         float64
	Avg_entry       float64
	Pos_size        float64
	pStart          float64
	pEnd            float64
	min_notable_amt float64
	//approx_start_entry           float64
	num_fractional_digits        int32
	pub_account_update_listen_id uint64
	pub_filled_order_listen_id   uint64
}

var (
	bot_long = bingx_bot1{
		ID:              "BingX1-WIF-Long",
		order_api:       &bingx_order_api{},
		bot_started:     false,
		symbol:          "WIF-USDT",
		is_long:         true,
		loi:             funk.Map(utils.Range(0.5, 4.5, 0.05), func(n float64) decimal.Decimal { return decimal.NewFromFloatWithExponent(n, -2) }).([]decimal.Decimal),
		Balance:         220,
		Avg_entry:       2.3,
		Pos_size:        0,
		pStart:          0.5,
		pEnd:            4.5,
		min_notable_amt: 4.45,
		//approx_start_entry:    2.15,
		num_fractional_digits: 2,
	}

	bot_short = bingx_bot1{
		ID:              "BingX1-WIF-Short",
		order_api:       &bingx_order_api{},
		bot_started:     false,
		symbol:          "WIF-USDT",
		is_long:         false,
		loi:             funk.Map(utils.Range(4.5, 0.5, 0.1), func(n float64) decimal.Decimal { return decimal.NewFromFloatWithExponent(n, -2) }).([]decimal.Decimal),
		Balance:         220,
		Avg_entry:       2.3,
		Pos_size:        0,
		pStart:          4.5,
		pEnd:            0.5,
		min_notable_amt: 4.45,
		//approx_start_entry:    2.1,
		num_fractional_digits: 2,
	}
)

func init() {
}

func (bot *bingx_bot1) addMessage(msg string, args ...interface{}) {
	db.AddMessageWithModule(msg, bot.ID, args)
}

func (bot *bingx_bot1) pos_size_rule(price decimal.Decimal) decimal.Decimal {
	//pStart := 0.25
	//pEnd := 5.25
	result := 1 - (price.InexactFloat64()-bot.pStart)/(bot.pEnd-bot.pStart)
	result *= 360
	return decimal.NewFromFloat(result)
}

func (bot *bingx_bot1) closest_loi(price decimal.Decimal) (prev decimal.Decimal, next decimal.Decimal, prev_index int, next_index int) {
	for i, p := range bot.loi {
		if p.Compare(price) == bot.cmp_sign() { //p > price for long; p < price for short
			return bot.loi[i-1], p, i - 1, i
		}
	}

	log.Fatalf("closest_loi failed for %v", price)
	return
}

func (bot *bingx_bot1) PosDirection() string {
	if bot.is_long {
		return "LONG"
	} else {
		return "SHORT"
	}
}

func (bot *bingx_bot1) BuySellSide(do_open bool) string {
	if bot.is_long == do_open {
		return "BUY"
	} else {
		return "SELL"
	}
}

func (bot *bingx_bot1) listen_account_balance() {
	bot.pub_account_update_listen_id = pub_account_update.Add(func(response bingx_aop_response) {
		bot.addMessage("listen_account_balance", map[string]string{"symbol": bot.symbol, "PosDirection": bot.PosDirection()}, response.Account)
		bot.Balance = utils.StrToFloat64(response.Account.BalanceInfo[0].WalletBalance)

		for _, ti := range response.Account.TradeInfo {
			if ti.Symbol == bot.symbol && ti.PosDirection == bot.PosDirection() {
				bot.Avg_entry = utils.StrToFloat64(ti.EntryPrice)
				bot.Pos_size, _ = strconv.ParseFloat(ti.Position, 64)
				bot.Pos_size = math.Abs(bot.Pos_size)
				bot.addMessage(fmt.Sprintf("BingX listen_account_balance: EntryPrice: %.2f, Pos_size: %.2f, Balance: %.2f", bot.Avg_entry, bot.Pos_size, bot.Balance), bot, response)
				break
			}
		}

	})
}

func (bot *bingx_bot1) FormatPrice(price float64) string {
	return decimal.NewFromFloatWithExponent(price, -bot.num_fractional_digits).String()
}

func (bot *bingx_bot1) stop() {
	bot.bot_started = false

	if bot.pub_account_update_listen_id != 0 {
		pub_account_update.Remove(bot.pub_account_update_listen_id)
		bot.pub_account_update_listen_id = 0
	}
	if bot.pub_filled_order_listen_id != 0 {
		pub_filled_order.Remove(bot.pub_filled_order_listen_id)
		bot.pub_filled_order_listen_id = 0
	}
	db.DB.Delete(&BotPlan{BotID: bot.ID})
}

func (bot *bingx_bot1) start() {
	if bot.bot_started {
		return
	}

	bot.bot_started = true
	bot.listen_account_balance() //gives position entry, size and balance

	//check if this bot was already started
	bot_plan := BotPlan{BotID: bot.ID, OpenOrderID: utils.Guid()}

	if db.DB.First(&bot_plan).RowsAffected == 1 {

		bot.wait_for_fill_one(bot_plan.OpenOrderID, bot_plan.CloseOrderID)
		bot.addMessage("Found existing bot plan, continuing")

	} else {
		//starting bot from scratch
		bot.wait_for_fill_one(bot_plan.OpenOrderID)

		//place market order and get the fill price
		//TODO: get real price instead of approx_start_entry
		market_price, err := get_latest_price_of_trading_pair(bot.symbol)

		if err != nil {
			db.AddError(err, "start "+bot.ID)
			return
		}

		approx_pos_size := bot.pos_size_rule(decimal.NewFromFloat(market_price))
		bot_plan.Description = fmt.Sprintf("%s: market order %.2f @ %.2f", bot.ID, approx_pos_size.InexactFloat64(), market_price)
		bot.Pos_size = approx_pos_size.InexactFloat64()
		bot.addMessage(fmt.Sprintf("Placing BingX market order %.2f @ %.2f", approx_pos_size.InexactFloat64(), market_price), map[string]interface{}{"market_price": market_price, "approx_pos_size": approx_pos_size, "pos_size": bot.Pos_size, "ClientOrderID": bot_plan.OpenOrderID})
		bot.order_api.place_futures_order(bot, swapOrderRequest{
			Symbol:        bot.symbol,
			Side:          bot.BuySellSide(true),
			PositionSide:  bot.PosDirection(),
			Type:          "MARKET",
			Quantity:      approx_pos_size.String(),
			ClientOrderID: bot_plan.OpenOrderID,
			Price:         strconv.FormatFloat(market_price, 'f', -1, 64), //needed for test only
		})

		db.DB.Save(&bot_plan)
	}
}

func Start_Bot1_Long() {
	bot_long.start()
}

func Stop_Bot1_Long() {
	bot_long.stop()
}
func Start_Bot1_Short() {
	bot_short.start()
}

func Stop_Bot1_Short() {
	bot_short.stop()
}

func (bot *bingx_bot1) place_2_orders(open_price decimal.Decimal, close_price decimal.Decimal, open_qty decimal.Decimal, close_qty decimal.Decimal) {
	OpenOrderID := utils.Guid()
	CloseOrderID := utils.Guid()

	bot_plan := BotPlan{
		BotID:        bot.ID,
		OpenOrderID:  OpenOrderID,
		CloseOrderID: CloseOrderID,
		Description:  fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_qty.InexactFloat64(), open_price.String(), close_qty.InexactFloat64(), close_price.String()),
	}
	db.DB.Save(&bot_plan)

	bot.wait_for_fill_one(OpenOrderID, CloseOrderID)

	bot.addMessage("Placing BingX limit orders: " + fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_qty.InexactFloat64(), open_price.String(), close_qty.InexactFloat64(), close_price.String()))

	bot.order_api.place_futures_order(bot, swapOrderRequest{
		Symbol:        bot.symbol,
		Side:          bot.BuySellSide(true),
		PositionSide:  bot.PosDirection(),
		Type:          "LIMIT",
		Price:         open_price.String(),
		Quantity:      open_qty.String(),
		ClientOrderID: OpenOrderID,
	}, swapOrderRequest{
		Symbol:        bot.symbol,
		Side:          bot.BuySellSide(false),
		PositionSide:  bot.PosDirection(),
		Type:          "LIMIT",
		Price:         close_price.String(),
		Quantity:      close_qty.String(),
		ClientOrderID: CloseOrderID,
	})

}

func (bot *bingx_bot1) wait_for_fill_one(ClientOrderIDs ...string) {
	ClientOrderIDs = funk.FilterString(ClientOrderIDs, func(s string) bool { return s != "" })

	bot.pub_filled_order_listen_id = pub_filled_order.Add(func(response bingx_aop_response) { //fires when the limit order is filled
		if !slices.Contains(ClientOrderIDs, response.Order.ClientOrderID) {
			return
		}
		pub_filled_order.Remove(bot.pub_filled_order_listen_id)
		bot.pub_filled_order_listen_id = 0

		//cancel non-filled orders
		for _, ClientOrderID := range ClientOrderIDs {
			if ClientOrderID != response.Order.ClientOrderID {
				bot.order_api.cancel_order(bot.symbol, ClientOrderID)
			}
		}

		fill_price, err := decimal.NewFromString(response.Order.AvgPrice)
		if err != nil {
			db.AddError(err, "bingx order_filled")
			return
		}

		prev, next, prev_index, next_index := bot.closest_loi(fill_price)

		//calc loi for opening
		open_price, open_amt := bot.calc_opening_loi(prev_index)

		//calc loi for closing
		close_price, close_amt := bot.calc_closing_loi(next_index)

		bot.addMessage("calc loi: "+fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_amt.InexactFloat64(), open_price.String(), close_amt.InexactFloat64(), close_price.String()), map[string]interface{}{"fill_price": fill_price, "prev": prev, "next": next, "prev_index": prev_index, "next_index": next_index, "open_price": open_price, "open_amt": open_amt, "close_price": close_price, "close_amt": close_amt, "bot": bot})
		bot.place_2_orders(open_price, close_price, open_amt, close_amt)

	})
}

func (bot *bingx_bot1) calc_opening_loi(lo_index int) (open_price decimal.Decimal, open_amt decimal.Decimal) {
	for i := lo_index; i >= 0; i-- {
		if amt := bot.calc_opening_amount(bot.loi[i]); amt > bot.min_notable_amt {
			open_price = bot.loi[i]
			open_amt = decimal.NewFromFloat(amt)
			return
		} else {
			//println(loi[i].String(), amt)
		}
	}
	return
}

func (bot *bingx_bot1) calc_closing_loi(hi_index int) (close_price decimal.Decimal, close_amt decimal.Decimal) {
	for i := hi_index; i < len(bot.loi); i++ {
		if amt := bot.calc_closing_amount(bot.loi[i]); amt > 0 {
			close_price = bot.loi[i]
			close_amt = decimal.NewFromFloat(amt)
			return
		} else {
			//println(loi[i].String(), amt)
		}
	}
	return
}

func (bot *bingx_bot1) calc_opening_amount(loi decimal.Decimal) float64 {
	if bot.pos_size_rule(loi).InexactFloat64() > bot.Pos_size {
		return bot.pos_size_rule(loi).InexactFloat64() - bot.Pos_size
	} else {
		return 0
	}
}

// 1 for long, -1 for short
func (bot *bingx_bot1) cmp_sign() (cmp_result int) {
	cmp_result = 1
	if !bot.is_long {
		cmp_result = -1
	}
	return
}

func (bot *bingx_bot1) calc_closing_amount(loi decimal.Decimal) float64 {

	if bot.Pos_size > bot.pos_size_rule(loi).InexactFloat64() && loi.Compare(decimal.NewFromFloat(bot.Avg_entry)) == bot.cmp_sign() {
		return bot.Pos_size - bot.pos_size_rule(loi).InexactFloat64()
	} else {
		return 0
	}
}
