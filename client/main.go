package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Success to connect to server.")

	go runReadLoop(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		if _, err := conn.Write([]byte(input)); err != nil {
			log.Println("Fail to write message.", err)
		}
	}
}

func runReadLoop(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Server closed")
				os.Exit(0)
			}
			log.Fatal("Fail to read message:", err.Error())
		}
		log.Print("Server relay:", message)
	}
}
