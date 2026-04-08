package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Message struct {
	Author string
	Text   string
}

type Client struct {
	Name        string
	Connection  net.Conn
	messageChan chan Message
}

func newClient(name string) *Client {
	return &Client{
		Name:        name,
		messageChan: make(chan Message, 50),
	}
}

func (c *Client) Start() error {
	var err error
	c.Connection, err = net.Dial("tcp", "chat-bee-bot.db-msk0.amvera.tech:27017")
	if err != nil {
		log.Println("Ошибка подключения:", err)
		return err
	}
	log.Println("Подключено к", c.Connection.RemoteAddr())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
		c.Connection.Close()
	}()
	c.Stop(ctx)
	return nil
}

func (c *Client) Stop(ctx context.Context) {
	<-ctx.Done()
	log.Println("Exit")
	c.Connection.Close()
}

func (c *Client) handleKeyboard(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		c.writeInConnection(text)
	}
}

func (c *Client) handleMessage(ctx context.Context) {
	reader := bufio.NewReader(c.Connection)
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
	_, err := c.Connection.Write([]byte(message))
	if err != nil {
		log.Printf("Failed write message: %v", err)
	}
}
func main() {
	cl := newClient("anton")
	cl.Start()
}
