package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeffc25/socket-chatroom/chatroom"
	"github.com/jeffc25/socket-chatroom/config"
)

func main() {
	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Starting server...")

	c, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config:", err)
		c.Host = "localhost"
		c.Port = "8000"
	}

	cr := chatroom.NewChatRoom(c)

	go cr.Run()

	<-notifyCh
	fmt.Println("Shutting down server...")
}
