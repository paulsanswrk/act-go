package bngx

import (
	"ACT_GO/utils"
	"fmt"
	"github.com/sacOO7/gowebsocket"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	path    = "wss://open-api-swap.bingx.com/swap-market?listenKey="                                                        //
	channel = `{"notice:":"no need to subscribe to  any specific channel,please check the hightlight msg in the api docs"}` //
)

var receivedMessage string
var conn *websocket.Conn

func TestAccount(t *testing.T) {
	interrupt := make(chan struct{})
	get_listen_key()

	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	var err error
	conn, _, err = websocket.DefaultDialer.Dial(path+listen_key, header)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(channel))
	if err != nil {
		log.Fatal("WebSocket write error:", err)
	}

	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				close(interrupt)
				return
			}

			handleMessage(messageType, message)
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
				if err != nil {
					log.Println("WebSocket write error:", err)
				}
			case <-interrupt:
				return
			}
		}
	}()

	<-interrupt
}

func handleMessage(messageType int, message []byte) {
	if messageType == websocket.TextMessage {
		//
		fmt.Println(string(message))
	} else if messageType == websocket.BinaryMessage {
		//
		decodedMsg, err := utils.DecodeGzip(message)
		if err != nil {
			log.Println("WebSocket decode error:", err)
			return
		}
		fmt.Println(decodedMsg)
		if decodedMsg == "Ping" {
			data := []byte("Pong")
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
			fmt.Println("Pong")
		}
		receivedMessage = decodedMsg
	}
}

func TestAccount2(t *testing.T) { //no auto response decompression in gowebsocket
	interrupt := make(chan struct{})
	get_listen_key()

	socket := gowebsocket.New(path + listen_key)

	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		UseSSL:         true,
		UseCompression: true,
		//Subprotocols: [] string{"chat","superchat"},
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Received connect error ", err)
		close(interrupt)
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
		close(interrupt)
	}

	socket.Connect()

	<-interrupt
	log.Println("interrupt")
	socket.Close()

}
