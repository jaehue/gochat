package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jaehue/gochat"
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

	var m gochat.Message
	if err := json.Unmarshal(buffer, &m); err != nil {
		log.Println("Fail to unmarshal message.", err)
	} else {
		log.Printf("[%s] %s", m.Sender, m.Text)
		messages <- buffer
	}

	handleConnection(conn)
}
