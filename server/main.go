package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	clients  []net.Conn
	messages = make(chan []byte)
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Fail to listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Start to listening.")

	go runWriteLoop(messages)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Fail to connecting:", err.Error())
			return
		}

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		clients = append(clients, c)
		go handleConnection(c)
	}
}

func runWriteLoop(in <-chan []byte) {
	for message := range in {
		for _, c := range clients {
			if _, err := c.Write(message); err != nil {
				log.Println("Fail to write message.", err)
			}
		}
	}
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer[:len(buffer)-1]))

	messages <- buffer

	handleConnection(conn)
}
