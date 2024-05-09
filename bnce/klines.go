package bnce

import (
	"ACT_GO/db"
	"ACT_GO/db/entities"
	. "ACT_GO/utils"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/thoas/go-funk"
	"log"
)

var (
	listen_interval       = "1m"
	listening_symbols_map = funk.Map(all_symbols, func(sym string) (string, string) {
		return sym, listen_interval
	}).(map[string]string)
)

func Listen_Binance_Klines() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in Listen_Binance_Klines:", r)
		}
	}()

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in wsKlineHandler:", r)
			}
		}()

		if !event.Kline.IsFinal {
			return
		}

		var price = &entities.Price{
			Symbol:        event.Symbol,
			Time:          event.Kline.StartTime,
			Open:          StrToFloat64(event.Kline.Open),
			High:          StrToFloat64(event.Kline.High),
			Low:           StrToFloat64(event.Kline.Low),
			Close:         StrToFloat64(event.Kline.Close),
			Volume:        StrToFloat64(event.Kline.Volume),
			WeightedPrice: (StrToFloat64(event.Kline.Low) + StrToFloat64(event.Kline.High)) / 2,
		}
		db.DB.Create(price)

	}
	errHandler := func(err error) {
		log.Println(err)
	}
	log.Println("Listen_Binance_Klines started")
	doneC, _, err := binance.WsCombinedKlineServe(listening_symbols_map, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
	log.Println("Listen_Binance_Klines finished")
}
