package hub

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Client struct {
	Connection    *websocket.Conn
	UserID        string
	token         string
	chat_token    string
	Room_id       string
	ServiceCancel context.CancelFunc
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
		Room_id: "general",
	}
}

func (c *Client) Start() error {
	fmt.Scan(&c.UserID)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := c.Connect(ctx)

	if err != nil {
		return err
	}

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

	var url string
	url = fmt.Sprintf("ws://localhost:8080/ws?token=%s&userID=%v", c.chat_token, c.UserID)

	dialCtx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()

	c.Connection, _, err = websocket.Dial(dialCtx, url, nil)
	if err != nil {
		fmt.Println("Authorization failed:", err)
		return err
	}

	log.Printf("Connected to server: %s", url)

	serviceCtx, serviceCancel := context.WithCancel(parentCtx)

	c.ServiceCancel = serviceCancel

	go c.handleKeyboard(serviceCtx)

	go c.handleMessage(serviceCtx)

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
		select {
		case <-ctx.Done():
			return
		default:
		}
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		if text == "/re" {
			c.Reconnect(ctx)
			continue
		}

		c.writeInConnection(text, ctx)
	}
}

func (c *Client) Reconnect(ctx context.Context) {
	c.Connection.CloseNow()
	err := c.Connect(ctx)
	if err != nil {
		log.Println(err.Error())
	}
}

func (c *Client) handleMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgType, message, err := c.Connection.Read(ctx)

		if err != nil {

			if errors.Is(err, io.EOF) {
				log.Println("f")
				c.Reconnect(ctx)
				return
			}

			if ctx.Err() != nil || websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			log.Printf("Connection error: %v", err)
			return
		}
		if msgType == websocket.MessageText {
			message = bytes.TrimSpace(message)

			if len(message) > 0 {
				fmt.Println(string(message))
			}
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
