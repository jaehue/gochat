package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jaehue/gochat"
	"github.com/labstack/echo"
)

var (
	clients   []net.Conn
	messages  = make(chan []byte)
	histories []gochat.Message
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, histories)
	})
	go e.Start(":8001")

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
		histories = append(histories, m)
		messages <- buffer
	}

	handleConnection(conn)
}
