package main

import (
	"bufio"
	"fmt"
	"net"
)

type client struct {
	Conn	net.Conn
	Address	string
	Name	string
}

//type message struct {
//	Client	client
//	Time	time.Time
//	Content	string
//}

func (c client) handleConnection(m chan string) {
	defer c.Conn.Close()

	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message")
		} else {
			m <- ("(" + c.Conn.RemoteAddr().String() + "): " + string(msg))
		}
	}
}

func main() {
	fmt.Println("Start server...")
	
//	var clients []client

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error binding to port")
	}

	defer ln.Close()

	m := make(chan string)

	go func() {
		for {
			fmt.Print(<- m)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
		}
		
		var c = client{
			Conn: conn,
			Name: "",
		}
		go c.handleConnection(m)
	}
}
