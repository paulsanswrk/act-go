package bnce

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/thoas/go-funk"
	"testing"
	"time"
)

var client = binance.NewClient(apiKey, secretKey)

func TestKlines(t *testing.T) {
	klines, err := client.NewKlinesService().Symbol("LTCBTC").
		Interval("1m").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, k := range klines {
		fmt.Println(k)
	}
}

func TestWS_Klines(t *testing.T) {
	var stopC chan struct{}

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event)
		stopC <- struct{}{}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, stopC, err := binance.WsKlineServe("LTCBTC", "1m", wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-stopC
}

func TestWS_Klines_Multi(t *testing.T) {
	var stopC chan struct{}

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		fmt.Println(event)
		//stopC <- struct{}{}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, stopC, err := binance.WsCombinedKlineServe(map[string]string{"LTCBTC": "1m", "ETHBTC": "1m"}, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	var ticker = time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	<-ticker.C
	stopC <- struct{}{}
	//<-stopC
}

func TestMap(t *testing.T) {
	var listening_symbols_map = funk.Map(all_symbols, func(sym string) (string, string) {
		return sym, listen_interval
	})
	fmt.Println(listening_symbols_map)

}

func TestListen_Binance_Klines(t *testing.T) {
	go Listen_Binance_Klines()
	time.Sleep(100 * time.Hour)
}
