package chatroom

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Client struct {
	Conn    net.Conn
	Address string
	Name    string
	Room    ChatRoom
}

func (c *Client) SendMessage() {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(c.Conn.RemoteAddr().String(), "disconnected")
				return
			}
			fmt.Println("Error reading message:", err)
			return
		}

		msg = strings.TrimSpace(msg)

		c.Room.Broadcast <- Message{
			Sender:  c,
			Content: msg,
		}
	}
}

func (c *Client) ReceiveMessage(msg Message) {
	_, err := c.Conn.Write([]byte(fmt.Sprintf("%s - %s: %s\n", msg.Timestamp.Local().Format(time.Kitchen), msg.Sender.Name, msg.Content)))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func (c *Client) GetName() {
	_, err := c.Conn.Write([]byte("What is your name?\n"))
	if err != nil {
		fmt.Println("Error asking for name:", err)
		return
	}

	reader := bufio.NewReader(c.Conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}

	c.Name = strings.TrimSpace(name)
}
