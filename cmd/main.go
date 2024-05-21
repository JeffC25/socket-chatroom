package main

import (
    "fmt"
    "io"
    "net"
    "strings"
    "sync"
    "time"
)

type message struct {
    Sender    string
    Timestamp time.Time
    Content   string
}

type client struct {
    Conn     net.Conn
    Name     string
    Messages chan message
}

var (
    broadcast = make(chan message)
    mutex     sync.Mutex
)

func (c *client) readMessages() {
    defer c.Conn.Close()

    // Send a prompt to the client to enter a name
    c.Conn.Write([]byte("Please enter your name: "))

    // Read messages without using bufio
    buf := make([]byte, 1024)
    for {
        n, err := c.Conn.Read(buf)
        if err == io.EOF {
            fmt.Println(c.Conn.RemoteAddr().String(), "disconnected")
            return
        } else if err != nil {
            fmt.Println("Error reading message:", err)
            return
        }

        msg := string(buf[:n])

        if c.Name == "" {
            c.Name = strings.TrimSpace(msg)
            fmt.Printf("New client connected with name: %s\n", c.Name)
            continue
        }

        mutex.Lock()
        broadcast <- message{
            Sender:    c.Name,
            Timestamp: time.Now(),
            Content:   msg,
        }
        mutex.Unlock()
    }
}

func (c *client) writeMessages() {
    defer c.Conn.Close()

    for msg := range c.Messages {
        _, err := c.Conn.Write([]byte(fmt.Sprintf("(%s %s): %s", msg.Sender, msg.Timestamp.Format(time.Stamp), msg.Content)))
        if err != nil {
            fmt.Println("Error writing message to client:", err)
            return
        }
    }
}

func main() {
    fmt.Println("Start server...")
    ln, err := net.Listen("tcp", ":8000")
    if err != nil {
        fmt.Println("Error binding to port:", err)
        return
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        clientMessages := make(chan message)
        client := &client{
            Conn:     conn,
            Messages: clientMessages,
        }

        go client.writeMessages()
        go client.readMessages()

        go func() {
            for {
                msg := <-broadcast
                clientMessages <- msg
            }
        }()
    }
}
