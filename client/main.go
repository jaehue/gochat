package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/jaehue/gochat"
)

var sender string

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("이름을 입력하세요: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Fail to read name.", err)
	}
	sender = input[:len(input)-1]

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Success to connect to server.")

	go runReadLoop(conn)

	for {
		input, _ := reader.ReadString('\n')

		m := gochat.Message{
			Sender:    sender,
			Text:      input[:len(input)-1],
			CreatedAt: time.Now(),
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Println("Fail to marshal message.", err)
			continue
		}

		if _, err := conn.Write(append(b, byte('\n'))); err != nil {
			log.Println("Fail to write message.", err)
		}
	}
}

func runReadLoop(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Server closed")
				os.Exit(0)
			}
			log.Fatal("Fail to read message:", err.Error())
		}

		var m gochat.Message
		if err := json.Unmarshal(message, &m); err != nil {
			log.Println("Fail to unmarshal message.", err)
			continue
		}

		log.Printf("[%s] %s\n", m.Sender, m.Text)
	}
}
