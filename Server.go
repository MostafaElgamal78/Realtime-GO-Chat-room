//go:build server
// +build server

package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Client struct {
	ID   int
	Chan chan string
}

type ChatServer struct {
	mu      sync.Mutex
	clients map[int]*Client
	nextID  int
}

type JoinArgs struct{}
type JoinReply struct {
	ID int
}

type MessageArgs struct {
	ID   int
	Text string
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[int]*Client),
	}
}

func (s *ChatServer) Join(_ JoinArgs, reply *JoinReply) error {
	s.mu.Lock()
	id := s.nextID
	s.nextID++

	client := &Client{
		ID:   id,
		Chan: make(chan string, 10),
	}

	s.clients[id] = client
	s.mu.Unlock()

	*reply = JoinReply{ID: id}

	go s.broadcast(fmt.Sprintf("User %d joined", id), id)

	return nil
}

func (s *ChatServer) SendMessage(args MessageArgs, _ *struct{}) error {
	s.broadcast(fmt.Sprintf("User %d: %s", args.ID, args.Text), args.ID)
	return nil
}

func (s *ChatServer) GetUpdates(clientID int, reply *[]string) error {
	s.mu.Lock()
	client := s.clients[clientID]
	s.mu.Unlock()

	if client == nil {
		*reply = []string{"Client not found"}
		return nil
	}

	msgs := []string{}
	for {
		select {
		case m := <-client.Chan:
			msgs = append(msgs, m)
		default:
			*reply = msgs
			return nil
		}
	}
}

func (s *ChatServer) broadcast(msg string, senderID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, cl := range s.clients {
		if id != senderID {
			cl.Chan <- msg
		}
	}
}

func main() {
	server := NewChatServer()

	err := rpc.Register(server)
	if err != nil {
		log.Fatal("RPC register error:", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Server start error:", err)
	}

	fmt.Println("Chat server running on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
