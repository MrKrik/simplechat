package main

import (
	"encoding/json"
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

	grpcClient, err := grpc.New(":44044", time.Duration(10*time.Second))
	if err != nil {
		log.Println("failed start grpc client")
	}

	hub.Rooms["general"] = []string{"user1", "user2"}

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		errMSG := serveWs(hub, w, r, grpcClient)
		if errMSG != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errMSG)
			return
		}
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request, grpc *grpc.Client) string {

	token := r.URL.Query().Get("token")

	log.Println("Try validate token")

	ok, errMSG := grpc.ValidateToken(token)
	if !ok {
		log.Println("Validate token failed")
		return errMSG
	}
	log.Println("Try validate successfully")

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	log.Println("New user", r.RemoteAddr, r.URL.Query().Get("login"))
	client := &chub.Client{
		UserID:     r.URL.Query().Get("login"), // Передаем в URL: ?login=user1
		ConnID:     string(rand.Int31()),
		Conn:       conn,
		Send:       make(chan []byte, 256),
		Hub:        hub,
		LastActive: time.Now(),
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
	return ""
}
