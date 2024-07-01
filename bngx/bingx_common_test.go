package bngx

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	params := map[string]interface{}{"b": 1, "a": "qqq"}
	res := build_and_sign_url(params)
	fmt.Println(res)
}

func TestMapStructure(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n\n", result)

	input = map[string]interface{}{}
	err = mapstructure.Decode(result, &input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", input)
}

func TestJsonUnmarshal(t *testing.T) {
	s := `{"e":"ACCOUNT_UPDATE","E":1719682643796,"a":{"m":"ORDER","B":[{"a":"USDT","wb":"211.61440542","cw":"199.96822542","bc":"0"}],"P":[{"s":"WIF-USDT","pa":"108.00000000","ep":"2.16070000","up":"0.06417210","mt":"cross","iw":"11.73195210","ps":"LONG"}]}}`
	var aop_response bingx_aop_response

	err := json.Unmarshal([]byte(s), &aop_response)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%#v\n", aop_response)
}

func TestGetLatestPriceOfATradingPair(t *testing.T) {

	price, err := get_latest_price_of_trading_pair("BTC-USDT")

	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", price)
}
