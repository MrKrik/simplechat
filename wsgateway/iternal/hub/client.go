package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Client struct {
	UserID     string
	ConnID     string // UUID для этого соединения (один user может иметь несколько устройств)
	Conn       *websocket.Conn
	Send       chan []byte
	Hub        *Hub
	LastPong   time.Time
	LastActive time.Time
	UserAgent  string
	Ip         string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB max message

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		var msg ClientMessage

		err := wsjson.Read(ctx, c.Conn, &msg)

		cancel()

		if err != nil {
			break
		}

		c.LastActive = time.Now()
		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.Close(websocket.StatusNormalClosure, "channel closed")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			err := c.Conn.Write(ctx, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}

		case <-ticker.C:
			// Send ping for continuation of connection
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := c.Conn.Ping(ctx); err != nil {
				log.Printf("Ping failed: %v", err)
				return
			}

		}
	}
}

func (c *Client) handleMessage(msg ClientMessage) {
	// TODO last or not
	c.LastActive = time.Now()
	switch msg.Type {
	case "chat_message":

		var data struct {
			Text   string `json:"text"`
			RoomID string `json:"room_id"`
		}
		if err := json.Unmarshal(msg.Payload, &data); err != nil {
			return
		}

		outbound := MessageToBroadcast{
			From:   c.UserID,
			Text:   data.Text,
			RoomID: data.RoomID,
			time:   time.Now(),
		}

		// Send
		c.Hub.broadcast <- outbound

	case "typing":
		// Todo

	case "join_room":
		// Todo
		// Join chat
		// c.hub.subscribe <- Subscription{Client: c, Room: "general"}

	default:
		log.Printf("Unkown message type: %s", msg.Type)
	}
}
