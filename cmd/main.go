package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Start server...")

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error binding to port")
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection")
	}

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message")
		} else {
			fmt.Print("Message Received:", string(msg))

		}
	}
}
