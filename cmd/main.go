package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, c chan string) {
	defer conn.Close()

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message")
		} else {
			c <- ("(" + conn.RemoteAddr().String() + "): " + string(msg))
		}
	}
}

func main() {
	fmt.Println("Start server...")

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error binding to port")
	}

	defer ln.Close()

	c := make(chan string)

	go func() {
		for {
			fmt.Print(<- c)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
		}
		go handleConnection(conn, c)
	}
}
