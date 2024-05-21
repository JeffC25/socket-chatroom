package chatroom

import (
	"fmt"
	"net"
	"time"
)

type Message struct {
	Sender    *Client
	Timestamp time.Time
	Content   string
}

type ChatRoom struct {
	Broadcast   chan Message
	Clients     map[*Client]bool
	Connects    chan *Client
	Disconnects chan *Client
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Broadcast:   make(chan Message),
		Clients:     make(map[*Client]bool),
		Connects:    make(chan *Client),
		Disconnects: make(chan *Client),
	}
}

func (cr *ChatRoom) Run() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error binding to port")
	}

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection")
			}
			c := Client{
				Conn:    conn,
				Address: conn.RemoteAddr().String(),
				Name:    "",
				Room:    *cr,
			}
			cr.Connects <- &c
			fmt.Println("New connection:", c.Address)
		}
	}()

	for {
		select {
		case msg := <-cr.Broadcast:
			fmt.Printf("%s - %s (%s): %s\n", msg.Timestamp, msg.Sender.Name, msg.Sender.Address, msg.Content)
			for c := range cr.Clients {
				go c.ReceiveMessage(msg)
			}
		case c := <-cr.Connects:
			cr.Clients[c] = true
			c.GetName()
			fmt.Printf("%s (%s) %s connected to chatroom\n", time.Now().Local().Format(time.RFC3339), c.Name, c.Address)
			go c.SendMessage()
		case conn := <-cr.Disconnects:
			delete(cr.Clients, conn)
			fmt.Printf("%s (%s) %s disconnected from chatroom\n", time.Now().Local().Format(time.RFC3339), conn.Name, conn.Address)
		}
	}
}
