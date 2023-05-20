package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
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
	_, err := http.Post("https://"+url+"/nodes", "application/json", bytes.NewBuffer([]byte(`{"ip":"`+os.Getenv("NODE_URL")+`","port":"8081"}`)))
	if err != nil {
		log.Panicln(err)
	}
	return err
}

func main() {
	url := os.Getenv("NODE_CONNECTOR_URL")
	err := registerInNodeConnector(url)

	wss := WebSocketServer{
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan []byte),
	}

	http.HandleFunc("/ws", wss.HandleWebSocketConnection)

	log.Println("Starting WebSocket server on :8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
