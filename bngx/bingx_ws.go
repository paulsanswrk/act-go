package bngx

import (
	"ACT_GO/db"
	"ACT_GO/db/entities"
	"ACT_GO/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"time"
)

type key_wrap struct {
	ListenKey string `json:"listenKey"`
}

var (
	listen_key         string
	pub_filled_order   utils.Publisher[bingx_aop_response]
	pub_account_update utils.Publisher[bingx_aop_response]
)

func Listen_Account_WS() {
	last_start := time.Now()
	n_soon_retry := 0

	//a loop for restarting websocket connection
	for {
		get_listen_key()
		listen_account_ws()

		if time.Now().Before(last_start.Add(60 * time.Second)) { //don't restart if stopped too soon
			n_soon_retry++
			if n_soon_retry > 50 {
				db.AddMessage(fmt.Sprintf("Listen_Account_WS: stopped retrying after %d attempts", n_soon_retry))
				//log.Println("stopped")
				break
			} else if n_soon_retry > 20 {
				time.Sleep(1000 * time.Millisecond)
			} else if n_soon_retry > 5 {
				time.Sleep(100 * time.Millisecond)
			}
		}
		last_start = time.Now()
	}
}

// subscribe to websocket and listen any account changes (orders, transfers, funding fees)
func listen_account_ws() {
	interrupt := make(chan struct{}) //stop on error and then restart in the caller

	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	//var err error
	conn, _, err := websocket.DefaultDialer.Dial("wss://open-api-swap.bingx.com/swap-market?listenKey="+listen_key, header)
	if err != nil {
		//log.Fatal("WebSocket connection error:", err)
		db.AddError(err, "BingX WebSocket connection error")
		log.Println("listen_account_ws close interrupt 1")
		close(interrupt)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("***"))
	if err != nil {
		//log.Fatal("WebSocket write error:", err)
		db.AddError(err, "BingX WebSocket write error")
		log.Println("listen_account_ws close interrupt 2")
		close(interrupt)
	}

	//run in parallel websocket listener and listen key renewer
	go func() {
		ping_ticker := time.NewTicker(5 * time.Second)
		defer ping_ticker.Stop()

		//message handling loop
		for {
			select { //ping/interrupt/read
			case <-ping_ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
				if err != nil {
					//log.Println("WebSocket write error:", err)
					db.AddError(err, "BingX WebSocket write error")
				}
			case <-interrupt:
				return
			default:
			}

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				//log.Println("WebSocket read error:", err)
				db.AddError(err, "BingX WebSocket read error")
				log.Println("listen_account_ws close interrupt 3")
				close(interrupt)
				return
			}

			if messageType == websocket.TextMessage { //never happens
				db.Add_Log(&entities.Log{Message: string(message), Tag: "listen_account_ws TextMessage"})
			} else if messageType == websocket.BinaryMessage {
				//
				decodedMsg, err := utils.DecodeGzip(message)
				if err != nil {
					//log.Println("WebSocket decode error:", err)
					db.AddError(err, "BingX WebSocket decode error")
					continue
				}

				//fmt.Println(decodedMsg)
				if decodedMsg == "Ping" {
					//log.Println("listen_account_ws ping received")
					err = conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
					if err != nil {
						//log.Println("WebSocket write error:", err)
						db.AddError(err, "BingX WebSocket write error")
						log.Println("listen_account_ws close interrupt 4")
						close(interrupt)
						return
					}
				} else if decodedMsg == "Pong" {
					//log.Println("listen_account_ws pong received")
					//nothing to do
				} else if json.Valid([]byte(decodedMsg)) {
					var aop_response bingx_aop_response
					err := json.Unmarshal([]byte(decodedMsg), &aop_response)
					if err != nil {
						db.AddError(err, "BingX WebSocket unmarshal error", decodedMsg, aop_response)
						continue
					}

					switch aop_response.EventType {
					case "listenKeyExpired":
						log.Println("listen_account_ws close interrupt 5")
						close(interrupt)
						return
					case "ORDER_TRADE_UPDATE":
						if aop_response.Order.Status == "FILLED" {
							pub_filled_order.RunAll(aop_response)
						}
						db.AddMessage("BingX ORDER_TRADE_UPDATE", nil, decodedMsg)
					case "ACCOUNT_UPDATE":
						run_count := pub_account_update.RunAll(aop_response)
						//if aop_response.Account.EventLaunchReason == "ORDER" {}
						db.AddMessage(fmt.Sprintf("BingX ACCOUNT_UPDATE: run_count: %d", run_count), decodedMsg, aop_response)
					case "ACCOUNT_CONFIG_UPDATE":
					default:
						db.Add_Log(&entities.Log{Message: string(decodedMsg), Tag: "listen_account_ws json BinaryMessage"})
					}
				} else {
					db.Add_Log(&entities.Log{Message: string(decodedMsg), Tag: "listen_account_ws non-json BinaryMessage"})
				}
			}
		}
	}()

	//renew listen key periodically
	go func() {
		listen_key_ticker := time.NewTicker(18 * time.Minute)
		defer listen_key_ticker.Stop()

		for {
			select {
			case <-listen_key_ticker.C:
				resp, body, errs := gorequest.New().
					Put("https://open-api.bingx.com/openApi/user/auth/userDataStream").
					Send(key_wrap{ListenKey: listen_key}).
					End()

				if errs != nil {
					db.AddErrors(errs)
				} else {
					db.Add_Log(&entities.Log{Message: "extend_listen_key success, response status: " + resp.Status, Response: body})
				}
			case <-interrupt:
				log.Println("listen_account_ws interrupt")
				return
			}
		}

	}()

	<-interrupt
	db.Add_Log(&entities.Log{Category: entities.LogWarning, Message: "BingX listen_account_ws interrupted"})
}

func get_listen_key() {
	var kw key_wrap

	_, _, errs := gorequest.New().Post("https://open-api.bingx.com/openApi/user/auth/userDataStream").
		AppendHeader("X-BX-APIKEY", apiKey).
		EndStruct(&kw)

	if errs != nil {
		db.AddErrors(errs)
	} else {
		db.Add_Log(&entities.Log{Message: "get_listen_key success"})
	}

	listen_key = kw.ListenKey
}
