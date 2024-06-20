package phmx

import (
	"ACT_GO/db"
	"ACT_GO/db/entities"
	"ACT_GO/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	//apiKey    = "19dbd945-9120-4c81-8b5f-9bee29f12115"
	//secretKey = "1fGvckB1eZrbvINKfA8eZ-tdtVQV_pvTPC4ECHbep0RkYTIzM2IxYy02YmQ1LTQyYzMtYWM1MS01ODYyNzRkMWRjYTA"
	apiKey    = "a68b4aa9-9394-4416-8aa4-213634f19014"
	secretKey = "7HaBhEyHs8UF2xNOXnE87_luLjur-n2Rf71evXt8Z5czODhmYzg4OS00ZWM1LTQzZmMtOTkzOS1iZTlmN2EyM2Q1MzM"
)

func Listen_Account_WS() {
	listen_account_ws()
}

func listen_account_ws() {
	interrupt := make(chan struct{})

	header := http.Header{}
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.phemex.com", header)
	if err != nil {
		//log.Fatal("WebSocket connection error:", err)
		db.AddError(err, "Phemex WebSocket connection error")
		//log.Println("listen_account_ws close interrupt 1")
		close(interrupt)
		return
	}
	defer c.Close()

	//auth
	expiry := time.Now().Unix() + 60
	raw := fmt.Sprintf("%s%v", apiKey, expiry)
	signedString := utils.ComputeHmac256(raw, secretKey)

	err = c.WriteJSON(map[string]interface{}{
		"method": "user.auth",
		"params": []interface{}{
			"API",
			apiKey,
			signedString,
			expiry,
		},
		"id": 100,
	})

	if err != nil {
		db.AddError(err, "Phemex WebSocket auth error")
		close(interrupt)
		return
	}
	_, _, err = c.ReadMessage()
	//db.AddMessage()

	err = c.WriteJSON(map[string]interface{}{
		"id":     0,
		"method": "aop_p.subscribe",
		"params": []interface{}{},
	})
	if err != nil {
		db.AddError(err, "Phemex WebSocket aop_p.subscribe error")
		close(interrupt)
		return
	}

	go func() {
		ping_ticker := time.NewTicker(5 * time.Second)
		msg_id := 1
		defer ping_ticker.Stop()

		for {
			select {
			case <-ping_ticker.C:
				msg_id++
				p, _ := json.Marshal(map[string]interface{}{
					"id":     msg_id,
					"method": "server.ping",
					"params": []string{},
				})
				err := c.WriteControl(websocket.PingMessage, p, time.Time{})
				if err != nil {
					db.AddError(err, "Phemex WebSocket WriteControl error")
					close(interrupt)
					return
				}
				//err := c.WriteMessage(websocket.TextMessage, []byte("Ping"))

			case <-interrupt:
				return
			default:
			}

			messageType, message, err := c.ReadMessage()
			if err != nil {
				db.AddError(err, "Phemex WebSocket read error")
				//log.Println("listen_account_ws close interrupt 3")
				close(interrupt)
				return
			}

			if messageType == websocket.TextMessage { //
				db.Add_Log(&entities.Log{Message: string(message), Tag: "Phemex listen_account_ws TextMessage"})
			} else if messageType == websocket.BinaryMessage {
				db.AddMessage("Phemex BinaryMessage")
			}

		}
	}()

	<-interrupt
	db.AddMessage("Phemex listen_account_ws interrupted")
}
