package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Fail to listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Start to listening.")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Fail to connecting:", err.Error())
			return
		}

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
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

	if _, err := conn.Write(buffer); err != nil {
		log.Println("Fail to write messafe.", err)
	}

	handleConnection(conn)
}
