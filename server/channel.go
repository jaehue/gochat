package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jaehue/gochat"
)

type Channel struct {
	Conn            net.Conn  `json:"-"`
	RemoteAddr      string    `json:"remoteAddr"`
	LastMessageTime time.Time `json:"lastMessageTime"`

	quit chan struct{}
}

func createChannel(conn net.Conn) *Channel {
	c := &Channel{
		Conn:            conn,
		RemoteAddr:      conn.RemoteAddr().String(),
		LastMessageTime: time.Now(),
		quit:            make(chan struct{}),
	}

	go c.handleConnection()
	go c.checkAlive()

	return c
}

func (c *Channel) handleConnection() {
	buffer, err := bufio.NewReader(c.Conn).ReadBytes('\n')
	if err != nil {
		log.Println(err)
		c.quit <- struct{}{}
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

	c.LastMessageTime = time.Now()

	c.handleConnection()
}

func (c *Channel) checkAlive() {
AliveLoop:
	for {

		tick := time.After(time.Second * 10)
		select {
		case <-tick:
			fmt.Println("check alive")
			if c.LastMessageTime.Add(time.Second * 10).Before(time.Now()) {
				fmt.Println("Gone")
				break AliveLoop
			}
		case <-c.quit:
			fmt.Println("Client left")
			break AliveLoop
		}
	}

	c.Conn.Close()
	for i := range channels {
		if channels[i] == c {
			channels = append(channels[:i], channels[i+1:]...)
			return
		}
	}
}
