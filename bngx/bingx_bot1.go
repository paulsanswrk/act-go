package bngx

import (
	"ACT_GO/db"
	"ACT_GO/utils"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"log"
	"slices"
	"strconv"
)

type bingx_bot1 struct {
	order_api   i_bingx_order_api
	bot_started bool
	symbol      string
	loi         []decimal.Decimal
	Balance     float64
	Avg_entry   float64
	Pos_size    float64
	pMin        float64
	pMax        float64
}

var (
	bot1 = bingx_bot1{
		order_api:   &bingx_order_api{},
		bot_started: false,
		symbol:      "WIF-USDT",
		Balance:     43.0,
		Avg_entry:   2.2,
		Pos_size:    0,
		pMin:        0.25,
		pMax:        5.25,
	}
)

func init() {
	var s_loi = []string{"0", "0.25", "0.5", "0.75", "1", "1.25", "1.5", "1.75", "2", "2.25", "2.5", "2.75", "3", "3.25", "3.5", "3.75", "4", "4.25", "4.5", "4.75", "5", "5.25", "5.5", "5.75", "6", "6.25", "6.5", "6.75", "7", "7.25"}
	bot1.loi = funk.Map(s_loi, func(s string) (p decimal.Decimal) { p, _ = decimal.NewFromString(s); return }).([]decimal.Decimal)
}

func (bot *bingx_bot1) pos_size_rule(price decimal.Decimal) decimal.Decimal {
	//pMin := 0.25
	//pMax := 5.25
	result := 1 - (price.InexactFloat64()-bot.pMin)/(bot.pMax-bot.pMin)
	result *= 10
	return decimal.NewFromFloat(result)
}

func (bot *bingx_bot1) closest_loi(price decimal.Decimal) (lo decimal.Decimal, hi decimal.Decimal, lo_index int, hi_index int) {
	for i, p := range bot.loi {
		if p.Compare(price) == 1 { //p > price
			return bot.loi[i-1], p, i - 1, i
		}
	}
	log.Fatalf("closest_loi failed for %v", price)
	return
}

func (bot *bingx_bot1) listen_account_balance() {
	pub_account_update.Add(func(response bingx_aop_response) {
		bot.Balance = response.Account.BalanceInfo.WalletBalance
		if response.Account.TradeInfo.Symbol == bot.symbol {
			bot.Avg_entry = response.Account.TradeInfo.EntryPrice
			bot.Pos_size, _ = strconv.ParseFloat(response.Account.TradeInfo.Position, 64)
		}
		db.AddMessage(fmt.Sprintf("BingX listen_account_balance: EntryPrice: %.2f, Pos_size: %.2f", bot.Avg_entry, bot.Pos_size), nil, bot)
	})
}

func (bot *bingx_bot1) start() {
	bot.bot_started = true
	bot.listen_account_balance()

	ClientOrderID := utils.Guid()

	bot.wait_for_fill_one(ClientOrderID)

	/*	var nSub uint64
		nSub = pub_filled_order.Add(func(response bingx_aop_response) { //fires when the market order is filled
			if response.Order.ClientOrderID != ClientOrderID {
				db.AddMessage("BingX Listen Market Order", ClientOrderID, response)
				return
			}
			pub_filled_order.Remove(nSub)

			fill_price, err := decimal.NewFromString(response.Order.AvgPrice)
			if err != nil {
				db.AddError(err, "bingx order_filled")
				return
			}

			db.AddMessage("bingx market order_filled")
			lo, hi, _, _ := bot.closest_loi(fill_price)
			go bot.place_2_orders(lo, hi, decimal.NewFromInt(1), decimal.NewFromInt(1))

		})*/

	//place market order and get the fill price
	approx_price := 2.15
	approx_pos_size := bot.pos_size_rule(decimal.NewFromFloat(approx_price))
	bot.Pos_size = approx_pos_size.InexactFloat64()
	db.AddMessage(fmt.Sprintf("Placing BingX market order %.2f @ %.2f", approx_pos_size.InexactFloat64(), approx_price), map[string]interface{}{"approx_price": approx_price, "approx_pos_size": approx_pos_size, "pos_size": bot.Pos_size, "ClientOrderID": ClientOrderID})
	bot.order_api.place_futures_order(swapOrderRequest{
		Symbol:        bot.symbol,
		Side:          "BUY",
		PositionSide:  "LONG",
		Type:          "MARKET",
		Quantity:      approx_pos_size.String(),
		ClientOrderID: ClientOrderID,
		Price:         strconv.FormatFloat(approx_price, 'f', -1, 64), //needed for test only
	})
}

func Start_Bot1() {
	if bot1.bot_started {
		return
	}
	bot1.start()
}

func (bot *bingx_bot1) bot_id() string {
	return "BingX_Bot_1"
}

func (bot *bingx_bot1) place_2_orders(open_price decimal.Decimal, close_price decimal.Decimal, open_qty decimal.Decimal, close_qty decimal.Decimal) {
	OpenOrderID := utils.Guid()
	CloseOrderID := utils.Guid()

	bot_plan := BotPlan{
		BotID:        bot.bot_id(),
		OpenOrderID:  OpenOrderID,
		CloseOrderID: CloseOrderID,
		Description:  fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_qty.InexactFloat64(), open_price.String(), close_qty.InexactFloat64(), close_price.String()),
	}
	db.DB.Save(&bot_plan)

	bot.wait_for_fill_one(OpenOrderID, CloseOrderID)

	db.AddMessage("Placing BingX limit orders: " + fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_qty.InexactFloat64(), open_price.String(), close_qty.InexactFloat64(), close_price.String()))

	bot.order_api.place_futures_order(swapOrderRequest{
		Symbol:        bot.symbol,
		Side:          "BUY",
		PositionSide:  "LONG",
		Type:          "LIMIT",
		Price:         open_price.String(),
		Quantity:      open_qty.String(),
		ClientOrderID: OpenOrderID,
	}, swapOrderRequest{
		Symbol:        bot.symbol,
		Side:          "SELL",
		PositionSide:  "LONG",
		Type:          "LIMIT",
		Price:         close_price.String(),
		Quantity:      close_qty.String(),
		ClientOrderID: CloseOrderID,
	})

}

func (bot *bingx_bot1) wait_for_fill_one(ClientOrderIDs ...string) {
	var nSub uint64
	nSub = pub_filled_order.Add(func(response bingx_aop_response) { //fires when the limit order is filled
		if !slices.Contains(ClientOrderIDs, response.Order.ClientOrderID) {
			return
		}
		pub_filled_order.Remove(nSub)

		bot.order_api.cancel_all_orders(bot.symbol)

		fill_price, err := decimal.NewFromString(response.Order.AvgPrice)
		if err != nil {
			db.AddError(err, "bingx order_filled")
			return
		}

		lo, hi, lo_index, hi_index := bot.closest_loi(fill_price)

		//calc loi for opening
		open_price, open_amt := bot.calc_opening_loi(lo_index)

		//calc loi for closing
		close_price, close_amt := bot.calc_closing_loi(hi_index)

		db.AddMessage("calc loi: "+fmt.Sprintf("Open %.2f @ %s, close %.2f @ %s", open_amt.InexactFloat64(), open_price.String(), close_amt.InexactFloat64(), close_price.String()), map[string]interface{}{"fill_price": fill_price, "lo": lo, "hi": hi, "lo_index": lo_index, "hi_index": hi_index, "open_price": open_price, "open_amt": open_amt, "close_price": close_price, "close_amt": close_amt, "bot": bot})
		bot.place_2_orders(open_price, close_price, open_amt, close_amt)

	})
}

func (bot *bingx_bot1) calc_opening_loi(lo_index int) (open_price decimal.Decimal, open_amt decimal.Decimal) {
	for i := lo_index; i >= 0; i-- {
		if amt := bot.calc_opening_amount(bot.loi[i]); amt > bot.Pos_size/100 {
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

func (bot *bingx_bot1) calc_closing_amount(loi decimal.Decimal) float64 {
	if bot.Pos_size > bot.pos_size_rule(loi).InexactFloat64() && loi.InexactFloat64() > bot.Avg_entry {
		return bot.Pos_size - bot.pos_size_rule(loi).InexactFloat64()
	} else {
		return 0
	}
}
