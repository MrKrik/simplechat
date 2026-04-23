package main

import (
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func echoServer(w http.ResponseWriter, r *http.Request) {
	// Принимаем соединение (Upgrade)
	// Options позволяют настроить сжатие или проверку Origin
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Для разработки (аналог CheckOrigin)
	})
	if err != nil {
		log.Printf("Ошибка приёма соединения: %v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	ctx := r.Context()

	for {
		// Чтение сообщения
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Printf("Ошибка чтения: %v", err)
			break
		}

		log.Printf("Получено: %s", message)

	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func main() {
	http.HandleFunc("/ws", echoServer)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Println("Сервер запущен на :8080")
	log.Fatal(server.ListenAndServe())
}
