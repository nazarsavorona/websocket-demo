package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WebSocketServer struct {
	upgrader  websocket.Upgrader
	clients   []*websocket.Conn
	broadcast chan []byte
}

func (wss *WebSocketServer) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(conn)

	wss.clients = append(wss.clients, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Message received: %s", message)
	}
}

func registerInNodeConnector(url string) error {
	_, err := http.Post("http://"+url+"/nodes", "application/json", bytes.NewBuffer([]byte(`{
            "hostname":"localhost:8081",
            "validator_key": [1,2,3]
}`)))
	if err != nil {
		log.Panicln(err)
	}
	return err
}

func main() {
	//url := os.Getenv("NODE_CONNECTOR_URL")
	url := "localhost:8080"
	err := registerInNodeConnector(url)

	wss := WebSocketServer{
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan []byte),
	}

	http.HandleFunc("/ws", wss.HandleWebSocketConnection)
	// ping endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wss.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}

		// set ping handler
		conn.SetPingHandler(func(appData string) error {
			log.Println("Received ping")

			err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(15*time.Second))
			if err != nil {
				log.Println("Error sending pong:", err)
			}

			return nil
		})

		defer func(conn *websocket.Conn) {
			err := conn.Close()
			println("closing connection")
			if err != nil {
				log.Println("Error closing connection11:", err)
			}
		}(conn)

		log.Println("Hello")

		// sync channel

		// read message
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		log.Println("message received:", string(message))
	})

	log.Println("Starting WebSocket server on :8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
