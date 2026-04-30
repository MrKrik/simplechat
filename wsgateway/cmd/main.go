package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"wsgateway/iternal/grpc"
	"wsgateway/iternal/hub"
	chub "wsgateway/iternal/hub"

	"github.com/coder/websocket"
)

func main() {
	hub := chub.NewHub()

	grpcClient, err := grpc.New(":44044", time.Duration(time.Duration.Minutes(1)))
	if err != nil {
		log.Println("failed start grpc client")
	}

	hub.Rooms["general"] = []string{"user1", "user2"}

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, grpcClient)
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request, grpc *grpc.Client) error {

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	token := r.URL.Query().Get("token")

	ok, err := grpc.ValidateToken(token)
	if !ok {
		return err
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

	return nil
}
