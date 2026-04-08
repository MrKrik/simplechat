package server

import (
	"bufio"
	"container/ring"
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Message struct {
	Author string
	Text   string
}

type Peer struct {
	Conn        net.Conn
	ConnectedAt time.Time
}

type Chat struct {
	messagesChan chan Message
}

type Server struct {
	Address       string
	Listener      net.Listener
	Chat_list     []string
	messagesChan  chan Message
	clients       map[net.Conn]*Peer
	deadClients   []net.Conn
	mu            sync.RWMutex
	messageBuffer *ring.Ring
}

func NewServer(address string) *Server {
	return &Server{
		messagesChan:  make(chan Message, 100),
		Address:       address,
		clients:       make(map[net.Conn]*Peer),
		deadClients:   make([]net.Conn, 0),
		messageBuffer: ring.New(10),
	}
}
func (s *Server) Start() error {
	var err error
	s.Listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		log.Printf("Some err %v", err)
		return err
	}
	log.Printf("Server started")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go s.acceptLoop(ctx)
	go s.Broadcast(ctx)
	defer func() {
		cancel()
		s.Listener.Close()
		close(s.messagesChan)
	}()
	s.Stop(ctx)
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	<-ctx.Done()
	for _, client := range s.clients {
		s.unregisterPeer(client.Conn)
	}
}

func (s *Server) acceptLoop(ctx context.Context) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("Accept loop stopped")
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}
		log.Printf("Welcome, %s", conn.RemoteAddr().String())
		go s.handleConn(conn)
	}
}

func (s *Server) registerPeer(conn net.Conn) {
	peer := &Peer{
		Conn:        conn,
		ConnectedAt: time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[conn] = peer
}

func (s *Server) addHistory(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messageBuffer.Value = msg
	s.messageBuffer = s.messageBuffer.Next()

}

func (s *Server) Broadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.messagesChan:
			go s.addHistory(msg)
			s.mu.RLock()
			message := fmt.Sprintf("%s: %s\n", msg.Author, msg.Text)
			for _, client := range s.clients {
				s.writeInConnection(client.Conn, message)
			}
			s.mu.RUnlock()

			for _, conn := range s.deadClients {
				s.unregisterPeer(conn)
			}
			s.deadClients = s.deadClients[:0]
		}
	}
}

func (s *Server) writeInConnection(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Failed write message: %v", err)
		s.deadClients = append(s.deadClients, conn)
	}
}

func (s *Server) unregisterPeer(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn.Close()
	delete(s.clients, conn)
	log.Printf("Client disconnected: %s", conn.RemoteAddr().String())
}

func (s *Server) sendHistory(conn net.Conn) {
	history := s.messageBuffer
	if history.Len() != 0 {
		history.Do(func(a any) {
			if a == nil {
				return
			}
			message := fmt.Sprintf("%s: %s\n", a.(Message).Author, a.(Message).Text)
			s.writeInConnection(conn, message)
		})
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	s.registerPeer(conn)
	reader := bufio.NewReader(conn)

	go s.sendHistory(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection error: %v", err)
			return
		}

		msg := &Message{
			Author: conn.RemoteAddr().String(),
			Text:   message,
		}
		s.messagesChan <- *msg
	}
}
