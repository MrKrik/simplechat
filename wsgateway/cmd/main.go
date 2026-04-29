package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"wsgateway/iternal/hub"
	chub "wsgateway/iternal/hub"

	"github.com/coder/websocket"
)

func main() {
	hub := chub.NewHub()

	hub.Rooms["general"] = []string{"user1", "user2"}

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New user", r.RemoteAddr, r.URL.Query().Get("userId"))
	client := &chub.Client{
		UserID:     r.URL.Query().Get("userId"), // Передаем в URL: ?userId=user1
		ConnID:     string(rand.Int31()),
		Conn:       conn,
		Send:       make(chan []byte, 256),
		Hub:        hub,
		LastActive: time.Now(),
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
