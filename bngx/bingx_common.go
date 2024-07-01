package bngx

import (
	"ACT_GO/utils"
	"fmt"
	"github.com/thoas/go-funk"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	apiKey    = "qTHuDyd8KB5ATOH7zgHyA8SoLa1CcyHJwQHwURclYW8oeJ7n2Af71H78wBulrhmkSsUgitzRiDPCJ2LYSfQ"
	secretKey = "f4efoSS2Rhl8RScD1tyG0EeTvJR19rR3risqHrbx34KTpSnaiZgJy9VhjnBg7ATbtgg04aJJx0aEwDTB6MHw"
)

// assuming that params map doesn't need URL encoding
func build_and_sign_url(params map[string]interface{}) string {

	keys := funk.Keys(params).([]string)
	slices.Sort(keys)

	timestamp := time.Now().UnixNano() / 1e6
	//timestamp = 1716810147094
	parts := funk.Map(keys, func(k string) string { return fmt.Sprintf("%s=%v", k, params[k]) }).([]string)
	url := strings.Join(parts, "&")
	url = fmt.Sprintf("%s&timestamp=%d", url, timestamp)
	sign := utils.ComputeHmac256(fmt.Sprintf("%v", url), secretKey)
	url = fmt.Sprintf("%s&signature=%s", url, sign)

	return url
}

func get_latest_price_of_trading_pair(pair string) (price float64, err error) {
	url := fmt.Sprintf("https://open-api.bingx.com/openApi/swap/v1/ticker/price?symbol=%s&timestamp=%d", pair, time.Now().UnixNano()/1e6)
	var res bingx_ticker_price_response
	_, err = utils.HTTP_Request(url, "GET", nil, &res)

	if err != nil {
		return
	}

	price, err = strconv.ParseFloat(res.Data.Price, 64)

	return
}
