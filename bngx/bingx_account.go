package bngx

import (
	"ACT_GO/db"
	"ACT_GO/db/entities"
	"ACT_GO/utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"time"
)

var (
	listen_key string
)

type key_wrap struct {
	ListenKey string `json:"listenKey"`
}

type bingx_response struct {
	EventType string `json:"e"`
	EventTime uint64 `json:"E"`
}

func Listen_Account_WS() {
	last_start := time.Now()
	get_listen_key()

	for {
		listen_account_ws()

		if time.Now().Before(last_start.Add(60 * time.Second)) { //don't restart if stopped too soon
			break
		}
		last_start = time.Now()
	}
}

func listen_account_ws() {
	interrupt := make(chan struct{})
	var resp_struct bingx_response

	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	//var err error
	conn, _, err := websocket.DefaultDialer.Dial("wss://open-api-swap.bingx.com/swap-market?listenKey="+listen_key, header)
	if err != nil {
		//log.Fatal("WebSocket connection error:", err)
		db.AddError(err, "WebSocket connection error")
		log.Println("listen_account_ws close interrupt")
		close(interrupt)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("***"))
	if err != nil {
		//log.Fatal("WebSocket write error:", err)
		db.AddError(err, "WebSocket write error")
		log.Println("listen_account_ws close interrupt")
		close(interrupt)
	}

	go func() {
		ping_ticker := time.NewTicker(5 * time.Second)
		defer ping_ticker.Stop()

		for {
			select {
			case <-ping_ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
				if err != nil {
					//log.Println("WebSocket write error:", err)
					db.AddError(err, "WebSocket write error")
				}
			default:
			}

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				//log.Println("WebSocket read error:", err)
				db.AddError(err, "WebSocket read error")
				log.Println("listen_account_ws close interrupt")
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
					db.AddError(err, "WebSocket decode error")
					continue
				}

				//fmt.Println(decodedMsg)
				if decodedMsg == "Ping" {
					//log.Println("listen_account_ws ping received")
					err = conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
					if err != nil {
						//log.Println("WebSocket write error:", err)
						db.AddError(err, "WebSocket write error")
						log.Println("listen_account_ws close interrupt")
						close(interrupt)
					}
				} else if decodedMsg == "Pong" {
					//log.Println("listen_account_ws pong received")
					//nothing to do
				} else if json.Valid([]byte(decodedMsg)) {
					json.Unmarshal([]byte(decodedMsg), &resp_struct)

					switch resp_struct.EventType {
					case "listenKeyExpired":
						log.Println("listen_account_ws close interrupt")
						close(interrupt)
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

	/*	go func() {
		ping_ticker := time.NewTicker(5 * time.Second)
		defer ping_ticker.Stop()

		for {
			select {
			case <-ping_ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
				if err != nil {
					//log.Println("WebSocket write error:", err)
					db.AddError(err, "WebSocket write error")
				}
			case <-interrupt:
				return
			}
		}
	}()*/

	go func() {
		listen_key_ticker := time.NewTicker(30 * time.Minute)
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
	db.Add_Log(&entities.Log{Category: entities.LogWarning, Message: "listen_account_ws interrupted"})
}

func get_listen_key() {
	request := gorequest.New() //.Timeout(100 * time.Minute)

	var kw key_wrap

	_, _, errs := request.Post("https://open-api.bingx.com/openApi/user/auth/userDataStream").
		AppendHeader("X-BX-APIKEY", apiKey).
		EndStruct(&kw)

	if errs != nil {
		db.AddErrors(errs)
	} else {
		db.Add_Log(&entities.Log{Message: "get_listen_key success"})
	}

	listen_key = kw.ListenKey
}
