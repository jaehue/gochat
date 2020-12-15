package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jaehue/gochat"
	"github.com/labstack/echo"
)

var (
	channels  []*Channel
	messages  = make(chan []byte)
	histories []gochat.Message
)

func main() {
	e := echo.New()
	e.GET("/histories", func(c echo.Context) error {
		return c.JSON(http.StatusOK, histories)
	})
	e.GET("/connections", func(c echo.Context) error {
		return c.JSON(http.StatusOK, channels)
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

		channels = append(channels, createChannel(c))
	}
}

func runWriteLoop(in <-chan []byte) {
	for message := range in {
		for _, c := range channels {
			if _, err := c.Conn.Write(message); err != nil {
				log.Println("Fail to write message.", err)
			}
		}
	}
}
