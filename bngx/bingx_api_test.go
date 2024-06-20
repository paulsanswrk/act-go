package bngx

import (
	"ACT_GO/utils"
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestPlaceTestOrder(t *testing.T) {

	dataStr := `{
    "uri": "/openApi/swap/v2/trade/order/test",
    "method": "POST",
    "protocol": "https"
}`
	payload := `{
    "symbol": "BTC-USDT",
    "side": "BUY",
    "positionSide": "LONG",
    "type": "LIMIT",
    "quantity": 5,
    "price": 60000
}`
	//"takeProfit": "{\"type\": \"TAKE_PROFIT_MARKET\", \"stopPrice\": 71968.0,\"price\": 71968.0,\"workingType\":\"MARK_PRICE\"}"

	TIMESTAMP := time.Now().UnixNano() / 1e6
	apiMap := getParameters(dataStr, payload, false, TIMESTAMP)
	sign := utils.ComputeHmac256(fmt.Sprintf("%v", apiMap["parameters"]), API_SECRET)
	fmt.Println("parameters:", fmt.Sprintf("%v", apiMap["parameters"]))
	fmt.Println("sign:", sign)
	parameters := ""
	contains := strings.ContainsAny(fmt.Sprintf("%v", apiMap["parameters"]), "[{")
	if contains {
		apiMap2 := getParameters(dataStr, payload, true, TIMESTAMP)
		parameters = fmt.Sprintf("%v&signature=%s", apiMap2["parameters"], sign)
	} else {
		parameters = fmt.Sprintf("%v&signature=%s", apiMap["parameters"], sign)
	}
	url := fmt.Sprintf("%v://%s%v?%s", apiMap["protocol"], HOST, apiMap["uri"], parameters)
	method := fmt.Sprintf("%v", apiMap["method"])
	client := &http.Client{}
	fmt.Println("url:", url)
	fmt.Println("method:", method)
	fmt.Println("apiMap[\"parameters\"]:", apiMap["parameters"])
	fmt.Println("TIMESTAMP:", TIMESTAMP)
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("X-BX-APIKEY", API_KEY)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func getParameters(dataStr string, payload string, urlEncode bool, timestemp int64) map[string]interface{} {

	var apiMap map[string]interface{}
	var payloadMap map[string]interface{}
	err := json.Unmarshal([]byte(dataStr), &apiMap)
	if err != nil {
		fmt.Printf("json to map error,err:%s", err)
		return apiMap
	}
	err = json.Unmarshal([]byte(payload), &payloadMap)
	if err != nil {
		fmt.Printf("json to map error,err:%s", err)
		return apiMap
	}

	//changed the example from docs
	keys := funk.Keys(payloadMap).([]string)
	slices.Sort(keys)

	parameters := ""
	for _, key := range keys {
		value := payloadMap[key]
		if urlEncode {
			encodedStr := url.QueryEscape(fmt.Sprintf("%v", value))
			encodedStr = strings.ReplaceAll(encodedStr, "+", "%20")
			parameters = parameters + key + "=" + encodedStr + "&"
		} else {
			parameters = parameters + key + "=" + fmt.Sprintf("%v", value) + "&"
		}
	}
	parameters += "timestamp=" + fmt.Sprintf("%d", timestemp)
	apiMap["parameters"] = fmt.Sprintf("%v", parameters)
	return apiMap
}
