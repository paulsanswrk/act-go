package phmx

import (
	"ACT_GO/db"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Krisa/go-phemex"
	socketio "github.com/googollee/go-socket.io"
	"github.com/parnurzeal/gorequest"
	"github.com/sacOO7/gowebsocket"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestGoPhemex(t *testing.T) {
	PhemexClient := phemex.NewClient(apiKey, secretKey)

	wsHandler := func(message interface{}) {
		switch message.(type) {
		case *phemex.WsAOP:
			// snapshots / increments
		case *phemex.WsPositionInfo:
			fmt.Printf("%v\n", message.(phemex.WsPositionInfo).PositionInfo.Symbol)
			// when a position is active
		case *phemex.WsError:
			// on connection
		}

		fmt.Printf("%v\n", message)
	}

	errHandler := func(err error) {
		// initiate reconnection with `once.Do...`
		fmt.Println("errHandler: ", err)
	}

	auth := PhemexClient.NewWsAuthService()
	auth = auth.URL("wss://ws.phemex.com")

	conn, err := auth.Do(context.Background())
	// err handling
	if err != nil {
		fmt.Println("auth.Do: ", err)
		return
	}

	err = PhemexClient.NewStartWsAOPService().SetID(1).Do(conn, wsHandler, errHandler)
	if err != nil {
		fmt.Println("NewStartWsAOPService: ", err)
		return
	}
	time.Sleep(60 * time.Second)
}

var id = 5

func inc_id() int {
	id++
	return id
}

func TestGoPhemex2(t *testing.T) {
	client := phemex.NewClient(apiKey, secretKey)

	openOrders, err := client.NewListOpenOrdersService().Symbol("SOLUSD").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, o := range openOrders {
		fmt.Println(o)
	}
}

func signString(raw string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(raw))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func TestRestListContractOrders(t *testing.T) {
	endpoint := "/exchange/order/v2/tradingList"
	//endpoint := "/g-accounts/positions"
	//endpoint := "/api-data/futures/v2/tradeAccountDetail"
	queryString := "currency=USDT&limit=100&offset=0" //&symbol=ENAUSDT
	bodyString := ""
	expiry := fmt.Sprintf("%v", time.Now().Unix()+60)
	raw := fmt.Sprintf("%s%s%s%s", endpoint, queryString, expiry, bodyString)
	signedString, _ := signString(raw)

	_, body, errs := gorequest.New().
		Get("https://api.phemex.com"+endpoint+"?"+queryString).
		AppendHeader("x-phemex-access-token", apiKey).
		AppendHeader("x-phemex-request-expiry", expiry).
		AppendHeader("x-phemex-request-signature", signedString).
		End()

	//fmt.Printf("%v\n", resp)
	fmt.Printf("%v\n", body)
	fmt.Printf("%v\n", errs)
}

func TestRestQueryWallets(t *testing.T) {
	endpoint := "/spot/wallets"
	queryString := "currency=" //USDT
	bodyString := ""
	expiry := fmt.Sprintf("%v", time.Now().Unix()+60)
	raw := fmt.Sprintf("%s%s%s%s", endpoint, queryString, expiry, bodyString)
	signedString, _ := signString(raw)

	_, body, errs := gorequest.New().
		Get("https://api.phemex.com"+endpoint+"?"+queryString).
		AppendHeader("x-phemex-access-token", apiKey).
		AppendHeader("x-phemex-request-expiry", expiry).
		AppendHeader("x-phemex-request-signature", signedString).
		End()

	//fmt.Printf("%v\n", resp)
	fmt.Printf("%v\n", body)
	fmt.Printf("%v\n", errs)
}

func TestWS2(t *testing.T) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("wss://ws.phemex.com/")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")

		socket.SendText(`{
  "id": 2,
  "method": "trade.subscribe",
  "params": [
    "BTCUSDT"
  ]
}`)

		go func() {
			for {
				select {
				case <-ticker.C:
					socket.SendText(fmt.Sprintf(`{
  "id": %d,
  "method": "server.ping",
  "params": []
}`, inc_id()))
				}
			}
		}()

	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Received connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Received message " + message)

	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Received binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}

func TestSocketIO(t *testing.T) {
	// Simple client to talk to default-http example
	uri := "wss://ws.phemex.com/"

	client, err := socketio.NewClient(uri, nil)
	if err != nil {
		panic(err)
	}

	// Handle an incoming event
	client.OnEvent("message", func(s socketio.Conn, msg string) {
		log.Println("Receive Message /reply: ", "reply", msg)
	})

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	client.Emit("kline.subscribe", `{
  "id": 2,
  "method": "kline.subscribe",
  "params": [
    "BTCUSDT",
    60
  ]
}`)

	time.Sleep(15 * time.Second)
	err = client.Close()
	if err != nil {
		panic(err)
	}
}

func TestSocketRaw(t *testing.T) {
	// create connection

	// schema can be ws:// or wss://

	// host, port – WebSocket server

	conn, err := websocket.Dial("wss://ws.phemex.com/", "", "")

	if err != nil {

		// handle error

	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// send message

	if err = websocket.JSON.Send(conn, `{
  "id": 2,
  "method": "kline.subscribe",
  "params": [
    "BTCUSDT",
    60
  ]
}`); err != nil {

		// handle error

	}

	// receive message

	// messageType initializes some type of message

	/*	message := messageType{}

		if err := websocket.JSON.Receive(conn, &message); err != nil {

			// handle error

		}*/
}

func TestWSRaw2(t *testing.T) {
	origin := "http://localhost/"
	url := "wss://ws.phemex.com/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte(`{
  "id": 2,
  "method": "kline.subscribe",
  "params": [
    "BTCUSDT",
    60
  ]
}`)); err != nil {
		log.Fatal(err)
	}

	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])

	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

func TestListen_Account_WS(t *testing.T) {
	db.TruncateLogs()
	listen_account_ws()
}
