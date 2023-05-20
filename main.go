package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
	defer conn.Close()

	// Add the new client connection to the clients list
	wss.clients = append(wss.clients, conn)

	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Broadcast the message to all connected clients
		wss.broadcast <- message
	}
}

func (wss *WebSocketServer) StartBroadcasting() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		log.Println("Waiting for messages...")
		select {
		case message := <-wss.broadcast:
			log.Println("Broadcasting message...")
			// Send the message to all connected clients
			for _, client := range wss.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error writing message:", err)
					return
				}
			}
		case <-ticker.C:
			log.Println("Default message")
			// Broadcast a predefined message every 10 seconds
			broadcastMessage := []byte("Broadcast message every 10 seconds")
			for _, client := range wss.clients {
				err := client.WriteMessage(websocket.TextMessage, broadcastMessage)
				if err != nil {
					log.Println("Error writing message:", err)
					return
				}
			}
		}
	}

	log.Println("Closing WebSocket server...")
}

func main() {
	wss := WebSocketServer{
		upgrader:  websocket.Upgrader{},
		broadcast: make(chan []byte),
	}

	http.HandleFunc("/ws", wss.HandleWebSocketConnection)

	go wss.StartBroadcasting()

	log.Println("Starting WebSocket server on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
