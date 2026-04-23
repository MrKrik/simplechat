package client

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
)

type Message struct {
	Author string
	Text   string
}

type Client struct {
	Connection  *websocket.Conn
	messageChan chan Message
	Token       string
}

func NewClient(name string) *Client {
	return &Client{
		messageChan: make(chan Message, 50),
	}
}

func (c *Client) Start() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c.Connection, _, err = websocket.Dial(ctx, "ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}

	if err != nil {
		log.Println("Ошибка подключения:", err)
		return err
	}

	if err != nil {
		fmt.Println("Authorization failed", err)
		return err
	}
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()
	go c.handleKeyboard(ctx)
	go c.handleMessage(ctx)
	defer func() {
		cancel()
		c.Connection.Close(websocket.StatusNormalClosure, "work end")
	}()
	c.Stop(ctx)
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
		c.writeInConnection(text)
	}
}

func (c *Client) handleMessage(ctx context.Context) {
	_, IOreader, _ := c.Connection.Reader(context.TODO())
	reader := bufio.NewReader(IOreader)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection error: %v", err)
			return
		}
		fmt.Println(message)
	}
}

func (c *Client) writeInConnection(message string) {
	err := c.Connection.Write(context.TODO(), websocket.MessageText, []byte(message))
	if err != nil {
		log.Printf("Failed write message: %v", err)
	}
}

func (c *Client) GetToken() string {
	var token string

	msg := getLoginAndPassword()
	data := msg

	// 2. Кодируем в JSON
	jsonData, _ := json.Marshal(data)

	// 3. Отправляем POST запрос
	resp, err := http.Post("http://localhost:8082/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // Обязательно закрываем тело ответа
	json.NewDecoder(resp.Body).Decode(&token)
	return token
}

func (c *Client) SetToken() {
	c.Token = c.GetToken()
}

func getLoginAndPassword() string {
	var login, password string
	fmt.Println("Login: ")
	fmt.Scan(&login)
	fmt.Println("Password: ")
	fmt.Scan(&password)
	return login + " " + password
}
