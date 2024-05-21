package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeffc25/socket-chatroom/chatroom"
)

func main() {
	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Starting server...")

	cr := chatroom.NewChatRoom()

	go cr.Run()

	<-notifyCh
	fmt.Println("Shutting down server...")
}
