package hub

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Message struct {
	Author string
	Text   string
}

type Client struct {
	Connection  *websocket.Conn
	messageChan chan Message
	token       string
	chat_token  string
	Room_id     string
}

type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatPayload struct {
	Room_id string `json:"room_id"`
	Text    string `json:"text"`
}

func NewClient() *Client {
	return &Client{
		messageChan: make(chan Message, 50),
		Room_id:     "general",
	}
}

func (c *Client) Start() error {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := c.Connect(ctx)

	if err != nil {
		return err
	}

	go c.handleKeyboard(ctx)

	go c.handleMessage(ctx)

	defer func() {
		cancel()
	}()

	c.Stop(ctx)

	return nil
}

func (c *Client) Connect(parentCtx context.Context) error {

	var err error

	for {
		if c.token == "" {
			c.SetToken()
		} else {
			c.SetChatToken()
			break
		}
	}

	var userID string
	fmt.Scan(&userID)

	var url string
	url = fmt.Sprintf("ws://localhost:8080/ws?token=%s&userID=%v", c.chat_token, userID)

	dialCtx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()

	c.Connection, _, err = websocket.Dial(dialCtx, url, nil)
	if err != nil {
		fmt.Println("Authorization failed:", err)
		return err
	}

	return nil
}

func (c *Client) Stop(ctx context.Context) {
	<-ctx.Done()
	log.Println("Exit")
	c.Connection.Close(websocket.StatusNormalClosure, "work end")
}

func (c *Client) handleKeyboard(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		c.writeInConnection(text, ctx)
	}
}

func (c *Client) handleMessage(ctx context.Context) {
	for {
		_, message, err := c.Connection.Read(ctx)

		if err != nil {
			if ctx.Err() != nil || websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			log.Printf("Connection error: %v", err)
			return
		}

		message = bytes.TrimSpace(message)

		if len(message) > 0 {
			fmt.Println(string(message))
		}
	}
}

func (c *Client) writeInConnection(message string, ctx context.Context) {

	payloadData := ChatPayload{
		Room_id: "general",
		Text:    message,
	}

	payloadBytes, _ := json.Marshal(payloadData)

	msg := ClientMessage{
		Type:    "chat_message",
		Payload: json.RawMessage(payloadBytes),
	}
	err := wsjson.Write(ctx, c.Connection, msg)
	if err != nil {
		log.Printf("Failed write message: %v", err)
	}
}

func (c *Client) getToken() string {
	var token string

	msg := getLoginAndPassword()
	data := msg

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8082/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&token)

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Error:", token)
		return ""
	}

	return token
}

func (c *Client) getChatToken() string {
	var token string

	data := c.token

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8082/getchattoken", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&token)

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Error:", token)
		return ""
	}

	return token
}

func (c *Client) SetToken() {
	c.token = c.getToken()
}

func getLoginAndPassword() string {
	var login, password string
	fmt.Println("Login: ")
	fmt.Scan(&login)
	fmt.Println("Password: ")
	fmt.Scan(&password)
	return login + " " + password
}

func (c *Client) SetChatToken() {
	c.chat_token = c.getChatToken()
}
