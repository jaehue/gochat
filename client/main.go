package main

import (
	"bufio"
	"fmt"
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

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if _, err := conn.Write([]byte(input)); err != nil {
			log.Println("Fail to write message.", err)
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}
}
