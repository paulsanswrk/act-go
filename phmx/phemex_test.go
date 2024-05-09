package phmx

import (
	"context"
	"fmt"
	"github.com/Krisa/go-phemex"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sacOO7/gowebsocket"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

var (
	//apiKey    = "19dbd945-9120-4c81-8b5f-9bee29f12115"
	//secretKey = "1fGvckB1eZrbvINKfA8eZ-tdtVQV_pvTPC4ECHbep0RkYTIzM2IxYy02YmQ1LTQyYzMtYWM1MS01ODYyNzRkMWRjYTA"
	apiKey    = "a68b4aa9-9394-4416-8aa4-213634f19014"
	secretKey = "7HaBhEyHs8UF2xNOXnE87_luLjur-n2Rf71evXt8Z5czODhmYzg4OS00ZWM1LTQzZmMtOTkzOS1iZTlmN2EyM2Q1MzM"
)

func TestWS(t *testing.T) {
	PhemexClient := phemex.NewClient(apiKey, secretKey)

	wsHandler := func(message interface{}) {
		switch message.(type) {
		case *phemex.WsAOP:
			// snapshots / increments
		case *phemex.WsPositionInfo:
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
